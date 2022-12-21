import { ActionIcon, Badge, Button, Flex, Switch, Table, Text } from '@mantine/core';
import { IconPencil, IconTrash } from '@tabler/icons';
import { useState } from 'react';

import { GreetingDrawer } from '../components/greetings/drawer';
import { type Greeting, useGreetingsManager } from '../services/api';

import { confirmDelete } from '@/components/confirmDelete';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import {useTranslation} from "next-i18next";

// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
    props: {
        ...(await serverSideTranslations(locale, ['greetings', 'layout'])),
    },
});

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableGreeting, setEditableGreeting] = useState<Greeting | undefined>();
  const { t } = useTranslation("greetings")

  const manager = useGreetingsManager();
  const { data: greetings } = manager.getAll();

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Text size="lg">{t("title")}</Text>
        <Button
          color="green"
          onClick={() => {
            setEditableGreeting(undefined);
            setEditDrawerOpened(true);
          }}
        >
            {t("create")}
        </Button>
      </Flex>
      <Table>
        <thead>
          <tr>
            <th>{t("userName")}</th>
            <th>{t("message")}</th>
            <th>{t("table.head.status")}</th>
            <th>{t("table.head.actions")}</th>
          </tr>
        </thead>
        <tbody>
          {greetings &&
            greetings.map((greeting, idx) => (
              <tr key={greeting.id}>
                <td>
                  <Badge>{greeting.userName}</Badge>
                </td>
                <td>
                  <Badge color="cyan">{greeting.text}</Badge>
                </td>
                <td>
                  <Switch
                    checked={greeting.enabled}
                    onChange={(event) => {
                      manager.patch(greeting.id, { enabled: event.currentTarget.checked });
                    }}
                  />
                </td>
                <td>
                <Flex direction="row" gap="xs">
                    <ActionIcon
                      onClick={() => {
                        setEditableGreeting(greetings[idx] as any);
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
                          onConfirm: () => manager.delete(greeting.id),
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

      <GreetingDrawer
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        greeting={editableGreeting}
      />
    </div>
  );
}
