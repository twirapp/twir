import { Meta, Story } from '@storybook/vue3';

// import AvatarPng from '@/assets/avatar.png?url';
import TswDropdown from '@/components/TswDropdown/TswDropdown.vue';

export default {
  title: 'TswDropdown',
  component: TswDropdown,
} as Meta;

const Template: Story<void> = (args) => ({
  components: { TswDropdown },
  setup() {
    return { args };
  },
  template: `<TswDropdown>
              <template #button={onClick}>
                <button class="text-white-100" @click="onClick">Toggle</button>
              </template>
              <template #menu>
                <ul>
                  <li>aboba</li>
                  <li>aboba2</li>
                </ul>
              </template>
            </TswDropdown>`,
});

export const Dropdown = Template.bind({});
