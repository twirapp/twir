import {
  ActionIcon,
  Alert,
  Anchor,
  Button,
  Flex,
  Grid,
  PasswordInput,
  Text,
  TextInput,
  Tooltip,
} from '@mantine/core';
import { IconDeviceFloppy, IconInfoCircle, IconLink } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { useEffect, useState } from 'react';

import { IntegrationCard } from './card';

import { useDonatelloIntegration } from '@/services/api/integrations/donatello';
import { useValorantIntegration } from '@/services/api/integrations/valorant';

export const ValorantIntegration: React.FC = () => {
  const manager = useValorantIntegration();
  const { data } = manager.useData();
  const { t: integrationsTranslate } = useTranslation('integrations');
  const { t } = useTranslation('common');
  const update = manager.usePost();

  const [username, setUsername] = useState<string>();

  useEffect(() => {
    if (typeof data?.username !== 'undefined') {
      setUsername(data.username);
    }
  }, [data]);

  async function save() {
    if (typeof username == 'undefined') return;
    await update.mutateAsync({ username });
  }

  return (
    <IntegrationCard
      title="Valorant"
      header={
        <Flex direction="row" gap="sm">
          <Button
            compact
            leftIcon={<IconDeviceFloppy />}
            variant="outline"
            color="green"
            onClick={save}
          >
            {t('save')}
          </Button>
        </Flex>
      }
    >
      <Grid align="flex-end">
        <Grid.Col span={12}>
          <TextInput
            label="Valorant username"
            value={username}
            onChange={(v) => setUsername(v.currentTarget.value)}
            placeholder={'Name#tag'}
          />
        </Grid.Col>
      </Grid>

      <Alert color={'lime'} icon={<IconInfoCircle />} mt={5}>
        <Text dangerouslySetInnerHTML={{ __html: integrationsTranslate('info.valorant') }} />
      </Alert>
    </IntegrationCard>
  );
};
