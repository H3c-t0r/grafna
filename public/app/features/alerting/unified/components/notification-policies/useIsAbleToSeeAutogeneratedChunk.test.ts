import { contextSrv } from 'app/core/services/context_srv';

import { isAbleToSeeAutogeneratedChunk } from './Policy';

jest.spyOn(contextSrv, 'hasRole');
const hasRoleSpy = jest.mocked(contextSrv.hasRole);

describe('useIsAbleToSeeAutogeneratedChunk hook', () => {
  it('returns true when user is admin', () => {
    hasRoleSpy.mockReturnValue(true);
    expect(isAbleToSeeAutogeneratedChunk()).toBe(true);
  });

  it('returns false when user is not admin', () => {
    hasRoleSpy.mockReturnValue(false);
    expect(isAbleToSeeAutogeneratedChunk()).toBe(false);
  });
});
