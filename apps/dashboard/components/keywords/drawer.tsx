import {
  Alert,
  Badge,
  Button,
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

import { useKeywordsManager } from '@/services/api';

type Props = {
  opened: boolean;
  keyword?: ChannelKeyword;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const KeywordDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<ChannelKeyword>({
    validate: {
      cooldown: (v) => v < 5 && 'Cooldown cannot be lower then 5.',
      text: (v) => !v?.length && 'Text cannot be empty',
    },
    initialValues: {
      channelId: '',
      id: '',
      cooldown: 5,
      cooldownExpireAt: null,
      enabled: true,
      isRegular: false,
      isReply: true,
      text: '',
      usages: 0,
    },
  });
  const viewPort = useViewportSize();

  useEffect(() => {
    form.reset();
    if (props.keyword) {
      form.setValues(props.keyword);
    }
  }, [props.keyword]);

  const manager = useKeywordsManager();

  async function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      console.log(validate.errors);
      return;
    }

    await manager.createOrUpdate(form.values);
    props.setOpened(false);
  }

  return (
    <Drawer
      opened={props.opened}
      onClose={() => props.setOpened(false)}
      title={
        <Button size="xs" color="green" onClick={onSubmit}>
          Save
        </Button>
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
            <NumberInput label="Cooldown" required {...form.getInputProps('cooldown')} />
            <NumberInput label="Used times" {...form.getInputProps('usages')} />

            <Switch label="Use twitch reply feature" {...form.getInputProps('isReply')} />
            {props.keyword?.id && (
              <Alert>
                You can access that counter in your commands via{' '}
                <b>
                  $(keywords.counter|
                  {props.keyword?.id})
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
