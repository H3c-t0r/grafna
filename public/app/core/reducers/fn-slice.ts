import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface FNState {
  FNDashboard: boolean;
  uid: string;
  slug: string;
  theme: string;
  controlsContainer: HTMLElement | undefined;
}

const initialState: FNState = {
  FNDashboard: false,
  uid: '',
  slug: '',
  theme: '',
  controlsContainer: undefined,
};

const fnSlice = createSlice({
  name: 'fnGlobleState',
  initialState,
  reducers: {
    setIntialMountState: (state, action: PayloadAction<FNState>) => {
      const { payload } = action.payload;
      return payload;
    },
    updateFnState: (state, action: PayloadAction<{ type: string; payload: string | boolean | HTMLElement }>) => {
      const { type, payload } = action.payload;
      return {
        ...state,
        [type]: payload,
      };
    },
  },
});

export const { updateFnState, setIntialMountState } = fnSlice.actions;
export const fnSliceReducer = fnSlice.reducer;
