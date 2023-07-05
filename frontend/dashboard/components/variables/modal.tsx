import { Alert, Button, Grid, Modal, NumberInput, Select, Textarea, TextInput, useMantineTheme } from '@mantine/core';
import { useForm } from '@mantine/form';
import Editor from '@monaco-editor/react';
import { type Variable, VariableType } from '@twir/grpc/generated/api/api/variables';
import { useTranslation } from 'next-i18next';
import { Fragment, useEffect, useRef } from 'react';

import { useVariablesManager } from '@/services/api';

type Props = {
  opened: boolean;
  variable?: Variable;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const VariableModal: React.FC<Props> = (props) => {
  const theme = useMantineTheme();
  const form = useForm<Variable>({
    initialValues: {
      description: '',
      evalValue: '',
      name: '',
      response: '',
      type: VariableType.TEXT,
			channelId: '',
    },
  });
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

  const manager = useVariablesManager();
  const updater = manager.update;
	const creator = manager.create;

  async function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      console.log(validate.errors);
      return;
    }

		const data = {
			...form.values,
			type: Number(form.values.type),
		};

		if (form.values.id) {
			await updater.mutateAsync({
				id: form.values.id,
				variable: data,
			});
		} else {
			await creator.mutateAsync({
				variable: data,
			});
		}

		props.setOpened(false);
		form.reset();
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
      size={form.values.type === VariableType.SCRIPT ? '80%' : 'xl'}
      overlayColor={theme.colorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[2]}
      overlayOpacity={0.55}
      overlayBlur={3}
			closeOnClickOutside={false}
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
                { value: VariableType.SCRIPT.toString(), label: 'Script' },
                { value: VariableType.TEXT.toString(), label: 'Text' },
                { value: VariableType.NUMBER.toString(), label: 'Number' },
              ]}
              {...form.getInputProps('type')}
              dropdownPosition={'bottom'}
              zIndex={999}
            />
          </Grid.Col>
          <Grid.Col span={12}>
            {form.values.type === VariableType.SCRIPT && (
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
            {form.values.type === VariableType.TEXT && (
              <Textarea label={t('response')} autosize={true} {...form.getInputProps('response')} />
            )}
            {form.values.type === VariableType.NUMBER && (
              <NumberInput label={t('response')} {...form.getInputProps('response')} />
            )}
          </Grid.Col>
        </Grid>
      </form>
    </Modal>
  );
};
