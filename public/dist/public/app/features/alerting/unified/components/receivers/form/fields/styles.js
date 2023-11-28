import { css } from '@emotion/css';
export const getReceiverFormFieldStyles = (theme) => ({
    collapsibleSection: css `
    margin: 0;
    padding: 0;
  `,
    wrapper: css `
    margin: ${theme.spacing(2, 0)};
    padding: ${theme.spacing(1)};
    border: solid 1px ${theme.colors.border.medium};
    border-radius: ${theme.shape.radius.default};
    position: relative;
  `,
    description: css `
    color: ${theme.colors.text.secondary};
    font-size: ${theme.typography.size.sm};
    font-weight: ${theme.typography.fontWeightRegular};
    margin: 0;
  `,
    deleteIcon: css `
    position: absolute;
    right: ${theme.spacing(1)};
    top: ${theme.spacing(1)};
  `,
    addButton: css `
    margin-top: ${theme.spacing(1)};
  `,
});
//# sourceMappingURL=styles.js.map