const withPWA = require('next-pwa')({
  disable: false,
  dest: 'public',
  register: true,
  skipWaiting: true,
});

const { i18n } = require('./next-i18next.config');

/** @type {import('next').NextConfig} */
const nextConfig = withPWA({
  reactStrictMode: false,
  swcMinify: true,
  basePath: '/dashboard',
  // trailingSlash: true,
  i18n,
  webpack: (config, { isServer }) => {
    config.module.rules.push({
      test: /\.svg$/,
      use: ['@svgr/webpack'],
    });

    if (!isServer) {
      config.resolve.fallback.fs = false;
    }

    return config;
  },
});

module.exports = nextConfig;
