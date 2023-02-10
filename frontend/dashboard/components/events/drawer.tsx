import {
  ActionIcon,
  Button, Card, Divider,
  Drawer,
  Flex,
  ScrollArea,
  TextInput,
  Center,
  useMantineTheme, Select, NumberInput, createStyles,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { IconArrowBigDownLines, IconHandFinger, IconPlus, IconX } from '@tabler/icons';
import { Event, EventType } from '@tsuwari/typeorm/entities/events/Event';
import { EventOperation, OperationType } from '@tsuwari/typeorm/entities/events/EventOperation';
import { useEffect } from 'react';

import { operationMapping } from '@/components/events/mapping';

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
  const viewPort = useViewportSize();

  const cardClasses = useStyles();

  useEffect(() => {
    form.reset();
    if (props.event) {
      form.setValues(props.event);
    }
  }, [props.event, props.opened]);

  async function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      console.log(validate.errors);
      return;
    }
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
        <form onSubmit={form.onSubmit((values) => console.log(values))}>
          <Flex direction="column" gap="md" justify="flex-start" align="flex-start" wrap="wrap">
            <TextInput label={'description'} required {...form.getInputProps('description')} />
            {form.values.operations.map((o, operationIndex) =>
              <div className={cardClasses.classes.root}>
                <div className={cardClasses.classes.label}>
                  <Flex gap={'xs'}>
                    <ActionIcon  variant={'default'}>
                      <IconHandFinger />
                    </ActionIcon>

                    <ActionIcon variant={'default'} onClick={() => form.removeListItem('operations', operationIndex)}>
                      <IconX />
                    </ActionIcon>
                  </Flex>
                </div>

                <Card shadow="sm" p="lg" radius="md" withBorder>
                <Card.Section p={'xs'} withBorder pt={20}>
                  <Select
                    searchable={true}
                    data={Object.keys(OperationType).map(t => ({
                      value: t,
                      label: operationMapping[t as OperationType]?.description || t,
                    }))}
                    onChange={(newValue) => {
                      form.setFieldValue(`operations.${operationIndex}.type`, newValue);
                    }}
                    value={form.values.operations[operationIndex]?.type}
                  />
                </Card.Section>
                <Card.Section p='sm'>
                  {operationMapping[o.type].haveInput && <TextInput
                      label={'Input for operation'}
                      required
                      {...form.getInputProps(`operations.${operationIndex}.input`)}
                  />}
                  <NumberInput
                    label={'Delay in seconds before operation executes'}
                    {...form.getInputProps(`operations.${operationIndex}.delay`)}
                  />
                  <NumberInput
                    label={'How many times execute that operation'}
                    {...form.getInputProps(`operations.${operationIndex}.repeat`)}
                  />
                </Card.Section>
              </Card>
                {operationIndex < form.values.operations.length-1 &&
                    <Center w={'100%'} mt={10}>
                      <IconArrowBigDownLines size={30} />
                    </Center>
                }
              </div>,
              )}

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
                Add new operation
              </Button>
            </Center>
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
