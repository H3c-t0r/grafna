+++
title = "Enterprise HTTP APIs"
description = "Grafana Enterprise APIs"
keywords = ["grafana", "enterprise", "api"]
aliases = ["/docs/grafana/latest/enterprise/http_api/"]
type = "docs"
[menu.docs]
name = "HTTP APIs"
parent = "enterprise"
weight = 1000
+++

# Enterprise HTTP APIs

## Authentication

To access the API you need to authenticate yourself. For more information about API authentication, refer to [Authentication API]({{< relref "../http_api/auth.md" >}}).

## Reporting API

### Send report

> Only available in Grafana Enterprise v7.0+.

`POST /api/reports/email`

Generate and send a report. This API waits for the report to be generated before returning. You are recommended to set the client's timeout to 60 seconds or more.

This API endpoint is experimental and may be deprecated in a future release. On deprecation, a migration strategy will be provided and the endpoint will remain functional until the next major release of Grafana.

**Example request:**

```http
POST /api/reports/email HTTP/1.1
Accept: application/json
Content-Type: application/json
Authorization: Bearer eyJrIjoiT0tTcG1pUlY2RnVKZTFVaDFsNFZXdE9ZWmNrMkZYbk

{
  "id":"3",
  "useEmailsFromReport": true
}
```

JSON Body Schema:

field name | type | description
---------- | ---- | -----------
id | string | ID of the report to send (same as in the URL when editing a report, not to be confused with the ID of the dashboard). Required.
emails | string | Comma-separated list of emails to which to send the report to. Overrides the emails from the report. Required if **useEmailsFromReport** is not present.
useEmailsFromReport | boolean | Send the report to the emails specified in the report. Required if **emails** is not present.

**Example response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 29

{"message":"Report was sent"}
```

Status Codes:

Code | Description
---- | -----------
200 | Report was sent
400 | Bad request (invalid json, missing content-type, missing or invalid fields, etc.)
401 | Authentication failed, refer to [Authentication API]({{< relref "../http_api/auth.md" >}})
403 | User is authenticated but is not authorized to generate the report
404 | Report not found
500 | Unexpected error or server misconfiguration. Refer to body and/or server logs for more details
