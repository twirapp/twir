const tsuwariTheme = require('@tsuwari/ui-theme');

module.exports = {
  content: ['./src/**/*.{vue,js,ts,jsx,tsx,astro}'],
  theme: {
    ...tsuwariTheme,
    extends: {},
  },
  plugins: [],
};
