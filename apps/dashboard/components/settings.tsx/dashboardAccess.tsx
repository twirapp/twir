import { ActionIcon, Alert, Avatar, Box, NavLink, Popover, ScrollArea, Text } from '@mantine/core';
import { IconPlus, IconShieldCheck, IconX } from '@tabler/icons';
import { useState } from 'react';

import { confirmDelete } from '../confirmDelete';
import { DashboardAccessDrawer } from '../settings/dashboardAcessDrawer';
import { SettingsCard } from './card';

import { useDashboardAccess } from '@/services/api/dashboardAcess';

export const DashboardAccess: React.FC = () => {
  const manager = useDashboardAccess();
  const [createDrawerOpened, setCreateDrawerOpened] = useState(false);

  const { data } = manager.getAll();

  return (
    <div>
      <SettingsCard
        title="Dashboard Access"
        icon={IconShieldCheck}
        header={
          <ActionIcon onClick={() => setCreateDrawerOpened(true)}>
            <IconPlus />
          </ActionIcon>
        }
      >
        <Box component={ScrollArea} sx={{ width: '100%' }}>
          {!data?.length && <Alert>No users added</Alert>}
          {!!data?.length &&
            data.map((d) => (
              <NavLink
                key={d.id}
                label={d.twitchUser.login}
                description={d.twitchUser.display_name}
                icon={<Avatar size={35} src={d.twitchUser.profile_image_url} />}
                rightSection={<IconX size={20} stroke={1.5} />}
                onClick={() => {
                  confirmDelete({
                    onConfirm: () => manager.delete(d.id),
                  });
                }}
                sx={{ width: '100%' }}
              />
            ))}
        </Box>
      </SettingsCard>

      <DashboardAccessDrawer opened={createDrawerOpened} setOpened={setCreateDrawerOpened} />
    </div>
  );
};
