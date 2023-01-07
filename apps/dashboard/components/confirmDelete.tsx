import { Text } from '@mantine/core';
import { openConfirmModal } from '@mantine/modals';

export type Props = {
  onCancel?: () => Promise<any> | any;
  onConfirm: () => Promise<any> | any;
  text?: string
  title?: string
};

export const confirmDelete = (props: Props) => {
  openConfirmModal({
    title: props.title ?? 'Please confirm your action',
    children: (
      <Text size="sm">{props.text ?? 'This action is so important that you are required to confirm it.'}</Text>
    ),
    labels: { confirm: 'Confirm', cancel: 'Cancel' },
    onCancel: props.onCancel,
    onConfirm: props.onConfirm,
  });
};
