import {
  Text,
  Flex,
  Button,
  Alert,
  PasswordInput,
  Grid,
} from '@mantine/core';
import { IconLogout, IconLogin, IconLink, IconDeviceFloppy, IconInfoCircle } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { useEffect, useState } from 'react';

import { ManualComponent } from './manualComponent';

import { useDonatePayIntegration, useDonationAlerts } from '@/services/api/integrations';

export const DonatePayIntegration: React.FC = () => {
  const manager = useDonatePayIntegration();
  const { data: apiKey } = manager.useData();
  const { t: integrationsTranslate } = useTranslation('integrations');
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
		<ManualComponent
			integrationKey={'donatepay'}
			save={save}
			imageSize={50}
			body={
				<Flex direction={'column'}>
					<PasswordInput
						label="Api key"
						value={key}
						onChange={(v) => setKey(v.currentTarget.value)}
					/>

					<Button
						variant={'light'}
						component={'a'}
						href={'https://donatepay.ru/page/api'}
						target={'_blank'}
						mt={5}
					>
						Get Api Key
					</Button>
				</Flex>
			}
		/>
  );
};
