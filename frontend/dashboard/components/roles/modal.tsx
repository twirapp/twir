import { Modal, Button, Group, Flex, Text, TextInput, Checkbox, Grid, Tabs } from '@mantine/core';
import { useForm } from '@mantine/form';
import { IconManualGearbox, IconUsers } from '@tabler/icons';
import { ChannelRole } from '@tsuwari/typeorm/entities/ChannelRole';
import { RolePermissionEnum } from '@tsuwari/typeorm/entities/RoleFlag';
import { useCallback, useEffect } from 'react';

import { chunk } from '../../util/chunk';

type Props = {
  opened: boolean;
  role: ChannelRole | undefined;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const RolesModal: React.FC<Props> = (props) => {
  const form = useForm<{
    name: string,
    flags: string[],
    users: string[],
  }>({
    initialValues: {
      name: '',
      flags: [],
      users: [],
    },
  });

  useEffect(() => {
    if (!props.role) return;
    form.setFieldValue('name', props.role.name || '');

    const flags = props.role.permissions?.map((permission) => permission?.flag!.flag as string) ?? [];

    form.setFieldValue('flags', flags);
  }, [props.role]);

  const createCheckboxes = useCallback(() => {
    if (!props.role) return (<></>);

    const checkboxes = Object.values(RolePermissionEnum).map((permission) => {
      const permissionName = permission.replace(/_/g, ' ').toLowerCase();
      const text = permissionName.split(' ').map((word) => {
        return word.charAt(0).toUpperCase() + word.slice(1);
      }).join(' ');

      return <Checkbox
        label={text}
        checked={form.values.flags.includes(permission)}
        onChange={(e) => {
          if (e.target.checked) {
            form.setFieldValue('flags', [...form.values.flags, permission]);
          } else {
            form.setFieldValue('flags', form.values.flags.filter((flag) => flag !== permission));
          }
        }}
      />;
      return;
    });

    const adminitratorCheckbox = checkboxes[0]!;
    checkboxes.splice(0, 1);

    const chunks = chunk(checkboxes, 2);
    return [
      <Grid.Col span={12}>{adminitratorCheckbox}</Grid.Col>,
        chunks.map(c => {
          return (
            <>
            <Grid.Col span={6}>
              {c[0]}
            </Grid.Col>
            <Grid.Col span={6}>
              {c[1]}
            </Grid.Col>
            </>
          );
        }),
    ];
  }, [props.role, form.values.flags]);

  return (
    <Modal
      opened={props.opened}
      onClose={() => props.setOpened(false)}
      title={'Edit role'}
      size={'lg'}
    >
      <Tabs defaultValue="settings">
        <Tabs.List grow>
          <Tabs.Tab value="settings" icon={<IconManualGearbox size={14} />}>Settings</Tabs.Tab>
          <Tabs.Tab value="users" icon={<IconUsers size={14} />}>Users</Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="settings" pt="xs">
          <TextInput label={'Name'} {...form.getInputProps('name')}/>

          <Grid mt={10} >
            {createCheckboxes()}
          </Grid>
        </Tabs.Panel>

        <Tabs.Panel value="users" pt="xs">
          <Text>Users</Text>
        </Tabs.Panel>
      </Tabs>

    </Modal>
  );
};