import { css, cx } from '@emotion/css';
import { clone } from 'lodash';
import React, { useMemo } from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import {
  SceneObjectState,
  SceneObjectBase,
  SceneComponentProps,
  SceneVariableValueChangedEvent,
  SceneObjectStateChangedEvent,
  SceneTimeRange,
  sceneUtils,
  AdHocFiltersVariable,
  SceneVariableSet,
} from '@grafana/scenes';
import { useStyles2, Tooltip, Stack, AdHocFilterItem } from '@grafana/ui';

import { DataTrail, DataTrailState, getTopSceneFor } from './DataTrail';
import { reportExploreMetrics } from './interactions';
import { VAR_FILTERS } from './shared';
import { getTrailFor, isSceneTimeRangeState } from './utils';

export interface DataTrailsHistoryState extends SceneObjectState {
  currentStep: number;
  steps: DataTrailHistoryStep[];
}

export interface DataTrailHistoryStep {
  description: string;
  type: TrailStepType;
  trailState: DataTrailState;
  parentIndex: number;
}

export type TrailStepType = 'filters' | 'time' | 'metric' | 'start';
export class DataTrailHistory extends SceneObjectBase<DataTrailsHistoryState> {
  public constructor(state: Partial<DataTrailsHistoryState>) {
    super({ steps: state.steps ?? [], currentStep: state.currentStep ?? 0 });

    this.addActivationHandler(this._onActivate.bind(this));
  }

  private stepTransitionInProgress = false;

  public _onActivate() {
    const trail = getTrailFor(this);

    if (this.state.steps.length === 0) {
      // We always want to ensure in initial 'start' step
      this.addTrailStep(trail, 'start');

      if (trail.state.metric) {
        // But if our current trail has a metric, we want to remove it and the topScene,
        // so that the "start" step always displays a metric select screen.

        // So we remove the metric and update the topscene for the "start" step
        const { metric, ...startState } = trail.state;
        startState.topScene = getTopSceneFor(undefined);
        this.state.steps[0].trailState = startState;

        // But must add a secondary step to represent the selection of the metric
        // for this restored trail state
        this.addTrailStep(trail, 'metric');
      }
    }

    trail.subscribeToState((newState, oldState) => {
      if (newState.metric !== oldState.metric) {
        if (this.state.steps.length === 1) {
          // For the first step we want to update the starting state so that it contains data
          this.state.steps[0].trailState = sceneUtils.cloneSceneObjectState(oldState, { history: this });
        }

        if (newState.metric || oldState.metric) {
          this.addTrailStep(trail, 'metric');
        }
      }
    });

    let respondingToFilterChange = false;
    let respondingToTimeChange = false;
    trail.subscribeToEvent(SceneObjectStateChangedEvent, (evt) => {
      if (evt.payload.changedObject instanceof SceneTimeRange && !respondingToTimeChange) {
        const { prevState, newState } = evt.payload;

        if (isSceneTimeRangeState(prevState) && isSceneTimeRangeState(newState)) {
          if (prevState.from === newState.from && prevState.to === newState.to) {
            return;
          }
        }

        respondingToTimeChange = true;
        const previousStep = this.state.currentStep;
        this.addTrailStep(trail, 'time');

        // Ensure the previous trail step keeps the previous state, using a clone of the previous time range
        const $timeRangePrevious = this.state.steps[previousStep].trailState.$timeRange;
        if ($timeRangePrevious) {
          this.state.steps[previousStep].trailState.$timeRange = new SceneTimeRange(
            sceneUtils.cloneSceneObjectState($timeRangePrevious.state, prevState)
          );
        }

        respondingToTimeChange = false;
      } else {
        if (
          evt.payload.changedObject instanceof AdHocFiltersVariable &&
          evt.payload.changedObject.state.name === VAR_FILTERS &&
          !respondingToFilterChange
        ) {
          const { prevState, newState } = evt.payload;

          if (hasFilters(prevState) && hasFilters(newState)) {
            if (prevState.filterExpression === newState.filterExpression) {
              return;
            }
          }

          respondingToFilterChange = true;
          const previousStep = this.state.currentStep;
          this.addTrailStep(trail, 'filters');

          if (hasFilters(prevState)) {
            this.state.steps[previousStep].trailState.initialFilters = clone(prevState.filters);
          }

          // // Ensure the previous trail step keeps the previous filter state
          const $variablesPrevious = this.state.steps[previousStep].trailState.$variables;
          if ($variablesPrevious) {
            const $variablesPreviousClone = $variablesPrevious.clone();

            const index = $variablesPreviousClone.state.variables.findIndex(
              (variable) => variable.state.name === VAR_FILTERS
            );

            const variables = [...$variablesPreviousClone.state.variables];
            variables[index] = new AdHocFiltersVariable(sceneUtils.cloneSceneObjectState(prevState));

            $variablesPreviousClone.setState({ variables });
            this.state.steps[previousStep].trailState.$variables = new SceneVariableSet(
              sceneUtils.cloneSceneObjectState($variablesPreviousClone.state)
            );
          }

          respondingToFilterChange = false;
        }
      }
    });
  }

