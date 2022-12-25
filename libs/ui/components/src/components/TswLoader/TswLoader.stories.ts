import { Meta, Story } from '@storybook/vue3';

import { LoaderSize } from '@/components/TswLoader/props.types';
import TswLoader from '@/components/TswLoader/TswLoader.vue';

export default {
  title: 'TswLoader',
  component: TswLoader,
  argTypes: {
    size: {
      type: {
        name: 'enum',
        value: ['lg', 'md', 'sm'] as LoaderSize[],
      },
      control: 'select',
    },
  },
} as Meta;

interface TswLoaderProps {
  size: LoaderSize;
}

const Template: Story<TswLoaderProps> = (args) => ({
  components: { TswLoader },
  setup() {
    return { args };
  },
  template: '<TswLoader v-bind="args" />',
});

export const Loader = Template.bind({});
Loader.args = { size: 'lg' };
