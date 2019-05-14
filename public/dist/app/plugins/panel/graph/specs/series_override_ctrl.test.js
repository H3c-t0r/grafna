import '../series_overrides_ctrl';
import { SeriesOverridesCtrl } from '../series_overrides_ctrl';
describe('SeriesOverridesCtrl', function () {
    var popoverSrv = {};
    var $scope;
    beforeEach(function () {
        $scope = {
            ctrl: {
                refresh: jest.fn(),
                render: jest.fn(),
                seriesList: [],
            },
            render: jest.fn(function () { }),
        };
        SeriesOverridesCtrl($scope, {}, popoverSrv);
    });
    describe('When setting an override', function () {
        beforeEach(function () {
            $scope.setOverride({ propertyName: 'lines' }, { value: true });
        });
        it('should set override property', function () {
            expect($scope.override.lines).toBe(true);
        });
        it('should update view model', function () {
            expect($scope.currentOverrides[0].name).toBe('Lines');
            expect($scope.currentOverrides[0].value).toBe('true');
        });
    });
    describe('When removing overide', function () {
        it('click should include option and value index', function () {
            $scope.setOverride(1, 0);
            $scope.removeOverride({ propertyName: 'lines' });
            expect($scope.currentOverrides.length).toBe(0);
        });
    });
});
//# sourceMappingURL=series_override_ctrl.test.js.map