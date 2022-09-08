package dota

import "github.com/samber/lo"

type Hero struct {
	LocalizedName string
	ID            int
	ShortName     *string
}

var DotaHeroes = []Hero{
	{
		ID:            0,
		LocalizedName: "Not Picked",
		ShortName:     nil,
	},
	{
		ID:            1,
		LocalizedName: "Anti-Mage",
		ShortName:     nil,
	},
	{
		ID:            2,
		LocalizedName: "Axe",
		ShortName:     nil,
	},
	{
		ID:            3,
		LocalizedName: "Bane",
		ShortName:     nil,
	},
	{
		ID:            4,
		LocalizedName: "Bloodseeker",
		ShortName:     nil,
	},
	{
		ID:            5,
		LocalizedName: "Crystal Maiden",
		ShortName:     lo.ToPtr("CM"),
	},
	{
		ID:            6,
		LocalizedName: "Drow Ranger",
		ShortName:     nil,
	},
	{
		ID:            7,
		LocalizedName: "Earthshaker",
		ShortName:     nil,
	},
	{
		ID:            8,
		LocalizedName: "Juggernaut",
		ShortName:     nil,
	},
	{
		ID:            9,
		LocalizedName: "Mirana",
		ShortName:     nil,
	},
	{
		ID:            10,
		LocalizedName: "Morphling",
		ShortName:     lo.ToPtr("Morph"),
	},
	{
		ID:            11,
		LocalizedName: "Shadow Fiend",
		ShortName:     lo.ToPtr("SF"),
	},
	{
		ID:            12,
		LocalizedName: "Phantom Lancer",
		ShortName:     lo.ToPtr("PL"),
	},
	{
		ID:            13,
		LocalizedName: "Puck",
		ShortName:     nil,
	},
	{
		ID:            14,
		LocalizedName: "Pudge",
		ShortName:     nil,
	},
	{
		ID:            15,
		LocalizedName: "Razor",
		ShortName:     nil,
	},
	{
		ID:            16,
		LocalizedName: "Sand King",
		ShortName:     lo.ToPtr("SK"),
	},
	{
		ID:            17,
		LocalizedName: "Storm Spirit",
		ShortName:     lo.ToPtr("Storm"),
	},
	{
		ID:            18,
		LocalizedName: "Sven",
		ShortName:     nil,
	},
	{
		ID:            19,
		LocalizedName: "Tiny",
		ShortName:     nil,
	},
	{
		ID:            20,
		LocalizedName: "Vengeful Spirit",
		ShortName:     nil,
	},
	{
		ID:            21,
		LocalizedName: "Windranger",
		ShortName:     lo.ToPtr("WR"),
	},
	{
		ID:            22,
		LocalizedName: "Zeus",
		ShortName:     nil,
	},
	{
		ID:            23,
		LocalizedName: "Kunkka",
		ShortName:     nil,
	},
	{
		ID:            25,
		LocalizedName: "Lina",
		ShortName:     nil,
	},
	{
		ID:            26,
		LocalizedName: "Lion",
		ShortName:     nil,
	},
	{
		ID:            27,
		LocalizedName: "Shadow Shaman",
		ShortName:     lo.ToPtr("Shaman"),
	},
	{
		ID:            28,
		LocalizedName: "Slardar",
		ShortName:     nil,
	},
	{
		ID:            29,
		LocalizedName: "Tidehunter",
		ShortName:     nil,
	},
	{
		ID:            30,
		LocalizedName: "Witch Doctor",
		ShortName:     nil,
	},
	{
		ID:            31,
		LocalizedName: "Lich",
		ShortName:     nil,
	},
	{
		ID:            32,
		LocalizedName: "Riki",
		ShortName:     nil,
	},
	{
		ID:            33,
		LocalizedName: "Enigma",
		ShortName:     nil,
	},
	{
		ID:            34,
		LocalizedName: "Tinker",
		ShortName:     nil,
	},
	{
		ID:            35,
		LocalizedName: "Sniper",
		ShortName:     nil,
	},
	{
		ID:            36,
		LocalizedName: "Necrophos",
		ShortName:     lo.ToPtr("Necr"),
	},
	{
		ID:            37,
		LocalizedName: "Warlock",
		ShortName:     nil,
	},
	{
		ID:            38,
		LocalizedName: "Beastmaster",
		ShortName:     nil,
	},
	{
		ID:            39,
		LocalizedName: "Queen of Pain",
		ShortName:     lo.ToPtr("QoP"),
	},
	{
		ID:            40,
		LocalizedName: "Venomancer",
		ShortName:     nil,
	},
	{
		ID:            41,
		LocalizedName: "Faceless Void",
		ShortName:     nil,
	},
	{
		ID:            42,
		LocalizedName: "Wraith King",
		ShortName:     nil,
	},
	{
		ID:            43,
		LocalizedName: "Death Prophet",
		ShortName:     nil,
	},
	{
		ID:            44,
		LocalizedName: "Phantom Assassin",
		ShortName:     lo.ToPtr("PA"),
	},
	{
		ID:            45,
		LocalizedName: "Pugna",
		ShortName:     nil,
	},
	{
		ID:            46,
		LocalizedName: "Templar Assassin",
		ShortName:     lo.ToPtr("TA"),
	},
	{
		ID:            47,
		LocalizedName: "Viper",
		ShortName:     nil,
	},
	{
		ID:            48,
		LocalizedName: "Luna",
		ShortName:     nil,
	},
	{
		ID:            49,
		LocalizedName: "Dragon Knight",
		ShortName:     lo.ToPtr("DK"),
	},
	{
		ID:            50,
		LocalizedName: "Dazzle",
		ShortName:     nil,
	},
	{
		ID:            51,
		LocalizedName: "Clockwerk",
		ShortName:     nil,
	},
	{
		ID:            52,
		LocalizedName: "Leshrac",
		ShortName:     nil,
	},
	{
		ID:            53,
		LocalizedName: "Nature's Prophet",
		ShortName:     lo.ToPtr("Furion"),
	},
	{
		ID:            54,
		LocalizedName: "Lifestealer",
		ShortName:     nil,
	},
	{
		ID:            55,
		LocalizedName: "Dark Seer",
		ShortName:     nil,
	},
	{
		ID:            56,
		LocalizedName: "Clinkz",
		ShortName:     nil,
	},
	{
		ID:            57,
		LocalizedName: "Omniknight",
		ShortName:     nil,
	},
	{
		ID:            58,
		LocalizedName: "Enchantress",
		ShortName:     lo.ToPtr("Encha"),
	},
	{
		ID:            59,
		LocalizedName: "Huskar",
		ShortName:     nil,
	},
	{
		ID:            60,
		LocalizedName: "Night Stalker",
		ShortName:     lo.ToPtr("NS"),
	},
	{
		ID:            61,
		LocalizedName: "Broodmother",
		ShortName:     lo.ToPtr("Brood"),
	},
	{
		ID:            62,
		LocalizedName: "Bounty Hunter",
		ShortName:     lo.ToPtr("BH"),
	},
	{
		ID:            63,
		LocalizedName: "Weaver",
		ShortName:     nil,
	},
	{
		ID:            64,
		LocalizedName: "Jakiro",
		ShortName:     nil,
	},
	{
		ID:            65,
		LocalizedName: "Batrider",
		ShortName:     nil,
	},
	{
		ID:            66,
		LocalizedName: "Chen",
		ShortName:     nil,
	},
	{
		ID:            67,
		LocalizedName: "Spectre",
		ShortName:     nil,
	},
	{
		ID:            68,
		LocalizedName: "Ancient Apparition",
		ShortName:     lo.ToPtr("AA"),
	},
	{
		ID:            69,
		LocalizedName: "Doom",
		ShortName:     nil,
	},
	{
		ID:            70,
		LocalizedName: "Ursa",
		ShortName:     nil,
	},
	{
		ID:            71,
		LocalizedName: "Spirit Breaker",
		ShortName:     nil,
	},
	{
		ID:            72,
		LocalizedName: "Gyrocopter",
		ShortName:     lo.ToPtr("Gyro"),
	},
	{
		ID:            73,
		LocalizedName: "Alchemist",
		ShortName:     lo.ToPtr("Alch"),
	},
	{
		ID:            74,
		LocalizedName: "Invoker",
		ShortName:     nil,
	},
	{
		ID:            75,
		LocalizedName: "Silencer",
		ShortName:     nil,
	},
	{
		ID:            76,
		LocalizedName: "Outworld Devourer",
		ShortName:     lo.ToPtr("OD"),
	},
	{
		ID:            77,
		LocalizedName: "Lycan",
		ShortName:     nil,
	},
	{
		ID:            78,
		LocalizedName: "Brewmaster",
		ShortName:     nil,
	},
	{
		ID:            79,
		LocalizedName: "Shadow Demon",
		ShortName:     nil,
	},
	{
		ID:            80,
		LocalizedName: "Lone Druid",
		ShortName:     nil,
	},
	{
		ID:            81,
		LocalizedName: "Chaos Knight",
		ShortName:     lo.ToPtr("CK"),
	},
	{
		ID:            82,
		LocalizedName: "Meepo",
		ShortName:     nil,
	},
	{
		ID:            83,
		LocalizedName: "Treant Protector",
		ShortName:     lo.ToPtr("Treant"),
	},
	{
		ID:            84,
		LocalizedName: "Ogre Magi",
		ShortName:     lo.ToPtr("Ogre"),
	},
	{
		ID:            85,
		LocalizedName: "Undying",
		ShortName:     nil,
	},
	{
		ID:            86,
		LocalizedName: "Rubick",
		ShortName:     nil,
	},
	{
		ID:            87,
		LocalizedName: "Disruptor",
		ShortName:     nil,
	},
	{
		ID:            88,
		LocalizedName: "Nyx Assassin",
		ShortName:     nil,
	},
	{
		ID:            89,
		LocalizedName: "Naga Siren",
		ShortName:     lo.ToPtr("Naga"),
	},
	{
		ID:            90,
		LocalizedName: "Keeper of the Light",
		ShortName:     lo.ToPtr("KOTL"),
	},
	{
		ID:            91,
		LocalizedName: "Io",
		ShortName:     nil,
	},
	{
		ID:            92,
		LocalizedName: "Visage",
		ShortName:     nil,
	},
	{
		ID:            93,
		LocalizedName: "Slark",
		ShortName:     nil,
	},
	{
		ID:            94,
		LocalizedName: "Medusa",
		ShortName:     nil,
	},
	{
		ID:            95,
		LocalizedName: "Troll Warlord",
		ShortName:     nil,
	},
	{
		ID:            96,
		LocalizedName: "Centaur Warrunner",
		ShortName:     lo.ToPtr("Centaur"),
	},
	{
		ID:            97,
		LocalizedName: "Magnus",
		ShortName:     nil,
	},
	{
		ID:            98,
		LocalizedName: "Timbersaw",
		ShortName:     nil,
	},
	{
		ID:            99,
		LocalizedName: "Bristleback",
		ShortName:     nil,
	},
	{
		ID:            100,
		LocalizedName: "Tusk",
		ShortName:     nil,
	},
	{
		ID:            101,
		LocalizedName: "Skywrath Mage",
		ShortName:     lo.ToPtr("Skywrath"),
	},
	{
		ID:            102,
		LocalizedName: "Abaddon",
		ShortName:     nil,
	},
	{
		ID:            103,
		LocalizedName: "Elder Titan",
		ShortName:     lo.ToPtr("Elder"),
	},
	{
		ID:            104,
		LocalizedName: "Legion Commander",
		ShortName:     nil,
	},
	{
		ID:            105,
		LocalizedName: "Techies",
		ShortName:     nil,
	},
	{
		ID:            106,
		LocalizedName: "Ember Spirit",
		ShortName:     lo.ToPtr("Ember"),
	},
	{
		ID:            107,
		LocalizedName: "Earth Spirit",
		ShortName:     nil,
	},
	{
		ID:            108,
		LocalizedName: "Underlord",
		ShortName:     nil,
	},
	{
		ID:            109,
		LocalizedName: "Terrorblade",
		ShortName:     nil,
	},
	{
		ID:            110,
		LocalizedName: "Phoenix",
		ShortName:     nil,
	},
	{
		ID:            111,
		LocalizedName: "Oracle",
		ShortName:     nil,
	},
	{
		ID:            112,
		LocalizedName: "Winter Wyvern",
		ShortName:     lo.ToPtr("Wyvern"),
	},
	{
		ID:            113,
		LocalizedName: "Arc Warden",
		ShortName:     lo.ToPtr("Arc"),
	},
	{
		ID:            114,
		LocalizedName: "Monkey King",
		ShortName:     nil,
	},
	{
		ID:            119,
		LocalizedName: "Dark Willow",
		ShortName:     nil,
	},
	{
		ID:            120,
		LocalizedName: "Pangolier",
		ShortName:     nil,
	},
	{
		ID:            121,
		LocalizedName: "Grimstroke",
		ShortName:     nil,
	},
	{
		ID:            123,
		LocalizedName: "Hoodwink",
		ShortName:     nil,
	},
	{
		ID:            129,
		LocalizedName: "Mars",
		ShortName:     nil,
	},
	{
		ID:            135,
		LocalizedName: "Dawnbreaker",
		ShortName:     nil,
	},
	{
		ID:            138,
		LocalizedName: "Marci",
		ShortName:     nil,
	},
}
