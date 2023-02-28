import {
  ActionIcon,
  Button,
  Card,
  Center,
  Checkbox,
  Code,
  CopyButton,
  createStyles,
  Drawer,
  Flex,
  Group,
  NumberInput,
  ScrollArea,
  Select,
  Text,
  Textarea,
  useMantineTheme,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { IconArrowBigDownLines, IconPlus, IconX } from '@tabler/icons';
import { Event, EventType } from '@tsuwari/typeorm/entities/events/Event';
import { OperationType } from '@tsuwari/typeorm/entities/events/EventOperation';
import { useTranslation } from 'next-i18next';
import React, { useEffect, useState } from 'react';
import { DragDropContext, Draggable, Droppable } from 'react-beautiful-dnd';

import { noop } from '../../util/chore';
import { eventsMapping } from './eventsMapping';

import { operationMapping } from '@/components/events/operationMapping';
import { RewardItem, RewardItemProps } from '@/components/reward';
import {
  commandsManager,
  eventsManager as useEventsManager,
  keywordsManager as useKeywordsManager,
  useRewards,
  variablesManager as useVariablesManager,
} from '@/services/api';
import { useObs } from '@/services/obs/hook';

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
      commandId: '',
      rewardId: '',
      keywordId: '',
      description: '',
      operations: [],
      enabled: true,
    },
    validate: {
      operations: {
        delay: (v) => v > 1800 ? 'Delay cannot be more then 1800' : null,
        repeat: (v) => v > 10 ? 'Repeat cannot be more then 10' : null,
        // input: (v, values, path) => {
        //   const pathWithoutInput = path.split('.input')[0];
        //   const operationIndex = pathWithoutInput.at(pathWithoutInput.length - 1);
        //   const operation = values.operations[Number(operationIndex)];
        //   if (!operation) return null;
        //   if (
        //     operation.type === OperationType.OBS_AUDIO_SET_VOLUME
        //     || operation.type === OperationType.OBS_AUDIO_INCREASE_VOLUME
        //     || operation.type === OperationType.OBS_AUDIO_DECREASE_VOLUME
        //   ) {
        //     const convertedValue = Number(v);
        //     if (Number.isNaN(convertedValue)) return 'Incorrect value. Can be from 0 to 20';
        //     if (convertedValue < 0 || convertedValue > 20) {
        //       return 'Volume can be from 0 to 20.';
        //     }
        //   }
        //
        //   return null;
        // },
      },
    },
  });
  const eventsManager = useEventsManager();
  const updater = eventsManager.useCreateOrUpdate();
  const { t } = useTranslation('events');
  const viewPort = useViewportSize();
  const cardClasses = useStyles();
  const [rewards, setRewards] = useState<RewardItemProps[]>([]);

  const commandManager = commandsManager();
  const commandList = commandManager.useGetAll();

  const keywordsManager = useKeywordsManager();
  const { data: keywords } = keywordsManager.useGetAll();

  const variablesManager = useVariablesManager();
  const { data: variables } = variablesManager.useGetAll();

  const rewardsManager = useRewards();
  const { data: rewardsData } = rewardsManager();

  const obsSocket = useObs();

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

  function getObsSourceByOperationType(type: OperationType) {
    if (
      type === OperationType.OBS_TOGGLE_AUDIO
      || type == OperationType.OBS_AUDIO_DECREASE_VOLUME
      || type == OperationType.OBS_AUDIO_INCREASE_VOLUME
      || type == OperationType.OBS_AUDIO_SET_VOLUME
      || type == OperationType.OBS_DISABLE_AUDIO
      || type == OperationType.OBS_ENABLE_AUDIO
    ) {
      return obsSocket.inputs.map((i) => ({
        label: i,
        value: i,
      }));
    }

    if (type == OperationType.OBS_SET_SCENE) {
      return Object.keys(obsSocket.scenes).map((s) => ({
        label: s,
        value: s,
      }));
    }

    if (type == OperationType.OBS_TOGGLE_SOURCE) {
      return Object.entries(obsSocket.scenes)
        .map((pair) => {
          return pair[1].sources.map((source) => ({
            label: `[${pair[0]}] ${source.name}`,
            value: source.name,
          }));
        })
        .flat();
    }

    return [];
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
                label: eventsMapping[e as EventType].description?.toUpperCase() || e.split('_').join(' '),
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
                {...form.getInputProps('commandId')}
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

            {form.values.type === EventType.KEYWORD_MATCHED && <Select
                label={t('triggerKeyword')}
                searchable={true}
                data={keywords?.map((c) => ({
                  value: c.id,
                  label: c.text,
                })) ?? []}
                onChange={(newValue) => {
                  form.setFieldValue(`keywordId`, newValue);
                }}
                value={form.values.keywordId}
                w={'100%'}
            />}

            <Flex direction={'column'} gap={'sm'}>
              <Text mb={5}>{t('availableVariables')}</Text>
              {eventsMapping[form.values.type]?.availableVariables?.map((variable, i) =>
                  <Text size={'sm'} key={i}>
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
                                    disabled: operationMapping[t as OperationType].dependsOnEvents
                                      ? !operationMapping[t as OperationType].dependsOnEvents?.some(e => e === form.values.type)
                                      : false,
                                  }))}
                                  onChange={(newValue) => {
                                    form.setFieldValue(`operations.${index}.type`, newValue);
                                  }}
                                  value={form.values.operations[index]?.type}
                                />
                              </Card.Section>
                              {(operationMapping[operation.type].haveInput || operationMapping[operation.type].producedVariables || operationMapping[operation.type].additionalValues) && <Card.Section p='sm'>
                                  {operationMapping[operation.type].haveInput && <Textarea
                                      label={t(`operations.inputDescription.${operation.type}`, t('operations.input'))}
                                      required
                                      w={'100%'}
                                      autosize={true}
                                      minRows={1}
                                      {...form.getInputProps(`operations.${index}.input`)}
                                  />}
                                  {form.values.operations && form.values.operations[index - 1]
                                    && operationMapping[form.values.operations[index - 1].type].producedVariables
                                    && <Flex direction={'column'} mt={5}>
                                          <Text size={'sm'}>Available variables from prev operation:</Text>
                                          <Flex direction={'row'}>
                                            {operationMapping[form.values.operations[index - 1].type].producedVariables!.map((v, i) => <CopyButton value={`{prevOperation.${v}}`}>
                                              {({ copied, copy }) => (
                                                <Text
                                                  onClick={copy}
                                                  style={{ cursor:'pointer' }}
                                                  size={'xs'}
                                                  key={i}
                                                >
                                                  {copied ? 'Copied' : `{prevOperation.${v}}`}
                                                </Text>
                                              )}
                                            </CopyButton>)}
                                          </Flex>
                                      </Flex>}
                                {operationMapping[operation.type].additionalValues?.map((v, i) => <Group key={i} mt={5}>
                                  {v === 'useAnnounce' && <Checkbox
                                    label={t('operations.additionalValues.useAnnounce')}
                                    labelPosition={'left'}
                                    {...form.getInputProps(`operations.${index}.useAnnounce`, { type: 'checkbox' })}
                                  />}
                                  {v === 'timeoutTime' && <NumberInput
                                    label={t('operations.additionalValues.timeoutTime')}
                                    {...form.getInputProps(`operations.${index}.timeoutTime`)}
                                  />}
                                  {
                                    v === 'target'
                                    && (operation.type === OperationType.CHANGE_VARIABLE
                                      || operation.type === OperationType.INCREMENT_VARIABLE
                                      || operation.type === OperationType.DECREMENT_VARIABLE
                                    )
                                    && <Select
                                    label={'Variable'}
                                    searchable
                                    data={variables?.map(v => ({
                                      value: v.id,
                                      label: v.name,
                                    })) ?? []}
                                    w={'100%'}
                                    {...form.getInputProps(`operations.${index}.target`)}
                                  />}
                                  {v === 'target' && operation.type.startsWith('OBS') && <Select
                                      label={'OBS Target'}
                                      searchable={true}
                                      data={getObsSourceByOperationType(operation.type)}
                                      w={'100%'}
                                      {...form.getInputProps(`operations.${index}.target`)}
                                  />}
                                  {v === 'target' && (
                                    operation.type === OperationType.ALLOW_COMMAND_TO_USER ||
                                    operation.type === OperationType.REMOVE_ALLOW_COMMAND_TO_USER ||
                                    operation.type === OperationType.DENY_COMMAND_TO_USER ||
                                    operation.type === OperationType.REMOVE_DENY_COMMAND_TO_USER) &&
                                      <Select
                                          label={'Command'}
                                          searchable={true}
                                          data={commandList.data?.map((c) => ({
                                            value: c.id,
                                            label: c.name,
                                          })) ?? []}
                                          {...form.getInputProps(`operations.${index}.target`)}
                                          w={'100%'}
                                      />}
                                </Group>)}
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
                  useAnnounce: false,
                  timeoutTime: 600,
                  target: '',
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
