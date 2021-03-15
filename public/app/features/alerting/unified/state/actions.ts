import { createAsyncThunk } from '@reduxjs/toolkit';
import { RuleNamespace } from 'app/types/unified-alerting/internal';
import { fetchRules } from '../api/rules';
import { withSerializedError } from '../utils/redux';

/*
 * Will need to be updated to:
 *
 * 1. Fetch grafana managed rules when the endpoint becomes available
 * 2. Reconcile with rules from the ruler where ruler is available
 */

export const fetchRulesAction = createAsyncThunk(
  'unifiedalerting/fetchRules',
  (datasourceName: string): Promise<RuleNamespace[]> => withSerializedError(fetchRules(datasourceName))
);
