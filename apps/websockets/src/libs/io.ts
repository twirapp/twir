import { User } from '@tsuwari/typeorm/entities/User';
import { Server, Socket } from 'socket.io';
import { ExtendedError } from 'socket.io/dist/namespace';

import { typeorm } from './typeorm';

export const io = new Server();

export const authMiddleware = (socket: Socket, next: (err?: ExtendedError) => void) => {
  const handshake = socket.handshake;
  const apiKey = handshake.auth?.apiKey as string | undefined;

  if (!apiKey) {
    return next(new Error('Apikey not provided'));
  }

  const user = typeorm.getRepository(User).findOneBy({
    apiKey,
  });

  if (!user) {
    return next(new Error('User with that token not found'));
  }

  return next();
};


io.use(authMiddleware);