import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { CurrentUser } from '../../types/user';

const userInitialState: CurrentUser = {
    email: "",
    first_name: "",
    last_name: "",
    is_superuser: false,
    is_template_manager: false,
    impersonated: false,
    last_login: "",
    created_at: "",
}

export const userSlice = createSlice({
    name: "user",
    initialState: userInitialState,
    reducers: {
        setUser: (_, action: PayloadAction<CurrentUser>) => {
            return action.payload;
        }
    }
});

export const { setUser } = userSlice.actions

export default userSlice.reducer
