import { ActionIcon, Avatar, Badge, Button, Flex, Group, Switch, Table, Text, TextInput } from '@mantine/core';
import { useDebouncedState, useViewportSize } from '@mantine/hooks';
import { IconPencil, IconSearch, IconTrash } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useState } from 'react';

import { confirmDelete } from '@/components/confirmDelete';
import { GreetingDrawer } from '@/components/greetings/drawer';
import { type Greeting, greetingsManager } from '@/services/api';


// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
    props: {
        ...(await serverSideTranslations(locale, ['greetings', 'layout'])),
    },
});

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableGreeting, setEditableGreeting] = useState<Greeting | undefined>();
  const { t } = useTranslation('greetings');
  const viewPort = useViewportSize();

  const { useGetAll, useDelete, usePatch } = greetingsManager();
  const { data: greetings } = useGetAll();
  const patcher = usePatch();
  const deleter = useDelete();

  const [searchInput, setSearchInput] = useDebouncedState('', 200);

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Group>
          <Text size="lg">{t('title')}</Text>
          <TextInput
            placeholder={'search...'}
            rightSection={<IconSearch size={18} />}
            onChange={(event) => setSearchInput(event.target.value)}
          />
        </Group>

        <Button
          color="green"
          onClick={() => {
            setEditableGreeting(undefined);
            setEditDrawerOpened(true);
          }}
        >
            {t('create')}
        </Button>
      </Flex>
      <Table>
        <thead>
          <tr>
            <th style={{ width:50 }}></th>
            <th>{t('userName')}</th>
            {viewPort.width > 550 && <th>{t('message')}</th>}
            <th>{t('table.head.status')}</th>
            <th>{t('table.head.actions')}</th>
          </tr>
        </thead>
        <tbody>
          {greetings &&
            greetings
              .filter((g) => g.userName.includes(searchInput))
              .map((greeting, idx) => (
              <tr key={greeting.id}>
                <td>
                  <Avatar src={greeting.avatar} style={{ borderRadius:111 }} />
                </td>
                <td>
                  <Badge>{greeting.userName}</Badge>
                </td>
                {viewPort.width > 550 && <td>
                    <Badge color="cyan">{greeting.text}</Badge>
                </td>}
                <td>
                  <Switch
                    checked={greeting.enabled}
                    onChange={(event) => {
                      patcher.mutate({ id: greeting.id, data: { enabled: event.currentTarget.checked } });
                    }}
                  />
                </td>
                <td>
                <Flex direction="row" gap="xs">
                    <ActionIcon
                      onClick={() => {
                        setEditableGreeting(greetings[idx] as any);
                        setEditDrawerOpened(true);
                      }}
                      variant="filled"
                      color="blue"
                    >
                      <IconPencil size={14} />
                    </ActionIcon>

                    <ActionIcon
                      onClick={() =>
                        confirmDelete({
                          onConfirm: () => deleter.mutate(greeting.id),
                        })
                      }
                      variant="filled"
                      color="red"
                    >
                      <IconTrash size={14} />
                    </ActionIcon>
                  </Flex>
                </td>
              </tr>
            ))}
        </tbody>
      </Table>

      <GreetingDrawer
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        greeting={editableGreeting}
      />
    </div>
  );
}
