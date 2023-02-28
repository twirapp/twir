import {
  ActionIcon, Alert,
  Button,
  Drawer,
  Flex,
  Menu,
  ScrollArea,
  Switch, Text,
  Textarea,
  TextInput,
  useMantineTheme,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { IconVariable } from '@tabler/icons';
import { ChannelGreeting } from '@tsuwari/typeorm/entities/ChannelGreeting';
import { useTranslation } from 'next-i18next';
import { useEffect } from 'react';

import { noop } from '../../util/chore';

import { type Greeting, greetingsManager, useVariables } from '@/services/api';

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
  const { t } = useTranslation('greetings');

  const { useCreateOrUpdate } = greetingsManager();
  const updater = useCreateOrUpdate();

  const variables = useVariables();

  useEffect(() => {
    form.reset();
    if (props.greeting) {
      form.setValues(props.greeting);
    }
  }, [props.greeting, props.opened]);

  async function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      console.log(validate.errors);
      return;
    }

    await updater.mutateAsync({
      id: form.values.id,
      data: form.values,
    }).then(() => {
      props.setOpened(false);
      form.reset();
    }).catch(noop);
  }

  return (
    <Drawer
      opened={props.opened}
      onClose={() => props.setOpened(false)}
      title={
        <Button size="xs" color="green" onClick={onSubmit}>
          {t('drawer.save')}
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
            <Alert><Text size={'xs'}>Tip: you can create tts for greeting in events.</Text></Alert>
            <TextInput label={t('userName')} required {...form.getInputProps('userName')} />
            <Textarea
              label={t('message')}
              required
              w="100%"
              autosize={true}
              minRows={1}
              rightSection={
                <Menu position="right-end" shadow="md" width={250}>
                  <Menu.Target>
                    <ActionIcon variant="filled">
                      <IconVariable size={18} />
                    </ActionIcon>
                  </Menu.Target>

                  <Menu.Dropdown>
                    <ScrollArea h={200} type={'always'} offsetScrollbars>
                      {variables.data?.length && variables.data.map(v => (
                        <Menu.Item key={v.name} onClick={() => {
                          const insertValue = `${v.example ? v.example : v.name}`;
                          form.setFieldValue(
                            `text`,
                            `${form.values.text} $(${insertValue})`,
                          );
                        }}>
                          <Flex direction={'column'}>
                            <Text>{v.name}</Text>
                            <Text size={'xs'}>{v.description}</Text>
                          </Flex>
                        </Menu.Item>
                      ))}

                    </ScrollArea>
                  </Menu.Dropdown>

                </Menu>
              }
              {...form.getInputProps('text')}
            />
            <Switch
              label={t('drawer.useReply')}
              {...form.getInputProps('isReply', { type: 'checkbox' })}
            />
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
