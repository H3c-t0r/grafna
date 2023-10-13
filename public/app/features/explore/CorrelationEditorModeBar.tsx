import { css } from '@emotion/css';
import React, { useEffect, useState } from 'react';
import { Prompt } from 'react-router-dom';
import { useBeforeUnload, useUnmount } from 'react-use';

import { GrafanaTheme2, colorManipulator } from '@grafana/data';
import { reportInteraction } from '@grafana/runtime';
import { Button, HorizontalGroup, Icon, Tooltip, useStyles2 } from '@grafana/ui';
import { CORRELATION_EDITOR_POST_CONFIRM_ACTION, ExploreItemState, useDispatch, useSelector } from 'app/types';

import { CorrelationUnsavedChangesModal } from './CorrelationUnsavedChangesModal';
import { saveCurrentCorrelation } from './state/correlations';
import { changeDatasource } from './state/datasource';
import { changeCorrelationHelperData } from './state/explorePane';
import { changeCorrelationEditorDetails, splitClose } from './state/main';
import { runQueries } from './state/query';
import { selectCorrelationDetails } from './state/selectors';

export const CorrelationEditorModeBar = ({ panes }: { panes: Array<[string, ExploreItemState]> }) => {
  const dispatch = useDispatch();
  const styles = useStyles2(getStyles);
  const correlationDetails = useSelector(selectCorrelationDetails);
  const [showSavePrompt, setShowSavePrompt] = useState(false);

  // handle refreshing and closing the tab
  useBeforeUnload(correlationDetails?.correlationDirty || false, 'Save correlation?');
  useBeforeUnload(
    (!correlationDetails?.correlationDirty && correlationDetails?.queryEditorDirty) || false,
    'The query editor was changed. Save correlation before continuing?'
  );

  // handle exiting (staying within explore)
  // if we are exiting and
  // we are closing the pane and either are dirty OR
  // we are changing the datasource and the query is dirty
  useEffect(() => {
    if (
      correlationDetails?.isExiting &&
      ((correlationDetails.postConfirmAction?.action === CORRELATION_EDITOR_POST_CONFIRM_ACTION.CLOSE_PANE &&
        (correlationDetails?.correlationDirty || correlationDetails.queryEditorDirty)) ||
        (correlationDetails.postConfirmAction?.action === CORRELATION_EDITOR_POST_CONFIRM_ACTION.CHANGE_DATASOURCE &&
          correlationDetails.queryEditorDirty))
    ) {
      setShowSavePrompt(true);
    } else if (correlationDetails?.isExiting && !correlationDetails?.correlationDirty) {
      dispatch(
        changeCorrelationEditorDetails({
          editorMode: false,
          correlationDirty: false,
          isExiting: false,
        })
      );
    } else if (
      correlationDetails?.isExiting &&
      correlationDetails.postConfirmAction?.action === CORRELATION_EDITOR_POST_CONFIRM_ACTION.CHANGE_DATASOURCE &&
      !correlationDetails.queryEditorDirty
    ) {
      const { exploreId, changeDatasourceUid } = correlationDetails?.postConfirmAction;
      changeCorrelationEditorDetails({
        isExiting: false,
      });
      if (exploreId && changeDatasourceUid) {
        dispatch(changeDatasource(exploreId, changeDatasourceUid, { importQueries: true }));
      }
    }
  }, [correlationDetails, dispatch]);

  // clear data when unmounted
  useUnmount(() => {
    dispatch(
      changeCorrelationEditorDetails({
        editorMode: false,
        isExiting: false,
        correlationDirty: false,
        label: undefined,
        description: undefined,
        canSave: false,
      })
    );

    panes.forEach((pane) => {
      dispatch(
        changeCorrelationHelperData({
          exploreId: pane[0],
          correlationEditorHelperData: undefined,
        })
      );
      dispatch(runQueries({ exploreId: pane[0] }));
    });
  });

  const resetEditor = () => {
    dispatch(
      changeCorrelationEditorDetails({
        editorMode: true,
        isExiting: false,
        correlationDirty: false,
        label: undefined,
        description: undefined,
        canSave: false,
      })
    );

    panes.forEach((pane) => {
      dispatch(
        changeCorrelationHelperData({
          exploreId: pane[0],
          correlationEditorHelperData: undefined,
        })
      );
      dispatch(runQueries({ exploreId: pane[0] }));
    });
  };

  const closePane = (exploreId: string) => {
    setShowSavePrompt(false);
    dispatch(splitClose(exploreId));
    reportInteraction('grafana_explore_split_view_closed');
  };

  const changeDatasourcePostAction = (exploreId: string, datasourceUid: string) => {
    setShowSavePrompt(false);
    dispatch(changeDatasource(exploreId, datasourceUid, { importQueries: true }));
    dispatch(
      changeCorrelationEditorDetails({
        editorMode: true,
        isExiting: false,
        correlationDirty: false,
        queryEditorDirty: false,
        label: undefined,
        description: undefined,
        canSave: false,
      })
    );
    panes.forEach((pane) => {
      dispatch(
        changeCorrelationHelperData({
          exploreId: pane[0],
          correlationEditorHelperData: undefined,
        })
      );
    });
  };

  const saveCorrelationPostAction = (skipPostConfirmAction: boolean) => {
    dispatch(
      saveCurrentCorrelation(
        correlationDetails?.label,
        correlationDetails?.description,
        correlationDetails?.transformations
      )
    );
    if (!skipPostConfirmAction && correlationDetails?.postConfirmAction !== undefined) {
      const { exploreId, action, changeDatasourceUid } = correlationDetails?.postConfirmAction;
      if (action === CORRELATION_EDITOR_POST_CONFIRM_ACTION.CLOSE_PANE) {
        closePane(exploreId);
        resetEditor();
      } else if (
        action === CORRELATION_EDITOR_POST_CONFIRM_ACTION.CHANGE_DATASOURCE &&
        changeDatasourceUid !== undefined
      ) {
        changeDatasource(exploreId, changeDatasourceUid);
      }
    } else {
      dispatch(changeCorrelationEditorDetails({ editorMode: false, correlationDirty: false, isExiting: false }));
    }
  };

  return (
    <>
      {/* Handle navigating outside of Explore */}
      <Prompt
        message={(location) => {
          if (
            location.pathname !== '/explore' &&
            (correlationDetails?.editorMode || false) &&
            (correlationDetails?.correlationDirty || false)
          ) {
            return 'You have unsaved correlation data. Continue?';
          } else {
            return true;
          }
        }}
      />

      {showSavePrompt && (
        <CorrelationUnsavedChangesModal
          onDiscard={() => {
            if (correlationDetails?.postConfirmAction !== undefined) {
              const { exploreId, action, changeDatasourceUid } = correlationDetails?.postConfirmAction;
              if (action === CORRELATION_EDITOR_POST_CONFIRM_ACTION.CLOSE_PANE) {
                closePane(exploreId);
                resetEditor();
              } else if (
                action === CORRELATION_EDITOR_POST_CONFIRM_ACTION.CHANGE_DATASOURCE &&
                changeDatasourceUid !== undefined
              ) {
                changeDatasourcePostAction(exploreId, changeDatasourceUid);
              }
            } else {
              // exit correlations mode
              // if we are discarding the in progress correlation, reset everything
              // this modal only shows if the editorMode is false, so we just need to update the dirty state
              dispatch(
                changeCorrelationEditorDetails({
                  editorMode: false,
                  correlationDirty: false,
                  isExiting: false,
                })
              );
            }
          }}
          onCancel={() => {
            // if we are cancelling the exit, set the editor mode back to true and hide the prompt
            dispatch(changeCorrelationEditorDetails({ isExiting: false }));
            setShowSavePrompt(false);
          }}
          onSave={() => {
            saveCorrelationPostAction(false);
          }}
          dirtyCorrelation={correlationDetails?.correlationDirty || false}
          dirtyQueryEditor={correlationDetails?.queryEditorDirty || false}
          action={correlationDetails?.postConfirmAction!.action}
          isActionLeft={correlationDetails?.postConfirmAction!.isActionLeft}
        />
      )}
      <div className={styles.correlationEditorTop}>
        <HorizontalGroup spacing="md" justify="flex-end">
          <Tooltip content="Correlations editor in Explore is an experimental feature.">
            <Icon className={styles.iconColor} name="info-circle" size="xl" />
          </Tooltip>
          <Button
            variant="secondary"
            disabled={!correlationDetails?.canSave}
            fill="outline"
            className={correlationDetails?.canSave ? styles.buttonColor : styles.disabledButtonColor}
            onClick={() => {
              saveCorrelationPostAction(true);
            }}
          >
            Save
          </Button>
          <Button
            variant="secondary"
            fill="outline"
            className={styles.buttonColor}
            icon="times"
            onClick={() => {
              dispatch(changeCorrelationEditorDetails({ isExiting: true }));
              reportInteraction('grafana_explore_correlation_editor_exit_pressed');
            }}
          >
            Exit correlation editor
          </Button>
        </HorizontalGroup>
      </div>
    </>
  );
};

const getStyles = (theme: GrafanaTheme2) => {
  const contrastColor = theme.colors.getContrastText(theme.colors.primary.main);
  const lighterBackgroundColor = colorManipulator.lighten(theme.colors.primary.main, 0.1);
  const darkerBackgroundColor = colorManipulator.darken(theme.colors.primary.main, 0.2);

  const disabledColor = colorManipulator.darken(contrastColor, 0.2);

  return {
    correlationEditorTop: css({
      backgroundColor: theme.colors.primary.main,
      marginTop: '3px',
      padding: theme.spacing(1),
    }),
    iconColor: css({
      color: contrastColor,
    }),
    buttonColor: css({
      color: contrastColor,
      borderColor: contrastColor,
      '&:hover': {
        color: contrastColor,
        borderColor: contrastColor,
        backgroundColor: lighterBackgroundColor,
      },
    }),
    // important needed to override disabled state styling
    disabledButtonColor: css({
      color: `${disabledColor} !important`,
      backgroundColor: `${darkerBackgroundColor} !important`,
    }),
  };
};
