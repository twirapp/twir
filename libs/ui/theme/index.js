/* eslint-disable no-undef */

/** @type { import('tailwindcss/defaultTheme') } */
module.exports = {
  fontFamily: require('@tsuwari/ui-fonts/tailwind'),
  container: {
    center: true,
    screens: {
      xs: '360px',
      sm: '568px',
      md: '768px',
      lg: '996px',
      xl: '1200px',
      '2xl': '1480px',
    },
  },
  colors: {
    black: {
      0: '#000000',
      10: '#151517',
      15: '#212123',
      17: '#272729',
      20: '#333335',
      25: '#3F3F42',
    },
    gray: {
      35: '#565658',
      50: '#7A7A7B',
      60: '#959596',
      70: '#B1B1B2',
    },
    white: {
      95: '#F3F3F3',
      100: '#FFFFFF',
    },
    purple: {
      95: '#F1F0F7',
      80: '#ADA2F4',
      70: '#816FEC',
      60: '#644EE8',
      55: '#5C47D8',
      25: '#332D59',
    },
    green: {
      45: '#24BCA0',
    },
    blue: {
      60: '#3284FF',
    },
    yellow: {
      65: '#F2AC59',
    },
    red: {
      60: '#ED3F4A',
      65: '#EF555E',
    },
  },
};
