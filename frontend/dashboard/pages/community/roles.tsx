import { Button, Card, Center, DEFAULT_THEME, Flex, Text } from '@mantine/core';
import { IconPlus } from '@tabler/icons';
import { ChannelRole } from '@tsuwari/typeorm/entities/ChannelRole';
import { NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useState } from 'react';

import { confirmDelete } from '@/components/confirmDelete';
import { RolesModal } from '@/components/roles/modal';
import { useRolesApi } from '@/services/api';

const Roles: NextPage = () => {
  const rolesManager = useRolesApi();
  const { data: roles } = rolesManager.useGetAll();
  const rolesDeleter = rolesManager.useDelete();

  const [modalOpened, setModalOpened] = useState(false);
  const [editableRole, setEditableRole] = useState<ChannelRole | undefined>();

  return (
    <>

      <Flex direction={'column'} align={'center'} gap={'lg'}>
        <Card
          shadow="sm"
          p="lg"
          radius="md"
          withBorder
          w={500}
          onMouseDown={() => {
            setEditableRole(undefined);
            setModalOpened(true);
          }}
          style={{
            cursor: 'pointer',
            backgroundColor: DEFAULT_THEME.colors.gray[7],
          }}
        >
          <Center>
            <IconPlus />
            <Text size={'lg'}>New role</Text>
          </Center>
        </Card>
        {!!roles?.length && roles.map((role) => (
          <Card
            key={role.id}
            shadow="sm"
            p="lg"
            radius="md"
            withBorder
            w={500}
            onMouseDown={() => {
              setEditableRole(role);
              setModalOpened(true);
            }}
            style={{ cursor: 'pointer' }}
          >
            <Card.Section p={'lg'}>
              <Flex direction={'row'} justify={'space-between'}>
                <Text size={'xl'}>{role.name}</Text>
                <Flex direction={'row'} gap={'xs'}>
                  <Button size={'xs'} variant={'light'} onMouseDown={(e) => {
                    e.stopPropagation();

                    setEditableRole(role);
                    setModalOpened(true);
                  }}>
                    Edit
                  </Button>
                  {!role.system &&
                    <Button
                      size={'xs'}
                      variant={'light'}
                      color={'red'}
                      onMouseDown={(e) => {
                        e.stopPropagation();
                        confirmDelete({
                          onConfirm: () => {
                            rolesDeleter.mutate(role.id);
                          },
                        });
                      }}
                    >
                      Remove
                    </Button>
                  }
                </Flex>
              </Flex>
            </Card.Section>
          </Card>
        ))}
      </Flex>
      <RolesModal
        opened={modalOpened}
        setOpened={setModalOpened}
        role={editableRole}
      />
    </>
  );
};

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['community', 'layout'])),
  },
});


export default Roles;