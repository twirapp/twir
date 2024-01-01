import type { Message } from './types.js';

export function normalizeDisplayName(userName: string, displayName: string) {
	if (userName === displayName.toLocaleLowerCase()) {
		return displayName;
	} else {
		return userName;
	}
}

export type MessageAlignType = 'center' | 'baseline';

export function getMessageAlign(direction: string): MessageAlignType {
	switch (direction) {
		case 'left':
		case 'right':
			return 'center';
		case 'top':
		case 'bottom':
		default:
			return 'baseline';
	}
}

export type Direction = 'horizontal' | 'vertical';

export function getChatDirection(direction: string): Direction {
	switch (direction) {
		case 'left':
		case 'right':
			return 'horizontal';
		case 'top':
		case 'bottom':
		default:
			return 'vertical';
	}
}

const DEFAULT_COLOR = '#a65ee8';

export function getUserColor(color?: string): string {
	if (color) {
		return color;
	}

	return DEFAULT_COLOR;
}

export function getColorFromMsg(msg: Message): string {
	if (msg.senderColor) {
		return msg.senderColor;
	}

	return DEFAULT_COLOR;
}
