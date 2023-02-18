import { getCookie } from 'cookies-next';
import OBSWebSocket from 'obs-websocket-js';
import { useCallback, useContext, useEffect, useState } from 'react';
import { io } from 'socket.io-client';

import { useProfile } from '@/services/api';
import { useObsModule } from '@/services/api/modules';
import { ObsWebsocketContext } from '@/services/obs/provider';


type OBSScene = {
  sources: Array<{
    name: string,
    type: string | null
  }>
}

type OBSScenes = {
  [x: string]: OBSScene
}

type OBSInputs = string[]

export const useObs = () => {
  const context = useContext(ObsWebsocketContext);
  const obsModule = useObsModule();
  const { data: obsSettings } = obsModule.useSettings();
  const profile = useProfile();

  const [scenes, setScenes] = useState<OBSScenes>({});
  const [inputs, setInputs] = useState<OBSInputs>([]);

  const connect = async () => {
    if (!profile.data) return;

    if (context.socket) {
      await disconnect();
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
      return;
    }

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

    context.setWebSocket(webSocket);

    context.webSocket!.connect();
    context.webSocket!.on('setScene', (data) => {
      console.log(context);
      setScene(data);
    });
  };

  const setScene = useCallback((data) => {
    console.log(context);
    context.socket?.call('GetCurrentProgramScene').then(console.log);
    context.socket?.call('SetCurrentProgramScene', { sceneName: data.sceneName })
      .catch(console.error);
  }, [context]);

  const disconnect = async () => {
    await context.socket?.disconnect();
    context.setConnected(false);
    context.setSocket(null);

    context.webSocket?.removeAllListeners();
    context.webSocket?.disconnect();
    context.setWebSocket(null);
  };

  useEffect(() => {
    if (context.connected) {
      getScenes().then((newScenes) => {
        if (newScenes) {
          setScenes(newScenes);
        }
      });
      getInputs().then((inputs) => {
        setInputs(inputs);
      });
    }
  }, [context.connected]);

  const getScenes = useCallback(async (): Promise<OBSScenes | undefined> => {
    const scenesReq = await context.socket?.call('GetSceneList');
    if (!scenesReq) return;

    const mappedScenesNames = scenesReq.scenes.map(s => s.sceneName as string);

    const itemsPromises = await Promise.all(mappedScenesNames.map((sceneName) => {
      return context.socket?.call('GetSceneItemList', { sceneName });
    }));

    const result: OBSScenes = {};

    await Promise.all(itemsPromises.map(async (item, index) => {
      if (!item) return;
      const sceneName = mappedScenesNames[index];
      result[sceneName] = {
        sources: item.sceneItems.filter(i => !i.isGroup).map((i) => ({
          name: i.sourceName as string,
          type: i.inputKind?.toString() || null,
        })),
      };

      const groups = item.sceneItems
        .filter(i => i.isGroup)
        .map(g => g.sourceName);

      await Promise.all(groups.map(async (g) => {
        const group = await context.socket?.call('GetGroupSceneItemList', { sceneName: g as string });
        if (!group) return;

        result[sceneName].sources = [
          ...result[sceneName].sources,
          ...group.sceneItems.filter(i => !i.isGroup).map((i) => ({
            name: i.sourceName as string,
            type: i.inputKind?.toString() || null,
          })),
        ];
      }));
    }));

    return result;
  }, [context.socket]);

  const getInputs = useCallback(async () => {
    const req = await context.socket?.call('GetInputList');

    return req?.inputs.map(i => i.inputName as string) ?? [];
  }, [context.socket]);

  return {
    connect,
    disconnect,
    connected: context.connected,
    scenes,
    inputs,
    setScene,
  };
};