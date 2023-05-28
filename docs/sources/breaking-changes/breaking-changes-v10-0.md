---
description: Breaking changes for Grafana v10.0
keywords:
  - grafana
  - breaking changes
  - documentation
  - '10.0'
  - release notes
title: Breaking changes in Grafana v10.0
weight: -1
---

# Breaking changes in Grafana v10.0

Following are breaking changes that you should be aware of when upgrading to Grafana v10.0.

For our purposes, a breaking change is any change that requires users or operators to do something. This includes:

- Changes in one part of the system that could cause other components to fail
- Deprecations or removal of a feature
- Changes to an API that could break automation
- Changes that affect some plugins or functions of Grafana
- Migrations that can’t be rolled back

For each change, the provided information:

- Helps you determine if you’re affected
- Describes the change or relevant background information
- Guides you in how to mitigate for the change or migrate
- Provides more learning resources

For release highlights and deprecations, refer to our [v10.0 What’s new]({{< relref "../whatsnew/whats-new-in-v10-0/" >}}). For the specific steps we recommend when you upgrade to v10.0, check out our [Upgrade guide]({{< relref "../upgrade-guide/upgrade-v10.0/index.md" >}}).

<!-- Template below
## Feature


> **Note:** You must use relative references when linking to docs within the Grafana repo. Please do not use absolute URLs. For more information about relrefs, refer to [Links and references](/docs/writers-toolkit/writing-guide/references/).-->

## General breaking changes

### Angular is deprecated and turned off by default for new Grafana Cloud stacks

#### You are affected if...

You create a new Grafana Cloud stack and intend to use any panel or data source plugins written using the Angular.js JavaScript framework. See the [list of Angular plugins]({{< relref "../../developers/angular_deprecation/angular-plugins/" >}}).

#### Background

