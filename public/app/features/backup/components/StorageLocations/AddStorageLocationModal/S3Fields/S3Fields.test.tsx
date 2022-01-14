import { TextInputField } from '@percona/platform-core';
import { shallow } from 'enzyme';
import React from 'react';

import { SecretToggler } from '../../../SecretToggler';

import { S3Fields } from './S3Fields';

describe('S3Fields', () => {
  it('should pass initial values', () => {
    const wrapper = shallow(
      <S3Fields bucketName="bucket" endpoint="/foo" accessKey="accessKey" secretKey="secretKey" />
    );
    const inputs = wrapper.find(TextInputField);
    expect(inputs.find({ initialValue: '/foo' }).exists()).toBeTruthy();
    expect(inputs.find({ initialValue: 'accessKey' }).exists()).toBeTruthy();
    expect(inputs.find({ initialValue: 'bucket' }).exists()).toBeTruthy();
    expect(wrapper.find(SecretToggler).last().prop('secret')).toBe('secretKey');
  });
});
