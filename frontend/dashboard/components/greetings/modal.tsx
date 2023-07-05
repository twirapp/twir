import {
  ActionIcon, Alert,
  Button,
  Flex, Grid,
  Menu, Modal, Paper,
  ScrollArea,
  Switch, Text,
  Textarea,
  TextInput,
  useMantineTheme,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { IconVariable } from '@tabler/icons';
import type { Greeting } from '@twir/grpc/generated/api/api/greetings';
import { useTranslation } from 'next-i18next';
import { useEffect } from 'react';

import { noop } from '../../util/chore';

import { useAllVariables } from '@/services/api/allVariables.js';
import { useGreetingsManager } from '@/services/api/index.js';

type Props = {
  opened: boolean;
  greeting?: Greeting;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const GreetingModal: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<Greeting>({
    initialValues: {
      id: '',
      channelId: '',
      processed: false,
      isReply: true,
      userId: '',
      text: '',
      enabled: true,
    },
  });

  const { t } = useTranslation('greetings');

  const greetingsManager = useGreetingsManager();
  const updater = greetingsManager.update!;
	const creator = greetingsManager.create;
	const allVariables = useAllVariables();

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

		if (form.values.id) {
			await updater.mutateAsync(form.values);
		} else {
			await creator.mutateAsync(form.values);
		}

		props.setOpened(false);
		form.reset();
  }

  return (
    <Modal
      opened={props.opened}
      onClose={() => props.setOpened(false)}
      title={
        <Button size="xs" color="green" onClick={onSubmit}>
          {t('drawer.save')}
        </Button>
      }
      padding="xl"
      size="xl"
      overlayColor={theme.colorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[2]}
      overlayOpacity={0.55}
      overlayBlur={3}
			closeOnClickOutside={false}
    >
        <form onSubmit={form.onSubmit((values) => console.log(values))}>
          <Grid>
            <Grid.Col span={12}>
              <TextInput label={t('userName')} required {...form.getInputProps('userName')} />
            </Grid.Col>
            <Grid.Col span={12}>
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
                        {allVariables.map(v => (
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
            </Grid.Col>
            <Grid.Col span={6}>
              <Paper shadow="xs" p="xs" withBorder>
                <Flex direction={'row'} align={'center'} justify={'space-between'}>
                  <Text>{t('drawer.useReply')}</Text>
                  <Switch
                    {...form.getInputProps('isReply', { type: 'checkbox' })}
                  />
                </Flex>
              </Paper>
            </Grid.Col>
          </Grid>
          <Alert mt={10}><Text size={'xs'}>Tip: you can create tts for greeting in events.</Text></Alert>
        </form>
    </Modal>
  );
};
