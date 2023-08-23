export const playAudio = async (buffer: ArrayBuffer, volume = 100) => {
	const audioContext = new (window.AudioContext || window.webkitAudioContext)();
	const gainNode = audioContext.createGain();

	const source = audioContext.createBufferSource();

	source.buffer = await audioContext.decodeAudioData(buffer);

	gainNode.gain.value = volume / 100;
	source.connect(gainNode);
	gainNode.connect(audioContext.destination);

	return new Promise((resolve) => {
		source.onended = () => {
			resolve(null);
		};

		source.start(0);
	});
};
