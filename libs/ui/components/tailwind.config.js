// eslint-disable-next-line @typescript-eslint/no-var-requires
const TsuwariTheme = require('@twir/ui-theme');

/** @type {import('tailwindcss').Config} **/
module.exports = {
  content: ['./src/components/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    ...TsuwariTheme,
    extend: {},
  },
  plugins: [],
};
