import {
  ActionIcon,
  Button,
  Card,
  Drawer,
  Flex,
  ScrollArea,
  TextInput,
  Center,
  Textarea,
  useMantineTheme,
  Select,
  NumberInput,
  createStyles, Group, Menu, Text, CopyButton, Code, Divider,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import {
  IconArrowBigDownLines,
  IconGripVertical,
  IconHandFinger, IconMinus,
  IconPlus,
  IconSearch,
  IconVariable,
  IconX,
} from '@tabler/icons';
import { Event, EventType } from '@tsuwari/typeorm/entities/events/Event';
import { EventOperation, OperationType } from '@tsuwari/typeorm/entities/events/EventOperation';
import { useTranslation } from 'next-i18next';
import React, { useEffect, useState } from 'react';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';

import { noop } from '../../util/chore';
import { eventsMapping } from './eventsMapping';

import { operationMapping } from '@/components/events/operationMapping';
import { RewardItem, RewardItemProps } from '@/components/reward';
import { commandsManager, eventsManager, useRewards } from '@/services/api';

type Props = {
  opened: boolean;
  event?: Event;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

const useStyles = createStyles(() => ({
  root: {
    position: 'relative',
    width: '100%',
  },

  label: {
    position: 'absolute',
    zIndex: 2,
    top: -15,
    right: 5,
  },
}));

export const EventsDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<Event>({
    initialValues: {
      id: '',
      type: '' as EventType,
      channelId: '',
      description: '',
      commandId: '',
      operations: [],
      rewardId: '',
    },
  });
  const manager = eventsManager();
  const updater = manager.useCreateOrUpdate();
  const { t } = useTranslation('events');
  const viewPort = useViewportSize();
  const cardClasses = useStyles();
  const [rewards, setRewards] = useState<RewardItemProps[]>([]);

  const commandManager = commandsManager();
  const commandList = commandManager.useGetAll();

  const rewardsManager = useRewards();
  const { data: rewardsData } = rewardsManager();

  useEffect(() => {
    form.reset();
    if (props.event) {
      form.setValues(props.event);
    }
  }, [props.event, props.opened]);

  useEffect(() => {
    if (rewardsData) {
      const data = rewardsData
        .sort((a, b) => (a.is_user_input_required === b.is_user_input_required ? 1 : -1))
        .map(
          (r) =>
            ({
              value: r.id,
              label: r.title,
              description: '',
              image: r.image?.url_4x || r.default_image?.url_4x,
              disabled: false,
            } as RewardItemProps),
        );

      setRewards(data);
    }
  }, [rewardsData]);

  async function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      console.log(validate.errors);
      return;
    }

    await updater.mutateAsync({
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
          Save
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
      <ScrollArea.Autosize maxHeight={viewPort.height - 120} type="auto" offsetScrollbars={true}>
        <form style={{ minHeight: viewPort.height - 150 }} onSubmit={form.onSubmit((values) => console.log(values))}>
          <Flex direction="column" gap="md" justify="flex-start" align="flex-start" wrap="wrap">
            <Select
              label={'Event'}
              searchable={true}
              disabled={!!form.values.id}
              data={Object.keys(eventsMapping).map((e) => ({
                value: e,
                label: eventsMapping[e as EventType].description || e.split('_').join(' '),
              })) ?? []}
              onChange={(newValue) => {
                form.setFieldValue(`type`, newValue as EventType);
              }}
              value={form.values.type}
              w={'100%'}
            />

            <Textarea
              label={t('description')}
              required
              w={'100%'}
              autosize={true}
              minRows={1}
              {...form.getInputProps('description')}
            />

            {form.values.type === EventType.COMMAND_USED && <Select
                label={t('triggerCommand')}
                searchable={true}
                data={commandList.data?.map((c) => ({
                  value: c.id,
                  label: c.name,
                })) ?? []}
                onChange={(newValue) => {
                  form.setFieldValue(`commandId`, newValue);
                }}
                value={form.values.commandId}
                w={'100%'}
            />}

            {form.values.type === EventType.REDEMPTION_CREATED && <Select
                label={t('triggerReward')}
                placeholder="..."
                searchable
                itemComponent={RewardItem}
                dropdownPosition={'bottom'}
                allowDeselect
                data={rewards}
                {...form.getInputProps('rewardId')}
                w={'100%'}
            />}

            <Flex direction={'column'} gap={'sm'}>
              <Text mb={5}>{t('availableVariables')}</Text>
              {eventsMapping[form.values.type]?.availableVariables?.map((variable) =>
                  <Text size={'sm'}>
                    <CopyButton value={`{${variable}}`}>
                      {({ copied, copy }) => (
                        <Code
                          onClick={copy}
                          style={{ cursor:'pointer' }}
                        >
                          {copied ? 'Copied' : `{${variable}}`}
                        </Code>
                      )}
                    </CopyButton>
                    {' '} {t(`variables.${variable}`)}
                  </Text>,
              )}
            </Flex>

            <DragDropContext
              onDragEnd={({ destination, source }) =>
                form.reorderListItem('operations', {
                  from: source.index,
                  to: destination!.index,
                })
              }
            >
              <Droppable droppableId="responses" direction="vertical">
                {(provided) => (
                  <div {...provided.droppableProps} ref={provided.innerRef} style={{ width: '100%' }}>
                    {form.values.operations?.map((operation, index) => (
                      <>
                      <Draggable key={index} index={index} draggableId={index.toString()}>
                        {(provided) => (
                          <div key={index} className={cardClasses.classes.root} {...provided.dragHandleProps} ref={provided.innerRef}>
                            <div className={cardClasses.classes.label}>
                              <Flex gap={'xs'}>
                                <ActionIcon variant={'default'} onClick={() => form.removeListItem('operations', index)}>
                                  <IconX />
                                </ActionIcon>
                              </Flex>
                            </div>

                            <Card
                              shadow="sm"
                              p="lg"
                              radius="md"
                              withBorder {...provided.draggableProps}
                              style={{ ...provided.draggableProps.style, position: 'static' }}
                            >
                              <Card.Section p={'xs'} withBorder pt={20}>
                                <Select
                                  searchable={true}
                                  data={Object.keys(OperationType).map(t => ({
                                    value: t,
                                    label: operationMapping[t as OperationType]?.description || t,
                                  }))}
                                  onChange={(newValue) => {
                                    form.setFieldValue(`operations.${index}.type`, newValue);
                                  }}
                                  value={form.values.operations[index]?.type}
                                />
                              </Card.Section>
                              {operationMapping[operation.type].haveInput && <Card.Section p='sm'>
                                  <TextInput
                                    label={t('operations.input')}
                                    required
                                    {...form.getInputProps(`operations.${index}.input`)}
                                  />
                                  {form.values.operations[index - 1]
                                    && operationMapping[form.values.operations[index - 1].type].producedVariables
                                    && <Flex direction={'column'} mt={5}>
                                          <Text size={'sm'}>Available variables from prev operation:</Text>
                                          <Flex direction={'row'}>
                                            {operationMapping[form.values.operations[index - 1].type].producedVariables!.map(v => <CopyButton value={`{prevOperation.${v}}`}>
                                              {({ copied, copy }) => (
                                                <Text
                                                  onClick={copy}
                                                  style={{ cursor:'pointer' }}
                                                  size={'xs'}
                                                >
                                                  {copied ? 'Copied' : `{prevOperation.${v}}`}
                                                </Text>
                                              )}
                                            </CopyButton>)}
                                          </Flex>
                                      </Flex>}
                              </Card.Section>}
                              <Card.Section p='sm' withBorder>
                                <NumberInput
                                  label={t('operations.delay')}
                                  {...form.getInputProps(`operations.${index}.delay`)}
                                />
                                <NumberInput
                                  label={t('operations.repeat')}
                                  {...form.getInputProps(`operations.${index}.repeat`)}
                                />
                              </Card.Section>
                            </Card>

                        </div>
                        )}
                      </Draggable>
                        {index < form.values.operations.length-1 &&
                            <Center w={'100%'} mt={10} mb={10}>
                                <IconArrowBigDownLines size={30} />
                            </Center>
                        }
                        </>
                    ))}

                    {provided.placeholder}
                  </div>
                )}
              </Droppable>
            </DragDropContext>


            <Center w={'100%'}>
              <Button variant={'light'} onClick={() => {
                form.setFieldValue('operations', [...form.values.operations, {
                  type: 'SEND_MESSAGE' as OperationType,
                  id: '',
                  delay: 0,
                  input: '',
                  eventId: '',
                  repeat: 1,
                  order: form.values.operations.length,
                }]);
              }}>
                <IconPlus size={30} />
                New
              </Button>
            </Center>
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
