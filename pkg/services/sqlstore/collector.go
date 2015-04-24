package sqlstore

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
	"strconv"
	"strings"
	"time"
)

func init() {
	bus.AddHandler("sql", GetCollectors)
	bus.AddHandler("sql", GetCollectorById)
	bus.AddHandler("sql", GetCollectorByName)
	bus.AddHandler("sql", AddCollector)
	bus.AddHandler("sql", UpdateCollector)
	bus.AddHandler("sql", DeleteCollector)
	bus.AddHandler("sql", GetCollectorHealthById)
	bus.AddHandler("sql", AddCollectorSession)
	bus.AddHandler("sql", DeleteCollectorSession)
	bus.AddHandler("sql", ClearCollectorSession)
	bus.AddHandler("sql", GetCollectorSessions)
}

type CollectorWithTag struct {
	Id        int64
	OrgId     int64
	Name      string
	Slug      string
	Tags      string
	Latitude  float64
	Longitude float64
	Public    bool
	Created   time.Time
	Updated   time.Time
	Online    bool
	Enabled   bool
}

func GetCollectorById(query *m.GetCollectorByIdQuery) error {
	sess := x.Table("collector")
	rawParams := make([]interface{}, 0)
	rawSql := `SELECT
		GROUP_CONCAT(DISTINCT(collector_tag.tag)) as tags,
		collector.*
	FROM collector
	LEFT JOIN collector_tag ON collector.id = collector_tag.collector_id AND collector_tag.org_id=?
	WHERE 
		(collector.public=1 OR collector.org_id=?)
	AND
		collector.id=?
	GROUP BY collector.id
	`
	rawParams = append(rawParams, query.OrgId, query.OrgId, query.Id)
	results := make([]CollectorWithTag, 0)
	err := sess.Sql(rawSql, rawParams...).Find(&results)
	if err != nil {
		return err
	}
	if len(results) < 1 {
		return m.ErrCollectorNotFound
	}

	result := results[0]

	tags := make([]string, 0)
	if result.Tags != "" {
		tags = strings.Split(result.Tags, ",")
	}

	query.Result = &m.CollectorDTO{
		Id:        result.Id,
		OrgId:     result.OrgId,
		Name:      result.Name,
		Slug:      result.Slug,
		Tags:      tags,
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
		Public:    result.Public,
		Online:    result.Online,
		Enabled:   result.Enabled,
	}

	return err
}

func GetCollectorByName(query *m.GetCollectorByNameQuery) error {
	sess := x.Table("collector")
	rawParams := make([]interface{}, 0)
	rawSql := `SELECT
		GROUP_CONCAT(DISTINCT(collector_tag.tag)) as tags,
		collector.*
	FROM collector
	LEFT JOIN collector_tag ON collector.id = collector_tag.collector_id AND collector_tag.org_id=?
	WHERE 
		(collector.public=1 OR collector.org_id=?)
	AND
		collector.name=?
	GROUP BY collector.id
	`
	rawParams = append(rawParams, query.OrgId, query.OrgId, query.Name)
	results := make([]CollectorWithTag, 0)
	err := sess.Sql(rawSql, rawParams...).Find(&results)
	if err != nil {
		return err
	}
	if len(results) < 1 {
		return m.ErrCollectorNotFound
	}

	result := results[0]

	tags := make([]string, 0)
	if result.Tags != "" {
		tags = strings.Split(result.Tags, ",")
	}

	query.Result = &m.CollectorDTO{
		Id:        result.Id,
		OrgId:     result.OrgId,
		Name:      result.Name,
		Slug:      result.Slug,
		Tags:      tags,
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
		Public:    result.Public,
		Online:    result.Online,
		Enabled:   result.Enabled,
	}

	return err
}

