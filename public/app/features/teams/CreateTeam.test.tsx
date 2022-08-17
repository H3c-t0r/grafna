import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';

import { BackendSrv, setBackendSrv } from '@grafana/runtime';
import * as hooks from 'app/core/components/RolePicker/hooks';

import { CreateTeam } from './CreateTeam';

beforeEach(() => {
  jest.clearAllMocks();
});

jest.mock('app/core/core', () => ({
  contextSrv: {
    licensedAccessControlEnabled: () => false,
    hasPermission: () => true,
    hasPermissionInMetadata: () => true,
    user: { orgId: 1 },
  },
}));

const mockPost = jest.fn(() => {
  return Promise.resolve({});
});

setBackendSrv({
  post: mockPost,
} as any as BackendSrv);

const setup = () => {
  jest.spyOn(hooks, 'useRoleOptions').mockReturnValue([{ roleOptions: [] }, jest.fn()]);
  return render(<CreateTeam />);
};

describe('Create team', () => {
  it('should render component', () => {
    setup();
    expect(screen.getByText(/new team/i)).toBeInTheDocument();
  });

  it('should send correct data to the server', async () => {
    setup();
    await userEvent.type(screen.getByRole('textbox', { name: /name/i }), 'Test team');
    await userEvent.type(screen.getByLabelText(/email/i), 'team@test.com');
    await userEvent.click(screen.getByRole('button', { name: /create/i }));
    await waitFor(() => {
      expect(mockPost).toHaveBeenCalledWith(expect.anything(), { name: 'Test team', email: 'team@test.com' });
    });
  });

  it('should validate required fields', async () => {
    setup();
    await userEvent.type(screen.getByLabelText(/email/i), 'team@test.com');
    await userEvent.click(screen.getByRole('button', { name: /create/i }));
    await waitFor(() => {
      expect(mockPost).not.toBeCalled();
    });
    expect(screen.getAllByRole('alert')).toHaveLength(1);
    expect(screen.getByText(/team name is required/i)).toBeInTheDocument();
  });
});
