import React, { useState, useEffect } from 'react';
import { connect, ConnectedProps } from 'react-redux';

import { ConfirmModal } from '@grafana/ui';
import { StoreState, useDispatch } from 'app/types';

import { TutorialOverlay } from './TutorialOverlay';
import { setCurrentStep } from './slice';

const TutorialProviderComponent = ({
  availableTutorials,
  currentStep,
  currentTutorial,
}: ConnectedProps<typeof connector>) => {
  const dispatch = useDispatch();
  const [showExitTutorialModal, setShowExitTutorialModal] = useState(false);
  const currentTutorialSteps = availableTutorials.find((t) => t.id === currentTutorial)?.steps;
  const step = currentStep !== null && currentTutorialSteps && currentTutorialSteps[currentStep];

  useEffect(() => {
    const handler = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setShowExitTutorialModal(true);
      }
    };

    // TODO: why doesn't this work on keydown?
    window.addEventListener('keyup', handler);

    return () => {
      window.removeEventListener('keyup', handler);
    };
  }, []);

  if (step) {
    return (
      <>
        <TutorialOverlay currentStep={currentStep} step={step} />
        <ConfirmModal
          confirmText="Stop tutorial"
          onDismiss={() => {}}
          isOpen={showExitTutorialModal}
          title={`Exit tutorial`}
          body={`Do you want to stop the tutorial?`}
          onConfirm={() =>
            new Promise((resolve) => {
              dispatch(setCurrentStep(null));
              resolve();
              setShowExitTutorialModal(false);
            })
          }
        />
      </>
    );
  }

  return null;
};

const mapStateToProps = (state: StoreState) => {
  return {
    ...state.tutorials,
  };
};

const connector = connect(mapStateToProps);

export const TutorialProvider = connector(TutorialProviderComponent);
