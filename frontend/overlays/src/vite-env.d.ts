/// <reference types="vite/client" />

declare global {
  interface Window {
    webkitAudioContext: typeof AudioContext
  }
}

/// <reference types="vite/client" />
/// <reference types="vite-svg-loader" />
declare module '*.vue';
