package sqlstore

import (
	"strconv"

	"github.com/go-xorm/xorm"

	"github.com/torkelo/grafana-pro/pkg/bus"
	m "github.com/torkelo/grafana-pro/pkg/models"
)

func init() {
	bus.AddHandler("sql", StarDashboard)
	bus.AddHandler("sql", UnstarDashboard)
	bus.AddHandler("sql", GetUserStars)
	bus.AddHandler("sql", IsStarredByUser)
}

func IsStarredByUser(query *m.IsStarredByUserQuery) error {
	rawSql := "SELECT 1 from star where user_id=? and dashboard_id=?"
	results, err := x.Query(rawSql, query.UserId, query.DashboardId)

	if err != nil {
		return err
	}
	if len(results) == 0 {
		return nil
	}

	query.Result, _ = strconv.ParseBool(string(results[0]["1"]))

	return nil
}

func StarDashboard(cmd *m.StarDashboardCommand) error {
	if cmd.DashboardId == 0 || cmd.UserId == 0 {
		return m.ErrCommandValidationFailed
	}

	return inTransaction(func(sess *xorm.Session) error {

		entity := m.Star{
			UserId:      cmd.UserId,
			DashboardId: cmd.DashboardId,
		}

		_, err := sess.Insert(&entity)
		return err
	})
}

func UnstarDashboard(cmd *m.UnstarDashboardCommand) error {
	if cmd.DashboardId == 0 || cmd.UserId == 0 {
		return m.ErrCommandValidationFailed
	}

	return inTransaction(func(sess *xorm.Session) error {
		var rawSql = "DELETE FROM star WHERE user_id=? and dashboard_id=?"
		_, err := sess.Exec(rawSql, cmd.UserId, cmd.DashboardId)
		return err
	})
}

func GetUserStars(query *m.GetUserStarsQuery) error {
	query.Result = make([]m.Star, 0)
	err := x.Where("user_id=?", query.UserId).Find(&query.Result)
	return err
}
