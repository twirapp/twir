import { ActionIcon, Badge, Button, Flex, Switch, Table, Text } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { IconPencil, IconTrash } from '@tabler/icons';
import { ChannelTimer } from '@twir/typeorm/entities/ChannelTimer';
import { useTranslation } from 'next-i18next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useState } from 'react';

import { confirmDelete } from '@/components/confirmDelete';
import { TimerModal } from '@/components/timers/modal';
import { useTimersManager } from '@/services/api';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['timers', 'layout'])),
  },
});

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableTimer, setEditableTimer] = useState<ChannelTimer | undefined>();
  const viewPort = useViewportSize();
  const { t } = useTranslation('timers');

  const { useGetAll, usePatch, useDelete } = useTimersManager();
  const { data: timers } = useGetAll();
  const patcher = usePatch();
  const deleter = useDelete();

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Text size="lg">{t('title')}</Text>
        <Button
          color="green"
          onClick={() => {
            setEditableTimer(undefined);
            setEditDrawerOpened(true);
          }}
        >
          {t('create')}
        </Button>
      </Flex>
      <Table>
        <thead>
          <tr>
            <th>{t('name')}</th>
            {viewPort.width > 550 && (
              <>
                <th>{t('responses')}</th>
                <th>{t('intervalTime')}</th>
                <th>{t('intervalMessages')}</th>
              </>
            )}
            {viewPort.width > 550 && <th>{t('table.head.status')}</th>}
            <th>{t('table.head.actions')}</th>
          </tr>
        </thead>
        <tbody>
          {timers &&
            timers.map((timer, idx) => (
              <tr key={timer.id}>
                <td>
                  <Badge>{timer.name}</Badge>
                </td>
                {viewPort.width > 550 && (
                  <>
                    <td>
                      {timer.responses.map((r, i) => (
                        <p key={i} style={{ margin: 0 }}>
                          {r.text}
                        </p>
                      ))}
                    </td>
                    <td>{timer.timeInterval} m.</td>
                    <td>{timer.messageInterval}</td>
                  </>
                )}

                <td>
                  <Switch
                    checked={timer.enabled}
                    onChange={(event) => {
                      patcher.mutate({
                        id: timer.id,
                        data: { enabled: event.currentTarget.checked },
                      });
                    }}
                  />
                </td>
                <td>
                  <Flex direction="row" gap="xs">
                    <ActionIcon
                      onClick={() => {
                        setEditableTimer(timers[idx] as any);
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
                          onConfirm: () => deleter.mutate(timer.id),
                        })
                      }
                      variant="filled"
                      color="red"
                    >
                      <IconTrash size={14} />
                    </ActionIcon>
                  </Flex>
                </td>
              </tr>
            ))}
        </tbody>
      </Table>

      <TimerModal opened={editDrawerOpened} setOpened={setEditDrawerOpened} timer={editableTimer} />
    </div>
  );
}
