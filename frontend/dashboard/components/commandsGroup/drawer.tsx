import {
  Button, ColorInput,
  Drawer,
  Flex,
  ScrollArea,
  TextInput,
  useMantineTheme,
} from '@mantine/core';
import { isNotEmpty, useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { Group } from '@twir/grpc/generated/api/api/commands_group';
import { useTranslation } from 'next-i18next';
import { useEffect } from 'react';

import { noop } from '../../util/chore';

import { colorPickerColors } from '@/components/commandsGroup/colors';
import { useCommandsGroupsManager } from '@/services/api';

type Props = {
  opened: boolean;
  group?: Group;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const ChannelCommandGroupDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<Group>({
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

  const group = useCommandsGroupsManager();
  const updater = group.update;
	const creator = group.create;

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

		if (!form.values.id) {
			await creator.mutateAsync({ group: form.values });
		} else {
			await updater.mutateAsync({
				id: form.values.id,
				group: form.values,
			});
		}

		props.setOpened(false);
		form.reset();
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
			closeOnClickOutside={false}
    >
      <ScrollArea.Autosize maxHeight={viewPort.height - 120} type="auto" offsetScrollbars={true}>
        <form onSubmit={form.onSubmit((values) => console.log(values))}>
          <Flex direction="column" gap="md" justify="flex-start" align="flex-start" wrap="wrap">
            <TextInput label={'Name'} required {...form.getInputProps('name')} />
            <ColorInput
              label="Color for group"
              format="rgba"
              value={form.values.color}
              onChange={(v) => {
                form.setFieldValue('color', v);
              }}
              swatches={colorPickerColors}
              withAsterisk
            />
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
