import {
  ActionIcon,
  Alert,
  Badge,
  Center,
  Drawer,
  Flex,
  Group,
  Menu,
  NumberInput,
  ScrollArea,
  Slider,
  Text,
  Textarea,
  TextInput,
  useMantineTheme,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { IconGripVertical, IconMinus, IconPlus, IconVariable } from '@tabler/icons';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import { useEffect } from 'react';
import { DragDropContext, Draggable, Droppable } from 'react-beautiful-dnd';

type Props = {
  opened: boolean;
  timer: ChannelTimer;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

const sliderSteps = new Array(6).fill(1).map((_, index) => {
  const value = (index + 1) * 15;
  return { value, label: value.toString() };
});

export const TimerDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<ChannelTimer>({});
  const viewPort = useViewportSize();

  useEffect(() => {
    form.setValues(props.timer);
  }, [props.timer]);

  return (
    <Drawer
      opened={props.opened}
      onClose={() => props.setOpened(false)}
      title={<Badge size="xl">{props.timer.name}</Badge>}
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
            <TextInput {...form.getInputProps('name')} label="Name" required></TextInput>
            <div style={{ width: '100%' }}>
              <Text>Time interval (seconds)</Text>
              <Flex direction="row" wrap="wrap" gap="md" justify="flex-start" align="flex-start">
                <Slider
                  w={'70%'}
                  style={{ marginTop: 9 }}
                  size="sm"
                  marks={sliderSteps}
                  value={form.getInputProps('timeInterval').value}
                  onChange={(v) => form.setFieldValue('timeInterval', v)}
                />
                <NumberInput w={'20%'} {...form.getInputProps('timeInterval')} />
              </Flex>
            </div>

            <NumberInput label="Messages Interval" {...form.getInputProps('messageInterval')} />

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
              {!form.getInputProps('responses').value?.length && <Alert>No responses added</Alert>}
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
                                  <Menu position="bottom-end" shadow="md" width={200}>
                                    <Menu.Target>
                                      <ActionIcon variant="filled">
                                        <IconVariable size={18} />
                                      </ActionIcon>
                                    </Menu.Target>
                                    <Menu.Dropdown>
                                      <Menu.Item>qwe</Menu.Item>
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
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
