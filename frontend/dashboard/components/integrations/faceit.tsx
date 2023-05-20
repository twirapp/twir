import { Alert, Avatar, Button, Flex, Group, Text } from '@mantine/core';
import { IconInfoCircle, IconLogin, IconLogout } from '@tabler/icons';
import { useTranslation } from 'next-i18next';

import { IntegrationCard } from './card';

import { useFaceit } from '@/services/api/integrations';

const noAvatar =
  'data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQECAgICAgICAgICAgMDAwMDAwMDAwP/2wBDAQEBAQEBAQEBAQECAgECAgMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwP/wAARCAAoACgDAREAAhEBAxEB/8QAGwAAAgIDAQAAAAAAAAAAAAAABQYEBwEDCAn/xAA3EAABAwIEAgYFDQAAAAAAAAABAgMEBREAEiFBBjETFCIjYXEVUYGxtAckMjNCRGNkkZShwdT/xAAUAQEAAAAAAAAAAAAAAAAAAAAA/8QAFBEBAAAAAAAAAAAAAAAAAAAAAP/aAAwDAQACEQMRAD8A8Hn6+9GkraZeShsBrKnoWFc2UKNitsqPaJ3wGtXEc9Q1kIt6+gi38rdCMAOf4kmpv84Rvzjxj69uh8MADl8WzEBVn273OvVol9L6ABjTlvgFo8WyZMsMPPtqbW3IzIDEdN8kd1ae0hpKgcwSdDgBs3iRKZThLmzVze+oaQB5AFOAjK4nbt9YANRooeI0vuMBGfrbqmFOpSsthBXe4+iASVZeZ7JwCNUeK0DN3g12BufL9cAntcWD0i13luxJ1KtTeK9t4j3YAjWapllvDPbRvzHctn+zgFOTW3EXyunzJB9ut98A/IrqhwS9VerulbEJcfIEq7ZQoRBISQc3QgKzqULZQFbDAc8y6285cF1RHnYba6W2wA6nVEmpM3WSckn4V718+RwFgcQ1TLUJAzbM6X01YaPvP84AVRuKKbSpbkqpUtNWAbAjtrcSlDLgUFF3K4262tWXQXTcWuMA4ufLbT2k5fQDxSBlyie0E5QCLW6ra1hgKc414upFekx5NKoqKKtKHUyw26hSZS1FPRuZGmmWkKQEquQm6r6nTAJtJn5qrHBV9iX8HIPuOAsiuNzZM16QwhDjTiWShXWIyL2YaSo5VvJWLLTuMApyYdV17hG/3uH4/j8tcADkQKsq9mEeyZC8fzHjgA71Lq5v3Df72F/pwGaZT6jFntSZDbTbLbcnOvrcRds0V5Cew2+pZutQGgPPAf/Z';

export const FaceitIntegration: React.FC = () => {
  const manager = useFaceit();
  const { data: profile } = manager.useData();
  const { t } = useTranslation('integrations');
  const auth = manager.useGetAuthLink();
  const logout = manager.useLogout();

  async function login() {
    if (auth.data) {
      window.location.replace(auth.data);
    }
  }

  return (
    <IntegrationCard
      title="Faceit"
      header={
        <Flex direction="row" gap="sm">
          {profile && (
            <Button
              compact
              leftIcon={<IconLogout />}
              variant="outline"
              color="red"
              onClick={() => logout.mutate()}
            >
              {t('logout')}
            </Button>
          )}
          <Button compact leftIcon={<IconLogin />} variant="outline" color="green" onClick={login}>
            {t('login')}
          </Button>
        </Flex>
      }
    >
      {!profile && <Alert>{t('notLoggedIn')}</Alert>}
      {profile && (
        <Group position="apart" mt={10}>
          <Text weight={500} size={30}>
            {profile.name}
          </Text>
          <Avatar src={profile.avatar ?? noAvatar} h={150} w={150} style={{ borderRadius: 900 }} />
        </Group>
      )}

      <Alert color={'lime'} icon={<IconInfoCircle />} mt={5}>
        <Text dangerouslySetInnerHTML={{ __html: t('info.faceit') }} />
      </Alert>
    </IntegrationCard>
  );
};
