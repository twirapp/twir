import { useCallback, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

import { useObs } from '../hooks/obs';

export const OBS: React.FC = () => {
  const { apiKey } = useParams();
  const obs = useObs();
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);

  const connect = () => {
    const url = `${`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}`}/socket/obs?apiKey=${apiKey}`;
    const socket = new WebSocket(url);
    return socket;
  };

  useEffect(() => {
    const conn = connect();
    setSocket(conn);

    return () => {
      conn?.close();
    };
  }, [apiKey]);


  useEffect(() => {
    if (!socket || !obs.instance) return;
    socket.onmessage = (msg) => {
      const { eventName, data } = JSON.parse(msg.data);

      switch (eventName) {
        case 'setScene': obs.setScene(data.sceneName); break;
        case 'toggleSource': obs.toggleSource(data.sourceName); break;
        case 'toggleAudioSource': obs.toggleAudioSource(data.audioSourceName); break;
        case 'setVolume': obs.setVolume(data.audioSourceName, data.volume); break;
        case 'increaseVolume': obs.changeVolume(data.audioSourceName, data.step, 'increase'); break;
        case 'decreaseVolume': obs.changeVolume(data.audioSourceName, data.step, 'decrease'); break;
        case 'enableAudio': obs.toggleAudioSource(data.audioSourceName, true); break;
        case 'disableAudio': obs.toggleAudioSource(data.audioSourceName, false); break;
        case 'startStart': obs.startStream(); break;
        case 'stopStream': obs.stopStream(); break;
      }
    };
  }, [socket, obs.instance]);

  useEffect(() => {
    if (!socket) return;
    socket.onopen = async () => {
      setConnected(true);
    };

    socket.onclose = (e) => {
      console.log('closed');
      setConnected(false);
      setSocket(null);
      setTimeout(() => {
        setSocket(connect());
      }, 1500);
    };
  }, [socket]);

  useEffect(() => {
    if (!connected) {
      obs.disconnect();
      return;
    }

    console.log('socket connected');
    obs.connect('localhost', 4455, '123456').then(() => console.log('obs connected'));

    return () => {
      obs.disconnect();
    };
  }, [connected]);

  useEffect(() => {
    if (!obs || !socket) return;

    obs.getSources().then((sources) => {
      if (!sources) return;
      socket.send(JSON.stringify({
        eventName: 'setSources',
        data: sources,
      }));
    });

    obs.getAudioSources().then((sources) => {
      if (!sources) return;
      socket.send(JSON.stringify({
        eventName: 'setAudioSources',
        data: sources,
      }));
    });
  }, [socket, obs.connected]);

  return <div></div>;
};