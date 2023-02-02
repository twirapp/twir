import { Alert, Button, Card, Group, Skeleton, Text } from '@mantine/core';
import { IconAlertCircle, IconCheck, IconLogin, IconLogout } from '@tabler/icons';
import { useTranslation } from 'next-i18next';

import { useBotApi } from '@/services/api/bot';
import { useCardStyles } from '@/styles/card';

export const BotManage = () => {
  const { t } = useTranslation('dashboard');
  const botApi = useBotApi();
  const { data: botInfo } = botApi.botInfo();
  const manager = botApi.useChangeState();
  const { classes } = useCardStyles();

  return (
    <Skeleton radius="md" visible={botInfo?.isMod === undefined} h={'100%'}>
      <Card withBorder radius="md" h={'100%'}>
        <Card.Section withBorder inheritPadding py="sm">
          <Group position="apart">
            <Text weight={500}>{t('widgets.bot.title')}</Text>
          </Group>
        </Card.Section>
        <Card.Section p="md" className={classes.card} h={'100%'}>
          {botInfo?.isMod ? (
            <Alert icon={<IconCheck size={16} />} color="teal" variant="outline">
              <span dangerouslySetInnerHTML={{ __html: t('widgets.bot.alert.true') }} />
            </Alert>
          ) : (
            <Alert icon={<IconAlertCircle size={16} />} color="red" variant="outline">
              <span
                dangerouslySetInnerHTML={{
                  __html: t('widgets.bot.alert.false', { botName: botInfo?.botName ?? '' }),
                }}
              />
            </Alert>
          )}
          <Button
            loading={manager.isLoading}
            mt="md"
            size="md"
            w="100%"
            color={botInfo?.enabled ? 'red' : 'teal'}
            leftIcon={botInfo?.enabled ? <IconLogout size={20} /> : <IconLogin size={20} />}
            onClick={() => {
              manager.mutate(botInfo?.enabled ? 'part' : 'join');
            }}
          >
            {botInfo?.enabled ? t('widgets.bot.actions.leave') : t('widgets.bot.actions.join')}
          </Button>
        </Card.Section>
      </Card>
    </Skeleton>
  );
};
