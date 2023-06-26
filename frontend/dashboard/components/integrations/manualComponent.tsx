import { Button, Flex, Modal, useMantineTheme } from '@mantine/core';
import { showNotification } from '@mantine/notifications';
import { IconLogin, IconManualGearbox, IconSettings } from '@tabler/icons';
import Image from 'next/image';
import { Fragment, useEffect, useState } from 'react';

type Props = {
	integrationKey: string
	body: React.ReactElement
	save?: () => void | Promise<void>
	imageSize?: number
	imageScale?: number
}

export const ManualComponent: React.FC<Props> = (props) => {
	const [opened, setOpened] = useState(false);
	const theme = useMantineTheme();

	const [uppercasedKey, setUppercasedKey] = useState<string>('');

	useEffect(() => {
		setUppercasedKey(props.integrationKey.charAt(0).toUpperCase() + props.integrationKey.slice(1));
	}, [props.integrationKey]);

	return <Fragment>
		<tr key={props.integrationKey} style={{ height: 65 }}>
			<td><Image
				src={`/dashboard/assets/icons/brands/${props.integrationKey}.svg`}
				height={props.imageSize ?? 30}
				width={props.imageSize ?? 30}
				alt={props.integrationKey}
				style={{
					scale: props.imageScale ?? 1,
				}}
			/>
			</td>

			<td></td>

			<td>
				<Flex direction='row' gap='sm' justify='flex-end'>
					<Button
						compact
						leftIcon={<IconSettings />}
						variant='light'
						color='blue'
						onClick={() => setOpened(true)}
					>
						Connect
					</Button>
				</Flex>
			</td>
		</tr>
		<Modal
			styles={{
				modal: {
					backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[3],
				},
			}}
			overlayColor={theme.colorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[2]}
			overlayOpacity={0.55}
			overlayBlur={3}
			opened={opened}
			onClose={() => setOpened(false)}
			title={uppercasedKey}
			closeOnClickOutside={false}
		>
			{props.body}
			{props.save && <Button
				mt={5}
				variant={'light'}
				color={'green'}
				w={'100%'}
				onClick={() => {
					props.save!();
					setOpened(false);
					showNotification({
						title: 'Saved',
						message: `${uppercasedKey} settings saved.`,
						color: 'green',
					});
				}}
			>Save</Button>}
		</Modal>
	</Fragment>;
};
