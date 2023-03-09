import { ActionIcon, Anchor, Button, Flex, Grid, PasswordInput, TextInput, Tooltip } from '@mantine/core';
import { IconDeviceFloppy, IconLink } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { useEffect, useState } from 'react';

import { IntegrationCard } from './card';

import { useDonatelloIntegration } from '@/services/api/integrations/donatello';

export const DonatelloIntegration: React.FC = () => {
  const manager = useDonatelloIntegration();
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
      title="Donatello"
      header={
        <Flex direction="row" gap="sm">
          <Button compact leftIcon={<IconDeviceFloppy/>} variant="outline" color="green" onClick={save}>
            {t('save')}
          </Button>
        </Flex>
      }
    >
      <Grid align="flex-end">
        <Grid.Col span={9} >
          <PasswordInput
            label="Api key"
            value={key}
            onChange={(v) => setKey(v.currentTarget.value)}
          />
        </Grid.Col>
        <Grid.Col span={'auto'}>
          <Button
            variant={'light'}
            component={'a'}
            href={'https://donatello.to/panel/doc-api'}
            target={'_blank'}>
            Get Api Key
          </Button>
        </Grid.Col>
      </Grid>
    </IntegrationCard>
  );
};
