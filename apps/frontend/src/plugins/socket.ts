import io, { ManagerOptions, Socket, SocketOptions } from 'socket.io-client';

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

export const NAMESPACES = {
  YOUTUBE: 'youtube',
};

export const socket = io(baseUrl, options);
export const nameSpaces: Map<string, Socket> = new Map();

nameSpaces.set(NAMESPACES.YOUTUBE, io(`${baseUrl}/youtube`, options));

function connect() {
  socket.removeAllListeners().disconnect().connect();

  nameSpaces.get(NAMESPACES.YOUTUBE)?.removeAllListeners().disconnect().connect();
}

selectedDashboardStore.subscribe(() => {
  connect();
});

userStore.subscribe(() => {
  connect();
});
