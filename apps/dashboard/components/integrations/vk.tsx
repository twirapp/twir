import { Group, Avatar, Text, Button, Flex, Alert } from '@mantine/core';
import { IconBrandVk, IconLogin, IconLogout } from '@tabler/icons';

import { IntegrationCard } from './card';

import { useVkIntegration } from '@/services/api/integrations';
import {useTranslation} from "next-i18next";

export const VKIntegration: React.FC = () => {
  const manager = useVkIntegration();
  const { data: profile } = manager.getProfile();
  const { t } = useTranslation("integrations")

  async function login() {
    const link = await manager.getAuthLink();
    if (link) {
      window.location.replace(link);
    }
  }

  return (
    <IntegrationCard
      title="VK"
      icon={IconBrandVk}
      iconColor="lightblue"
      header={
        <Flex direction="row" gap="sm">
          {profile && (
            <Button
              compact
              leftIcon={<IconLogout />}
              variant="outline"
              color="red"
              onClick={manager.logout}
            >
              {t("logout")}
            </Button>
          )}
          <Button compact leftIcon={<IconLogin />} variant="outline" color="green" onClick={login}>
            {t("login")}
          </Button>
        </Flex>
      }
    >
      {!profile && <Alert>{t("notLoggedIn")}</Alert>}
      {profile && (
        <Group position="apart" mt={10}>
          <Text weight={500} size={30}>
            {profile.first_name} {profile.last_name}
          </Text>
          <Avatar src={profile.photo_max_orig} h={150} w={150} style={{ borderRadius: 900 }} />
        </Group>
      )}
    </IntegrationCard>
  );
};
