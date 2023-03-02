import { Badge, Flex, Table, Text, Tooltip } from '@mantine/core';
import { IconCategory, IconCurrencyDollar, IconDiamond, IconSword, IconVideo } from '@tabler/icons-react';
import { useQuery } from '@tanstack/react-query';
import { GetServerSideProps, NextPage } from 'next';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import { useUsersByNames } from '@/services/users';

export type Command = {
  id: string
  name: string
  responses: string[]
  permissions: string[]
  cooldown: number
  cooldownType: string
  aliases: string[]
  description: null | string
  group: string | null
  module: string
}

const rolesMapping: Record<string, JSX.Element> = {
  'BROADCASTER': <IconVideo color={'#db4f4f'} />,
  'MODERATOR': <IconSword color={'green'} />,
  'VIP': <IconDiamond color={'pink'} />,
  'SUBSCRIBER': <IconCurrencyDollar color={'cyan'} />,
};


const Commands: NextPage = () => {
  const router = useRouter();
  const { data: users } = useUsersByNames([router.query.channelName as string]);
  const [commands, setCommands] = useState<Record<string, Command[]>>();

  const {
    data: rawCommands,
  } = useQuery({
    queryKey: ['commands', users?.at(0)?.id],
    queryFn: async (): Promise<Command[]> => {
      const req = await fetch(`/api/v1/p/commands/${users?.at(0)?.id}`);

      return req.json();
    },
    enabled: !!users?.at(0)?.id,
  });

  useEffect(() => {
    if (!rawCommands) return;
    const result = {} as Record<string, Command[]>;

    const customCommands = rawCommands.filter(c => c.module === 'CUSTOM');
    const otherCommands = rawCommands.filter(c => c.module !== 'CUSTOM');

    console.log(customCommands, otherCommands);
    for (const command of customCommands) {
      const key = command.group ?? command.module;
      if (!result[key]) result[key] = [];
      result[key].push(command);
    }

    for (const command of otherCommands) {
      const key = command.module;
      if (!result[key]) result[key] = [];
      result[key].push(command);
    }

    setCommands(result);
  }, [rawCommands]);

  return (
    <>
    <Table highlightOnHover>
      <thead>
      <tr>
        <th>Name</th>
        <th>Description</th>
        <th>Permissions</th>
        <th>Cooldown</th>
      </tr>
      </thead>
      <tbody>
      {commands && Object.keys(commands).map((module) => {
        const cmds = commands[module];

        return (<>
          {module !== 'CUSTOM' && <tr style={{ height: 50 }}>
              <td colSpan={4}>
                  <Text size={'md'}>
                    <IconCategory size={17} style={{ paddingTop: 4 }} />
                    {' '} {module.charAt(0).toUpperCase() + module.slice(1).toLowerCase()}
                  </Text>
              </td>
          </tr>}
          {cmds.map((c, commandIndex) => (<tr key={commandIndex}>
              <td style={{
                whiteSpace: 'nowrap',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                maxWidth: 150,
              }}>
                <Tooltip label={[c?.name, ...c.aliases || []].join(', ')}>
                  <Text truncate>
                    {[c?.name, ...c.aliases || []].join(', ')}
                  </Text>
                </Tooltip>
              </td>
              <td>{c.description ? c.description : c?.responses?.map((r, responseIndex) => <Text
                key={responseIndex}
                title={r}
                lineClamp={1}
                style={{ textOverflow: 'ellipsis', overflow: 'hidden' }}
              >
                {r}
              </Text>)}</td>
              <td>
                <Flex direction={'column'} gap={'xs'} align={'center'}>
                  <Flex direction={'row'} gap={'xs'}>
                    {c?.permissions?.map((p, i) => {
                      if (rolesMapping[p]) {
                        return <Tooltip label={p} key={i}>{rolesMapping[p]}</Tooltip>;
                      }
                    })}
                  </Flex>
                  <Flex direction={'row'} gap={'xs'}>
                    {c?.permissions?.map((p, i) => {
                      if (!rolesMapping[p]) {
                        return <Badge color={'green'} size={'sm'}>{p}</Badge>;
                      }
                    })}
                  </Flex>
                </Flex>
              </td>
              <td>{c?.cooldown} ({c?.cooldownType?.toLowerCase().replace('_', ' ')})</td>
            </tr>))}
        </>);
      })}
      </tbody>
  </Table>
  </>);
};

export default Commands;