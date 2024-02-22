type ReturnToPreviousHook = () => (title: string, href?: string) => void;

let rtpHook: ReturnToPreviousHook | undefined = undefined;

export const setReturnToPreviousHook = (hook: ReturnToPreviousHook) => {
  rtpHook = hook;
};

/**
 * Guidelines:
 * - Do only use the ‘Return to previous’ functionality when the user is sent to another context, such as from Alerting to a dashboard.
 * - Do specify a button title that identifies the page to return to in the most understandable way. Do not use text such as ‘Back to the previous page’. Be specific.
 *
 * {@link https://docs.google.com/document/d/1bSh5W_v5okEwjVPvP4bOmILwfKwmYLkf6Y752Ay_02I ReturnToPrevious: Documentation for devs}
 */
export const useReturnToPrevious: ReturnToPreviousHook = () => {
  if (!rtpHook) {
    if (process.env.NODE_ENV !== 'production') {
      throw new Error('useReturnToPrevious hook not found in @grafana/runtime');
    }
    return () => console.error('ReturnToPrevious hook not found');
  }

  return rtpHook();
};
