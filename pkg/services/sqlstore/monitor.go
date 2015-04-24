package sqlstore

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/events"
	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
)

func init() {
	bus.AddHandler("sql", GetMonitors)
	bus.AddHandler("sql", GetMonitorById)
	bus.AddHandler("sql", GetMonitorTypes)
	bus.AddHandler("sql", AddMonitor)
	bus.AddHandler("sql", UpdateMonitor)
	bus.AddHandler("sql", DeleteMonitor)
	bus.AddHandler("sql", UpdateMonitorCollectorState)
	bus.AddHandler("sql", GetMonitorHealthById)
}

type MonitorWithCollectorDTO struct {
	Id            int64
	EndpointId    int64
	OrgId         int64
	Namespace     string
	MonitorTypeId int64
	CollectorIds  string
	CollectorTags string
	TagCollectors string
	State         int64
	StateChange   time.Time
	Settings      []*m.MonitorSettingDTO
	Frequency     int64
	Enabled       bool
	Offset        int64
	Updated       time.Time
	Created       time.Time
}

func GetMonitorById(query *m.GetMonitorByIdQuery) error {
	sess := x.Table("monitor")
	rawParams := make([]interface{}, 0)
	rawSql := `
SELECT
    GROUP_CONCAT(DISTINCT(monitor_collector.collector_id)) as collector_ids,
    GROUP_CONCAT(DISTINCT(monitor_collector_tag.tag)) as collector_tags,
    GROUP_CONCAT(DISTINCT(collector_tags.collector_id)) as tag_collectors,
    monitor.*
FROM monitor
    LEFT JOIN monitor_collector ON monitor.id = monitor_collector.monitor_id
    LEFT JOIN monitor_collector_tag ON monitor.id = monitor_collector_tag.monitor_id
    LEFT JOIN 
        (SELECT
            collector.id AS collector_id,
            collector_tag.tag as tag
        FROM collector
        LEFT JOIN collector_tag ON collector.id = collector_tag.collector_id
        WHERE (collector.public=1 OR collector.org_id=?) AND (collector_tag.org_id=? OR collector_tag.id is NULL)) as collector_tags
    ON collector_tags.tag = monitor_collector_tag.tag
WHERE monitor.id=?
	`
	rawParams = append(rawParams, query.OrgId, query.OrgId, query.Id)

	if !query.IsGrafanaAdmin {
		rawSql += "AND monitor.org_id=?\n"
		rawParams = append(rawParams, query.OrgId)
	}
	rawSql += "GROUP BY monitor.id"

	//store the results into an array of maps.
	results := make([]*MonitorWithCollectorDTO, 0)
	err := sess.Sql(rawSql, rawParams...).Find(&results)
	if err != nil {
		return err
	}
	if len(results) < 1 {
		return m.ErrMonitorNotFound
	}
	result := results[0]

	monitorCollectorIds := make([]int64, 0)
	monitorCollectorsMap := make(map[int64]bool)
	if result.CollectorIds != "" {
		for _, l := range strings.Split(result.CollectorIds, ",") {
			i, err := strconv.ParseInt(l, 10, 64)
			if err != nil {
				return err
			}
			monitorCollectorIds = append(monitorCollectorIds, i)
			monitorCollectorsMap[i] = true
		}
	}

	monitorCollectorTags := make([]string, 0)
	if result.CollectorTags != "" {
		monitorCollectorTags = strings.Split(result.CollectorTags, ",")
		for _, l := range strings.Split(result.TagCollectors, ",") {
			i, err := strconv.ParseInt(l, 10, 64)
			if err != nil {
				return err
			}
			monitorCollectorsMap[i] = true
		}
	}

	mergedCollectors := make([]int64, len(monitorCollectorsMap))
	count := 0
	for k := range monitorCollectorsMap {
		mergedCollectors[count] = k
		count += 1
	}

	query.Result = &m.MonitorDTO{
		Id:            result.Id,
		EndpointId:    result.EndpointId,
		OrgId:         result.OrgId,
		Namespace:     result.Namespace,
		MonitorTypeId: result.MonitorTypeId,
		CollectorIds:  monitorCollectorIds,
		CollectorTags: monitorCollectorTags,
		Collectors:    mergedCollectors,
		State:         result.State,
		StateChange:   result.StateChange,
		Settings:      result.Settings,
		Frequency:     result.Frequency,
		Enabled:       result.Enabled,
		Offset:        result.Offset,
		Updated:       result.Updated,
	}

	return nil
}

