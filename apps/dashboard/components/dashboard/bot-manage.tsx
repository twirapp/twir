import { Alert, Button, Card, createStyles, Skeleton, Text } from '@mantine/core';
import { IconAlertCircle, IconCheck, IconLogin, IconLogout } from '@tabler/icons';
import { useTranslation } from 'next-i18next';

import { useBotApi } from '@/services/api/bot';

const useStyles = createStyles((theme) => ({
  card: {
    backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.white,
  },
  title: {
    lineHeight: 1,
  },
}));

export const BotManage = () => {
  const { classes } = useStyles();
  const { t } = useTranslation('dashboard');
  const botApi = useBotApi();
  const { data: botInfo } = botApi.botInfo();
  const manager = botApi.useChangeState();

  return (
    <Skeleton visible={botInfo?.isMod === undefined}>
      <Card withBorder radius="md" p="xl" className={classes.card}>
        <Text size="lg" className={classes.title} weight={500}>
          {t('widgets.bot.title')}
        </Text>

        <Card.Section pt="lg" p="lg">
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
            mt="lg"
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
