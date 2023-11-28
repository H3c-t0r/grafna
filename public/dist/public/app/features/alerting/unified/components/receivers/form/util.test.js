import { normalizeFormValues } from './util';
describe('normalizeFormValues', () => {
    it('should leave the older config alone', () => {
        const config = createContactPoint({ bearer_token: 'token' });
        expect(normalizeFormValues(config)).toEqual(config);
    });
    it('should leave the older config alone', () => {
        const config = createContactPoint({ bearer_token_file: 'file' });
        expect(normalizeFormValues(config)).toEqual(config);
    });
    it('should normalize newer config', () => {
        const config = createContactPoint({
            authorization: {
                type: 'bearer',
                credentials: 'token',
            },
        });
        expect(normalizeFormValues(config)).toEqual(createContactPoint({ bearer_token: 'token' }));
    });
    it('should normalize newer config', () => {
        const config = createContactPoint({
            authorization: {
                type: 'bearer',
                credentials_file: 'file',
            },
        });
        expect(normalizeFormValues(config)).toEqual(createContactPoint({ bearer_token_file: 'file' }));
    });
    it('should normalize even if authorization is not defined', () => {
        const config = createContactPoint({});
        expect(normalizeFormValues(config)).toEqual(createContactPoint({}));
    });
});
function createContactPoint(httpConfig) {
    const config = {
        name: 'My Contact Point',
        items: [
            {
                __id: '',
                type: '',
                secureSettings: {},
                secureFields: {},
                settings: {
                    http_config: Object.assign({}, httpConfig),
                },
            },
        ],
    };
    return config;
}
//# sourceMappingURL=util.test.js.map