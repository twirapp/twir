import { Grid, Alert, Flex, Tooltip, ActionIcon } from '@mantine/core';
import { IconMusicOff, IconPlaylist, IconSettings } from '@tabler/icons';
import { useRouter } from 'next/router';

export const AlertPlayerDisabled = () => {
  const router = useRouter();

  return (
    <Grid.Col>
      <Alert
        icon={<IconMusicOff size={16} />}
        color="red"
        radius="md"
        variant="outline"
        style={{ flexDirection: 'row', alignItems: 'center' }}
      >
        <Flex direction="row" align="center" justify="space-between">
          Song requests are currently disabled{' '}
          <Tooltip withinPortal position="left" label="Song requests settings">
            <ActionIcon variant="default" onClick={() => router.push('settings')}>
              <IconSettings size={14} />
            </ActionIcon>
          </Tooltip>
        </Flex>
      </Alert>
    </Grid.Col>
  );
};

export const AlertQueueEmpty = () => {
  return (
    <Grid.Col>
      <Alert
        icon={<IconPlaylist size={16} />}
        radius="md"
        variant="outline"
        style={{ flexDirection: 'row', alignItems: 'center' }}
      >
        <Flex direction="row" align="center" justify="space-between">
          Queue is empty
        </Flex>
      </Alert>
    </Grid.Col>
  );
};
