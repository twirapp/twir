import { getCookie } from 'cookies-next';
import { useAtom } from 'jotai';
import { useCallback, useEffect, useState } from 'react';
import { io } from 'socket.io-client';

import { externalObsWsAtom, internalObsWsAtom, MyOBSWebsocket } from '../../stores/obs';

import { useProfile } from '@/services/api';
import { useObsModule } from '@/services/api/modules';

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


export const useInternalObsWs = () => {
  const [ws, setWs] = useAtom(internalObsWsAtom);
  const profile = useProfile();
  const obs = useObs();

  const connect = () => {
    disconnect();

    const webSocket = io(
      `${`${window.location.protocol == 'https:' ? 'wss' : 'ws'}://${
        window.location.host
      }`}/obs`,
      {
        transports: ['websocket'],
        autoConnect: true,
        auth: (cb) => {
          cb({ apiKey: profile.data?.apiKey, channelId: getCookie('dashboard_id') });
        },
      },
    );

    setWs(webSocket);

    webSocket?.off('setScene').on('setScene', (data) => {
      console.log(data, obs.connected);
      obs.setScene(data.sceneName);
    });
  };

  const disconnect = () => {
    ws?.removeAllListeners();
    ws?.disconnect();
    setWs(null);
  };

  return {
    connect,
    disconnect,
    connected: ws?.connected,
  };
};

export const useObs = () => {
  const [obs, setObs] = useAtom(externalObsWsAtom);
  const obsModule = useObsModule();
  const { data: obsSettings } = obsModule.useSettings();

  const [scenes, setScenes] = useState<OBSScenes>({});
  const [inputs, setInputs] = useState<OBSInputs>([]);

  const setScene = useCallback((sceneName: string) => {
    obs?.call('GetCurrentProgramScene').then(console.log);
    obs?.call('SetCurrentProgramScene', { sceneName })
      .catch(console.error);
  }, [obs]);

  const connect = async () => {
    await disconnect();

    if (!obsSettings) return;

    const newObs = new MyOBSWebsocket();
    await newObs.connect(`ws://${obsSettings.serverAddress}:${obsSettings.serverPort}`, obsSettings.serverPassword);
    setObs(newObs);
  };

  const disconnect = async () => {
    if (!obs) return;

    await obs.disconnect();
    setObs(null);
  };

  useEffect(() => {
    if (obs?.connected) {
      getScenes().then((newScenes) => {
        if (newScenes) {
          setScenes(newScenes);
        }
      });
      getInputs().then((inputs) => {
        setInputs(inputs);
      });
    }
  }, [obs?.connected]);

  const getScenes = useCallback(async (): Promise<OBSScenes | undefined> => {
    const scenesReq = await obs?.call('GetSceneList');
    if (!scenesReq) return;

    const mappedScenesNames = scenesReq.scenes.map(s => s.sceneName as string);

    const itemsPromises = await Promise.all(mappedScenesNames.map((sceneName) => {
      return obs?.call('GetSceneItemList', { sceneName });
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
        const group = await obs?.call('GetGroupSceneItemList', { sceneName: g as string });
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
  }, [obs]);

  const getInputs = useCallback(async () => {
    const req = await obs?.call('GetInputList');

    return req?.inputs.map(i => i.inputName as string) ?? [];
  }, [obs]);

  return {
    connected: obs?.connected,
    disconnect,
    connect,
    scenes,
    inputs,
    setScene,
  };
};