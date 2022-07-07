export type DotaGame = {
  activate_time: number,
  lobby_type: number,
  game_mode: number,
  average_mmr: number,
  players: Array<{ account_id: number, hero_id: number }>,
  weekend_tourney_bracket_round?: string,
  weekend_tourney_skill_level?: string,
  match_id?: string,
  lobby_id: string
}

export const gameModes = [
  { id: 1, name: 'All Pick' },
  { id: 2, name: 'Captains Mode' },
  { id: 3, name: 'Random Draft' },
  { id: 4, name: 'Single Draft' },
  { id: 5, name: 'All Random' },
  { id: 6, name: 'Intro' },
  { id: 7, name: 'Diretide' },
  { id: 8, name: 'Reverse Captains Mode' },
  { id: 9, name: 'The Greeviling' },
  { id: 10, name: 'Tutorial' },
  { id: 11, name: 'Mid Only' },
  { id: 12, name: 'Least Played' },
  { id: 13, name: 'New Player Pool' },
  { id: 14, name: 'Compendium Matchmaking' },
  { id: 15, name: 'Custom Game' },
  { id: 16, name: 'Captains Draft' }, { id: 18, name: 'Ability Draft' },
  { id: 19, name: 'Event Game' },
  { id: 20, name: 'All Random Deathmatch' },
  { id: 21, name: '1v1 Mid Only' },
  { id: 22, name: 'Ranked' },
  { id: 23, name: 'Turbo' },
  { id: 24, name: 'Mutation' },
];

