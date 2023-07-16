export const copyToClipBoard = async (text: string) => {
	await navigator.clipboard.writeText(text);
};
