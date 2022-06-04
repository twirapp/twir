import type { ClientToServerEvents, EventParams, ServerToClientEvents } from '@tsuwari/shared';
import jwtDecode from 'jwt-decode';
import { io, Socket } from 'socket.io-client';

import { refreshAccessToken } from '../functions/refreshAccessToken.js';
import { selectedDashboardStore } from '../stores/userStore.js';

const url = import.meta.env.DEV ? 'http://localhost:3002' : `ws://${window.location.host}/api`

export const socket: Socket<ServerToClientEvents, ClientToServerEvents> = io(url, {
  auth: (cb) => {
    cb({ token: localStorage.getItem('accessToken') });
  },
});

selectedDashboardStore.subscribe(v => {
  if (!v?.channelId || !v?.userId) return
  socket.io.opts.query = {
    channelId: v.channelId,
    userId: v.userId,
  };
  socket.disconnect().connect();
});

export async function socketEmit<EV extends keyof ClientToServerEvents>(event: EV, ...params: EventParams<ClientToServerEvents, EV>) {
  const accessToken = localStorage.getItem('accessToken');
  if (!accessToken) {
    throw new Error('No access token');
  }
  const decodedJwt = jwtDecode(accessToken) as any;

  if (!decodedJwt.exp) throw new Error('Token have to exp');

  const exp = decodedJwt.exp * 1000;
  if (new Date().getTime() >= exp) {
    await refreshAccessToken();
    socket.disconnect().connect().emit(event, ...params);
  } else {
    socket.emit(event, ...params);
  }
}

socket.on('connect_error', (err) => {
  console.log(`connect_error due to ${err}`);
});