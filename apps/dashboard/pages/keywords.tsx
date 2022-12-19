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
import { useState } from 'react';

import { confirmDelete } from '@/components/confirmDelete';
import { KeywordDrawer } from '@/components/keywords/drawer';
import { useKeywordsManager } from '@/services/api';

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableKeyword, setEditableKeyword] = useState<ChannelKeyword | undefined>();

  const manager = useKeywordsManager();
  const { data: keywords } = manager.getAll();

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Text size="lg">Keywords</Text>
        <Button
          color="green"
          onClick={() => {
            setEditableKeyword(undefined);
            setEditDrawerOpened(true);
          }}
        >
          Create
        </Button>
      </Flex>
      <Table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Usages</th>
            <th>Status</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {keywords &&
            keywords.map((keyword, idx) => (
              <tr key={keyword.id}>
                <td>
                  <Badge>{keyword.text}</Badge>
                </td>
                <td>
                  <Badge>{keyword.usages}</Badge>
                </td>
                <td>
                  <Switch
                    checked={keyword.enabled}
                    onChange={(event) => {
                      manager.patch(keyword.id, { enabled: event.currentTarget.checked });
                    }}
                  />
                </td>
                <td>
                  <Flex direction="row" gap="xs">
                    <CopyButton value={`$(keywords.counter|${keyword.id})`}>
                      {({ copied, copy }) => (
                        <Tooltip
                          label="Copy variable id for use in commands"
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
                          onConfirm: () => manager.delete(keyword.id),
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
