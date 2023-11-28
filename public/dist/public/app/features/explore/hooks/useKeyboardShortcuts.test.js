import { render } from '@testing-library/react';
import React from 'react';
import { dateTime, EventBusSrv } from '@grafana/data';
import { getAppEvents } from '@grafana/runtime';
import { AbsoluteTimeEvent, ShiftTimeEvent, ShiftTimeEventDirection, ZoomOutEvent } from 'app/types/events';
import { TestProvider } from '../../../../test/helpers/TestProvider';
import { configureStore } from '../../../store/configureStore';
import { initialExploreState } from '../state/main';
import { makeExplorePaneState } from '../state/utils';
import { useKeyboardShortcuts } from './useKeyboardShortcuts';
const testEventBus = new EventBusSrv();
jest.mock('@grafana/runtime', () => {
    return Object.assign(Object.assign({}, jest.requireActual('@grafana/runtime')), { reportInteraction: jest.fn(), getAppEvents: () => testEventBus });
});
const NOW = new Date('2020-10-10T00:00:00.000Z');
function daysFromNow(daysDiff) {
    return new Date(NOW.getTime() + daysDiff * 86400000);
}
function setup() {
    const store = configureStore({
        explore: Object.assign(Object.assign({}, initialExploreState), { panes: {
                left: makeExplorePaneState({
                    range: {
                        from: dateTime(),
                        to: dateTime(),
                        raw: { from: 'now-1d', to: 'now' },
                    },
                }),
                right: makeExplorePaneState({
                    range: {
                        from: dateTime(),
                        to: dateTime(),
                        raw: { from: 'now-2d', to: 'now' },
                    },
                }),
            } }),
    });
    const Wrapper = () => {
        useKeyboardShortcuts();
        return React.createElement("div", null);
    };
    render(React.createElement(TestProvider, { store: store },
        React.createElement(Wrapper, null)));
    return store;
}
describe('useKeyboardShortcuts', () => {
    beforeEach(() => {
        jest.useFakeTimers().setSystemTime(NOW);
    });
    afterEach(() => {
        jest.useRealTimers();
    });
    it('changes both panes to absolute time range', () => {
        const store = setup();
        getAppEvents().publish(new AbsoluteTimeEvent({ updateUrl: false }));
        const exploreState = store.getState().explore;
        const panes = Object.values(exploreState.panes);
        expect(panes[0].absoluteRange.from).toBe(daysFromNow(-1).getTime());
        expect(panes[0].absoluteRange.to).toBe(daysFromNow(0).getTime());
        expect(panes[1].absoluteRange.from).toBe(daysFromNow(-2).getTime());
        expect(panes[1].absoluteRange.to).toBe(daysFromNow(0).getTime());
    });
    it('shifts time range in both panes', () => {
        const store = setup();
        getAppEvents().publish(new ShiftTimeEvent({ direction: ShiftTimeEventDirection.Left }));
        const exploreState = store.getState().explore;
        const panes = Object.values(exploreState.panes);
        expect(panes[0].absoluteRange.from).toBe(daysFromNow(-1.5).getTime());
        expect(panes[0].absoluteRange.to).toBe(daysFromNow(-0.5).getTime());
        expect(panes[1].absoluteRange.from).toBe(daysFromNow(-3).getTime());
        expect(panes[1].absoluteRange.to).toBe(daysFromNow(-1).getTime());
    });
    it('zooms out the time range in both panes', () => {
        const store = setup();
        getAppEvents().publish(new ZoomOutEvent({ scale: 2 }));
        const exploreState = store.getState().explore;
        const panes = Object.values(exploreState.panes);
        expect(panes[0].absoluteRange.from).toBe(daysFromNow(-1.5).getTime());
        expect(panes[0].absoluteRange.to).toBe(daysFromNow(0.5).getTime());
        expect(panes[1].absoluteRange.from).toBe(daysFromNow(-3).getTime());
        expect(panes[1].absoluteRange.to).toBe(daysFromNow(1).getTime());
    });
});
//# sourceMappingURL=useKeyboardShortcuts.test.js.map