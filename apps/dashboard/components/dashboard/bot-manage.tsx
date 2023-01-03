import { useBotApi } from '@/services/api/bot';
import { Alert, Button, Card, createStyles, Grid, Skeleton, Text } from '@mantine/core';
import { IconAlertCircle, IconCheck, IconLogin, IconLogout } from '@tabler/icons';
import { useTranslation } from 'next-i18next';

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
  const { data: isBotMod } = botApi.isMod();
  const manager = botApi.useChangeState();

  return (
    <Skeleton visible={isBotMod === undefined}>
      <Card withBorder radius="md" p="xl" className={classes.card}>
        <Text size="lg" className={classes.title} weight={500}>
          {t('widgets.bot.title')}
        </Text>

        <Card.Section pt="lg" p="lg">
          {isBotMod ? (
            <Alert icon={<IconCheck size={16} />} color="teal" variant="outline">
              <span dangerouslySetInnerHTML={{ __html: t('widgets.bot.alert.true') }} />
            </Alert>
          ) : (
            <Alert icon={<IconAlertCircle size={16} />} color="red" variant="outline">
              <span dangerouslySetInnerHTML={{ __html: t('widgets.bot.alert.false') }} />
            </Alert>
          )}

          <Grid grow mt="xs">
            <Grid.Col span={4}>
              <Button
                size="md"
                w="100%"
                color="red"
                leftIcon={<IconLogin />}
                onClick={() => {
                  manager.mutate('part');
                }}
              >
                {t('widgets.bot.actions.leave')}
              </Button>
            </Grid.Col>
            <Grid.Col span={4}>
              <Button
                size="md"
                w="100%"
                color="teal"
                leftIcon={<IconLogout />}
                onClick={() => {
                  manager.mutate('join');
                }}
              >
                {t('widgets.bot.actions.join')}
              </Button>
            </Grid.Col>
          </Grid>
        </Card.Section>
      </Card>
    </Skeleton>
  );
};
