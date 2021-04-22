import { Observable, of } from 'rxjs';
import { cloneDeep } from 'lodash';
import { AnnotationEvent, AnnotationQuery } from '@grafana/data';

import { DashboardQueryRunnerOptions, DashboardQueryRunnerWorker, DashboardQueryRunnerWorkerResult } from './types';
import { emptyResult } from './utils';

export class SnapshotWorker implements DashboardQueryRunnerWorker {
  canWork({ dashboard }: DashboardQueryRunnerOptions): boolean {
    return dashboard.annotations.list.some((a) => a.enable && Boolean(a.snapshotData));
  }

  work(options: DashboardQueryRunnerOptions): Observable<DashboardQueryRunnerWorkerResult> {
    if (!this.canWork(options)) {
      return emptyResult();
    }

    const annotations = this.getAnnotationsFromSnapshot(options);
    return of({ annotations, alertStates: [] });
  }

  getAnnotationsFromSnapshot(options: DashboardQueryRunnerOptions): AnnotationEvent[] {
    const { dashboard } = options;
    const dashAnnotations = dashboard.annotations.list.filter((a) => a.enable);
    const snapshots = dashAnnotations.filter((a) => Boolean(a.snapshotData));
    const annotations = snapshots.reduce((acc, curr) => {
      return acc.concat(SnapshotWorker.translateQueryResult(curr, curr.snapshotData));
    }, [] as AnnotationEvent[]);

    return annotations;
  }

  private static translateQueryResult(annotation: AnnotationQuery, results: AnnotationEvent[]): AnnotationEvent[] {
    annotation = cloneDeep(annotation);
    delete annotation.snapshotData;

    for (const item of results) {
      item.source = annotation;
      item.color = annotation.iconColor;
      item.type = annotation.name;
      item.isRegion = Boolean(item.timeEnd && item.time !== item.timeEnd);
    }

    return results;
  }
}
