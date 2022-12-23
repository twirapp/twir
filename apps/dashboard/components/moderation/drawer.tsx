import {
  ActionIcon,
  Badge,
  Drawer,
  Flex,
  Grid,
  Group,
  NumberInput,
  ScrollArea,
  Switch,
  TextInput,
  useMantineTheme,
  Text,
  Input, Button,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { IconMinus, IconPlus } from '@tabler/icons';
import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { getCookie } from 'cookies-next';
import { useTranslation } from 'next-i18next';
import { useEffect, useState } from 'react';

import { typesMapping } from './mapping';

import { queryClient, useModerationSettings } from '@/services/api';
import { SELECTED_DASHBOARD_KEY } from '@/services/dashboard';

type Props = {
  opened: boolean;
  settings: ChannelModerationSetting;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const ModerationDrawer: React.FC<Props> = (props) => {
  const { t } = useTranslation('moderation');
  const theme = useMantineTheme();
  const form = useForm<ChannelModerationSetting>({
    initialValues: {
      id: '',
      type: '' as any,
      enabled: true,
      banMessage: '',
      warningMessage: '',
      banTime: 0,
      blackListSentences: [],
      channelId: '',
      checkClips: true,
      maxPercentage: 50,
      subscribers: false,
      vips: false,
      triggerLength: 300,
    },
  });
  const viewPort = useViewportSize();
  const [title, setTitle] = useState('');

  useEffect(() => {
    form.setValues(props.settings);
    const titleValue = props.settings.type
      ? typesMapping[props.settings.type]?.name ?? props.settings.type
      : '';
    setTitle(titleValue);
  }, [props.settings]);

  const manager = useModerationSettings();

  function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      return;
    }

    const current = queryClient.getQueryData<ChannelCommand[]>([`/api/v1/channels/${getCookie(SELECTED_DASHBOARD_KEY)}/moderation`]);
    const currentIndex = current?.findIndex(t => t.id === props.settings.id);

    if (!current) {
      return;
    }
    current![currentIndex!] = form.values;

    manager.update.mutateAsync(current).then(() => {
      props.setOpened(false);
      form.reset();
    }).catch(() => {});
  }

  return (
    <div>
      <Drawer
        opened={props.opened}
        onClose={() => props.setOpened(false)}
        title={<Button size="xs" color="green" onClick={onSubmit}>
          Save
        </Button>}
        padding="xl"
        size="xl"
        position="right"
        transition="slide-left"
        overlayColor={theme.colorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[2]}
        overlayOpacity={0.55}
        overlayBlur={3}
      >
        <ScrollArea.Autosize maxHeight={viewPort.height - 120} type="auto" offsetScrollbars={true}>
          <form onSubmit={form.onSubmit((values) => console.log(values))}>
            <Flex direction="column" gap="md" justify="flex-start" align="flex-start" wrap="wrap">
              <Grid>
                <Grid.Col xs={12} sm={8} md={4} lg={4} xl={4}>
                  <TextInput label={t('drawer.timeout.message')} {...form.getInputProps('banMessage')} />
                </Grid.Col>
                <Grid.Col xs={12} sm={3} md={4} lg={4} xl={4}>
                  <NumberInput label={t('drawer.timeout.time')} {...form.getInputProps('banTime')} />
                </Grid.Col>
              </Grid>
              <NumberInput label={t('drawer.warning.message')} {...form.getInputProps('warningMessage')} />
              <Grid>
                {props.settings.type === 'links' && (
                  <Grid.Col>
                    <Switch
                      label={t('drawer.filters.clips')}
                      labelPosition="left"
                      {...form.getInputProps('checkClips', { type: 'checkbox' })}
                    />
                  </Grid.Col>
                )}
                <Grid.Col>
                  <Switch
                    label={t('drawer.filters.vips')}
                    labelPosition="left"
                    {...form.getInputProps('vips', { type: 'checkbox' })}
                  />
                </Grid.Col>
                <Grid.Col>
                  <Switch
                    label={t('drawer.filters.subs')}
                    labelPosition="left"
                    {...form.getInputProps('subscribers', { type: 'checkbox' })}
                  />
                </Grid.Col>
              </Grid>
              {props.settings.type === 'emotes' && (
                <NumberInput
                  label={t('drawer.maxEmotes')}
                  required
                  {...form.getInputProps('triggerLength')}
                />
              )}
              {props.settings.type === 'symbols' && (
                <NumberInput
                  label={t('drawer.maxSymbols')}
                  required
                  {...form.getInputProps('maxPercentage')}
                />
              )}
              {props.settings.type === 'caps' && (
                <NumberInput
                  label={t('drawer.maxCaps')}
                  required
                  {...form.getInputProps('maxPercentage')}
                />
              )}
              {props.settings.type === 'longMessage' && (
                <NumberInput
                  label={t('drawer.maxLength')}
                  required
                  {...form.getInputProps('triggerLength')}
                />
              )}
              {props.settings.type === 'blacklists' && <div>
                <Flex direction="row" gap="xs">
                  <Text>{t('drawer.denyListWords')}</Text>
                  <ActionIcon variant="light" color="green" size="xs">
                    <IconPlus
                        size={18}
                        onClick={() => {
                          form.insertListItem('blackListSentences', '');
                        }}
                    />
                  </ActionIcon>
                </Flex>

                <Grid grow gutter="xs" style={{ margin: 0, gap: 8 }}>
                  {form.values.blackListSentences?.map((_, i) => (
                      <Input
                          key={i}
                          placeholder="word"
                          {...form.getInputProps(`blackListSentences.${i}`)}
                          rightSection={
                            <ActionIcon
                                variant="filled"
                                onClick={() => {
                                  form.removeListItem('blackListSentences', i);
                                }}
                            >
                              <IconMinus size={18} />
                            </ActionIcon>
                          }
                      />
                  ))}
                </Grid>
              </div>}
            </Flex>
          </form>
        </ScrollArea.Autosize>
      </Drawer>
    </div>
  );
};
