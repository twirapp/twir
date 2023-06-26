import {
	ActionIcon,
	Anchor,
	Button,
	CopyButton,
	Stepper,
	TextInput,
} from '@mantine/core';
import { showNotification } from '@mantine/notifications';
import { IconCopy } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { Fragment, useCallback, useEffect, useState } from 'react';

import { ManualComponent } from './manualComponent';

import { useDonateStreamIntegration } from '@/services/api/integrations';

export default () => {
	const manager = useDonateStreamIntegration();
	const { data: apiKey } = manager.useData();
	const secretPost = manager.usePost();

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
		<ManualComponent
			integrationKey={'donate.stream'}
			imageSize={50}
			imageScale={1.5}
			body={
				<Stepper active={0} allowNextStepsSelect={false} orientation='vertical'>
					<Stepper.Step
						label='Step 1'
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
									onChange={() => {
									}}
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
						label='Step 2'
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
						label='Step 3'
						description={'Back to donate.stream and click "confirm" button'}
					/>
				</Stepper>
			}
		/>
	);
};
