import { __assign } from "tslib";
import { dateTime } from '@grafana/data';
import { alertRulesReducer, initialChannelState, initialState, loadAlertRules, loadedAlertRules, notificationChannelReducer, setSearchQuery, notificationChannelLoaded, } from './reducers';
import { reducerTester } from '../../../../test/core/redux/reducerTester';
describe('Alert rules', function () {
    var realDateNow = Date.now.bind(global.Date);
    var anchorUnix = dateTime('2019-09-04T10:01:01+02:00').valueOf();
    var dateNowStub = jest.fn(function () { return anchorUnix; });
    global.Date.now = dateNowStub;
    var newStateDate = dateTime().subtract(1, 'y');
    var newStateDateFormatted = newStateDate.format('YYYY-MM-DD');
    var newStateDateAge = newStateDate.fromNow(true);
    var payload = [
        {
            id: 2,
            dashboardId: 7,
            dashboardUid: 'ggHbN42mk',
            dashboardSlug: 'alerting-with-testdata',
            panelId: 4,
            name: 'TestData - Always Alerting',
            state: 'alerting',
            newStateDate: newStateDateFormatted + "T10:00:30+02:00",
            evalDate: '0001-01-01T00:00:00Z',
            evalData: { evalMatches: [{ metric: 'A-series', tags: null, value: 215 }] },
            executionError: '',
            url: '/d/ggHbN42mk/alerting-with-testdata',
        },
        {
            id: 1,
            dashboardId: 7,
            dashboardUid: 'ggHbN42mk',
            dashboardSlug: 'alerting-with-testdata',
            panelId: 3,
            name: 'TestData - Always OK',
            state: 'ok',
            newStateDate: newStateDateFormatted + "T10:01:01+02:00",
            evalDate: '0001-01-01T00:00:00Z',
            evalData: {},
            executionError: '',
            url: '/d/ggHbN42mk/alerting-with-testdata',
        },
        {
            id: 3,
            dashboardId: 7,
            dashboardUid: 'ggHbN42mk',
            dashboardSlug: 'alerting-with-testdata',
            panelId: 3,
            name: 'TestData - ok',
            state: 'ok',
            newStateDate: newStateDateFormatted + "T10:01:01+02:00",
            evalDate: '0001-01-01T00:00:00Z',
            evalData: {},
            executionError: 'error',
            url: '/d/ggHbN42mk/alerting-with-testdata',
        },
        {
            id: 4,
            dashboardId: 7,
            dashboardUid: 'ggHbN42mk',
            dashboardSlug: 'alerting-with-testdata',
            panelId: 3,
            name: 'TestData - Paused',
            state: 'paused',
            newStateDate: newStateDateFormatted + "T10:01:01+02:00",
            evalDate: '0001-01-01T00:00:00Z',
            evalData: {},
            executionError: 'error',
            url: '/d/ggHbN42mk/alerting-with-testdata',
        },
        {
            id: 5,
            dashboardId: 7,
            dashboardUid: 'ggHbN42mk',
            dashboardSlug: 'alerting-with-testdata',
            panelId: 3,
            name: 'TestData - Ok',
            state: 'ok',
            newStateDate: newStateDateFormatted + "T10:01:01+02:00",
            evalDate: '0001-01-01T00:00:00Z',
            evalData: {
                noData: true,
            },
            executionError: 'error',
            url: '/d/ggHbN42mk/alerting-with-testdata',
        },
    ];
    afterAll(function () {
        global.Date.now = realDateNow;
    });
    describe('when loadAlertRules is dispatched', function () {
        it('then state should be correct', function () {
            reducerTester()
                .givenReducer(alertRulesReducer, __assign({}, initialState))
                .whenActionIsDispatched(loadAlertRules())
                .thenStateShouldEqual(__assign(__assign({}, initialState), { isLoading: true }));
        });
    });
    describe('when setSearchQuery is dispatched', function () {
        it('then state should be correct', function () {
            reducerTester()
                .givenReducer(alertRulesReducer, __assign({}, initialState))
                .whenActionIsDispatched(setSearchQuery('query'))
                .thenStateShouldEqual(__assign(__assign({}, initialState), { searchQuery: 'query' }));
        });
    });
    describe('when loadedAlertRules is dispatched', function () {
        it('then state should be correct', function () {
            reducerTester()
                .givenReducer(alertRulesReducer, __assign(__assign({}, initialState), { isLoading: true }))
                .whenActionIsDispatched(loadedAlertRules(payload))
                .thenStateShouldEqual(__assign(__assign({}, initialState), { isLoading: false, items: [
                    {
                        dashboardId: 7,
                        dashboardSlug: 'alerting-with-testdata',
                        dashboardUid: 'ggHbN42mk',
                        evalData: {
                            evalMatches: [
                                {
                                    metric: 'A-series',
                                    tags: null,
                                    value: 215,
                                },
                            ],
                        },
                        evalDate: '0001-01-01T00:00:00Z',
                        executionError: '',
                        id: 2,
                        name: 'TestData - Always Alerting',
                        newStateDate: newStateDateFormatted + "T10:00:30+02:00",
                        panelId: 4,
                        state: 'alerting',
                        stateAge: newStateDateAge,
                        stateClass: 'alert-state-critical',
                        stateIcon: 'heart-break',
                        stateText: 'ALERTING',
                        url: '/d/ggHbN42mk/alerting-with-testdata',
                    },
                    {
                        dashboardId: 7,
                        dashboardSlug: 'alerting-with-testdata',
                        dashboardUid: 'ggHbN42mk',
                        evalData: {},
                        evalDate: '0001-01-01T00:00:00Z',
                        executionError: '',
                        id: 1,
                        name: 'TestData - Always OK',
                        newStateDate: newStateDateFormatted + "T10:01:01+02:00",
                        panelId: 3,
                        state: 'ok',
                        stateAge: newStateDateAge,
                        stateClass: 'alert-state-ok',
                        stateIcon: 'heart',
                        stateText: 'OK',
                        url: '/d/ggHbN42mk/alerting-with-testdata',
                    },
                    {
                        dashboardId: 7,
                        dashboardSlug: 'alerting-with-testdata',
                        dashboardUid: 'ggHbN42mk',
                        evalData: {},
                        evalDate: '0001-01-01T00:00:00Z',
                        executionError: 'error',
                        id: 3,
                        info: 'Execution Error: error',
                        name: 'TestData - ok',
                        newStateDate: newStateDateFormatted + "T10:01:01+02:00",
                        panelId: 3,
                        state: 'ok',
                        stateAge: newStateDateAge,
                        stateClass: 'alert-state-ok',
                        stateIcon: 'heart',
                        stateText: 'OK',
                        url: '/d/ggHbN42mk/alerting-with-testdata',
                    },
                    {
                        dashboardId: 7,
                        dashboardSlug: 'alerting-with-testdata',
                        dashboardUid: 'ggHbN42mk',
                        evalData: {},
                        evalDate: '0001-01-01T00:00:00Z',
                        executionError: 'error',
                        id: 4,
                        name: 'TestData - Paused',
                        newStateDate: newStateDateFormatted + "T10:01:01+02:00",
                        panelId: 3,
                        state: 'paused',
                        stateAge: newStateDateAge,
                        stateClass: 'alert-state-paused',
                        stateIcon: 'pause',
                        stateText: 'PAUSED',
                        url: '/d/ggHbN42mk/alerting-with-testdata',
                    },
                    {
                        dashboardId: 7,
                        dashboardSlug: 'alerting-with-testdata',
                        dashboardUid: 'ggHbN42mk',
                        evalData: {
                            noData: true,
                        },
                        evalDate: '0001-01-01T00:00:00Z',
                        executionError: 'error',
                        id: 5,
                        info: 'Query returned no data',
                        name: 'TestData - Ok',
                        newStateDate: newStateDateFormatted + "T10:01:01+02:00",
                        panelId: 3,
                        state: 'ok',
                        stateAge: newStateDateAge,
                        stateClass: 'alert-state-ok',
                        stateIcon: 'heart',
                        stateText: 'OK',
                        url: '/d/ggHbN42mk/alerting-with-testdata',
                    },
                ] }));
        });
    });
});
describe('Notification channel', function () {
    var notifiers = [
        {
            type: 'webhook',
            name: 'webhook',
            heading: 'Webhook settings',
            description: 'Sends HTTP POST request to a URL',
            info: '',
            options: [
                {
                    element: 'input',
                    inputType: 'text',
                    label: 'Url',
                    description: '',
                    placeholder: '',
                    propertyName: 'url',
                    showWhen: { field: '', is: '' },
                    required: true,
                    validationRule: '',
                    secure: false,
                },
                {
                    element: 'select',
                    inputType: '',
                    label: 'Http Method',
                    description: '',
                    placeholder: '',
                    propertyName: 'httpMethod',
                    selectOptions: [
                        { value: 'POST', label: 'POST' },
                        { value: 'PUT', label: 'PUT' },
                    ],
                    showWhen: { field: '', is: '' },
                    required: false,
                    validationRule: '',
                    secure: false,
                },
                {
                    element: 'input',
                    inputType: 'text',
                    label: 'Username',
                    description: '',
                    placeholder: '',
                    propertyName: 'username',
                    showWhen: { field: '', is: '' },
                    required: false,
                    validationRule: '',
                    secure: false,
                },
                {
                    element: 'input',
                    inputType: 'password',
                    label: 'Password',
                    description: '',
                    placeholder: '',
                    propertyName: 'password',
                    showWhen: { field: '', is: '' },
                    required: false,
                    validationRule: '',
                    secure: true,
                },
            ],
        },
    ];
    describe('Load notification channel', function () {
        it('should migrate non secure settings to secure fields', function () {
            var payload = {
                id: 2,
                uid: '9L3FrrHGk',
                name: 'Webhook test',
                type: 'webhook',
                isDefault: false,
                sendReminder: false,
                disableResolveMessage: false,
                frequency: '',
                created: '2020-08-28T08:49:24Z',
                updated: '2020-08-28T08:49:24Z',
                settings: {
                    autoResolve: true,
                    httpMethod: 'POST',
                    password: 'asdf',
                    severity: 'critical',
                    uploadImage: true,
                    url: 'http://localhost.webhook',
                    username: 'asdf',
                },
            };
            var expected = {
                id: 2,
                uid: '9L3FrrHGk',
                name: 'Webhook test',
                type: 'webhook',
                isDefault: false,
                sendReminder: false,
                disableResolveMessage: false,
                frequency: '',
                created: '2020-08-28T08:49:24Z',
                updated: '2020-08-28T08:49:24Z',
                secureSettings: {
                    password: 'asdf',
                },
                settings: {
                    autoResolve: true,
                    httpMethod: 'POST',
                    password: '',
                    severity: 'critical',
                    uploadImage: true,
                    url: 'http://localhost.webhook',
                    username: 'asdf',
                },
            };
            reducerTester()
                .givenReducer(notificationChannelReducer, __assign(__assign({}, initialChannelState), { notifiers: notifiers }))
                .whenActionIsDispatched(notificationChannelLoaded(payload))
                .thenStateShouldEqual(__assign(__assign({}, initialChannelState), { notifiers: notifiers, notificationChannel: expected }));
        });
        it('should handle already secure field', function () {
            var payload = {
                id: 2,
                uid: '9L3FrrHGk',
                name: 'Webhook test',
                type: 'webhook',
                isDefault: false,
                sendReminder: false,
                disableResolveMessage: false,
                frequency: '',
                created: '2020-08-28T08:49:24Z',
                updated: '2020-08-28T08:49:24Z',
                secureFields: {
                    password: true,
                },
                settings: {
                    autoResolve: true,
                    httpMethod: 'POST',
                    password: '',
                    severity: 'critical',
                    uploadImage: true,
                    url: 'http://localhost.webhook',
                    username: 'asdf',
                },
            };
            var expected = {
                id: 2,
                uid: '9L3FrrHGk',
                name: 'Webhook test',
                type: 'webhook',
                isDefault: false,
                sendReminder: false,
                disableResolveMessage: false,
                frequency: '',
                created: '2020-08-28T08:49:24Z',
                updated: '2020-08-28T08:49:24Z',
                secureFields: {
                    password: true,
                },
                settings: {
                    autoResolve: true,
                    httpMethod: 'POST',
                    password: '',
                    severity: 'critical',
                    uploadImage: true,
                    url: 'http://localhost.webhook',
                    username: 'asdf',
                },
            };
            reducerTester()
                .givenReducer(notificationChannelReducer, __assign(__assign({}, initialChannelState), { notifiers: notifiers }))
                .whenActionIsDispatched(notificationChannelLoaded(payload))
                .thenStateShouldEqual(__assign(__assign({}, initialChannelState), { notifiers: notifiers, notificationChannel: expected }));
        });
    });
});
//# sourceMappingURL=reducers.test.js.map