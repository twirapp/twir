import {
  Drawer,
  Flex,
  ScrollArea,
  TextInput,
  Text,
  Grid,
  ActionIcon,
  Input,
  Select,
  Alert,
  NumberInput,
  Switch,
  Group,
  Center,
  Textarea,
  Menu,
  Button,
  useMantineTheme,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { IconGripVertical, IconMinus, IconPlus, IconVariable } from '@tabler/icons';
import type {
  ChannelCommand,
  CommandModule,
  CommandPermission,
  CooldownType,
} from '@tsuwari/typeorm/entities/ChannelCommand';
import { useTranslation } from 'next-i18next';
import { useEffect } from 'react';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';

import { noop } from '../../util/chore';

import { commandsManager, useVariables } from '@/services/api';

type Props = {
  opened: boolean;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
  command?: ChannelCommand;
};

const COMMAND_PERMS: Array<keyof typeof CommandPermission> = [
  'BROADCASTER',
  'MODERATOR',
  'VIP',
  'SUBSCRIBER',
  'FOLLOWER',
  'VIEWER',
];

const switches: Array<{
  prop: keyof ChannelCommand;
}> = [
  { prop: 'isReply' },
  { prop: 'visible' },
  { prop: 'keepResponsesOrder' },
];

export const CommandDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<ChannelCommand>({
    validate: {
      name: (value) => (!value.length ? 'Name cannot be empty' : null),
      aliases: (values) => (values.some((s) => !s.length) ? 'Aliase cannot be empty' : null),
      cooldown: (value) => (value && value < 0 ? 'Cooldown cannot be lower then 0' : null),
      permission: (v) => (!COMMAND_PERMS.includes(v as any) ? 'Unknown permission' : null),
      responses: {
        text: (value) => (value && !value.length ? 'Response cannot be empty' : null),
      },
    },
    initialValues: {
      aliases: [],
      name: '',
      cooldown: 0,
      cooldownType: 'GLOBAL' as CooldownType,
      default: false,
      defaultName: null,
      description: '',
      enabled: true,
      isReply: true,
      keepResponsesOrder: true,
      module: 'CUSTOM' as CommandModule,
      permission: 'VIEWER' as CommandPermission,
      visible: true,
      responses: [],
      channelId: '',
      id: '',
    },
  });

  const { t } = useTranslation('commands');
  const viewPort = useViewportSize();
  const { useCreateOrUpdate } = commandsManager();
  const updater = useCreateOrUpdate();

  const variables = useVariables();

  useEffect(() => {
    form.reset();
    if (props.command) {
      form.setValues(props.command);
    }
  }, [props.command, props.opened]);

  function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      return;
    }

    updater.mutateAsync({
      id: form.values.id,
      data: form.values,
    }).then(() => {
      props.setOpened(false);
      form.reset();
    }).catch(noop);
  }


  return (
    <Drawer
      opened={props.opened}
      onClose={() => props.setOpened(false)}
      title={
        <Button size="xs" color="green" onClick={onSubmit}>
          {t('drawer.save')}
        </Button>
      }
      padding="xl"
      size="xl"
      position="right"
      transition="slide-left"
      overlayColor={theme.colorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[2]}
      overlayOpacity={0.55}
      overlayBlur={3}
    >
      <ScrollArea.Autosize maxHeight={viewPort.height - 100} type="auto" offsetScrollbars={true}>
        <form>
          <Flex direction="column" gap="md" justify="flex-start" align="flex-start" wrap="wrap">
            <div>
              <TextInput
                label={t('name')}
                placeholder="coolcommand"
                withAsterisk
                {...form.getInputProps('name')}
              />
            </div>

            <div>
              <Flex direction="row" gap="xs">
                <Text>{t('drawer.aliases.name')}</Text>
                <ActionIcon variant="light" color="green" size="xs">
                  <IconPlus
                    size={18}
                    onClick={() => {
                      form.insertListItem('aliases', '');
                    }}
                  />
                </ActionIcon>
              </Flex>
              {!form.values.aliases?.length && <Alert>{t('drawer.aliases.emptyAlert')}</Alert>}
              <ScrollArea.Autosize maxHeight={100} mx="auto" type="auto" offsetScrollbars={true}>
                <Grid grow gutter="xs" style={{ margin: 0, gap: 8 }}>
                  {form.values.aliases?.map((_, i) => (
                    <Grid.Col style={{ padding: 0 }} key={i} xs={4} sm={4} md={4} lg={4} xl={4}>
                      <Input
                        placeholder="aliase"
                        {...form.getInputProps(`aliases.${i}`)}
                        rightSection={
                          <ActionIcon
                            variant="filled"
                            onClick={() => {
                              form.removeListItem('aliases', i);
                            }}
                          >
                            <IconMinus size={18} />
                          </ActionIcon>
                        }
                      />
                    </Grid.Col>
                  ))}
                </Grid>
              </ScrollArea.Autosize>
            </div>

            <div>
              <Flex direction="row" gap={5} wrap="wrap">
                <NumberInput
                  defaultValue={0}
                  placeholder="0"
                  label={t('drawer.cooldown.time')}
                  {...form.getInputProps('cooldown')}
                />

                <Select
                  label={t('drawer.cooldown.type')}
                  defaultValue="GLOBAL"
                  {...form.getInputProps('cooldownType')}
                  data={[
                    { value: 'GLOBAL', label: 'Global' },
                    { value: 'PER_USER', label: 'Per User' },
                  ]}
                />
              </Flex>
            </div>

            <div>
              <Select
                label={t('drawer.permission')}
                {...form.getInputProps('permission')}
                data={COMMAND_PERMS.map((value) => ({
                  value,
                  label: value,
                }))}
              />
            </div>

            <div>
              <Grid>
                {switches.map(({ prop }, i) => (
                  <Grid.Col key={i} span={12}>
                    <Switch
                      key={i}
                      labelPosition="left"
                      label={t(`drawer.switches.${prop}.name`)}
                      description={t(`drawer.switches.${prop}.description`)}
                      {...form.getInputProps(prop, { type: 'checkbox' })}
                    />
                  </Grid.Col>
                ))}
              </Grid>
            </div>

            {form.values.module === 'CUSTOM' && (
              <div style={{ width: 450 }}>
                <Flex direction="row" gap="xs">
                  <Text>{t('responses')}</Text>
                  <ActionIcon variant="light" color="green" size="xs">
                    <IconPlus
                      size={18}
                      onClick={() => {
                        form.insertListItem('responses', { text: '' });
                      }}
                    />
                  </ActionIcon>
                </Flex>
                {!form.getInputProps('responses').value?.length && (
                  <Alert>{t('drawer.responses.emptyAlert')}</Alert>
                )}
                <DragDropContext
                  onDragEnd={({ destination, source }) =>
                    form.reorderListItem('responses', {
                      from: source.index,
                      to: destination!.index,
                    })
                  }
                >
                  <Droppable droppableId="responses" direction="vertical">
                    {(provided) => (
                      <div {...provided.droppableProps} ref={provided.innerRef}>
                        {form.values.responses?.map((_, index) => (
                          <Draggable key={index} index={index} draggableId={index.toString()}>
                            {(provided) => (
                              <Group
                                style={{ width: '100%' }}
                                ref={provided.innerRef}
                                mt="xs"
                                {...provided.draggableProps}
                              >
                                <Textarea
                                  w={'80%'}
                                  placeholder="response"
                                  autosize={true}
                                  minRows={1}
                                  rightSection={
                                    <Menu position="bottom-end" shadow="md" width={250}>
                                      <Menu.Target>
                                        <ActionIcon variant="filled">
                                          <IconVariable size={18} />
                                        </ActionIcon>
                                      </Menu.Target>

                                      <Menu.Dropdown>
                                        <ScrollArea h={200} type={'always'} offsetScrollbars>
                                        {variables.data?.length && variables.data.map(v => (
                                          <Menu.Item key={v.name} onClick={() => {
                                            const insertValue = `${v.example ? v.example : v.name}`;
                                            form.setFieldValue(
                                              `responses.${index}.text`,
                                              `${form.values.responses![index]!.text} $(${insertValue})`,
                                            );
                                          }}>
                                            <Flex direction={'column'}>
                                              <Text>{v.name}</Text>
                                              <Text size={'xs'}>{v.description}</Text>
                                            </Flex>
                                          </Menu.Item>
                                        ))}

                                        </ScrollArea>
                                      </Menu.Dropdown>

                                    </Menu>
                                  }
                                  {...form.getInputProps(`responses.${index}.text`)}
                                />
                                <Center {...provided.dragHandleProps}>
                                  <IconGripVertical size={18} />
                                </Center>
                                <ActionIcon>
                                  <IconMinus
                                    size={18}
                                    onClick={() => {
                                      form.removeListItem('responses', index);
                                    }}
                                  />
                                </ActionIcon>
                              </Group>
                            )}
                          </Draggable>
                        ))}

                        {provided.placeholder}
                      </div>
                    )}
                  </Droppable>
                </DragDropContext>
              </div>
            )}
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
