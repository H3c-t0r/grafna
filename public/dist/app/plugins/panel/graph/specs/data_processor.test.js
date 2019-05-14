import { DataProcessor } from '../data_processor';
describe('Graph DataProcessor', function () {
    var panel = {
        xaxis: {},
    };
    var dashboard = {};
    var processor = new DataProcessor(panel, dashboard);
    describe('Given default xaxis options and query that returns docs', function () {
        beforeEach(function () {
            panel.xaxis.mode = 'time';
            panel.xaxis.name = 'hostname';
            panel.xaxis.values = [];
            processor.getSeriesList({
                dataList: [
                    {
                        type: 'docs',
                        datapoints: [{ hostname: 'server1', avg: 10 }],
                    },
                ],
            });
        });
        it('Should automatically set xaxis mode to field', function () {
            expect(panel.xaxis.mode).toBe('field');
        });
    });
    describe('getDataFieldNames(', function () {
        var dataList = [
            {
                type: 'docs',
                datapoints: [
                    {
                        hostname: 'server1',
                        valueField: 11,
                        nested: {
                            prop1: 'server2',
                            value2: 23,
                        },
                    },
                ],
            },
        ];
        it('Should return all field names', function () {
            var fields = processor.getDataFieldNames(dataList, false);
            expect(fields).toContain('hostname');
            expect(fields).toContain('valueField');
            expect(fields).toContain('nested.prop1');
            expect(fields).toContain('nested.value2');
        });
        it('Should return all number fields', function () {
            var fields = processor.getDataFieldNames(dataList, true);
            expect(fields).toContain('valueField');
            expect(fields).toContain('nested.value2');
        });
    });
});
//# sourceMappingURL=data_processor.test.js.map