import { Button, Drawer, Flex, TextInput, useMantineTheme } from '@mantine/core';
import { useForm } from '@mantine/form';
import { useEffect } from 'react';

import { useDashboardAccess } from '@/services/api';

type Props = {
  opened: boolean;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const DashboardAccessDrawer: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<{ userName: string }>({
    initialValues: {
      userName: '',
    },
  });

  const manager = useDashboardAccess();

  useEffect(() => {
    form.reset();
  }, [props.opened]);

  async function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      console.log(validate.errors);
      return;
    }

    await manager.create(form.values.userName);
    props.setOpened(false);
  }

  return (
    <Drawer
      opened={props.opened}
      onClose={() => props.setOpened(false)}
      title={
        <Button size="xs" color="green" onClick={onSubmit}>
          Save
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
          <TextInput {...form.getInputProps('userName')} label="User name" required></TextInput>
        </Flex>
      </form>
    </Drawer>
  );
};
