import { WrapTextIcon, ListCollapseIcon, BoldIcon, ItalicIcon, LinkIcon, StrikethroughIcon, UnderlineIcon, ListIcon, Heading1Icon, Heading2Icon, ImageIcon, QuoteIcon } from 'lucide-vue-next';
import { computed, type Component } from 'vue';

import { useNotificationsForm } from './use-notifications-form';

const TEXTAREA_MODIFIERS = {
	h1: '<h1 id="h1">|</h1>',
	h2: '<h2 id="h2">|</h2>',
	br: '<br>',
	bold: '<b>|</b>',
	italic: '<i>|</i>',
	strikethrough: '<s>|</s>',
	underline: '<u>|</u>',
	link: '<a href="#" id="a">|</a>',
	img: '<img src="#">',
	blockquote: '<blockquote id="bq">|</blockquote>',
	list: '<ul id="ul"><li>|</li></ul>',
	spoiler: '<details><summary>|</summary> </details>',
} as const;

interface TextareaButton {
	name: keyof typeof TEXTAREA_MODIFIERS;
	title?: string;
	icon: Component;
}

export const textareaButtons: TextareaButton[] = [
	{
		name: 'h1',
		title: 'Heading 1',
		icon: Heading1Icon,
	},
	{
		name: 'h2',
		title: 'Heading 2',
		icon: Heading2Icon,
	},
	{
		name: 'br',
		title: 'Break',
		icon: WrapTextIcon,
	},
	{
		name: 'bold',
		title: 'Bold',
		icon: BoldIcon,
	},
	{
		name: 'italic',
		title: 'Italic',
		icon: ItalicIcon,
	},
	{
		name: 'strikethrough',
		title: 'Strikethrough',
		icon: StrikethroughIcon,
	},
	{
		name: 'underline',
		title: 'Underline',
		icon: UnderlineIcon,
	},
	{
		name: 'link',
		title: 'Link',
		icon: LinkIcon,
	},
	{
		name: 'img',
		title: 'Image',
		icon: ImageIcon,
	},
	{
		name: 'blockquote',
		title: 'Blockquote',
		icon: QuoteIcon,
	},
	{
		name: 'list',
		title: 'List',
		icon: ListIcon,
	},
	{
		name: 'spoiler',
		title: 'Spoiler',
		icon: ListCollapseIcon,
	},
];

export const useTextarea = () => {
	const notificationForm = useNotificationsForm();
	const textareaRef = computed({
		get() {
			return notificationForm.messageField.fieldRef;
		},
		set(el: any) {
			notificationForm.messageField.fieldRef = el?.$el;
		},
	});

	function getCursorPosition() {
		const el = notificationForm.messageField.fieldRef!;
		return { start: el.selectionStart ?? 0, end: el.selectionEnd ?? 0 };
	}

	function updateTextarea(newValue: string) {
		notificationForm.messageField.fieldModel = newValue;
	}

	function applyModifier(modifier: keyof typeof TEXTAREA_MODIFIERS) {
		const mod = TEXTAREA_MODIFIERS[modifier];
		if (!mod) {
			throw new Error('Modifier not implemented');
		}

		const { start, end } = getCursorPosition();
		const msg = notificationForm.formValues.message ?? '';
		const selection = msg.slice(start, end);
		updateTextarea(`${msg.slice(0, start)}${mod.replace('|', selection ?? '')}${msg.slice(end)}`);
	}

	return {
		textareaRef,
		applyModifier,
	};
};
