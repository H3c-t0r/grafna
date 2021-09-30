package pipeline

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/grafana/grafana/pkg/models"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/live"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	service     = "grafana"
	environment = "dev"
	id          = 1
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}

// TODO: some things to implement and ideas to consider.
// * Can return ChannelData from Frame Outputter – thus starting corresponding rule from processing input step.
// * Better error communication with WS/HTTP? Currently we can't return nothing beyond InternalError.
// * Better copy data and Frame on redirects? (at least for safety, but no real problem exists yet).
// * How to deal with auth on redirects? What if publish auth configured for a rule in redirect? Currently auth skipped.
// * Avoid ProcessorVars, OutputVars – use Vars everywhere.
// * Rename:
//		* Converter => FrameConverter
//		* Outputter => FrameOutputter
//		* Processor => FrameProcessor
//		* ConditionChecker => FrameConditionChecker.
// * Put SignedInUser into Vars.

// ChannelData is a wrapper over raw data with additional channel information.
// Channel is used for rule routing, if the channel is empty then data processing
// stops. If channel is not empty then data processing will be redirected to a
// corresponding channel rule.
type ChannelData struct {
	Channel string
	Data    []byte
}

// ChannelFrame is a wrapper over data.Frame with additional channel information.
// Channel is used for rule routing, if the channel is empty then frame processing
// will try to take current rule Processor and Outputter. If channel is not empty
// then frame processing will be redirected to a corresponding channel rule.
type ChannelFrame struct {
	Channel string
	Frame   *data.Frame
}

// Vars has some helpful things pipeline entities could use.
type Vars struct {
	OrgID     int64
	Channel   string
	Scope     string
	Namespace string
	Path      string
}

// ProcessorVars has some helpful things Processor entities could use.
type ProcessorVars struct {
	Vars
}

// OutputVars has some helpful things Outputter entities could use.
type OutputVars struct {
	ProcessorVars
}

// DataOutputter can output incoming data before conversion to frames.
type DataOutputter interface {
	Type() string
	OutputData(ctx context.Context, vars Vars, data []byte) ([]*ChannelData, error)
}

// Converter converts raw bytes to slice of ChannelFrame. Each element
// of resulting slice will be then individually processed and outputted
// according configured channel rules.
type Converter interface {
	Type() string
	Convert(ctx context.Context, vars Vars, body []byte) ([]*ChannelFrame, error)
}

// Processor can modify data.Frame in a custom way before it will be outputted.
type Processor interface {
	Type() string
	Process(ctx context.Context, vars ProcessorVars, frame *data.Frame) (*data.Frame, error)
}

// Outputter outputs data.Frame to a custom destination. Or simply
// do nothing if some conditions not met.
type Outputter interface {
	Type() string
	Output(ctx context.Context, vars OutputVars, frame *data.Frame) ([]*ChannelFrame, error)
}

// Subscriber can handle channel subscribe events.
type Subscriber interface {
	Type() string
	Subscribe(ctx context.Context, vars Vars) (models.SubscribeReply, backend.SubscribeStreamStatus, error)
}

// PublishAuthChecker checks whether current user can publish to a channel.
type PublishAuthChecker interface {
	CanPublish(ctx context.Context, u *models.SignedInUser) (bool, error)
}

// SubscribeAuthChecker checks whether current user can subscribe to a channel.
type SubscribeAuthChecker interface {
	CanSubscribe(ctx context.Context, u *models.SignedInUser) (bool, error)
}

// LiveChannelRule is an in-memory representation of each specific rule, with Converter, Processor
// and Outputter to be executed by Pipeline.
type LiveChannelRule struct {
	OrgId          int64
	Pattern        string
	PublishAuth    PublishAuthChecker
	SubscribeAuth  SubscribeAuthChecker
	DataOutputters []DataOutputter
	Converter      Converter
	Processors     []Processor
	Outputters     []Outputter
	Subscribers    []Subscriber
}

// Label ...
type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"` // Can be JSONPath or Goja script.
}

