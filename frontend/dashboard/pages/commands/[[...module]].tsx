import {
  ActionIcon,
  Badge,
  Button,
  Flex,
  Group,
  Switch,
  Table,
  Text,
  TextInput,
} from '@mantine/core';
import { useDebouncedState, useViewportSize } from '@mantine/hooks';
import { IconPencil, IconSearch, IconTrash } from '@tabler/icons';
import { ChannelCommand } from '@tsuwari/typeorm/entities/ChannelCommand';
import { GetServerSideProps, NextPage } from 'next';
import { useTranslation } from 'next-i18next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useRouter } from 'next/router';
import { useState } from 'react';

import { CommandDrawer } from '@/components/commands/drawer';
import { confirmDelete } from '@/components/confirmDelete';
import { commandsManager } from '@/services/api';

export const getServerSideProps: GetServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale!, ['commands', 'layout'])),
  },
});

const Commands: NextPage = () => {
  const router = useRouter();
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableCommand, setEditableCommand] = useState<ChannelCommand | undefined>();
  const viewPort = useViewportSize();
  const { t } = useTranslation('commands');

  const { useGetAll, usePatch, useDelete } = commandsManager();
  const patcher = usePatch();
  const deleter = useDelete();
  const { data: commands } = useGetAll();

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
            setEditableCommand(undefined);
            setEditDrawerOpened(true);
          }}
        >
          {t('create')}
        </Button>
      </Flex>
      <Table style={{ tableLayout: 'fixed', width: '100%' }}>
        <thead>
          <tr>
            <th>{t('name')}</th>
            {viewPort.width > 550 && <th>{t('responses')}</th>}
            <th>{t('table.head.status')}</th>
            <th>{t('table.head.actions')}</th>
          </tr>
        </thead>

        <tbody>
          {commands &&
            commands.length &&
            commands
              .filter((c) => {
                const isActiveTab =
                  c.module.toLowerCase() ===
                  (router.query.module ? router.query.module[0] : 'custom');
                const nameIncludes = c.name.includes(searchInput);
                const aliasesIncludes = c.aliases.some((a) => a.includes(searchInput));
                return isActiveTab && (nameIncludes || aliasesIncludes);
              })
              .map((command) => (
                <tr key={command.id}>
                  <td>
                    <Badge>{command.name}</Badge>
                  </td>
                  {viewPort.width > 550 && (
                    <td>
                      {command.module != 'CUSTOM' && <Badge>{t('builtInBadge')}</Badge>}
                      {command.module === 'CUSTOM' &&
                        (command.responses?.map((r) => (
                          <Text
                            title={r.text!}
                            lineClamp={1}
                            style={{ textOverflow: 'ellipsis', overflow: 'hidden' }}
                          >
                            {r.text}
                          </Text>
                        )) || <Badge>No Response</Badge>)}
                    </td>
                  )}
                  <td>
                    <Switch
                      checked={command.enabled}
                      onChange={() =>
                        patcher.mutate({ id: command.id, data: { enabled: !command.enabled } })
                      }
                    />
                  </td>
                  <td>
                    <Flex direction="row" gap="xs">
                      <ActionIcon
                        onClick={() => {
                          setEditableCommand(commands.find((c) => c.id === command.id)!);
                          setEditDrawerOpened(true);
                        }}
                        variant="filled"
                        color="blue"
                      >
                        <IconPencil size={14} />
                      </ActionIcon>
                      {command.module === 'CUSTOM' && (
                        <ActionIcon
                          onClick={() =>
                            confirmDelete({
                              onConfirm: () => deleter.mutate(command.id),
                            })
                          }
                          variant="filled"
                          color="red"
                        >
                          <IconTrash size={14} />
                        </ActionIcon>
                      )}
                    </Flex>
                  </td>
                </tr>
              ))}
        </tbody>
      </Table>

      <CommandDrawer
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        command={editableCommand}
      />
    </div>
  );
};

export default Commands;
