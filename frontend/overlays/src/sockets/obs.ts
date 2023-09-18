import OBSWebSocket from 'obs-websocket-js';
import { ref } from 'vue';

export const useObs = () => {
  const obs = ref<OBSWebSocket>(new OBSWebSocket());
  const connected = ref(false);

  const connect = async (address: string, port: number | string, password: string) => {
    if (!address || !port || !password) {
      return;
    }

    await obs.value.disconnect();

    try {
      await obs.value.connect(`ws://${address}:${port}`, password);
      connected.value = true;
    } catch (e) {
      connected.value = false;
      throw e;
    }
  };

  const disconnect = async () => {
    await obs.value.disconnect();
  };

  const setScene = (sceneName: string) => {
    obs.value.call('SetCurrentProgramScene', { sceneName });
  };

  const toggleSource = async (sourceName: string) => {
    const currentSceneReq = await obs.value.call('GetCurrentProgramScene');
    if (!currentSceneReq) return;

    const [currentStateReq, idReq] = await Promise.all([
      obs.value.call('GetSourceActive', { sourceName }),
      obs.value.call('GetSceneItemId', { sourceName, sceneName: currentSceneReq.currentProgramSceneName }),
    ]);
    if (!currentStateReq || !idReq) return;

    await obs.value.call('SetSceneItemEnabled', {
      sceneName: currentSceneReq.currentProgramSceneName,
      sceneItemId: idReq.sceneItemId,
      sceneItemEnabled: !currentStateReq.videoShowing,
    });
  };

  const toggleAudioSource = async (sourceName: string, muted?: boolean) => {
    if (typeof muted !== 'undefined') {
      await obs.value.call('SetInputMute', { inputName: sourceName, inputMuted: !muted });
    } else {
      await obs.value.call('ToggleInputMute', { inputName: sourceName });
    }
  };

  const setVolume = async (inputName: string, volume: number) => {
    await obs.value.call('SetInputVolume', {
      inputName,
      inputVolumeDb: volume * 3 - 60,
    });
  };

  const changeVolume =async (inputName: string, step: number, operation: 'increase' | 'decrease') => {
    const currentVolumeReq = await obs.value.call('GetInputVolume', { inputName });
    if (!currentVolumeReq) return;

    if (currentVolumeReq.inputVolumeDb === 0 && operation === 'increase') {
      return;
    }

    if (currentVolumeReq.inputVolumeDb <= 95 && operation === 'decrease') {
      return;
    }

    const newVolume = currentVolumeReq.inputVolumeDb + (operation === 'increase' ? step : -step);

    await obs.value.call('SetInputVolume', {
      inputName,
      inputVolumeDb: newVolume,
    });
  };

  const startStream = async () => {
    await obs.value.call('StartStream');
  };

  const stopStream = async () => {
    await obs.value.call('StopStream');
  };

  const getSources = async () => {
    const scenesReq = await obs.value.call('GetSceneList');
    if (!scenesReq) return;

    const mappedScenesNames = scenesReq.scenes.map(s => s.sceneName as string);

    const itemsPromises = await Promise.all(mappedScenesNames.map((sceneName) => {
      return obs.value.call('GetSceneItemList', { sceneName });
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
        const group = await obs.value.call('GetGroupSceneItemList', { sceneName: g as string });
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
  };

  const getAudioSources = async () => {
    const req = await obs.value.call('GetInputList');

    return req?.inputs.map(i => i.inputName as string) ?? [];
  };

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
    instance: obs,
  };
};

type ObsSource = {
  name: string,
  type: string | null
}

type OBSScenes = {
  [x: string]: ObsSource[]
}
