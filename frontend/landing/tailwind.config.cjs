const tsuwariTheme = require('@twir/ui-theme');

module.exports = {
  content: ['./src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    ...tsuwariTheme,
    extends: {},
  },
  plugins: [require('tailwindcss-animation-delay')],
};