func GetMonitors(query *m.GetMonitorsQuery) error {
	sess := x.Table("monitor")
	rawParams := make([]interface{}, 0)
	rawSql := `
SELECT
    GROUP_CONCAT(DISTINCT(monitor_collector.collector_id)) as collector_ids,
    GROUP_CONCAT(DISTINCT(monitor_collector_tag.tag)) as collector_tags,
    GROUP_CONCAT(DISTINCT(collector_tags.collector_id)) as tag_collectors,
    monitor.*
FROM monitor
    LEFT JOIN monitor_collector ON monitor.id = monitor_collector.monitor_id
    LEFT JOIN monitor_collector_tag ON monitor.id = monitor_collector_tag.monitor_id
    LEFT JOIN 
        (SELECT
            collector.id AS collector_id,
            collector_tag.tag as tag
        FROM collector
        LEFT JOIN collector_tag ON collector.id = collector_tag.collector_id
        WHERE (collector.public=1 OR collector.org_id=?) AND (collector_tag.org_id=? OR collector_tag.id is NULL)) as collector_tags
    ON collector_tags.tag = monitor_collector_tag.tag
`
	rawParams = append(rawParams, query.OrgId, query.OrgId)
	whereSql := make([]string, 0)
	if !query.IsGrafanaAdmin {
		whereSql = append(whereSql, "monitor.org_id=?")
		rawParams = append(rawParams, query.OrgId)
	}

	if len(query.EndpointId) > 0 {
		p := make([]string, len(query.EndpointId))
		for i, e := range query.EndpointId {
			p[i] = "?"
			rawParams = append(rawParams, e)
		}
		whereSql = append(whereSql, fmt.Sprintf("monitor.endpoint_id IN (%s)", strings.Join(p, ",")))
	}

	if len(query.CollectorId) > 0 {
		rawSql += "LEFT JOIN monitor_collector AS mc ON mc.monitor_id = monitor.id\n"
		rawSql += `LEFT JOIN 
        (SELECT
            collector.id AS collector_id,
            collector_tag.tag as tag
        FROM collector
        LEFT JOIN collector_tag ON collector.id = collector_tag.collector_id
        WHERE (collector.public=1 OR collector.org_id=?) AND (collector_tag.org_id=? OR collector_tag.id is NULL)) as ct
		ON ct.tag = monitor_collector_tag.tag
		`
		rawParams = append(rawParams, query.OrgId, query.OrgId)

		p := make([]string, len(query.CollectorId))
		for i, c := range query.CollectorId {
			p[i] = "?"
			rawParams = append(rawParams, c)
		}

		p2 := make([]string, len(query.CollectorId))
		for i, e := range query.CollectorId {
			p2[i] = "?"
			rawParams = append(rawParams, e)
		}
		whereSql = append(whereSql, fmt.Sprintf("(mc.collector_id IN (%s) OR ct.collector_id IN (%s))", strings.Join(p, ","), strings.Join(p2, ",")))
	}

	if query.Modulo > 0 {
		whereSql = append(whereSql, "(monitor.id % ?) = ?")
		rawParams = append(rawParams, query.Modulo, query.ModuloOffset)
	}

	if len(whereSql) > 0 {
		rawSql += "WHERE " + strings.Join(whereSql, " AND ")
	}
	rawSql += " GROUP BY monitor.id"

	result := make([]*MonitorWithCollectorDTO, 0)
	err := sess.Sql(rawSql, rawParams...).Find(&result)
	if err != nil {
		return err
	}

	monitors := make([]*m.MonitorDTO, 0)
	//iterate through all of the results and build out our checks model.
	for _, row := range result {
		monitorCollectorIds := make([]int64, 0)
		monitorCollectorsMap := make(map[int64]bool)
		if row.CollectorIds != "" {
			for _, l := range strings.Split(row.CollectorIds, ",") {
				i, err := strconv.ParseInt(l, 10, 64)
				if err != nil {
					return err
				}
				monitorCollectorIds = append(monitorCollectorIds, i)
				monitorCollectorsMap[i] = true
			}
		}

		monitorCollectorTags := make([]string, 0)
		if row.CollectorTags != "" {
			monitorCollectorTags = strings.Split(row.CollectorTags, ",")
			for _, l := range strings.Split(row.TagCollectors, ",") {
				i, err := strconv.ParseInt(l, 10, 64)
				if err != nil {
					return err
				}
				monitorCollectorsMap[i] = true
			}
		}

		mergedCollectors := make([]int64, len(monitorCollectorsMap))
		count := 0
		for k := range monitorCollectorsMap {
			mergedCollectors[count] = k
			count += 1
		}

		monitors = append(monitors, &m.MonitorDTO{
			Id:            row.Id,
			EndpointId:    row.EndpointId,
			OrgId:         row.OrgId,
			Namespace:     row.Namespace,
			MonitorTypeId: row.MonitorTypeId,
			CollectorIds:  monitorCollectorIds,
			CollectorTags: monitorCollectorTags,
			Collectors:    mergedCollectors,
			State:         row.State,
			StateChange:   row.StateChange,
			Settings:      row.Settings,
			Frequency:     row.Frequency,
			Enabled:       row.Enabled,
			Offset:        row.Offset,
			Updated:       row.Updated,
		})
	}
	query.Result = monitors

	return nil

}

