import {
  Alert,
  Button,
  Drawer,
  Flex,
  Grid,
  Modal,
  NumberInput,
  ScrollArea,
  Select,
  Textarea,
  TextInput,
  useMantineTheme,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useViewportSize } from '@mantine/hooks';
import Editor from '@monaco-editor/react';
import { ChannelCustomvar, CustomVarType } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { useTranslation } from 'next-i18next';
import { Fragment, useEffect, useRef } from 'react';

import { noop } from '../../util/chore';

import { variablesManager } from '@/services/api';

type Props = {
  opened: boolean;
  variable?: ChannelCustomvar;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const VariableModal: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<ChannelCustomvar>({
    initialValues: {
      id: '',
      description: '',
      evalValue: '',
      name: '',
      response: '',
      type: 'TEXT' as CustomVarType,
    },
  });
  const viewPort = useViewportSize();
  const editorRef = useRef(null);
  const { t } = useTranslation('variables');

  function handleEditorDidMount(editor: any) {
    editorRef.current = editor;
  }

  useEffect(() => {
    form.reset();
    if (props.variable) {
      form.setValues(props.variable);
    }
  }, [props.variable, props.opened]);

  const { useCreateOrUpdate } = variablesManager();
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
        data: {
          ...form.values,
          response: form.values.response.toString(),
        },
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
      size={form.values.type === 'SCRIPT' ? '80%' : 'xl'}
      transition="slide-left"
      overlayColor={theme.colorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[2]}
      overlayOpacity={0.55}
      overlayBlur={3}
    >
      <form onSubmit={form.onSubmit((values) => console.log(values))}>
        <Grid w={'100%'}>
          <Grid.Col span={6}>
            <TextInput label={t('name')} required {...form.getInputProps('name')} />
          </Grid.Col>
          <Grid.Col span={6}>
            <Select
              label={t('type')}
              data={[
                { value: 'SCRIPT', label: 'Script' },
                { value: 'TEXT', label: 'Text' },
                { value: 'NUMBER', label: 'Number' },
              ]}
              {...form.getInputProps('type')}
              dropdownPosition={'bottom'}
              zIndex={999}
            />
          </Grid.Col>
          <Grid.Col span={12}>
            {form.values.type === 'SCRIPT' && (
              <Fragment>
                <Alert mb={10}>{t('drawer.scriptAlert')}</Alert>
                <Editor
                  height="50vh"
                  defaultLanguage="javascript"
                  theme={theme.colorScheme === 'dark' ? 'vs-dark' : 'light'}
                  defaultValue="//"
                  onMount={handleEditorDidMount}
                  value={form.values.evalValue}
                  onChange={(v) => {
                    form.values.evalValue = v ?? '';
                  }}
                />
              </Fragment>
            )}
            {form.values.type === 'TEXT' && (
              <Textarea label={t('response')} autosize={true} {...form.getInputProps('response')} />
            )}
            {form.values.type === 'NUMBER' && (
              <NumberInput label={t('response')} {...form.getInputProps('response')} />
            )}
          </Grid.Col>
        </Grid>
      </form>
    </Modal>
  );
};
