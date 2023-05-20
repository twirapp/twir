import { ActionIcon, Badge, Button, CopyButton, Flex, Table, Text, Tooltip } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { IconCopy, IconPencil, IconTrash } from '@tabler/icons';
import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { useTranslation } from 'next-i18next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useState } from 'react';

import { confirmDelete } from '@/components/confirmDelete';
import { VariableModal } from '@/components/variables/modal';
import { variablesManager } from '@/services/api';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['variables', 'layout'])),
  },
});

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableVariable, setEditableVariable] = useState<ChannelCustomvar | undefined>();
  const { t } = useTranslation('variables');
  const viewPort = useViewportSize();

  const { useGetAll, useDelete } = variablesManager();
  const { data: variables } = useGetAll();
  const deleter = useDelete();

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Text size="lg">{t('title')}</Text>
        <Button
          color="green"
          onClick={() => {
            setEditableVariable(undefined);
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
            <th>{t('type')}</th>
            {viewPort.width > 550 && <th>{t('response')}</th>}
            <th>{t('table.head.actions')}</th>
          </tr>
        </thead>
        <tbody>
          {variables &&
            variables.map((variable, idx) => (
              <tr key={variable.id}>
                <td>
                  <Badge>{variable.name}</Badge>
                </td>
                <td>
                  <Badge color="cyan">{variable.type}</Badge>
                </td>
                {viewPort.width > 550 && (
                  <td>
                    {variable.type !== 'SCRIPT' && <Badge>{variable.response}</Badge>}
                    {variable.type === 'SCRIPT' && (
                      <Badge color="red">{t('table.scriptAlert')}</Badge>
                    )}
                  </td>
                )}
                <td>
                  <Flex direction="row" gap="xs">
                    <CopyButton value={`$(customvar|${variable.name})`}>
                      {({ copied, copy }) => (
                        <Tooltip label={t('table.copy')} withArrow position="bottom">
                          <ActionIcon
                            color={copied ? 'teal' : 'blue'}
                            variant="filled"
                            onClick={copy}
                          >
                            <IconCopy size={14} />
                          </ActionIcon>
                        </Tooltip>
                      )}
                    </CopyButton>
                    <ActionIcon
                      onClick={() => {
                        setEditableVariable(variables[idx] as any);
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
                          onConfirm: () => deleter.mutate(variable.id),
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

      <VariableModal
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        variable={editableVariable}
      />
    </div>
  );
}
