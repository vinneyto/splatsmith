import { configureStore } from "@reduxjs/toolkit";
import authReducer from "@/store/slices/authSlice";
import { splatmakerApi } from "@/store/api/splatmakerApi";

export const store = configureStore({
  reducer: {
    auth: authReducer,
    [splatmakerApi.reducerPath]: splatmakerApi.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(splatmakerApi.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
