import React from 'react';
import { css } from 'emotion';
import { NavModel } from '@grafana/data';
import Page from '../../core/components/Page/Page';
import { LicenseChrome } from './LicenseChrome';
import { Forms } from '@grafana/ui';
import { hot } from 'react-hot-loader';
import { StoreState } from '../../types';
import { getNavModel } from '../../core/selectors/navModel';
import { connect } from 'react-redux';

interface Props {
  navModel: NavModel;
}

export const UpgradePage: React.FC<Props> = ({ navModel }) => {
  return (
    <Page navModel={navModel}>
      <Page.Contents>
        <UpgradeInfo />
      </Page.Contents>
    </Page>
  );
};

const titleStyles = { fontWeight: 500, fontSize: '26px', lineHeight: '123%' };

const UpgradeInfo: React.FC = () => {
  return (
    <LicenseChrome header="Grafana Enterprise">
      <FeatureInfo />
      <ServiceInfo />
    </LicenseChrome>
  );
};

const GetEnterprise: React.FC = () => {
  return (
    <div style={{ marginTop: '60px', marginBottom: '120px' }}>
      <h2 style={titleStyles}>Get Grafana Enterprise</h2>
      <CallToAction />
      <p style={{ paddingTop: '12px' }}>
        You can use the trial version for free for <strong>30 days</strong>. We will remind you about it
        <strong>5 days before the trial period ends</strong>.
      </p>
    </div>
  );
};

const CallToAction: React.FC = () => {
  return (
    <Forms.Button
      variant="primary"
      size="lg"
      onClick={() =>
        (window.location.href = 'https://grafana.com/contact?about=grafana-enterprise&utm_source=grafana-upgrade-page')
      }
    >
      Contact us and get a free trial
    </Forms.Button>
  );
};

const ServiceInfo: React.FC = () => {
  return (
    <div>
      <h4>At your service</h4>

      <List>
        <Item title="Premium Plugins" image="/public/img/licensing/plugin_enterprise.svg" />
        <Item title="Critical SLA: 2 hours" image="/public/img/licensing/sla.svg" />
        <Item title="Unlimited Expert Support" image="/public/img/licensing/customer_support.svg">
          24x7x365 support via
          <List nested={true}>
            <Item title="Email" />
            <Item title="Private slack channel" />
            <Item title="Phone" />
          </List>
        </Item>
        <Item title="Hand-in-hand support" image="/public/img/licensing/handinhand_support.svg">
          in the upgrade process
        </Item>
      </List>

      <GetEnterprise />
    </div>
  );
};

const FeatureInfo: React.FC = () => {
  return (
    <div style={{ paddingRight: '11px' }}>
      <h4>Enhanced Functionality</h4>
      <FeatureListing />

      <div style={{ marginTop: '20px' }}>
        <strong>Also included:</strong>
        <br />
        Indemnification, working with Grafana Labs on future prioritization, and training from the core Grafana team.
      </div>
    </div>
  );
};

const FeatureListing: React.FC = () => {
  return (
    <List>
      <Item title="Data source permissions" />
      <Item title="Reporting" />
      <Item title="SAML Authentication" />
      <Item title="Enhanced LDAP Integration" />
      <Item title="Team Sync">LDAP, GitHub OAuth, Auth Proxy</Item>
      <Item title="White labeling" />
      <Item title="Premium Plugins">
        <List nested={true}>
          <Item title="Oracle" />
          <Item title="Splunk" />
          <Item title="Service Now" />
          <Item title="Dynatrace" />
          <Item title="DataDog" />
          <Item title="AppDynamics" />
        </List>
      </Item>
    </List>
  );
};

interface ListProps {
  nested?: boolean;
  children: React.ReactNode;
}

const List: React.FC<ListProps> = ({ children, nested }) => {
  const listStyle = css`
    display: flex;
    flex-direction: column;
    padding-top: 8px;

    > div {
      margin-bottom: ${nested ? 0 : 8}px;
    }
  `;

  return <div className={listStyle}>{children}</div>;
};

interface ItemProps {
  title: string;
  image?: string;
  children?: React.ReactNode;
}

const Item: React.FC<ItemProps> = ({ children, title, image }) => {
  const imageUrl = image ? image : '/public/img/licensing/checkmark.svg';
  const itemStyle = css`
    display: flex;

    > img {
      display: block;
      height: 22px;
      flex-grow: 0;
      padding-right: 12px;
    }
  `;

  return (
    <div className={itemStyle}>
      <img src={imageUrl} />
      <div>
        <div style={{ fontWeight: 500, lineHeight: 1.7 }}>{title}</div>
        {children && children}
      </div>
    </div>
  );
};

const mapStateToProps = (state: StoreState) => ({
  navModel: getNavModel(state.navIndex, 'upgrading'),
});

export default hot(module)(connect(mapStateToProps)(UpgradePage));
