// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
import { css } from '@emotion/css';
import React from 'react';
import { IconButton, useStyles2 } from '@grafana/ui';
const getStyles = () => {
    return {
        TimelineCollapser: css `
      align-items: center;
      display: flex;
      flex: none;
      justify-content: center;
      margin-right: 0.5rem;
    `,
    };
};
export function TimelineCollapser(props) {
    const { onExpandAll, onExpandOne, onCollapseAll, onCollapseOne } = props;
    const styles = useStyles2(getStyles);
    return (React.createElement("div", { className: styles.TimelineCollapser, "data-testid": "TimelineCollapser" },
        React.createElement(IconButton, { tooltip: "Expand +1", size: "xl", tooltipPlacement: "top", name: "angle-down", onClick: onExpandOne }),
        React.createElement(IconButton, { tooltip: "Collapse +1", size: "xl", tooltipPlacement: "top", name: "angle-right", onClick: onCollapseOne }),
        React.createElement(IconButton, { tooltip: "Expand All", size: "xl", tooltipPlacement: "top", name: "angle-double-down", onClick: onExpandAll }),
        React.createElement(IconButton, { tooltip: "Collapse All", size: "xl", tooltipPlacement: "top", name: "angle-double-right", onClick: onCollapseAll })));
}
//# sourceMappingURL=TimelineCollapser.js.map