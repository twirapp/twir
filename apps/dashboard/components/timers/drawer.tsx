import {
  Badge,
  Drawer,
  Flex,
  Group,
  ScrollArea,
  Slider,
  Text,
  TextInput,
  useMantineTheme,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import { useEffect } from 'react';

type Props = {
  opened: boolean;
  timer: ChannelTimer;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

const sliderSteps = new Array(6).fill(1).map((_, index) => {
  const value = (index + 1) * 15;
  return { value, label: value.toString() };
});

export const TimerDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<ChannelTimer>({});
  const viewPort = useViewportSize();

  useEffect(() => {
    form.setValues(props.timer);
  }, [props.timer]);

  return (
    <Drawer
      opened={props.opened}
      onClose={() => props.setOpened(false)}
      title={<Badge size="xl">{props.timer.name}</Badge>}
      padding="xl"
      size="xl"
      position="right"
      transition="slide-left"
      overlayColor={theme.colorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[2]}
      overlayOpacity={0.55}
      overlayBlur={3}
    >
      {/* <ScrollArea.Autosize maxHeight={viewPort.height - 120} type="auto" offsetScrollbars={true}> */}
      <form onSubmit={form.onSubmit((values) => console.log(values))}>
        <Flex direction="column" gap="md" justify="flex-start" align="flex-start" wrap="wrap">
          <TextInput {...form.getInputProps('name')} label="Name" required></TextInput>
          <div style={{ width: '100%' }}>
            <Text>Time interval (seconds)</Text>
            <Flex direction="row" wrap="wrap" gap="md">
              <Slider
                w={'70%'}
                style={{ marginTop: 9 }}
                size="sm"
                marks={sliderSteps}
                value={form.getInputProps('timeInterval').value}
                onChange={(v) => form.setFieldValue('timeInterval', v)}
              />
              <TextInput w={'20%'} {...form.getInputProps('timeInterval')}></TextInput>
            </Flex>
          </div>
        </Flex>
      </form>
      {/* </ScrollArea.Autosize> */}
    </Drawer>
  );
};
