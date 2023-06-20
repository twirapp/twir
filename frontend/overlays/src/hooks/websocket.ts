export function useWebSocket(apiKey: string) {
    const [ws, setWs] = useState<WebSocket | null>(null);

    const connect = useCallback(() => {
        const url = `${`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}`}/socket/tts?apiKey=${apiKey}`;
        setWs(new WebSocket(url));
    }, [apiKey]);

    useEffect(() => {
        connect();
    }, [apiKey]);

    return {
        ws,
        connect,
    }
}