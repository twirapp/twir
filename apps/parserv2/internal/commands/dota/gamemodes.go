package dota

type GameMode struct {
	ID   int
	Name string
}

var DotaGameModes = [...]GameMode{
	{
		ID:   1,
		Name: "All Pick",
	},
	{
		ID:   2,
		Name: "Captains Mode",
	},
	{
		ID:   3,
		Name: "Random Draft",
	},
	{
		ID:   4,
		Name: "Single Draft",
	},
	{
		ID:   5,
		Name: "All Random",
	},
	{
		ID:   6,
		Name: "Intro",
	},
	{
		ID:   7,
		Name: "Diretide",
	},
	{
		ID:   8,
		Name: "Reverse Captains Mode",
	},
	{
		ID:   9,
		Name: "The Greeviling",
	},
	{
		ID:   10,
		Name: "Tutorial",
	},
	{
		ID:   11,
		Name: "Mid Only",
	},
	{
		ID:   12,
		Name: "Least Played",
	},
	{
		ID:   13,
		Name: "New Player Pool",
	},
	{
		ID:   14,
		Name: "Compendium Matchmaking",
	},
	{
		ID:   15,
		Name: "Custom Game",
	},
	{
		ID:   16,
		Name: "Captains Draft",
	},
	{
		ID:   18,
		Name: "Ability Draft",
	},
	{
		ID:   19,
		Name: "Event Game",
	},
	{
		ID:   20,
		Name: "All Random Deathmatch",
	},
	{
		ID:   21,
		Name: "1v1 Mid Only",
	},
	{
		ID:   22,
		Name: "Ranked",
	},
	{
		ID:   23,
		Name: "Turbo",
	},
	{
		ID:   24,
		Name: "Mutation",
	},
}
