import type { Message } from './types.js';

export function normalizeDisplayName(userName: string, displayName: string) {
	if (userName === displayName.toLocaleLowerCase()) {
		return displayName;
	} else {
		return userName;
	}
}

export function getMessageAlign(direction: string): 'stretch' | 'center' {
	switch (direction) {
		case 'left':
		case 'right':
			return 'center';
		case 'top':
		case 'bottom':
			return 'stretch';
		default:
			return 'stretch';
	}
}

const DEFAULT_COLOR = '#a65ee8';
export function getColorFromMsg(msg: Message): string {
	if (msg.senderColor) {
		return msg.senderColor;
	}

	return DEFAULT_COLOR;
}
