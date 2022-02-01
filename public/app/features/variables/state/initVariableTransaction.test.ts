import { getPreloadedState, getRootReducer, RootReducerType } from './helpers';
import { variableAdapters } from '../adapters';
import { createQueryVariableAdapter } from '../query/adapter';
import { createConstantVariableAdapter } from '../constant/adapter';
import { reduxTester } from '../../../../test/core/redux/reduxTester';
import {
  addVariable,
  changeVariableProp,
  setCurrentVariableValue,
  variableStateCompleted,
  variableStateFetching,
  variableStateNotStarted,
} from './sharedReducer';
import { adHocBuilder, constantBuilder, datasourceBuilder, queryBuilder } from '../shared/testing/builders';
import { cleanEditorState } from '../editor/reducer';
import {
  initialTransactionState,
  variablesClearTransaction,
  variablesCompleteTransaction,
  variablesInitTransaction,
} from './transactionReducer';
import { cleanPickerState } from '../pickers/OptionsPicker/reducer';
import { cleanVariables } from './variablesReducer';
import { createAdHocVariableAdapter } from '../adhoc/adapter';
import { createDataSourceVariableAdapter } from '../datasource/adapter';
import { DataSourceRef, LoadingState } from '@grafana/data/src';
import { setDataSourceSrv } from '@grafana/runtime/src';
import { TransactionStatus, VariableModel } from '../types';
import { toAsyncOfResult } from '../../query/state/DashboardQueryRunner/testHelpers';
import { setVariableQueryRunner } from '../query/VariableQueryRunner';
import { createDataSourceOptions } from '../datasource/reducer';
import { initVariablesTransaction } from './actions';
import { toKeyedAction } from './keyedVariablesReducer';
import { toVariablePayload } from '../utils';

variableAdapters.setInit(() => [
  createQueryVariableAdapter(),
  createConstantVariableAdapter(),
  createAdHocVariableAdapter(),
  createDataSourceVariableAdapter(),
]);

function getTestContext(variables?: VariableModel[]) {
  const uid = 'uid';
  const constant = constantBuilder().withId('constant').withName('constant').build();
  const templating = { list: variables ?? [constant] };
  const getInstanceSettingsMock = jest.fn().mockReturnValue(undefined);
  setDataSourceSrv({
    get: jest.fn().mockResolvedValue({}),
    getList: jest.fn().mockReturnValue([]),
    getInstanceSettings: getInstanceSettingsMock,
    reload: jest.fn(),
  });
  const variableQueryRunner: any = {
    cancelRequest: jest.fn(),
    queueRequest: jest.fn(),
    getResponse: () => toAsyncOfResult({ state: LoadingState.Done, identifier: { type: 'query', id: 'query' } }),
    destroy: jest.fn(),
  };
  setVariableQueryRunner(variableQueryRunner);

  const dashboard: any = { title: 'Some dash', uid, templating };

  return { constant, getInstanceSettingsMock, templating, uid, dashboard };
}

