import type { NextConfig } from "next";

const apiProxyTarget = process.env.SPLATMAKER_API_PROXY_TARGET ?? "http://localhost:8080";

const nextConfig: NextConfig = {
  reactStrictMode: true,
  async rewrites() {
    return [
      {
        source: "/v1/:path*",
        destination: `${apiProxyTarget}/v1/:path*`,
      },
      {
        source: "/healthz",
        destination: `${apiProxyTarget}/healthz`,
      },
    ];
  },
};

export default nextConfig;
