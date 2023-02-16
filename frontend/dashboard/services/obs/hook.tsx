import OBSWebSocket from 'obs-websocket-js';
import { useCallback, useContext } from 'react';

import { useObsModule } from '@/services/api/modules';
import { ObsWebsocketContext } from '@/services/obs/provider';


export const useObsSocket = () => {
  const context = useContext(ObsWebsocketContext);
  const obsModule = useObsModule();
  const { data: obsSettings } = obsModule.useSettings();

  const connect = useCallback(async () => {
    if (context.socket) {
      await context.socket.disconnect();
      context.setSocket(null);
    }

    if (!obsSettings || !obsSettings.serverAddress || !obsSettings.serverPort) {
      return;
    }

    const newSocket = new OBSWebSocket();
    try {
      await newSocket.connect(`ws://${obsSettings.serverAddress}:${obsSettings.serverPort}`, obsSettings.serverPassword);
      context.setSocket(newSocket);
      context.setConnected(true);
    } catch (e) {
      console.log(e);
      context.setSocket(null);
      context.setConnected(false);
    }
  }, [context.socket, obsSettings]);

  const disconnect = useCallback(() => {
    context.socket?.disconnect().then(() => context.setConnected(false));
  }, [context.socket]);

  // const getScenes = useCallback(() => {
  //   return context.socket?.call('GetSceneList');
  // }, [context.socket]);
  //
  // const getScenesItems = useCallback((scene: string) => {
  //   return context.socket?.call('GetSceneItemList', { sceneName: 'Scene' });
  // }, [context.socket]);

  const getScenes = useCallback(async () => {
    const scenesReq = await context.socket?.call('GetSceneList');
    if (!scenesReq) return;

    const mappedScenesNames = scenesReq.scenes.map(s => s.sceneName as string);

    const itemsPromises = await Promise.all(mappedScenesNames.map((sceneName) => {
      return context.socket?.call('GetSceneItemList', { sceneName });
    }));

    return itemsPromises
      .flat()
      .map(item => {
        return item!.sceneItems.map((i) => ({
          name: i.sourceName,
        }));
      })
      .flat();
  }, [context.socket]);

  return {
    connect,
    disconnect,
    connected: context.connected,
    getScenes,
  };
};