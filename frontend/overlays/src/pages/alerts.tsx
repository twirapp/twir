import { useCallback, useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import useWebSocket from 'react-use-websocket';

declare global {
	interface Window {
		webkitAudioContext: typeof AudioContext
	}
}

export const Alerts: React.FC = () => {
	const [url, setUrl] = useState<string | null>(null);
	const { apiKey } = useParams();
	const { lastMessage } = useWebSocket(url, {
		shouldReconnect: () => true,
		onOpen: () => console.log('Opened'),
		reconnectInterval: 500,
	});

	useEffect(() => {
		if (!lastMessage) return;
		const parsedData = JSON.parse(lastMessage.data);

		if (parsedData.eventName === 'trigger') {
			queueRef.current.push(parsedData.data);

			if (queueRef.current.length === 1) {
				processQueue();
			}
		}
	}, [lastMessage]);

	useEffect(() => {
		if (!apiKey) return;

		setUrl(`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/socket/alerts?apiKey=${apiKey}`);
	}, [apiKey]);

	const queueRef = useRef<Array<{
		id: string,
		channel_id: string,
		audio_id: string,
		audio_volume: number
	}>>([]);
	const currentAudioBuffer = useRef<AudioBufferSourceNode | null>(null);

	const processQueue = useCallback(async () => {
		if (queueRef.current.length === 0) {
			return;
		}

		const current = queueRef.current[0];
		if (current.audio_id) {
			await playAudio(current.channel_id, current.audio_id, current.audio_volume);
		}

		// change next val
		queueRef.current = queueRef.current.slice(1);

		// Process the next item in the queue
		processQueue();
	}, []);

	const playAudio = async (channelId: string, audioId: string, volume: number) => {
		const req = await fetch(`${window.location.origin}/cdn/twir/channels/${channelId}/${audioId}`);
		if (!req.ok) {
			console.error(await req.text());
			return;
		}

		const audioContext = new (window.AudioContext || window.webkitAudioContext)();
		const gainNode = audioContext.createGain();

		const data = await req.arrayBuffer();

		const source = audioContext.createBufferSource();
		currentAudioBuffer.current = source;

		source.buffer = await audioContext.decodeAudioData(data);

		gainNode.gain.value = volume / 100;
		source.connect(gainNode);
		gainNode.connect(audioContext.destination);

		return new Promise((resolve) => {
			source.onended = () => {
				currentAudioBuffer.current = null;
				resolve(null);
			};

			source.start(0);
		});
	};

	return <></>;
};
