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
  const [scenes, setScenes] = useState<OBSScenes>({});
  const [inputs, setInputs] = useState<OBSInputs>([]);
  const profile = useProfile();

  const setVolume = useCallback(async (opts: { audioSourceName: string, volume: number}) => {
    await context.socket?.call('SetInputVolume', {
      inputName: opts.audioSourceName,
      inputVolumeDb: opts.volume,
    });
  }, [context.socket]);

  const setScene = useCallback(async (opts: { sceneName: string }) => {
    context.socket?.call('GetCurrentProgramScene').then(console.log);
    context.socket?.call('SetCurrentProgramScene', { sceneName: opts.sceneName })
      .catch(console.error);
  }, [context.socket, context.connected]);

  const connect = useCallback(async () => {
    if (!profile.data) return;

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

    if (context.webSocket.current) {
      context.webSocket.current?.removeAllListeners();
      context.webSocket.current?.disconnect();
    }

    context.webSocket.current = io(
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

    context.webSocket.current.connect();
    context.webSocket.current.on('setScene', (data) => {
      setScene(data);
    });
  }, [context.socket, obsSettings, profile.data]);

  const disconnect = useCallback(() => {
    context.socket?.disconnect().then(() => context.setConnected(false));
    context.webSocket.current?.removeAllListeners();
    context.webSocket.current?.disconnect();
  }, [context.socket]);

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
    setVolume,
  };
};