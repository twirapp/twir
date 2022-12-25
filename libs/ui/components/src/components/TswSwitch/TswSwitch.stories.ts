import { Meta, Story } from '@storybook/vue3';

import TswSwitch from '@/components/TswSwitch/TswSwitch.vue';

export default {
  title: 'TswSwitch',
  component: TswSwitch,
} as Meta;

const Template: Story<{
  isChecked: boolean;
}> = (args) => ({
  components: { TswSwitch },
  setup() {
    return { args };
  },
  data: () => ({ isChecked: args.isChecked }),
  template: '<TswSwitch v-model:isChecked="isChecked" />',
});

export const Switch = Template.bind({});
Switch.args = {
  isChecked: false,
};
