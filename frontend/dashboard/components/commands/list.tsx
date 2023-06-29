import { ActionIcon, Badge, Flex, Switch, Table, Text } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { IconPencil, IconTrash } from '@tabler/icons';
import { type ChannelCommand } from '@twir/typeorm/entities/ChannelCommand';
import { useTranslation } from 'next-i18next';
import { FC, Fragment, useState } from 'react';

import { CommandsModal } from '@/components/commands/modal';
import { confirmDelete } from '@/components/confirmDelete';
import { commandsManager } from '@/services/api';

type Props = {
  commands: ChannelCommand[]
}

export const CommandsList: FC<Props> = (props) => {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableCommand, setEditableCommand] = useState<ChannelCommand | undefined>();

  const { t } = useTranslation('commands');
  const viewPort = useViewportSize();

  const { usePatch, useDelete } = commandsManager();
  const patcher = usePatch();
  const deleter = useDelete();

  return (
    <Fragment>
      <Table style={{ tableLayout: 'fixed', width: '100%' }}>
        <thead>
        <tr>
          <th style={{ width: '15%' }}>{t('name')}</th>
          {viewPort.width > 550 && <th style={{ width: '70%' }}>{t('responses')}</th>}
          <th style={{ width: '10%' }}>{t('table.head.status')}</th>
          <th style={{ width: '10%' }}>{t('table.head.actions')}</th>
        </tr>
        </thead>
        <tbody>
        {props.commands.map((command) => <tr key={command.id}>
          <td style={{
            maxWidth: 100,
            paddingLeft: 10,
          }}
          >
            <Badge>
              <Text truncate>
                {command.name}
              </Text>
            </Badge>
          </td>
          {viewPort.width > 550 && (
            <td>
              {command.module != 'CUSTOM' && <Text dangerouslySetInnerHTML={{ __html: command.description || '' }}/>}
              {command.module === 'CUSTOM' &&
                (command.responses?.map((r) => (
                  <Text
                    title={r.text!}
                    lineClamp={1}
                    style={{ textOverflow: 'ellipsis', overflow: 'hidden' }}
                    key={r.id}
                  >
                    {r.text}
                  </Text>
                )) || <Badge>No Response</Badge>)}
            </td>
          )}
          <td>
            <Switch
              checked={command.enabled}
              onChange={() =>
                patcher.mutate({ id: command.id, data: { enabled: !command.enabled } })
              }
              color={command.enabled ? 'teal' : 'gray'}
            />
          </td>
          <td>
            <Flex direction="row" gap="xs">
              <ActionIcon
                onClick={() => {
                  setEditableCommand(props.commands!.find((c) => c.id === command.id)!);
                  setEditDrawerOpened(true);
                }}
                variant="filled"
                color={'blue'}
              >
                <IconPencil size={14}/>
              </ActionIcon>
              {command.module === 'CUSTOM' && (
                <ActionIcon
                  onClick={() =>
                    confirmDelete({
                      onConfirm: () => deleter.mutate(command.id),
                    })
                  }
                  variant="filled"
                  color="red"
                >
                  <IconTrash size={14}/>
                </ActionIcon>
              )}
            </Flex>
          </td>
        </tr>)}
        </tbody>
      </Table>

      <CommandsModal
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        command={editableCommand}
      />
    </Fragment>
  );
};
