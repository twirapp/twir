import {
  ActionIcon,
  Button, ColorPicker,
  Drawer,
  Flex,
  Menu,
  ScrollArea,
  Switch, Text,
  Textarea,
  TextInput,
  useMantineTheme,
} from '@mantine/core';
import { isNotEmpty, useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { IconVariable } from '@tabler/icons';
import { ChannelCommandGroup } from '@tsuwari/typeorm/entities/ChannelCommandGroup';
import { useTranslation } from 'next-i18next';
import { useEffect } from 'react';

import { noop } from '../../util/chore';

import { commandsGroupManager } from '@/services/api';

type Props = {
  opened: boolean;
  group?: ChannelCommandGroup;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const ChannelCommandGroupDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<ChannelCommandGroup>({
    validate: {
      name: isNotEmpty('Name cannot be empty'),
      color: isNotEmpty('Color cannot be empty'),
    },
    initialValues: {
      id: '',
      channelId: '',
      color: 'rgba(37, 38, 43, 1)',
      name: '',
    },
  });
  const viewPort = useViewportSize();
  const { t } = useTranslation('greetings');

  const group = commandsGroupManager();
  const updater = group.useCreateOrUpdate();

  useEffect(() => {
    form.reset();
    if (props.group) {
      form.setValues(props.group);
    }
  }, [props.group, props.opened]);

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
        <Button size="xs" color="green" onClick={onSubmit}>Save</Button>
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
            <TextInput label={'Name'} required {...form.getInputProps('name')} />
            <Text>Color</Text>
            <ColorPicker format="rgba" value={form.values.color} onChange={(v) => {
              form.setFieldValue('color', v);
            }} />
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
