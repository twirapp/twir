import { Avatar, Box, Group, Text, UnstyledButton, useMantineTheme } from '@mantine/core';
import { TablerIcon } from '@tabler/icons';
import React from 'react';

type Props = {
  image?: string,
  text: string,
  onClick?: () => void
  icon?: TablerIcon
}

export const YouTubeSettingsListButtonButton: React.FC<Props> = (props) => {
  const theme = useMantineTheme();

  return <UnstyledButton
    sx={{
      display: 'block',
      width: '100%',
      padding: theme.spacing.xs,
      borderRadius: theme.radius.sm,
      color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,
      backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.colors.gray[3],
      '&:hover': {
        backgroundColor:
          theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[6],
      },
    }}
    onClick={props.onClick}
  >
    <Group>
      {props.image && <Avatar
          src={props.image}
          radius="xl"
      />}
      <Box sx={{ flex: 1 }}>
        <Text size="sm" weight={500}>
          {props.text}
        </Text>
      </Box>

      {props.icon && <props.icon size={18}/>}
    </Group>
  </UnstyledButton>;
};