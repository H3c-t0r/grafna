var backendSrv = {
    get: jest.fn(),
    getDashboard: jest.fn(),
    getDashboardByUid: jest.fn(),
    getFolderByUid: jest.fn(),
    post: jest.fn(),
};
export function getBackendSrv() {
    return backendSrv;
}
//# sourceMappingURL=backend_srv.js.map