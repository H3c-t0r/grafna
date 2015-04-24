package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/components/apikeygen"
	"github.com/grafana/grafana/pkg/events"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/middleware"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/eventpublisher"
	"github.com/grafana/grafana/pkg/services/metricpublisher"
	"reflect"
	"runtime"
	"time"
)

var server *socketio.Server
var bufCh chan m.MetricDefinition
var localSockets map[string]socketio.Socket

type CollectorContext struct {
	*m.SignedInUser
	Collector *m.CollectorDTO
	Socket    socketio.Socket
	SocketId  string
}

type BinMessage struct {
	Payload *socketio.Attachment
}

func register(so socketio.Socket) (*CollectorContext, error) {
	req := so.Request()
	req.ParseForm()
	keyString := req.Form.Get("apiKey")
	name := req.Form.Get("name")
	if name == "" {
		return nil, errors.New("collector name not provided.")
	}
	if keyString != "" {
		// base64 decode key
		decoded, err := apikeygen.Decode(keyString)
		if err != nil {
			return nil, m.ErrInvalidApiKey
		}
		// fetch key
		keyQuery := m.GetApiKeyByNameQuery{KeyName: decoded.Name, OrgId: decoded.OrgId}
		if err := bus.Dispatch(&keyQuery); err != nil {
			return nil, m.ErrInvalidApiKey
		} else {
			apikey := keyQuery.Result

			// validate api key
			if !apikeygen.IsValid(decoded, apikey.Key) {
				return nil, m.ErrInvalidApiKey
			}
			// lookup collector
			colQuery := m.GetCollectorByNameQuery{Name: name, OrgId: apikey.OrgId}
			if err := bus.Dispatch(&colQuery); err != nil {
				return nil, m.ErrCollectorNotFound
			} else {

				sess := &CollectorContext{
					SignedInUser: &m.SignedInUser{
						IsGrafanaAdmin: apikey.IsAdmin,
						OrgRole:        apikey.Role,
						ApiKeyId:       apikey.Id,
						OrgId:          apikey.OrgId,
						Name:           apikey.Name,
					},
					Collector: colQuery.Result,
					Socket:    so,
					SocketId:  so.Id(),
				}
				if err := sess.Save(); err != nil {
					return nil, err
				}
				return sess, nil
			}
		}
	}
	return nil, m.ErrInvalidApiKey
}

func InitCollectorController() {
	cmd := m.ClearCollectorSessionCommand{ProcessId: 0}
	if err := bus.Dispatch(&cmd); err != nil {
		log.Fatal(4, "failed to clear collectorSessions.", err)
	}
	bufCh = make(chan m.MetricDefinition, runtime.NumCPU())
	go metricpublisher.ProcessBuffer(bufCh)
}

func init() {
	localSockets = make(map[string]socketio.Socket)
	var err error
	server, err = socketio.NewServer([]string{"polling", "websocket"})
	if err != nil {
		log.Fatal(4, "failed to initialize socketio.", err)
		return
	}
	server.On("connection", func(so socketio.Socket) {
		sess, err := register(so)
		if err != nil {
			log.Error(0, "Failed to initialize collector.", err)
			so.Emit("authFailed", err.Error())
			return
		}
		sess.OnConnection()
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Error(0, "socket emitted error", err)
	})
	//tap into the update/add/Delete events emitted when monitors are modified.
	bus.AddEventListener(EmitUpdateMonitor)
	bus.AddEventListener(EmitAddMonitor)
	bus.AddEventListener(EmitDeleteMonitor)
}

func (c *CollectorContext) Save() error {
	cmd := &m.AddCollectorSessionCommand{
		CollectorId: c.Collector.Id,
		SocketId:    c.Socket.Id(),
		OrgId:       c.OrgId,
	}
	if err := bus.Dispatch(cmd); err != nil {
		return err
	}
	localSockets[c.SocketId] = c.Socket
	return nil
}

func (c *CollectorContext) Remove() error {
	log.Info(fmt.Sprintf("removing socket with Id %s", c.SocketId))
	cmd := &m.DeleteCollectorSessionCommand{
		SocketId:    c.SocketId,
		OrgId:       c.OrgId,
		CollectorId: c.Collector.Id,
	}
	err := bus.Dispatch(cmd)
	delete(localSockets, c.SocketId)
	return err
}

func (c *CollectorContext) OnConnection() {
	log.Info(fmt.Sprintf("New connection for %s owned by OrgId: %d", c.Collector.Name, c.OrgId))
	c.Socket.Emit("ready", c.Collector)
	c.Socket.On("event", c.OnEvent)
	c.Socket.On("results", c.OnResults)
	c.Socket.On("disconnection", c.OnDisconnection)
	RefreshCollector(c.Collector.Id)
}

func (c *CollectorContext) OnDisconnection() {
	log.Info(fmt.Sprintf("%s disconnected", c.Collector.Name))
	if err := c.Remove(); err != nil {
		log.Error(4, fmt.Sprintf("Failed to remove collectorSession. %s", c.Collector.Name), err)
	}
}