export const dotaHeroes = [
  {
    'localized_name': 'Not Picked',
    'id': 0,
  },
  {
    'localized_name': 'Anti-Mage',
    'id': 1,
  },
  {
    'localized_name': 'Axe',
    'id': 2,
  },
  {
    'localized_name': 'Bane',
    'id': 3,
  },
  {
    'localized_name': 'Bloodseeker',
    'id': 4,
  },
  {
    'localized_name': 'Crystal Maiden',
    'id': 5,
  },
  {
    'localized_name': 'Drow Ranger',
    'id': 6,
  },
  {
    'localized_name': 'Earthshaker',
    'id': 7,
  },
  {
    'localized_name': 'Juggernaut',
    'id': 8,
  },
  {
    'localized_name': 'Mirana',
    'id': 9,
  },
  {
    'localized_name': 'Morphling',
    'id': 10,
  },
  {
    'localized_name': 'Shadow Fiend',
    'id': 11,
  },
  {
    'localized_name': 'Phantom Lancer',
    'id': 12,
  },
  {
    'localized_name': 'Puck',
    'id': 13,
  },
  {
    'localized_name': 'Pudge',
    'id': 14,
  },
  {
    'localized_name': 'Razor',
    'id': 15,
  },
  {
    'localized_name': 'Sand King',
    'id': 16,
  },
  {
    'localized_name': 'Storm Spirit',
    'id': 17,
  },
  {
    'localized_name': 'Sven',
    'id': 18,
  },
  {
    'localized_name': 'Tiny',
    'id': 19,
  },
  {
    'localized_name': 'Vengeful Spirit',
    'id': 20,
  },
  {
    'localized_name': 'Windranger',
    'id': 21,
  },
  {
    'localized_name': 'Zeus',
    'id': 22,
  },
  {
    'localized_name': 'Kunkka',
    'id': 23,
  },
  {
    'localized_name': 'Lina',
    'id': 25,
  },
  {
    'localized_name': 'Lion',
    'id': 26,
  },
  {
    'localized_name': 'Shadow Shaman',
    'id': 27,
  },
  {
    'localized_name': 'Slardar',
    'id': 28,
  },
  {
    'localized_name': 'Tidehunter',
    'id': 29,
  },
  {
    'localized_name': 'Witch Doctor',
    'id': 30,
  },
  {
    'localized_name': 'Lich',
    'id': 31,
  },
  {
    'localized_name': 'Riki',
    'id': 32,
  },
  {
    'localized_name': 'Enigma',
    'id': 33,
  },
  {
    'localized_name': 'Tinker',
    'id': 34,
  },
  {
    'localized_name': 'Sniper',
    'id': 35,
  },
  {
    'localized_name': 'Necrophos',
    'id': 36,
  },
  {
    'localized_name': 'Warlock',
    'id': 37,
  },
  {
    'localized_name': 'Beastmaster',
    'id': 38,
  },
  {
    'localized_name': 'Queen of Pain',
    'id': 39,
  },
  {
    'localized_name': 'Venomancer',
    'id': 40,
  },
  {
    'localized_name': 'Faceless Void',
    'id': 41,
  },
  {
    'localized_name': 'Wraith King',
    'id': 42,
  },
  {
    'localized_name': 'Death Prophet',
    'id': 43,
  },
  {
    'localized_name': 'Phantom Assassin',
    'id': 44,
  },
  {
    'localized_name': 'Pugna',
    'id': 45,
  },
  {
    'localized_name': 'Templar Assassin',
    'id': 46,
  },
  {
    'localized_name': 'Viper',
    'id': 47,
  },
  {
    'localized_name': 'Luna',
    'id': 48,
  },
  {
    'localized_name': 'Dragon Knight',
    'id': 49,
  },
  {
    'localized_name': 'Dazzle',
    'id': 50,
  },
  {
    'localized_name': 'Clockwerk',
    'id': 51,
  },
  {
    'localized_name': 'Leshrac',
    'id': 52,
  },
  {
    'localized_name': 'Nature\'s Prophet',
    'id': 53,
  },
  {
    'localized_name': 'Lifestealer',
    'id': 54,
  },
  {
    'localized_name': 'Dark Seer',
    'id': 55,
  },
  {
    'localized_name': 'Clinkz',
    'id': 56,
  },
  {
    'localized_name': 'Omniknight',
    'id': 57,
  },
  {
    'localized_name': 'Enchantress',
    'id': 58,
  },
  {
    'localized_name': 'Huskar',
    'id': 59,
  },
  {
    'localized_name': 'Night Stalker',
    'id': 60,
  },
  {
    'localized_name': 'Broodmother',
    'id': 61,
  },
  {
    'localized_name': 'Bounty Hunter',
    'id': 62,
  },
  {
    'localized_name': 'Weaver',
    'id': 63,
  },
  {
    'localized_name': 'Jakiro',
    'id': 64,
  },
  {
    'localized_name': 'Batrider',
    'id': 65,
  },
  {
    'localized_name': 'Chen',
    'id': 66,
  },
  {
    'localized_name': 'Spectre',
    'id': 67,
  },
  {
    'localized_name': 'Ancient Apparition',
    'id': 68,
  },
  {
    'localized_name': 'Doom',
    'id': 69,
  },
  {
    'localized_name': 'Ursa',
    'id': 70,
  },
  {
    'localized_name': 'Spirit Breaker',
    'id': 71,
  },
  {
    'localized_name': 'Gyrocopter',
    'id': 72,
  },
  {
    'localized_name': 'Alchemist',
    'id': 73,
  },
  {
    'localized_name': 'Invoker',
    'id': 74,
  },
  {
    'localized_name': 'Silencer',
    'id': 75,
  },
  {
    'localized_name': 'Outworld Devourer',
    'id': 76,
  },
  {
    'localized_name': 'Lycan',
    'id': 77,
  },
  {
    'localized_name': 'Brewmaster',
    'id': 78,
  },
  {
    'localized_name': 'Shadow Demon',
    'id': 79,
  },
  {
    'localized_name': 'Lone Druid',
    'id': 80,
  },
  {
    'localized_name': 'Chaos Knight',
    'id': 81,
  },
  {
    'localized_name': 'Meepo',
    'id': 82,
  },
  {
    'localized_name': 'Treant Protector',
    'id': 83,
  },
  {
    'localized_name': 'Ogre Magi',
    'id': 84,
  },
  {
    'localized_name': 'Undying',
    'id': 85,
  },
  {
    'localized_name': 'Rubick',
    'id': 86,
  },
  {
    'localized_name': 'Disruptor',
    'id': 87,
  },
  {
    'localized_name': 'Nyx Assassin',
    'id': 88,
  },
  {
    'localized_name': 'Naga Siren',
    'id': 89,
  },
  {
    'localized_name': 'Keeper of the Light',
    'id': 90,
  },
  {
    'localized_name': 'Io',
    'id': 91,
  },
  {
    'localized_name': 'Visage',
    'id': 92,
  },
  {
    'localized_name': 'Slark',
    'id': 93,
  },
  {
    'localized_name': 'Medusa',
    'id': 94,
  },
  {
    'localized_name': 'Troll Warlord',
    'id': 95,
  },
  {
    'localized_name': 'Centaur Warrunner',
    'id': 96,
  },
  {
    'localized_name': 'Magnus',
    'id': 97,
  },
  {
    'localized_name': 'Timbersaw',
    'id': 98,
  },
  {
    'localized_name': 'Bristleback',
    'id': 99,
  },
  {
    'localized_name': 'Tusk',
    'id': 100,
  },
  {
    'localized_name': 'Skywrath Mage',
    'id': 101,
  },
  {
    'localized_name': 'Abaddon',
    'id': 102,
  },
  {
    'localized_name': 'Elder Titan',
    'id': 103,
  },
  {
    'localized_name': 'Legion Commander',
    'id': 104,
  },
  {
    'localized_name': 'Techies',
    'id': 105,
  },
  {
    'localized_name': 'Ember Spirit',
    'id': 106,
  },
  {
    'localized_name': 'Earth Spirit',
    'id': 107,
  },
  {
    'localized_name': 'Underlord',
    'id': 108,
  },
  {
    'localized_name': 'Terrorblade',
    'id': 109,
  },
  {
    'localized_name': 'Phoenix',
    'id': 110,
  },
  {
    'localized_name': 'Oracle',
    'id': 111,
  },
  {
    'localized_name': 'Winter Wyvern',
    'id': 112,
  },
  {
    'localized_name': 'Arc Warden',
    'id': 113,
  },
  {
    'localized_name': 'Monkey King',
    'id': 114,
  },
  {
    'localized_name': 'Dark Willow',
    'id': 119,
  },
  {
    'localized_name': 'Pangolier',
    'id': 120,
  },
  {
    'localized_name': 'Grimstroke',
    'id': 121,
  },
  {
    'localized_name': 'Hoodwink',
    'id': 123,
  },
  {
    'localized_name': 'Mars',
    'id': 129,
  },
  {
    'localized_name': 'Dawnbreaker',
    'id': 135,
  },
  {
    'localized_name': 'Marci',
    'id': 138,
  },
];