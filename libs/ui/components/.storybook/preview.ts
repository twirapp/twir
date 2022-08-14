import '@tsuwari/ui-fonts';
import '../src/index.css';

import { theme } from './manager.js';

export const parameters = {
  docs: {
    theme,
  },
  actions: { argTypesRegex: '^on.*' },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
  backgrounds: {
    default: 'dark',
    values: [
      {
        name: 'dark',
        value: '#0f0f0f',
      },
    ],
  },
};
