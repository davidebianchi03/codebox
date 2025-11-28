import { configureStore } from '@reduxjs/toolkit';
import userReducer from "./slices/user";
import settingsReducer from "./slices/settings";

const store = configureStore({
    reducer: {
        user: userReducer,
        settings: settingsReducer,
    },
});

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
export default store;