import {
  Alert,
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
import { useTranslation } from 'next-i18next';
import { useEffect } from 'react';

import { noop } from '../../util/chore';

import { keywordsManager } from '@/services/api';

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
  const { t } = useTranslation('keywords');

  useEffect(() => {
    form.reset();
    if (props.keyword) {
      form.setValues(props.keyword);
    }
  }, [props.keyword]);

  const manager = keywordsManager();

  async function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      console.log(validate.errors);
      return;
    }

    await manager.createOrUpdate.mutateAsync({
      id: form.values.id,
      data: form.values,
    })
      .then(() => {
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
            <TextInput {...form.getInputProps('text')} label={t('trigger')} required w="100%" />
            <Switch label={t('drawer.isRegular')} {...form.getInputProps('isRegular')} />
            {form.values.isRegular && (
              <Alert>
                {t('drawer.expressionAlert')}
              </Alert>
            )}
            <Textarea
              {...form.getInputProps('response')}
              label={t('response')}
              autosize={true}
              w="100%"
            />
            <NumberInput label={t('cooldown')} required {...form.getInputProps('cooldown')} />
            <NumberInput label={t('usages')} {...form.getInputProps('usages')} />

            <Switch label={t('drawer.useReply')} {...form.getInputProps('isReply')} />
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
