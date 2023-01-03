import { createStyles, Card, Group, Switch, Text } from '@mantine/core';

const useStyles = createStyles((theme) => ({
  card: {
    backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.white,
  },

  switch: {
    '& *': {
      cursor: 'pointer',
    },
  },

  title: {
    lineHeight: 1,
  },
}));

export function ServiceManage() {
  const { classes } = useStyles();

  return (
    <Card withBorder radius="md" p="xl" className={classes.card}>
      <Text size="lg" pb="lg" className={classes.title} weight={500}>
        Services
      </Text>
      <Group position="apart" noWrap spacing="xl">
        <div>
          <Text>Foo</Text>
          <Text size="xs" color="dimmed">
            Bar
          </Text>
        </div>
        <Switch onLabel="ON" offLabel="OFF" className={classes.switch} size="lg" />
      </Group>
      <Group position="apart" noWrap spacing="xl">
        <div>
          <Text>Foo</Text>
          <Text size="xs" color="dimmed">
            Bar
          </Text>
        </div>
        <Switch onLabel="ON" offLabel="OFF" className={classes.switch} size="lg" />
      </Group>
    </Card>
  );
}
