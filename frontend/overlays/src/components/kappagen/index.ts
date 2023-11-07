export interface Emote {
	url: string;
	// zero width emote
	zwe?: Omit<Emote, 'zwe'>[];
}

interface AnimationParams {
	style: string;
	prefs: Record<string, any>;
	count: number;
}

declare global {
	interface Window {
		startup: () => void;
		emote: {
			addToShowList: (emote: Emote[]) => void;
			showEmotes: () => void;
		};
		kappagen: {
			show: (emotes: Emote[], params: AnimationParams) => Promise<void>;
		};
		random: (num: number) => number;
	}
}

window.addEventListener('load', window.startup);