type MonitorTypeWithSettingsDTO struct {
	Id           int64
	Name         string
	Variable     string
	Description  string
	Required     bool
	DataType     string
	Conditions   map[string]interface{}
	DefaultValue string
}

func GetMonitorTypes(query *m.GetMonitorTypesQuery) error {
	sess := x.Table("monitor_type")
	sess.Limit(100, 0).Asc("name")
	sess.Join("LEFT", "monitor_type_setting", "monitor_type_setting.monitor_type_id=monitor_type.id")
	sess.Cols("monitor_type.id", "monitor_type.name",
		"monitor_type_setting.variable", "monitor_type_setting.description",
		"monitor_type_setting.required", "monitor_type_setting.data_type",
		"monitor_type_setting.conditions", "monitor_type_setting.default_value")

	result := make([]*MonitorTypeWithSettingsDTO, 0)
	err := sess.Find(&result)
	if err != nil {
		return err
	}
	types := make(map[int64]*m.MonitorTypeDTO)
	//iterate through all of the results and build out our checks model.
	for _, row := range result {
		if _, ok := types[row.Id]; ok != true {
			//this is the first time we have seen this monitorId
			var typeSettings []*m.MonitorTypeSettingDTO
			types[row.Id] = &m.MonitorTypeDTO{
				Id:       row.Id,
				Name:     row.Name,
				Settings: typeSettings,
			}
		}

		types[row.Id].Settings = append(types[row.Id].Settings, &m.MonitorTypeSettingDTO{
			Variable:     row.Variable,
			Description:  row.Description,
			Required:     row.Required,
			DataType:     row.DataType,
			Conditions:   row.Conditions,
			DefaultValue: row.DefaultValue,
		})
	}

	query.Result = make([]*m.MonitorTypeDTO, len(types))
	count := 0
	for _, v := range types {
		query.Result[count] = v
		count++
	}

	return nil
}

func DeleteMonitor(cmd *m.DeleteMonitorCommand) error {
	return inTransaction2(func(sess *session) error {
		q := m.GetMonitorByIdQuery{
			Id:    cmd.Id,
			OrgId: cmd.OrgId,
		}
		err := GetMonitorById(&q)
		if err != nil {
			return err
		}
		var rawSql = "DELETE FROM monitor WHERE id=? and org_id=?"
		_, err = sess.Exec(rawSql, cmd.Id, cmd.OrgId)
		if err != nil {
			return err
		}
		rawSql = "DELETE FROM monitor_collector WHERE monitor_id=?"
		_, err = sess.Exec(rawSql, cmd.Id)
		if err != nil {
			return err
		}
		rawSql = "DELETE FROM monitor_collector_tag WHERE monitor_id=?"
		_, err = sess.Exec(rawSql, cmd.Id)
		if err != nil {
			return err
		}
		rawSql = "DELETE FROM monitor_collector_state WHERE monitor_id=?"
		_, err = sess.Exec(rawSql, cmd.Id)
		if err != nil {
			return err
		}

		sess.publishAfterCommit(&events.MonitorRemoved{
			Timestamp:     time.Now(),
			Id:            q.Result.Id,
			EndpointId:    q.Result.EndpointId,
			OrgId:         q.Result.OrgId,
			CollectorIds:  q.Result.CollectorIds,
			CollectorTags: q.Result.CollectorTags,
			Collectors:    q.Result.Collectors,
		})
		return nil
	})
}

