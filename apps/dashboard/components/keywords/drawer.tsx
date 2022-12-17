import {
  Alert,
  Badge,
  Drawer,
  Flex,
  NumberInput,
  ScrollArea,
  Switch,
  Textarea,
  TextInput,
  useMantineTheme,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import { useEffect } from 'react';

type Props = {
  opened: boolean;
  keyword: ChannelKeyword;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const KeywordDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<ChannelKeyword>({});
  const viewPort = useViewportSize();

  useEffect(() => {
    form.setValues(props.keyword);
  }, [props.keyword]);

  return (
    <Drawer
      opened={props.opened}
      onClose={() => props.setOpened(false)}
      title={<Badge size="xl">{props.keyword.text}</Badge>}
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
            <TextInput {...form.getInputProps('text')} label="Trigger" required w="100%" />
            <Switch label="Regular expression" {...form.getInputProps('isRegular')} />
            {form.values.isRegular && (
              <Alert>
                We use <b>Golang</b> as backend. So your expressions also should be for golang.
              </Alert>
            )}
            <Textarea
              {...form.getInputProps('response')}
              label="Response"
              autosize={true}
              w="100%"
            />
            <NumberInput label="Cooldown" {...form.getInputProps('cooldown')} />
            <NumberInput label="Used times" {...form.getInputProps('usages')} />

            <Switch label="Use twitch reply feature" {...form.getInputProps('isReply')} />
            {props.keyword.id && (
              <Alert>
                You can access that counter in your commands via{' '}
                <b>
                  $(keywords.counter|
                  {props.keyword.id})
                </b>{' '}
                variable.
              </Alert>
            )}
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
