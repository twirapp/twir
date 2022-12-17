import { Badge, Drawer, Flex, ScrollArea, useMantineTheme } from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { useEffect } from 'react';

type Props = {
  opened: boolean;
  settings: ChannelModerationSetting;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const ModerationDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<ChannelModerationSetting>({});
  const viewPort = useViewportSize();

  useEffect(() => {
    form.setValues(props.settings);
  }, [props.settings]);

  return (
    <Drawer
      opened={props.opened}
      onClose={() => props.setOpened(false)}
      title={
        <Badge size="xl">
          {props.settings.type?.charAt(0).toUpperCase() + props.settings.type?.slice(1)}
        </Badge>
      }
      padding="xl"
      size="xl"
      position="right"
      transition="slide-left"
      overlayColor={theme.colorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[2]}
      overlayOpacity={0.55}
      overlayBlur={3}
    >
      <ScrollArea.Autosize maxHeight={viewPort.height - 120} type="auto" offsetScrollbars={true}>
        <form onSubmit={form.onSubmit((values) => console.log(values))}>
          <Flex direction="column" gap="md" justify="flex-start" align="flex-start" wrap="wrap">
            qwe
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