// store collector list query result
type collectorList struct {
	Id int64
}

func AddMonitor(cmd *m.AddMonitorCommand) error {
	return inTransaction2(func(sess *session) error {
		return addMonitorTransaction(cmd, sess)
	})
}

func addMonitorTransaction(cmd *m.AddMonitorCommand, sess *session) error {
	//validate Endpoint.
	endpointQuery := m.GetEndpointByIdQuery{
		Id:    cmd.EndpointId,
		OrgId: cmd.OrgId,
	}
	err := GetEndpointByIdTransaction(&endpointQuery, sess)
	if err != nil {
		return err
	}

	filtered_collectors := make([]*collectorList, 0, len(cmd.CollectorIds))
	if len(cmd.CollectorIds) > 0 {
		sess.Table("collector")
		sess.In("id", cmd.CollectorIds).Where("org_id=? or public=1", cmd.OrgId)
		sess.Cols("id")
		err = sess.Find(&filtered_collectors)

		if err != nil {
			return err
		}
	}

	if len(filtered_collectors) < len(cmd.CollectorIds) {
		return m.ErrMonitorCollectorsInvalid
	}

	//get settings definition for thie monitorType.
	var typeSettings []*m.MonitorTypeSetting
	sess.Table("monitor_type_setting")
	sess.Where("monitor_type_id=?", cmd.MonitorTypeId)
	err = sess.Find(&typeSettings)
	if err != nil {
		return nil
	}

	// push the typeSettings into a Map with the variable name as key
	settingMap := make(map[string]*m.MonitorTypeSetting)
	for _, s := range typeSettings {
		settingMap[s.Variable] = s
	}

	//validate the settings.
	seenMetrics := make(map[string]bool)
	for _, v := range cmd.Settings {
		def, ok := settingMap[v.Variable]
		if ok != true {
			log.Info("Unkown variable %s passed.", v.Variable)
			return m.ErrMonitorSettingsInvalid
		}
		//TODO:(awoods) make sure the value meets the definition.
		seenMetrics[def.Variable] = true
		log.Info("%s present in settings", def.Variable)
	}

	//make sure all required variables were provided.
	//add defaults for missing optional variables.
	for k, s := range settingMap {
		if _, ok := seenMetrics[k]; ok != true {
			log.Info("%s not in settings", k)
			if s.Required {
				// required setting variable missing.
				return m.ErrMonitorSettingsInvalid
			}
			cmd.Settings = append(cmd.Settings, &m.MonitorSettingDTO{
				Variable: k,
				Value:    s.DefaultValue,
			})
		}
	}

	if cmd.Namespace == "" {
		label := strings.ToLower(endpointQuery.Result.Name)
		re := regexp.MustCompile("[^\\w-]+")
		re2 := regexp.MustCompile("\\s")
		slug := re2.ReplaceAllString(re.ReplaceAllString(label, "_"), "-")
		cmd.Namespace = slug
	}

	mon := &m.Monitor{
		EndpointId:    cmd.EndpointId,
		OrgId:         cmd.OrgId,
		Namespace:     cmd.Namespace,
		MonitorTypeId: cmd.MonitorTypeId,
		Offset:        rand.Int63n(cmd.Frequency - 1),
		Settings:      cmd.Settings,
		Created:       time.Now(),
		Updated:       time.Now(),
		Frequency:     cmd.Frequency,
		Enabled:       cmd.Enabled,
		State:         -1,
		StateChange:   time.Now(),
	}

	if _, err := sess.Insert(mon); err != nil {
		return err
	}

	if len(cmd.CollectorIds) > 0 {
		monitor_collectors := make([]*m.MonitorCollector, len(cmd.CollectorIds))
		for i, l := range cmd.CollectorIds {
			monitor_collectors[i] = &m.MonitorCollector{
				MonitorId:   mon.Id,
				CollectorId: l,
			}
		}
		sess.Table("monitor_collector")
		if _, err := sess.Insert(&monitor_collectors); err != nil {
			return err
		}
	}

	if len(cmd.CollectorTags) > 0 {
		monitor_collector_tags := make([]*m.MonitorCollectorTag, len(cmd.CollectorTags))
		for i, t := range cmd.CollectorTags {
			monitor_collector_tags[i] = &m.MonitorCollectorTag{
				MonitorId: mon.Id,
				Tag:       t,
			}
		}

		sess.Table("monitor_collector_tag")
		if _, err := sess.Insert(&monitor_collector_tags); err != nil {
			return err
		}
	}
	// get collectorIds from tags
	tagCollectors, err := getCollectorIdsFromTags(cmd.OrgId, cmd.CollectorTags, sess)
	if err != nil {
		return err
	}

	collectorIdMap := make(map[int64]bool)
	collectorList := make([]int64, 0)
	for _, id := range cmd.CollectorIds {
		collectorIdMap[id] = true
		collectorList = append(collectorList, id)
	}

	for _, id := range tagCollectors {
		if _, ok := collectorIdMap[id]; !ok {
			collectorList = append(collectorList, id)
		}
	}

	if len(collectorList) > 0 {
		monitor_collector_states := make([]*m.MonitorCollectorState, len(collectorList))
		for i, c := range collectorList {
			monitor_collector_states[i] = &m.MonitorCollectorState{
				OrgId:       mon.OrgId,
				EndpointId:  mon.EndpointId,
				MonitorId:   mon.Id,
				CollectorId: c,
				State:       -1,
				Updated:     time.Now(),
			}
		}
		sess.Table("monitor_collector_state")
		if _, err := sess.Insert(&monitor_collector_states); err != nil {
			return err
		}
	}

	cmd.Result = &m.MonitorDTO{
		Id:            mon.Id,
		EndpointId:    mon.EndpointId,
		OrgId:         mon.OrgId,
		Namespace:     mon.Namespace,
		MonitorTypeId: mon.MonitorTypeId,
		CollectorIds:  cmd.CollectorIds,
		CollectorTags: cmd.CollectorTags,
		Collectors:    collectorList,
		Settings:      mon.Settings,
		Frequency:     mon.Frequency,
		Enabled:       mon.Enabled,
		State:         mon.State,
		StateChange:   mon.StateChange,
		Offset:        mon.Offset,
		Updated:       mon.Updated,
	}
	sess.publishAfterCommit(&events.MonitorCreated{
		Timestamp: mon.Updated,
		MonitorPayload: events.MonitorPayload{
			Id:            mon.Id,
			EndpointId:    mon.EndpointId,
			OrgId:         mon.OrgId,
			Namespace:     mon.Namespace,
			MonitorTypeId: mon.MonitorTypeId,
			CollectorIds:  cmd.CollectorIds,
			CollectorTags: cmd.CollectorTags,
			Collectors:    collectorList,
			Settings:      mon.Settings,
			Frequency:     mon.Frequency,
			Enabled:       mon.Enabled,
			Offset:        mon.Offset,
			Updated:       mon.Updated,
		},
	})
	return nil
}

