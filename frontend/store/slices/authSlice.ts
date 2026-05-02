import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

type AuthState = {
  token: string | null;
};

const initialState: AuthState = {
  token: null,
};

const authSlice = createSlice({
  name: "auth",
  initialState,
  reducers: {
    setAuth(state, action: PayloadAction<{ token: string }>) {
      state.token = action.payload.token;
    },
    logout(state) {
      state.token = null;
    },
  },
});

export const { setAuth, logout } = authSlice.actions;
export default authSlice.reducer;
