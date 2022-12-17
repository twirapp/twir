import { jsx } from '@emotion/react';
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
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { IconGripVertical, IconInfoCircle, IconMinus, IconPlus } from '@tabler/icons';
import type { ChannelCommand, CommandPermission } from '@tsuwari/typeorm/entities/ChannelCommand';
import { useEffect } from 'react';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';

type Props = {
  opened: boolean;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
  command: ChannelCommand;
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
  label: string;
  description: string;
  prop: keyof ChannelCommand;
}> = [
  { label: 'Reply', description: 'Bot will send command response as reply', prop: 'isReply' },
  { label: 'Visible', description: 'Bot will send command response as reply', prop: 'visible' },
  {
    label: 'Keep Order',
    description: 'Bot will send command response as reply',
    prop: 'keepResponsesOrder',
  },
];

export const CommandDrawer: React.FC<Props> = (props) => {
  const form = useForm<ChannelCommand>({
    validate: {
      name: (value) => (!value.length ? 'Name cannot be empty' : null),
      aliases: (value) => (value.some((s) => !s.length) ? 'Aliase cannot be empty' : null),
    },
  });

  useEffect(() => {
    form.setValues(props.command);
  }, [props.command]);

  return (
    <Drawer
      opened={props.opened}
      onClose={() => props.setOpened(false)}
      title={`Edit !${props.command.name}`}
      padding="xl"
      size="xl"
      position="right"
      transition="slide-left"
      style={{ minWidth: 250, maxWidth: 500 }}
    >
      <form onSubmit={form.onSubmit((values) => console.log(values))}>
        <Flex direction="column" gap="md" justify="flex-start" align="flex-start" wrap="wrap">
          <div>
            <TextInput
              label="Name"
              placeholder="coolcommand"
              withAsterisk
              {...form.getInputProps('name')}
            />
          </div>

          {/* <div>
            <Flex direction="row" gap="xs">
              <Text>Responses</Text>
              <ActionIcon variant="light" color="green" size="xs">
                <IconPlus
                  size={18}
                  onClick={() => {
                    form.insertListItem('responses', { text: '' });
                  }}
                />
              </ActionIcon>
            </Flex>
            <ScrollArea mah={200} type="auto">
              {!form.values.responses?.length && (
                <Alert icon={<IconInfoCircle size={16} />} color="cyan">
                  There is no responses
                </Alert>
              )}
                {form.values.responses?.map((_, i) => (
                   <Input
                   key={i}
                      component="textarea"
                      autosize="true"
                      placeholder="response"
                      {...form.getInputProps(`responses.${i}.text`)}
                      rightSection={
                        <ActionIcon
                          variant="filled"
                          onClick={() => {
                            form.removeListItem('responses', i);
                          }}
                        >
                          <IconMinus size={18} />
                        </ActionIcon>
                      }
                    />
                ))}
              
            </ScrollArea>
          </div> */}
          <div style={{ width: 450 }}>
            <Flex direction="row" gap="xs">
              <Text>Responses</Text>
              <ActionIcon variant="light" color="green" size="xs">
                <IconPlus
                  size={18}
                  onClick={() => {
                    form.insertListItem('responses', { text: '' });
                  }}
                />
              </ActionIcon>
            </Flex>
            <ScrollArea mah={200} type="auto">
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
            </ScrollArea>
          </div>

          <div>
            <Flex direction="row" gap="xs">
              <Text>Aliases</Text>
              <ActionIcon variant="light" color="green" size="xs">
                <IconPlus
                  size={18}
                  onClick={() => {
                    form.insertListItem('aliases', '');
                  }}
                />
              </ActionIcon>
            </Flex>
            <ScrollArea mah={100} type="auto">
              {!form.values.aliases?.length && (
                <Alert icon={<IconInfoCircle size={16} />} color="cyan">
                  There is no aliases
                </Alert>
              )}
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
            </ScrollArea>
          </div>

          <div>
            <Flex direction="row" gap={5} wrap="wrap">
              <NumberInput
                defaultValue={0}
                placeholder="0"
                label="Cooldown time (seconds)"
                withAsterisk
                {...form.getInputProps('cooldown')}
              />

              <Select
                label="Cooldown Type"
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
              label="Permission"
              {...form.getInputProps('permission')}
              data={COMMAND_PERMS.map((value) => ({
                value,
                label: value,
              }))}
            />
          </div>

          <div>
            <Flex direction="row" gap={5} wrap="wrap">
              {switches.map(({ prop, ...rest }, i) => (
                <Switch
                  key={i}
                  labelPosition="left"
                  {...rest}
                  {...form.getInputProps(prop, { type: 'checkbox' })}
                />
              ))}
            </Flex>
          </div>
        </Flex>
      </form>
    </Drawer>
  );
};
