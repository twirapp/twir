import { Badge, Table, Text } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { NextPage } from 'next';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';

import { useUsersByNames } from '@/services/users';

type Command = {
  name: string
  responses: string[]
  permission: string
  cooldown: number
  cooldownType: string
}


const Commands: NextPage = () => {
  const router = useRouter();
  const { data: users } = useUsersByNames([router.query.channelName as string]);

  const {
    data: commands,
  } = useQuery({
    queryKey: ['commands', users?.at(0)?.id],
    queryFn: async (): Promise<Command[]> => {
      const req = await fetch(`/api/v1/p/commands/${users?.at(0)?.id}`);

      return req.json();
    },
    enabled: !!users?.at(0)?.id,
  });

  return (<Table highlightOnHover>
    <thead>
    <tr>
      <th>Name</th>
      <th>Responses</th>
      <th>Permission</th>
      <th>Cooldown</th>
    </tr>
    </thead>
    <tbody>
    {commands?.map(c => <tr>
      <td>{c.name}</td>
      <td>{c.responses.map(r => <Text
        title={r}
        lineClamp={1}
        style={{ textOverflow: 'ellipsis', overflow: 'hidden' }}
      >
        {r}
      </Text>)}</td>
      <td><Badge>{c.permission}</Badge></td>
      <td>{c.cooldown} ({c.cooldownType.toLowerCase().replace('_', ' ')})</td>
    </tr>)}
    </tbody>
  </Table>);
};

export default Commands;