import { getBackendSrv, setBackendSrv } from '@grafana/runtime';
import { getTimeSrv } from 'app/features/dashboard/services/TimeSrv';
import { TemplateSrv } from 'app/features/templating/template_srv';
import { ResourcesAPI } from '../resources/ResourcesAPI';
import { CloudWatchSettings, setupMockedTemplateService } from './CloudWatchDataSource';
export function setupMockedResourcesAPI({ variables, response, getMock, } = {}) {
    let templateService = variables ? setupMockedTemplateService(variables) : new TemplateSrv();
    const timeSrv = getTimeSrv();
    const api = new ResourcesAPI(CloudWatchSettings, templateService);
    let resourceRequestMock = getMock ? getMock : jest.fn().mockReturnValue(response);
    setBackendSrv(Object.assign(Object.assign({}, getBackendSrv()), { get: resourceRequestMock }));
    return { api, resourceRequestMock, templateService, timeSrv };
}
//# sourceMappingURL=ResourcesAPI.js.map