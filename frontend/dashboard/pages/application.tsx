import {
  ActionIcon,
  Alert,
  Badge,
  Button,
  Card, createStyles,
  Divider,
  Flex,
  NumberInput, PasswordInput,
  Text,
  TextInput,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useStore } from 'jotai';
import { NextPage } from 'next';
import { useTranslation } from 'next-i18next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useEffect, useState } from 'react';

import { externalObsWsAtom } from '../stores/obs';

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
    ...(await serverSideTranslations(locale, ['application', 'layout'])),
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
  const store = useStore();
  const [, promptInstall] = useAddToHomescreenPrompt();
  const obsSocket = useObs();
  const obsSettingsManager = useObsModule();
  const obsSettingsUpdater = obsSettingsManager.useUpdate();
  const { data: obsSettings } = obsSettingsManager.useSettings();
  const obsStyles = useObsStyles();

  const { t } = useTranslation('application');

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
      {t('title')}
      <Button onClick={() => promptInstall()} size={'xs'} variant={'outline'} ml={5}>Install</Button>
    </Text>
    <Divider mt={5} />
    <Card shadow="sm" radius="md" withBorder mt={5}>
      <Card.Section withBorder p={'sm'}>
        <Flex direction={'row'} justify={'space-between'}>
          <Flex direction={'column'}>
            <Text>OBS Websocket</Text>
            <Text size={'xs'}>{t('obs.title')}</Text>
          </Flex>
          <Badge color={store.get(externalObsWsAtom)?.connected ? 'green' : 'red'}>
            {store.get(externalObsWsAtom)?.connected ? 'Connected' : 'Disconnected'}
          </Badge>
        </Flex>
      </Card.Section>
      <Card.Section p={'xs'}>
        <Alert color="cyan" mb={5}>
          <Text>
            {t('obs.info')}
          </Text>
        </Alert>
      </Card.Section>
      <Card.Section p={'sm'} withBorder className={obsStyles.classes.card}>
          <TextInput
            label={t('obs.address')}
            {...obsSettingsForm.getInputProps('serverAddress')} withAsterisk
          />
          <NumberInput label={t('obs.port')} {...obsSettingsForm.getInputProps('serverPort')} withAsterisk />
          <PasswordInput label={t('obs.password')} {...obsSettingsForm.getInputProps('serverPassword')} withAsterisk />
      </Card.Section>
      <Card.Section p={'sm'}>
        <Button color={'green'} onClick={saveObsSettings}>Save</Button>
      </Card.Section>
    </Card>
  </>);
};

export default Application;