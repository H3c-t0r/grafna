---
description: Install guide for Grafana on Red Hat, RHEL, and Fedora.
title: Install Grafana on Red Hat, RHEL, or Fedora
menuTitle: Redhat, RHEL, or Fedora
weight: 400
---

# Install Grafana on Red Hat, RHEL, or Fedora

This topic explains how to install Grafana dependencies, install Grafana on Redhat, RHEL, or Fedora, and start the Grafana server on your system.

You can install Grafana using a YUM repository, using RPM, or by downloading a binary `.tar.gz` file.

If you install via RPM or the `.tar.gz` file, then you must manually update Grafana for each new version.

## Install Grafana from the YUM respository

If you install from the YUM repository, then Grafana is automatically updated every time you run `sudo yum update`.

| Grafana Version    | Package            | Repository                |
| ------------------ | ------------------ | ------------------------- |
| Grafana Enterprise | grafana-enterprise | `https://rpm.grafana.com` |
| Grafana OSS        | grafana            | `https://rpm.grafana.com` |

> **Note:** Grafana Enterprise is the recommended and default edition. It is available for free and includes all the features of the OSS edition. You can also upgrade to the [full Enterprise feature set](https://grafana.com/products/enterprise/?utm_source=grafana-install-page), which has support for [Enterprise plugins](https://grafana.com/grafana/plugins/?enterprise=1&utcm_source=grafana-install-page).

To install Grafana using a YUM repository, complete the following steps:

1. Add a file to your YUM repository using the method of your choice.

   The following example uses `nano` to add a file to the YUM repo.

   ```bash
   sudo nano /etc/yum.repos.d/grafana.repo
   ```

   ```bash
   [grafana]
   name=grafana
   baseurl=https://rpm.grafana.com
   repo_gpgcheck=1
   enabled=1
   gpgcheck=1
   gpgkey=https://rpm.grafana.com/gpg.key
   sslverify=1
   sslcacert=/etc/pki/tls/certs/ca-bundle.crt
   ```

1. To prevent beta versions from being installed, add the following exclude line to your `.repo` file.

   ```bash
   exclude=*beta*
   ```

1. To install Grafana OSS, run the following command:

   ```bash
   sudo yum install grafana
   ```

1. To install Grafana Enterprise, run the following command:

   ```bash
   sudo yum install grafana-enterprise
   ```

## Install the Grafana RPM package manually

If you install Grafana manually using YUM or RPM, then you must manually update Grafana for each new version. This method varies according to which Linux OS you are running.

**Note:** The RPM files are signed. You can verify the signature with this [public GPG key](https://rpm.grafana.com/gpg.key).

1. On the [Grafana download page](https://grafana.com/grafana/download), select the Grafana version you want to install.
   - The most recent Grafana version is selected by default.
   - The **Version** field displays only finished releases. If you want to install a beta version, click **Nightly Builds** and then select a version.
1. Select an **Edition**.
   - **Enterprise** - Recommended download. Functionally identical to the open source version, but includes features you can unlock with a license if you so choose.
   - **Open Source** - Functionally identical to the Enterprise version, but you will need to download the Enterprise version if you want Enterprise features.
1. Depending on which system you are running, click **Linux** or **ARM**.
1. Copy and paste the RPM package URL and the local RPM package information from the installation page into the pattern shown below, then run the commands.

   ```bash
   sudo yum install initscripts urw-fonts wget
   wget <rpm package url>
   sudo rpm -Uvh <local rpm package>
   ```

## Install Grafana as a standalone binary

Complete the following steps to install Grafana using the standalone binaries:

1. Navigate to the [Grafana download page](https://grafana.com/grafana/download).
1. Select the Grafana version you want to install.
   - The most recent Grafana version is selected by default.
   - The **Version** field displays only tagged releases. If you want to install a nightly build, click **Nightly Builds** and then select a version.
1. Select an **Edition**.
   - **Enterprise:** This is the recommended version. It is functionally identical to the open-source version but includes features you can unlock with a license if you so choose.
   - **Open Source:** This version is functionally identical to the Enterprise version, but you will need to download the Enterprise version if you want Enterprise features.
1. Depending on which system you are running, click the **Linux** or **ARM** tab on the download page.
1. Copy and paste the code from the installation page into your command line and run.

## 2. Start the server

The following sections provide instructions for starting the `grafana-server` process as the `grafana` user created during the package installation.

If you installed using the YUM repository or as an RPM package, then you can start the server using `systemd` or `init.d`. If you installed a binary `.tar.gz` file, then you need to execute the binary.

### Start the Grafana server with systemd

Complete the following steps to start the Grafana server with systemd and verify that it is running:

1. To start the service, run the following commands:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl start grafana-server
   sudo systemctl status grafana-server
   ```

1. To verify that the service is running, run the following command:

   ```
   sudo systemctl status grafana-server
   ```

1. To configure the Grafana server to start at boot, run the following command:

   ```bash
   sudo systemctl enable grafana-server.service
   ```

#### Serving Grafana on a port < 1024

{{< docs/shared "systemd/bind-net-capabilities.md" >}}

#### Serving Grafana behind a proxy

When serving Grafana behind a proxy, you need to configure the `http_proxy` and `https_proxy` environment variables.

### Start the server with init.d

Complete the following steps to start the Grafana service and verify that it is running:

1. To start the Grafana server, run the following commands:

   ```bash
   sudo service grafana-server start
   sudo service grafana-server status
   ```

1. To verify that the service is running, run the following command:

   ```
   sudo service grafana-server status
   ```

1. To configure the Grafana server to start at boot, run the following command:

   ```bash
   sudo update-rc.d grafana-server defaults
   ```

### Start the server using the binary

The `grafana-server` binary .tar.gz needs the working directory to be the root install directory where the binary and the `public` folder are located.

To start the Grafana server, run the following command:

```bash
./bin/grafana-server
```

## Next steps

Refer to the [Getting Started]({{< relref "../../../getting-started/build-first-dashboard/" >}}) guide for information about logging in, setting up data sources, and so on.

## Configure Grafana

Refer to the [Configuration]({{< relref "../../configure-grafana/" >}}) page for details on options for customizing your environment, logging, database, and so on.