func (c *CollectorContext) OnEvent(msg *m.EventDefinition) {
	log.Info(fmt.Sprintf("recieved event from %s", c.Collector.Name))
	if !c.IsGrafanaAdmin {
		msg.OrgId = c.OrgId
	}
	msgString, err := json.Marshal(msg)
	if err != nil {
		log.Error(0, "Failed to marshal event.", err)
	}
	// send to RabbitMQ
	routingKey := fmt.Sprintf("EVENT.%s.%s", msg.Severity, msg.EventType)
	go eventpublisher.Publish(routingKey, msgString)
	if msg.EventType == "monitor_state" {
		updateState(msg)
	}
}

func (c *CollectorContext) OnResults(results []*m.MetricDefinition) {
	for _, r := range results {
		if !c.IsGrafanaAdmin {
			r.OrgId = c.OrgId
		}
		bufCh <- *r
	}
}

func RefreshCollector(collectorId int64) {
	log.Info(fmt.Sprintf("Collector %d refreshing", collectorId))
	//step 1. get list of collectorSessions for this collector.
	q := m.GetCollectorSessionsQuery{CollectorId: collectorId}
	if err := bus.Dispatch(&q); err != nil {
		log.Error(0, "failed to get list of collectorSessions.", err)
		return
	}
	totalSessions := len(q.Result)
	//step 2. for each session
	for pos, sess := range q.Result {
		//step 3. get list of monitors configured for this colletor.
		monQuery := m.GetMonitorsQuery{
			CollectorId:    []int64{collectorId},
			IsGrafanaAdmin: true,
			Modulo:         int64(totalSessions),
			ModuloOffset:   int64(pos),
		}
		if err := bus.Dispatch(&monQuery); err != nil {
			log.Error(0, "failed to get list of monitors.", err)
			return
		}
		log.Info("sending refresh to " + sess.SocketId)
		//step 5. get socket
		socket, ok := localSockets[sess.SocketId]
		if !ok {
			log.Error(0, "socket"+sess.SocketId+" is not local.", nil)
		} else {
			//step 6. send
			socket.Emit("refresh", monQuery.Result)
		}
	}
}

func SocketIO(c *middleware.Context) {
	if server == nil {
		log.Fatal(4, "socket.io server not initialized.", nil)
	}

	server.ServeHTTP(c.Resp, c.Req.Request)
}

func EmitUpdateMonitor(event *events.MonitorUpdated) error {
	log.Info("processing monitorUpdated event.")
	seenCollectors := make(map[int64]bool)
	for _, collectorId := range event.Collectors {
		seenCollectors[collectorId] = true
		if err := EmitEvent(collectorId, "updated", event); err != nil {
			return err
		}
	}
	for _, collectorId := range event.LastState.Collectors {
		if _, ok := seenCollectors[collectorId]; !ok {
			if err := EmitEvent(collectorId, "removed", event); err != nil {
				return err
			}
		}
	}
	return nil
}

func EmitAddMonitor(event *events.MonitorCreated) error {
	log.Info("processing monitorCreated event")
	for _, collectorId := range event.Collectors {
		if err := EmitEvent(collectorId, "created", event); err != nil {
			return err
		}
	}
	return nil
}

func EmitDeleteMonitor(event *events.MonitorRemoved) error {
	log.Info("processing monitorRemoved event")
	for _, collectorId := range event.Collectors {
		if err := EmitEvent(collectorId, "removed", event); err != nil {
			return err
		}
	}
	return nil
}

func EmitEvent(collectorId int64, eventName string, event interface{}) error {
	q := m.GetCollectorSessionsQuery{CollectorId: collectorId}
	if err := bus.Dispatch(&q); err != nil {
		return err
	}
	totalSessions := int64(len(q.Result))
	if totalSessions < 1 {
		return nil
	}
	eventId := reflect.ValueOf(event).Elem().FieldByName("Id").Int()
	log.Info(fmt.Sprintf("emitting %s event for MonitorId %d totalSessions: %d", eventName, eventId, totalSessions))
	pos := eventId % totalSessions
	socketId := q.Result[pos].SocketId
	socket, ok := localSockets[socketId]
	if !ok {
		log.Info("socket" + socketId + " is not local.")
	} else {
		//step 6. send
		socket.Emit(eventName, event)
	}
	return nil
}

func updateState(event *m.EventDefinition) {
	collector, ok := event.Extra["collector_id"]
	if !ok {
		log.Error(0, "event does not have collector_id set.", nil)
		return
	}
	endpoint, ok := event.Extra["endpoint_id"]
	if !ok {
		log.Error(0, "event does not have endpoint_id set.", nil)
		return
	}
	monitor, ok := event.Extra["monitor_id"]
	if !ok {
		log.Error(0, "event does not have monitor_id set.", nil)
		return
	}
	log.Debug(fmt.Sprintf("updating state of monitor: %v from %v", monitor, event.Extra["collector"]))
	cmd := m.UpdateMonitorCollectorStateCommand{
		OrgId:       event.OrgId,
		EndpointId:  int64(endpoint.(float64)),
		MonitorId:   int64(monitor.(float64)),
		CollectorId: int64(collector.(float64)),
		Updated:     time.Unix(0, event.Timestamp*int64(time.Millisecond)),
	}
	// update the check state
	switch event.Severity {
	case "OK":
		cmd.State = 0
	case "WARN":
		cmd.State = 1
	case "ERROR":
		cmd.State = 2
	default:
		cmd.State = -1
	}

	if err := bus.Dispatch(&cmd); err != nil {
		log.Error(0, "faile to update MonitorcollectorState", err)
	}
}
