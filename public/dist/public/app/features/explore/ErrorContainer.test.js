import { render, screen } from '@testing-library/react';
import React from 'react';
import { ErrorContainer } from './ErrorContainer';
describe('ErrorContainer', () => {
    it('should render component and show message', () => {
        const props = {
            queryError: {
                data: {
                    message: 'Error data message',
                    error: 'Error data content',
                },
                message: 'Error message',
                status: 500,
                statusText: 'Error status text',
                refId: 'A',
            },
        };
        render(React.createElement(ErrorContainer, Object.assign({}, props)));
        const alertEl = screen.getByRole('alert');
        expect(alertEl).toBeInTheDocument();
        expect(alertEl).toHaveTextContent(/query error/i);
        expect(alertEl).toHaveTextContent('Error message');
    });
    it('should render component and show message if message is in data only', () => {
        const props = {
            queryError: {
                data: {
                    message: 'Error data message',
                    error: 'Error data content',
                },
                status: 500,
                statusText: 'Error status text',
                refId: 'A',
            },
        };
        render(React.createElement(ErrorContainer, Object.assign({}, props)));
        const alertEl = screen.getByRole('alert');
        expect(alertEl).toBeInTheDocument();
        expect(alertEl).toHaveTextContent(/query error/i);
        expect(alertEl).toHaveTextContent('Error data message');
    });
    it('should have hidden unknown error if prop is not passed in', () => {
        const props = {};
        render(React.createElement(ErrorContainer, Object.assign({}, props)));
        const alertEl = screen.getByRole('alert', { hidden: true });
        expect(alertEl).toBeInTheDocument();
        expect(alertEl).toHaveTextContent('Unknown error');
    });
});
//# sourceMappingURL=ErrorContainer.test.js.map