---
description: Feature and improvement highlights for Grafana v10.0
keywords:
  - grafana
  - new
  - documentation
  - '10.0'
  - release notes
title: What's new in Grafana v10.0
weight: -37
---

# What’s new in Grafana v10.0

Welcome to Grafana 10.0! Read on to learn about changes to search and navigation, dashboards and visualizations, and authentication and security. For even more detail about all the changes in this release, refer to the [changelog](https://github.com/grafana/grafana/blob/master/CHANGELOG.md).

<!-- Template below
## Feature
[Generally available | Available in private/public preview | Experimental] in Grafana [Open Source, Enterprise, Cloud Free, Cloud Pro, Cloud Advanced]
Description. Include an overview of the feature and problem it solves, and where to learn more (like a link to the docs).
> **Note:** You must use relative references when linking to docs within the Grafana repo. Please do not use absolute URLs. For more information about relrefs, refer to [Links and references](/docs/writers-toolkit/writing-guide/references/).-->

## Authentication and authorization

### No longer possible to disable Role-based access control

_Generally available in all editions of Grafana._

It is no longer possible to disable [Role-based access control (RBAC)](({{< relref "../administration/roles-and-permissions/access-control/" >}})) using a configuration option in Grafana.

We understand that this may affect some users who have relied on the ability to disable RBAC in the past. However, we believe that this change is necessary to ensure the best possible security and user experience for our community and customers.

If you have disabled RBAC for you Grafana instance, refer to our [Upgrade Guide]({{< relref "../upgrade-guide/upgrade-v10.0/index.md" >}}) to check if the change is impacting you and what you can do to mitigate any potential issues.

### SAML UI

_Generally available in Grafana Enterprise, Cloud Pro, and Cloud Advanced._

You can now configure SAML using our new user interface, making the process easier and more convenient than ever before.
With the new user interface, you can now configure SAML without needing to restart Grafana and you can control access to the configuration UI by using [role-based access control (RBAC)]({{< relref "../administration/roles-and-permissions/access-control/" >}}). which makes the process much faster and more efficient.

The SAML UI is available in Grafana Enterprise and Grafana Cloud Pro and Advanced. It is intuitive and user-friendly, with clear instructions and helpful prompts to guide you through the process.

For more information on how to set up SAML using the Grafana UI, refer to [Configure SAML authentication using the Grafana user interface]({{< relref "../setup-grafana/configure-security/configure-authentication/saml-ui/" >}}).

### Case-insensitive usernames and email addresses

_Generally available in all editions of Grafana._

Usernames and email addresses are now treated as case-insensitive, which means that you will no longer need to worry about capitalization when logging in or creating an account.

From now on, whether you type your username or email address in uppercase, lowercase, or a combination of both, Grafana will treat them as the same. This will simplify the login process, reduce the risk of typos and identity conflicts when changing authentication providers.

To help you with dealing with potential user identity conflicts, we have built a [Grafana CLI user identity conflict resolver tool](https://grafana.com/blog/2022/12/12/guide-to-using-the-new-grafana-cli-user-identity-conflict-tool-in-grafana-9.3/) which is available since Grafana 9.3.

Note that if you are running Grafana with MySQL as a database, this change does not have any impact as MySQL users were already treated as case-insensitive.

## Creating subfolders

_Available in preview in all editions of Grafana._

You can now create subfolders in Grafana to help you better organize your dashboards and alerts. This new feature allows you to create, read, update, and delete subfolders, making it easier to sort resources by business units, departments, and teams.

You can also set up permissions using Role-Based Access Control (RBAC). Folder permissions will cascade, being inherited from the parent folder, which simplifies access management.

It's worth noting that the the ability to add subfolders is currently in preview. As such, it's recommended to enable it only on test or development instances, rather than in production environments.

To get started with creating subfolders, you'll need to enable the `nestedFolders` feature toggle. If you’re using Grafana Cloud, and would like to enable this feature, please contact customer support.

In subsequent releases, we’ll be refining and enhancing the user interface for managing dashboards and folders, to provide a more streamlined user experience.

{{< figure src="/media/docs/grafana/screenshot-grafana-10.0-nested-folders.png" max-width="750px" caption="Subfolders in Grafana" >}}

## Correlations

_Available in preview in all editions of Grafana._

You can now bring context from multiple data sources into the Explore experience. Correlations is an extension of our existing Data Links functionality and now enables you to link from any data source, to any data source.

Correlations enable you to seamlessly jump from one data source to another. You define relationships between your different data sources, and when Exploring simply click a button next to a related field in one data source and Grafana will run the corresponding query in the other datasource.

It's worth noting that Correlations is currently in preview. As such, it's recommended to enable it only on test or development instances, rather than in production environments.

To try out Correlations, you'll need to enable the `correlations` feature toggle. If you’re using Grafana Cloud, and would like to enable this feature, please contact customer support.

In subsequent releases, we’ll be refining and enhancing the user interface for Correlations, to provide a more streamlined user experience.

## Dashboards and visualizations

### Enhanced data source selection experience: simplifying querying for users

_Available in preview in all editions of Grafana._

Concepts like data sources and dashboards panels are hard to grasp and it’s a struggle to go from Grafana’s “empty state” to a working dashboard that displays data. Our latest advancements streamline the process of selecting the ideal data source in Grafana, prioritizing recent usage and providing labels and supplementary descriptions.

With this flow, selecting a data source has been greatly simplified, providing a clear overview of available data sources and allowing users to quickly connect to a new one when needed. Additionally, the flow now enables quick uploading of CSV files and access to built-in data sources.

{{< video-embed src="/media/docs/grafana/screen-recording-ds-picker-whats-new-10.mp4" max-width="750px" caption="Datasource picker flow" >}}

To try it out, you'll need to enable the `advancedDataSourcePicker` feature toggle. If you’re using Grafana Cloud, and would like to enable this feature, please contact customer support.

### Public dashboards

_Available in preview in all editions of Grafana._
Public dashboards allow you to share your Grafana dashboard with anyone without requiring them to log in to Grafana. This is useful when you want to make your dashboard available to the world.

With this update, we've made the following improvements:

- The time picker and annotations can be toggled on or off from the configuration modal
- You can see a list of all your public dashboards at **Dashboards > Public dashboards**
- Improved UI: new modal design, paused and not found pages
- Added support for collapsed rows, hidden queries, and zoom into panels

To try it out, you'll need to enable the `publicDashboards` feature toggle. If you’re using Grafana Cloud, and would like to enable this feature, please contact customer support.

To learn more, check out our public dashboards [documentation]({{< relref "../../dashboards/dashboard-public" >}}).

### Public dashboards insights

_Available in preview in Grafana Enterprise, Cloud Pro, and Cloud Advanced._

The public dashboards insights provide valuable information about your public dashboard usage. You can easily access and view important metrics such as the daily query count, the number of views in the last 30 days, and the number of errors in the last 30 days.

To try it out, you'll need to enable the `publicDashboards` feature toggle. If you’re using Grafana Cloud, and would like to enable this feature, please contact customer support.

### Email sharing for public dashboards

_Available in preview in Grafana Enterprise, Cloud Pro, and Cloud Advanced._

Our email sharing feature allows you to easily share your public dashboards with specific individuals. When you add their email addresses, they will receive a one-time link to access the dashboard. This provides you with greater control over who can visualize your data.

We've also added a **Public dashboard users** tab in **Administration > Users** where you can view a list of users who have accessed your public dashboards by way of a email sharing.

To try it out, you'll need to enable the `publicDashboardsEmailSharing` and `publicDashboards` feature toggles. If you’re using Grafana Cloud, and would like to enable this feature, please contact customer support.

### Datagrid

We introduced a new panel, which allows users to edit their data within a Grafana dashboard: The Datagrid
The Datagrid is a new beta panel that allows users to edit their data within a Grafana dashboard. Imagine having a spreadsheet-like view where you can fine-tune data pulled from a datasource or create your own dataset from scratch and use it within your dashboard to update your panels in real time.

In combination with other panels using the Datagrid as a datasource, you can enhance the experience even more.

By default, the Datagrid panel will not allow a user to edit data, but only show data like a standard visualization. To make full use of this new panel editing functionality you will need to enable the “enableDatagridEditing” feature toggle.
But for now the Datagrid Panel will support the following features in G10:

- Creating and deleting rows and columns
- Data and column header edit or delete
- Search functionality
- Column freezing
- Grid selection actions (copy/paste/delete)
- Draggable columns and rows
- Series select when pulling data from a datasource

To learn more about the Datagrid panel, please refer to the official documentation.

The list of features is already growing and of course we will look into those. So stay tuned for more features coming in the near future that will give an even better experience using the datagrid.

### Drag and Drop

Finally connecting your local data is made easier! Welcome drag & drop in Grafana!
The drag & drop functionality allows users to connect their csv, excel or numbers files by simply dragging and dropping them on a dashboard, or even into the editor directly.

Once dropped a table will be created by default containing the data that is living in the file itself.

By default, the drag & drop feature is not available out of the box. In order to utilize this functionality panel you will need to enable the `featuretogglename` feature toggle.

As of now drag & drop supports the following scenarios in G10:

- Drag & drop files on an empty dashboard
- Drag & drop files into the panel editor
- Replace files in the panel editor
- Default table panel creation

As of now the data is being stored on the dashboard json and has a 1MB size limit. Obviously it’s our mission to increase the limit without decreasing the performance in the near future.
To learn more about drag & drop, please refer to the official documentation.

The list of requests is already growing and of course we will look into those. So stay tuned for more features coming in the near future that will give an even better experience using the the drag & drop datasource.

### Canvas GA

We are promoting the canvas panel out of beta and into general availability. Over the past several months we have introduced several enhancements to the panel such as the ability to draw connections between elements, the ability to add data links, and a myriad of UX improvements. We are excited to include canvas as a first class citizen in Grafana’s core panel library. To learn more about the Canvas panel, refer to [Canvas]({{< relref "../panels-visualizations/visualizations/canvas" >}}).

### New Trend Panel

_Available in **beta** in all editions of Grafana_

We are excited to introduce a new panel that allows you to display trends where the x axis is not time. This new beta panel addresses gaps that were not solved by either the Time series or XY Chart panels.

### Time series time region support

We have implemented support for adding time regions to the time series panel. Time regions provide a more contextualized experience, enabling you to highlight certain days of the week, such as Monday-Friday to display work weeks right alongside your data. Time regions are a useful way to highlight specific parts of a day like night, work hours, or whatever you define for each day. They allow the viewer to quickly orient themselves in parts of the day and/or ignore highlighted parts of the time series.

{{< figure src="/media/docs/grafana/screenshot-grafana-10-0-time-regions.png" max-width="750px" caption="Time regions" >}}

### Annotation filtering

We’ve improved the way you can configure annotations by adding the possibility to apply annotations to all panels, selected panels or to select the panels by exclusion.

{{< figure src="/media/docs/grafana/screenshot-grafana-10-0-annotation-filtering.png" max-width="750px" caption="Annotation filtering" >}}

### Redesigned and improved Log Context

Our team has recently made enhancements to Grafana's log context feature, resulting in a more seamless and consistent user experience. With the updated user interface, users can expect the same level of functionality and usability in log context as they would in any other logs panel.

Notably, we've added the following new features that streamline the log context experience:

1. Log details with actions including a **Copy** button to easily copy lines and an eye icon to display only selected labels, allowing users to focus on specific information without leaving the log context panel.
2. A **Wrap Lines** toggle to automatically wrap long lines of text for easier reading and analysis of log entry context directly in log context.
3. An **Open in split view** button to execute the context query for a log entry in a split screen in the Explore view.
4. Only for Loki: A quick filter menu that lets you easily refine the context query by selecting and removing labels.

{{< figure src="/media/docs/grafana/log-context-loki-new-whats-new-10-0.png" max-width="750px" caption="Grafana Log Context" >}}

These improvements make working with log context in Grafana more intuitive and efficient, ultimately improving the overall user experience.

## Scenes

_This is an experimental library_

Scenes is a new front-end library by Grafana that empowers application engineers to effortlessly build stunning dashboard experiences right into their products. With Scenes, you can easily create apps that mirror the Grafana dashboarding experience, complete with template variable support, flexible layouts, dynamic panel rendering, and so much more.

To try it out, please check [@grafana/scenes](https://github.com/grafana/scenes).

<!--
- TODO: Add link to docs and resources
-->

## Query multiple data sources in Explore

_Available in preview in all editions of Grafana using the exploreMixedDatasource feature toggle in all editions of Grafana._

You can now query multiple data sources in Explore. Select "Mixed" from the data source picker and specify a data source for each query.

Mixed data source in Explore is gradually rolling out to all users on Grafana Cloud. If you’re using Grafana Open Source and Enterprise, you can enable this feature using the exploreMixedDatasource feature toggle.

## Tracing

### Span filtering for traces

_Available in preview in all editions of Grafana._

You can now work much more efficiently with traces that consist of a large number of spans.

The span filters exist above the trace view, and allow you to filter the spans that are shown in the trace view. The more filters you add, the more specific are the filtered spans.

Currently, you can add one or more of the following filters:

- Service name,
- Span name,
- Duration,
- Tags (which include tags, process tags, and log fields).

To try it out, you'll need to enable the `newTraceView` feature toggle. If you’re using Grafana Cloud, and would like to enable this feature, please contact customer support.

### OpenTelemetry replacing OpenTracing

We have started our work to migrate to OpenTelemetry in Grafana 8.4 and now we are removing OpenTracing and replacing it, for those that still have it configured, with OpenTelemetry under the hood. The changes are made in backwards compatible way so that users of Grafana do not need to change anything and they will continue working in the same way as before.

## Datasources

### Azure Monitor data source

The Azure Monitor datasource now supports visualizing Application Insights Traces. A new query type `Traces` has been added to the service list. This can be utilised against Application Insights resources to query and visualize traces in both a tabular format and using the built-in Traces visualization.

This also includes support for a new Azure API that will correlate trace ID's against all Application Insights resources that are accessible to the principal that the datasource is configured with. To support this feature a new query builder has been added with support for querying the Application Insights resource using an `Operation ID` or visualizing and filtering the data based on the event type and a subset of the properties available on the trace.

### Prometheus dashboard performance improvements

_This is an experimental feature_

As of Grafana 10, the Prometheus datasource supports delta (incremental) querying, in which values from data frames are cached and leveraged to modify future requests to avoid requesting duplicate values in dashboards with now-relative (i.e. any dashboard querying until "now") queries. This feature is disabled by default as it is still experimental, but can be turned on and configured in the Prometheus data source configuration.
This will reduce network load, and speed up now-relative dashboards, especially for dashboards returning lots of data.

### Phlare renamed to Grafana Pyroscope

As a part of ongoing unification of Phlare and Pyroscope project we are renaming the Phlare data source to Grafana Pyroscope data source. This data source will support both Phlare and Pyroscope backends. Existing instances of the data source should not be affected. When creating new instance of the data source, backend type will be autodetected on the config page, or you can select it manually.

## Alerting

### State history view

Use the improved State history view to get insight into how your alert instances behave over time. View information on when a state change occurred, what the previous state was, the current state, any other alert instances that changed their state at the same time as well as what the query value was that triggered the change.

### Preview notification templates

Preview how your notification templates will look before using them in your contact poitns.

### OnCall contact point

Use the on-call contact point to notify on-call teams on-call about firing alerts.

## Security

### Trusted Types support

_In development, available in all editions of Grafana._

Use [trusted types](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/trusted-types) to reduce the risk of XSS vulnerabilities. This is an experimental web technology with limited browser support. One aspect of it is sanitization of third party libraries or plugins that have not explicitly done sanitization.

To use it in report only mode:

- Enable `content_security_policy_report_only` in the configuration.
- Add `require-trusted-types-for 'script';` to the `content_security_policy_report_only_template`.

To use it in enforce mode:

- Enable `content_security_policy` in the configuration.
- Add `require-trusted-types-for 'script';` to the `content_security_policy_template`.
