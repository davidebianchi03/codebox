import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { CurrentUser } from '../../types/user';

const userInitialState: CurrentUser | null = null;

export const userSlice = createSlice({
    name: "user",
    initialState: userInitialState as CurrentUser | null,
    reducers: {
        setUser: (_, action: PayloadAction<CurrentUser>) => {
            return action.payload;
        }
    }
});

export const { setUser } = userSlice.actions

export default userSlice.reducer
