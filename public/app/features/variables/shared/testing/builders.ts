import { AdHocVariableBuilder } from './adHocVariableBuilder';
import { IntervalVariableBuilder } from './intervalVariableBuilder';
import { DatasourceVariableBuilder } from './datasourceVariableBuilder';
import { OptionsVariableBuilder } from './optionsVariableBuilder';
import { initialQueryVariableModelState } from '../../query/reducer';
import { initialAdHocVariableModelState } from '../../adhoc/reducer';
import { initialDataSourceVariableModelState } from '../../datasource/reducer';
import { initialIntervalVariableModelState } from '../../interval/reducer';
import { initialTextBoxVariableModelState } from '../../textbox/reducer';
import { initialCustomVariableModelState } from '../../custom/reducer';
import { MultiVariableBuilder } from './multiVariableBuilder';
import { initialConstantVariableModelState } from '../../constant/reducer';

export const adHoc = () => new AdHocVariableBuilder(initialAdHocVariableModelState);
export const interval = () => new IntervalVariableBuilder(initialIntervalVariableModelState);
export const datasource = () => new DatasourceVariableBuilder(initialDataSourceVariableModelState);
export const query = () => new DatasourceVariableBuilder(initialQueryVariableModelState);
export const textbox = () => new OptionsVariableBuilder(initialTextBoxVariableModelState);
export const custom = () => new MultiVariableBuilder(initialCustomVariableModelState);
export const constant = () => new OptionsVariableBuilder(initialConstantVariableModelState);
