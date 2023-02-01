import { ActionIcon, Anchor, Button, Flex, TextInput, Tooltip } from '@mantine/core';
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
      <TextInput
        label="Api key"
        value={key}
        onChange={(v) => setKey(v.currentTarget.value)}
        rightSection={<Tooltip label="Get api key" color="violet" withArrow>
          <Anchor href={'https://donatello.to/panel/doc-api'} target={'_blank'}>
            <ActionIcon>
              <IconLink/>
            </ActionIcon>
          </Anchor>
        </Tooltip>}
      />
    </IntegrationCard>
  );
};
