import { showNotification } from '@mantine/notifications';

export const printError = (message: string | string[]) => {
  showNotification({
    title: 'Oops',
    message: (
      <div>
        {Array.isArray(message) &&
          message.map((m) => m.charAt(0).toUpperCase() + m.slice(1)).join(', ')}
        {!Array.isArray(message) && message}
      </div>
    ),
    color: 'red',
  });
};
