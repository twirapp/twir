import { Card, Divider, Flex, Grid, Group, Text } from '@mantine/core';
import { TablerIcon } from '@tabler/icons';

type Props = React.PropsWithChildren<{
  title: string;
  icon?: TablerIcon;
  iconColor?: string;
  header?: React.ReactNode;
}>;

export const ModerationCard: React.FC<Props> = (props) => {
  return (
    <Grid grow>
      <Grid.Col span={4}>
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
