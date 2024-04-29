import { difference, groupBy, take, trim, upperFirst } from 'lodash';
import { ReactNode } from 'react';

import { config } from '@grafana/runtime';
import {
  AlertManagerCortexConfig,
  GrafanaManagedContactPoint,
  GrafanaManagedReceiverConfig,
  MatcherOperator,
  Receiver,
  Route,
} from 'app/plugins/datasource/alertmanager/types';
import { NotifierDTO, NotifierStatus, ReceiversStateDTO } from 'app/types';

import { OnCallIntegrationDTO } from '../../api/onCallApi';
import { computeInheritedTree } from '../../utils/notification-policies';
import { extractReceivers } from '../../utils/receivers';
import { ReceiverTypes } from '../receivers/grafanaAppReceivers/onCall/onCall';
import { getOnCallMetadata, ReceiverPluginMetadata } from '../receivers/grafanaAppReceivers/useReceiversMetadata';

import { RECEIVER_META_KEY, RECEIVER_PLUGIN_META_KEY, RECEIVER_STATUS_KEY } from './useContactPoints';

const AUTOGENERATED_RECEIVER_POLICY_MATCHER_KEY = '__grafana_receiver__';

export function isProvisioned(contactPoint: GrafanaManagedContactPoint) {
  // for some reason the provenance is on the receiver and not the entire contact point
  const provenance = contactPoint.grafana_managed_receiver_configs?.find((receiver) => receiver.provenance)?.provenance;

  return Boolean(provenance);
}

