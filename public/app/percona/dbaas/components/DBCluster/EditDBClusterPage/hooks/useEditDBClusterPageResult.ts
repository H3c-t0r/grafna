import { useSelector } from 'app/types';

import { getAddDbCluster, getUpdateDbCluster } from '../../../../../shared/core/selectors';
import { DBClusterPageMode } from '../EditDBClusterPage.types';

export const useEditDBClusterPageResult = (mode: DBClusterPageMode): [any] => {
  const { result } = useSelector(mode === 'create' ? getAddDbCluster : getUpdateDbCluster);
  return [result];
};