func GetCollectors(query *m.GetCollectorsQuery) error {
	sess := x.Table("collector")
	rawParams := make([]interface{}, 0)
	rawSql := `SELECT
		GROUP_CONCAT(DISTINCT(collector_tag.tag)) as tags,
		collector.*
	FROM collector
	LEFT JOIN collector_tag ON collector.id = collector_tag.collector_id AND collector_tag.org_id=?
	`
	rawParams = append(rawParams, query.OrgId)
	whereSql := make([]string, 0)
	whereSql = append(whereSql, "(collector.public=1 OR collector.org_id=?)")
	rawParams = append(rawParams, query.OrgId)
	if len(query.Tag) > 0 {
		// this is a bit complicated because we want to
		// match only collectors that have the tag(s),
		// but we still need to return all of the tags that
		// the collector has.
		rawSql += "LEFT JOIN collector_tag AS lt ON lt.collector_id = collector.id AND collector.org_id = collector_tag.org_id\n"
		p := make([]string, len(query.Tag))
		for i, t := range query.Tag {
			p[i] = "?"
			rawParams = append(rawParams, t)
		}
		whereSql = append(whereSql, fmt.Sprintf("lt.tag IN (%s)", strings.Join(p, ",")))
	}
	if len(query.Name) > 0 {
		p := make([]string, len(query.Name))
		for i, t := range query.Name {
			p[i] = "?"
			rawParams = append(rawParams, t)
		}
		whereSql = append(whereSql, fmt.Sprintf("collector.name IN (%s)", strings.Join(p, ",")))
	}
	if query.Public != "" {
		if p, err := strconv.ParseBool(query.Public); err == nil {
			whereSql = append(whereSql, "collector.public=?")
			rawParams = append(rawParams, p)
		} else {
			return err
		}
	}

	rawSql += "WHERE " + strings.Join(whereSql, " AND ")
	rawSql += " GROUP BY collector.id"

	result := make([]CollectorWithTag, 0)
	err := sess.Sql(rawSql, rawParams...).Find(&result)
	if err != nil {
		return err
	}

	collectors := make([]*m.CollectorDTO, len(result))

	//iterate through all of the results and build out our collectors model.
	for i, row := range result {
		tags := make([]string, 0)
		if row.Tags != "" {
			tags = strings.Split(row.Tags, ",")
		}
		collectors[i] = &m.CollectorDTO{
			Id:        row.Id,
			OrgId:     row.OrgId,
			Name:      row.Name,
			Slug:      row.Slug,
			Latitude:  row.Latitude,
			Longitude: row.Longitude,
			Tags:      tags,
			Public:    row.Public,
			Online:    row.Online,
			Enabled:   row.Enabled,
		}
	}

	query.Result = collectors
	return nil
}

func DeleteCollector(cmd *m.DeleteCollectorCommand) error {
	return inTransaction(func(sess *xorm.Session) error {
		//Query the collector to make sure we own it.
		collectorQuery := m.GetCollectorByIdQuery{
			Id:    cmd.Id,
			OrgId: cmd.OrgId,
		}
		err := GetCollectorById(&collectorQuery)
		if err != nil {
			return err
		}
		if collectorQuery.OrgId != cmd.OrgId {
			return errors.New("Permision Denined. You do not own this Collector.")
		}

		var rawSql = "DELETE FROM collector_tag WHERE collector_id=?"
		if _, err := sess.Exec(rawSql, cmd.Id); err != nil {
			return err
		}
		rawSql = "DELETE FROM collector WHERE id=?"
		if _, err := sess.Exec(rawSql, cmd.Id); err != nil {
			return err
		}
		return nil
	})
}

func AddCollector(cmd *m.AddCollectorCommand) error {

	return inTransaction(func(sess *xorm.Session) error {
		l := &m.Collector{
			OrgId:     cmd.OrgId,
			Name:      cmd.Name,
			Public:    cmd.Public,
			Latitude:  cmd.Latitude,
			Longitude: cmd.Longitude,
			Created:   time.Now(),
			Updated:   time.Now(),
			Online:    cmd.Online,
			Enabled:   cmd.Enabled,
		}
		l.UpdateCollectorSlug()
		sess.UseBool("public")
		sess.UseBool("online")
		sess.UseBool("enabled")
		if _, err := sess.Insert(l); err != nil {
			return err
		}
		collectorTags := make([]m.CollectorTag, 0, len(cmd.Tags))
		for _, tag := range cmd.Tags {
			collectorTags = append(collectorTags, m.CollectorTag{
				OrgId:       cmd.OrgId,
				CollectorId: l.Id,
				Tag:         tag,
			})
		}
		if len(collectorTags) > 0 {
			sess.Table("collector_tag")
			if _, err := sess.Insert(&collectorTags); err != nil {
				return err
			}
		}

		cmd.Result = &m.CollectorDTO{
			Id:        l.Id,
			OrgId:     l.OrgId,
			Name:      l.Name,
			Slug:      l.Slug,
			Tags:      cmd.Tags,
			Latitude:  l.Latitude,
			Longitude: l.Longitude,
			Public:    l.Public,
			Online:    l.Online,
			Enabled:   l.Enabled,
		}
		return nil
	})
}

