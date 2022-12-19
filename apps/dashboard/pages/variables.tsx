import { ActionIcon, Badge, Button, Flex, Table, Text } from '@mantine/core';
import { IconPencil, IconTrash } from '@tabler/icons';
import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { useState } from 'react';

import { confirmDelete } from '@/components/confirmDelete';
import { VariableDrawer } from '@/components/variables/drawer';
import { useVariablesManager } from '@/services/api';

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableVariable, setEditableVariable] = useState<ChannelCustomvar | undefined>();

  const manager = useVariablesManager();
  const { data: variables } = manager.getCreated();

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Text size="lg">Variables</Text>
        <Button
          color="green"
          onClick={() => {
            setEditableVariable(undefined);
            setEditDrawerOpened(true);
          }}
        >
          Create
        </Button>
      </Flex>
      <Table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Response</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {variables &&
            variables.map((variable, idx) => (
              <tr key={variable.id}>
                <td>
                  <Badge>{variable.name}</Badge>
                </td>
                <td>
                  <Badge color="cyan">{variable.type}</Badge>
                </td>
                <td>
                  {variable.type === 'TEXT' && <Badge>{variable.response}</Badge>}
                  {variable.type !== 'TEXT' && (
                    <Badge color="red">Script cannot be displayed</Badge>
                  )}
                </td>
                <td>
                  <Flex direction="row" gap="xs">
                    <ActionIcon
                      onClick={() => {
                        setEditableVariable(variables[idx] as any);
                        setEditDrawerOpened(true);
                      }}
                      variant="filled"
                      color="blue"
                    >
                      <IconPencil size={14} />
                    </ActionIcon>

                    <ActionIcon
                      onClick={() =>
                        confirmDelete({
                          onConfirm: () => manager.delete(variable.id),
                        })
                      }
                      variant="filled"
                      color="red"
                    >
                      <IconTrash size={14} />
                    </ActionIcon>
                  </Flex>
                </td>
              </tr>
            ))}
        </tbody>
      </Table>

      <VariableDrawer
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        variable={editableVariable}
      />
    </div>
  );
}
