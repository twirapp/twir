// eslint-disable-next-line @typescript-eslint/no-var-requires
const TsuwariTheme = require('@tsuwari/ui-theme');

/** @type {import('tailwindcss').Config} **/
module.exports = {
  content: ['./src/components/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    ...TsuwariTheme,
    extend: {},
  },
  plugins: [],
};
