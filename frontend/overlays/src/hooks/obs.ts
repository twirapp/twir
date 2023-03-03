import OBSWebSocket from 'obs-websocket-js';
import { useCallback, useState } from 'react';

type ObsSource = {
  name: string,
  type: string | null
}

type OBSScenes = {
  [x: string]: ObsSource[]
}

export const useObs = () => {
  const [obs, setObs] = useState<OBSWebSocket | null>(null);
  const [connected, setConnected] = useState(false);

  const connect = useCallback(async (address: string, port: number | string, password: string) => {
    if (!address || !port || !password) {
      return;
    }

    
    const newSocket = new OBSWebSocket();
    try {
      await newSocket.connect(`ws://${address}:${port}`, password);
      setObs(newSocket);
      setConnected(true);
    } catch (e) {
      setObs(null);
      setConnected(false);
      console.error(e);
      return;
    }
  }, []);

  const disconnect = useCallback(async () => {
    obs?.disconnect();
    setObs(null);
    setConnected(false);
  }, [obs]);

  const setScene = useCallback((sceneName: string) => {
    obs?.call('SetCurrentProgramScene', { sceneName })
      .catch(console.error);
  }, [obs]);

  const toggleSource = useCallback(async (sourceName: string) => {
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
  }, [obs]);

  const toggleAudioSource = useCallback(async (sourceName: string, muted?: boolean) => {
    if (typeof muted !== 'undefined') {
      await obs?.call('SetInputMute', { inputName: sourceName, inputMuted: !muted });
    } else {
      await obs?.call('ToggleInputMute', { inputName: sourceName });
    }
  }, [obs]);

  const setVolume = useCallback(async (inputName: string, volume: number) => {
    await obs?.call('SetInputVolume', {
      inputName,
      inputVolumeDb: volume * 3 - 60,
    });
  }, [obs]);

  const changeVolume = useCallback(async (inputName: string, step: number, operation: 'increase' | 'decrease') => {
    const currentVolumeReq = await obs?.call('GetInputVolume', { inputName });
    if (!currentVolumeReq) return;

    if (currentVolumeReq.inputVolumeDb === 0 && operation === 'increase') {
      return;
    }

    if (currentVolumeReq.inputVolumeDb <= 95 && operation === 'decrease') {
      return;
    }

    const newVolume = currentVolumeReq.inputVolumeDb + (operation === 'increase' ? step : -step);

    await obs?.call('SetInputVolume', {
      inputName,
      inputVolumeDb: newVolume,
    });
  }, [obs]);

  const startStream = useCallback(async () => {
    await obs?.call('StartStream');
  }, [obs]);

  const stopStream = useCallback(async () => {
    await obs?.call('StopStream');
  }, [obs]);

  const getSources = useCallback(async () => {
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
      result[sceneName] = item.sceneItems.filter(i => !i.isGroup).map((i) => ({
        name: i.sourceName as string,
        type: i.inputKind?.toString() || null,
      }));

      const groups = item.sceneItems
        .filter(i => i.isGroup)
        .map(g => g.sourceName);

      await Promise.all(groups.map(async (g) => {
        const group = await obs?.call('GetGroupSceneItemList', { sceneName: g as string });
        if (!group) return;

        result[sceneName] = [
          ...result[sceneName],
          ...group.sceneItems.filter(i => !i.isGroup).map((i) => ({
            name: i.sourceName as string,
            type: i.inputKind?.toString() || null,
          })),
        ];
      }));
    }));

    return result;
  }, [obs]);

  const getAudioSources = useCallback(async () => {
    const req = await obs?.call('GetInputList');

    return req?.inputs.map(i => i.inputName as string) ?? [];
  }, [obs]);

  return {
    connect,
    disconnect,
    connected,
    setScene,
    toggleSource,
    toggleAudioSource,
    setVolume,
    changeVolume,
    startStream,
    stopStream,
    getSources,
    getAudioSources,
  };
};