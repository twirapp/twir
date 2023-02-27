import { User } from '@tsuwari/typeorm/entities/User';
import { Server, Socket } from 'socket.io';
import { type ExtendedError } from 'socket.io/dist/namespace';

import { typeorm } from './typeorm.js';

export const io = new Server();

export const authMiddleware = async (socket: Socket, next: (err?: ExtendedError) => void) => {
  const handshake = socket.handshake;
  const apiKey = handshake.auth?.apiKey as string | undefined;

  if (!apiKey) {
    return next(new Error('Apikey not provided'));
  }

  const user = await typeorm.getRepository(User).findOneBy({
    apiKey,
  }).catch(() => null);

  if (!user) {
    return next(new Error('User with that token not found'));
  }

  socket.data.channelId = user.id;

  return next();
};


io.use(authMiddleware);