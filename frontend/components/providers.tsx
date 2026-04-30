"use client";

import type { ReactNode } from "react";
import { Provider } from "react-redux";
import { ThemeProvider } from "@/components/theme-provider";
import { store } from "@/store/store";

export function Providers({ children }: { children: ReactNode }) {
  return (
    <ThemeProvider>
      <Provider store={store}>{children}</Provider>
    </ThemeProvider>
  );
}
