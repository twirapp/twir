import OBSWebSocket from 'obs-websocket-js';
import { createContext, Dispatch, SetStateAction, useCallback, useContext, useState } from 'react';

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
      }}
    >
    {children}
  </ObsWebsocketContext.Provider>
);
}

export const useObsSocket = () => {
  const context = useContext(ObsWebsocketContext);

  const connect = useCallback(async () => {
    if (context.socket) {
      await context.socket.disconnect();
      context.setSocket(null);
    }

    const newSocket = new OBSWebSocket();
    await newSocket.connect('ws://localhost:4455', '123456');
    context.setSocket(newSocket);
    context.setConnected(true);
  }, [context.socket]);

  const disconnect = useCallback(() => {
    context.socket?.disconnect().then(() => context.setConnected(false));
  }, [context.socket]);

  return {
    connect,
    disconnect,
    connected: context.connected,
  };
};