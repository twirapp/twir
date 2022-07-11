import SteamID from 'steamid';

import { RichPresence } from '../types.js';

export function converUsers(users: Record<string, { richPresence: RichPresence }>) {
  return Object.entries(users).map(user => {
    const USERRP = user[1].richPresence;
    return {
      userId: new SteamID(user[0]).accountid,
      steamId: user[0],
      richPresence: {
        ...USERRP,
        watching_server: USERRP.watching_server ? new SteamID(USERRP.watching_server).getSteamID64() : USERRP.watching_server,
        createdAt: new Date(),
        lobbyId: USERRP.WatchableGameID,
      },
    };
  });
}