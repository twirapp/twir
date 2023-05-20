import {
  Flex,
  ScrollArea,
  TextInput,
  Text,
  Grid,
  ActionIcon,
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
  MultiSelect,
  Modal,
  Title,
  Divider,
} from '@mantine/core';
import { isNotEmpty, useForm } from '@mantine/form';
import { useDebouncedState, useViewportSize } from '@mantine/hooks';
import {
  IconChevronDown,
  IconGripVertical,
  IconMinus,
  IconPlus,
  IconSearch,
  IconShieldHalfFilled,
  IconVariable,
} from '@tabler/icons';
import type {
  ChannelCommand,
  CommandModule,
  CooldownType,
} from '@tsuwari/typeorm/entities/ChannelCommand';
import { useTranslation } from 'next-i18next';
import { useEffect, useState } from 'react';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';

import { noop } from '../../util/chore';

import { commandsGroupManager, commandsManager, useVariables, useRolesApi } from '@/services/api';

type Props = {
  opened: boolean;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
  command?: ChannelCommand;
};

const switches: Array<{
  prop: keyof ChannelCommand;
}> = [
  { prop: 'isReply' },
  { prop: 'visible' },
  { prop: 'keepResponsesOrder' },
  { prop: 'onlineOnly' },
];

type ChannelCommandForm = Omit<ChannelCommand, 'deniedUsersIds' | 'allowedUsersIds'> & {
  deniedUsersIds: Array<{ name: string }>;
  allowedUsersIds: Array<{ name: string }>;
};