describe('initVariablesTransaction', () => {
  describe('when called and the previous dashboard has completed', () => {
    it('then correct actions are dispatched', async () => {
      const { constant, uid, dashboard } = getTestContext();
      const tester = await reduxTester<RootReducerType>()
        .givenRootReducer(getRootReducer())
        .whenAsyncActionIsDispatched(initVariablesTransaction(uid, dashboard));

      tester.thenDispatchedActionsPredicateShouldEqual((dispatchedActions) => {
        expect(dispatchedActions[0]).toEqual(toKeyedAction(uid, variablesInitTransaction({ uid })));
        expect(dispatchedActions[1].payload.action.type).toEqual(addVariable.type);
        expect(dispatchedActions[1].payload.action.payload.id).toEqual('__dashboard');
        expect(dispatchedActions[2].payload.action.type).toEqual(addVariable.type);
        expect(dispatchedActions[2].payload.action.payload.id).toEqual('__org');
        expect(dispatchedActions[3].payload.action.type).toEqual(addVariable.type);
        expect(dispatchedActions[3].payload.action.payload.id).toEqual('__user');
        expect(dispatchedActions[4]).toEqual(
          toKeyedAction(uid, addVariable(toVariablePayload(constant, { global: false, index: 0, model: constant })))
        );
        expect(dispatchedActions[5]).toEqual(toKeyedAction(uid, variableStateNotStarted(toVariablePayload(constant))));
        expect(dispatchedActions[6]).toEqual(toKeyedAction(uid, variableStateCompleted(toVariablePayload(constant))));

        expect(dispatchedActions[7]).toEqual(toKeyedAction(uid, variablesCompleteTransaction({ uid })));
        return dispatchedActions.length === 8;
      });
    });

    describe('and there are variables that have data source that need to be migrated', () => {
      it('then correct actions are dispatched', async () => {
        const legacyDs = ('${ds}' as unknown) as DataSourceRef;
        const ds = datasourceBuilder().withId('ds').withDashboardUid('uid').withName('ds').withQuery('prom').build();
        const query = queryBuilder()
          .withId('query')
          .withDashboardUid('uid')
          .withName('query')
          .withDatasource(legacyDs)
          .build();
        const adhoc = adHocBuilder()
          .withId('adhoc')
          .withDashboardUid('uid')
          .withName('adhoc')
          .withDatasource(legacyDs)
          .build();
        const { uid, dashboard } = getTestContext([ds, query, adhoc]);
        const tester = await reduxTester<RootReducerType>()
          .givenRootReducer(getRootReducer())
          .whenAsyncActionIsDispatched(initVariablesTransaction(uid, dashboard));

        tester.thenDispatchedActionsPredicateShouldEqual((dispatchedActions) => {
          expect(dispatchedActions[0]).toEqual(toKeyedAction(uid, variablesInitTransaction({ uid })));
          expect(dispatchedActions[1].payload.action.type).toEqual(addVariable.type);
          expect(dispatchedActions[1].payload.action.payload.id).toEqual('__dashboard');
          expect(dispatchedActions[2].payload.action.type).toEqual(addVariable.type);
          expect(dispatchedActions[2].payload.action.payload.id).toEqual('__org');
          expect(dispatchedActions[3].payload.action.type).toEqual(addVariable.type);
          expect(dispatchedActions[3].payload.action.payload.id).toEqual('__user');
          expect(dispatchedActions[4]).toEqual(
            toKeyedAction(uid, addVariable(toVariablePayload(ds, { global: false, index: 0, model: ds })))
          );
          expect(dispatchedActions[5]).toEqual(
            toKeyedAction(uid, addVariable(toVariablePayload(query, { global: false, index: 1, model: query })))
          );
          expect(dispatchedActions[6]).toEqual(
            toKeyedAction(uid, addVariable(toVariablePayload(adhoc, { global: false, index: 2, model: adhoc })))
          );
          expect(dispatchedActions[7]).toEqual(toKeyedAction(uid, variableStateNotStarted(toVariablePayload(ds))));
          expect(dispatchedActions[8]).toEqual(toKeyedAction(uid, variableStateNotStarted(toVariablePayload(query))));
          expect(dispatchedActions[9]).toEqual(toKeyedAction(uid, variableStateNotStarted(toVariablePayload(adhoc))));
          expect(dispatchedActions[10]).toEqual(
            toKeyedAction(
              uid,
              changeVariableProp(toVariablePayload(query, { propName: 'datasource', propValue: { uid: '${ds}' } }))
            )
          );
          expect(dispatchedActions[11]).toEqual(
            toKeyedAction(
              uid,
              changeVariableProp(toVariablePayload(adhoc, { propName: 'datasource', propValue: { uid: '${ds}' } }))
            )
          );
          expect(dispatchedActions[12]).toEqual(toKeyedAction(uid, variableStateFetching(toVariablePayload(ds))));
          expect(dispatchedActions[13]).toEqual(toKeyedAction(uid, variableStateCompleted(toVariablePayload(adhoc))));
          expect(dispatchedActions[14]).toEqual(
            toKeyedAction(uid, createDataSourceOptions(toVariablePayload(ds, { sources: [], regex: undefined })))
          );
          expect(dispatchedActions[15]).toEqual(
            toKeyedAction(
              uid,
              setCurrentVariableValue(
                toVariablePayload(ds, { option: { selected: false, text: 'No data sources found', value: '' } })
              )
            )
          );
          expect(dispatchedActions[16]).toEqual(toKeyedAction(uid, variableStateCompleted(toVariablePayload(ds))));
          expect(dispatchedActions[17]).toEqual(toKeyedAction(uid, variableStateFetching(toVariablePayload(query))));
          expect(dispatchedActions[18]).toEqual(toKeyedAction(uid, variableStateCompleted(toVariablePayload(query))));
          expect(dispatchedActions[19]).toEqual(toKeyedAction(uid, variablesCompleteTransaction({ uid })));

          return dispatchedActions.length === 20;
        });
      });
    });
  });

  describe('when called and the previous dashboard is still processing variables', () => {
    it('then correct actions are dispatched', async () => {
      const { constant, uid, dashboard } = getTestContext();
      const transactionState = { ...initialTransactionState, uid: 'previous-uid', status: TransactionStatus.Fetching };
      const preloadedState = getPreloadedState(uid, { transaction: transactionState });

      const tester = await reduxTester<RootReducerType>({ preloadedState })
        .givenRootReducer(getRootReducer())
        .whenAsyncActionIsDispatched(initVariablesTransaction(uid, dashboard));

      tester.thenDispatchedActionsPredicateShouldEqual((dispatchedActions) => {
        expect(dispatchedActions[0]).toEqual(toKeyedAction(uid, cleanVariables()));
        expect(dispatchedActions[1]).toEqual(toKeyedAction(uid, cleanEditorState()));
        expect(dispatchedActions[2]).toEqual(toKeyedAction(uid, cleanPickerState()));
        expect(dispatchedActions[3]).toEqual(toKeyedAction(uid, variablesClearTransaction()));
        expect(dispatchedActions[4]).toEqual(toKeyedAction(uid, variablesInitTransaction({ uid })));
        expect(dispatchedActions[5].payload.action.type).toEqual(addVariable.type);
        expect(dispatchedActions[5].payload.action.payload.id).toEqual('__dashboard');
        expect(dispatchedActions[6].payload.action.type).toEqual(addVariable.type);
        expect(dispatchedActions[6].payload.action.payload.id).toEqual('__org');
        expect(dispatchedActions[7].payload.action.type).toEqual(addVariable.type);
        expect(dispatchedActions[7].payload.action.payload.id).toEqual('__user');
        expect(dispatchedActions[8]).toEqual(
          toKeyedAction(uid, addVariable(toVariablePayload(constant, { global: false, index: 0, model: constant })))
        );
        expect(dispatchedActions[9]).toEqual(toKeyedAction(uid, variableStateNotStarted(toVariablePayload(constant))));
        expect(dispatchedActions[10]).toEqual(toKeyedAction(uid, variableStateCompleted(toVariablePayload(constant))));
        expect(dispatchedActions[11]).toEqual(toKeyedAction(uid, variablesCompleteTransaction({ uid })));
        return dispatchedActions.length === 12;
      });
    });
  });
});
