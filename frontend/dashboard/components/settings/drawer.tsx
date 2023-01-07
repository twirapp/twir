import { Button, Drawer, Flex, TextInput, useMantineTheme } from '@mantine/core';
import { useForm } from '@mantine/form';
import { useTranslation } from 'next-i18next';
import { useEffect } from 'react';

import { noop } from '../../util/chore';

import { dashboardAccessManager } from '@/services/api';

type Props = {
  opened: boolean;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const DashboardAccessDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<{ userName: string, id?: string }>({
    initialValues: {
      id: '',
      userName: '',
    },
  });

  const { t } = useTranslation('settings');
  const manager = dashboardAccessManager();
  const updater = manager.useCreateOrUpdate();

  useEffect(() => {
    form.reset();
  }, [props.opened]);

  async function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      console.log(validate.errors);
      return;
    }

    updater.mutateAsync({
      id: form.values.id,
      data: form.values as any,
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
          {t('dashboardAccess.drawer.save')}
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
      <form>
        <Flex direction="column" gap="md" justify="flex-start" align="flex-start" wrap="wrap">
          <TextInput {...form.getInputProps('userName')} label={t('dashboardAccess.drawer.userName')} required></TextInput>
        </Flex>
      </form>
    </Drawer>
  );
};
