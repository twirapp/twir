import { IconUserPlus } from '@tabler/icons-vue';
import { FunctionalComponent } from 'vue';

export const EVENTS: Record<string, {
	name: string,
	icon?: FunctionalComponent
}> = {
	'FOLLOW': {
		name: 'Follow',
		icon: IconUserPlus,
	},
};
