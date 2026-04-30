import type { ReactNode } from "react";
import { Providers } from "@/components/providers";

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body style={{ margin: 0, fontFamily: "Inter, Arial, sans-serif", background: "#f5f7fb" }}>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
