import React, { useEffect, useState, ChangeEvent } from 'react';
import {
  TabContent,
  Label,
  Button,
  Select,
  Input,
  Spinner,
  TabsBar,
  Tab,
  StringValueEditor,
  useTheme2,
  stylesFactory,
} from '@grafana/ui';
import { GrafanaTheme2, SelectableValue } from '../../../../../../packages/grafana-data/src';
import { getBackendSrv } from '@grafana/runtime';
import Cards from './Cards';
import SVG from 'react-inlinesvg';
import { css } from '@emotion/css';

interface Props {
  value: string; //img/icons/unicons/0-plus.svg
  onChange: (value: string) => void;
  mediaType: 'icon' | 'image';
}

function ResourcePicker(props: Props) {
  const { value, onChange, mediaType } = props;
  const folders = (mediaType === 'icon' ? ['img/icons/unicons', 'img/icons/iot'] : ['img/bg']).map((v) => ({
    label: v,
    value: v,
  }));
  const folderOfCurrentValue = folders.filter((folder) => value.indexOf(folder.value) > -1)[0];
  const [currentFolder, setCurrentFolder] = useState<SelectableValue<string>>(folderOfCurrentValue);
  const [tabs, setTabs] = useState([
    { label: 'Select', active: true },
    // { label: 'Upload', active: false },
  ]);
  const [directoryIndex, setDirectoryIndex] = useState<SelectableValue[]>([]);
  const [defaultList, setDefaultList] = useState<SelectableValue[]>([]);
  const theme = useTheme2();
  const styles = getStyles(theme, theme.colors.text.primary);

  useEffect(() => {
    // we don't want to load everything before picking a folder
    if (currentFolder) {
      getBackendSrv()
        .get(`public/${currentFolder?.value}/index.json`)
        .then((data) => {
          const cards = data.files.map((v: string) => ({
            value: v,
            label: v,
            imgUrl: `public/${currentFolder?.value}/${v}`,
          }));
          setDirectoryIndex(cards);
          setDefaultList(cards);
        })
        .catch((e) => console.error(e));
    } else {
      return;
    }
  }, [currentFolder]);

  const onChangeSearch = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.value) {
      const filtered = directoryIndex.filter((card) =>
        card.value
          // exclude file type (.svg) in the search
          .substr(0, card.value.length - 4)
          .toLocaleLowerCase()
          .includes(e.target.value.toLocaleLowerCase())
      );
      setDirectoryIndex(filtered);
    } else {
      setDirectoryIndex(defaultList);
    }
  };
  const imgSrc = value.indexOf(':/') > 0 ? value : 'public/' + value;

  return (
    <>
      <Label>Current Item</Label>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', columnGap: '2px' }}>
        {mediaType === 'icon' && <SVG src={imgSrc} width="40" height="40" fill={theme.colors.text.primary} />}
        {mediaType === 'image' && <img src={imgSrc} width="40" height="40" />}
        <StringValueEditor value={value} onChange={onChange} item={{} as any} context={{} as any} />
        <Button variant="secondary">Apply</Button>
      </div>
      <TabsBar>
        {tabs.map((tab, index) => (
          <Tab
            label={tab.label}
            key={index}
            active={tab.active}
            onChangeTab={() => setTabs(tabs.map((tab, idx) => ({ ...tab, active: idx === index })))}
          />
        ))}
      </TabsBar>
      <TabContent>
        {tabs[0].active && (
          <div className={styles.tabContent}>
            <Select options={folders} onChange={setCurrentFolder} value={currentFolder} />
            <Input placeholder="Search" onChange={onChangeSearch} />
            {directoryIndex ? (
              <Cards cards={directoryIndex} onChange={onChange} currentFolder={currentFolder} />
            ) : (
              <Spinner />
            )}
          </div>
        )}
        {/* TODO: add file upload
          {tabs[1].active && (
          <FileUpload
            onFileUpload={({ currentTarget }) => console.log('file', currentTarget?.files && currentTarget.files[0])}
            className={styles.tabContent}
          />
        )} */}
      </TabContent>
    </>
  );
}

const getStyles = stylesFactory((theme: GrafanaTheme2, color) => {
  return {
    tabContent: css`
      margin-top: 20px;
      & > :nth-child(2) {
        margin-top: 10px;
      },
    `,
  };
});
export default ResourcePicker;
