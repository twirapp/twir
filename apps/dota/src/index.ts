import { dirname, resolve } from 'path';
import { fileURLToPath } from 'url';

import { config } from '@tsuwari/config';
import protobufjs from 'protobufjs';
import SteamUser from 'steam-user';
import SteamID from 'steamid';

const client = new SteamUser();

const proto = await new protobufjs.Root().load(resolve(dirname(fileURLToPath(import.meta.url)), '..', 'protos', 'dota2', 'dota_gcmessages_client_watch.proto'), {
  keepCase: true,
});

client.logOn({
  accountName: config.STEAM_USERNAME,
  password: config.STEAM_PASSWORD,
});

client.on('loggedOn', () => {
  console.log('Logged into Steam as ' + client.steamID?.getSteam3RenderedID());
  client.gamesPlayed([570]);
});

type RP = SteamUser.RichPresence & {
  watching_server?: string,
  param0?: string, // lobby type
  param1?: string,
  param2?: string, // hero
  WatchableGameID?: string,
}
function converUsers(users: Record<string, { richPresence: RP }>) {
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

client.on('appLaunched', async (appid) => {
  console.log('Dota launched');

  client.requestRichPresence(570, [911977148, 1102609846].map(SteamID.fromIndividualAccountID), 'english', (error, data) => {
    const users = converUsers(data.users);

    const lobbyIds = users.filter(u => u.richPresence.watching_server).map(u => u.richPresence.watching_server);

    if (!lobbyIds.length) return;
    const type = proto.root.lookupType('CMsgClientToGCFindTopSourceTVGames');
    const newMsg = type.encode({
      search_key: '',
      league_id: 0,
      hero_id: 0,
      start_game: 0,
      game_list_index: 0,
      lobby_ids: lobbyIds,
    }).finish();

    client.sendToGC(570, 8009, {}, Buffer.from(newMsg));
  });
});

client.on('receivedFromGC', (app, msg, payload) => {
  const type = proto.root.lookupType('CMsgGCToClientFindTopSourceTVGamesResponse');

  const data = type.decode(payload).toJSON() as {
    game_list?: Array<{
      lobby_type: number,
      game_mode: number,
      average_mmr: number,
      players: Array<{ account_id: number, hero_id: number }>,
      weekend_tourney_bracket_round?: string,
      weekend_tourney_skill_level?: string,
      match_id?: string
    }>
  };

  if (data.game_list) {
    console.dir({
      app,
      msg,
      games: data.game_list,
    }, { depth: null });
  } else {
    console.log(data);
  }
});

client.on('error', (e) => {
  console.log(e);
});