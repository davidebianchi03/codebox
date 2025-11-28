import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { InstanceSettings } from '../../types/settings';
import { AppDispatch } from '../store';
import { RetrieveInstanceSettings } from '../../api/common';

const settingsInitialState: InstanceSettings = {
    version: "",
    use_subdomains: false,
    external_url: "",
    wildcard_domain: "",
    recommended_runner_version: "",
}

export const settingsSlice = createSlice({
    name: "settings",
    initialState: settingsInitialState,
    reducers: {
        setSettings: (_, action: PayloadAction<InstanceSettings>) => {
            return action.payload;
        }
    }
});

export const { setSettings } = settingsSlice.actions

export default settingsSlice.reducer

export const FetchSettings = () => async (dispatch: AppDispatch) => {
    const r = await RetrieveInstanceSettings();
    if (r) {
        dispatch(setSettings(r));
    }
}
