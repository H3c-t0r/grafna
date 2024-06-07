import { Observable, ReplaySubject, Subject, firstValueFrom, map, scan, startWith } from 'rxjs';

import { PluginPreloadResult } from '../../pluginPreloader';
import { deepFreeze } from '../utils';

type ConstructorOptions<T extends Record<string | symbol, unknown>> = {
  initialState: T;
};

// This is the base-class used by the separate specific registries.
export class Registry<T extends Record<string | symbol, unknown>> {
  private resultSubject: Subject<PluginPreloadResult>;
  private registrySubject: ReplaySubject<T>;

  constructor(options: ConstructorOptions<T>) {
    const { initialState } = options;
    this.resultSubject = new Subject<PluginPreloadResult>();
    // This is the subject that we expose.
    // (It will buffer the last value on the stream - the registry - and emit it to new subscribers immediately.)
    this.registrySubject = new ReplaySubject<T>(1);

    this.resultSubject
      .pipe(
        scan(this.mapToRegistry, initialState),
        // Emit an empty registry to start the stream (it is only going to do it once during construction, and then just passes down the values)
        startWith(initialState),
        map((registry) => deepFreeze(registry))
      )
      // Emitting the new registry to `this.registrySubject`
      .subscribe(this.registrySubject);
  }

  mapToRegistry(registry: T, item: PluginPreloadResult): T {
    return registry;
  }

  register(result: PluginPreloadResult): void {
    this.resultSubject.next(result);
  }

  asObservable(): Observable<T> {
    return this.registrySubject.asObservable();
  }

  getState(): Promise<T> {
    return firstValueFrom(this.asObservable());
  }
}
