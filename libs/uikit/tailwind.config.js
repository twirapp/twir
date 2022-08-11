module.exports = {
  content: ['./src/components/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    fontFamily: {
      sans: ['Golos', 'sans-serif'],
    },
    colors: {
      black: {
        100: '#0F0F0F',
        90: '#212123',
        80: '#343438',
        75: '#3F3F44',
      },
      gray: {
        70: '#515159',
        60: '#6E6F7B',
        50: '#878998',
        40: '#9EA2B3',
        30: '#BEC0CD',
      },
      white: {
        20: '#D8DAE5',
        10: '#EBECF2',
        0: '#FFFFFF',
      },
      purple: {
        45: '#5E3FD6',
        40: '#6946ED',
        35: '#7C5DF4',
      },
    },
    extends: {},
  },
  plugins: [],
};
