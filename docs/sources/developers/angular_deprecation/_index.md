---
keywords:
  - grafana
  - documentation
  - developers
  - resources
labels:
  products:
    - enterprise
    - oss
title: AngularJS support deprecation
weight: 500
---

# AngularJS support deprecation

AngularJS plugin support is deprecated and will be removed in a future release. There are still many community plugins that rely on Grafana's AngularJS plugin support to work. The same is true for many internal (private) plugins that have been developed by Grafana users over the years. Grafana version 9 has a server configuration option that is global to the entire instance and controls whether AngularJS plugin support is available or not. By default, AngularJS support is still enabled, but that will change soon once we complete the migration of all AngularJS code in the core product.

## Why are we deprecating AngularJS support?

AngularJS is an old frontend framework whose active development stopped many years ago. Therefore, it poses a security risk. As Grafana itself had already started migrating to React in v5, this presented the most logical framework for our plugin platform. AngularJS also requires unsafe-eval in the CSP (Content Security Policy) settings, which also reduces the security of running JavaScript in the browser.

## When will AngularJS plugins stop working?

Our goal is to transfer all the remaining AngularJS code to the core of Grafana before Grafana 10 is released in Summer 2023. Once this is done, the option "[angular_support_enabled](https://github.com/grafana/grafana/blob/d61bcdf4ca5e69489e0067c56fbe7f0bfdf84ee4/conf/defaults.ini#L362)" will be disabled by default for new Grafana Cloud users, resulting in the inability to use AngularJS plugins. In case you still rely on AngularJS-based plugins developed internally or by the community, you will need to enable this option to continue using them. Following the release of Grafana 10 we will be migrating Grafana Cloud users where possible and disabling AngularJS support when appropriate, we will also be introducing new features to help all users identify how they are impacted and to warn of the use of deprecated plugins within the Grafana UI.

## When will we remove AngularJS support completely?

Our plan is to completely remove support for AngularJS plugins in version 11, which will be released in 2024. This means that all plugins that depend on AngularJS will stop working and the temporary option introduced in version 10 to enable AngularJS will be removed.

## How do I migrate an AngularJS plugin to React?

Depending on if it’s a data source plugin, panel plugin, or app plugin the process will differ.

For panels, the rendering logic could in some cases be easily preserved but all options need to be redone to use the declarative options framework. For data source plugins the query editor and config options will likely need a total rewrite.

## How do I encourage a community plugin to migrate?

We encourage you to locate the repository of the corresponding plugin and create or upvote an Issue within it if you are using a plugin that is still based on AngularJS.

### Links

- [Migrate from AngularJS to React](https://grafana.com/developers/plugin-tools/migration-guides/migrate-angularjs-to-react)
- [Build a panel plugin](https://grafana.com/tutorials/build-a-panel-plugin/)
- [Build a data source plugin](https://grafana.com/tutorials/build-a-data-source-plugin/)
- [List of plugins using AngularJS]({{< relref "./angular-plugins/" >}})
