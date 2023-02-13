import { Button, Text } from '@mantine/core';
import { NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import OBSWebSocket from 'obs-websocket-js';
import { useContext, useEffect, useState } from 'react';

import { ObsWebsocketContext, useObsSocket } from '@/services/obsWebsocket';

interface IBeforeInstallPromptEvent extends Event {
  readonly platforms: string[];
  readonly userChoice: Promise<{
    outcome: 'accepted' | 'dismissed';
    platform: string;
  }>;
  prompt(): Promise<void>;
}

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['layout'])),
  },
});

export function useAddToHomescreenPrompt(): [
    IBeforeInstallPromptEvent | null,
  () => void
] {
  const [prompt, setState] = useState<IBeforeInstallPromptEvent | null>(
    null,
  );

  const promptToInstall = () => {
    if (prompt) {
      return prompt.prompt();
    }

    // return Promise.reject(
    //   new Error(
    //     'Tried installing before browser sent "beforeinstallprompt" event',
    //   ),
    // );
  };

  useEffect(() => {
    const ready = (e: IBeforeInstallPromptEvent) => {
      e.preventDefault();
      setState(e);
    };

    window.addEventListener('beforeinstallprompt', ready as any);

    return () => {
      window.removeEventListener('beforeinstallprompt', ready as any);
    };
  }, []);

  return [prompt, promptToInstall];
}


const Application: NextPage = () => {
  const [, promptInstall] = useAddToHomescreenPrompt();
  const obsSocket = useObsSocket();

  useEffect(() => {
    obsSocket.connect();
  }, []);

  return (<>
    <Text>{obsSocket.connected ? 'Connected' : 'Disconnected'}</Text>
    <Text>
      You can install site as application on your system. In this case you will be able to use dashboard without actual browser opened, and also it brings OBS Websocket support!
    </Text>
    <Button onClick={() => promptInstall()}>install</Button>
    <Button onClick={() => obsSocket.disconnect()}>disconnect</Button>
    <Button onClick={() => obsSocket.connect()}>connect</Button>
  </>);
};

export default Application;