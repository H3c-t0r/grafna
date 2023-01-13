---
aliases:
  - /docs/grafana/latest/guides/whats-new-in-v9-4/
description: Feature and improvement highlights for Grafana v9.4
keywords:
  - grafana
  - new
  - documentation
  - '9.4'
  - release notes
title: What's new in Grafana v9.4
weight: -33
---

# What’s new in Grafana v9.4

Welcome to Grafana 9.4! Read on to learn about [add short list of what's included in this release]. For even more detail about all the changes in this release, refer to the [changelog](https://github.com/grafana/grafana/blob/master/CHANGELOG.md).

## Feature

[Generally available | Available in experimental/beta] in Grafana [Open Source, Enterprise, Cloud Free, Cloud Pro, Cloud Advanced]

Description. Include an overview of the feature and problem it solves, and where to learn more (like a link to the docs).

> **Note:** You must use relative references when linking to docs within the Grafana repo. Please do not use absolute URLs. For more information about relrefs, refer to [Links and references](/docs/writers-toolkit/writing-guide/references/).

## Alert email templating

We've improved the design and functionality of email templates to make template creation much easier and more customizable. The email template framework utilizes MJML to define and compile the final email HTML output. Sprig functions in the email templates provide more customizable template functions.

{{< figure src="/static/img/docs/alerting/alert-templates-whats-new-v9.3.png" max-width="750px" caption="Email template redesign" >}}

## Service account expiration dates

We have included a new configuration option, disabled by default. This will allow us to require an expiration date limit for all newly created service account tokens.

This will not affect existing tokens, however newly created tokens will require an expiration date that doesn't exceed the configuration option `token_expiration_day_limit`.
