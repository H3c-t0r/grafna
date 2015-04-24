package api

import (
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/middleware"
	m "github.com/grafana/grafana/pkg/models"
	_ "github.com/grafana/grafana/pkg/services/endpointdiscovery"
)

func GetEndpointById(c *middleware.Context) {
	id := c.ParamsInt64(":id")

	query := m.GetEndpointByIdQuery{Id: id, OrgId: c.OrgId}

	err := bus.Dispatch(&query)
	if err != nil {
		c.JsonApiErr(404, "Endpoint not found", nil)
		return
	}

	c.JSON(200, query.Result)
}

func getEndpointHealthById(c *middleware.Context) {
	id := c.ParamsInt64(":id")
	query := m.GetEndpointHealthByIdQuery{
		Id:    id,
		OrgId: c.OrgId,
	}
	err := bus.Dispatch(&query)
	if err != nil {
		c.JsonApiErr(500, "Failed to query endpoint health", err)
		return
	}

	c.JSON(200, query.Result)
}

func GetEndpoints(c *middleware.Context, query m.GetEndpointsQuery) {
	query.OrgId = c.OrgId

	if err := bus.Dispatch(&query); err != nil {
		c.JsonApiErr(500, "Failed to query endpoints", err)
		return
	}
	c.JSON(200, query.Result)
}

func DeleteEndpoint(c *middleware.Context) {
	id := c.ParamsInt64(":id")

	cmd := &m.DeleteEndpointCommand{Id: id, OrgId: c.OrgId}

	err := bus.Dispatch(cmd)
	if err != nil {
		c.JsonApiErr(500, "Failed to delete endpoint", err)
		return
	}

	c.JsonOK("endpoint deleted")
}

func AddEndpoint(c *middleware.Context, cmd m.AddEndpointCommand) {
	cmd.OrgId = c.OrgId
	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Failed to add endpoint", err)
		return
	}

	c.JSON(200, cmd.Result)
}

func UpdateEndpoint(c *middleware.Context, cmd m.UpdateEndpointCommand) {
	cmd.OrgId = c.OrgId

	err := bus.Dispatch(&cmd)
	if err != nil {
		c.JsonApiErr(500, "Failed to update endpoint", err)
		return
	}

	c.JsonOK("Endpoint updated")
}

func DiscoverEndpoint(c *middleware.Context, cmd m.EndpointDiscoveryCommand) {
	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Failed to discover endpoint", err)
		return
	}
	c.JSON(200, cmd.Result)
}
