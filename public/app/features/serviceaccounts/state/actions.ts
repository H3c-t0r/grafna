import { ThunkResult } from '../../../types';
import { getBackendSrv } from '@grafana/runtime';
import { ServiceAccountDTO } from 'app/types';
import { serviceAccountLoaded, serviceAccountsLoaded, serviceAccountTokensLoaded } from './reducers';

const BASE_URL = `/api/serviceaccounts`;

export function loadServiceAccount(id: number): ThunkResult<void> {
  return async (dispatch) => {
    try {
      const response = await getBackendSrv().get(`${BASE_URL}/${id}`);
      dispatch(serviceAccountLoaded(response));
    } catch (error) {
      console.error(error);
    }
  };
}

export function deleteServiceAccountToken(id: number): ThunkResult<void> {
  return async (dispatch, getStore) => {
    const serviceAccount = getStore().serviceAccountProfile.serviceAccount;
    await getBackendSrv().delete(`${BASE_URL}/${serviceAccount.userId}/tokens/${id}`);
    dispatch(loadServiceAccountTokens(serviceAccount.userId));
  };
}

export function loadServiceAccountTokens(saID: number): ThunkResult<void> {
  return async (dispatch) => {
    try {
      const response = await getBackendSrv().get(`${BASE_URL}/${saID}/tokens`);
      dispatch(serviceAccountTokensLoaded(response));
    } catch (error) {
      console.error(error);
    }
  };
}

export function loadServiceAccounts(): ThunkResult<void> {
  return async (dispatch) => {
    try {
      const response = await getBackendSrv().get(BASE_URL);
      dispatch(serviceAccountsLoaded(response));
    } catch (error) {
      console.error(error);
    }
  };
}

export function updateServiceAccount(serviceAccount: ServiceAccountDTO): ThunkResult<void> {
  return async (dispatch) => {
    // TODO: implement on backend
    await getBackendSrv().patch(`${BASE_URL}/${serviceAccount.userId}`, {});
    dispatch(loadServiceAccounts());
  };
}

export function removeServiceAccount(serviceAccountId: number): ThunkResult<void> {
  return async (dispatch) => {
    await getBackendSrv().delete(`${BASE_URL}/${serviceAccountId}`);
    dispatch(loadServiceAccounts());
  };
}
