import { ActionIcon, Badge, Button, Flex, Switch, Table, Tabs, Text } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import {
  IconClipboardCopy,
  IconPencilPlus,
  IconSword,
  IconUser,
  IconTrash,
  IconPencil,
} from '@tabler/icons';
import type { ChannelCommand, CommandModule } from '@tsuwari/typeorm/entities/ChannelCommand';
import { useTranslation } from 'next-i18next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useState } from 'react';

import { CommandDrawer } from '../components/commands/drawer';

import { confirmDelete } from '@/components/confirmDelete';
import { useCommandManager } from '@/services/api';

type Module = keyof typeof CommandModule;

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

  const manager = useCommandManager();
  const { data: commands } = manager.getAll();

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Text size="lg">{t('title')}</Text>
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
      <Tabs onTabChange={(n) => setActiveTab(n as Module)}>
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
            {viewPort.width > 550 && (
              <>
                <th>{t('responses')}</th>
                <th>{t('table.head.status')}</th>
              </>
            )}
            <th>{t('table.head.actions')}</th>
          </tr>
        </thead>

        <tbody>
          {commands &&
            commands.length &&
            commands
              .filter((c) => c.module === activeTab)
              .map((command) => (
                <tr key={command.id}>
                  <td>
                    <Badge>{command.name}</Badge>
                  </td>
                  {viewPort.width > 550 && (
                    <>
                      <td>
                        {command.module != 'CUSTOM' && <Badge>{t("builtInBadge")}</Badge>}
                        {command.module === 'CUSTOM' &&
                          (command.responses?.map((r, i) => (
                            <p key={i} style={{ margin: 0 }}>
                              {r.text}
                            </p>
                          )) || <Badge>No Response</Badge>)}
                      </td>
                      <td>
                        <Switch
                          checked={command.enabled}
                          onChange={(event) => {
                            manager.patch(command.id, { enabled: event.currentTarget.checked });
                          }}
                        />
                      </td>
                    </>
                  )}
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
                              onConfirm: () => manager.delete(command.id),
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
