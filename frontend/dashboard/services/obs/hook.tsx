import { getCookie } from 'cookies-next';
import { useStore } from 'jotai';
import { useEffect, useState } from 'react';
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
  const jotaiStore = useStore();

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

    jotaiStore.set(internalObsWsAtom, webSocket);

    webSocket.off('setScene').on('setScene', (data) => {
      obs.setScene(data.sceneName);
    });

    webSocket.off('toggleSource').on('toggleSource', (data) => {
      // sourceName
      obs.toggleSource(data.sourceName);
    });

    webSocket.off('toggleAudioSource').on('toggleAudioSource', (data) => {
      obs.toggleAudioSource(data.audioSourceName);
    });

    webSocket.off('setVolume').on('setVolume', (data) => {
      obs.setVolume(data.audioSourceName, data.volume);
    });

    webSocket.off('increaseVolume').on('increaseVolume', (data) => {
      obs.changeVolume(data.audioSourceName, data.step, 'increase');
    });

    webSocket.off('decreaseVolume').on('decreaseVolume', (data) => {
      obs.changeVolume(data.audioSourceName, data.step, 'decrease');
    });

    webSocket.off('enableAudio').on('enableAudio', (data) => {
      obs.toggleAudioSource(data.audioSourceName, true);
    });

    webSocket.off('disableAudio').on('disableAudio', (data) => {
      obs.toggleAudioSource(data.audioSourceName, false);
    });
  };

  const disconnect = () => {
    const ws = jotaiStore.get(internalObsWsAtom);

    ws?.removeAllListeners();
    ws?.disconnect();
    jotaiStore.set(internalObsWsAtom, null);
  };

  return {
    connect,
    disconnect,
    connected: jotaiStore.get(internalObsWsAtom)?.connected,
  };
};

export const useObs = () => {
  const jotaiStore = useStore();

  const obsModule = useObsModule();
  const { data: obsSettings } = obsModule.useSettings();

  const [scenes, setScenes] = useState<OBSScenes>({});
  const [inputs, setInputs] = useState<OBSInputs>([]);

  const setScene = (sceneName: string) => {
    const obs = jotaiStore.get(externalObsWsAtom);

    obs?.call('SetCurrentProgramScene', { sceneName })
      .catch(console.error);
  };

  const toggleSource = async (sourceName: string) => {
    const obs = jotaiStore.get(externalObsWsAtom);

    const currentSceneReq = await obs?.call('GetCurrentProgramScene');
    if (!currentSceneReq) return;

    const [currentStateReq, idReq] = await Promise.all([
      obs?.call('GetSourceActive', { sourceName }),
      obs?.call('GetSceneItemId', { sourceName, sceneName: currentSceneReq.currentProgramSceneName }),
    ]);
    if (!currentStateReq || !idReq) return;

    await obs?.call('SetSceneItemEnabled', {
      sceneName: currentSceneReq.currentProgramSceneName,
      sceneItemId: idReq.sceneItemId,
      sceneItemEnabled: !currentStateReq.videoShowing,
    });
  };

  const toggleAudioSource = async (sourceName: string, muted?: boolean) => {
    const obs = jotaiStore.get(externalObsWsAtom);

    if (typeof muted !== 'undefined') {
      await obs?.call('SetInputMute', { inputName: sourceName, inputMuted: !muted });
    } else {
      await obs?.call('ToggleInputMute', { inputName: sourceName });
    }
  };

  const setVolume = async (inputName: string, volume: number) => {
    const obs = jotaiStore.get(externalObsWsAtom);

    await obs?.call('SetInputVolume', {
      inputName,
      inputVolumeDb: volume * 3 - 60,
    });
  };

  const changeVolume = async (inputName: string, step: number, operation: 'increase' | 'decrease') => {
    const obs = jotaiStore.get(externalObsWsAtom);

    const currentVolumeReq = await obs?.call('GetInputVolume', { inputName });
    if (!currentVolumeReq) return;

    if (currentVolumeReq.inputVolumeDb === 0 && operation === 'increase') {
      return;
    }

    if (currentVolumeReq.inputVolumeDb <= 95 && operation === 'decrease') {
      return;
    }

    const newVolume = currentVolumeReq.inputVolumeDb + (operation === 'increase' ? step : -step);
    console.log(newVolume, operation, step);

    await obs?.call('SetInputVolume', {
      inputName,
      inputVolumeDb: newVolume,
    });
  };

  const connect = async () => {
    await disconnect();

    if (!obsSettings) return;

    const newObs = new MyOBSWebsocket();
    await newObs.connect(`ws://${obsSettings.serverAddress}:${obsSettings.serverPort}`, obsSettings.serverPassword);

    jotaiStore.set(externalObsWsAtom, newObs);
  };

  const disconnect = async () => {
    const obs = jotaiStore.get(externalObsWsAtom);

    if (!obs) return;

    await obs.disconnect();
    jotaiStore.set(externalObsWsAtom, null);
  };

  useEffect(() => {
    const obs = jotaiStore.get(externalObsWsAtom);

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
  }, [jotaiStore.get(externalObsWsAtom)]);

  const getScenes = async (): Promise<OBSScenes | undefined> => {
    const obs = jotaiStore.get(externalObsWsAtom);

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
  };

  const getInputs = async () => {
    const obs = jotaiStore.get(externalObsWsAtom);

    const req = await obs?.call('GetInputList');

    return req?.inputs.map(i => i.inputName as string) ?? [];
  };

  return {
    connected: jotaiStore.get(externalObsWsAtom)?.connected,
    disconnect,
    connect,
    scenes,
    inputs,
    setScene,
    toggleSource,
    toggleAudioSource,
    setVolume,
    changeVolume,
  };
};