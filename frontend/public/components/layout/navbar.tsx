import {
  ActionIcon,
  Burger,
  Container,
  createStyles,
  Flex,
  Group,
  Header,
  Text,
  Box,
} from '@mantine/core';
import { IconMoonStars, IconSun } from '@tabler/icons';
import Image from 'next/image';
import { Dispatch, SetStateAction } from 'react';

import { useTheme } from '../../hooks/useTheme';

const useStyles = createStyles((theme) => ({
  header: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    height: '100%',
  },

  hiddenMobile: {
    pointerEvents: 'none',
    userSelect: 'none',

    [theme.fn.smallerThan('sm')]: {
      display: 'none',
    },
  },

  hiddenDesktop: {
    [theme.fn.largerThan('sm')]: {
      display: 'none',
    },
  },
}));

export function NavBar({
  opened,
  setOpened,
}: {
  setOpened: Dispatch<SetStateAction<boolean>>;
  opened: boolean;
}) {
  const { classes } = useStyles();
  const { theme, colorScheme, toggleColorScheme } = useTheme();

  return (
    <Header height={60}>
      <Container maw="unset" className={classes.header}>
        <Flex gap="sm" justify="flex-start" align="center" direction="row">
          <Burger
            className={classes.hiddenDesktop}
            opened={opened}
            onClick={() => setOpened(!opened)}
            size="sm"
            color={theme.colors.gray[6]}
            mr="xl"
          />
          <Box display="flex" className={classes.hiddenMobile}>
            <Image src="/dashboard/TsuwariInCircle.svg" width={30} height={30} alt="Tsuwari Logo" />
            <Text
              component="span"
              ml="sm"
              sx={{
                color: 'white',
                fontFamily: 'Golos Text, sans-serif',
              }}
              fz="xl"
              fw={500}
            >
              Tsuwari
            </Text>
          </Box>
        </Flex>
        <Group position="center">
          <ActionIcon
            size="lg"
            variant="default"
            color={colorScheme === 'dark' ? 'yellow' : 'blue'}
            onClick={() => toggleColorScheme()}
            title="Toggle color scheme"
          >
            {colorScheme === 'dark' ? <IconSun size={18} /> : <IconMoonStars size={18} />}
          </ActionIcon>
        </Group>
      </Container>
    </Header>
  );
}
