import {
  ActionIcon,
  Badge,
  Button,
  CopyButton,
  Flex,
  Switch,
  Table,
  Text,
  Tooltip,
} from '@mantine/core';
import { IconCopy, IconPencil, IconTrash } from '@tabler/icons';
import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import { useTranslation } from 'next-i18next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useReducer, useState } from 'react';

import { confirmDelete } from '@/components/confirmDelete';
import { KeywordDrawer } from '@/components/keywords/drawer';
import { keywordsManager } from '@/services/api';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
    props: {
        ...(await serverSideTranslations(locale, ['keywords', 'layout'])),
    },
});

export default function () {
  const [, forceUpdate] = useReducer(x => x + 1, 0);
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableKeyword, setEditableKeyword] = useState<ChannelKeyword | undefined>();
  const { t } = useTranslation('keywords');

  const { data: keywords } = keywordsManager.getAll;

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Text size="lg">{t('title')}</Text>
        <Button
          color="green"
          onClick={() => {
            setEditableKeyword(undefined);
            setEditDrawerOpened(true);
          }}
        >
            {t('create')}
        </Button>
      </Flex>
      <Table>
        <thead>
          <tr>
            <th>{t('trigger')}</th>
            <th>{t('response')}</th>
            <th>{t('usages')}</th>
            <th>{t('table.head.status')}</th>
            <th>{t('table.head.actions')}</th>
          </tr>
        </thead>
        <tbody>
          {keywords &&
            keywords.map((keyword, idx) => (
              <tr key={keyword.id}>
                <td>
                  <Badge>{keyword.text}</Badge>
                </td>
                <td>{keyword.response}</td>
                <td>
                  <Badge>{keyword.usages}</Badge>
                </td>
                <td>
                  <Switch
                    checked={keyword.enabled}
                    onChange={(event) => {
                      keywordsManager.patch.mutate({ id: keyword.id, data: { enabled: event.currentTarget.checked } });
                        // .then(() => forceUpdate());
                    }}
                  />
                </td>
                <td>
                  <Flex direction="row" gap="xs">
                    <CopyButton value={`$(keywords.counter|${keyword.id})`}>
                      {({ copied, copy }) => (
                        <Tooltip
                          label={t('copy')}
                          withArrow
                          position="bottom"
                        >
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
                        setEditableKeyword(keywords[idx] as any);
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
                          onConfirm: () => keywordsManager.delete.mutate(keyword.id),
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

      <KeywordDrawer
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        keyword={editableKeyword}
      />
    </div>
  );
}
