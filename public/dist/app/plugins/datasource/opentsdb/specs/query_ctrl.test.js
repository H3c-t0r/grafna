import { OpenTsQueryCtrl } from '../query_ctrl';
describe('OpenTsQueryCtrl', function () {
    var ctx = {
        target: { target: '' },
        datasource: {
            tsdbVersion: '',
            getAggregators: function () { return Promise.resolve([]); },
            getFilterTypes: function () { return Promise.resolve([]); },
        },
    };
    ctx.panelCtrl = {
        panel: {
            targets: [ctx.target],
        },
        refresh: function () { },
    };
    OpenTsQueryCtrl.prototype = Object.assign(OpenTsQueryCtrl.prototype, ctx);
    beforeEach(function () {
        ctx.ctrl = new OpenTsQueryCtrl({}, {});
    });
    describe('init query_ctrl variables', function () {
        it('filter types should be initialized', function () {
            expect(ctx.ctrl.filterTypes.length).toBe(7);
        });
        it('aggregators should be initialized', function () {
            expect(ctx.ctrl.aggregators.length).toBe(8);
        });
        it('fill policy options should be initialized', function () {
            expect(ctx.ctrl.fillPolicies.length).toBe(4);
        });
    });
    describe('when adding filters and tags', function () {
        it('addTagMode should be false when closed', function () {
            ctx.ctrl.addTagMode = true;
            ctx.ctrl.closeAddTagMode();
            expect(ctx.ctrl.addTagMode).toBe(false);
        });
        it('addFilterMode should be false when closed', function () {
            ctx.ctrl.addFilterMode = true;
            ctx.ctrl.closeAddFilterMode();
            expect(ctx.ctrl.addFilterMode).toBe(false);
        });
        it('removing a tag from the tags list', function () {
            ctx.ctrl.target.tags = { tagk: 'tag_key', tagk2: 'tag_value2' };
            ctx.ctrl.removeTag('tagk');
            expect(Object.keys(ctx.ctrl.target.tags).length).toBe(1);
        });
        it('removing a filter from the filters list', function () {
            ctx.ctrl.target.filters = [
                {
                    tagk: 'tag_key',
                    filter: 'tag_value2',
                    type: 'wildcard',
                    groupBy: true,
                },
            ];
            ctx.ctrl.removeFilter(0);
            expect(ctx.ctrl.target.filters.length).toBe(0);
        });
        it('adding a filter when tags exist should generate error', function () {
            ctx.ctrl.target.tags = { tagk: 'tag_key', tagk2: 'tag_value2' };
            ctx.ctrl.addFilter();
            expect(ctx.ctrl.errors.filters).toBe('Please remove tags to use filters, tags and filters are mutually exclusive.');
        });
        it('adding a tag when filters exist should generate error', function () {
            ctx.ctrl.target.filters = [
                {
                    tagk: 'tag_key',
                    filter: 'tag_value2',
                    type: 'wildcard',
                    groupBy: true,
                },
            ];
            ctx.ctrl.addTag();
            expect(ctx.ctrl.errors.tags).toBe('Please remove filters to use tags, tags and filters are mutually exclusive.');
        });
    });
});
//# sourceMappingURL=query_ctrl.test.js.map