// Field description.
type Field struct {
	Name   string            `json:"name"`
	Type   data.FieldType    `json:"type"`
	Value  string            `json:"value"` // Can be JSONPath or Goja script.
	Labels []Label           `json:"labels,omitempty"`
	Config *data.FieldConfig `json:"config,omitempty"`
}

type ChannelRuleGetter interface {
	Get(orgID int64, channel string) (*LiveChannelRule, bool, error)
}

// Pipeline allows processing custom input data according to user-defined rules.
// This includes:
// * transforming custom input to data.Frame objects
// * do some processing on these frames
// * output resulting frames to various destinations.
type Pipeline struct {
	ruleGetter ChannelRuleGetter
	tracer     trace.Tracer
}

// New creates new Pipeline.
func New(ruleGetter ChannelRuleGetter) (*Pipeline, error) {
	logger.Info("Live pipeline initialization")
	p := &Pipeline{
		ruleGetter: ruleGetter,
	}

	if os.Getenv("GF_LIVE_PIPELINE_TRACE") != "" {
		// Traces for development only at the moment.
		// Start local Jaeger and then run Grafana with GF_LIVE_PIPELINE_TRACE:
		// docker run --rm -it --name jaeger -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 14250:14250 -p 9411:9411 jaegertracing/all-in-one:1.26
		// Then visit http://localhost:16686/ where Jaeger UI is served.
		tp, err := tracerProvider("http://localhost:14268/api/traces")
		if err != nil {
			return nil, err
		}
		tracer := tp.Tracer("gf.live.pipeline")
		p.tracer = tracer
	}

	if os.Getenv("GF_LIVE_PIPELINE_DEV") != "" {
		go postTestData() // TODO: temporary for development, remove before merge.
	}

	return p, nil
}

func (p *Pipeline) Get(orgID int64, channel string) (*LiveChannelRule, bool, error) {
	return p.ruleGetter.Get(orgID, channel)
}

func (p *Pipeline) ProcessInput(ctx context.Context, orgID int64, channelID string, body []byte) (bool, error) {
	var span trace.Span
	if p.tracer != nil {
		ctx, span = p.tracer.Start(ctx, "live.pipeline.process_input")
		span.SetAttributes(
			attribute.Int64("orgId", orgID),
			attribute.String("channel", channelID),
			attribute.String("body", string(body)),
		)
		defer span.End()
	}
	ok, err := p.processInput(ctx, orgID, channelID, body, nil)
	if err != nil {
		if p.tracer != nil && span != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		return ok, err
	}
	return ok, err
}

func (p *Pipeline) processInput(ctx context.Context, orgID int64, channelID string, body []byte, visitedChannels map[string]struct{}) (bool, error) {
	var span trace.Span
	if p.tracer != nil {
		ctx, span = p.tracer.Start(ctx, "live.pipeline.process_input_"+channelID)
		span.SetAttributes(
			attribute.Int64("orgId", orgID),
			attribute.String("channel", channelID),
			attribute.String("body", string(body)),
		)
		defer span.End()
	}
	rule, ok, err := p.ruleGetter.Get(orgID, channelID)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	if visitedChannels == nil {
		visitedChannels = map[string]struct{}{}
	}
	if len(rule.DataOutputters) > 0 {
		channelDataList := []*ChannelData{{Channel: channelID, Data: body}}
		err = p.processChannelDataList(ctx, orgID, channelID, channelDataList, visitedChannels)
		if err != nil {
			return false, err
		}
	}
	if rule.Converter == nil {
		return false, nil
	}
	channelFrames, err := p.dataToChannelFrames(ctx, *rule, orgID, channelID, body)
	if err != nil {
		return false, err
	}
	err = p.processChannelFrames(ctx, orgID, channelID, channelFrames, nil)
	if err != nil {
		return false, fmt.Errorf("error processing frame: %w", err)
	}
	return true, nil
}

