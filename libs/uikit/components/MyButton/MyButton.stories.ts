import MyButton from '@/components/MyButton/MyButton.vue';

export default {
  title: 'MyButton',
  component: MyButton,
};

export const Primary = () => ({
  components: { MyButton },
  template: '<MyButton />',
});