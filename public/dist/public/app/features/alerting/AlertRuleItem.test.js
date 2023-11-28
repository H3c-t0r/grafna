import { render, screen, fireEvent } from '@testing-library/react';
import React from 'react';
import AlertRuleItem from './AlertRuleItem';
const setup = (propOverrides) => {
    const props = {
        rule: {
            id: 1,
            dashboardId: 1,
            panelId: 1,
            name: 'Some rule',
            state: 'Open',
            stateText: 'state text',
            stateIcon: 'anchor',
            stateClass: 'state class',
            stateAge: 'age',
            url: 'https://something.something.darkside',
        },
        search: '',
        onTogglePause: jest.fn(),
    };
    Object.assign(props, propOverrides);
    return render(React.createElement(AlertRuleItem, Object.assign({}, props)));
};
describe('AlertRuleItem', () => {
    it('should render component', () => {
        const mockToggle = jest.fn();
        setup({ onTogglePause: mockToggle });
        expect(screen.getByText('Some rule')).toBeInTheDocument();
        expect(screen.getByText('state text')).toBeInTheDocument();
        expect(screen.getByText('Pause')).toBeInTheDocument();
        expect(screen.getByText('Edit alert')).toBeInTheDocument();
        fireEvent.click(screen.getByText('Pause'));
        expect(mockToggle).toHaveBeenCalled();
    });
});
//# sourceMappingURL=AlertRuleItem.test.js.map