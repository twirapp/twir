import SteamUser from 'steam-user';

export type RP = SteamUser.RichPresence & {
  watching_server?: string,
  param0?: string, // lobby type
  param1?: string,
  param2?: string, // hero
  WatchableGameID?: string,
}