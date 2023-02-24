import {
  Modal,
  Text,
  TextInput,
  Checkbox,
  Grid,
  Tabs,
  Card,
  Avatar,
  Flex,
  Button,
  ActionIcon,
  Select,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useDebouncedState } from '@mantine/hooks';
import { IconManualGearbox, IconPlus, IconTrash, IconUsers } from '@tabler/icons';
import { ChannelRole } from '@tsuwari/typeorm/entities/ChannelRole';
import { RoleFlags } from '@tsuwari/typeorm/entities/ChannelRole';
import { HelixUserData } from '@twurple/api';
import { useCallback, useEffect, useState } from 'react';

import { chunk } from '../../util/chunk';

import { useRolesUsers, useTwitch } from '@/services/api';

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

  const [newUser, setNewUser] = useState('');

  const usersManager = useRolesUsers();
  const { data: users } = usersManager.useGetAll(props.role?.id || '');

  useEffect(() => {
    if (!props.role) return;
    form.setFieldValue('name', props.role.name || '');
    form.setFieldValue('flags', props.role.permissions);
  }, [props.role]);

  useEffect(() => {
    if (!users?.length) return;

    form.setFieldValue('users', users.map((user) => user.userName));
  }, [users]);

  const createCheckboxes = useCallback(() => {
    if (!props.role) return (<></>);

    const checkboxes = Object.values(RoleFlags).map((permission) => {
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
          <Grid>

            <Grid.Col span={12}>
              <Grid align="center">
                <Grid.Col span={10}>
                  <TextInput
                    placeholder={'Add new user by name'}
                    value={newUser}
                    onChange={(event) => setNewUser(event.currentTarget.value)}
                  />
                </Grid.Col>
                <Grid.Col span={2}>
                  <Button fullWidth variant={'light'}><IconPlus /></Button>
                </Grid.Col>
              </Grid>
            </Grid.Col>
            {users?.map((user, index) => {
              return (
                <Grid.Col span={6}>
                  <Card>
                    <Card.Section p={'lg'}>
                      <Flex direction={'row'} justify={'space-between'} align={'center'}>
                        <Flex direction={'row'} gap={'md'} align={'center'}>
                          <Avatar src={user.userAvatar} radius="xl" />
                          <Text>
                            {user.userDisplayName}
                          </Text>
                        </Flex>
                        <ActionIcon
                          color={'red'}
                          onClick={() => {
                            form.removeListItem('users', index);
                          }}
                        >
                          <IconTrash size={25} />
                        </ActionIcon>
                      </Flex>
                    </Card.Section>
                  </Card>
                </Grid.Col>
              );
            })
            }
          </Grid>
        </Tabs.Panel>
      </Tabs>

    </Modal>
  );
};