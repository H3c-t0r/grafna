import { Databases, DATABASE_LABELS } from 'app/percona/shared/core';
import { Operators } from './EditDBClusterPage/DBClusterBasicOptions/DBClusterBasicOptions.types';
import { PSMDBService } from './PSMDB.service';
import { XtraDBService } from './XtraDB.service';
export const ADVANCED_SETTINGS_URL = '/graph/settings/advanced-settings';
export const DATABASE_OPTIONS = [
    {
        value: Databases.mysql,
        label: DATABASE_LABELS.mysql,
    },
    {
        value: Databases.mongodb,
        label: DATABASE_LABELS.mongodb,
    },
];
export const SERVICE_MAP = {
    [Databases.mysql]: new XtraDBService(),
    [Databases.mongodb]: new PSMDBService(),
};
export const THOUSAND = 1000;
export const BILLION = Math.pow(10, 9);
export const RESOURCES_PRECISION = 2;
export const DATABASE_OPERATORS = {
    [Operators.pxc]: Databases.mysql,
    [Operators.psmdb]: Databases.mongodb,
};
export const GET_CLUSTERS_CANCEL_TOKEN = 'getClusters';
//# sourceMappingURL=DBCluster.constants.js.map