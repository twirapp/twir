module.exports = {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
    './node_modules/tw-elements/dist/js/**/*.js',
    // './node_modules/@tsuwari/uikit/dist/**/*.js',
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('tw-elements/dist/plugin'),
  ],
  variants: {
    scrollbar: ['rounded'],
  },
};
