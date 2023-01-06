import { createStyles, Group, Switch, Text } from '@mantine/core';

import { DashboardCard } from '@/components/common';

const useStyles = createStyles(() => ({
  switch: {
    '& *': {
      cursor: 'pointer',
    },
  },
}));

export function ServiceManage() {
  const { classes } = useStyles();

  return (
    <DashboardCard title="Services">
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
    </DashboardCard>
  );
}
