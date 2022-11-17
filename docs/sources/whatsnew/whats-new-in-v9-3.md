---
_build:
  list: false
aliases:
  - /docs/grafana/latest/guides/whats-new-in-v9-3/
description: Feature and improvement highlights for Grafana v9.3
keywords:
  - grafana
  - new
  - documentation
  - '9.3'
  - release notes
title: What's new in Grafana v9.3
weight: -33
---

# What’s new in Grafana v9.3 (Beta)

Welcome to Grafana v9.3. If you’d prefer to dig into the details, check out the complete [changelog](https://github.com/grafana/grafana/blob/master/CHANGELOG.md).

## What's New feature template

Use this template to add your what's new section.

[Generally available | Available in experimental/beta] in Grafana [Open Source, Enterprise, Cloud Free, Cloud Pro, Cloud Advanced]

Description. Include an overview of the feature and problem it solves, and where to learn more (like a link to the docs).

Add a screenshot to this enablement material folder and link to it here.

## Public Dashboards - Audit Table

Available in experimental in Grafana Open Source, Enterprise, Cloud Pro, Cloud Advanced

We have a new feature as part of our Public Dashboard efforts. We have introduced a new menu item under Dashboards → Public Dashboards. This new view is a list of all of the public dashboards in your Grafana instance. From here you can navigate to the underlying dashboard, see if it is enabled, get a quick link out to the public version of the dashboard, or get quick access to the configuration. You will be able to see a dashboard if you have view access, you will be able to navigate to the configuration if you have public dashboard privileges (RBAC “Dashboard Public” or ADMIN basic).

[image public-dashboard-audit-table.png]

## Public Dashboards - Annotations

Available in experimental in Grafana Open Source, Enterprise, Cloud Pro, Cloud Advanced

Annotations are now supported in public dashboards, with the exception of query annotations. They are turned off by default, but can be turned on in the public dashboard settings modal.

## Transformations - Partition by values

Available in experimental in Grafana Open Source, Enterprise, Cloud Pro, Cloud Advanced

This new transformation can help eliminate the need for multiple queries to the same datasource with different WHERE clauses when graphing multiple series. Consider a metrics SQL table with the following data:

| Time                | Region | Value |
| ------------------- | ------ | ----- |
| 2022-10-20 12:00:00 | US     | 1520  |
| 2022-10-20 12:00:00 | EU     | 2936  |
| 2022-10-20 01:00:00 | US     | 1327  |
| 2022-10-20 01:00:00 | EU     | 912   |

Prior to v9.3. if you wanted to plot a red trendline for US and a blue one for EU in the same TimeSeries panel, you would likely have to split this into two queries:

```
   SELECT Time, Value FROM metrics WHERE Time > ‘2022-10-20’ AND Region=’US’
   SELECT Time, Value FROM metrics WHERE Time > ‘2022-10-20’ AND Region=’EU’
```

This also requires you to know ahead of time which regions actually exist in the metrics table.

With the Partition by values transformer, you can now issue a single query and split the results by unique (enum) values from one or more columns (fields) of your choosing. In this case, Region.

```
   SELECT Time, Region, Value FROM metrics WHERE Time > ‘2022-10-20’
```

| Time                | Region | Value |
| ------------------- | ------ | ----- |
| 2022-10-20 12:00:00 | US     | 1520  |
| 2022-10-20 01:00:00 | US     | 1327  |

| Time                | Region | Value |
| ------------------- | ------ | ----- |
| 2022-10-20 12:00:00 | EU     | 2936  |
| 2022-10-20 01:00:00 | EU     | 912   |

## Authentication - OAuth token handling improvements

This feature is generally available in Grafana Open Source, Enterprise, Cloud Free, Cloud Pro, and Cloud Advanced.

As part of our efforts to improve security of Grafana, we introduce a long-awaited feature which enhances Grafana's OAuth 2.0 compatibility. When a user logs in using an OAuth provider, Grafana verifies on each request that the access token has not expired. When an access token expires, Grafana uses the provided refresh token (if any exists) to obtain a new access token.

This feature is behind the `accessTokenExpirationCheck` feature toggle and it is disabled by default.

Complete documentation on how to configure obtaining a refresh token can be found on the [authentication configuration page]({{< relref "../setup-grafana/configure-security/configure-authentication/" >}}), in the instructions for your Oauth identity provider.

## Grafana CLI tool user management of conflicts

This new CLI command allows you to resolve user identity conflicts of users within Grafana. The email or login field are case sensitive which can cause two or more accounts created for the same user which we call a user identity conflict. You can now resolve these accounts using this CLI tool.

> Note “As a Grafana Cloud user, there are no user identity conflicts, or if you are running a Grafana instance with MySQL.”

```bash
# lists all the conflicting users
$ grafana-cli user-manager conflicts list

# creates a conflict patch file to edit
$ grafana-cli user-manager conflicts generate-file

# reads edited conflict patch file for validation
$ grafana-cli user-manager conflicts validate-file <filepath>

# ingests the conflict users file. Can be executed once per file and will change the state of the database.
$ grafana-cli user-manager conflicts ingest-file <filepath>
```

## Terraform - Extended support for provisioning permissions

All Grafana users can now use the latest release of [Terraform's Grafana provider](https://registry.terraform.io/providers/grafana/grafana/latest/docs) (version 1.31.0) to provision [user and team access to service accounts](https://grafana.com/docs/grafana/latest/administration/service-accounts/#manage-users-and-teams-permissions-for-a-service-account-in-grafana). This allows full management of service accounts through Terraform - from creating a service account and allowing users to access it to assigning roles to the service account and generating service account tokens.

Grafana Enterprise and Cloud Pro and Advanced users can now provision [access to data sources](https://grafana.com/docs/grafana/latest/administration/data-source-management/#data-source-permissions) for Grafana's `Viewer`, `Editor` and `Admin` basic roles, as well as assign `Edit` permission.

We have also added [documentation on provisioning RBAC roles and role assignments](https://grafana.com/docs/grafana/latest/administration/roles-and-permissions/access-control/rbac-terraform-provisioning/) to guide our Grafana Enterprise and Cloud Pro and Advanced users through this process.

Finally, we have fixed several access control related bugs to ensure a smoother provisioning experience.
