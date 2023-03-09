import { Registry } from '@grafana/data';

import { fieldNameByRegexMatcherItem } from './FieldNameByRegexMatcherEditor';
import { fieldNameMatcherItem } from './FieldNameMatcherEditor';
import { fieldNamesMatcherItem } from './FieldNamesMatcherEditor';
import { fieldTypeMatcherItem } from './FieldTypeMatcherEditor';
import { fieldValuesMatcherItem } from './FieldValuesMatcher';
import { fieldsByFrameRefIdItem } from './FieldsByFrameRefIdMatcher';
import { FieldMatcherUIRegistryItem } from './types';

export const fieldMatchersUI = new Registry<FieldMatcherUIRegistryItem<any>>(() => [
  fieldNameMatcherItem,
  fieldNameByRegexMatcherItem,
  fieldTypeMatcherItem,
  fieldsByFrameRefIdItem,
  fieldNamesMatcherItem,
  fieldValuesMatcherItem,
]);
