import {
	ActionIcon,
	Alert,
	Anchor,
	Button,
	CopyButton,
	Stepper,
	Text,
	TextInput,
} from '@mantine/core';
import { showNotification } from '@mantine/notifications';
import { IconCopy, IconInfoCircle } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { Fragment, useCallback, useEffect, useState } from 'react';

import { IntegrationCard } from './card';

import { useDonateStreamIntegration } from '@/services/api/integrations';

export default () => {
	const manager = useDonateStreamIntegration();
	const { data: apiKey } = manager.useData();
	const secretPost = manager.usePost();
	const { t: integrationsTranslate } = useTranslation('integrations');

	const [key, setKey] = useState<string>('');
	const [secret, setSecret] = useState<string>('');

	useEffect(() => {
		if (typeof apiKey !== 'undefined') {
			setKey(`${window.location.origin}/api/webhooks/integrations/donate-stream/${apiKey}`);
		}
	}, [apiKey]);

	const onSubmit = useCallback(() => {
		if (apiKey && secret) {
			secretPost.mutate({ id: apiKey, secret });
		}
	}, [apiKey, secret]);

	return (
		<IntegrationCard title="Donate.stream">
			<Stepper active={0} allowNextStepsSelect={false} orientation="vertical">
				<Stepper.Step
					label="Step 1"
					description={
						<Fragment>
							Paste that link into input on the{' '}
							<Anchor href={'https://lk.donate.stream/settings/api-key'} target={'_blank'}>
								https://lk.donate.stream/settings/api
							</Anchor>{' '}
							page
							<TextInput
								mt={5}
								value={key}
								onChange={() => {}}
								variant={'filled'}
								rightSection={
									<CopyButton value={key!}>
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
					label="Step 2"
					description={
						<Fragment>
							Paste the{' '}
							<Anchor href={'https://i.imgur.com/OtW97pV.png'} target={'_blank'}>
								secret key
							</Anchor>{' '}
							from page and click button for save
							<TextInput
								mt={5}
								value={secret}
								placeholder={'paste secret here'}
								onChange={(e) => setSecret(e.target.value)}
								variant={'filled'}
							/>
							<Button mt={5} size={'sm'} color={'green'} variant={'light'} onClick={onSubmit}>
								Send
							</Button>
						</Fragment>
					}
				/>
				<Stepper.Step
					label="Step 3"
					description={'Back to donate.stream and click "confirm" button'}
				/>
			</Stepper>

			<Alert color={'lime'} icon={<IconInfoCircle />} mt={5}>
				<Text dangerouslySetInnerHTML={{ __html: integrationsTranslate('info.donations') }} />
			</Alert>
		</IntegrationCard>
	);
};
