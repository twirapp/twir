import { addons } from '@storybook/addons';
import { create } from '@storybook/theming';

export const theme = create({
  base: 'dark',
  appContentBg: '#151517',
  colorSecondary: '#816FEC',
  appBg: '#151517',
  inputBg: '#333335',
  appBorderColor: '#3F3F44',
  fontBase: '"Golos Text", sans-serif',
  inputBorder: '#3F3F42',
  barBg: '#333335',
  barSelectedColor: '#9989F6',
  textColor: 'white',
  barTextColor: '#B0B0B0',
  textMutedColor: '#B0B0B0',
  brandUrl: 'https://tsuwari.tk',
  brandTitle: 'Tsuwari Storybook',
  inputTextColor: '#F1F0F7',
});

addons.setConfig({
  theme,
});
