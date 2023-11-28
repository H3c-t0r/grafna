import { __rest } from "tslib";
/* eslint-disable @typescript-eslint/no-explicit-any, @typescript-eslint/consistent-type-assertions */
import React from 'react';
import { Select } from '@grafana/ui';
import { LabelCore } from 'app/percona/shared/components/Form/LabelCore';
import { withSelectStyles } from '../withSelectStyles/withSelectStyles';
const SelectFieldWrapper = (_a) => {
    var { label, name, required, inputId, tooltipLink, tooltipText, tooltipLinkText, tooltipDataTestId, tooltipIcon, tooltipLinkTarget } = _a, props = __rest(_a, ["label", "name", "required", "inputId", "tooltipLink", "tooltipText", "tooltipLinkText", "tooltipDataTestId", "tooltipIcon", "tooltipLinkTarget"]);
    return (React.createElement(React.Fragment, null,
        React.createElement(LabelCore, { name: name, label: label, required: required, inputId: inputId, tooltipLink: tooltipLink, tooltipLinkText: tooltipLinkText, tooltipText: tooltipText, tooltipDataTestId: tooltipDataTestId, tooltipLinkTarget: tooltipLinkTarget, tooltipIcon: tooltipIcon }),
        React.createElement(Select, Object.assign({}, props))));
};
export const SelectField = withSelectStyles(SelectFieldWrapper);
//# sourceMappingURL=SelectField.js.map