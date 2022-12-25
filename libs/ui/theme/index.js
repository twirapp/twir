/* eslint-disable no-undef */

/** @type { import('tailwindcss/defaultTheme') } */
module.exports = {
  fontFamily: require('@tsuwari/ui-fonts/tailwind'),
  container: {
    center: true,
    screens: {
      xs: { max: '360px' },
      sm: { max: '568px' },
      md: { max: '768px' },
      lg: { max: '996px' },
      xl: { max: '1200px' },
      '2xl': { max: '1480px' },
    },
  },
  screens: {
    'max-xs': { max: '359.98px' },
    'min-xs': { min: '360px' },
    'max-sm': { max: '575.98px' },
    'min-sm': { min: '576px' },
    'max-md': { max: '767.98px' },
    'min-md': { min: '768px' },
    'max-lg': { max: '995.98px' },
    'min-lg': { min: '996px' },
    'max-xl': { max: '1199.98px' },
    'min-xl': { min: '1200px' },
    'max-2xl': { max: '1479.98px' },
    'min-2xl': { min: '1480px' },
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
      30: '#4A4A4F',
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
      50: '#513EBF',
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
  extend: {
    keyframes: {
      fadeInDown: {
        from: { opacity: '0', transform: 'translateY(-20px)' },
        to: { opacity: '1', transform: 'translateY(0)' },
      },
      fadeIn: {
        from: { opacity: '0' },
        to: { opacity: '1' },
      },
    },
    animation: {
      fadeInDown: 'fadeInDown 1s forwards',
      fadeIn: 'fadeIn 1s forwards',
      fadeInLong: 'fadeIn 2s forwards',
    },
  },
};
