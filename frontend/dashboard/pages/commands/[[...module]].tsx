import {
  ActionIcon,
  Badge,
  Button,
  Flex,
  Group,
  Switch,
  Table,
  Text,
  TextInput,
} from '@mantine/core';
import { useDebouncedState, useViewportSize } from '@mantine/hooks';
import {
  IconCaretDown,
  IconCaretUp,
  IconFolder,
  IconPencil,
  IconSearch,
  IconTrash,
} from '@tabler/icons';
import { ChannelCommand } from '@twir/typeorm/entities/ChannelCommand';
import { ChannelCommandGroup } from '@twir/typeorm/entities/ChannelCommandGroup';
import { GetServerSideProps, NextPage } from 'next';
import { useTranslation } from 'next-i18next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useRouter } from 'next/router';
import { Fragment, useEffect, useState } from 'react';

import { CommandsModal } from '@/components/commands/modal';
import { ChannelCommandGroupDrawer } from '@/components/commandsGroup/drawer';
import { confirmDelete } from '@/components/confirmDelete';
import { commandsGroupManager, commandsManager } from '@/services/api';

export const getServerSideProps: GetServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale!, ['commands', 'layout'])),
  },
});

const Commands: NextPage = () => {
  const router = useRouter();
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [groupEditDrawerOpened, setGroupEditDrawerOpened] = useState(false);

  const [editableCommand, setEditableCommand] = useState<ChannelCommand | undefined>();
  const [editableGroup, setEditableGroup] = useState<ChannelCommandGroup | undefined>();

  const viewPort = useViewportSize();
  const { t } = useTranslation('commands');

  const { useGetAll, usePatch, useDelete } = commandsManager();
  const patcher = usePatch();
  const deleter = useDelete();
  const { data: commands } = useGetAll();

  const groupManager = commandsGroupManager();
  const groupDeleter = groupManager.useDelete();

  const [commandsWithGroups, setCommandsWithGroups] = useState<{
    [x: string]: {
      list: ChannelCommand[];
      show: boolean;
      id: string;
      color?: string;
    };
  }>({
    default: {
      list: [],
      show: true,
      id: '',
    },
  });

  function setCommandsGroup() {
    if (!commands) return;

    setCommandsWithGroups({});

    for (const command of commands) {
      if (
        command.module.toLowerCase() !== (router.query.module ? router.query.module[0] : 'custom')
      )
        continue;

      if (!command.group) {
        setCommandsWithGroups((prev) => ({
          ...prev,
          default: {
            ...prev.default,
            list: [...(prev.default?.list || []), command],
          },
        }));
        continue;
      } else {
        setCommandsWithGroups((prev) => ({
          ...prev,
          [command.group!.name]: {
            list: [...(prev[command.group!.name]?.list || []), command],
            show: prev[command.group!.name]?.show || false,
            id: command.group!.id,
            color: command.group!.color,
          },
        }));
      }
    }
  }

  useEffect(() => {
    setCommandsGroup();
  }, [commands]);

  useEffect(() => {
    setCommandsGroup();
  }, [router.query.module]);

  const [searchInput, setSearchInput] = useDebouncedState('', 200);

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Group>
          <Text size="lg">{t('title')}</Text>
          <TextInput
            placeholder={'search...'}
            rightSection={<IconSearch size={18} />}
            onChange={(event) => setSearchInput(event.target.value)}
          />
        </Group>
        <Group>
          <Button
            color="green"
            onClick={() => {
              setEditableGroup(undefined);
              setGroupEditDrawerOpened(true);
            }}
          >
            {t('createGroup')}
          </Button>
          <Button
            color="green"
            onClick={() => {
              setEditableCommand(undefined);
              setEditDrawerOpened(true);
            }}
          >
            {t('createCommand')}
          </Button>
        </Group>
      </Flex>

      <Table style={{ tableLayout: 'fixed', width: '100%' }}>
        <thead>
          <tr>
            <th style={{ width: '15%' }}>{t('name')}</th>
            {viewPort.width > 550 && <th style={{ width: '70%' }}>{t('responses')}</th>}
            <th style={{ width: '5%' }}>{t('table.head.status')}</th>
            <th style={{ width: '10%' }}>{t('table.head.actions')}</th>
          </tr>
        </thead>

        <tbody>
          {Object.keys(commandsWithGroups).map((group, groupIndex) => (
            <Fragment key={groupIndex}>
              <tr
                style={{
                  padding: 5,
                  display: group === 'default' ? 'none' : undefined,
                  backgroundColor:
                    group !== 'default' && commandsWithGroups[group].show
                      ? commandsWithGroups[group].color!
                      : undefined,
                }}
              >
                <td
                  style={{
                    cursor: 'pointer',
                  }}
                  onClick={() => {
                    setCommandsWithGroups((prev) => ({
                      ...prev,
                      [group]: {
                        ...prev[group],
                        show: !prev[group].show,
                      },
                    }));
                  }}
                >
                  <Text size={'md'}>
                    <IconFolder size={15} /> {group}
                  </Text>
                </td>
                <td></td>
                <td>
                  {!commandsWithGroups[group].show && <IconCaretDown size={25} />}
                  {commandsWithGroups[group].show && <IconCaretUp size={25} />}
                </td>
                <td>
                  <Flex direction={'row'} gap="xs">
                    <ActionIcon
                      onClick={() => {
                        setEditableGroup(
                          commands?.find((c) => c.groupId === commandsWithGroups[group]!.id)?.group,
                        );
                        setGroupEditDrawerOpened(true);
                      }}
                      variant="filled"
                      color="blue"
                    >
                      <IconPencil size={14} />
                    </ActionIcon>
                    <ActionIcon
                      onClick={() =>
                        confirmDelete({
                          onConfirm: () => groupDeleter.mutate(commandsWithGroups[group].id),
                          text: 'This will delete only group, not commands. Delete group?',
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
              {commandsWithGroups[group].list
                .filter((c) => {
                  const someAliase = c.aliases.some((a) => a.includes(searchInput));
                  return someAliase || c.name.includes(searchInput);
                })
                .map((command, commandIndex) => (
                  <tr
                    key={command.id}
                    style={{
                      borderTop:
                        groupIndex === Object.keys(commandsWithGroups).length - 1 &&
                        commandIndex === 0
                          ? '1px solid #373A40'
                          : undefined,
                      display:
                        group !== 'default'
                          ? commandsWithGroups[group].show
                            ? undefined
                            : 'none'
                          : undefined,
                      marginLeft: group !== 'default' && commandsWithGroups[group].show ? 10 : 10,
                      backgroundColor:
                        group !== 'default' && commandsWithGroups[group].show
                          ? commandsWithGroups[group].color!
                          : undefined,
                    }}
                  >
                    <td
                      style={{
                        whiteSpace: 'nowrap',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        maxWidth: 100,
                        paddingLeft: 10,
                      }}
                    >
                      <Badge>
                        <Text truncate>{command.name}</Text>
                      </Badge>
                    </td>
                    {viewPort.width > 550 && (
                      <td>
                        {command.module != 'CUSTOM' && (
                          <Text dangerouslySetInnerHTML={{ __html: command.description || '' }} />
                        )}
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
                      />
                    </td>
                    <td>
                      <Flex direction="row" gap="xs">
                        <ActionIcon
                          onClick={() => {
                            setEditableCommand(commands!.find((c) => c.id === command.id)!);
                            setEditDrawerOpened(true);
                          }}
                          variant="filled"
                          color="blue"
                        >
                          <IconPencil size={14} />
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
                            <IconTrash size={14} />
                          </ActionIcon>
                        )}
                      </Flex>
                    </td>
                  </tr>
                ))}
            </Fragment>
          ))}
        </tbody>
      </Table>

      <CommandsModal
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        command={editableCommand}
      />

      <ChannelCommandGroupDrawer
        opened={groupEditDrawerOpened}
        setOpened={setGroupEditDrawerOpened}
        group={editableGroup}
      />
    </div>
  );
};

export default Commands;
