import { ActionIcon, Alert, Avatar, Box, NavLink, ScrollArea } from '@mantine/core';
import { IconPlus, IconShieldCheck, IconX } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { useState } from 'react';

import { resolveUserName } from '../../util/resolveUserName';
import { confirmDelete } from '../confirmDelete';
import { SettingsCard } from './card';
import { DashboardAccessDrawer } from './drawer';

import { dashboardAccessManager } from '@/services/api';

export const DashboardAccess: React.FC = () => {
  const [createDrawerOpened, setCreateDrawerOpened] = useState(false);
  const { t } = useTranslation('settings');

  const manager = dashboardAccessManager();
  const deleter = manager.useDelete();
  const { data } = manager.useGetAll();

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
                  label={resolveUserName(d.twitchUser.login, d.twitchUser.display_name)}
                  icon={<Avatar size={35} src={d.twitchUser.profile_image_url} />}
                  rightSection={<IconX size={20} stroke={1.5} />}
                  onClick={() => {
                    confirmDelete({
                      onConfirm: () => deleter.mutate(d.id),
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
