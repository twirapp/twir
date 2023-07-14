import { SVGProps } from '@tabler/icons-vue';
import { FunctionalComponent, h } from 'vue';

type Opts = {
	color?: string
}

export const renderIcon = (icon: (props: SVGProps) => FunctionalComponent<SVGProps>, opts?: Opts) => {
	return () => h(icon, opts, { default: () => h(icon) });
};
