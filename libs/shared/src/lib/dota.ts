export type DotaGame = {
  activate_time: number,
  lobby_type: number,
  game_mode: number,
  average_mmr: number,
  players: Array<{ account_id: number, hero_id: number }>,
  weekend_tourney_bracket_round?: string,
  weekend_tourney_skill_level?: string,
  match_id?: string
}