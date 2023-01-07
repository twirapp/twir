import { ActionIcon, Card, Divider, Flex, Grid, Group, Text } from '@mantine/core';
import { IconPencil } from '@tabler/icons';
import type { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { useTranslation } from 'next-i18next';

import { typesMapping } from './mapping';

type Props = React.PropsWithChildren<{
  settings: ChannelModerationSetting;
  setEditableSettings: React.Dispatch<React.SetStateAction<ChannelModerationSetting>>;
  setEditDrawerOpened: React.Dispatch<React.SetStateAction<boolean>>;
}>;

export const ModerationCard: React.FC<Props> = (props) => {
  const { t } = useTranslation('moderation');

  return (
    <Grid grow>
      <Grid.Col span={4}>
        <Card>
          <Card.Section p="">
            <Flex gap="xs" direction="row" justify="space-between">
              <Group position="left">
                {typesMapping[props.settings.type]['icon']({
                  color: typesMapping[props.settings.type].iconColor,
                })}
                <Text size="lg">
                  {typesMapping[props.settings.type].name ??
                    props.settings.type.charAt(0).toUpperCase() + props.settings.type.slice(1)}
                </Text>
              </Group>
              <ActionIcon
                onClick={() => {
                  props.setEditDrawerOpened(true);
                  props.setEditableSettings(props.settings);
                }}
              >
                <IconPencil size={18}></IconPencil>
              </ActionIcon>
            </Flex>
          </Card.Section>
          <Divider />
          <Card.Section p="lg">{t(`descriptions.${props.settings.type}`)}</Card.Section>
        </Card>
      </Grid.Col>
    </Grid>
  );
};
