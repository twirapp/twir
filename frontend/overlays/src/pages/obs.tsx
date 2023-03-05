import { useCallback, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

import { useObs } from '../hooks/obs';

export const OBS: React.FC = () => {
  const { apiKey } = useParams();
  const obs = useObs();
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [settings, setSettings] = useState<Record<string, any>>({});

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


  const onMessage = useCallback((msg: MessageEvent) => {
    const { eventName, data } = JSON.parse(msg.data);

    switch (eventName) {
      case 'settings': setSettings(data); break;
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
  }, [socket, obs.connected]);

  useEffect(() => {
    if (!socket) return;
    socket.onmessage = (msg) => onMessage(msg);
    socket.onopen = async () => {
      socket.send(JSON.stringify({ eventName: 'requestSettings' }));
    };

    socket.onclose = (e) => {
      console.log('closed');
      setSocket(null);
      setTimeout(() => {
        setSocket(connect());
      }, 1500);
    };
  }, [socket]);

  useEffect(() => {
    if (!settings) {
      obs.disconnect();
      return;
    }

    obs.connect(settings.serverAddress, settings.serverPort, settings.serverPassword).then(() => {
      console.log('obs connected');
    });

    return () => {
      obs.disconnect();
    };
  }, [settings]);

  useEffect(() => {
    if (!obs.connected || !socket) return;

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

    const scenesHandler = async () => {
      const sources = await obs.getSources();
      socket.send(JSON.stringify({
        eventName: 'setSources',
        data: sources,
      }));
    };

    const audioHandler = async () => {
      const sources = await obs.getAudioSources();
      socket.send(JSON.stringify({
        eventName: 'setAudioSources',
        data: sources,
      }));
    };

    obs.instance.current
      .on('SceneListChanged', scenesHandler)

      .on('InputCreated', audioHandler)
      .on('InputRemoved', audioHandler)
      .on('InputNameChanged', audioHandler)

      .on('SceneItemCreated', scenesHandler)
      .on('SceneItemRemoved', scenesHandler);
  }, [socket, obs.connected]);

  return <div></div>;
};