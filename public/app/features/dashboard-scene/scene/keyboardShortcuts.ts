import { SetPanelAttentionEvent } from '@grafana/data';
import { locationService } from '@grafana/runtime';
import { sceneGraph, VizPanel } from '@grafana/scenes';
import appEvents from 'app/core/app_events';
import { KeybindingSet } from 'app/core/services/KeybindingSet';

import { ShareModal } from '../sharing/ShareModal';
import { dashboardSceneGraph } from '../utils/dashboardSceneGraph';
import { getEditPanelUrl, getInspectUrl, getViewPanelUrl, tryGetExploreUrlForPanel } from '../utils/urlBuilders';
import { getPanelIdForVizPanel } from '../utils/utils';

import { DashboardScene } from './DashboardScene';
import { onRemovePanel, toggleVizPanelLegend } from './PanelMenuBehavior';

export function setupKeyboardShortcuts(scene: DashboardScene) {
  const keybindings = new KeybindingSet();
  let vizPanelKey: string | null = null;

  const panelAttentionSubscription = appEvents.subscribe(SetPanelAttentionEvent, (event) => {
    console.log('panelattention');
    if (typeof event.payload.panelId === 'string') {
      vizPanelKey = event.payload.panelId;
    }
  });

  function withFocusedPanel(scene: DashboardScene, fn: (vizPanel: VizPanel) => void) {
    return () => {
      const vizPanel = sceneGraph.findObject(scene, (o) => o.state.key === vizPanelKey);
      if (vizPanel && vizPanel instanceof VizPanel) {
        fn(vizPanel);
        return;
      }
    };
  }

  // View panel
  keybindings.addBinding({
    key: 'v',
    onTrigger: withFocusedPanel(scene, (vizPanel: VizPanel) => {
      if (!scene.state.viewPanelScene) {
        locationService.push(getViewPanelUrl(vizPanel));
      }
    }),
  });

  // Panel edit
  keybindings.addBinding({
    key: 'e',
    onTrigger: withFocusedPanel(scene, async (vizPanel: VizPanel) => {
      const sceneRoot = vizPanel.getRoot();
      if (sceneRoot instanceof DashboardScene) {
        const panelId = getPanelIdForVizPanel(vizPanel);
        if (!scene.state.editPanel) {
          locationService.push(getEditPanelUrl(panelId));
        }
      }
    }),
  });

  // Panel share
  keybindings.addBinding({
    key: 'p s',
    onTrigger: withFocusedPanel(scene, async (vizPanel: VizPanel) => {
      scene.showModal(new ShareModal({ panelRef: vizPanel.getRef(), dashboardRef: scene.getRef() }));
    }),
  });

  // Panel inspect
  keybindings.addBinding({
    key: 'i',
    onTrigger: withFocusedPanel(scene, async (vizPanel: VizPanel) => {
      locationService.push(getInspectUrl(vizPanel));
    }),
  });

  // Got to Explore for panel
  keybindings.addBinding({
    key: 'p x',
    onTrigger: withFocusedPanel(scene, async (vizPanel: VizPanel) => {
      const url = await tryGetExploreUrlForPanel(vizPanel);
      if (url) {
        locationService.push(url);
      }
    }),
  });

  // Toggle legend
  keybindings.addBinding({
    key: 'p l',
    onTrigger: withFocusedPanel(scene, toggleVizPanelLegend),
  });

  // Refresh
  keybindings.addBinding({
    key: 'd r',
    onTrigger: () => sceneGraph.getTimeRange(scene).onRefresh(),
  });

  // Zoom out
  keybindings.addBinding({
    key: 't z',
    onTrigger: () => {
      handleZoomOut(scene);
    },
  });
  keybindings.addBinding({
    key: 'ctrl+z',
    onTrigger: () => {
      handleZoomOut(scene);
    },
  });

  keybindings.addBinding({
    key: 't left',
    onTrigger: () => {
      handleTimeRangeShift(scene, 'left');
    },
  });

  keybindings.addBinding({
    key: 't right',
    onTrigger: () => {
      handleTimeRangeShift(scene, 'right');
    },
  });

  // Dashboard settings
  keybindings.addBinding({
    key: 'd s',
    onTrigger: scene.onOpenSettings,
  });

  keybindings.addBinding({
    key: 'mod+s',
    onTrigger: () => scene.openSaveDrawer({}),
  });

  // toggle all panel legends (TODO)
  // delete panel
  keybindings.addBinding({
    key: 'p r',
    onTrigger: withFocusedPanel(scene, (vizPanel: VizPanel) => {
      if (scene.state.isEditing) {
        onRemovePanel(scene, vizPanel);
      }
    }),
  });

  // duplicate panel
  keybindings.addBinding({
    key: 'p d',
    onTrigger: withFocusedPanel(scene, (vizPanel: VizPanel) => {
      if (scene.state.isEditing) {
        scene.duplicatePanel(vizPanel);
      }
    }),
  });

  // toggle all exemplars (TODO)
  // collapse all rows (TODO)
  // expand all rows (TODO)

  return () => {
    keybindings.removeAll();
    panelAttentionSubscription.unsubscribe();
  };
}

function handleZoomOut(scene: DashboardScene) {
  const timePicker = dashboardSceneGraph.getTimePicker(scene);
  timePicker?.onZoom();
}

function handleTimeRangeShift(scene: DashboardScene, direction: 'left' | 'right') {
  const timePicker = dashboardSceneGraph.getTimePicker(scene);

  if (!timePicker) {
    return;
  }

  if (direction === 'left') {
    timePicker.onMoveBackward();
  }
  if (direction === 'right') {
    timePicker.onMoveForward();
  }
}
