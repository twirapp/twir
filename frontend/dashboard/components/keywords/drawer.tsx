import {
  Alert,
  Button,
  Divider,
  Drawer,
  Flex,
  Grid,
  Modal,
  NumberInput,
  ScrollArea,
  Switch,
  Textarea,
  TextInput,
  Title,
  Text,
  useMantineTheme,
  Paper,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import { ChannelKeyword } from '@twir/typeorm/entities/ChannelKeyword';
import { useTranslation } from 'next-i18next';
import { useEffect } from 'react';

import { noop } from '../../util/chore';

import { useKeywordsManager } from '@/services/api';

type Props = {
  opened: boolean;
  keyword?: ChannelKeyword;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const KeywordModal: React.FC<Props> = (props) => {
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
  const { t } = useTranslation('keywords');

  useEffect(() => {
    form.reset();
    if (props.keyword) {
      form.setValues(props.keyword);
    }
  }, [props.keyword, props.opened]);

  const { useCreateOrUpdate } = useKeywordsManager();
  const updater = useCreateOrUpdate();

  async function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      console.log(validate.errors);
      return;
    }

    await updater
      .mutateAsync({
        id: form.values.id,
        data: form.values,
      })
      .then(() => {
        props.setOpened(false);
        form.reset();
      })
      .catch(noop);
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
        <Flex
          w={'100%'}
          direction="column"
          gap="md"
          justify="flex-start"
          align="flex-start"
          wrap="wrap"
        >
          <TextInput
            {...form.getInputProps('text')}
            label={t('trigger')}
            required
            w="100%"
            placeholder={'Type trigger'}
          />
          <Switch
            label={t('drawer.isRegular')}
            {...form.getInputProps('isRegular', { type: 'checkbox' })}
          />
          {form.values.isRegular && <Alert w={'100%'}>{t('drawer.expressionAlert')}</Alert>}
          <Textarea
            {...form.getInputProps('response')}
            label={t('response')}
            autosize={true}
            w="100%"
            placeholder={'Text here'}
          />

          <Divider w={'100%'} label={<Title order={3}>Settings</Title>} />

          <Grid w={'100%'}>
            <Grid.Col span={6}>
              <NumberInput label={t('cooldown')} required {...form.getInputProps('cooldown')} />
            </Grid.Col>
            <Grid.Col span={6}>
              <NumberInput label={t('usages')} {...form.getInputProps('usages')} />
            </Grid.Col>
            <Grid.Col span={6}>
              <Paper shadow="xs" p="xs" withBorder>
                <Flex direction={'row'} align={'center'} justify={'space-between'}>
                  <Text>{t('drawer.useReply')}</Text>
                  <Switch {...form.getInputProps('isReply', { type: 'checkbox' })} />
                </Flex>
              </Paper>
            </Grid.Col>
          </Grid>
        </Flex>
      </form>
    </Modal>
  );
};
