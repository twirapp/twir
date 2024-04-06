import { defineStore } from 'pinia';
import { ref } from 'vue';

const badgeUrl = (name: string) => `http://localhost:3005/twir-${name}.png`;

export const useBadges = defineStore('admin-panel/badges', () => {
	const badges = ref([
		{
			name: 'Twir Owner',
			image: badgeUrl('owner'),
		},
		{
			name: 'Twir Contributor',
			image: badgeUrl('contributor'),
		},
		{
			name: 'Twir Translator',
			image: badgeUrl('translator'),
		},
	]);

	return {
		badges,
	};
});
