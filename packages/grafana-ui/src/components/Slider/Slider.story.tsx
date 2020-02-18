import React from 'react';
import { Slider } from './Slider';
import { select, number, boolean } from '@storybook/addon-knobs';

const getStory = (title: string, component: any) => ({
  title,
  parameters: {
    component,
  },
});

export default getStory('General/Slider', Slider);

const getKnobs = () => {
  return {
    min: number('min', 0),
    max: number('max', 100),
    lengthOfSlider: number('lengthOfSlider', 200),
    orientation: select('orientation', ['horizontal', 'vertical'], 'vertical'),
    reverse: boolean('reverse', true),
  };
};

const SliderWrapper = () => {
  const { min, max, orientation, reverse, lengthOfSlider } = getKnobs();
  return <Slider min={min} max={max} orientation={orientation} lengthOfSlider={lengthOfSlider} reverse={reverse} />;
};

export const basic = () => <SliderWrapper />;
