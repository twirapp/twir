package handlers

import (
	"strconv"
	"strings"
)

// raw: "25:4-8,10-14/22639:16-23"
// messageText: йцу Kappa Kappa
func (c *Handlers) ParseEmotes(messageText, raw string) []MessageEmote {
	var emotes []MessageEmote
	if raw == "" {
		return emotes
	}

	emotesRaw := strings.Split(raw, "/")
	for _, emoteRaw := range emotesRaw {
		emote := MessageEmote{}

		emoteRawParts := strings.Split(emoteRaw, ":")
		emote.ID = emoteRawParts[0]

		emotePositions := emoteRawParts[1]
		emote.Count = strings.Count(emotePositions, "-")

		emotePositionsRaw := strings.Split(emotePositions, ",")

		for _, emotePositionRaw := range emotePositionsRaw {
			emotePositionRawParts := strings.Split(emotePositionRaw, "-")
			start, _ := strconv.Atoi(emotePositionRawParts[0])
			end, _ := strconv.Atoi(emotePositionRawParts[1])
			emote.Positions = append(
				emote.Positions,
				EmotePosition{
					Start: start,
					End:   end,
				},
			)

			emote.Name = messageText[start : end+1]
		}

		emotes = append(emotes, emote)
	}

	return emotes
}
