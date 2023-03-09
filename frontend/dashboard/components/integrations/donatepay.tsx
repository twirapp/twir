import {
  Group,
  Avatar,
  Text,
  Flex,
  Button,
  Alert,
  TextInput,
  Tooltip,
  ActionIcon,
  Anchor,
  PasswordInput,
  Grid,
} from '@mantine/core';
import { IconLogout, IconLogin, IconLink, IconDeviceFloppy } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { useEffect, useState } from 'react';

import { IntegrationCard } from './card';

import { useDonatePayIntegration, useDonationAlerts } from '@/services/api/integrations';

export const DonatePayIntegration: React.FC = () => {
  const manager = useDonatePayIntegration();
  const { data: apiKey } = manager.useData();
  const { t } = useTranslation('common');
  const update = manager.usePost();

  const [key, setKey] = useState<string>();

  useEffect(() => {
    if (typeof apiKey !== 'undefined') {
      setKey(apiKey);
    }
  }, [apiKey]);

  async function save() {
    if (typeof key === 'undefined') return;
    await update.mutateAsync({ apiKey: key });
  }

  return (
    <IntegrationCard
      title="DonatePay"
      header={
        <Flex direction="row" gap="sm">
          <Button compact leftIcon={<IconDeviceFloppy />} variant="outline" color="green" onClick={save}>
            {t('save')}
          </Button>
        </Flex>
      }
    >
      <Grid align="flex-end" gutter="xs" justify={'space-beetwen'}>
        <Grid.Col span={9}>
          <PasswordInput
            label='Api key'
            value={key}
            onChange={(v) => setKey(v.currentTarget.value)}
          />
        </Grid.Col>

        <Grid.Col span={'auto'}>
          <Button
            variant={'light'}
            component={'a'}
            href={'https://donatepay.ru/page/api'}
            target={'_blank'}>
            Get Api Key
          </Button>
        </Grid.Col>
      </Grid>
    </IntegrationCard>
  );
};