func (p *Pipeline) dataToChannelFrames(ctx context.Context, rule LiveChannelRule, orgID int64, channelID string, body []byte) ([]*ChannelFrame, error) {
	var span trace.Span
	if p.tracer != nil {
		ctx, span = p.tracer.Start(ctx, "live.pipeline.convert_"+rule.Converter.Type())
		span.SetAttributes(
			attribute.Int64("orgId", orgID),
			attribute.String("channel", channelID),
		)
		defer span.End()
	}

	channel, err := live.ParseChannel(channelID)
	if err != nil {
		logger.Error("Error parsing channel", "error", err, "channel", channelID)
		return nil, err
	}

	vars := Vars{
		OrgID:     orgID,
		Channel:   channelID,
		Scope:     channel.Scope,
		Namespace: channel.Namespace,
		Path:      channel.Path,
	}

	frames, err := rule.Converter.Convert(ctx, vars, body)
	if err != nil {
		logger.Error("Error converting data", "error", err)
		return nil, err
	}

	return frames, nil
}

var errChannelRecursion = errors.New("channel recursion")

func (p *Pipeline) processChannelDataList(ctx context.Context, orgID int64, channelID string, channelDataList []*ChannelData, visitedChannels map[string]struct{}) error {
	for _, channelData := range channelDataList {
		var nextChannel = channelID
		if channelData.Channel != "" {
			nextChannel = channelData.Channel
		}
		if _, ok := visitedChannels[nextChannel]; ok {
			return fmt.Errorf("%w: %s", errChannelRecursion, nextChannel)
		}
		visitedChannels[nextChannel] = struct{}{}
		newChannelDataList, err := p.processData(ctx, orgID, nextChannel, channelData.Data)
		if err != nil {
			return err
		}
		if len(newChannelDataList) > 0 {
			for _, cd := range newChannelDataList {
				_, err := p.processInput(ctx, orgID, cd.Channel, cd.Data, visitedChannels)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *Pipeline) processChannelFrames(ctx context.Context, orgID int64, channelID string, channelFrames []*ChannelFrame, visitedChannels map[string]struct{}) error {
	if visitedChannels == nil {
		visitedChannels = map[string]struct{}{}
	}
	for _, channelFrame := range channelFrames {
		var processorChannel = channelID
		if channelFrame.Channel != "" {
			processorChannel = channelFrame.Channel
		}
		if _, ok := visitedChannels[processorChannel]; ok {
			return fmt.Errorf("%w: %s", errChannelRecursion, processorChannel)
		}
		visitedChannels[processorChannel] = struct{}{}
		frames, err := p.processFrame(ctx, orgID, processorChannel, channelFrame.Frame)
		if err != nil {
			return err
		}
		if len(frames) > 0 {
			err := p.processChannelFrames(ctx, orgID, processorChannel, frames, visitedChannels)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Pipeline) processFrame(ctx context.Context, orgID int64, channelID string, frame *data.Frame) ([]*ChannelFrame, error) {
	var span trace.Span
	if p.tracer != nil {
		table, err := frame.StringTable(32, 32)
		if err != nil {
			return nil, err
		}
		ctx, span = p.tracer.Start(ctx, "live.pipeline.process_frame_"+channelID)
		span.SetAttributes(
			attribute.Int64("orgId", orgID),
			attribute.String("channel", channelID),
			attribute.String("frame", table),
		)
		defer span.End()
	}
	rule, ruleOk, err := p.ruleGetter.Get(orgID, channelID)
	if err != nil {
		logger.Error("Error getting rule", "error", err)
		return nil, err
	}
	if !ruleOk {
		logger.Debug("Rule not found", "channel", channelID)
		return nil, err
	}

	ch, err := live.ParseChannel(channelID)
	if err != nil {
		logger.Error("Error parsing channel", "error", err, "channel", channelID)
		return nil, err
	}

	vars := ProcessorVars{
		Vars: Vars{
			OrgID:     orgID,
			Channel:   channelID,
			Scope:     ch.Scope,
			Namespace: ch.Namespace,
			Path:      ch.Path,
		},
	}

	if len(rule.Processors) > 0 {
		for _, proc := range rule.Processors {
			frame, err = p.execProcessor(ctx, proc, vars, frame)
			if err != nil {
				logger.Error("Error processing frame", "error", err)
				return nil, err
			}
			if frame == nil {
				return nil, nil
			}
		}
	}

	outputVars := OutputVars{
		ProcessorVars: vars,
	}

	if len(rule.Outputters) > 0 {
		var resultingFrames []*ChannelFrame
		for _, out := range rule.Outputters {
			frames, err := p.processFrameOutput(ctx, out, outputVars, frame)
			if err != nil {
				logger.Error("Error outputting frame", "error", err)
				return nil, err
			}
			resultingFrames = append(resultingFrames, frames...)
		}
		return resultingFrames, nil
	}

	return nil, nil
}

func (p *Pipeline) execProcessor(ctx context.Context, proc Processor, vars ProcessorVars, frame *data.Frame) (*data.Frame, error) {
	var span trace.Span
	if p.tracer != nil {
		ctx, span = p.tracer.Start(ctx, "live.pipeline.apply_processor_"+proc.Type())
		table, err := frame.StringTable(32, 32)
		if err != nil {
			return nil, err
		}
		span.SetAttributes(
			attribute.Int64("orgId", vars.OrgID),
			attribute.String("channel", vars.Channel),
			attribute.String("frame", table),
			attribute.String("processor", proc.Type()),
		)
		// Note: we can also visualize resulting frame here.
		defer span.End()
	}
	return proc.Process(ctx, vars, frame)
}

func (p *Pipeline) processFrameOutput(ctx context.Context, out Outputter, vars OutputVars, frame *data.Frame) ([]*ChannelFrame, error) {
	var span trace.Span
	if p.tracer != nil {
		ctx, span = p.tracer.Start(ctx, "live.pipeline.frame_output_"+out.Type())
		table, err := frame.StringTable(32, 32)
		if err != nil {
			return nil, err
		}
		span.SetAttributes(
			attribute.Int64("orgId", vars.OrgID),
			attribute.String("channel", vars.Channel),
			attribute.String("frame", table),
			attribute.String("output", out.Type()),
		)
		defer span.End()
	}
	return out.Output(ctx, vars, frame)
}

func (p *Pipeline) processData(ctx context.Context, orgID int64, channelID string, data []byte) ([]*ChannelData, error) {
	var span trace.Span
	if p.tracer != nil {
		ctx, span = p.tracer.Start(ctx, "live.pipeline.process_data_"+channelID)
		span.SetAttributes(
			attribute.Int64("orgId", orgID),
			attribute.String("channel", channelID),
			attribute.String("data", string(data)),
		)
		defer span.End()
	}
	rule, ruleOk, err := p.ruleGetter.Get(orgID, channelID)
	if err != nil {
		logger.Error("Error getting rule", "error", err)
		return nil, err
	}
	if !ruleOk {
		logger.Debug("Rule not found", "channel", channelID)
		return nil, err
	}

	ch, err := live.ParseChannel(channelID)
	if err != nil {
		logger.Error("Error parsing channel", "error", err, "channel", channelID)
		return nil, err
	}

	vars := Vars{
		OrgID:     orgID,
		Channel:   channelID,
		Scope:     ch.Scope,
		Namespace: ch.Namespace,
		Path:      ch.Path,
	}

	if len(rule.DataOutputters) > 0 {
		var resultingChannelDataList []*ChannelData
		for _, out := range rule.DataOutputters {
			channelDataList, err := p.processDataOutput(ctx, out, vars, data)
			if err != nil {
				logger.Error("Error outputting frame", "error", err)
				return nil, err
			}
			resultingChannelDataList = append(resultingChannelDataList, channelDataList...)
		}
		return resultingChannelDataList, nil
	}

	return nil, nil
}

func (p *Pipeline) processDataOutput(ctx context.Context, out DataOutputter, vars Vars, data []byte) ([]*ChannelData, error) {
	var span trace.Span
	if p.tracer != nil {
		ctx, span = p.tracer.Start(ctx, "live.pipeline.data_output_"+out.Type())
		span.SetAttributes(
			attribute.Int64("orgId", vars.OrgID),
			attribute.String("channel", vars.Channel),
			attribute.String("data", string(data)),
			attribute.String("output", out.Type()),
		)
		defer span.End()
	}
	return out.OutputData(ctx, vars, data)
}
