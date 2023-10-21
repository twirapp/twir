import { IconAbc, IconAsteriskSimple, IconLanguageOff, IconLinkOff, IconListLetters, IconMessageOff, IconMoodOff, type SVGProps } from '@tabler/icons-vue';
import type { Item } from '@twir/grpc/generated/api/api/moderation';
import type { FunctionalComponent } from 'vue';

export const availableSettingsMapping: Record<string, Item & {
	icon: FunctionalComponent<SVGProps>
}> = {
	links: {
		icon: IconLinkOff,
	},
	language: {
		icon: IconLanguageOff,
	},
	deny_list: {
		icon: IconListLetters,
	},
	long_message: {
		icon: IconMessageOff,
	},
	caps: {
		icon: IconAbc,
	},
	emotes: {
		icon: IconMoodOff,
	},
	symbols: {
		icons: IconAsteriskSimple,
	},
};
