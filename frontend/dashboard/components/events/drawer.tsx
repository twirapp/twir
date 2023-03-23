import {
  ActionIcon, Alert,
  Button,
  Card,
  Center,
  Checkbox,
  Code,
  CopyButton,
  createStyles, Divider,
  Drawer,
  Flex,
  Grid,
  Group,
  NumberInput,
  ScrollArea,
  Select, Switch,
  Text,
  Textarea, TextInput,
  useMantineTheme,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { IconArrowBigDownLines, IconArrowDown, IconArrowUp, IconPlus, IconTrash, IconX } from '@tabler/icons';
import { Event, EventType } from '@tsuwari/typeorm/entities/events/Event';
import { OperationType } from '@tsuwari/typeorm/entities/events/EventOperation';
import { useTranslation } from 'next-i18next';
import React, { Fragment, useEffect, useState } from 'react';

import { noop } from '../../util/chore';
import { eventsMapping } from './eventsMapping';

import { filtersMapping } from '@/components/events/filtersMapping';
import { operationMapping } from '@/components/events/operationMapping';
import { RewardItem, RewardItemProps } from '@/components/reward';
import {
  commandsManager,
  eventsManager as useEventsManager,
  keywordsManager as useKeywordsManager,
  useRewards,
  variablesManager as useVariablesManager,
} from '@/services/api';
import { useObsModule } from '@/services/api/modules';

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
      onlineOnly: false,
    },
    validate: {
      operations: {
        delay: (v) => v > 1800 ? 'Delay cannot be more then 1800' : null,
        repeat: (v) => v > 10 ? 'Repeat cannot be more then 10' : null,
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

  const obsManager = useObsModule();
  const { data: obsData } = obsManager.useData();

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
      return obsData?.audioSources.map((s) => ({
        label: s,
        value: s,
      })) ?? [];
    }

    if (type == OperationType.OBS_SET_SCENE) {
      return obsData?.scenes.map((s) => ({
        label: s,
        value: s,
      })) ?? [];
    }

    if (type == OperationType.OBS_TOGGLE_SOURCE) {
      return obsData?.sources.map((s) => ({
        label: s,
        value: s,
      })) ?? [];
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
      size={'50%'}
      position="right"
      transition="slide-left"
      overlayColor={theme.colorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[2]}
      overlayOpacity={0.55}
      overlayBlur={3}
    >
      <ScrollArea.Autosize maxHeight={viewPort.height - 120} type="auto" offsetScrollbars={true}>
        <form style={{ minHeight: viewPort.height - 150 }} onSubmit={form.onSubmit((values) => console.log(values))}>
          <Flex direction="column" gap="md" justify="flex-start" align="flex-start" wrap="wrap">
            <Grid>
              <Grid.Col span={6}>
                <Flex direction={'column'}>
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
                    withinPortal={true}
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

                  <Switch
                    mt={10}
                    label={t('onlineOnly')}
                    labelPosition={'right'}
                    w={'100%'}
                    {...form.getInputProps('onlineOnly', { type: 'checkbox' })}
                  />
                </Flex>
              </Grid.Col>
              <Grid.Col span={6}>
                <Flex direction={'column'} gap={'sm'}>
                  <Text mb={5}>{t('availableVariables')}</Text>
                  {eventsMapping[form.values.type]?.availableVariables?.map((variable, i) =>
                    <Text size={'sm'} key={i}>
                      <CopyButton value={`{${variable}}`}>
                        {({ copied, copy }) => (
                          <Code
                            onClick={copy}
                            style={{ cursor: 'pointer' }}
                          >
                            {copied ? 'Copied' : `{${variable}}`}
                          </Code>
                        )}
                      </CopyButton>
                      {' '} {t(`variables.${variable}`)}
                    </Text>,
                  )}
                  <Alert><Text size={'xs'}>{t('availableVariablesDescription')}</Text></Alert>
                </Flex>
              </Grid.Col>
            </Grid>

            {form.values.operations?.map((operation, operationIndex) => (
              <Fragment key={operationIndex}>
                <div className={cardClasses.classes.root}>
                  <div className={cardClasses.classes.label}>
                    <Flex gap={'xs'}>
                      {operationIndex > 0 && <ActionIcon
                          variant={'default'}
                          onClick={() => form.reorderListItem('operations', {
                            from: operationIndex,
                            to: operationIndex - 1,
                          })}
                      >
                          <IconArrowUp/>
                      </ActionIcon>}
                      {(form.values.operations.length > 1 && operationIndex + 1 !== form.values.operations.length) &&
                          <ActionIcon
                              variant={'default'}
                              onClick={() => form.reorderListItem('operations', {
                                from: operationIndex,
                                to: operationIndex + 1,
                              })}
                          >
                              <IconArrowDown/>
                          </ActionIcon>}
                      <ActionIcon variant={'default'} onClick={() => form.removeListItem('operations', operationIndex)}>
                        <IconX/>
                      </ActionIcon>
                    </Flex>
                  </div>

                  <Card
                    shadow="sm"
                    p="lg"
                    radius="md"
                    withBorder
                  >
                    <Card.Section p={'lg'}>
                      <Grid>
                        <Grid.Col span={6} w={'100%'}>
                          <Select
                            label={'Operation'}
                            searchable={true}
                            data={Object.keys(OperationType).map(t => ({
                              value: t,
                              label: operationMapping[t as OperationType]?.description || t,
                              disabled: operationMapping[t as OperationType].dependsOnEvents
                                ? !operationMapping[t as OperationType].dependsOnEvents?.some(e => e === form.values.type)
                                : false,
                            }))}
                            onChange={(newValue) => {
                              form.setFieldValue(`operations.${operationIndex}.type`, newValue);
                            }}
                            value={form.values.operations[operationIndex]?.type}
                            w={'100%'}
                            withinPortal={true}
                          />

                          {(operationMapping[operation.type].haveInput || operationMapping[operation.type].producedVariables || operationMapping[operation.type].additionalValues) &&
                              <Fragment>
                                {operationMapping[operation.type].haveInput && <Textarea
                                    label={t(`operations.inputDescription.${operation.type}`, t('operations.input'))}
                                    required
                                    autosize={true}
                                    minRows={1}
                                    {...form.getInputProps(`operations.${operationIndex}.input`)}
                                    w={'100%'}
                                />}
                                {form.values.operations && form.values.operations[operationIndex - 1]
                                  && operationMapping[form.values.operations[operationIndex - 1].type].producedVariables
                                  && <Flex direction={'column'}>
                                        <Text size={'sm'}>Available variables from prev operation:</Text>
                                        <Flex direction={'row'}>
                                          {operationMapping[form.values.operations[operationIndex - 1].type].producedVariables!.map((v, i) =>
                                            <CopyButton value={`{prevOperation.${v}}`}>
                                              {({ copied, copy }) => (
                                                <Text
                                                  onClick={copy}
                                                  style={{ cursor: 'pointer' }}
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
                                      {...form.getInputProps(`operations.${operationIndex}.useAnnounce`, { type: 'checkbox' })}
                                  />}
                                  {v === 'timeoutTime' && <NumberInput
                                      label={t('operations.additionalValues.timeoutTime')}
                                      {...form.getInputProps(`operations.${operationIndex}.timeoutTime`)}
                                      w={'100%'}
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
                                          {...form.getInputProps(`operations.${operationIndex}.target`)}
                                          w={'100%'}
                                          withinPortal={true}
                                      />}
                                  {v === 'target' && operation.type.startsWith('OBS') && <Select
                                      label={'OBS Target'}
                                      searchable={true}
                                      data={getObsSourceByOperationType(operation.type)}
                                      {...form.getInputProps(`operations.${operationIndex}.target`)}
                                      w={'100%'}
                                      withinPortal={true}
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
                                          {...form.getInputProps(`operations.${operationIndex}.target`)}
                                          w={'100%'}
                                          withinPortal={true}
                                      />}
                                </Group>)}
                              </Fragment>}
                        </Grid.Col>
                        <Grid.Col span={6}>
                          <NumberInput
                            label={t('operations.delay')}
                            {...form.getInputProps(`operations.${operationIndex}.delay`)}
                          />
                          <NumberInput
                            label={t('operations.repeat')}
                            {...form.getInputProps(`operations.${operationIndex}.repeat`)}
                          />
                        </Grid.Col>
                      </Grid>

                      <Divider mt={5}/>
                      <Text mt={5} size={'lg'}>Filters</Text>

                      {operation.filters?.map((f, filterIndex) => <Grid columns={20} key={filterIndex}>
                        <Grid.Col span={6}>
                          <TextInput
                            {...form.getInputProps(`operations.${operationIndex}.filters.${filterIndex}.left`)}
                            placeholder={'Enter value'}
                          />
                        </Grid.Col>
                        <Grid.Col span={6}>
                          <Select
                            searchable
                            placeholder={'Select filter'}
                            data={Object.entries(filtersMapping).map(f => ({
                                label: f[1].description,
                                value: f[0],
                              }),
                            )}
                            withinPortal={true}
                            onChange={(v) => {
                              form.setFieldValue(
                                `operations.${operationIndex}.filters.${filterIndex}.right`,
                                '',
                              );
                              form.setFieldValue(
                                `operations.${operationIndex}.filters.${filterIndex}.left`,
                                '',
                              );
                              form.setFieldValue(
                                `operations.${operationIndex}.filters.${filterIndex}.type`,
                                v,
                              );
                            }}
                            value={f.type}
                          />
                        </Grid.Col>
                        <Grid.Col span={6}>
                          <TextInput
                            {...form.getInputProps(`operations.${operationIndex}.filters.${filterIndex}.right`)}
                            placeholder={'Enter value'}
                            disabled={filtersMapping[f.type].withoutRight}
                          />
                        </Grid.Col>
                        <Grid.Col span={2}>
                          <Button
                            variant={'light'}
                            color={'red'}
                            onClick={() => {
                              form.removeListItem(`operations.${operationIndex}.filters`, filterIndex);
                            }}
                          >
                            <IconTrash />
                          </Button>
                        </Grid.Col>
                      </Grid>)}

                      <Button mt={10} variant={'light'} size={'xs'} onClick={() => {
                        form.insertListItem(`operations.${operationIndex}.filters`, {
                          right: '',
                          left: '',
                          type: 'CONTAINS',
                        });
                      }}
                      >
                        <IconPlus size={30}/>
                        New filter
                      </Button>
                    </Card.Section>
                  </Card>

                </div>
                {operationIndex < form.values.operations.length - 1 &&
                    <Center w={'100%'} mt={10} mb={10}>
                        <IconArrowBigDownLines size={30}/>
                    </Center>
                }
              </Fragment>
            ))}

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
                  filters: [],
                }]);
              }}
              >
                <IconPlus size={30}/>
                New
              </Button>
            </Center>
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