func UpdateMonitor(cmd *m.UpdateMonitorCommand) error {
	return inTransaction2(func(sess *session) error {
		//validate Endpoint.
		endpointQuery := m.GetEndpointByIdQuery{
			Id:    cmd.EndpointId,
			OrgId: cmd.OrgId,
		}
		err := GetEndpointById(&endpointQuery)
		if err != nil {
			return err
		}
		currentEndpoint := endpointQuery.Result

		q := m.GetMonitorByIdQuery{
			Id:    cmd.Id,
			OrgId: cmd.OrgId,
		}
		err = GetMonitorById(&q)
		if err != nil {
			return err
		}
		lastState := q.Result

		if lastState.EndpointId != cmd.EndpointId {
			return m.ErrorEndpointCantBeChanged
		}

		//validate collectors.
		filtered_collectors := make([]*collectorList, 0, len(cmd.CollectorIds))
		if len(cmd.CollectorIds) > 0 {
			sess.Table("collector")
			sess.In("id", cmd.CollectorIds).Where("org_id=? or public=1", cmd.OrgId)
			sess.Cols("id")
			err = sess.Find(&filtered_collectors)

			if err != nil {
				return err
			}
		}

		if len(filtered_collectors) < len(cmd.CollectorIds) {
			return m.ErrMonitorCollectorsInvalid
		}

		//get settings definition for thie monitorType.
		var typeSettings []*m.MonitorTypeSetting
		sess.Table("monitor_type_setting")
		sess.Where("monitor_type_id=?", cmd.MonitorTypeId)
		err = sess.Find(&typeSettings)
		if err != nil {
			return nil
		}
		if len(typeSettings) < 1 {
			log.Info("no monitorType settings found for type: %d", cmd.MonitorTypeId)
			return m.ErrMonitorSettingsInvalid
		}

		// push the typeSettings into a Map with the variable name as key
		settingMap := make(map[string]*m.MonitorTypeSetting)
		for _, s := range typeSettings {
			settingMap[s.Variable] = s
		}

		//validate the settings.
		seenMetrics := make(map[string]bool)
		for _, v := range cmd.Settings {
			def, ok := settingMap[v.Variable]
			if ok != true {
				log.Info("Unkown variable %s passed.", v.Variable)
				return m.ErrMonitorSettingsInvalid
			}
			//TODO:(awoods) make sure the value meets the definition.
			seenMetrics[def.Variable] = true
		}

		//make sure all required variables were provided.
		//add defaults for missing optional variables.
		for k, s := range settingMap {
			if _, ok := seenMetrics[k]; ok != true {
				log.Info("%s not in settings", k)
				if s.Required {
					// required setting variable missing.
					return m.ErrMonitorSettingsInvalid
				}
				cmd.Settings = append(cmd.Settings, &m.MonitorSettingDTO{
					Variable: k,
					Value:    s.DefaultValue,
				})
			}
		}

		if cmd.Namespace == "" {
			label := strings.ToLower(currentEndpoint.Name)
			re := regexp.MustCompile("[^\\w-]+")
			re2 := regexp.MustCompile("\\s")
			slug := re2.ReplaceAllString(re.ReplaceAllString(label, "_"), "-")
			cmd.Namespace = slug
		}

		mon := &m.Monitor{
			Id:            cmd.Id,
			EndpointId:    cmd.EndpointId,
			OrgId:         cmd.OrgId,
			Namespace:     cmd.Namespace,
			MonitorTypeId: cmd.MonitorTypeId,
			Settings:      cmd.Settings,
			Updated:       time.Now(),
			Enabled:       cmd.Enabled,
			State:         lastState.State,
			StateChange:   lastState.StateChange,
			Frequency:     cmd.Frequency,
		}

		//check if we need to update the time offset for when the monitor should run.
		if mon.Offset >= mon.Frequency {
			mon.Offset = rand.Int63n(mon.Frequency - 1)
		}

		var rawSql = "DELETE FROM monitor_collector WHERE monitor_id=?"
		if _, err := sess.Exec(rawSql, cmd.Id); err != nil {
			return err
		}
		if len(cmd.CollectorIds) > 0 {
			monitor_collectors := make([]*m.MonitorCollector, len(cmd.CollectorIds))
			for i, l := range cmd.CollectorIds {
				monitor_collectors[i] = &m.MonitorCollector{
					MonitorId:   cmd.Id,
					CollectorId: l,
				}
			}
			sess.Table("monitor_collector")
			if _, err := sess.Insert(&monitor_collectors); err != nil {
				return err
			}
		}

		rawSql = "DELETE FROM monitor_collector_tag WHERE monitor_id=?"
		if _, err := sess.Exec(rawSql, cmd.Id); err != nil {
			return err
		}
		if len(cmd.CollectorTags) > 0 {
			monitor_collector_tags := make([]*m.MonitorCollectorTag, len(cmd.CollectorTags))
			for i, t := range cmd.CollectorTags {
				monitor_collector_tags[i] = &m.MonitorCollectorTag{
					MonitorId: cmd.Id,
					Tag:       t,
				}
			}

			sess.Table("monitor_collector_tag")
			if _, err := sess.Insert(&monitor_collector_tags); err != nil {
				return err
			}
		}

		// get collectorIds from tags
		tagCollectors, err := getCollectorIdsFromTags(cmd.OrgId, cmd.CollectorTags, sess)
		if err != nil {
			return err
		}

		collectorIdMap := make(map[int64]bool)
		collectorList := make([]int64, 0)
		lastCollectors := make(map[int64]bool)
		for _, id := range lastState.Collectors {
			lastCollectors[id] = false
		}

		for _, id := range cmd.CollectorIds {
			collectorIdMap[id] = true
			collectorList = append(collectorList, id)
		}

		for _, id := range tagCollectors {
			if _, ok := collectorIdMap[id]; !ok {
				collectorList = append(collectorList, id)
			}
		}

		stateChange, err := updateCollectorState(mon, collectorList, sess)
		if err != nil {
			return err
		}
		if stateChange {
			fmt.Println("TODO: propagate stateChange to endpoint.")
		}

		sess.Table("monitor")
		sess.UseBool("enabled")
		if _, err = sess.Where("id=? and org_id=?", mon.Id, mon.OrgId).Update(mon); err != nil {
			return err
		}

		sess.publishAfterCommit(&events.MonitorUpdated{
			MonitorPayload: events.MonitorPayload{
				Id:            mon.Id,
				EndpointId:    mon.EndpointId,
				OrgId:         mon.OrgId,
				Namespace:     mon.Namespace,
				MonitorTypeId: mon.MonitorTypeId,
				CollectorIds:  cmd.CollectorIds,
				CollectorTags: cmd.CollectorTags,
				Collectors:    collectorList,
				Settings:      mon.Settings,
				Frequency:     mon.Frequency,
				Enabled:       mon.Enabled,
				Offset:        mon.Offset,
				Updated:       mon.Updated,
			},
			Timestamp: mon.Updated,
			LastState: &events.MonitorPayload{
				Id:            lastState.Id,
				EndpointId:    lastState.EndpointId,
				OrgId:         lastState.OrgId,
				Namespace:     lastState.Namespace,
				MonitorTypeId: lastState.MonitorTypeId,
				CollectorIds:  lastState.CollectorIds,
				CollectorTags: lastState.CollectorTags,
				Collectors:    lastState.Collectors,
				Settings:      lastState.Settings,
				Frequency:     lastState.Frequency,
				Enabled:       lastState.Enabled,
				Offset:        lastState.Offset,
				Updated:       lastState.Updated,
			},
		})

		return err
	})
}