func CopyPublicCollectorTags(orgId int64, sess *session) error {
	sess.Table("collector_tag")
	sess.Join("INNER", "collector", "collector.id=collector_tag.collector_id")
	sess.Where("collector.public=1").And("collector.org_id=collector_tag.org_id")
	result := make([]*m.CollectorTag, 0)
	err := sess.Find(&result)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		collectorTags := make([]m.CollectorTag, len(result))
		for i, collectorTag := range result {
			collectorTags[i] = m.CollectorTag{
				OrgId:       orgId,
				CollectorId: collectorTag.CollectorId,
				Tag:         collectorTag.Tag,
			}
		}
		sess.Table("collector_tag")
		if _, err := sess.Insert(&collectorTags); err != nil {
			return err
		}
	}
	return nil

}

func UpdateCollector(cmd *m.UpdateCollectorCommand) error {

	return inTransaction(func(sess *xorm.Session) error {
		//Query the collector to make sure we own it.
		collectorQuery := m.GetCollectorByIdQuery{
			Id:    cmd.Id,
			OrgId: cmd.OrgId,
		}
		err := GetCollectorById(&collectorQuery)
		if err != nil {
			return err
		}

		//the collector can only be edited by those who own it.
		if collectorQuery.Result.OrgId == cmd.OrgId {
			l := &m.Collector{
				OrgId:     cmd.OrgId,
				Latitude:  cmd.Latitude,
				Longitude: cmd.Longitude,
				Public:    cmd.Public,
				Updated:   time.Now(),
				Enabled:   cmd.Enabled,
			}

			sess.UseBool("public")
			_, err := sess.Id(cmd.Id).Update(l)
			if err != nil {
				return err
			}
		}

		rawSql := "DELETE FROM collector_tag WHERE collector_id=? and org_id=?"
		if _, err := sess.Exec(rawSql, cmd.Id, cmd.OrgId); err != nil {
			return err
		}

		collectorTags := make([]m.CollectorTag, 0, len(cmd.Tags))
		for _, tag := range cmd.Tags {
			collectorTags = append(collectorTags, m.CollectorTag{
				OrgId:       cmd.OrgId,
				CollectorId: cmd.Id,
				Tag:         tag,
			})
		}
		if len(collectorTags) > 0 {
			sess.Table("collector_tag")
			if _, err := sess.Insert(&collectorTags); err != nil {
				return err
			}
		}

		return nil
	})
}

func GetCollectorHealthById(query *m.GetCollectorHealthByIdQuery) error {
	sess := x.Table("monitor_collector_state")
	sess.Where("collector_id=?", query.Id).And("org_id=?", query.OrgId)
	err := sess.Find(&query.Result)
	if err != nil {
		return err
	}
	return nil
}

func AddCollectorSession(cmd *m.AddCollectorSessionCommand) error {
	return inTransaction(func(sess *xorm.Session) error {
		collectorSess := m.CollectorSession{
			OrgId:       cmd.OrgId,
			CollectorId: cmd.CollectorId,
			SocketId:    cmd.SocketId,
			Updated:     time.Now(),
		}
		if _, err := sess.Insert(&collectorSess); err != nil {
			return err
		}
		rawSql := "UPDATE collector set online=1 where id=?"
		if _, err := sess.Exec(rawSql, cmd.CollectorId); err != nil {
			return err
		}
		return nil
	})
}

func GetCollectorSessions(query *m.GetCollectorSessionsQuery) error {
	sess := x.Table("collector_session")
	fmt.Printf("searching for sessions for collector %d\n", query.CollectorId)
	err := sess.Where("collector_id=?", query.CollectorId).OrderBy("updated").Find(&query.Result)
	return err
}

func DeleteCollectorSession(cmd *m.DeleteCollectorSessionCommand) error {
	return inTransaction(func(sess *xorm.Session) error {
		var rawSql = "DELETE FROM collector_session WHERE org_id=? AND socket_id=?"
		if _, err := sess.Exec(rawSql, cmd.OrgId, cmd.SocketId); err != nil {
			return err
		}
		q := m.GetCollectorSessionsQuery{CollectorId: cmd.CollectorId}
		if err := GetCollectorSessions(&q); err != nil {
			return err
		}
		if len(q.Result) < 1 {
			rawSql := "UPDATE collector set online=0 where id=?"
			if _, err := sess.Exec(rawSql, cmd.CollectorId); err != nil {
				return err
			}
		}
		return nil
	})
}

func ClearCollectorSession(cmd *m.ClearCollectorSessionCommand) error {
	return inTransaction(func(sess *xorm.Session) error {
		var rawSql = "DELETE FROM collector_session"
		if _, err := sess.Exec(rawSql); err != nil {
			return err
		}
		rawSql = "UPDATE collector set online=0"
		if _, err := sess.Exec(rawSql); err != nil {
			return err
		}
		return nil
	})
}
