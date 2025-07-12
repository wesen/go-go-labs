import { configureStore, createSlice } from '@reduxjs/toolkit';
import { widgetsApi } from './services/widgetsApi.ts';
import featuredWidgetReducer from './features/featuredWidget/featuredWidgetSlice.ts';

export const store = configureStore({
  reducer: {
    featuredWidget: featuredWidgetReducer,
    [widgetsApi.reducerPath]: widgetsApi.reducer,
  },
  middleware: getDefaultMiddleware => 
    getDefaultMiddleware().concat(widgetsApi.middleware),
});
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch; 