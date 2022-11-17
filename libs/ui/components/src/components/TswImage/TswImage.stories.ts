import { Meta, Story } from '@storybook/vue3';

import TswImage from '@/components/TswImage/TswImage.vue';

export default {
  title: 'TswImage',
  component: TswImage,
} as Meta;

const Template: Story<{
  src: string;
  width?: number;
  height?: number;
  lazy?: boolean;
  renderType: 'bg-image' | 'image';
}> = (args) => ({
  components: { TswImage },
  setup() {
    return { args };
  },
  template: '<TswImage v-bind="args" />',
});

export const Image = Template.bind({});
Image.args = {
  src: 'https://ichef.bbci.co.uk/news/976/cpsprodpb/B911/production/_120877374_gettyimages-1232608943.jpg',
  width: 976,
  height: 549,
  lazy: true,
  renderType: 'bg-image',
};
