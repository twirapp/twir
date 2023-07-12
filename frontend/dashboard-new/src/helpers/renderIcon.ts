import { SVGProps } from '@tabler/icons-vue';
import { FunctionalComponent, h } from 'vue';

export const renderIcon = (icon: (props: SVGProps) => FunctionalComponent<SVGProps>) => {
	return () => h(icon, null, { default: () => h(icon) });
};
