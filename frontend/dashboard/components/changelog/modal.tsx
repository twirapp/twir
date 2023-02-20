import { Button, Group, Text, Modal, Flex } from '@mantine/core';
import { useAtom } from 'jotai';
import { useEffect, useState } from 'react';

import { groupBy } from '../../util/groupBy';

import { modalOpenedAtomic } from '@/components/changelog/store';
import { useChangelog, type Commit } from '@/services/api';

export const ChangelogModal = () => {
  const [opened, setOpened] = useAtom(modalOpenedAtomic);
  const { data } = useChangelog();
  const [mappedData, setMappedData] = useState<{ [x: string]: Commit[]}>({});

  useEffect(() => {
    if (!data) return;

    const groupedData = groupBy(
      data.map(v => ({
        ...v,
        commit: {
          ...v.commit,
          author: {
            ...v.commit.author,
            date: new Date(v.commit.author.date).toDateString(),
          },
        },
      })),
    (v) => v.commit.author.date);

    setMappedData(groupedData);
  }, [data]);

  return (
      <Modal
        opened={opened}
        onClose={() => setOpened(false)}
        title="Changelog"
        size={'lg'}
      >
        {Object.entries(mappedData).map(([date, commits], i) => (
          <Flex key={date} direction={'column'}>
            <Text size={'xl'} mt={i === 0 ? 0 : 20}>{date}</Text>
            <Flex direction={'column'}>
              {commits.map((c) => <Text key={c.sha} size={'sm'}>{c.commit.message.split('\n')[0]}</Text>)}
            </Flex>
          </Flex>
        ))}
      </Modal>
  );
};