import { createContext } from 'react';
import { DATASOURCES_ROUTES } from '../constants';
// The purpose of this context is to be able to override the data-sources routes (used for links for example) used under
// the app/features/datasources modules, so we can reuse them more easily in different parts of the application (e.g. under Connections)
export const DataSourcesRoutesContext = createContext(DATASOURCES_ROUTES);
//# sourceMappingURL=contexts.js.map