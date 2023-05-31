import { ActionIcon, Alert, Anchor, CopyButton, Stepper, Text, TextInput } from '@mantine/core';
import { showNotification } from '@mantine/notifications';
import { IconCopy, IconDeviceFloppy, IconInfoCircle, IconLink } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { Fragment, useEffect, useState } from 'react';

import { IntegrationCard } from './card';

import { useDonatelloIntegration } from '@/services/api/integrations/donatello';

export default () => {
	const manager = useDonatelloIntegration();
	const { data: apiKey } = manager.useData();
	const { t: integrationsTranslate } = useTranslation('integrations');
	const update = manager.usePost();

	const [key, setKey] = useState<string>();

	const siteUrl = `${window.location.origin}/api/webhooks/integrations/donatello`;

	useEffect(() => {
		if (apiKey) {
			setKey(apiKey);
		}
	}, [apiKey]);

	async function save() {
		if (typeof key === 'undefined') return;
		await update.mutateAsync({ apiKey: key });
	}

	return (
		<IntegrationCard title="Donatello">
			<Stepper active={0} allowNextStepsSelect={false} orientation="vertical">
				<Stepper.Step
					label="Step 1"
					description={
						<Fragment>
							Go to{' '}
							<Anchor href={'https://donatello.to/panel/settings'} target={'_blank'}>
								https://donatello.to/panel/settings
							</Anchor>{' '}
							and scroll to "Вихідний API" section
						</Fragment>
					}
				/>
				<Stepper.Step
					label="Step 2"
					description={
						<Fragment>
							Copy api key and paste into "Api Key" input
							<TextInput
								mt={5}
								value={apiKey}
								onChange={() => {}}
								variant={'filled'}
								rightSection={
									<CopyButton value={apiKey!}>
										{({ copy }) => (
											<ActionIcon
												onClick={() => {
													copy();
													showNotification({
														color: 'green',
														title: 'Copied',
														message: 'Copied to clipboard',
													});
												}}
											>
												<IconCopy />
											</ActionIcon>
										)}
									</CopyButton>
								}
							/>
						</Fragment>
					}
				/>
				<Stepper.Step
					label="Step 3"
					description={
						<Fragment>
							Copy link and paste into link field
							<TextInput
								mt={5}
								value={siteUrl}
								onChange={() => {}}
								variant={'filled'}
								rightSection={
									<CopyButton value={siteUrl!}>
										{({ copy }) => (
											<ActionIcon
												onClick={() => {
													copy();
													showNotification({
														color: 'green',
														title: 'Copied',
														message: 'Copied to clipboard',
													});
												}}
											>
												<IconCopy />
											</ActionIcon>
										)}
									</CopyButton>
								}
							/>
						</Fragment>
					}
				/>
			</Stepper>

			<Alert color={'lime'} icon={<IconInfoCircle />} mt={5}>
				<Text dangerouslySetInnerHTML={{ __html: integrationsTranslate('info.donations') }} />
			</Alert>
		</IntegrationCard>
	);
};