type CollectorId struct {
	CollectorId int64
}

func getCollectorIdsFromTags(orgId int64, tags []string, sess *session) ([]int64, error) {
	result := make([]int64, 0)
	if len(tags) < 1 {
		return result, nil
	}
	params := make([]interface{}, 0)
	rawSql := `SELECT DISTINCT(collector.id) AS collector_id
	FROM collector
	INNER JOIN collector_tag ON collector.id = collector_tag.collector_id 
	WHERE (collector.public=1 OR collector.org_id=?) 
		AND collector_tag.org_id=?
	`

	params = append(params, orgId, orgId)

	p := make([]string, len(tags))
	for i, t := range tags {
		p[i] = "?"
		params = append(params, t)
	}
	rawSql += fmt.Sprintf("AND collector_tag.tag IN (%s)", strings.Join(p, ","))

	results := make([]CollectorId, 0)
	if err := sess.Sql(rawSql, params...).Find(&results); err != nil {
		return result, err
	}

	if len(results) > 0 {
		for _, r := range results {
			result = append(result, r.CollectorId)
		}
	}

	return result, nil
}

func updateCollectorState(mon *m.Monitor, collectorList []int64, sess *session) (bool, error) {
	states := make([]*m.MonitorCollectorState, 0)
	sess.Table("monitor_collector_state")
	err := sess.Where("monitor_id=?", mon.Id).Find(&states)
	if err != nil {
		return false, err
	}

	collectorMap := make(map[int64]bool)
	for _, id := range collectorList {
		collectorMap[id] = false
	}

	collectorsToDel := make([]int64, 0)

	state := int64(-1)
	okCount, unknownCount, warnCount, errorCount := 0, 0, 0, 0
	totalCount := len(collectorList)

	if len(states) > 0 {
		for _, row := range states {
			if _, ok := collectorMap[row.CollectorId]; !ok {
				collectorsToDel = append(collectorsToDel, row.CollectorId)
				continue
			}
			collectorMap[row.CollectorId] = true

			switch row.State {
			case 0:
				okCount++
			case 1:
				warnCount++
			case 2:
				errorCount++
			default:
				unknownCount++
			}
		}
	}

	for _, seen := range collectorMap {
		if !seen {
			unknownCount++
		}
	}

	if len(collectorsToDel) > 0 {
		params := make([]interface{}, len(collectorsToDel)+1)
		params[0] = mon.Id
		p := make([]string, len(collectorsToDel))
		for i, c := range collectorsToDel {
			p[i] = "?"
			params[i+1] = c
		}
		rawSql := fmt.Sprintf("DELETE FROM monitor_collector_state WHERE monitor_id=? and collector_id IN (%s)", strings.Join(p, ","))
		if _, err := sess.Exec(rawSql, params...); err != nil {
			return false, err
		}
	}

	// determine our state
	if okCount < totalCount/2 || (totalCount-okCount) >= 3 {
		//state is not OK.
		if errorCount > 0 {
			state = 2
		} else if warnCount > 0 {
			state = 1
		} else {
			state = -1
		}
	} else {
		state = 0
	}
	stateChange := false
	if mon.State != state {
		mon.State = state
		mon.StateChange = time.Now()
		mon.Updated = time.Now()
		stateChange = true
	}
	return stateChange, nil
}

