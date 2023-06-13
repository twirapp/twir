import {
  ActionIcon,
  Alert,
  Button,
  Card,
  CopyButton,
  Divider,
  Flex,
  Modal, MultiSelect, NumberInput, PasswordInput, Select,
  Switch,
  Text, Textarea, TextInput,
  Tooltip,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { IconBroadcast, IconCopy, IconInfoSquare, IconInfoSquareRounded, IconSettings, IconSpeakerphone, IconTooltip } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { useEffect, useState } from 'react';

import { useProfile } from '@/services/api';
import { type OBS, useObsModule } from '@/services/api/modules';

export const OBSOverlay: React.FC = () => {
  const form = useForm<OBS['GET']>({
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

  const { t } = useTranslation('application');

  const { data: profile } = useProfile();
  const [modalOpened, setModalOpened] = useState(false);

  const obsManager = useObsModule();
  const updater = obsManager.useUpdate();
  const { data: obsSettings } = obsManager.useSettings();
  const { data: obsData } = obsManager.useData();

  function saveObsSettings() {
    const validate = form.validate();
    if (validate.hasErrors) return;

    updater.mutateAsync(form.values);
  }

  useEffect(() => {
    if (!obsSettings) {
      return;
    }

    Object.entries(obsSettings).forEach((e) => {
      if (!e[1]) return;
      form.setFieldValue(e[0], e[1]);
    });
  }, [obsSettings]);

  return (<>
    <Card shadow="sm" p="lg" radius="md" w={200} withBorder>
      <Card.Section>
        <Flex direction={'row'} justify={'space-between'}>
          <div></div>
          <Flex direction={'row'} gap={0}>
            <Tooltip zIndex={500} label={'Connects obs with twir'} withArrow arrowSize={5} color={'dark'}>
              <ActionIcon
                color={'dark'}
              >
                <IconInfoSquareRounded/>
              </ActionIcon>
            </Tooltip>
            <CopyButton
              value={'window' in globalThis
                ? `${window.location.origin}/overlays/${profile?.apiKey}/obs`
                : ''}
            >
              {({ copied, copy }) => (
                <Tooltip label={'Copy link to overlay'} withArrow arrowSize={5} color={'dark'}>
                  <ActionIcon
                    color={'dark'}
                    onClick={copy}
                    disabled={!!obsSettings === false}
                  >
                    <IconCopy/>
                  </ActionIcon>
                </Tooltip>
              )}
            </CopyButton>
            <ActionIcon color={'dark'} onClick={() => setModalOpened(true)}><IconSettings/></ActionIcon>
          </Flex>
        </Flex>
      </Card.Section>
      <Card.Section>
        <Flex direction={'column'} align={'center'}>
          <IconBroadcast size={100}/>
          <Text size={45}>OBS</Text>
        </Flex>
      </Card.Section>
    </Card>

    <Modal
      opened={modalOpened}
      onClose={() => setModalOpened(false)}
      title={<Button size={'sm'} variant={'light'} onClick={saveObsSettings} color={'green'}>Save</Button>}
			closeOnClickOutside={false}
    >
      <Divider />
      <Flex mt={10} direction={'column'} gap={'md'}>
        <Alert>
          <Text size={'xs'}>This overlay used for connect TwirApp with your obs. Paste overlay link as browser source in obs.</Text>
        </Alert>
        <TextInput
          label={t('obs.address')}
          {...form.getInputProps('serverAddress')} withAsterisk
        />
        <NumberInput label={t('obs.port')} {...form.getInputProps('serverPort')} withAsterisk />
        <PasswordInput label={t('obs.password')} {...form.getInputProps('serverPassword')} withAsterisk />

        <Alert title={'Sources'} variant="outline" color={obsData?.sources?.length ? 'green' : 'red'}>
          {obsData?.sources?.join(', ') || 'Empty'}
        </Alert>
        <Alert title={'Audio sources'} variant="outline" color={obsData?.audioSources?.length ? 'green' : 'red'}>
          {obsData?.audioSources?.join(', ') || 'Empty'}
        </Alert>
        <Alert title={'Scenes'} variant="outline" color={obsData?.scenes?.length ? 'green' : 'red'}>
          {obsData?.scenes?.join(', ') || 'Empty'}
        </Alert>
      </Flex>
    </Modal>
  </>);
};
