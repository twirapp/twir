import { Alert, Button, Card, Divider, Flex, Grid, Group, Text } from '@mantine/core';
import { IconLogin, IconLogout } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import React from 'react';

import { useBotApi } from '@/services/api/bot';

export const BotWidget: React.FC = () => {
  const { t } = useTranslation('dashboard');
  const botApi = useBotApi();
  const { data: isBotMod } = botApi.isMod();

  return <Grid grow>
    <Grid.Col span={4}>
      <Card>
        <Card.Section p="sm">
          <Flex gap="xs" direction="row" justify="space-between">
            {/*{props.icon && <props.icon color={props.iconColor} />}*/}
            <Text size="md">{t('widgets.bot.title')}</Text>
          </Flex>
        </Card.Section>
        <Divider />
        <Card.Section p="lg">
          {isBotMod && <Alert color="teal"><span dangerouslySetInnerHTML={{ __html:t('widgets.bot.alert.true') }} /></Alert>}
          {!isBotMod && <Alert color="red"><span dangerouslySetInnerHTML={{ __html:t('widgets.bot.alert.false') }} /></Alert>}
        </Card.Section>
        <Card.Section p={'lg'} pt={0}>
          <Grid grow style={{ marginTop: 5 }}>
            <Grid.Col span={4}>
              <Button variant="subtle" size="lg" w={'100%'} color={'red'} leftIcon={<IconLogin />}>Leave</Button>
            </Grid.Col>
            <Grid.Col span={4}>
              <Button variant="subtle" size="lg" w={'100%'} color={'teal'} leftIcon={<IconLogout />}>Join</Button>
            </Grid.Col>
          </Grid>
        </Card.Section>
      </Card>
    </Grid.Col>
  </Grid>;
};
