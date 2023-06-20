import { useEffect, useState, useCallback } from 'react';

export function useWebSocket(apiKey: string | undefined) {
	const [ws, setWs] = useState<WebSocket | null>(null);

	const close = useCallback(() => {
		ws?.close();
		setWs(null);
	}, [ws]);

	const connect = useCallback(() => {
		const url = `${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/socket/tts?apiKey=${apiKey}`;
		setWs(new WebSocket(url));
	}, [apiKey]);

	useEffect(() => {
		connect();
	}, [apiKey]);

	return {
		ws,
		close,
		connect,
	};
}
