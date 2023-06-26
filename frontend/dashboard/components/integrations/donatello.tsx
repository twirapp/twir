import { ActionIcon, Alert, Anchor, CopyButton, Stepper, Text, TextInput } from '@mantine/core';
import { showNotification } from '@mantine/notifications';
import { IconCopy, IconDeviceFloppy, IconInfoCircle, IconLink } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { Fragment, useEffect, useState } from 'react';

import { ManualComponent } from './manualComponent';

import { useDonatelloIntegration } from '@/services/api/integrations/donatello';

export default () => {
	const manager = useDonatelloIntegration();
	const { data: apiKey } = manager.useData();

	const [key, setKey] = useState<string>();

	const siteUrl = `${window.location.origin}/api/webhooks/integrations/donatello`;

	useEffect(() => {
		if (apiKey) {
			setKey(apiKey);
		}
	}, [apiKey]);


	return (
		<ManualComponent
			integrationKey={'donatello'}
			imageSize={50}
			body={
				<Stepper
					active={0}
					allowNextStepsSelect={false}
					orientation="vertical"
					w={'100%'}
					styles={{
						step: {
							width: '100%',
						},
						stepBody: {
							width: '100%',
						},
					}}
				>
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
							<div>
								Copy api key and paste into "Api Key" input
								<TextInput
									mt={5}
									value={key}
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
							</div>
						}
					/>
					<Stepper.Step
						label="Step 3"
						description={
							<div>
								Copy link and paste into link field
								<TextInput
									styles={{
										input: {
											width: '100%',
										},
									}}
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
							</div>
						}
					/>
				</Stepper>
			}
		/>
	);
};
