import React, { useState } from 'react';

import { GenAIButton } from './GenAIButton';

export const GenAIPanelDescription = () => {
  const [isDescriptionGenerating, setIsDescriptionGenerating] = useState(false);
  const [genAIDescriptionResult, setGenAIDescriptionResult] = useState('');

  const onGenAIButtonClick = () => {
    setIsDescriptionGenerating(true);
    console.log('generating... ', genAIDescriptionResult);
    setTimeout(() => setIsDescriptionGenerating(false), 3000);
  };

  return (
    <GenAIButton
      text={isDescriptionGenerating ? 'Generating description' : 'Generate description'}
      onClick={onGenAIButtonClick}
      loading={isDescriptionGenerating}
    />
  );
};
