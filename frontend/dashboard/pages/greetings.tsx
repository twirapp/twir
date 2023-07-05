import { ActionIcon, Avatar, Badge, Button, Flex, Group, Switch, Table, Text, TextInput } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { IconPencil, IconSearch, IconTrash } from '@tabler/icons';
import type { Greeting } from '@twir/grpc/generated/api/api/greetings';
import { useTranslation } from 'next-i18next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useState } from 'react';

import { confirmDelete } from '@/components/confirmDelete';
import { GreetingModal } from '@/components/greetings/modal';
import { useGreetingsManager, useTwitchUsers } from '@/services/api';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
    props: {
        ...(await serverSideTranslations(locale, ['greetings', 'layout'])),
    },
});

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableGreeting, setEditableGreeting] = useState<Greeting | undefined>();
  const { t } = useTranslation('greetings');
  const viewPort = useViewportSize();

  const greetingsManager = useGreetingsManager();

  const { data: greetings } = greetingsManager.getAll({});
  const patcher = greetingsManager.patch!;
  const deleter = greetingsManager.deleteOne!;

	const twitchUsers = useTwitchUsers([], greetings?.greetings?.map(g => g.userId) ?? []);

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Group>
          <Text size="lg">{t('title')}</Text>
          <TextInput
            placeholder={'search...'}
            rightSection={<IconSearch size={18} />}
          />
        </Group>

        <Button
          color="green"
          onClick={() => {
            setEditableGreeting(undefined);
            setEditDrawerOpened(true);
          }}
        >
            {t('create')}
        </Button>
      </Flex>
      <Table>
        <thead>
          <tr>
            <th style={{ width:50 }}></th>
            <th>{t('userName')}</th>
            {viewPort.width > 550 && <th>{t('message')}</th>}
            <th>{t('table.head.status')}</th>
            <th>{t('table.head.actions')}</th>
          </tr>
        </thead>
        <tbody>
          {greetings?.greetings
              .map((greeting, idx) => {
								const twitchUser = twitchUsers.data?.users.find(u => u.id === greeting.userId);

								return <tr key={greeting.id}>
									<td>
										<Avatar src={twitchUser?.profileImageUrl} style={{ borderRadius:111 }} />
									</td>
									<td>
										<Badge>{twitchUser?.login}</Badge>
									</td>
									{viewPort.width > 550 && <td>
										<Badge color="cyan">{greeting.text}</Badge>
									</td>}
									<td>
										<Switch
											checked={greeting.enabled}
											onChange={(event) => {
												patcher.mutate({ id: greeting.id, enabled: event.currentTarget.checked });
											}}
										/>
									</td>
									<td>
										<Flex direction="row" gap="xs">
											<ActionIcon
												onClick={() => {
													setEditableGreeting(greetings!.greetings[idx] as any);
													setEditDrawerOpened(true);
												}}
												variant="filled"
												color="blue"
											>
												<IconPencil size={14} />
											</ActionIcon>

											<ActionIcon
												onClick={() =>
													confirmDelete({
														onConfirm: () => deleter.mutate({ id: greeting.id }),
													})
												}
												variant="filled"
												color="red"
											>
												<IconTrash size={14} />
											</ActionIcon>
										</Flex>
									</td>
								</tr>;
							})}
        </tbody>
      </Table>

      <GreetingModal
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        greeting={editableGreeting}
      />
    </div>
  );
}
