import OBSWebSocket from 'obs-websocket-js';
import { createContext, Dispatch, SetStateAction, useCallback, useContext, useState } from 'react';

export const ObsWebsocketContext = createContext({} as {
  socket: OBSWebSocket | null,
  setSocket: Dispatch<SetStateAction<OBSWebSocket | null>>
  connected: boolean,
  setConnected: Dispatch<SetStateAction<boolean>>,
});

export const useObsSocket = () => {
  // needed for global state
  const [socket, setSocket] = useState<OBSWebSocket | null>(null);
  const [connected, setConnected] = useState(false);

  const context = useContext(ObsWebsocketContext);

  const connect = useCallback(() => {
    const newSocket = new OBSWebSocket();
    newSocket.connect(`ws://127.0.0.1:4455`, '123456').then(async () => {
      context.setSocket(newSocket);
      context.setConnected(true);
      setConnected(true);
    });
  }, []);

  const disconnect = useCallback(() => {
    context.socket?.disconnect().then(() => setConnected(false));
  }, [context.socket]);

  return {
    socket,
    setSocket,
    connect,
    connected,
    setConnected,
    disconnect,
  };
};