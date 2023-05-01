import React from 'react';

import { Page } from 'app/core/components/Page/Page';
import { GrafanaRouteComponentProps } from 'app/core/navigation/types';

import { EditDataSource } from '../components/EditDataSource';
import { EditDataSourceActions } from '../components/EditDataSourceActions';
import { EditDataSourceSubtitle } from '../components/EditDataSourceSubtitle';
import { EditDataSourceTitle } from '../components/EditDataSourceTitle';
import { useDataSourceSettingsNav } from '../state';

export interface Props extends GrafanaRouteComponentProps<{ uid: string }> {}

export function EditDataSourcePage(props: Props) {
  const uid = props.match.params.uid;
  const params = new URLSearchParams(props.location.search);
  const pageId = params.get('page');
  const nav = useDataSourceSettingsNav(uid, pageId);

  return (
    <Page
      navId="datasources"
      pageNav={nav.main}
      renderTitle={(title) => (
        <EditDataSourceTitle
          dataSource={nav.dataSource}
          title={title}
          readOnly={nav.dataSource.readOnly}
          onUpdate={nav.dataSourceHeader.onUpdate}
          onNameChange={nav.dataSourceHeader.onNameChange}
        />
      )}
      subTitle={
        <EditDataSourceSubtitle
          dataSourcePluginName={nav.main.dataSourcePluginName}
          isDefault={nav.dataSource.isDefault || false}
          alertingSupported={nav.dataSourceHeader.alertingSupported}
          onDefaultChange={nav.dataSourceHeader.onDefaultChange}
        />
      }
      actions={<EditDataSourceActions uid={uid} />}
    >
      <Page.Contents>
        <EditDataSource uid={uid} pageId={pageId} />
      </Page.Contents>
    </Page>
  );
}

export default EditDataSourcePage;
