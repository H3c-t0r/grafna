import { useEffect } from 'react';
import { Subject } from 'rxjs';
import { useForceUpdate } from '@grafana/ui';
export class StateManagerBase {
    constructor(state) {
        this._subject = new Subject();
        this._state = state;
    }
    useState() {
        // eslint-disable-next-line react-hooks/rules-of-hooks
        return useLatestState(this);
    }
    get state() {
        return this._state;
    }
    setState(update) {
        this._state = Object.assign(Object.assign({}, this._state), update);
        this._subject.next(this._state);
    }
    /**
     * Subscribe to the scene state subject
     **/
    subscribeToState(observerOrNext) {
        return this._subject.subscribe(observerOrNext);
    }
}
/**
 * This hook is always returning model.state instead of a useState that remembers the last state emitted on the subject
 * The reason for this is so that if the model instance change this function will always return the latest state.
 */
function useLatestState(model) {
    const forceUpdate = useForceUpdate();
    useEffect(() => {
        const s = model.subscribeToState({ next: forceUpdate });
        return () => s.unsubscribe();
    }, [model, forceUpdate]);
    return model.state;
}
//# sourceMappingURL=StateManagerBase.js.map