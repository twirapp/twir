import { ActionIcon, Badge, Button, Flex, Group, Switch, Table, Tabs, Text, TextInput } from '@mantine/core';
import { useDebouncedValue, useViewportSize } from '@mantine/hooks';
import {
  IconClipboardCopy,
  IconPencilPlus,
  IconSword,
  IconUser,
  IconTrash,
  IconPencil, IconSearch,
} from '@tabler/icons';
import type { ChannelCommand, CommandModule } from '@tsuwari/typeorm/entities/ChannelCommand';
import { useTranslation } from 'next-i18next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useState } from 'react';

import { CommandDrawer } from '@/components/commands/drawer';
import { confirmDelete } from '@/components/confirmDelete';
import { commandsManager } from '@/services/api';

type Module = keyof typeof CommandModule;

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['commands', 'layout'])),
  },
});

export default function Commands() {
  const [activeTab, setActiveTab] = useState<Module | null>('CUSTOM');
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableCommand, setEditableCommand] = useState<ChannelCommand | undefined>();
  const viewPort = useViewportSize();
  const { t } = useTranslation('commands');

  const { useGetAll, usePatch, useDelete } = commandsManager();
  const patcher = usePatch();
  const deleter = useDelete();
  const { data: commands } = useGetAll();

  const [searchInput, setSearchInput] = useState('');
  const [debouncedSearchInput] = useDebouncedValue(searchInput, 200);

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Group>
          <Text size="lg">{t('title')}</Text>
          <TextInput
            placeholder={'search...'}
            rightSection={<IconSearch size={18} />}
            onChange={event => setSearchInput(event.target.value)}
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
      <Tabs defaultValue={'CUSTOM'} onTabChange={(n) => setActiveTab(n as Module)}>
        <Tabs.List>
          <Tabs.Tab value="CUSTOM" icon={<IconPencilPlus size={14} />}>
            {t('tabs.custom')}
          </Tabs.Tab>
          <Tabs.Tab value="CHANNEL" icon={<IconUser size={14} />}>
            {t('tabs.channel')}
          </Tabs.Tab>

          <Tabs.Tab value="MODERATION" icon={<IconSword size={14} />}>
            {t('tabs.moderation')}
          </Tabs.Tab>
          <Tabs.Tab value="MANAGE" icon={<IconClipboardCopy size={14} />}>
            {t('tabs.manage')}
          </Tabs.Tab>
          <Tabs.Tab value="DOTA">{t('tabs.dota')}</Tabs.Tab>
        </Tabs.List>
      </Tabs>
      <Table>
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
                const isActiveTab = c.module === activeTab;
                const nameIncludes = c.name.includes(debouncedSearchInput);
                const aliasesIncludes = c.aliases.some(a => a.includes(debouncedSearchInput));
                return isActiveTab && (nameIncludes || aliasesIncludes);
              })
              .map((command) => (
                <tr key={command.id}>
                  <td>
                    <Badge>{command.name}</Badge>
                  </td>
                  {viewPort.width > 550 && <td>
                    {command.module != 'CUSTOM' && <Badge>{t('builtInBadge')}</Badge>}
                    {command.module === 'CUSTOM' &&
                      (command.responses?.map((r, i) => (
                        <p key={i} style={{ margin: 0 }}>
                          {r.text}
                        </p>
                      )) || <Badge>No Response</Badge>)}
                  </td>}
                  <td>
                    <Switch
                      checked={command.enabled}
                      onChange={() => patcher.mutate({ id: command.id, data: { enabled: !command.enabled } })}
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
}
