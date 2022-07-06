import SteamUser from 'steam-user';

export type RP = SteamUser.RichPresence & {
  watching_server?: string,
  param0?: string, // lobby type
  param1?: string,
  param2?: string, // hero
  WatchableGameID?: string,
}

export type Game = {
  lobby_type: number,
  game_mode: number,
  average_mmr: number,
  players: Array<{ account_id: number, hero_id: number }>,
  weekend_tourney_bracket_round?: string,
  weekend_tourney_skill_level?: string,
  match_id?: string
}