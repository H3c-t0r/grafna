import React from 'react';
import { render, screen } from '@testing-library/react';
import { ErrorContainer, ErrorContainerProps } from './ErrorContainer';

describe('ErrorContainer', () => {
  it('should render component and show message', () => {
    const props: ErrorContainerProps = {
      queryError: {
        data: {
          message: 'Error data message',
          error: 'Error data content',
        },
        message: 'Error message',
        status: 'Error status',
        statusText: 'Error status text',
        refId: 'A',
      },
    };
    render(<ErrorContainer {...props} />);
    expect(screen.getByText('Query error')).toBeInTheDocument();
    expect(screen.getByText('Error message')).toBeInTheDocument();
  });

  it('should render component and show message if message is in data only', () => {
    const props: ErrorContainerProps = {
      queryError: {
        data: {
          message: 'Error data message',
          error: 'Error data content',
        },
        status: 'Error status',
        statusText: 'Error status text',
        refId: 'A',
      },
    };
    render(<ErrorContainer {...props} />);
    expect(screen.getByText('Query error')).toBeInTheDocument();
    expect(screen.getByText('Error data message')).toBeInTheDocument();
  });

  it('should show unknown error if prop is not passed in', () => {
    const props: ErrorContainerProps = {};
    render(<ErrorContainer {...props} />);
    expect(screen.getByText('Unknown error')).toBeInTheDocument();
  });
});