export const CommandsModal: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<ChannelCommandForm>({
    validate: {
      name: (value) => {
        if (!value.length || value.trim().length == 0) return 'Name cannot be empty';
        return null;
      },
      deniedUsersIds: {
        name: isNotEmpty('User name cannot be empty'),
      },
      allowedUsersIds: {
        name: isNotEmpty('User name cannot be empty'),
      },
      cooldown: (value) => (value && value < 0 ? 'Cooldown cannot be lower then 0' : null),
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
      rolesIds: [],
      visible: true,
      responses: [],
      channelId: '',
      id: '',
      deniedUsersIds: [],
      allowedUsersIds: [],
      onlineOnly: false,
    },
  });

  const [aliases, setAliases] = useState<Array<string>>([]);
  const [aliasesSearch, setAliasesSearch] = useState('');

  const { t } = useTranslation('commands');
  const viewPort = useViewportSize();
  const { useCreateOrUpdate } = commandsManager();
  const updater = useCreateOrUpdate();

  const variables = useVariables();

  const rolesManager = useRolesApi();
  const { data: roles } = rolesManager.useGetAll();

  const groupsManager = commandsGroupManager();
  const { data: groups } = groupsManager.useGetAll();

  useEffect(() => {
    form.reset();
    setAliasesSearch('');
    setAliases([]);

    if (props.command) {
      form.setValues({
        ...props.command,
        deniedUsersIds: props.command.deniedUsersIds.map((a) => ({ name: a })) ?? [],
        allowedUsersIds: props.command.allowedUsersIds.map((a) => ({ name: a })) ?? [],
      });

      setAliases(props.command.aliases);
    }
  }, [props.command, props.opened]);

  function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      console.log(validate.errors);
      return;
    }

    updater
      .mutateAsync({
        id: form.values.id,
        data: {
          ...form.values,
          aliases: aliases,
          deniedUsersIds: form.values.deniedUsersIds.map((a) => a.name),
          allowedUsersIds: form.values.allowedUsersIds.map((a) => a.name),
        },
      })
      .then(() => {
        props.setOpened(false);
      })
      .catch(noop);
  }

  const [variablesSearchInput, setVariablesSearchInput] = useDebouncedState('', 200);

  return (
    <Modal
      opened={props.opened}
      onClose={() => props.setOpened(false)}
      title={
        <Button size="xs" color="green" onClick={onSubmit}>
          {t('drawer.save')}
        </Button>
      }
      padding="xl"
      size="xl"
      overlayColor={theme.colorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[2]}
      overlayOpacity={0.55}
      overlayBlur={3}
    >
      <form>
        <Flex
          direction="column"
          gap="md"
          justify="flex-start"
          align="flex-start"
          wrap="wrap"
          w={'100%'}
        >
          <Grid w={'100%'}>
            <Grid.Col span={4}>
              <TextInput
                label={t('name')}
                placeholder="coolcommand"
                withAsterisk
                {...form.getInputProps('name')}
              />
            </Grid.Col>
            <Grid.Col span={8}>
              <MultiSelect
                label={t('drawer.aliases.name')}
                data={aliases}
                value={aliases}
                placeholder={t('drawer.aliases.placeholder')!}
                searchable
                creatable
                withinPortal
                getCreateLabel={(query) => `+ Create ${query}`}
                onChange={(data) => {
                  setAliases(data);
                }}
                searchValue={aliasesSearch}
                onSearchChange={setAliasesSearch}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' || e.key === ';' || e.key === ',') {
                    if (aliases.includes(aliasesSearch)) return;
                    setAliases((data) => [...data, aliasesSearch]);
                    setAliasesSearch('');
                  }
                }}
              />
            </Grid.Col>
          </Grid>

          {form.values.module === 'CUSTOM' && (
            <div style={{ width: '100%' }}>
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
                            <Grid
                              {...provided.draggableProps}
                              ref={provided.innerRef}
                              mt="xs"
                              style={{
                                ...provided.draggableProps.style,
                                position: 'static',
                                width: '100%',
                              }}
                            >
                              <Grid.Col span={10}>
                                <Textarea
                                  placeholder="response"
                                  autosize={true}
                                  minRows={1}
                                  rightSection={
                                    <Menu
                                      position="bottom-end"
                                      shadow="md"
                                      width={380}
                                      offset={15}
                                      onClose={() => {
                                        setVariablesSearchInput('');
                                      }}
                                    >
                                      <Menu.Target>
                                        <ActionIcon color="blue" variant="filled">
                                          <IconVariable size={18} />
                                        </ActionIcon>
                                      </Menu.Target>

                                      <Menu.Dropdown>
                                        <TextInput
                                          placeholder={'search...'}
                                          size={'sm'}
                                          rightSection={<IconSearch size={12} />}
                                          onChange={(event) =>
                                            setVariablesSearchInput(event.target.value)
                                          }
                                        />
                                        <ScrollArea
                                          h={350}
                                          type={'always'}
                                          offsetScrollbars
                                          style={{ marginTop: 5 }}
                                        >
                                          {variables.data?.length &&
                                            variables.data
                                              .filter(
                                                (v) =>
                                                  v.name.includes(variablesSearchInput) ||
                                                  v.description?.includes(variablesSearchInput),
                                              )
                                              .map((v) => (
                                                <Menu.Item
                                                  key={v.name}
                                                  onClick={() => {
                                                    const insertValue = `${
                                                      v.example ? v.example : v.name
                                                    }`;
                                                    form.setFieldValue(
                                                      `responses.${index}.text`,
                                                      `${
                                                        form.values.responses![index]!.text
                                                      } $(${insertValue})`,
                                                    );
                                                    setVariablesSearchInput('');
                                                  }}
                                                >
                                                  <Flex direction={'column'}>
                                                    <Text>{v.name}</Text>
                                                    <Text size={'xs'} c="dimmed">
                                                      {v.description}
                                                    </Text>
                                                  </Flex>
                                                </Menu.Item>
                                              ))}
                                        </ScrollArea>
                                      </Menu.Dropdown>
                                    </Menu>
                                  }
                                  {...form.getInputProps(`responses.${index}.text`)}
                                />
                              </Grid.Col>
                              <Grid.Col span={'auto'}>
                                <Flex
                                  direction={'row'}
                                  align={'center'}
                                  justify={'center'}
                                  gap={20}
                                  w={'100%'}
                                  mt={5}
                                >
                                  <Center {...provided.dragHandleProps}>
                                    <IconGripVertical size={18} />
                                  </Center>
                                  <ActionIcon color={'red'} variant={'filled'}>
                                    <IconMinus
                                      size={18}
                                      onClick={() => {
                                        form.removeListItem('responses', index);
                                      }}
                                    />
                                  </ActionIcon>
                                </Flex>
                              </Grid.Col>
                            </Grid>
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

          <Textarea
            label={t('drawer.description.label')}
            placeholder={t('drawer.description.placeholder') ?? ''}
            {...form.getInputProps('description')}
            w={'100%'}
            autosize={true}
            minRows={1}
          />

          <Divider label={<Title order={3}>{t('drawer.permission')}</Title>} w={'100%'} />
          <Grid w={'100%'}>
            <Grid.Col span={12}>
              <MultiSelect
                data={
                  roles?.map((r) => ({
                    value: r.id,
                    label: r.name,
                    group: r.type !== 'CUSTOM' ? 'System' : 'Custom',
                  })) ?? []
                }
                icon={<IconShieldHalfFilled size={18} />}
                label={'Roles'}
                placeholder="That roles will access to command."
                description={'Leave blank for everyone.'}
                clearButtonLabel="Clear selection"
                clearable
                {...form.getInputProps('rolesIds')}
              />
            </Grid.Col>
            <Grid.Col span={6}>
              <Flex direction="row" gap="xs">
                <Text>Denied users</Text>
                <ActionIcon variant="light" color="green" size="xs">
                  <IconPlus
                    size={18}
                    onClick={() => {
                      form.insertListItem('deniedUsersIds', { name: '' });
                    }}
                  />
                </ActionIcon>
              </Flex>

              {!form.values.deniedUsersIds?.length && <Alert>No users added</Alert>}
              <ScrollArea.Autosize maxHeight={100} mx="auto" type="auto" offsetScrollbars={true}>
                <Grid grow gutter="xs" style={{ margin: 0, gap: 8 }}>
                  {form.values.deniedUsersIds?.map((_, i) => (
                    <Grid.Col style={{ padding: 0 }} key={i} xs={4} sm={4} md={4} lg={4} xl={4}>
                      <TextInput
                        placeholder="username"
                        rightSection={
                          <ActionIcon
                            variant="filled"
                            onClick={() => {
                              form.removeListItem('deniedUsersIds', i);
                            }}
                          >
                            <IconMinus size={18} />
                          </ActionIcon>
                        }
                        {...form.getInputProps(`deniedUsersIds.${i}.name`)}
                      />
                    </Grid.Col>
                  ))}
                </Grid>
              </ScrollArea.Autosize>
            </Grid.Col>
            <Grid.Col span={'auto'}>
              <Flex direction="row" gap="xs">
                <Text>Allowed users</Text>
                <ActionIcon variant="light" color="green" size="xs">
                  <IconPlus
                    size={18}
                    onClick={() => {
                      form.insertListItem('allowedUsersIds', { name: '' });
                    }}
                  />
                </ActionIcon>
              </Flex>

              {!form.values.allowedUsersIds?.length && <Alert>No users added</Alert>}
              <ScrollArea.Autosize maxHeight={100} mx="auto" type="auto" offsetScrollbars={true}>
                <Grid grow gutter="xs" style={{ margin: 0, gap: 8 }}>
                  {form.values.allowedUsersIds?.map((_, i) => (
                    <Grid.Col style={{ padding: 0 }} key={i} xs={4} sm={4} md={4} lg={4} xl={4}>
                      <TextInput
                        placeholder="username"
                        rightSection={
                          <ActionIcon
                            variant="filled"
                            onClick={() => {
                              form.removeListItem('allowedUsersIds', i);
                            }}
                          >
                            <IconMinus size={18} />
                          </ActionIcon>
                        }
                        {...form.getInputProps(`allowedUsersIds.${i}.name`)}
                      />
                    </Grid.Col>
                  ))}
                </Grid>
              </ScrollArea.Autosize>
            </Grid.Col>
          </Grid>

          <Divider label={<Title order={3}>Cooldown</Title>} w={'100%'} />
          <Grid w={'100%'}>
            <Grid.Col span={6}>
              <NumberInput
                defaultValue={0}
                placeholder="0"
                label={t('drawer.cooldown.time')}
                {...form.getInputProps('cooldown')}
              />
            </Grid.Col>
            <Grid.Col span={6}>
              <Select
                label={t('drawer.cooldown.type')}
                defaultValue="GLOBAL"
                {...form.getInputProps('cooldownType')}
                data={[
                  { value: 'GLOBAL', label: 'Global' },
                  { value: 'PER_USER', label: 'Per User' },
                ]}
              />
            </Grid.Col>
          </Grid>

          <Divider label={<Title order={3}>Settings</Title>} w={'100%'} />
          <Grid w={'100%'}>
            {switches.map(({ prop }, i) => (
              <Grid.Col key={i} span={6} w={'100%'}>
                <Grid>
                  <Grid.Col span={10}>
                    <Text>{t(`drawer.switches.${prop}.name`)}</Text>
                  </Grid.Col>
                  <Grid.Col span={2}>
                    <Switch {...form.getInputProps(prop, { type: 'checkbox' })} />
                  </Grid.Col>
                  <Grid.Col span={12} mt={-10}>
                    <Text c="dimmed" size={'xs'}>
                      {t(`drawer.switches.${prop}.description`)}
                    </Text>
                  </Grid.Col>
                </Grid>
              </Grid.Col>
            ))}
          </Grid>

          <Divider label={<Title order={3}>Misc</Title>} w={'100%'} />
          <Grid w={'100%'}>
            <Grid.Col span={6}>
              <Select
                label="Group"
                {...form.getInputProps('groupId')}
                data={
                  groups?.map((value) => ({
                    value: value.id,
                    label: value.name,
                  })) ?? []
                }
                placeholder={'Command group'}
                allowDeselect={true}
              />
            </Grid.Col>
          </Grid>
        </Flex>
      </form>
    </Modal>
  );
};