// TODO we should really add some type information to these receiver settings...
export function getReceiverDescription(receiver: ReceiverConfigWithMetadata): ReactNode | undefined {
  if (!receiver.settings) {
    return undefined;
  }
  switch (receiver.type) {
    case 'email': {
      const hasEmailAddresses = 'addresses' in receiver.settings; // when dealing with alertmanager email_configs we don't normalize the settings
      return hasEmailAddresses ? summarizeEmailAddresses(receiver.settings['addresses']) : undefined;
    }
    case 'slack': {
      const recipient: string | undefined = receiver.settings['recipient'];
      if (!recipient) {
        return;
      }

      // Slack channel name might have a "#" in the recipient already
      const channelName = recipient.replace(/^#/, '');
      return `#${channelName}`;
    }
    case 'kafka': {
      const topicName: string | undefined = receiver.settings['kafkaTopic'];
      return topicName;
    }
    case 'webhook': {
      const url: string | undefined = receiver.settings['url'];
      return url;
    }
    case ReceiverTypes.OnCall: {
      return receiver[RECEIVER_PLUGIN_META_KEY]?.description;
    }
    default:
      return receiver[RECEIVER_META_KEY]?.description;
  }
}

// input: foo+1@bar.com, foo+2@bar.com, foo+3@bar.com, foo+4@bar.com
// output: foo+1@bar.com, foo+2@bar.com, +2 more
export function summarizeEmailAddresses(addresses: string): string {
  const MAX_ADDRESSES_SHOWN = 3;
  const SUPPORTED_SEPARATORS = /,|;|\n+/g;

  // split all email addresses
  const emails = addresses.trim().split(SUPPORTED_SEPARATORS).map(trim);

  // grab the first 3 and the rest
  const summary = take(emails, MAX_ADDRESSES_SHOWN);
  const rest = difference(emails, summary);

  if (rest.length) {
    summary.push(`+${rest.length} more`);
  }

  return summary.join(', ');
}

// Grafana Managed contact points have receivers with additional diagnostics
export interface ReceiverConfigWithMetadata extends GrafanaManagedReceiverConfig {
  // we're using a symbol here so we'll never have a conflict on keys for a receiver
  // we also specify that the diagnostics might be "undefined" for vanilla Alertmanager
  [RECEIVER_STATUS_KEY]?: NotifierStatus | undefined;
  [RECEIVER_META_KEY]: {
    name: string;
    description?: string;
  };
  // optional metadata that comes from a particular plugin (like Grafana OnCall)
  [RECEIVER_PLUGIN_META_KEY]?: ReceiverPluginMetadata;
}

export interface ContactPointWithMetadata extends GrafanaManagedContactPoint {
  policies?: RouteReference[]; // now is optional as we don't have the data from the read-only endpoint
  grafana_managed_receiver_configs: ReceiverConfigWithMetadata[];
}

/**
 * This function adds the status information for each of the integrations (contact point types) in a contact point
 * 1. we iterate over all contact points
 * 2. for each contact point we "enhance" it with the status or "undefined" for vanilla Alertmanager
 * contactPoints: list of contact points
 * alertmanagerConfiguration: optional as is passed when we need to get number of policies for each contact point
 * and we prefer using the data from the read-only endpoint.
 */
export function enhanceContactPointsWithMetadata(
  status: ReceiversStateDTO[] = [],
  notifiers: NotifierDTO[] = [],
  onCallIntegrations: OnCallIntegrationDTO[] | undefined | null,
  contactPoints: Receiver[],
  alertmanagerConfiguration?: AlertManagerCortexConfig
): ContactPointWithMetadata[] {
  // compute the entire inherited tree before finding what notification policies are using a particular contact point
  const fullyInheritedTree = computeInheritedTree(alertmanagerConfiguration?.alertmanager_config?.route ?? {});
  const usedContactPoints = getUsedContactPoints(fullyInheritedTree);
  const usedContactPointsByName = groupBy(usedContactPoints, 'receiver');

  const contactPointsList = alertmanagerConfiguration
    ? alertmanagerConfiguration?.alertmanager_config.receivers ?? []
    : contactPoints ?? [];

  return contactPointsList.map((contactPoint) => {
    const receivers = extractReceivers(contactPoint);
    const statusForReceiver = status.find((status) => status.name === contactPoint.name);

    return {
      ...contactPoint,
      policies:
        alertmanagerConfiguration && usedContactPointsByName && (usedContactPointsByName[contactPoint.name] ?? []),
      grafana_managed_receiver_configs: receivers.map((receiver, index) => {
        const isOnCallReceiver = receiver.type === ReceiverTypes.OnCall;
        return {
          ...receiver,
          [RECEIVER_STATUS_KEY]: statusForReceiver?.integrations[index],
          [RECEIVER_META_KEY]: getNotifierMetadata(notifiers, receiver),
          // if OnCall plugin is installed, we'll add it to the receiver's plugin metadata
          [RECEIVER_PLUGIN_META_KEY]: isOnCallReceiver ? getOnCallMetadata(onCallIntegrations, receiver) : undefined,
        };
      }),
    };
  });
}

export function isAutoGeneratedPolicy(route: Route) {
  const simplifiedRoutingToggleEnabled = config.featureToggles.alertingSimplifiedRouting ?? false;
  if (!simplifiedRoutingToggleEnabled) {
    return false;
  }
  if (!route.object_matchers) {
    return false;
  }
  return (
    route.object_matchers.some((objectMatcher) => {
      return (
        objectMatcher[0] === AUTOGENERATED_RECEIVER_POLICY_MATCHER_KEY && objectMatcher[1] === MatcherOperator.equal
      );
    }) ?? false
  );
}

export interface RouteReference {
  receiver: string;
  route: {
    type: 'auto-generated' | 'normal';
  };
}

export function getUsedContactPoints(route: Route): RouteReference[] {
  const childrenContactPoints = route.routes?.flatMap((route) => getUsedContactPoints(route)) ?? [];

  if (route.receiver) {
    return [
      {
        receiver: route.receiver,
        route: {
          type: isAutoGeneratedPolicy(route) ? 'auto-generated' : 'normal',
        },
      },
      ...childrenContactPoints,
    ];
  }

  return childrenContactPoints;
}

function getNotifierMetadata(notifiers: NotifierDTO[], receiver: GrafanaManagedReceiverConfig) {
  const match = notifiers.find((notifier) => notifier.type === receiver.type);

  return {
    name: match?.name ?? upperFirst(receiver.type),
    description: match?.description,
  };
}