AngularJS is an old frontend framework that stopped active development many years ago. Because of that, it is a security risk. AngularJS also requires unsafe-eval in the [CSP (Content Security Policy)](https://developer.mozilla.org/en/Security/CSP) settings which also reduces the security level of how javascript is executed in the browser.

Angular plugin support in Grafana is deprecated, meaning it will be removed in a future release. There are still some community and private plugins built using Angular. Starting with v9.0, Grafana has a server configuration option called [angular_support_enabled]({{< relref "../setup-grafana/configure-grafana/#angular_support_enabled" >}}), that controls whether Angular plugin support is available.

#### Change in Grafana v10.0

Angular will be disabled by default for all new Grafana Cloud stacks. Existing stacks in Grafana Cloud, Grafana Enterprise on-premise instances, and Grafana OSS instances will not be automatically impacted.

#### How to mitigate

If you are using Angular plugins in Grafana, try an alternative panel or data source that does not use Angular. Here is a list of detected [Angular plugins](({{< relref "../developers/angular_deprecation/angular-plugins/" >}})), and their alternatives where they exist.

If you are a developer maintaining a plugin that uses Angular, we recommend you refactor or rewrite its frontend using React.js instead. The team at Grafana is ready to help and provide tools and guidance; you can reach us in [Slack](https://grafana.slack.com/archives/C3HJV5PNE) or our [Forum](https://community.grafana.com/c/plugin-development/30).

#### Learn more

[Angular Support Deprecation docs]({{< relref "../developers/angular_deprecation/" >}})

### Grafana legacy alerting is deprecated and no longer accepts internal or external contributions

#### You are affected if...

You use Grafana legacy alerting and have requested new features or changes to it.

#### Description

Grafana legacy alerting (dashboard alerts) has been deprecated since Grafana 9.0, in favor of the new, improved Grafana Alerting. In Grafana 10, the legacy alerting codebase which depends on Angular will still be available, but we will no longer contribute to it or accept external contributions for it. We will continue to work on the migration path from legacy alerting to Grafana Alerting for our remaining users that need to migrate.

#### Migration path:

The new Grafana Alerting was introduced in Grafana 8 and is a superset of legacy alerting. Learn how to migrate your alerts in the [Migrating Alerts doc]({{< relref "../alerting/migrating-alerts/" >}}).

### API keys are migrating to Service Accounts

#### You're affected if:

You use Grafana API Keys and haven't yet migrated to Service Accounts

#### Description:

In Grafana 8.5 we introduced [Service Accounts](https://grafana.com/blog/2022/08/24/new-in-grafana-9.1-service-accounts-are-now-ga/), which are a superset of API keys that support token rotation and Role-Based Access Control. They became GA in v9.1, with an option to manually migrate API keys to service accounts via UI and API since then. When you upgrade to Grafana 10, Grafana will automatically migrate all API keys to service accounts and hide the API keys screen in the Admin section.

This is a "breaking" change because if users are used to seeing and interacting with API keys, they will not see that page in nav anymore and will need to navigate to the Service Accounts page instead. However, your existing API tokens will remain fully functional and migrated to service accounts, so no automation will break. If you roll back to a previous version of Grafana, your API keys will remain intact.

Grafana's [HTTP API endpoints for generating and managing API Keys]({{< relref "../developers/http_api/auth/#create-api-token" >}}) remain functional, but we recommend you begin using the [Service Account HTTP API]({{< relref "../developers/http_api/serviceaccount/" >}}) to generate and manage machine authentication instead.

#### Migration path:

While upgrading to Grafana 10, you do not need to take any action; your API keys will be automatically migrated. To test or perform the migration from API keys to Service Accounts before upgrading to Grafana 10, follow our [migration documentation]({{< relref "../administration/api-keys/#migrate-api-keys-to-grafana-service-account" >}}").

#### Learn more:

- [Documentation]({{< relref "../administration/api-keys/" >}}) from API keys to service accounts

- [Blog post announcement](https://grafana.com/blog/2022/08/24/new-in-grafana-9.1-service-accounts-are-now-ga/) with a video demo including how to migrate

### The experimental "dashboard previews" feature is removed

#### You're affected if:

You have the dashboardPreviews feature toggle switched on.

#### Description:

We are removing the experimental dashboard previews feature due to user permission and performance issues that prevented us from continuing its development. Dashboard previews have been behind a feature toggle since v9.0.

#### Migration path:

The entire feature has been removed from Grafana. For users who enabled the dashboardPreviews feature flag, Grafana will continue to work as-is but without dashboard previews. We recommend you remove the dashboardPreviews feature flag from your Grafana configuration.

#### Relevant links to learn more:

- [Documentation]({{< relref "../search/dashboard-previews/" >}})

### RBAC is on by default in Grafana Enterprise and cannot be disabled

#### You're affected if:

You have actively disabled RBAC in Grafana's configuration.

#### Description:

Starting from Grafana 9, RBAC has been enabled by default. An option to disable RBAC was introduced as an emergency mechanism and has only been used in rare cases. With Grafana 10 we will remove the configuration option to disable RBAC. This will ensure that Grafana works consistently across different features and decrease the risk of having issues with Grafana's legacy access control. Additionally, access control for some of the Grafana 10 features will only work with the RBAC, so disabling it is no longer an option.

#### Migration path:

No action is needed - migration is automatic. Users' current roles, permissions, SSO mapping, and other authorization functionality will continue to work as before.

### Usernames are now case-insensitive by default 

#### You're affected if:

You run Grafana with a Postgres or sqlite database, you import users from different SSO identity providers (like Google and Active Directory), and users who signed in from different providers were previously created twice because of case differences (for example "<Myra@Grafana.com>" vs "<myra@grafana.com>."

#### Background:

When someone signs up for a Grafana account, they can do it using an email or a login. The fields are case sensitive, which can lead to two or more accounts being created for the same user. Additionally, Grafana allows users to set up an authentication provider, and that provider might return an individual's sign-up email with an uppercased domain name or some combination of uppercase lowercase letters.

Having several accounts leads to split user permissions, confusion among signup flows, and unused "zombie" accounts in your database. Plus, multiple accounts can introduce issues when switching between authentication providers. We refer to these inconsistencies in user uniqueness as a user identity conflict.

#### Change in Grafana 10:

Grafana will start matching users without regard for case. Conflicting users will not be able to sign in.

#### Migration path:

We have built a [CLI tool](https://grafana.com/blog/2022/12/12/guide-to-using-the-new-grafana-cli-user-identity-conflict-tool-in-grafana-9.3/) which users can use to resolve any potential conflicts beforehand. The breaking change was communicated already with Grafana v9.3

#### Relevant links to learn more:

- [Blog post describing usage of the new CLI command](https://grafana.com/blog/2022/12/12/guide-to-using-the-new-grafana-cli-user-identity-conflict-tool-in-grafana-9.3/)

### The "Alias" field in the CloudWatch data source is removed

#### You're affected if:

You are using the "alias" field in the CloudWatch data source, instead of using dynamic labels.

#### Description:

Alias patterns in the CloudWatch query editor were replaced by Label (dynamic labels) behind a feature toggle. Starting from Grafana 9, Label has been enabled by default. With Grafana 10 we remove the option to disable Label and remove the Alias field entirely.

#### Migration path:

Open and save each dashboard that uses the Alias field. Alias is migrated to Label automatically when you load the dashboard.

#### Relevant links to learn more:

[Grafana CloudWatch documentation about the change]({{< relref "../datasources/aws-cloudwatch/query-editor/#label" >}})

## For Plugin Developers

### Upgrading to React 18

#### You're affected if:

You have developed a Grafana plugin that uses the React JS framework (this includes most app, panel, and data source plugins).

#### Description:

Grafana has been upgraded to React 18 and now leverages the new React client rendering API. Plugin authors in particular should be aware, as there could be unintended side effects due to the changes around automatic batching of state updates and consistent `useEffect` timings.

#### Migration path:

As a plugin developer: once a canary build has been released (post-9.5.0), test your plugin against one of the grafana-dev docker builds [here](https://hub.docker.com/r/grafana/grafana-dev/tags?page=1) ([this one](https://hub.docker.com/layers/grafana/grafana-dev/10.0.0-111404pre/images/sha256-ac78acf54b44bd2ce7e68b796b1df47030da7f35e53b02bc3eec3f4de05f780f?context=explore) for example). Because of the extra optimisations made by React 18, some changes may be needed to maintain previous behavior.

#### Relevant links to learn more:

- [React 18 release notes](https://react.dev/blog/2022/03/29/react-v18)

- [React 18 upgrade guide](https://react.dev/blog/2022/03/08/react-18-upgrade-guide)

### Deprecated logs-related functions and variables have been removed from the "@grafana/data" and "@grafana/ui" packages

#### You're affected if:

You are a plugin author and you use one of the following logs-related functions:

- The interface type LogsParser in grafana-data

- The following functions and classes related to logs in the grafana-ui package: LogLabels, LogMessageAnsi, LogRows, getLogRowStyles., getLogLevel, getLogLevelFromKey, addLogLevelToSeries, LogsParsers, calculateFieldStats, calculateLogsLabelStats, calculateStats, getParser, sortInAscendingOrder, sortInDescendingOrder, sortLogsResult, sortLogRows, checkLogsError, escapeUnescapedString.

#### Description:

Certain rarely-used logs-related functions and variables are being moved from grafana-packages to the core grafana codebase. These items have been marked as deprecated previously, and a deprecation-notice was issued in the grafana changelog. Plugin authors who relied on these functions and variables being available will have to adjust their codebase.

#### Migration path:

If you have written a data source or panel plugin, make sure it does not use the functions and variables which have been deprecated and removed.

#### Relevant links to learn more:

<https://github.com/grafana/grafana/issues/65779>

<https://github.com/grafana/grafana/issues/65778>

### DataFrame: Use Array<T> or Vector<T> for field values

#### You're affected if:

You are a plugin author and you have implemented your own version of Vector for data frames.

#### Description:

Working with DataFrames is more complicated than it should be because data is held in a Vector interface rather than a simple Array.  In Grafana 10, the interface has been changed so developers can use simple arrays or Vectors.

#### Migration path:

Any code using Vectors will continue to work without issue.  If you have implemented your own version of Vector it will need to be updated -- the easiest approach is to refactor code so it extends MutableVector.  This will work in both Grafana 9 and 10.

#### Relevant links to learn more:

<https://github.com/grafana/grafana/issues/66480>

### grafana/toolkit CLI commands have been removed and migrated to the create-plugin package

#### You're affected if:

You develop plugins using @grafana/toolkit CLI commands.

#### Description:

We announced the deprecation of the grafana/toolkit in 9.3 (November 2022) and have a new tool called create-plugin now available as a replacement. We encourage you to migrate and use our updated tooling. All grafana/toolkit commands except `build` are removed in Grafana 10.

Here are some of the benefits of create-plugin:

- More flexible: With @grafana/create-plugin, you have more control over your plugins and their dependencies, and can more easily customize the tooling to fit your specific needs.

- Faster development time: With its out-of-the-box development environment @grafana/create-plugin can significantly reduce development time compared to using @grafana/toolkit.

- Improved testing capabilities: Testing plugins with @grafana/create-plugin is much easier with github workflows that automate unit and e2e test runs whenever changes are pushed to github.

- Better documentation: The [documentation](https://grafana.github.io/plugin-tools/docs/creating-a-plugin) for @grafana/create-plugin is more comprehensive and easier to discover than @grafana/toolkit.

#### Migration path:

You may already be using the new tooling. If you have an existing plugin previously created using the @grafana/toolkit you can use the following command to migrate it to the new build tooling:

```

# Run this command from the root of your plugin

cd ./my-plugin

npx @grafana/create-plugin@latest migrate

```

See more details in our [migration guide](https://grafana.github.io/plugin-tools/docs/migrating-from-toolkit/).

#### Relevant links to learn more:

- See more details in our [migration guide](https://grafana.github.io/plugin-tools/docs/migrating-from-toolkit/).
