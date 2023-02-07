import { Avatar, Center, createStyles, Group, Table, UnstyledButton, Text, Flex, Button } from '@mantine/core';
import { IconChevronDown, IconChevronUp, IconSelector } from '@tabler/icons';
import { NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useState } from 'react';

import { confirmDelete } from '@/components/confirmDelete';
import { SortyByField, useCommunity } from '@/services/api';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['community', 'layout'])),
  },
});

const HOUR = 1000 * 60 * 60;

const useStyles = createStyles((theme) => ({
  th: {
    padding: '0 !important',
  },

  control: {
    width: '100%',
    padding: `${theme.spacing.xs}px ${theme.spacing.md}px`,

    '&:hover': {
      backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[0],
    },
  },

  icon: {
    width: 21,
    height: 21,
    borderRadius: 21,
  },
}));

const sortableColumns = ['watched', 'messages', 'emotes'] as const;

const Community: NextPage = () => {
  const community = useCommunity();
  const statsResetter = community.useResetStats();

  const [sortBy, setSortBy] = useState<SortyByField>('watched');
  const [order, setOrder] = useState<'desc' | 'asc'>('desc');
  const [reverseSortDirection, setReverseSortDirection] = useState(false);
  const users = community.useUsers(50, 1, sortBy, order);

  const setSorting = (field: SortyByField) => {
    const reversed = field === sortBy ? !reverseSortDirection : false;
    setReverseSortDirection(reversed);
    setSortBy(field);

    reversed ? setOrder('asc') : setOrder('desc');
  };

  return (
    <div>
      <Flex justify={'space-between'} mb={15}>
        <div></div>
        <Group>
          {sortableColumns.map((item, i) => <Button
            onClick={() => confirmDelete({ onConfirm: () => statsResetter.mutate(item) })}
            variant={'light'}
          >
            Reset {item}
          </Button>)}
        </Group>
      </Flex>
      <Table>
        <thead>
        <tr>
          <th style={{ width:50 }}></th>
          <th>Name</th>
          {sortableColumns.map((item, i) => <Th
            sorted={sortBy === item.toLowerCase()}
            reversed={reverseSortDirection}
            onSort={() => setSorting(item)}
            key={i}
          >
            {item.charAt(0).toUpperCase() + item.slice(1)}
          </Th>)}
        </tr>
        </thead>
        <tbody>
        {users.data?.map(u => <tr key={u.id}>
          <td><Avatar src={u.avatarUrl} style={{ borderRadius:111 }} /></td>
          <td>{u.name === u.displayName.toLowerCase() ? u.displayName : `${u.displayName} (${u.name})`}</td>
          <td>{(u.watched / HOUR).toFixed(1)}h</td>
          <td>{u.messages}</td>
          <td>{u.emotes}</td>
        </tr>)}
        </tbody>
      </Table>
    </div>
  );
};

interface ThProps {
  children: React.ReactNode;
  reversed: boolean;
  sorted: boolean;
  onSort(): void;
}

function Th({ children, reversed, sorted, onSort }: ThProps) {
  const { classes } = useStyles();
  const Icon = sorted ? (reversed ? IconChevronUp : IconChevronDown) : IconSelector;
  return (
    <th className={classes.th}>
      <UnstyledButton onClick={onSort} className={classes.control}>
        <Group position="apart">
          <Text weight={500} size="sm">
            {children}
          </Text>
          <Center className={classes.icon}>
            <Icon size={14} stroke={1.5} />
          </Center>
        </Group>
      </UnstyledButton>
    </th>
  );
}

export default Community;