  public addTrailStep(trail: DataTrail, type: TrailStepType) {
    if (this.stepTransitionInProgress) {
      // Do not add trail steps when step transition is in progress
      return;
    }

    const stepIndex = this.state.steps.length;
    const parentIndex = type === 'start' ? -1 : this.state.currentStep;

    this.setState({
      currentStep: stepIndex,
      steps: [
        ...this.state.steps,
        {
          description: 'Test',
          type,
          trailState: sceneUtils.cloneSceneObjectState(trail.state, { history: this }),
          parentIndex,
        },
      ],
    });
  }

  public goBackToStep(stepIndex: number) {
    if (stepIndex === this.state.currentStep) {
      return;
    }

    this.stepTransitionInProgress = true;
    const step = this.state.steps[stepIndex];
    const type = step.type === 'metric' && step.trailState.metric === undefined ? 'metric-clear' : step.type;
    reportExploreMetrics('history_step_clicked', { type });

    this.setState({ currentStep: stepIndex });
    // The URL will update

    this.stepTransitionInProgress = false;
  }

  renderStepTooltip(step: DataTrailHistoryStep) {
    return (
      <Stack direction="column">
        <div>{step.type}</div>
        {step.type === 'metric' && <div>{step.trailState.metric || 'Select new metric'}</div>}
      </Stack>
    );
  }

  public static Component = ({ model }: SceneComponentProps<DataTrailHistory>) => {
    const { steps, currentStep } = model.useState();
    const styles = useStyles2(getStyles);

    const { ancestry, alternatePredecessorStyle } = useMemo(() => {
      const ancestry = new Set<number>();

      let cursor = currentStep;
      while (cursor >= 0) {
        const step = steps[cursor];
        if (!step) {
          break;
        }
        ancestry.add(cursor);
        cursor = step.parentIndex;
      }

      const alternatePredecessorStyle = new Map<number, string>();

      ancestry.forEach((index) => {
        const parent = steps[index].parentIndex;
        if (parent + 1 !== index) {
          alternatePredecessorStyle.set(index, createAlternatePredecessorStyle(index, parent));
        }
      });

      return { ancestry, alternatePredecessorStyle };
    }, [currentStep, steps]);

    return (
      <div className={styles.container}>
        <div className={styles.heading}>History</div>
        {steps.map((step, index) => {
          let stepType = step.type;

          if (stepType === 'metric' && step.trailState.metric === undefined) {
            // If we're resetting the metric, we want it to look like a start node
            stepType = 'start';
          }

          return (
            <Tooltip content={() => model.renderStepTooltip(step)} key={index}>
              <button
                className={cx(
                  // Base for all steps
                  styles.step,
                  // Specifics per step type
                  styles.stepTypes[stepType],
                  // To highlight selected step
                  model.state.currentStep === index ? styles.stepSelected : '',
                  // To alter the look of steps with distant non-directly preceding parent
                  alternatePredecessorStyle.get(index) ?? '',
                  // To remove direct link for steps that don't have a direct parent
                  index !== step.parentIndex + 1 ? styles.stepOmitsDirectLeftLink : '',
                  // To remove the direct parent link on the start node as well
                  index === 0 ? styles.stepOmitsDirectLeftLink : '',
                  // To darken steps that aren't the current step's ancesters
                  !ancestry.has(index) ? styles.stepIsNotAncestorOfCurrent : ''
                )}
                onClick={() => model.goBackToStep(index)}
              ></button>
            </Tooltip>
          );
        })}
      </div>
    );
  };
}

