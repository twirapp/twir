import { getCookie } from 'cookies-next';
import OBSWebSocket from 'obs-websocket-js';
import {
  createContext,
  Dispatch, MutableRefObject,
  SetStateAction, useCallback, useEffect,
  useRef,
  useState,
} from 'react';
import { io, Socket } from 'socket.io-client';

import { useProfile } from '@/services/api';
import { useObsModule } from '@/services/api/modules';
import { useObs } from '@/services/obs/hook';

export const ObsWebsocketContext = createContext({} as {
  obs: OBSWebSocket | null,
  setObs: Dispatch<SetStateAction<OBSWebSocket | null>>
  
  connected: boolean,
  setConnected: Dispatch<SetStateAction<boolean>>,
  
  connect: () => Promise<void>
  disconnect: () => Promise<void>
});

export function OBSWebsocketProvider({ children }: { children: React.ReactElement }) {
  const [obs, setObs] = useState<OBSWebSocket | null>(null);
  const [connected, setConnected] = useState(false);

  const obsModule = useObsModule();
  const { data: obsSettings } = obsModule.useSettings();

  const connect = async () => {
    if (!obsSettings || !obsSettings.serverAddress || !obsSettings.serverPort) {
      return;
    }

    const newSocket = new OBSWebSocket();
    try {
      await newSocket.connect(`ws://${obsSettings.serverAddress}:${obsSettings.serverPort}`, obsSettings.serverPassword);
      setObs(newSocket);
      setConnected(true);
    } catch (e) {
      console.log(e);
      setObs(null);
      setConnected(false);
      return;
    }
  };

  const disconnect = async () => {
    obs?.disconnect();
    setConnected(false);
    setObs(null);
  };


  return (
      <ObsWebsocketContext.Provider
        value={{
          connected,
          setConnected,
          obs,
          setObs,
          connect,
          disconnect,
      }}>{children}</ObsWebsocketContext.Provider>
  );
}

export const InternalObsWebsocketContext = createContext({} as {
  socket: Socket | null,
  setSocket: Dispatch<SetStateAction<Socket | null>>,
});

export function InternalObsWebsocketProvider({ children }: { children: React.ReactElement }) {
  const obs = useObs();
  const [socket, setSocket] = useState<Socket | null>(null);
  const profile = useProfile();

  useEffect(() => {
    if (!profile.data) return;

    const webSocket = io(
      `${`${window.location.protocol == 'https:' ? 'wss' : 'ws'}://${
        window.location.host
      }`}/obs`,
      {
        transports: ['websocket'],
        autoConnect: false,
        auth: (cb) => {
          cb({ apiKey: profile.data?.apiKey, channelId: getCookie('dashboard_id') });
        },
      },
    );

    webSocket.connect();

    setSocket(webSocket);

    webSocket.on('setScene', (data) => {
      console.log(data);
      obs.setScene(data.sceneName);
    });

    return () => {
      webSocket.removeAllListeners();
      webSocket.disconnect();
      setSocket(null);
    };
  }, [profile.data]);

  return (
    <InternalObsWebsocketContext.Provider
      value={{
        socket,
        setSocket,
      }}>{children}</InternalObsWebsocketContext.Provider>
  );
}