import { useCallback, useContext, useEffect, useState } from 'react';

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

  const [scenes, setScenes] = useState<OBSScenes>({});
  const [inputs, setInputs] = useState<OBSInputs>([]);

  const setScene = useCallback((sceneName: string) => {
    console.log(context);
    context.obs?.call('GetCurrentProgramScene').then(console.log);
    context.obs?.call('SetCurrentProgramScene', { sceneName })
      .catch(console.error);
  }, [context.obs]);

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
    const scenesReq = await context.obs?.call('GetSceneList');
    if (!scenesReq) return;

    const mappedScenesNames = scenesReq.scenes.map(s => s.sceneName as string);

    const itemsPromises = await Promise.all(mappedScenesNames.map((sceneName) => {
      return context.obs?.call('GetSceneItemList', { sceneName });
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
        const group = await context.obs?.call('GetGroupSceneItemList', { sceneName: g as string });
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
  }, [context.obs]);

  const getInputs = useCallback(async () => {
    const req = await context.obs?.call('GetInputList');

    return req?.inputs.map(i => i.inputName as string) ?? [];
  }, [context.obs]);

  return {
    connect: context.connect,
    connected: context.connected,
    disconnect: context.disconnect,
    scenes,
    inputs,
    setScene,
  };
};