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
import { useState } from 'react';

import { CommandDrawer } from '../components/commands/drawer';

import { confirmDelete } from '@/components/confirmDelete';
import { useCommands, useDeleteCommand } from '@/services/api';

type Module = keyof typeof CommandModule;

export default function Commands() {
  const [activeTab, setActiveTab] = useState<Module | null>('CUSTOM');
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableCommand, setEditableCommand] = useState<ChannelCommand | undefined>();
  const viewPort = useViewportSize();

  const { data: commands } = useCommands();
  const deleteCommand = useDeleteCommand();

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Text size="lg">Commands</Text>
        <Button
          color="green"
          onClick={() => {
            setEditableCommand(undefined);
            setEditDrawerOpened(true);
          }}
        >
          Create
        </Button>
      </Flex>
      <Tabs onTabChange={(n) => setActiveTab(n as Module)}>
        <Tabs.List>
          <Tabs.Tab value="CUSTOM" icon={<IconPencilPlus size={14} />}>
            Custom
          </Tabs.Tab>
          <Tabs.Tab value="CHANNEL" icon={<IconUser size={14} />}>
            Channel
          </Tabs.Tab>

          <Tabs.Tab value="MODERATION" icon={<IconSword size={14} />}>
            Moderation
          </Tabs.Tab>
          <Tabs.Tab value="MANAGE" icon={<IconClipboardCopy size={14} />}>
            Manage
          </Tabs.Tab>
          <Tabs.Tab value="DOTA">Dota</Tabs.Tab>
        </Tabs.List>
      </Tabs>
      <Table>
        <thead>
          <tr>
            <th>Name</th>
            {viewPort.width > 550 && (
              <>
                <th>Responses</th>
                <th>Status</th>
              </>
            )}
            <th>Actions</th>
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
                        {command.module != 'CUSTOM' && <Badge>This is built-in command</Badge>}
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
                          onChange={(event) => (command.enabled = event.currentTarget.checked)}
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
                              onConfirm: () => deleteCommand(command.id),
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
