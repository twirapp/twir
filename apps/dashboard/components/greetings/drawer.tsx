import {
  Badge,
  Button,
  Drawer,
  Flex,
  ScrollArea,
  Switch,
  Textarea,
  TextInput,
  useMantineTheme,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import { useEffect } from 'react';

import { useGreetingsManager, type Greeting } from '@/services/api';
import {useTranslation} from "next-i18next";

type Props = {
  opened: boolean;
  greeting?: ChannelGreeting & { userName: string };
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const GreetingDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<Greeting>({
    initialValues: {
      id: '',
      channelId: '',
      processed: false,
      isReply: true,
      userId: '',
      userName: '',
      text: '',
      enabled: true,
    },
  });
  const viewPort = useViewportSize();
  const { t } = useTranslation("greetings")

  const manager = useGreetingsManager();

  useEffect(() => {
    form.reset();
    if (props.greeting) {
      form.setValues(props.greeting);
    }
  }, [props.greeting]);

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
          {t("drawer.save")}
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
            <TextInput label={t("userName")} required {...form.getInputProps('userName')} />
            <Textarea
              label={t("message")}
              required
              w="100%"
              {...form.getInputProps('text')}
            />
            <Switch
              label={t("drawer.useReply")}
              {...form.getInputProps('isReply', { type: 'checkbox' })}
            />
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
