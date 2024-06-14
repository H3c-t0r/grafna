import { delay, http, HttpResponse } from 'msw';

export const MOCK_GRAFANA_ALERT_RULE_TITLE = 'Test alert';

import {
  RulerGrafanaRuleDTO,
  RulerRuleGroupDTO,
  RulerRulesConfigDTO,
} from '../../../../../../types/unified-alerting-dto';
import { AlertGroupUpdated } from '../../../api/alertRuleApi';
import { grafanaRulerRule, namespaceByUid, namespaces } from '../../alertRuleApi';
import { HandlerOptions } from '../configure';

export const rulerRulesHandler = () => {
  return http.get(`/api/ruler/grafana/api/v1/rules`, () => {
    const response = Object.entries(namespaces).reduce<RulerRulesConfigDTO>((acc, [namespaceUid, groups]) => {
      acc[namespaceByUid[namespaceUid].name] = groups;
      return acc;
    }, {});

    return HttpResponse.json<RulerRulesConfigDTO>(response);
  });
};

export const getRulerRuleNamespaceHandler = () =>
  http.get<{ folderUid: string }>(`/api/ruler/grafana/api/v1/rules/:folderUid`, ({ params: { folderUid } }) => {
    // This mimic API response as closely as possible - Invalid folderUid returns 403
    const namespace = namespaces[folderUid];
    if (!namespace) {
      return new HttpResponse(null, { status: 403 });
    }

    return HttpResponse.json<RulerRulesConfigDTO>({
      [namespaceByUid[folderUid].name]: namespaces[folderUid],
    });
  });

export const updateRulerRuleNamespaceHandler = (options?: HandlerOptions) =>
  http.post<{ folderUid: string }>(`/api/ruler/grafana/api/v1/rules/:folderUid`, async ({ params }) => {
    const { folderUid } = params;

    // @TODO make this more generic so we can use this in other endpoints too
    if (options?.delay) {
      await delay(options.delay ?? 0);
    }

    if (options?.error) {
      return options.error;
    }

    // This mimic API response as closely as possible.
    // Invalid folderUid returns 403 but invalid group will return 202 with empty list of rules
    const namespace = namespaces[folderUid];
    if (!namespace) {
      return new HttpResponse(null, { status: 403 });
    }

    return HttpResponse.json<AlertGroupUpdated>({
      message: 'updated',
      updated: [],
    });
  });

export const rulerRuleGroupHandler = () => {
  return http.get<{ folderUid: string; groupName: string }>(
    `/api/ruler/grafana/api/v1/rules/:folderUid/:groupName`,
    ({ params: { folderUid, groupName } }) => {
      // This mimic API response as closely as possible.
      // Invalid folderUid returns 403 but invalid group will return 202 with empty list of rules
      const namespace = namespaces[folderUid];
      if (!namespace) {
        return new HttpResponse(null, { status: 403 });
      }

      const matchingGroup = namespace.find((group) => group.name === groupName);
      return HttpResponse.json<RulerRuleGroupDTO>({
        name: groupName,
        interval: matchingGroup?.interval,
        rules: matchingGroup?.rules ?? [],
      });
    }
  );
};

export const rulerRuleHandler = () => {
  const grafanaRules = new Map<string, RulerGrafanaRuleDTO>(
    [grafanaRulerRule].map((rule) => [rule.grafana_alert.uid, rule])
  );

  return http.get<{ uid: string }>(`/api/ruler/grafana/api/v1/rule/:uid`, ({ params: { uid } }) => {
    const rule = grafanaRules.get(uid);
    if (!rule) {
      return new HttpResponse(null, { status: 404 });
    }
    return HttpResponse.json(rule);
  });
};

const handlers = [
  rulerRulesHandler(),
  getRulerRuleNamespaceHandler(),
  updateRulerRuleNamespaceHandler(),
  rulerRuleGroupHandler(),
  rulerRuleHandler(),
];
export default handlers;
