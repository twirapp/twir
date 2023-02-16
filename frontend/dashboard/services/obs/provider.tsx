import OBSWebSocket from 'obs-websocket-js';
import { createContext, Dispatch, SetStateAction, useState } from 'react';

export const ObsWebsocketContext = createContext({} as {
  socket: OBSWebSocket | null,
  setSocket: Dispatch<SetStateAction<OBSWebSocket | null>>
  connected: boolean,
  setConnected: Dispatch<SetStateAction<boolean>>,
});

export function OBSWebsocketProvider({ children }: { children: React.ReactElement }) {
  const [socket, setSocket] = useState<OBSWebSocket | null>(null);
  const [connected, setConnected] = useState(false);

  return (
      <ObsWebsocketContext.Provider
        value={{
        connected,
        setConnected,
        setSocket,
        socket,
      }}>{children}</ObsWebsocketContext.Provider>
  );
}