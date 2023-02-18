import OBSWebSocket from 'obs-websocket-js';
import {
  createContext,
  Dispatch, MutableRefObject,
  SetStateAction,
  useRef,
  useState,
} from 'react';
import { Socket } from 'socket.io-client';


export const ObsWebsocketContext = createContext({} as {
  socket: OBSWebSocket | null,
  setSocket: Dispatch<SetStateAction<OBSWebSocket | null>>
  connected: boolean,
  setConnected: Dispatch<SetStateAction<boolean>>,
  webSocket: MutableRefObject<Socket | null>
});

export function OBSWebsocketProvider({ children }: { children: React.ReactElement }) {
  const [socket, setSocket] = useState<OBSWebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const webSocket = useRef<Socket | null>(null);

  return (
      <ObsWebsocketContext.Provider
        value={{
        connected,
        setConnected,
        setSocket,
        socket,
        webSocket,
      }}>{children}</ObsWebsocketContext.Provider>
  );
}