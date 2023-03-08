import { useCallback, useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';

declare global {
  interface Window {
    webkitAudioContext: typeof AudioContext
  }
}

export const TTS: React.FC = () => {
  const { apiKey } = useParams();
  const [tts, setTTS] = useState<WebSocket | null>(null);

  const connect = () => {
    const url = `${`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}`}/socket/tts?apiKey=${apiKey}`;
    const socket = new WebSocket(url);
    return socket;
  };

  useEffect(() => {
    setTTS(connect());

    return () => {
      tts?.close();
    };
  }, []);

  const queueRef = useRef<Array<Record<string, string>>>([]);
  const currentAudioBuffer = useRef<AudioBufferSourceNode | null>(null);

  useEffect(() => {
    if (tts) {
      tts.onopen = () => {
        console.log('connected');
      };

      tts.onmessage = (msg) => {
        const parsedData = JSON.parse(msg.data);
        console.log(parsedData);

        if (parsedData.eventName === 'say') {
          queueRef.current.push(parsedData.data);

          if (queueRef.current.length === 1) {
            processQueue();
          }
        }

        if (parsedData.eventName === 'skip') {
          currentAudioBuffer.current?.stop();
        }
      };

      tts.onclose = () => {
        setTTS(null);
        setTimeout(() => {
          setTTS(connect());
        }, 1500);
      };
    }
  }, [tts]);

  const processQueue = useCallback(async () => {
    if (queueRef.current.length === 0) {
      return;
    }

    await say(queueRef.current[0]);
    queueRef.current = queueRef.current.slice(1);

    // Process the next item in the queue
    processQueue();
  }, []);

  const say = async (data: Record<string, string>) => {
    if (!apiKey || !data.text) return;

    const query = new URLSearchParams(data);

    const audioContext = new (window.AudioContext || window.webkitAudioContext)();
    const gainNode = audioContext.createGain();

    const req = await fetch(`/api/v1/tts/say?${query}`, {
      headers: {
        'Api-Key': apiKey,
      },
    });
    if (!req.ok) {
      currentAudioBuffer.current = null;
      return;
    }
    const arrayBuffer = await req.arrayBuffer();

    const source = audioContext.createBufferSource();
    currentAudioBuffer.current = source;

    source.buffer = await audioContext.decodeAudioData(arrayBuffer);

    gainNode.gain.value = parseInt(data.volume) / 100;
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
