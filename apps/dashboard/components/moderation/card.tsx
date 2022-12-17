import { ActionIcon, Card, Divider, Flex, Grid, Group, Text } from '@mantine/core';
import {
  IconLambda,
  IconLetterCaseUpper,
  IconLink,
  IconMoodSmile,
  IconPencil,
  IconPlaylistX,
  IconTextWrapDisabled,
  TablerIcon,
} from '@tabler/icons';
import type {
  ChannelModerationSetting,
  SettingsType,
} from '@tsuwari/typeorm/entities/ChannelModerationSetting';

type Props = React.PropsWithChildren<{
  settings: ChannelModerationSetting;
  setEditableSettings: React.Dispatch<React.SetStateAction<ChannelModerationSetting>>;
  setEditDrawerOpened: React.Dispatch<React.SetStateAction<boolean>>;
}>;

const typesMapping: Record<
  keyof typeof SettingsType,
  {
    icon: TablerIcon;
    iconColor?: string;
    name?: string;
    description: string;
  }
> = {
  links: {
    icon: IconLink,
    iconColor: 'cyan',
    description: `Remove messages containing any links you haven't whitelisted.`,
  },
  caps: {
    icon: IconLetterCaseUpper,
    iconColor: 'orange',
    description: `Remove messages containing excessive amounts of capital letters.`,
  },
  emotes: {
    icon: IconMoodSmile,
    iconColor: 'yellow',
    description: 'Remove messages containing an excessive amount of emotes.',
  },
  longMessage: {
    icon: IconTextWrapDisabled,
    name: 'Long Messages',
    description: `Remove lengthy messages.`,
  },
  blacklists: {
    icon: IconPlaylistX,
    name: 'Deny List',
    description: 'Remove denied words from chat.',
  },
  symbols: {
    icon: IconLambda,
    iconColor: 'green',
    description: `Remove messages containing disruptive or excessive use of symbols.`,
  },
};

export const ModerationCard: React.FC<Props> = (props) => {
  return (
    <Grid grow>
      <Grid.Col span={4}>
        <Card>
          <Card.Section p="">
            <Flex gap="xs" direction="row" justify="space-between">
              <Group position="left">
                {/* <props.icon color={props.iconColor} /> */}
                {typesMapping[props.settings.type]['icon']({
                  color: typesMapping[props.settings.type].iconColor,
                })}
                <Text size="lg">
                  {typesMapping[props.settings.type].name ??
                    props.settings.type.charAt(0).toUpperCase() + props.settings.type.slice(1)}
                </Text>
              </Group>
              <ActionIcon onClick={() => props.setEditDrawerOpened(true)}>
                <IconPencil size={18}></IconPencil>
              </ActionIcon>
            </Flex>
          </Card.Section>
          <Divider />
          <Card.Section p="lg">{typesMapping[props.settings.type].description}</Card.Section>
        </Card>
      </Grid.Col>
    </Grid>
  );
};
