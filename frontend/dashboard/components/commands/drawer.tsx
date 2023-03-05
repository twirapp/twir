import {
  Drawer,
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
} from '@mantine/core';
import { isNotEmpty, useForm } from '@mantine/form';
import { useDebouncedState, useViewportSize } from '@mantine/hooks';
import { IconChevronDown, IconGripVertical, IconMinus, IconPlus, IconSearch, IconShieldHalfFilled, IconVariable } from '@tabler/icons';
import type {
  ChannelCommand,
  CommandModule,
  CooldownType,
} from '@tsuwari/typeorm/entities/ChannelCommand';
import { useTranslation } from 'next-i18next';
import { useEffect, useState } from 'react';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';

import { noop } from '../../util/chore';

import {
  commandsGroupManager,
  commandsManager,
  useVariables,
  useRolesApi,
} from '@/services/api';

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
];

type ChannelCommandForm = Omit<ChannelCommand, 'deniedUsersIds' | 'allowedUsersIds'> & {
  deniedUsersIds: Array<{ name: string }>,
  allowedUsersIds: Array<{ name: string }>,
}

export const CommandDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<ChannelCommandForm>({
    validate: {
      name: (value) => {
        if (!value.length) return 'Name cannot be empty';
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
        deniedUsersIds: props.command.deniedUsersIds.map(a => ({ name: a })) ?? [],
        allowedUsersIds: props.command.allowedUsersIds.map(a => ({ name: a })) ?? [],
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

    updater.mutateAsync({
      id: form.values.id,
      data: {
        ...form.values,
        aliases: aliases,
        deniedUsersIds: form.values.deniedUsersIds.map(a => a.name),
        allowedUsersIds: form.values.allowedUsersIds.map(a => a.name),
      },
    }).then(() => {
      props.setOpened(false);
    }).catch(noop);
  }

  const [variablesSearchInput, setVariablesSearchInput] = useDebouncedState('', 200);

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
              w={'100%'}
            />

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
                                {...provided.draggableProps}
                                ref={provided.innerRef}
                                mt="xs"
                                style={{ ...provided.draggableProps.style, position: 'static', width: '100%'  }}
                              >
                                <Textarea
                                  w={'80%'}
                                  placeholder="response"
                                  autosize={true}
                                  minRows={1}
                                  rightSection={
                                    <Menu position="bottom-end" shadow="md" width={380} offset={15} onClose={() => {
                                      setVariablesSearchInput('');
                                    }}>
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
                                          onChange={event => setVariablesSearchInput(event.target.value)}
                                        />
                                        <ScrollArea h={350} type={'always'} offsetScrollbars style={{ marginTop: 5 }}>
                                          {variables.data?.length && variables.data
                                            .filter(v => v.name.includes(variablesSearchInput)
                                              || v.description?.includes(variablesSearchInput))
                                            .map(v => (
                                              <Menu.Item key={v.name} onClick={() => {
                                                const insertValue = `${v.example ? v.example : v.name}`;
                                                form.setFieldValue(
                                                  `responses.${index}.text`,
                                                  `${form.values.responses![index]!.text} $(${insertValue})`,
                                                );
                                                setVariablesSearchInput('');
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

            <div>
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
            </div>

            <div>
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
            </div>

            <div style={{ width: '100%' }}>
              <Textarea
                  label={t('drawer.description.label')}
                  placeholder={t('drawer.description.placeholder') ?? ''}
                  {...form.getInputProps('description')}
                  w={'100%'}
                  autosize={true}
                  minRows={1}
              />
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

            <div style={{ width: '100%' }}>
              <MultiSelect
                data={roles?.map(r => ({
                  value: r.id,
                  label: r.name,
                  group: r.type !== 'CUSTOM' ? 'System' : 'Custom',
                })) ?? []}
                icon={<IconShieldHalfFilled size={18} />}
                label={t('drawer.permission')}
                placeholder="That roles will access to command."
                description={'Leave blank for everyone.'}
                clearButtonLabel="Clear selection"
                clearable
                {...form.getInputProps('rolesIds')}
              />
            </div>

            <div>
              <Select
                label='Group'
                {...form.getInputProps('groupId')}
                data={groups?.map((value) => ({
                  value: value.id,
                  label: value.name,
                })) ?? []}
                allowDeselect={true}
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

          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
