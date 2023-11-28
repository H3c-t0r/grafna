import { RouteParamsProvider } from '../core/navigation/patch/RouteParamsProvider';
import { RouteProvider } from '../core/navigation/patch/RouteProvider';
import { AngularLocationWrapper } from './AngularLocationWrapper';
import { coreModule } from './core_module';
// Neutralizing Angular’s location tampering
// https://stackoverflow.com/a/19825756
const tamperAngularLocation = () => {
    coreModule.config([
        '$provide',
        ($provide) => {
            $provide.decorator('$browser', [
                '$delegate',
                ($delegate) => {
                    $delegate.onUrlChange = () => { };
                    $delegate.url = () => '';
                    return $delegate;
                },
            ]);
        },
    ]);
};
// Intercepting $location service with implementation based on history
const interceptAngularLocation = () => {
    coreModule.config([
        '$provide',
        ($provide) => {
            $provide.decorator('$location', [
                '$delegate',
                ($delegate) => {
                    $delegate = new AngularLocationWrapper();
                    return $delegate;
                },
            ]);
        },
    ]);
    coreModule.provider('$route', RouteProvider);
    coreModule.provider('$routeParams', RouteParamsProvider);
};
export function initAngularRoutingBridge() {
    tamperAngularLocation();
    interceptAngularLocation();
}
//# sourceMappingURL=bridgeReactAngularRouting.js.map