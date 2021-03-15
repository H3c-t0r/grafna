import { RuleNamespace } from 'app/types/alerting-unified/internal';
import { PromRulesResponse } from 'app/types/alerting-unified/dto';
import { getDatasourceByName } from '../utils/config';
import { datasourceRequest, DataSourceType } from '../utils/datasource';

// @TODO currently uses datasource proxy. Will need rework when unified alerting API is operational
export async function fetchRules(datasourceName: string): Promise<RuleNamespace[]> {
  const datasource = getDatasourceByName(datasourceName);
  if (!datasource) {
    throw new Error(`Datasource ${datasourceName} not found`);
  }
  if (datasource.type !== DataSourceType.Loki && datasource.type !== DataSourceType.Prometheus) {
    throw new Error(`Unexpected datasource type ${datasource.type}`);
  }
  const response = await datasourceRequest<PromRulesResponse>(
    datasource.name,
    datasource.type === DataSourceType.Loki ? '/prometheus/api/v1/rules' : '/api/v1/rules'
  );

  if (response.status === 200 && response.data.status === 'success') {
    const nsMap: { [key: string]: RuleNamespace } = {};
    response.data.data.groups.forEach((group) => {
      if (!nsMap[group.file]) {
        nsMap[group.file] = {
          name: group.file,
          groups: [group],
        };
      } else {
        nsMap[group.file].groups.push(group);
      }
    });

    return Object.values(nsMap);
  } else if (response.status === 404) {
    return [];
  } else {
    throw new Error(`http error status=${response.status} body=${JSON.stringify(response.data)}`);
  }
}

export async function checkIfPromIsCortex(datasourceName: string): Promise<boolean> {
  // if it's prometheus, it returns 200 with html, because this api endpoint does not exist and it defaults to index page
  // if it's cortex, it returns 200 with rules in yaml, or if there are no rules, 404
  try {
    const response = await datasourceRequest(datasourceName, '/rules');
    return response.headers.get('content-type') === 'application/yaml';
  } catch (e) {
    if (e.status === 404) {
      return true;
    }
    throw e;
  }
}
