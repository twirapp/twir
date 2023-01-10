import { Badge, Card, Flex, ScrollArea, Table, Text } from '@mantine/core';
import { createStyles } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { useRouter } from 'next/router';

import { useUsersByNames } from '@/services/users';

type Command = {
  name: string;
  responses: string[];
  permission: string;
  cooldown: number;
  cooldownType: string;
};

export const useStyles = createStyles((theme) => ({
  card: {
    backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[0],
  },
  table: {
    overflow: 'scroll',
    tableLayout: 'fixed',
    width: '100%',
  },
  thead: {
    position: 'sticky',
    top: 0,
    zIndex: 1,
    backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.white,
    transition: 'box-shadow 150ms ease',

    '&::after': {
      content: '""',
      position: 'absolute',
      left: 0,
      right: 0,
      bottom: 0,
      borderBottom: `1px solid ${
        theme.colorScheme === 'dark' ? theme.colors.dark[4] : theme.colors.gray[2]
      }`,
    },
  },
}));

const Commands = () => {
  const router = useRouter();
  const { classes } = useStyles();
  const { data: users } = useUsersByNames([router.query.channelName as string]);

  const { data: commands } = useQuery({
    queryKey: ['commands', users?.at(0)?.id],
    queryFn: async (): Promise<Command[]> => {
      const req = await fetch(`/api/v1/p/commands/${users?.at(0)?.id}`);

      return req.json();
    },
    initialData: [],
    enabled: !!users?.at(0)?.id,
  });

  return (
    <Card withBorder radius="md" p="md">
      <Card.Section className={classes.card}>
        <ScrollArea.Autosize maxHeight="80vh">
          <Table highlightOnHover className={classes.table}>
            <thead className={classes.thead}>
              <tr>
                <th>Name</th>
                <th>Responses</th>
                <th>Permission</th>
                <th>Cooldown</th>
              </tr>
            </thead>
            <tbody>
              {commands.map((command) => (
                <tr key={command.name}>
                  <td>!{command.name}</td>
                  <td>
                    {command.responses?.map((response, key) => (
                      <Text
                        key={key}
                        title={response}
                        lineClamp={1}
                        style={{ textOverflow: 'ellipsis', overflow: 'hidden' }}
                      >
                        {response}
                      </Text>
                    ))}
                  </td>
                  <td>
                    <Badge>{command.permission}</Badge>
                  </td>
                  <td>
                    <Flex gap="sm">
                      <Badge color={command.cooldown === 0 ? 'green' : 'teal'}>
                        {command.cooldown}
                      </Badge>
                      <Badge>{command.cooldownType?.toLowerCase().replace('_', ' ')}</Badge>
                    </Flex>
                  </td>
                </tr>
              ))}
            </tbody>
          </Table>
        </ScrollArea.Autosize>
      </Card.Section>
    </Card>
  );
};

export default Commands;
