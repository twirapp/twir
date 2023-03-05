import {
  ActionIcon,
  Alert,
  Button,
  Card,
  CopyButton,
  Divider,
  Flex,
  Modal,
  MultiSelect,
  NumberInput,
  Select,
  Switch, Tabs,
  Text,
  Textarea,
  Tooltip,
} from '@mantine/core';
import { isNotEmpty, useForm, isInRange } from '@mantine/form';
import { IconCommand, IconCopy, IconSettings, IconSpeakerphone } from '@tabler/icons';
import { Fragment, useCallback, useEffect, useState } from 'react';

declare global {
  interface Window {
    webkitAudioContext: typeof AudioContext
  }
}

import { noop } from '../../util/chore';

import { CommandsList } from '@/components/commands/list';
import { authFetch, commandsManager, useProfile } from '@/services/api';
import { TTS, useTtsModule } from '@/services/api/modules';

import { CommandModule } from '@tsuwari/typeorm/entities/ChannelCommand';

export const TTSOverlay: React.FC = () => {
  const form = useForm<TTS['POST']>({
    initialValues: {
      enabled: false,
      pitch: 50,
      rate: 50,
      volume: 80,
      voice: '',
      allow_users_choose_voice_in_main_command: false,
      max_symbols: 0,
      disallowed_voices: [],
      do_not_read_emoji: false,
      do_not_read_twitch_emotes: false,
    },
    validate: {
      voice: isNotEmpty('Voice is required'),
      pitch: isInRange({ min: 1, max: 100 }, 'Pitch must be between 1 and 100'),
      rate: isInRange({ min: 1, max: 100 }, 'Rate must be between 1 and 100'),
      volume: isInRange({ min: 1, max: 100 }, 'Volume must be between 1 and 100'),
    },
  });

  const [modalOpened, setModalOpened] = useState(false);
  const [testText, setTestText] = useState('');

  const [availableVoices, setAvailableVoices] = useState<Array<{ value: string, label: string }>>([]);

  const tts = useTtsModule();
  const { data: ttsSettings } = tts.useSettings();
  const ttsInfo = tts.useInfo();
  const updater = tts.useUpdate();
  const { data: profile } = useProfile();

  const cmdsManager = commandsManager();
  const { data: commands } = cmdsManager.useGetAll();

  const [activeTab, setActiveTab] = useState<string | null>('settings');

  useEffect(() => {
    if (ttsSettings) {
      form.setValues(ttsSettings);
    }
  }, [ttsSettings]);

  useEffect(() => {
    if (!ttsInfo.data) return;

    const voices = Object.keys(ttsInfo.data.rhvoice_wrapper_voices_info)
      .sort((a, b) => {
        const dataA = ttsInfo.data?.rhvoice_wrapper_voices_info[a];
        const dataB = ttsInfo.data?.rhvoice_wrapper_voices_info[b];
        if (dataA.country === dataB.country) {
          return dataA.name.localeCompare(dataB.name);
        }
        return dataA.country.localeCompare(dataB.country);
      })
      .map((key) => {
        const data = ttsInfo.data?.rhvoice_wrapper_voices_info[key];

        return {
          value: key,
          label: `[${data.country}] ${data.name}`,
        };
      });

    setAvailableVoices(voices);
  }, [ttsInfo.data]);

  async function onSubmit() {
    if (form.validate().hasErrors) return;

    updater.mutateAsync(form.values)
      .then(() => setModalOpened(false))
      .catch(noop);
  }

  const testSpeak = useCallback(async () => {
    if (!testText) return;

    const query = new URLSearchParams({
      ...form.values as any,
      text: testText,
    });

    const audioContext = new (window.AudioContext || window.webkitAudioContext)();

    const req = await authFetch(`/api/v1/tts/say?${query}`);
    const arrayBuffer = await req.arrayBuffer();

    const source = audioContext.createBufferSource();
    source.buffer = await audioContext.decodeAudioData(arrayBuffer);
    source.connect(audioContext.destination);
    source.start(0);
  }, [form.values, testText]);

  return (
    <Fragment>
      <Card shadow="sm" p="lg" radius="md" w={200} withBorder>
        <Card.Section>
          <Flex direction={'row'} justify={'space-between'}>
            <div></div>
            <Flex direction={'row'} gap={0}>
                <CopyButton
                  value={'window' in globalThis
                    ? `${window.location.origin}/overlays/${profile?.apiKey}/tts`
                    : ''}
                >
                  {({ copied, copy }) => (
                    <Tooltip label={'Copy link to overlay'} withArrow arrowSize={5} color={'dark'}>
                        <ActionIcon
                          color={'dark'}
                          onClick={copy}
                          disabled={!!ttsSettings === false}
                        >
                          <IconCopy />
                        </ActionIcon>
                    </Tooltip>
                  )}
                </CopyButton>
              <ActionIcon color={'dark'} onClick={() => setModalOpened(true)}><IconSettings /></ActionIcon>
            </Flex>
          </Flex>
        </Card.Section>
        <Card.Section>
          <Flex direction={'column'} align={'center'}>
            <IconSpeakerphone size={100} />
            <Text size={45}>TTS</Text>
          </Flex>
        </Card.Section>
      </Card>

      <Modal
        opened={modalOpened}
        onClose={() => setModalOpened(false)}
        title={activeTab === 'settings' && <Button size={'sm'} variant={'light'} onClick={onSubmit} color={'green'}>Save</Button>}
        size={'xl'}
      >
        <Divider />
        <Tabs value={activeTab} onTabChange={setActiveTab} defaultValue="settings" radius={0}>
          <Tabs.List grow>
            <Tabs.Tab value="settings" icon={<IconSettings size={14} />}>Settings</Tabs.Tab>
            <Tabs.Tab value="commands" icon={<IconCommand size={14} />}>Commands</Tabs.Tab>
          </Tabs.List>
          <Tabs.Panel value="settings" pt="xs">
            <Flex mt={10} direction={'column'} gap={'md'}>
              <Alert><Text size={'xs'}>Hint: you can use events system to trigger tts on reward.</Text></Alert>
              <Switch
                label={'Enabled'}
                labelPosition={'left'}
                {...form.getInputProps('enabled', { type: 'checkbox' })}
              />
              <Select
                label="Voice"
                required
                data={availableVoices}
                {...form.getInputProps('voice')}
              />
              <NumberInput label={'Pitch'} max={100} min={1} required {...form.getInputProps('pitch')} />
              <NumberInput label={'Rate'} max={100} min={1} required {...form.getInputProps('rate')} />
              <NumberInput label={'Volume'} max={100} min={1} required {...form.getInputProps('volume')} />
              <Switch
                label={'Allow users use different voices in main (!tts) command'}
                labelPosition={'left'}
                {...form.getInputProps('allow_users_choose_voice_in_main_command', { type: 'checkbox' })}
              />
              <Switch
                label={'Do not read emoji'}
                labelPosition={'left'}
                {...form.getInputProps('do_not_read_emoji', { type: 'checkbox' })}
              />
              <Switch
                label={'Do read twitch emotes. Including 7tv, ffz, bttv.'}
                labelPosition={'left'}
                {...form.getInputProps('do_not_read_twitch_emotes', { type: 'checkbox' })}
              />
              <NumberInput
                label={'Max message length for tts. If setted to 0 then there is no restriction'}
                max={500}
                min={0}
                {...form.getInputProps('max_symbols')}
              />
              <MultiSelect
                label={'Disallowed for usage voices'}
                data={availableVoices}
                clearable
                {...form.getInputProps('disallowed_voices')}
              />
            </Flex>

            <Divider mt={10} mb={5} />

            <Textarea placeholder={'enter text for test'} value={testText} onChange={e => setTestText(e.target.value)} />
            <Button variant={'light'} onClick={testSpeak} fullWidth mt={10}>Test</Button>
          </Tabs.Panel>
          <Tabs.Panel value="commands" pt="xs">
            <CommandsList commands={commands?.filter(c => c.module === CommandModule.TTS) ?? []} />
          </Tabs.Panel>
        </Tabs>
      </Modal>
    </Fragment>
  );
};