func UpdateMonitorCollectorState(cmd *m.UpdateMonitorCollectorStateCommand) error {
	return inTransaction2(func(sess *session) error {
		sess.Table("monitor_collector_state")
		results := make([]*m.MonitorCollectorState, 0)
		sess.Where("monitor_id=?", cmd.MonitorId)
		sess.And("org_id=?", cmd.OrgId).And("endpoint_id=?", cmd.EndpointId).And("collector_id=?", cmd.CollectorId)
		if err := sess.Find(&results); err != nil {
			return err
		}
		stateChange := false
		if len(results) < 1 {
			//need to insert
			state := &m.MonitorCollectorState{
				OrgId:       cmd.OrgId,
				MonitorId:   cmd.MonitorId,
				EndpointId:  cmd.EndpointId,
				CollectorId: cmd.CollectorId,
				State:       cmd.State,
				Updated:     cmd.Updated,
			}
			sess.UseBool("state")
			if _, err := sess.Insert(state); err != nil {
				return err
			}
			stateChange = true
		} else if results[0].State != cmd.State {
			//need to update
			state := results[0]
			state.State = cmd.State
			state.Updated = cmd.Updated
			sess.UseBool("state")
			if _, err := sess.Id(state.Id).Update(state); err != nil {
				return err
			}
			stateChange = true
		}

		if stateChange {
			//update state of monitor.
			q := m.GetMonitorByIdQuery{
				Id:    cmd.MonitorId,
				OrgId: cmd.OrgId,
			}
			err := GetMonitorById(&q)
			if err != nil {
				return err
			}
			monView := q.Result
			mon := &m.Monitor{
				Id:            monView.Id,
				EndpointId:    monView.EndpointId,
				OrgId:         monView.OrgId,
				Namespace:     monView.Namespace,
				MonitorTypeId: monView.MonitorTypeId,
				Offset:        monView.Offset,
				Settings:      monView.Settings,
				Updated:       monView.StateChange,
				Frequency:     monView.Frequency,
				Enabled:       monView.Enabled,
				State:         monView.State,
				StateChange:   monView.StateChange,
			}
			monStateChange, err := updateCollectorState(mon, monView.Collectors, sess)
			if err != nil {
				return err
			}
			if monStateChange {
				sess.Table("monitor")
				sess.UseBool("state")
				if _, err := sess.Id(mon.Id).Update(mon); err != nil {
					return err
				}
				fmt.Println("TODO: propagate stateChange to endpoint.")
			}
		}
		return nil
	})
}

func GetMonitorHealthById(query *m.GetMonitorHealthByIdQuery) error {
	sess := x.Table("monitor_collector_state")
	sess.Where("monitor_id=?", query.Id).And("org_id=?", query.OrgId)
	//query.result = make([]*m.MonitorCollectorState, 0)
	err := sess.Find(&query.Result)
	if err != nil {
		return err
	}
	return nil
}
