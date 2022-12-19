import {
  Badge,
  Button,
  Drawer,
  Flex,
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
import { useEffect, useRef } from 'react';

import { useVariablesManager } from '@/services/api';

type Props = {
  opened: boolean;
  variable?: ChannelCustomvar;
  setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

export const VariableDrawer: React.FC<Props> = (props) => {
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

  function handleEditorDidMount(editor: any) {
    editorRef.current = editor;
  }

  useEffect(() => {
    form.reset();
    if (props.variable) {
      form.setValues(props.variable);
    }
  }, [props.variable]);

  const manager = useVariablesManager();
  async function onSubmit() {
    const validate = form.validate();
    if (validate.hasErrors) {
      console.log(validate.errors);
      return;
    }

    await manager.createOrUpdate(form.values);
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
      <ScrollArea.Autosize maxHeight={viewPort.height - 120} type="auto" offsetScrollbars={true}>
        <form onSubmit={form.onSubmit((values) => console.log(values))}>
          <Flex direction="column" gap="md" justify="flex-start" align="flex-start" wrap="wrap">
            <TextInput label="Name" required {...form.getInputProps('name')} />
            <Select
              label="Variable type"
              data={[
                { value: 'SCRIPT', label: 'Script' },
                { value: 'TEXT', label: 'Text' },
              ]}
              {...form.getInputProps('type')}
            />
            {form.values.type === 'SCRIPT' && (
              <Editor
                height="50vh"
                defaultLanguage="javascript"
                theme={theme.colorScheme === 'dark' ? 'vs-dark' : 'light'}
                defaultValue="// some comment"
                onMount={handleEditorDidMount}
                value={form.values.evalValue}
                onChange={(v) => {
                  form.values.evalValue = v ?? '';
                }}
              />
            )}
            {form.values.type === 'TEXT' && (
              <Textarea label="Response" autosize={true} {...form.getInputProps('response')} />
            )}
          </Flex>
        </form>
      </ScrollArea.Autosize>
    </Drawer>
  );
};
