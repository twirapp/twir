import {
  ActionIcon,
  Alert,
  Badge,
  Button,
  Card, createStyles,
  Divider,
  Flex,
  Group,
  NumberInput,
  Text,
  TextInput,
  Tooltip,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { IconInfoCircle, IconInfoSquareRounded, IconQuestionMark } from '@tabler/icons';
import { NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useEffect, useState } from 'react';

import { useObsModule, type OBS } from '@/services/api/modules';
import { useObs } from '@/services/obs/hook';

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

    return Promise.reject(
      new Error(
        'Tried installing before browser sent "beforeinstallprompt" event',
      ),
    );
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

export const useObsStyles = createStyles((theme) => ({
  card: {
    backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.white,
  },
}));


const Application: NextPage = () => {
  const [, promptInstall] = useAddToHomescreenPrompt();
  const obsSocket = useObs();
  const obsSettingsManager = useObsModule();
  const obsSettingsUpdater = obsSettingsManager.useUpdate();
  const { data: obsSettings } = obsSettingsManager.useSettings();
  const obsStyles = useObsStyles();

  const obsSettingsForm = useForm<OBS['GET']>({
    initialValues: {
      serverAddress: 'localhost',
      serverPort: 4455,
      serverPassword: '',
    },
    validate: {
      serverAddress: (v) => !v.length ? 'Cannot be empty' : null,
      serverPort: (v) => !v ? 'Cannot be empty' : null,
      serverPassword: (v) => !v.length ? 'Cannot be empty' : null,
    },
  });

  useEffect(() => {
    if (!obsSettings) {
      return;
    }

    Object.entries(obsSettings).forEach((e) => {
      if (!e[1]) return;
      obsSettingsForm.setFieldValue(e[0], e[1]);
    });
  }, [obsSettings]);


  function saveObsSettings() {
    const validate = obsSettingsForm.validate();
    if (validate.hasErrors) return;

    obsSettingsUpdater.mutateAsync(obsSettingsForm.values);
  }

  return (<>
    <Text>
      You can install site as application on your system. In this case you will be able to use dashboard without actual browser opened, and also it brings OBS Websocket support
      <Button onClick={() => promptInstall()} size={'xs'} variant={'outline'} ml={5}>Install</Button>
    </Text>
    <Divider mt={5} />
    <Card shadow="sm" radius="md" withBorder mt={5}>
      <Card.Section withBorder p={'sm'}>
        <Flex direction={'row'} justify={'space-between'}>
          <Flex direction={'column'}>
            <Text>OBS Websocket</Text>
            <Text size={'xs'}>It brings support to events for hide/show scenes, mute/unmute audio and some other obs control things</Text>
          </Flex>
          <Badge color={obsSocket.connected ? 'green' : 'red'}>{obsSocket.connected ? 'Connected' : 'Disconnected'}</Badge>
        </Flex>
      </Card.Section>
      <Card.Section p={'xs'}>
        <Alert color="cyan" mb={5}>
          <Text>
            For working with obs we need you to keep site OPENED. Otherwise connection to obs will be closed.
          </Text>
          <Text>
            You can install the site as application for a more comfortable experience.
          </Text>
        </Alert>
      </Card.Section>
      <Card.Section p={'sm'} withBorder className={obsStyles.classes.card}>
          <TextInput
            label={`Address. Usually it's localhost, but if you advanced user you know what to do with that field.`}
            {...obsSettingsForm.getInputProps('serverAddress')} withAsterisk
          />
          <NumberInput label={'Port'} {...obsSettingsForm.getInputProps('serverPort')} withAsterisk />
          <TextInput label={'Password'} {...obsSettingsForm.getInputProps('serverPassword')} withAsterisk />
      </Card.Section>
      <Card.Section p={'sm'}>
        <Button color={'green'} onClick={saveObsSettings}>Save</Button>
      </Card.Section>
    </Card>
  </>);
};

export default Application;