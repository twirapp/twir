import { Meta, Story } from '@storybook/vue3';

// import AvatarPng from '@/assets/avatar.png?url';
import TswAvatar from '@/components/TswAvatar/TswAvatar.vue';

export default {
  title: 'TswAvatar',
  component: TswAvatar,
} as Meta;

const Template: Story<{
  src: string;
  href?: string;
  size?: number;
  lazy?: boolean;
}> = (args) => ({
  components: { TswAvatar },
  setup() {
    return { args };
  },
  template: '<TswAvatar v-bind="args" />',
});

export const Avatar = Template.bind({});
Avatar.args = {
  src:
    'https://miuipolska.pl/forum/uploads/monthly_2016_01/1988433-steam_avatar_advtime.jpg.f53695d43d11e9002aab38c2cac6cd8c.thumb.jpg.2b12f6b3c7311b157159b07225be12d6.jpg',
};
