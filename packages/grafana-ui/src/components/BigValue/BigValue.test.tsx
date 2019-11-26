import React from 'react';
import { shallow } from 'enzyme';
import { BigValue, Props, BigValueColorMode, BigValueGraphMode } from './BigValue';
import { getTheme } from '../../themes/index';

jest.mock('jquery', () => ({
  plot: jest.fn(),
}));

const setup = (propOverrides?: object) => {
  const props: Props = {
    height: 300,
    width: 300,
    colorMode: BigValueColorMode.Value,
    graphMode: BigValueGraphMode.Line,
    value: {
      text: '25',
      numeric: 25,
    },
    theme: getTheme(),
  };

  Object.assign(props, propOverrides);

  const wrapper = shallow(<BigValue {...props} />);
  const instance = wrapper.instance() as BigValue;

  return {
    instance,
    wrapper,
  };
};

describe('Render SingleStat with basic options', () => {
  it('should render', () => {
    const { wrapper } = setup();
    expect(wrapper).toBeDefined();
    // expect(wrapper).toMatchSnapshot();
  });
});
