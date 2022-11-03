import io, { ManagerOptions, SocketOptions } from 'socket.io-client';

import { selectedDashboardStore, userStore } from '@/stores/userStore';

const baseUrl = `${window.location.protocol == 'https:' ? 'wss' : 'ws'}://${window.location.host}`;

const authCb = (cb: (data: Record<string, any>) => void) => {
  const user = userStore.get();
  if (!user) {
    return;
  }
  const dashboard = selectedDashboardStore.get();
  if (!dashboard) {
    throw new Error('Dashboard not selected');
  }

  cb({ apiKey: user.apiKey, channelId: dashboard.channelId });
};

const options: Partial<ManagerOptions & SocketOptions> = {
  transports: ['websocket'],
  autoConnect: false,
  auth: authCb,
};

export const socket = io(baseUrl, options);
const youtube = io(`${baseUrl}/youtube`, options);

function connect() {
  socket.removeAllListeners().disconnect().connect();

  youtube.removeAllListeners().disconnect().connect();
  youtube.on('currentQueue', (d) => console.log('getting', d));
}

selectedDashboardStore.subscribe(() => {
  connect();
});

userStore.subscribe(() => {
  connect();
});
