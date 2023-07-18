import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import useWebSocket from 'react-use-websocket';

import { useObs } from '../hooks/obs';

export const OBS: React.FC = () => {
  const { apiKey } = useParams();
  const obs = useObs();

  const [settings, setSettings] = useState<Record<string, any>>({});

  const [url, setUrl] = useState<string | null>(null);
  useEffect(() => {
    if (!apiKey) return;

    const urlPrefix = window.location.protocol === 'https:' ? 'wss' : 'ws';
    setUrl(`${urlPrefix}://${window.location.host}/socket/obs?apiKey=${apiKey}`);
  }, [apiKey]);

  const { lastMessage, sendMessage } = useWebSocket(url, {
    shouldReconnect: () => true,
    onOpen: () => console.log('Twir socket opened'),
    reconnectInterval: 500,
  });

  useEffect(() => {
    if (!lastMessage) return;
    const { eventName, data } = JSON.parse(lastMessage.data);

    if (eventName === 'connected') {
      sendMessage(JSON.stringify({ eventName: 'requestSettings' }));
    }

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
  }, [lastMessage]);

  useEffect(() => {
    if (!settings || !Object.keys(settings).length) {
      obs.disconnect();
      return;
    }

    obs.connect(settings.serverAddress, settings.serverPort, settings.serverPassword).then(() => {
      console.log('Twir obs socket opened');
      sendMessage(JSON.stringify({ eventName: 'obsConnected' }));

      obs.getSources().then((sources) => {
        if (!sources) return;
        sendMessage(JSON.stringify({
          eventName: 'setSources',
          data: sources,
        }));
      });

      obs.getAudioSources().then((sources) => {
        if (!sources) return;
        sendMessage(JSON.stringify({
          eventName: 'setAudioSources',
          data: sources,
        }));
      });

      const scenesHandler = async () => {
        const sources = await obs.getSources();
        sendMessage(JSON.stringify({
          eventName: 'setSources',
          data: sources,
        }));
      };

      const audioHandler = async () => {
        const sources = await obs.getAudioSources();
        sendMessage(JSON.stringify({
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
    });

    return () => {
      obs.disconnect();
    };
  }, [settings]);

  return <div></div>;
};
