import { ActionIcon, Alert, Avatar, Box, NavLink, ScrollArea } from '@mantine/core';
import { IconPlus, IconShieldCheck, IconX } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { useState } from 'react';

import { confirmDelete } from '../confirmDelete';
import { SettingsCard } from './card';
import { DashboardAccessDrawer } from './drawer';

import { dashboardAccessManager } from '@/services/api';

export const DashboardAccess: React.FC = () => {
  const [createDrawerOpened, setCreateDrawerOpened] = useState(false);
  const { t } = useTranslation('settings');

  const manager = dashboardAccessManager();
  const { data } = manager.getAll;

  return (
    <div>
      <SettingsCard
        title={t('dashboardAccess.title')}
        icon={IconShieldCheck}
        header={
          <ActionIcon onClick={() => setCreateDrawerOpened(true)}>
            <IconPlus />
          </ActionIcon>
        }
      >
        <ScrollArea type={'always'}>
          <Box sx={{ width: '100%', height: 400 }}>
            {!data?.length && <Alert>{t('dashboardAccess.emptyAlert')}</Alert>}
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
                      onConfirm: () => manager.delete.mutate(d.id),
                    });
                  }}
                  sx={{ width: '100%' }}
                />
              ))}
          </Box>
        </ScrollArea>

      </SettingsCard>

      <DashboardAccessDrawer opened={createDrawerOpened} setOpened={setCreateDrawerOpened} />
    </div>
  );
};
