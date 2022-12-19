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
  Input,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { IconMinus, IconPlus } from '@tabler/icons';
import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { useEffect, useState } from 'react';

import { typesMapping } from './mapping';

type Props = {
  opened: boolean;
  settings: ChannelModerationSetting;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const ModerationDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<ChannelModerationSetting>({});
  const viewPort = useViewportSize();
  const [title, setTitle] = useState('');

  useEffect(() => {
    form.setValues(props.settings);
    const titleValue = props.settings.type
      ? typesMapping[props.settings.type]?.name ?? props.settings.type
      : '';
    setTitle(titleValue);
  }, [props.settings]);

  return (
    <div>
      <Drawer
        opened={props.opened}
        onClose={() => props.setOpened(false)}
        title={<Badge size="xl">{title.charAt(0).toUpperCase() + title.slice(1)}</Badge>}
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
                  <TextInput label="Timeout message" {...form.getInputProps('banMessage')} />
                </Grid.Col>
                <Grid.Col xs={12} sm={3} md={4} lg={4} xl={4}>
                  <NumberInput label="Timeout time" {...form.getInputProps('banTime')} />
                </Grid.Col>
              </Grid>
              <NumberInput label="Warning message" {...form.getInputProps('warningMessage')} />
              <Group grow>
                {props.settings.type === 'links' && (
                  <Switch
                    label="Moderate clips"
                    labelPosition="left"
                    {...form.getInputProps('checkClips', { type: 'checkbox' })}
                  />
                )}
                <Switch
                  label="Moderate vips"
                  labelPosition="left"
                  {...form.getInputProps('vips', { type: 'checkbox' })}
                />
                <Switch
                  label="Moderate subscribers"
                  labelPosition="left"
                  {...form.getInputProps('subscribers', { type: 'checkbox' })}
                />
              </Group>
              {props.settings.type === 'emotes' && (
                <NumberInput
                  label="Max emotes in message"
                  required
                  {...form.getInputProps('triggerLength')}
                />
              )}
              {props.settings.type === 'symbols' && (
                <NumberInput
                  label="Max symbols in message (percent)"
                  required
                  {...form.getInputProps('maxPercentage')}
                />
              )}
              {props.settings.type === 'caps' && (
                <NumberInput
                  label="Max caps in message (percent)"
                  required
                  {...form.getInputProps('maxPercentage')}
                />
              )}
              {props.settings.type === 'longMessage' && (
                <NumberInput
                  label="Max message length"
                  required
                  {...form.getInputProps('triggerLength')}
                />
              )}
              <div>
                <Flex direction="row" gap="xs">
                  <Text>Deny list words</Text>
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
              </div>
            </Flex>
          </form>
        </ScrollArea.Autosize>
      </Drawer>
    </div>
  );
};
