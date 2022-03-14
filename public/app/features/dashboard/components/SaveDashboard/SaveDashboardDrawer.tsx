import React, { useMemo, useState } from 'react';
import { css } from '@emotion/css';
import {
  Button,
  Checkbox,
  CustomScrollbar,
  Drawer,
  Field,
  HorizontalGroup,
  Input,
  Spinner,
  Tab,
  TabContent,
  TabsBar,
  TextArea,
  useStyles2,
} from '@grafana/ui';
import { useForm } from 'react-hook-form';
import { SaveDashboardModalProps } from './types';
import { GrafanaTheme2 } from '@grafana/data';
import { DiffViewer } from '../VersionHistory/DiffViewer';
import { jsonDiff } from '../VersionHistory/utils';
import { DiffGroup } from '../VersionHistory/DiffGroup';
import { selectors } from '@grafana/e2e-selectors';
import { useAsync } from 'react-use';
import { backendSrv } from 'app/core/services/backend_srv';
import { useDashboardSave } from './useDashboardSave';

interface FormDTO {
  message?: string; // the commit message

  newDashboardTitle?: string;
}

export const SaveDashboardDrawer = ({ dashboard, onDismiss }: SaveDashboardModalProps) => {
  const styles = useStyles2(getStyles);
  const {
    handleSubmit,
    register,
    formState: { errors },
  } = useForm<FormDTO>({ defaultValues: {} });

  const hasTimeChanged = useMemo(() => dashboard.hasTimeChanged(), [dashboard]);
  const hasVariableChanged = useMemo(() => dashboard.hasVariableValuesChanged(), [dashboard]);

  const [saveTimeRange, setSaveTimeRange] = useState(false);
  const [saveVariables, setSaveVariables] = useState(false);

  const status = useMemo(() => {
    const isProvisioned = dashboard.meta.provisioned;
    const isNew = dashboard.version === 0 && !dashboard.uid;
    const isChanged = dashboard.version > 0;

    return {
      isProvisioned,
      isNew,
      isChanged,
    };
  }, [dashboard]);

  // Fetch the prevous version rather than depend on redux store
  const previous = useAsync(async () => {
    if (status.isNew) {
      return undefined;
    }

    const result = await backendSrv.getDashboardByUid(dashboard.uid);
    return result.dashboard;
  }, [dashboard, status]);

  const data = useMemo(() => {
    const clone = dashboard.getSaveModelClone({
      saveTimerange: Boolean(saveTimeRange),
      saveVariables: Boolean(saveVariables),
    });

    if (!previous.value) {
      return { clone };
    }

    const cloneJSON = JSON.stringify(clone, null, 2);
    const cloneSafe = JSON.parse(cloneJSON); // avoids undefined issues

    const diff = jsonDiff(previous.value, cloneSafe);
    let diffCount = 0;
    for (const d of Object.values(diff)) {
      diffCount += d.length;
    }

    return {
      clone,
      cloneJSON, // for diff
      diff,
      diffCount,
      hasChanges: diffCount > 0,
    };
  }, [dashboard, previous.value, saveTimeRange, saveVariables]);

  const [showDiff, setShowDiff] = useState(false);
  const [saving, setSaving] = useState(false);

  // reuse current save path
  const { onDashboardSave } = useDashboardSave(dashboard);
  const doSave = async (dto: FormDTO) => {
    setSaving(true);

    const body = data.clone;
    if (status.isNew && dto.newDashboardTitle) {
      body.title = dto.newDashboardTitle;
    }

    const result = await onDashboardSave(
      body,
      {
        message: dto.message,
        saveTimerange: saveTimeRange,
        saveVariables: saveVariables,
      },
      dashboard
    );

    if (result.status === 'success') {
      if (saveVariables) {
        dashboard.resetOriginalVariables();
      }
      if (saveTimeRange) {
        dashboard.resetOriginalTime();
      }
    }
    onDismiss();
    setSaving(false);
  };

  const renderBody = () => {
    if (showDiff) {
      return (
        <>
          {previous.loading && <Spinner />}
          {!data.hasChanges && <div>No changes made to this dashboard</div>}
          {data.diff && data.hasChanges && (
            <div>
              <div className={styles.spacer}>
                {Object.entries(data.diff).map(([key, diffs]) => (
                  <DiffGroup diffs={diffs} key={key} title={key} />
                ))}
              </div>

              <h4>JSON Diff</h4>
              <DiffViewer oldValue={JSON.stringify(previous.value, null, 2)} newValue={data.cloneJSON!} />
            </div>
          )}
        </>
      );
    }
    if (status.isProvisioned) {
      return <div>TODO: the provisioning feedback message</div>;
    }

    return (
      <form onSubmit={handleSubmit(doSave)}>
        {status.isNew && (
          <>
            <Field label="Dashboard title" invalid={!!errors.newDashboardTitle} error="Title is required">
              <Input {...register('newDashboardTitle', { required: true })} placeholder="Set dashboard title" />
            </Field>
            <Field label="Dashboard folder">
              <div>TODO: folder selectoin</div>
            </Field>
          </>
        )}

        <Field label="Message" invalid={!!errors.message} error="Message is required">
          <TextArea
            {...register('message', { required: false })}
            rows={5}
            placeholder="Add a note to describe your changes."
          />
        </Field>

        <HorizontalGroup>
          <Button
            type="submit"
            aria-label="Save dashboard button"
            disabled={!data.hasChanges && !status.isNew}
            icon={saving ? 'fa fa-spinner' : undefined}
          >
            {saving ? '' : 'Save'}
          </Button>
          <Button type="button" variant="secondary" onClick={onDismiss} fill="outline">
            Cancel
          </Button>
        </HorizontalGroup>
        {!data.hasChanges && !status.isNew && <div className={styles.nothing}>No changes to save</div>}
      </form>
    );
  };

  return (
    <Drawer
      title={dashboard.title}
      onClose={onDismiss}
      width={'40%'}
      subtitle={
        <>
          {hasTimeChanged && (
            <Checkbox
              checked={saveTimeRange}
              onChange={() => setSaveTimeRange(!saveTimeRange)}
              label="Save current time range as dashboard default"
              aria-label={selectors.pages.SaveDashboardModal.saveTimerange}
            />
          )}
          {hasVariableChanged && (
            <Checkbox
              checked={saveVariables}
              onChange={() => setSaveVariables(!saveVariables)}
              label="Save current variable values as dashboard default"
              aria-label={selectors.pages.SaveDashboardModal.saveVariables}
            />
          )}
          <TabsBar className={styles.tabsBar}>
            <Tab label={'Save'} active={!showDiff} onChangeTab={() => setShowDiff(false)} />
            {data.hasChanges && (
              <Tab label={'Changes'} active={showDiff} onChangeTab={() => setShowDiff(true)} counter={data.diffCount} />
            )}
          </TabsBar>
        </>
      }
      expandable
    >
      <CustomScrollbar autoHeightMin="100%">
        <TabContent>{renderBody()}</TabContent>
      </CustomScrollbar>
    </Drawer>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  tabsBar: css`
    padding-left: ${theme.v1.spacing.md};
    margin: ${theme.v1.spacing.lg} -${theme.v1.spacing.sm} -${theme.v1.spacing.lg} -${theme.v1.spacing.lg};
  `,
  spacer: css`
    margin-bottom: ${theme.v1.spacing.xl};
  `,
  nothing: css`
    margin: ${theme.v1.spacing.sm};
    color: ${theme.colors.secondary.shade};
  `,
});
