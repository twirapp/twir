import { Badge, Button, Switch, Table, Tabs } from '@mantine/core';
import { useLocalStorage, useViewportSize } from '@mantine/hooks';
import { IconClipboardCopy, IconPencilPlus, IconSword, IconUser } from '@tabler/icons';
import { Dashboard } from '@tsuwari/shared';
import type { ChannelCommand, CommandModule } from '@tsuwari/typeorm/entities/ChannelCommand';
import { useState } from 'react';
import useSWR from 'swr';

import { CommandDrawer } from '../components/commands/drawer';
import { swrFetcher } from '../services/swrFetcher';

type Module = keyof typeof CommandModule;

export default function Commands() {
  const [activeTab, setActiveTab] = useState<Module | null>('CUSTOM');
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableCommand, setEditableCommand] = useState<ChannelCommand>({} as any);
  const [selectedDashboard] = useLocalStorage<Dashboard>({
    key: 'selectedDashboard',
    serialize: (v) => JSON.stringify(v),
    deserialize: (v) => JSON.parse(v),
  });
  const viewPort = useViewportSize();

  const { data: commands } = useSWR<ChannelCommand[]>(
    selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/commands` : null,
    swrFetcher,
  );

  return (
    <div>
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
                <th>Response</th>
                <th>Status</th>
              </>
            )}
            <th>Actions</th>
          </tr>
        </thead>

        <tbody>
          {commands &&
            commands
              .filter((c) => c.module === activeTab)
              .map((element, idx) => (
                <tr key={element.id}>
                  <td>
                    <Badge>{element.name}</Badge>
                  </td>
                  {viewPort.width > 550 && (
                    <>
                      <td>
                        {element.module != 'CUSTOM' && <Badge>This is built-in command</Badge>}
                        {element.module === 'CUSTOM' &&
                          (element.responses?.map((r, i) => (
                            <p key={i} style={{ margin: 0 }}>
                              {r.text}
                            </p>
                          )) || <Badge>No Response</Badge>)}
                      </td>
                      <td>
                        <Switch
                          checked={element.enabled}
                          onChange={(event) => (element.enabled = event.currentTarget.checked)}
                        />
                      </td>
                    </>
                  )}
                  <td>
                    <Button
                      onClick={() => {
                        setEditableCommand(commands[idx] as any);
                        setEditDrawerOpened(true);
                      }}
                    >
                      Edit
                    </Button>
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
