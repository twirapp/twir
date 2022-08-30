const tsuwariTheme = require('@tsuwari/ui-theme');

module.exports = {
  content: ['./src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    ...tsuwariTheme,
    extends: {},
  },
  plugins: [],
};
