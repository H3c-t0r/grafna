import React, { FC } from 'react';

import { useStyles } from '@grafana/ui';

import { SecretToggler } from '../SecretToggler';

import { Messages } from './KeysBlock.messages';
import { getStyles } from './KeysBlock.styles';
import { KeysBlockProps } from './KeysBlock.types';

export const KeysBlock: FC<KeysBlockProps> = ({ accessKey, secretKey }) => {
  const styles = useStyles(getStyles);

  return (
    <div className={styles.keysWrapper}>
      <div data-qa="access-key">
        <span className={styles.keyLabel}>{Messages.accessKey}</span>
        {accessKey}
      </div>
      <div data-qa="secret-key">
        <span className={styles.keyLabel}>{Messages.secretKey}</span>
        <span className={styles.secretTogglerWrapper}>
          <SecretToggler small secret={secretKey} />
        </span>
      </div>
    </div>
  );
};