function getStyles(theme: GrafanaTheme2) {
  const visTheme = theme.visualization;

  return {
    container: css({
      display: 'flex',
      gap: 10,
      alignItems: 'center',
    }),
    heading: css({}),
    step: css({
      flexGrow: 0,
      cursor: 'pointer',
      border: 'none',
      boxShadow: 'none',
      padding: 0,
      margin: 0,
      width: 8,
      height: 8,
      opacity: 0.7,
      borderRadius: theme.shape.radius.circle,
      background: theme.colors.primary.main,
      position: 'relative',
      '&:hover': {
        opacity: 1,
      },
      '&:hover:before': {
        // We only want the node to hover, not its connection to its parent
        opacity: 0.7,
      },
      '&:before': {
        content: '""',
        position: 'absolute',
        width: 10,
        height: 2,
        left: -10,
        top: 3,
        background: theme.colors.primary.border,
        pointerEvents: 'none',
      },
    }),
    stepSelected: css({
      '&:after': {
        content: '""',
        borderStyle: `solid`,
        borderWidth: 2,
        borderRadius: '50%',
        position: 'absolute',
        width: 16,
        height: 16,
        left: -4,
        top: -4,
        boxShadow: `0px 0px 0px 2px inset ${theme.colors.background.canvas}`,
      },
    }),
    stepOmitsDirectLeftLink: css({
      '&:before': {
        background: 'none',
      },
    }),
    stepIsNotAncestorOfCurrent: css({
      opacity: 0.2,
      '&:hover:before': {
        opacity: 0.2,
      },
    }),
    stepTypes: {
      start: generateStepTypeStyle(visTheme.getColorByName('green')),
      filters: generateStepTypeStyle(visTheme.getColorByName('purple')),
      metric: generateStepTypeStyle(visTheme.getColorByName('orange')),
      time: generateStepTypeStyle(theme.colors.primary.main),
    },
  };
}

function generateStepTypeStyle(color: string) {
  return css({
    background: color,
    '&:before': {
      background: color,
      borderColor: color,
    },
    '&:after': {
      borderColor: color,
    },
  });
}

function createAlternatePredecessorStyle(index: number, parent: number) {
  const difference = index - parent;

  const NODE_DISTANCE = 18;
  const distanceToParent = difference * NODE_DISTANCE;

  return css({
    '&:before': {
      content: '""',
      width: distanceToParent + 2,
      height: 10,
      borderStyle: 'solid',
      borderWidth: 2,
      borderBottom: 'none',
      borderTopLeftRadius: 8,
      borderTopRightRadius: 8,
      top: -10,
      left: 3 - distanceToParent,
      background: 'none',
    },
  });
}

type StateWithFilterExpression = SceneObjectState & { filters: AdHocFilterItem[]; filterExpression: string };

function hasFilters(state: SceneObjectState): state is StateWithFilterExpression {
  if ('filterExpression' in state && typeof state.filterExpression === 'string') {
    return true;
  }
  return false;
}
