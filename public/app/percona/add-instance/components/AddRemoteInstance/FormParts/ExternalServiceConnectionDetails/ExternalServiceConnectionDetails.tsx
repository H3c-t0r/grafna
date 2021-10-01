import React, { FC, useCallback, useEffect, useMemo } from 'react';
import { PasswordInputField, RadioButtonGroupField, TextInputField, validators } from '@percona/platform-core';
import Validators from 'app/percona/shared/helpers/validators';
import { LinkTooltip } from 'app/percona/shared/components/Elements/LinkTooltip/LinkTooltip';
import { useTheme } from '@grafana/ui';
import { FormPartProps, MetricsParameters, Schema } from '../FormParts.types';
import { getStyles } from '../FormParts.styles';
import { Messages } from '../FormParts.messages';
import { metricsParametersOptions, schemaOptions } from '../FormParts.constants';

export const ExternalServiceConnectionDetails: FC<FormPartProps> = ({ form }) => {
  const theme = useTheme();
  const styles = getStyles(theme);
  const formValues = form.getState().values;
  const selectedOption = formValues?.metricsParameters;
  const urlValue = formValues?.url;
  const portValidators = useMemo(() => [validators.required, Validators.validatePort], []);

  const trim = useCallback((value) => (value ? value.trim() : value), []);
  const getUrlParts = () => {
    try {
      const url = new URL(form.getState().values.url);
      const protocol = url.protocol.replace(':', '');

      form.mutators?.setValue('schema', protocol);
      form.mutators?.setValue('address', url.hostname);
      form.mutators?.setValue('port', url.port || (protocol === 'https' ? '443' : '80'));
      form.mutators?.setValue('metrics_path', url.pathname);
      form.mutators?.setValue('username', url.username);
      form.mutators?.setValue('password', url.password);
    } catch (e) {
      form.mutators?.setValue('schema', Schema.HTTPS);
      form.mutators?.setValue('address', '');
      form.mutators?.setValue('port', '443');
      form.mutators?.setValue('metrics_path', '');
      form.mutators?.setValue('username', '');
      form.mutators?.setValue('password', '');
    }
  };

  useEffect(getUrlParts, [urlValue]);

  return (
    <div className={styles.groupWrapper}>
      <h4 className={styles.sectionHeader}>{Messages.form.titles.connectionDetails}</h4>
      <div className={styles.labelWrapper} data-testid="username-label">
        <span>{Messages.form.labels.externalService.serviceName}</span>
        <LinkTooltip tooltipText={Messages.form.tooltips.externalService.serviceName} icon="info-circle" />
      </div>
      <TextInputField name="serviceName" placeholder={Messages.form.placeholders.externalService.serviceName} />
      <div className={styles.labelWrapper} data-testid="username-label">
        <span>{Messages.form.labels.externalService.group}</span>
        <LinkTooltip tooltipText={Messages.form.tooltips.externalService.group} icon="info-circle" />
      </div>
      <TextInputField name="group" />
      <div className={styles.labelWrapper} data-testid="address-label">
        <span>{Messages.form.labels.externalService.connectionParameters}</span>
        <LinkTooltip tooltipText={Messages.form.tooltips.externalService.url} icon="info-circle" />
      </div>
      <RadioButtonGroupField
        name="metricsParameters"
        data-testid="metrics-parameters-field"
        options={metricsParametersOptions}
      />
      {selectedOption === MetricsParameters.parsed && (
        <>
          <div className={styles.labelWrapper} data-testid="address-label">
            <span>{Messages.form.labels.externalService.url}</span>
            <LinkTooltip tooltipText={Messages.form.tooltips.externalService.url} icon="info-circle" />
          </div>
          <div className={styles.urlFieldWrapper}>
            <TextInputField
              name="url"
              placeholder={Messages.form.placeholders.externalService.url}
              validators={[Validators.validateUrl, validators.required]}
            />
          </div>
        </>
      )}
      {selectedOption === MetricsParameters.manually && (
        <>
          <div className={styles.labelWrapper} data-testid="address-label">
            <span>{Messages.form.labels.externalService.schema}</span>
            <LinkTooltip tooltipText={Messages.form.tooltips.externalService.schema} icon="info-circle" />
          </div>
          <RadioButtonGroupField name="schema" data-testid="http-schema-field" options={schemaOptions} />
          <div className={styles.labelWrapper} data-testid="address-label">
            <span>{Messages.form.labels.externalService.address}</span>
            <LinkTooltip tooltipText={Messages.form.tooltips.externalService.address} icon="info-circle" />
          </div>
          <TextInputField
            name="address"
            placeholder={Messages.form.placeholders.externalService.address}
            validators={[validators.required]}
          />
          <div className={styles.labelWrapper} data-testid="service-name-label">
            <span>{Messages.form.labels.externalService.metricsPath}</span>
            <LinkTooltip tooltipText={Messages.form.tooltips.externalService.metricsPath} icon="info-circle" />
          </div>
          <TextInputField name="metrics_path" placeholder={Messages.form.placeholders.externalService.metricsPath} />
          <div className={styles.labelWrapper} data-testid="port-label">
            <span>{Messages.form.labels.externalService.port}</span>
            <LinkTooltip tooltipText={Messages.form.tooltips.externalService.port} icon="info-circle" />
          </div>
          <TextInputField name="port" placeholder="Port" validators={portValidators} />
          <div className={styles.labelWrapper} data-testid="username-label">
            <span>{Messages.form.labels.externalService.username}</span>
            <LinkTooltip tooltipText={Messages.form.tooltips.externalService.username} icon="info-circle" />
          </div>
          <TextInputField
            name="username"
            placeholder={Messages.form.placeholders.externalService.username}
            format={trim}
          />
          <div className={styles.labelWrapper} data-testid="password-label">
            <span>{Messages.form.labels.externalService.password}</span>
            <LinkTooltip tooltipText={Messages.form.tooltips.externalService.password} icon="info-circle" />
          </div>
          <PasswordInputField
            name="password"
            placeholder={Messages.form.placeholders.externalService.password}
            format={trim}
          />
        </>
      )}
    </div>
  );
};
