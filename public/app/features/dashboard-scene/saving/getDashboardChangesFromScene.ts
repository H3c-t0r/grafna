import { DashboardScene } from '../scene/DashboardScene';
import { transformSceneToSaveModel } from '../serialization/transformSceneToSaveModel';

import { getDashboardChanges as getDashboardSaveModelChanges } from './getDashboardChanges';

/**
 * Get changes between the initial save model and the current scene.
 * It also checks if the folder has changed.
 * @param scene DashboardScene object
 * @param saveTimeRange if true, compare the time range
 * @param saveVariables if true, compare the variables
 * @returns
 */
export function getDashboardChangesFromScene(scene: DashboardScene, saveTimeRange?: boolean, saveVariables?: boolean) {
  const changeInfo = getDashboardSaveModelChanges(
    scene.getInitialSaveModel()!,
    transformSceneToSaveModel(scene),
    saveTimeRange,
    saveVariables
  );
  const hasFolderChanges = scene.getInitialState()?.meta.folderUid !== scene.state.meta.folderUid;

  return {
    ...changeInfo,
    hasFolderChanges,
    hasChanges: changeInfo.hasChanges || hasFolderChanges,
  };
}
