import { Scope, ScopeSpec, ScopeNode } from '@grafana/data';
import { getBackendSrv } from '@grafana/runtime';
import { ScopedResourceClient } from 'app/features/apiserver/client';
import { NodesMap } from 'app/features/dashboard-scene/scene/Scopes/types';

import { group, namespace, version } from './common';

const nodesEndpoint = `/apis/${group}/${version}/namespaces/${namespace}/find/scope_node_children`;

const client = new ScopedResourceClient<ScopeSpec, 'Scope'>({
  group,
  version,
  resource: 'scopes',
});

const cache = new Map<string, Promise<Scope>>();

async function fetchScopeTreeItems(parent: string, query: string): Promise<ScopeNode[]> {
  try {
    return (await getBackendSrv().get<{ items: ScopeNode[] }>(nodesEndpoint, { parent, query }))?.items ?? [];
  } catch (err) {
    return [];
  }
}

export async function fetchNodes(parent: string, query: string): Promise<NodesMap> {
  return (await fetchScopeTreeItems(parent, query)).reduce<NodesMap>((acc, { metadata: { name }, spec }) => {
    acc[name] = {
      name,
      ...spec,
      isExpandable: spec.nodeType === 'container',
      isSelectable: spec.linkType === 'scope',
      isExpanded: false,
      query: '',
      nodes: {},
    };
    return acc;
  }, {});
}

export async function fetchScope(name: string): Promise<Scope> {
  if (cache.has(name)) {
    return cache.get(name)!;
  }

  const response = new Promise<Scope>(async (resolve) => {
    const basicScope: Scope = {
      metadata: { name },
      spec: {
        filters: [],
        title: name,
        type: '',
        category: '',
        description: '',
      },
    };

    try {
      const serverScope = await client.get(name);

      const scope = {
        ...basicScope,
        metadata: {
          ...basicScope.metadata,
          ...serverScope.metadata,
        },
        spec: {
          ...basicScope.spec,
          ...serverScope.spec,
        },
      };

      resolve(scope);
    } catch (err) {
      cache.delete(name);

      resolve(basicScope);
    }
  });

  cache.set(name, response);

  return response;
}

export async function fetchScopes(names: string[]): Promise<Scope[]> {
  return await Promise.all(names.map(fetchScope));
}
