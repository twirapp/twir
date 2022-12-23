import { Text, Card, Grid, Divider, Flex, Group } from '@mantine/core';
import { TablerIcon } from '@tabler/icons';
import React from 'react';

type Props = React.PropsWithChildren<{
  title: string;
  icon?: TablerIcon;
  iconColor?: string;
  header?: React.ReactNode;
}>;

export const IntegrationCard: React.FC<Props> = (props) => {
  return (
    <Grid grow>
      <Grid.Col span={4} >
        <Card>
          <Card.Section p="">
            <Flex gap="xs" direction="row" justify="space-between">
              <Group position="left">
                {props.icon && <props.icon color={props.iconColor} />}
                <Text size="lg">{props.title}</Text>
              </Group>
              {props.header && props.header}
            </Flex>
          </Card.Section>
          <Divider />
          <Card.Section p="lg">{props.children}</Card.Section>
        </Card>
      </Grid.Col>
    </Grid>
  );
};
