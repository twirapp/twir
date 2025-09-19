package locales

import twiri18n "github.com/twirapp/twir/libs/i18n"

type KeysCommandsChannelGameSuccessGameSetVars struct {
	Game any
}
type KeysCommandsChannelGameSuccessGameSet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelGameSuccessGameSet) IsTranslationKey() {
}
func (k KeysCommandsChannelGameSuccessGameSet) GetPath() string {
	return "commands.channel.game.success.game_set"
}
func (k KeysCommandsChannelGameSuccessGameSet) GetPathSlice() []string {
	return []string{"commands", "channel", "game", "success", "game_set"}
}
func (k KeysCommandsChannelGameSuccessGameSet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelGameSuccessGameSet) SetVars(vars KeysCommandsChannelGameSuccessGameSetVars) twiri18n.TranslationKey[KeysCommandsChannelGameSuccessGameSetVars] {
	k.Vars = twiri18n.Vars{"game": vars.Game}
	return k
}

type KeysCommandsChannelGameSuccess struct {
	GameSet KeysCommandsChannelGameSuccessGameSet
}
type KeysCommandsChannelGameDescriptionVars struct {
}
type KeysCommandsChannelGameDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelGameDescription) IsTranslationKey() {
}
func (k KeysCommandsChannelGameDescription) GetPath() string {
	return "commands.channel.game.description"
}
func (k KeysCommandsChannelGameDescription) GetPathSlice() []string {
	return []string{"commands", "channel", "game", "description"}
}
func (k KeysCommandsChannelGameDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelGameDescription) SetVars(vars KeysCommandsChannelGameDescriptionVars) twiri18n.TranslationKey[KeysCommandsChannelGameDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelGameErrorsCannotSetGameVars struct {
}
type KeysCommandsChannelGameErrorsCannotSetGame struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelGameErrorsCannotSetGame) IsTranslationKey() {
}
func (k KeysCommandsChannelGameErrorsCannotSetGame) GetPath() string {
	return "commands.channel.game.errors.cannot_set_game"
}
func (k KeysCommandsChannelGameErrorsCannotSetGame) GetPathSlice() []string {
	return []string{"commands", "channel", "game", "errors", "cannot_set_game"}
}
func (k KeysCommandsChannelGameErrorsCannotSetGame) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelGameErrorsCannotSetGame) SetVars(vars KeysCommandsChannelGameErrorsCannotSetGameVars) twiri18n.TranslationKey[KeysCommandsChannelGameErrorsCannotSetGameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelGameErrorsGameNotFoundVars struct {
}
type KeysCommandsChannelGameErrorsGameNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelGameErrorsGameNotFound) IsTranslationKey() {
}
func (k KeysCommandsChannelGameErrorsGameNotFound) GetPath() string {
	return "commands.channel.game.errors.game_not_found"
}
func (k KeysCommandsChannelGameErrorsGameNotFound) GetPathSlice() []string {
	return []string{"commands", "channel", "game", "errors", "game_not_found"}
}
func (k KeysCommandsChannelGameErrorsGameNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelGameErrorsGameNotFound) SetVars(vars KeysCommandsChannelGameErrorsGameNotFoundVars) twiri18n.TranslationKey[KeysCommandsChannelGameErrorsGameNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelGameErrors struct {
	CannotSetGame	KeysCommandsChannelGameErrorsCannotSetGame
	GameNotFound	KeysCommandsChannelGameErrorsGameNotFound
}
type KeysCommandsChannelGame struct {
	Success		KeysCommandsChannelGameSuccess
	Description	KeysCommandsChannelGameDescription
	Errors		KeysCommandsChannelGameErrors
}
type KeysCommandsChannelDescriptionVars struct {
}
type KeysCommandsChannelDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelDescription) IsTranslationKey() {
}
func (k KeysCommandsChannelDescription) GetPath() string {
	return "commands.channel.description"
}
func (k KeysCommandsChannelDescription) GetPathSlice() []string {
	return []string{"commands", "channel", "description"}
}
func (k KeysCommandsChannelDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelDescription) SetVars(vars KeysCommandsChannelDescriptionVars) twiri18n.TranslationKey[KeysCommandsChannelDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelTitleDescriptionVars struct {
}
type KeysCommandsChannelTitleDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelTitleDescription) IsTranslationKey() {
}
func (k KeysCommandsChannelTitleDescription) GetPath() string {
	return "commands.channel.title.description"
}
func (k KeysCommandsChannelTitleDescription) GetPathSlice() []string {
	return []string{"commands", "channel", "title", "description"}
}
func (k KeysCommandsChannelTitleDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelTitleDescription) SetVars(vars KeysCommandsChannelTitleDescriptionVars) twiri18n.TranslationKey[KeysCommandsChannelTitleDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelTitleErrorsCannotSetTitleVars struct {
}
type KeysCommandsChannelTitleErrorsCannotSetTitle struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelTitleErrorsCannotSetTitle) IsTranslationKey() {
}
func (k KeysCommandsChannelTitleErrorsCannotSetTitle) GetPath() string {
	return "commands.channel.title.errors.cannot_set_title"
}
func (k KeysCommandsChannelTitleErrorsCannotSetTitle) GetPathSlice() []string {
	return []string{"commands", "channel", "title", "errors", "cannot_set_title"}
}
func (k KeysCommandsChannelTitleErrorsCannotSetTitle) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelTitleErrorsCannotSetTitle) SetVars(vars KeysCommandsChannelTitleErrorsCannotSetTitleVars) twiri18n.TranslationKey[KeysCommandsChannelTitleErrorsCannotSetTitleVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelTitleErrorsTitleTooLongVars struct {
}
type KeysCommandsChannelTitleErrorsTitleTooLong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelTitleErrorsTitleTooLong) IsTranslationKey() {
}
func (k KeysCommandsChannelTitleErrorsTitleTooLong) GetPath() string {
	return "commands.channel.title.errors.title_too_long"
}
func (k KeysCommandsChannelTitleErrorsTitleTooLong) GetPathSlice() []string {
	return []string{"commands", "channel", "title", "errors", "title_too_long"}
}
func (k KeysCommandsChannelTitleErrorsTitleTooLong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelTitleErrorsTitleTooLong) SetVars(vars KeysCommandsChannelTitleErrorsTitleTooLongVars) twiri18n.TranslationKey[KeysCommandsChannelTitleErrorsTitleTooLongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelTitleErrors struct {
	CannotSetTitle	KeysCommandsChannelTitleErrorsCannotSetTitle
	TitleTooLong	KeysCommandsChannelTitleErrorsTitleTooLong
}
type KeysCommandsChannelTitleSuccessTitleSetVars struct {
	Title any
}
type KeysCommandsChannelTitleSuccessTitleSet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelTitleSuccessTitleSet) IsTranslationKey() {
}
func (k KeysCommandsChannelTitleSuccessTitleSet) GetPath() string {
	return "commands.channel.title.success.title_set"
}
func (k KeysCommandsChannelTitleSuccessTitleSet) GetPathSlice() []string {
	return []string{"commands", "channel", "title", "success", "title_set"}
}
func (k KeysCommandsChannelTitleSuccessTitleSet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelTitleSuccessTitleSet) SetVars(vars KeysCommandsChannelTitleSuccessTitleSetVars) twiri18n.TranslationKey[KeysCommandsChannelTitleSuccessTitleSetVars] {
	k.Vars = twiri18n.Vars{"title": vars.Title}
	return k
}

type KeysCommandsChannelTitleSuccess struct {
	TitleSet KeysCommandsChannelTitleSuccessTitleSet
}
type KeysCommandsChannelTitle struct {
	Description	KeysCommandsChannelTitleDescription
	Errors		KeysCommandsChannelTitleErrors
	Success		KeysCommandsChannelTitleSuccess
}
type KeysCommandsChannel struct {
	Game		KeysCommandsChannelGame
	Description	KeysCommandsChannelDescription
	Title		KeysCommandsChannelTitle
}
type KeysCommandsClipDescriptionVars struct {
}
type KeysCommandsClipDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsClipDescription) IsTranslationKey() {
}
func (k KeysCommandsClipDescription) GetPath() string {
	return "commands.clip.description"
}
func (k KeysCommandsClipDescription) GetPathSlice() []string {
	return []string{"commands", "clip", "description"}
}
func (k KeysCommandsClipDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsClipDescription) SetVars(vars KeysCommandsClipDescriptionVars) twiri18n.TranslationKey[KeysCommandsClipDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsClipErrorsCannotFindChannelVars struct {
}
type KeysCommandsClipErrorsCannotFindChannel struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsClipErrorsCannotFindChannel) IsTranslationKey() {
}
func (k KeysCommandsClipErrorsCannotFindChannel) GetPath() string {
	return "commands.clip.errors.cannot_find_channel"
}
func (k KeysCommandsClipErrorsCannotFindChannel) GetPathSlice() []string {
	return []string{"commands", "clip", "errors", "cannot_find_channel"}
}
func (k KeysCommandsClipErrorsCannotFindChannel) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsClipErrorsCannotFindChannel) SetVars(vars KeysCommandsClipErrorsCannotFindChannelVars) twiri18n.TranslationKey[KeysCommandsClipErrorsCannotFindChannelVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsClipErrorsCannotCreateBroadcasterClientVars struct {
}
type KeysCommandsClipErrorsCannotCreateBroadcasterClient struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsClipErrorsCannotCreateBroadcasterClient) IsTranslationKey() {
}
func (k KeysCommandsClipErrorsCannotCreateBroadcasterClient) GetPath() string {
	return "commands.clip.errors.cannot_create_broadcaster_client"
}
func (k KeysCommandsClipErrorsCannotCreateBroadcasterClient) GetPathSlice() []string {
	return []string{"commands", "clip", "errors", "cannot_create_broadcaster_client"}
}
func (k KeysCommandsClipErrorsCannotCreateBroadcasterClient) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsClipErrorsCannotCreateBroadcasterClient) SetVars(vars KeysCommandsClipErrorsCannotCreateBroadcasterClientVars) twiri18n.TranslationKey[KeysCommandsClipErrorsCannotCreateBroadcasterClientVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsClipErrorsCannotCreateClipVars struct {
}
type KeysCommandsClipErrorsCannotCreateClip struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsClipErrorsCannotCreateClip) IsTranslationKey() {
}
func (k KeysCommandsClipErrorsCannotCreateClip) GetPath() string {
	return "commands.clip.errors.cannot_create_clip"
}
func (k KeysCommandsClipErrorsCannotCreateClip) GetPathSlice() []string {
	return []string{"commands", "clip", "errors", "cannot_create_clip"}
}
func (k KeysCommandsClipErrorsCannotCreateClip) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsClipErrorsCannotCreateClip) SetVars(vars KeysCommandsClipErrorsCannotCreateClipVars) twiri18n.TranslationKey[KeysCommandsClipErrorsCannotCreateClipVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsClipErrorsCannotGetClipVars struct {
}
type KeysCommandsClipErrorsCannotGetClip struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsClipErrorsCannotGetClip) IsTranslationKey() {
}
func (k KeysCommandsClipErrorsCannotGetClip) GetPath() string {
	return "commands.clip.errors.cannot_get_clip"
}
func (k KeysCommandsClipErrorsCannotGetClip) GetPathSlice() []string {
	return []string{"commands", "clip", "errors", "cannot_get_clip"}
}
func (k KeysCommandsClipErrorsCannotGetClip) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsClipErrorsCannotGetClip) SetVars(vars KeysCommandsClipErrorsCannotGetClipVars) twiri18n.TranslationKey[KeysCommandsClipErrorsCannotGetClipVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsClipErrorsEmptyClipUrlVars struct {
}
type KeysCommandsClipErrorsEmptyClipUrl struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsClipErrorsEmptyClipUrl) IsTranslationKey() {
}
func (k KeysCommandsClipErrorsEmptyClipUrl) GetPath() string {
	return "commands.clip.errors.empty_clip_url"
}
func (k KeysCommandsClipErrorsEmptyClipUrl) GetPathSlice() []string {
	return []string{"commands", "clip", "errors", "empty_clip_url"}
}
func (k KeysCommandsClipErrorsEmptyClipUrl) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsClipErrorsEmptyClipUrl) SetVars(vars KeysCommandsClipErrorsEmptyClipUrlVars) twiri18n.TranslationKey[KeysCommandsClipErrorsEmptyClipUrlVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsClipErrors struct {
	CannotFindChannel		KeysCommandsClipErrorsCannotFindChannel
	CannotCreateBroadcasterClient	KeysCommandsClipErrorsCannotCreateBroadcasterClient
	CannotCreateClip		KeysCommandsClipErrorsCannotCreateClip
	CannotGetClip			KeysCommandsClipErrorsCannotGetClip
	EmptyClipUrl			KeysCommandsClipErrorsEmptyClipUrl
}
type KeysCommandsClipSuccessClipCreatedVars struct {
	Url any
}
type KeysCommandsClipSuccessClipCreated struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsClipSuccessClipCreated) IsTranslationKey() {
}
func (k KeysCommandsClipSuccessClipCreated) GetPath() string {
	return "commands.clip.success.clip_created"
}
func (k KeysCommandsClipSuccessClipCreated) GetPathSlice() []string {
	return []string{"commands", "clip", "success", "clip_created"}
}
func (k KeysCommandsClipSuccessClipCreated) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsClipSuccessClipCreated) SetVars(vars KeysCommandsClipSuccessClipCreatedVars) twiri18n.TranslationKey[KeysCommandsClipSuccessClipCreatedVars] {
	k.Vars = twiri18n.Vars{"url": vars.Url}
	return k
}

type KeysCommandsClipSuccess struct {
	ClipCreated KeysCommandsClipSuccessClipCreated
}
type KeysCommandsClip struct {
	Description	KeysCommandsClipDescription
	Errors		KeysCommandsClipErrors
	Success		KeysCommandsClipSuccess
}
type KeysCommandsDotaDescriptionVars struct {
}
type KeysCommandsDotaDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDotaDescription) IsTranslationKey() {
}
func (k KeysCommandsDotaDescription) GetPath() string {
	return "commands.dota.description"
}
func (k KeysCommandsDotaDescription) GetPathSlice() []string {
	return []string{"commands", "dota", "description"}
}
func (k KeysCommandsDotaDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDotaDescription) SetVars(vars KeysCommandsDotaDescriptionVars) twiri18n.TranslationKey[KeysCommandsDotaDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDotaStatsDescriptionVars struct {
}
type KeysCommandsDotaStatsDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDotaStatsDescription) IsTranslationKey() {
}
func (k KeysCommandsDotaStatsDescription) GetPath() string {
	return "commands.dota.stats.description"
}
func (k KeysCommandsDotaStatsDescription) GetPathSlice() []string {
	return []string{"commands", "dota", "stats", "description"}
}
func (k KeysCommandsDotaStatsDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDotaStatsDescription) SetVars(vars KeysCommandsDotaStatsDescriptionVars) twiri18n.TranslationKey[KeysCommandsDotaStatsDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDotaStatsErrorsCannotGetStatsVars struct {
}
type KeysCommandsDotaStatsErrorsCannotGetStats struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDotaStatsErrorsCannotGetStats) IsTranslationKey() {
}
func (k KeysCommandsDotaStatsErrorsCannotGetStats) GetPath() string {
	return "commands.dota.stats.errors.cannot_get_stats"
}
func (k KeysCommandsDotaStatsErrorsCannotGetStats) GetPathSlice() []string {
	return []string{"commands", "dota", "stats", "errors", "cannot_get_stats"}
}
func (k KeysCommandsDotaStatsErrorsCannotGetStats) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDotaStatsErrorsCannotGetStats) SetVars(vars KeysCommandsDotaStatsErrorsCannotGetStatsVars) twiri18n.TranslationKey[KeysCommandsDotaStatsErrorsCannotGetStatsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDotaStatsErrorsPlayerNotFoundVars struct {
}
type KeysCommandsDotaStatsErrorsPlayerNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDotaStatsErrorsPlayerNotFound) IsTranslationKey() {
}
func (k KeysCommandsDotaStatsErrorsPlayerNotFound) GetPath() string {
	return "commands.dota.stats.errors.player_not_found"
}
func (k KeysCommandsDotaStatsErrorsPlayerNotFound) GetPathSlice() []string {
	return []string{"commands", "dota", "stats", "errors", "player_not_found"}
}
func (k KeysCommandsDotaStatsErrorsPlayerNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDotaStatsErrorsPlayerNotFound) SetVars(vars KeysCommandsDotaStatsErrorsPlayerNotFoundVars) twiri18n.TranslationKey[KeysCommandsDotaStatsErrorsPlayerNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDotaStatsErrors struct {
	CannotGetStats	KeysCommandsDotaStatsErrorsCannotGetStats
	PlayerNotFound	KeysCommandsDotaStatsErrorsPlayerNotFound
}
type KeysCommandsDotaStatsSuccessStatsDisplayVars struct {
	Player	any
	Stats	any
}
type KeysCommandsDotaStatsSuccessStatsDisplay struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDotaStatsSuccessStatsDisplay) IsTranslationKey() {
}
func (k KeysCommandsDotaStatsSuccessStatsDisplay) GetPath() string {
	return "commands.dota.stats.success.stats_display"
}
func (k KeysCommandsDotaStatsSuccessStatsDisplay) GetPathSlice() []string {
	return []string{"commands", "dota", "stats", "success", "stats_display"}
}
func (k KeysCommandsDotaStatsSuccessStatsDisplay) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDotaStatsSuccessStatsDisplay) SetVars(vars KeysCommandsDotaStatsSuccessStatsDisplayVars) twiri18n.TranslationKey[KeysCommandsDotaStatsSuccessStatsDisplayVars] {
	k.Vars = twiri18n.Vars{"player": vars.Player, "stats": vars.Stats}
	return k
}

type KeysCommandsDotaStatsSuccess struct {
	StatsDisplay KeysCommandsDotaStatsSuccessStatsDisplay
}
type KeysCommandsDotaStats struct {
	Description	KeysCommandsDotaStatsDescription
	Errors		KeysCommandsDotaStatsErrors
	Success		KeysCommandsDotaStatsSuccess
}
type KeysCommandsDotaMatchesDescriptionVars struct {
}
type KeysCommandsDotaMatchesDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDotaMatchesDescription) IsTranslationKey() {
}
func (k KeysCommandsDotaMatchesDescription) GetPath() string {
	return "commands.dota.matches.description"
}
func (k KeysCommandsDotaMatchesDescription) GetPathSlice() []string {
	return []string{"commands", "dota", "matches", "description"}
}
func (k KeysCommandsDotaMatchesDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDotaMatchesDescription) SetVars(vars KeysCommandsDotaMatchesDescriptionVars) twiri18n.TranslationKey[KeysCommandsDotaMatchesDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDotaMatchesErrorsCannotGetMatchesVars struct {
}
type KeysCommandsDotaMatchesErrorsCannotGetMatches struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDotaMatchesErrorsCannotGetMatches) IsTranslationKey() {
}
func (k KeysCommandsDotaMatchesErrorsCannotGetMatches) GetPath() string {
	return "commands.dota.matches.errors.cannot_get_matches"
}
func (k KeysCommandsDotaMatchesErrorsCannotGetMatches) GetPathSlice() []string {
	return []string{"commands", "dota", "matches", "errors", "cannot_get_matches"}
}
func (k KeysCommandsDotaMatchesErrorsCannotGetMatches) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDotaMatchesErrorsCannotGetMatches) SetVars(vars KeysCommandsDotaMatchesErrorsCannotGetMatchesVars) twiri18n.TranslationKey[KeysCommandsDotaMatchesErrorsCannotGetMatchesVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDotaMatchesErrors struct {
	CannotGetMatches KeysCommandsDotaMatchesErrorsCannotGetMatches
}
type KeysCommandsDotaMatchesSuccessMatchesDisplayVars struct {
	Matches any
}
type KeysCommandsDotaMatchesSuccessMatchesDisplay struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDotaMatchesSuccessMatchesDisplay) IsTranslationKey() {
}
func (k KeysCommandsDotaMatchesSuccessMatchesDisplay) GetPath() string {
	return "commands.dota.matches.success.matches_display"
}
func (k KeysCommandsDotaMatchesSuccessMatchesDisplay) GetPathSlice() []string {
	return []string{"commands", "dota", "matches", "success", "matches_display"}
}
func (k KeysCommandsDotaMatchesSuccessMatchesDisplay) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDotaMatchesSuccessMatchesDisplay) SetVars(vars KeysCommandsDotaMatchesSuccessMatchesDisplayVars) twiri18n.TranslationKey[KeysCommandsDotaMatchesSuccessMatchesDisplayVars] {
	k.Vars = twiri18n.Vars{"matches": vars.Matches}
	return k
}

type KeysCommandsDotaMatchesSuccess struct {
	MatchesDisplay KeysCommandsDotaMatchesSuccessMatchesDisplay
}
type KeysCommandsDotaMatches struct {
	Description	KeysCommandsDotaMatchesDescription
	Errors		KeysCommandsDotaMatchesErrors
	Success		KeysCommandsDotaMatchesSuccess
}
type KeysCommandsDota struct {
	Description	KeysCommandsDotaDescription
	Stats		KeysCommandsDotaStats
	Matches		KeysCommandsDotaMatches
}
type KeysCommandsDudesStatsSuccessStatsDisplayVars struct {
	Stats any
}
type KeysCommandsDudesStatsSuccessStatsDisplay struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesStatsSuccessStatsDisplay) IsTranslationKey() {
}
func (k KeysCommandsDudesStatsSuccessStatsDisplay) GetPath() string {
	return "commands.dudes.stats.success.stats_display"
}
func (k KeysCommandsDudesStatsSuccessStatsDisplay) GetPathSlice() []string {
	return []string{"commands", "dudes", "stats", "success", "stats_display"}
}
func (k KeysCommandsDudesStatsSuccessStatsDisplay) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesStatsSuccessStatsDisplay) SetVars(vars KeysCommandsDudesStatsSuccessStatsDisplayVars) twiri18n.TranslationKey[KeysCommandsDudesStatsSuccessStatsDisplayVars] {
	k.Vars = twiri18n.Vars{"stats": vars.Stats}
	return k
}

type KeysCommandsDudesStatsSuccess struct {
	StatsDisplay KeysCommandsDudesStatsSuccessStatsDisplay
}
type KeysCommandsDudesStatsDescriptionVars struct {
}
type KeysCommandsDudesStatsDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesStatsDescription) IsTranslationKey() {
}
func (k KeysCommandsDudesStatsDescription) GetPath() string {
	return "commands.dudes.stats.description"
}
func (k KeysCommandsDudesStatsDescription) GetPathSlice() []string {
	return []string{"commands", "dudes", "stats", "description"}
}
func (k KeysCommandsDudesStatsDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesStatsDescription) SetVars(vars KeysCommandsDudesStatsDescriptionVars) twiri18n.TranslationKey[KeysCommandsDudesStatsDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesStatsErrorsCannotGetStatsVars struct {
}
type KeysCommandsDudesStatsErrorsCannotGetStats struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesStatsErrorsCannotGetStats) IsTranslationKey() {
}
func (k KeysCommandsDudesStatsErrorsCannotGetStats) GetPath() string {
	return "commands.dudes.stats.errors.cannot_get_stats"
}
func (k KeysCommandsDudesStatsErrorsCannotGetStats) GetPathSlice() []string {
	return []string{"commands", "dudes", "stats", "errors", "cannot_get_stats"}
}
func (k KeysCommandsDudesStatsErrorsCannotGetStats) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesStatsErrorsCannotGetStats) SetVars(vars KeysCommandsDudesStatsErrorsCannotGetStatsVars) twiri18n.TranslationKey[KeysCommandsDudesStatsErrorsCannotGetStatsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesStatsErrors struct {
	CannotGetStats KeysCommandsDudesStatsErrorsCannotGetStats
}
type KeysCommandsDudesStats struct {
	Success		KeysCommandsDudesStatsSuccess
	Description	KeysCommandsDudesStatsDescription
	Errors		KeysCommandsDudesStatsErrors
}
type KeysCommandsDudesDescriptionVars struct {
}
type KeysCommandsDudesDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesDescription) IsTranslationKey() {
}
func (k KeysCommandsDudesDescription) GetPath() string {
	return "commands.dudes.description"
}
func (k KeysCommandsDudesDescription) GetPathSlice() []string {
	return []string{"commands", "dudes", "description"}
}
func (k KeysCommandsDudesDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesDescription) SetVars(vars KeysCommandsDudesDescriptionVars) twiri18n.TranslationKey[KeysCommandsDudesDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesPlayDescriptionVars struct {
}
type KeysCommandsDudesPlayDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesPlayDescription) IsTranslationKey() {
}
func (k KeysCommandsDudesPlayDescription) GetPath() string {
	return "commands.dudes.play.description"
}
func (k KeysCommandsDudesPlayDescription) GetPathSlice() []string {
	return []string{"commands", "dudes", "play", "description"}
}
func (k KeysCommandsDudesPlayDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesPlayDescription) SetVars(vars KeysCommandsDudesPlayDescriptionVars) twiri18n.TranslationKey[KeysCommandsDudesPlayDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesPlayErrorsCannotStartGameVars struct {
}
type KeysCommandsDudesPlayErrorsCannotStartGame struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesPlayErrorsCannotStartGame) IsTranslationKey() {
}
func (k KeysCommandsDudesPlayErrorsCannotStartGame) GetPath() string {
	return "commands.dudes.play.errors.cannot_start_game"
}
func (k KeysCommandsDudesPlayErrorsCannotStartGame) GetPathSlice() []string {
	return []string{"commands", "dudes", "play", "errors", "cannot_start_game"}
}
func (k KeysCommandsDudesPlayErrorsCannotStartGame) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesPlayErrorsCannotStartGame) SetVars(vars KeysCommandsDudesPlayErrorsCannotStartGameVars) twiri18n.TranslationKey[KeysCommandsDudesPlayErrorsCannotStartGameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesPlayErrorsGameAlreadyRunningVars struct {
}
type KeysCommandsDudesPlayErrorsGameAlreadyRunning struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesPlayErrorsGameAlreadyRunning) IsTranslationKey() {
}
func (k KeysCommandsDudesPlayErrorsGameAlreadyRunning) GetPath() string {
	return "commands.dudes.play.errors.game_already_running"
}
func (k KeysCommandsDudesPlayErrorsGameAlreadyRunning) GetPathSlice() []string {
	return []string{"commands", "dudes", "play", "errors", "game_already_running"}
}
func (k KeysCommandsDudesPlayErrorsGameAlreadyRunning) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesPlayErrorsGameAlreadyRunning) SetVars(vars KeysCommandsDudesPlayErrorsGameAlreadyRunningVars) twiri18n.TranslationKey[KeysCommandsDudesPlayErrorsGameAlreadyRunningVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesPlayErrors struct {
	CannotStartGame		KeysCommandsDudesPlayErrorsCannotStartGame
	GameAlreadyRunning	KeysCommandsDudesPlayErrorsGameAlreadyRunning
}
type KeysCommandsDudesPlaySuccessGameStartedVars struct {
}
type KeysCommandsDudesPlaySuccessGameStarted struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesPlaySuccessGameStarted) IsTranslationKey() {
}
func (k KeysCommandsDudesPlaySuccessGameStarted) GetPath() string {
	return "commands.dudes.play.success.game_started"
}
func (k KeysCommandsDudesPlaySuccessGameStarted) GetPathSlice() []string {
	return []string{"commands", "dudes", "play", "success", "game_started"}
}
func (k KeysCommandsDudesPlaySuccessGameStarted) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesPlaySuccessGameStarted) SetVars(vars KeysCommandsDudesPlaySuccessGameStartedVars) twiri18n.TranslationKey[KeysCommandsDudesPlaySuccessGameStartedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesPlaySuccess struct {
	GameStarted KeysCommandsDudesPlaySuccessGameStarted
}
type KeysCommandsDudesPlay struct {
	Description	KeysCommandsDudesPlayDescription
	Errors		KeysCommandsDudesPlayErrors
	Success		KeysCommandsDudesPlaySuccess
}
type KeysCommandsDudesStopDescriptionVars struct {
}
type KeysCommandsDudesStopDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesStopDescription) IsTranslationKey() {
}
func (k KeysCommandsDudesStopDescription) GetPath() string {
	return "commands.dudes.stop.description"
}
func (k KeysCommandsDudesStopDescription) GetPathSlice() []string {
	return []string{"commands", "dudes", "stop", "description"}
}
func (k KeysCommandsDudesStopDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesStopDescription) SetVars(vars KeysCommandsDudesStopDescriptionVars) twiri18n.TranslationKey[KeysCommandsDudesStopDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesStopErrorsCannotStopGameVars struct {
}
type KeysCommandsDudesStopErrorsCannotStopGame struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesStopErrorsCannotStopGame) IsTranslationKey() {
}
func (k KeysCommandsDudesStopErrorsCannotStopGame) GetPath() string {
	return "commands.dudes.stop.errors.cannot_stop_game"
}
func (k KeysCommandsDudesStopErrorsCannotStopGame) GetPathSlice() []string {
	return []string{"commands", "dudes", "stop", "errors", "cannot_stop_game"}
}
func (k KeysCommandsDudesStopErrorsCannotStopGame) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesStopErrorsCannotStopGame) SetVars(vars KeysCommandsDudesStopErrorsCannotStopGameVars) twiri18n.TranslationKey[KeysCommandsDudesStopErrorsCannotStopGameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesStopErrorsNoGameRunningVars struct {
}
type KeysCommandsDudesStopErrorsNoGameRunning struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesStopErrorsNoGameRunning) IsTranslationKey() {
}
func (k KeysCommandsDudesStopErrorsNoGameRunning) GetPath() string {
	return "commands.dudes.stop.errors.no_game_running"
}
func (k KeysCommandsDudesStopErrorsNoGameRunning) GetPathSlice() []string {
	return []string{"commands", "dudes", "stop", "errors", "no_game_running"}
}
func (k KeysCommandsDudesStopErrorsNoGameRunning) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesStopErrorsNoGameRunning) SetVars(vars KeysCommandsDudesStopErrorsNoGameRunningVars) twiri18n.TranslationKey[KeysCommandsDudesStopErrorsNoGameRunningVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesStopErrors struct {
	CannotStopGame	KeysCommandsDudesStopErrorsCannotStopGame
	NoGameRunning	KeysCommandsDudesStopErrorsNoGameRunning
}
type KeysCommandsDudesStopSuccessGameStoppedVars struct {
}
type KeysCommandsDudesStopSuccessGameStopped struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesStopSuccessGameStopped) IsTranslationKey() {
}
func (k KeysCommandsDudesStopSuccessGameStopped) GetPath() string {
	return "commands.dudes.stop.success.game_stopped"
}
func (k KeysCommandsDudesStopSuccessGameStopped) GetPathSlice() []string {
	return []string{"commands", "dudes", "stop", "success", "game_stopped"}
}
func (k KeysCommandsDudesStopSuccessGameStopped) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesStopSuccessGameStopped) SetVars(vars KeysCommandsDudesStopSuccessGameStoppedVars) twiri18n.TranslationKey[KeysCommandsDudesStopSuccessGameStoppedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesStopSuccess struct {
	GameStopped KeysCommandsDudesStopSuccessGameStopped
}
type KeysCommandsDudesStop struct {
	Description	KeysCommandsDudesStopDescription
	Errors		KeysCommandsDudesStopErrors
	Success		KeysCommandsDudesStopSuccess
}
type KeysCommandsDudes struct {
	Stats		KeysCommandsDudesStats
	Description	KeysCommandsDudesDescription
	Play		KeysCommandsDudesPlay
	Stop		KeysCommandsDudesStop
}
type KeysCommandsMarkerErrorsCannotFindChannelVars struct {
}
type KeysCommandsMarkerErrorsCannotFindChannel struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsMarkerErrorsCannotFindChannel) IsTranslationKey() {
}
func (k KeysCommandsMarkerErrorsCannotFindChannel) GetPath() string {
	return "commands.marker.errors.cannot_find_channel"
}
func (k KeysCommandsMarkerErrorsCannotFindChannel) GetPathSlice() []string {
	return []string{"commands", "marker", "errors", "cannot_find_channel"}
}
func (k KeysCommandsMarkerErrorsCannotFindChannel) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsMarkerErrorsCannotFindChannel) SetVars(vars KeysCommandsMarkerErrorsCannotFindChannelVars) twiri18n.TranslationKey[KeysCommandsMarkerErrorsCannotFindChannelVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsMarkerErrorsCannotCreateBroadcasterClientVars struct {
}
type KeysCommandsMarkerErrorsCannotCreateBroadcasterClient struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsMarkerErrorsCannotCreateBroadcasterClient) IsTranslationKey() {
}
func (k KeysCommandsMarkerErrorsCannotCreateBroadcasterClient) GetPath() string {
	return "commands.marker.errors.cannot_create_broadcaster_client"
}
func (k KeysCommandsMarkerErrorsCannotCreateBroadcasterClient) GetPathSlice() []string {
	return []string{"commands", "marker", "errors", "cannot_create_broadcaster_client"}
}
func (k KeysCommandsMarkerErrorsCannotCreateBroadcasterClient) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsMarkerErrorsCannotCreateBroadcasterClient) SetVars(vars KeysCommandsMarkerErrorsCannotCreateBroadcasterClientVars) twiri18n.TranslationKey[KeysCommandsMarkerErrorsCannotCreateBroadcasterClientVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsMarkerErrorsCannotCreateMarkerVars struct {
}
type KeysCommandsMarkerErrorsCannotCreateMarker struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsMarkerErrorsCannotCreateMarker) IsTranslationKey() {
}
func (k KeysCommandsMarkerErrorsCannotCreateMarker) GetPath() string {
	return "commands.marker.errors.cannot_create_marker"
}
func (k KeysCommandsMarkerErrorsCannotCreateMarker) GetPathSlice() []string {
	return []string{"commands", "marker", "errors", "cannot_create_marker"}
}
func (k KeysCommandsMarkerErrorsCannotCreateMarker) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsMarkerErrorsCannotCreateMarker) SetVars(vars KeysCommandsMarkerErrorsCannotCreateMarkerVars) twiri18n.TranslationKey[KeysCommandsMarkerErrorsCannotCreateMarkerVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsMarkerErrors struct {
	CannotFindChannel		KeysCommandsMarkerErrorsCannotFindChannel
	CannotCreateBroadcasterClient	KeysCommandsMarkerErrorsCannotCreateBroadcasterClient
	CannotCreateMarker		KeysCommandsMarkerErrorsCannotCreateMarker
}
type KeysCommandsMarkerSuccessMarkerCreatedVars struct {
}
type KeysCommandsMarkerSuccessMarkerCreated struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsMarkerSuccessMarkerCreated) IsTranslationKey() {
}
func (k KeysCommandsMarkerSuccessMarkerCreated) GetPath() string {
	return "commands.marker.success.marker_created"
}
func (k KeysCommandsMarkerSuccessMarkerCreated) GetPathSlice() []string {
	return []string{"commands", "marker", "success", "marker_created"}
}
func (k KeysCommandsMarkerSuccessMarkerCreated) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsMarkerSuccessMarkerCreated) SetVars(vars KeysCommandsMarkerSuccessMarkerCreatedVars) twiri18n.TranslationKey[KeysCommandsMarkerSuccessMarkerCreatedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsMarkerSuccess struct {
	MarkerCreated KeysCommandsMarkerSuccessMarkerCreated
}
type KeysCommandsMarkerDescriptionVars struct {
}
type KeysCommandsMarkerDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsMarkerDescription) IsTranslationKey() {
}
func (k KeysCommandsMarkerDescription) GetPath() string {
	return "commands.marker.description"
}
func (k KeysCommandsMarkerDescription) GetPathSlice() []string {
	return []string{"commands", "marker", "description"}
}
func (k KeysCommandsMarkerDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsMarkerDescription) SetVars(vars KeysCommandsMarkerDescriptionVars) twiri18n.TranslationKey[KeysCommandsMarkerDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsMarker struct {
	Errors		KeysCommandsMarkerErrors
	Success		KeysCommandsMarkerSuccess
	Description	KeysCommandsMarkerDescription
}
type KeysCommandsSongDescriptionVars struct {
}
type KeysCommandsSongDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongDescription) IsTranslationKey() {
}
func (k KeysCommandsSongDescription) GetPath() string {
	return "commands.song.description"
}
func (k KeysCommandsSongDescription) GetPathSlice() []string {
	return []string{"commands", "song", "description"}
}
func (k KeysCommandsSongDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongDescription) SetVars(vars KeysCommandsSongDescriptionVars) twiri18n.TranslationKey[KeysCommandsSongDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongCurrentDescriptionVars struct {
}
type KeysCommandsSongCurrentDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongCurrentDescription) IsTranslationKey() {
}
func (k KeysCommandsSongCurrentDescription) GetPath() string {
	return "commands.song.current.description"
}
func (k KeysCommandsSongCurrentDescription) GetPathSlice() []string {
	return []string{"commands", "song", "current", "description"}
}
func (k KeysCommandsSongCurrentDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongCurrentDescription) SetVars(vars KeysCommandsSongCurrentDescriptionVars) twiri18n.TranslationKey[KeysCommandsSongCurrentDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongCurrentErrorsCannotGetCurrentSongVars struct {
}
type KeysCommandsSongCurrentErrorsCannotGetCurrentSong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongCurrentErrorsCannotGetCurrentSong) IsTranslationKey() {
}
func (k KeysCommandsSongCurrentErrorsCannotGetCurrentSong) GetPath() string {
	return "commands.song.current.errors.cannot_get_current_song"
}
func (k KeysCommandsSongCurrentErrorsCannotGetCurrentSong) GetPathSlice() []string {
	return []string{"commands", "song", "current", "errors", "cannot_get_current_song"}
}
func (k KeysCommandsSongCurrentErrorsCannotGetCurrentSong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongCurrentErrorsCannotGetCurrentSong) SetVars(vars KeysCommandsSongCurrentErrorsCannotGetCurrentSongVars) twiri18n.TranslationKey[KeysCommandsSongCurrentErrorsCannotGetCurrentSongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongCurrentErrorsNoSongPlayingVars struct {
}
type KeysCommandsSongCurrentErrorsNoSongPlaying struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongCurrentErrorsNoSongPlaying) IsTranslationKey() {
}
func (k KeysCommandsSongCurrentErrorsNoSongPlaying) GetPath() string {
	return "commands.song.current.errors.no_song_playing"
}
func (k KeysCommandsSongCurrentErrorsNoSongPlaying) GetPathSlice() []string {
	return []string{"commands", "song", "current", "errors", "no_song_playing"}
}
func (k KeysCommandsSongCurrentErrorsNoSongPlaying) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongCurrentErrorsNoSongPlaying) SetVars(vars KeysCommandsSongCurrentErrorsNoSongPlayingVars) twiri18n.TranslationKey[KeysCommandsSongCurrentErrorsNoSongPlayingVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongCurrentErrors struct {
	CannotGetCurrentSong	KeysCommandsSongCurrentErrorsCannotGetCurrentSong
	NoSongPlaying		KeysCommandsSongCurrentErrorsNoSongPlaying
}
type KeysCommandsSongCurrentSuccessCurrentSongVars struct {
	Song	any
	Artist	any
}
type KeysCommandsSongCurrentSuccessCurrentSong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongCurrentSuccessCurrentSong) IsTranslationKey() {
}
func (k KeysCommandsSongCurrentSuccessCurrentSong) GetPath() string {
	return "commands.song.current.success.current_song"
}
func (k KeysCommandsSongCurrentSuccessCurrentSong) GetPathSlice() []string {
	return []string{"commands", "song", "current", "success", "current_song"}
}
func (k KeysCommandsSongCurrentSuccessCurrentSong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongCurrentSuccessCurrentSong) SetVars(vars KeysCommandsSongCurrentSuccessCurrentSongVars) twiri18n.TranslationKey[KeysCommandsSongCurrentSuccessCurrentSongVars] {
	k.Vars = twiri18n.Vars{"song": vars.Song, "artist": vars.Artist}
	return k
}

type KeysCommandsSongCurrentSuccess struct {
	CurrentSong KeysCommandsSongCurrentSuccessCurrentSong
}
type KeysCommandsSongCurrent struct {
	Description	KeysCommandsSongCurrentDescription
	Errors		KeysCommandsSongCurrentErrors
	Success		KeysCommandsSongCurrentSuccess
}
type KeysCommandsSongSkipDescriptionVars struct {
}
type KeysCommandsSongSkipDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongSkipDescription) IsTranslationKey() {
}
func (k KeysCommandsSongSkipDescription) GetPath() string {
	return "commands.song.skip.description"
}
func (k KeysCommandsSongSkipDescription) GetPathSlice() []string {
	return []string{"commands", "song", "skip", "description"}
}
func (k KeysCommandsSongSkipDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongSkipDescription) SetVars(vars KeysCommandsSongSkipDescriptionVars) twiri18n.TranslationKey[KeysCommandsSongSkipDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongSkipErrorsCannotSkipSongVars struct {
}
type KeysCommandsSongSkipErrorsCannotSkipSong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongSkipErrorsCannotSkipSong) IsTranslationKey() {
}
func (k KeysCommandsSongSkipErrorsCannotSkipSong) GetPath() string {
	return "commands.song.skip.errors.cannot_skip_song"
}
func (k KeysCommandsSongSkipErrorsCannotSkipSong) GetPathSlice() []string {
	return []string{"commands", "song", "skip", "errors", "cannot_skip_song"}
}
func (k KeysCommandsSongSkipErrorsCannotSkipSong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongSkipErrorsCannotSkipSong) SetVars(vars KeysCommandsSongSkipErrorsCannotSkipSongVars) twiri18n.TranslationKey[KeysCommandsSongSkipErrorsCannotSkipSongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongSkipErrorsNoSongToSkipVars struct {
}
type KeysCommandsSongSkipErrorsNoSongToSkip struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongSkipErrorsNoSongToSkip) IsTranslationKey() {
}
func (k KeysCommandsSongSkipErrorsNoSongToSkip) GetPath() string {
	return "commands.song.skip.errors.no_song_to_skip"
}
func (k KeysCommandsSongSkipErrorsNoSongToSkip) GetPathSlice() []string {
	return []string{"commands", "song", "skip", "errors", "no_song_to_skip"}
}
func (k KeysCommandsSongSkipErrorsNoSongToSkip) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongSkipErrorsNoSongToSkip) SetVars(vars KeysCommandsSongSkipErrorsNoSongToSkipVars) twiri18n.TranslationKey[KeysCommandsSongSkipErrorsNoSongToSkipVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongSkipErrors struct {
	CannotSkipSong	KeysCommandsSongSkipErrorsCannotSkipSong
	NoSongToSkip	KeysCommandsSongSkipErrorsNoSongToSkip
}
type KeysCommandsSongSkipSuccessSongSkippedVars struct {
}
type KeysCommandsSongSkipSuccessSongSkipped struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongSkipSuccessSongSkipped) IsTranslationKey() {
}
func (k KeysCommandsSongSkipSuccessSongSkipped) GetPath() string {
	return "commands.song.skip.success.song_skipped"
}
func (k KeysCommandsSongSkipSuccessSongSkipped) GetPathSlice() []string {
	return []string{"commands", "song", "skip", "success", "song_skipped"}
}
func (k KeysCommandsSongSkipSuccessSongSkipped) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongSkipSuccessSongSkipped) SetVars(vars KeysCommandsSongSkipSuccessSongSkippedVars) twiri18n.TranslationKey[KeysCommandsSongSkipSuccessSongSkippedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongSkipSuccess struct {
	SongSkipped KeysCommandsSongSkipSuccessSongSkipped
}
type KeysCommandsSongSkip struct {
	Description	KeysCommandsSongSkipDescription
	Errors		KeysCommandsSongSkipErrors
	Success		KeysCommandsSongSkipSuccess
}
type KeysCommandsSongVolumeErrorsCannotSetVolumeVars struct {
}
type KeysCommandsSongVolumeErrorsCannotSetVolume struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongVolumeErrorsCannotSetVolume) IsTranslationKey() {
}
func (k KeysCommandsSongVolumeErrorsCannotSetVolume) GetPath() string {
	return "commands.song.volume.errors.cannot_set_volume"
}
func (k KeysCommandsSongVolumeErrorsCannotSetVolume) GetPathSlice() []string {
	return []string{"commands", "song", "volume", "errors", "cannot_set_volume"}
}
func (k KeysCommandsSongVolumeErrorsCannotSetVolume) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongVolumeErrorsCannotSetVolume) SetVars(vars KeysCommandsSongVolumeErrorsCannotSetVolumeVars) twiri18n.TranslationKey[KeysCommandsSongVolumeErrorsCannotSetVolumeVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongVolumeErrorsInvalidVolumeVars struct {
}
type KeysCommandsSongVolumeErrorsInvalidVolume struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongVolumeErrorsInvalidVolume) IsTranslationKey() {
}
func (k KeysCommandsSongVolumeErrorsInvalidVolume) GetPath() string {
	return "commands.song.volume.errors.invalid_volume"
}
func (k KeysCommandsSongVolumeErrorsInvalidVolume) GetPathSlice() []string {
	return []string{"commands", "song", "volume", "errors", "invalid_volume"}
}
func (k KeysCommandsSongVolumeErrorsInvalidVolume) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongVolumeErrorsInvalidVolume) SetVars(vars KeysCommandsSongVolumeErrorsInvalidVolumeVars) twiri18n.TranslationKey[KeysCommandsSongVolumeErrorsInvalidVolumeVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongVolumeErrors struct {
	CannotSetVolume	KeysCommandsSongVolumeErrorsCannotSetVolume
	InvalidVolume	KeysCommandsSongVolumeErrorsInvalidVolume
}
type KeysCommandsSongVolumeSuccessVolumeSetVars struct {
	Volume any
}
type KeysCommandsSongVolumeSuccessVolumeSet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongVolumeSuccessVolumeSet) IsTranslationKey() {
}
func (k KeysCommandsSongVolumeSuccessVolumeSet) GetPath() string {
	return "commands.song.volume.success.volume_set"
}
func (k KeysCommandsSongVolumeSuccessVolumeSet) GetPathSlice() []string {
	return []string{"commands", "song", "volume", "success", "volume_set"}
}
func (k KeysCommandsSongVolumeSuccessVolumeSet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongVolumeSuccessVolumeSet) SetVars(vars KeysCommandsSongVolumeSuccessVolumeSetVars) twiri18n.TranslationKey[KeysCommandsSongVolumeSuccessVolumeSetVars] {
	k.Vars = twiri18n.Vars{"volume": vars.Volume}
	return k
}

type KeysCommandsSongVolumeSuccess struct {
	VolumeSet KeysCommandsSongVolumeSuccessVolumeSet
}
type KeysCommandsSongVolumeDescriptionVars struct {
}
type KeysCommandsSongVolumeDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongVolumeDescription) IsTranslationKey() {
}
func (k KeysCommandsSongVolumeDescription) GetPath() string {
	return "commands.song.volume.description"
}
func (k KeysCommandsSongVolumeDescription) GetPathSlice() []string {
	return []string{"commands", "song", "volume", "description"}
}
func (k KeysCommandsSongVolumeDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongVolumeDescription) SetVars(vars KeysCommandsSongVolumeDescriptionVars) twiri18n.TranslationKey[KeysCommandsSongVolumeDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongVolume struct {
	Errors		KeysCommandsSongVolumeErrors
	Success		KeysCommandsSongVolumeSuccess
	Description	KeysCommandsSongVolumeDescription
}
type KeysCommandsSong struct {
	Description	KeysCommandsSongDescription
	Current		KeysCommandsSongCurrent
	Skip		KeysCommandsSongSkip
	Volume		KeysCommandsSongVolume
}
type KeysCommandsStatsDescriptionVars struct {
}
type KeysCommandsStatsDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsDescription) IsTranslationKey() {
}
func (k KeysCommandsStatsDescription) GetPath() string {
	return "commands.stats.description"
}
func (k KeysCommandsStatsDescription) GetPathSlice() []string {
	return []string{"commands", "stats", "description"}
}
func (k KeysCommandsStatsDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsDescription) SetVars(vars KeysCommandsStatsDescriptionVars) twiri18n.TranslationKey[KeysCommandsStatsDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsFollowageSuccessFollowageDisplayVars struct {
	User		any
	Duration	any
}
type KeysCommandsStatsFollowageSuccessFollowageDisplay struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsFollowageSuccessFollowageDisplay) IsTranslationKey() {
}
func (k KeysCommandsStatsFollowageSuccessFollowageDisplay) GetPath() string {
	return "commands.stats.followage.success.followage_display"
}
func (k KeysCommandsStatsFollowageSuccessFollowageDisplay) GetPathSlice() []string {
	return []string{"commands", "stats", "followage", "success", "followage_display"}
}
func (k KeysCommandsStatsFollowageSuccessFollowageDisplay) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsFollowageSuccessFollowageDisplay) SetVars(vars KeysCommandsStatsFollowageSuccessFollowageDisplayVars) twiri18n.TranslationKey[KeysCommandsStatsFollowageSuccessFollowageDisplayVars] {
	k.Vars = twiri18n.Vars{"user": vars.User, "duration": vars.Duration}
	return k
}

type KeysCommandsStatsFollowageSuccess struct {
	FollowageDisplay KeysCommandsStatsFollowageSuccessFollowageDisplay
}
type KeysCommandsStatsFollowageDescriptionVars struct {
}
type KeysCommandsStatsFollowageDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsFollowageDescription) IsTranslationKey() {
}
func (k KeysCommandsStatsFollowageDescription) GetPath() string {
	return "commands.stats.followage.description"
}
func (k KeysCommandsStatsFollowageDescription) GetPathSlice() []string {
	return []string{"commands", "stats", "followage", "description"}
}
func (k KeysCommandsStatsFollowageDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsFollowageDescription) SetVars(vars KeysCommandsStatsFollowageDescriptionVars) twiri18n.TranslationKey[KeysCommandsStatsFollowageDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsFollowageErrorsCannotGetFollowageVars struct {
}
type KeysCommandsStatsFollowageErrorsCannotGetFollowage struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsFollowageErrorsCannotGetFollowage) IsTranslationKey() {
}
func (k KeysCommandsStatsFollowageErrorsCannotGetFollowage) GetPath() string {
	return "commands.stats.followage.errors.cannot_get_followage"
}
func (k KeysCommandsStatsFollowageErrorsCannotGetFollowage) GetPathSlice() []string {
	return []string{"commands", "stats", "followage", "errors", "cannot_get_followage"}
}
func (k KeysCommandsStatsFollowageErrorsCannotGetFollowage) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsFollowageErrorsCannotGetFollowage) SetVars(vars KeysCommandsStatsFollowageErrorsCannotGetFollowageVars) twiri18n.TranslationKey[KeysCommandsStatsFollowageErrorsCannotGetFollowageVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsFollowageErrorsUserNotFollowingVars struct {
}
type KeysCommandsStatsFollowageErrorsUserNotFollowing struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsFollowageErrorsUserNotFollowing) IsTranslationKey() {
}
func (k KeysCommandsStatsFollowageErrorsUserNotFollowing) GetPath() string {
	return "commands.stats.followage.errors.user_not_following"
}
func (k KeysCommandsStatsFollowageErrorsUserNotFollowing) GetPathSlice() []string {
	return []string{"commands", "stats", "followage", "errors", "user_not_following"}
}
func (k KeysCommandsStatsFollowageErrorsUserNotFollowing) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsFollowageErrorsUserNotFollowing) SetVars(vars KeysCommandsStatsFollowageErrorsUserNotFollowingVars) twiri18n.TranslationKey[KeysCommandsStatsFollowageErrorsUserNotFollowingVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsFollowageErrors struct {
	CannotGetFollowage	KeysCommandsStatsFollowageErrorsCannotGetFollowage
	UserNotFollowing	KeysCommandsStatsFollowageErrorsUserNotFollowing
}
type KeysCommandsStatsFollowage struct {
	Success		KeysCommandsStatsFollowageSuccess
	Description	KeysCommandsStatsFollowageDescription
	Errors		KeysCommandsStatsFollowageErrors
}
type KeysCommandsStatsUptimeDescriptionVars struct {
}
type KeysCommandsStatsUptimeDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsUptimeDescription) IsTranslationKey() {
}
func (k KeysCommandsStatsUptimeDescription) GetPath() string {
	return "commands.stats.uptime.description"
}
func (k KeysCommandsStatsUptimeDescription) GetPathSlice() []string {
	return []string{"commands", "stats", "uptime", "description"}
}
func (k KeysCommandsStatsUptimeDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsUptimeDescription) SetVars(vars KeysCommandsStatsUptimeDescriptionVars) twiri18n.TranslationKey[KeysCommandsStatsUptimeDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsUptimeErrorsCannotGetUptimeVars struct {
}
type KeysCommandsStatsUptimeErrorsCannotGetUptime struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsUptimeErrorsCannotGetUptime) IsTranslationKey() {
}
func (k KeysCommandsStatsUptimeErrorsCannotGetUptime) GetPath() string {
	return "commands.stats.uptime.errors.cannot_get_uptime"
}
func (k KeysCommandsStatsUptimeErrorsCannotGetUptime) GetPathSlice() []string {
	return []string{"commands", "stats", "uptime", "errors", "cannot_get_uptime"}
}
func (k KeysCommandsStatsUptimeErrorsCannotGetUptime) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsUptimeErrorsCannotGetUptime) SetVars(vars KeysCommandsStatsUptimeErrorsCannotGetUptimeVars) twiri18n.TranslationKey[KeysCommandsStatsUptimeErrorsCannotGetUptimeVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsUptimeErrorsStreamOfflineVars struct {
}
type KeysCommandsStatsUptimeErrorsStreamOffline struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsUptimeErrorsStreamOffline) IsTranslationKey() {
}
func (k KeysCommandsStatsUptimeErrorsStreamOffline) GetPath() string {
	return "commands.stats.uptime.errors.stream_offline"
}
func (k KeysCommandsStatsUptimeErrorsStreamOffline) GetPathSlice() []string {
	return []string{"commands", "stats", "uptime", "errors", "stream_offline"}
}
func (k KeysCommandsStatsUptimeErrorsStreamOffline) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsUptimeErrorsStreamOffline) SetVars(vars KeysCommandsStatsUptimeErrorsStreamOfflineVars) twiri18n.TranslationKey[KeysCommandsStatsUptimeErrorsStreamOfflineVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsUptimeErrors struct {
	CannotGetUptime	KeysCommandsStatsUptimeErrorsCannotGetUptime
	StreamOffline	KeysCommandsStatsUptimeErrorsStreamOffline
}
type KeysCommandsStatsUptimeSuccessUptimeDisplayVars struct {
	Duration any
}
type KeysCommandsStatsUptimeSuccessUptimeDisplay struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsUptimeSuccessUptimeDisplay) IsTranslationKey() {
}
func (k KeysCommandsStatsUptimeSuccessUptimeDisplay) GetPath() string {
	return "commands.stats.uptime.success.uptime_display"
}
func (k KeysCommandsStatsUptimeSuccessUptimeDisplay) GetPathSlice() []string {
	return []string{"commands", "stats", "uptime", "success", "uptime_display"}
}
func (k KeysCommandsStatsUptimeSuccessUptimeDisplay) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsUptimeSuccessUptimeDisplay) SetVars(vars KeysCommandsStatsUptimeSuccessUptimeDisplayVars) twiri18n.TranslationKey[KeysCommandsStatsUptimeSuccessUptimeDisplayVars] {
	k.Vars = twiri18n.Vars{"duration": vars.Duration}
	return k
}

type KeysCommandsStatsUptimeSuccess struct {
	UptimeDisplay KeysCommandsStatsUptimeSuccessUptimeDisplay
}
type KeysCommandsStatsUptime struct {
	Description	KeysCommandsStatsUptimeDescription
	Errors		KeysCommandsStatsUptimeErrors
	Success		KeysCommandsStatsUptimeSuccess
}
type KeysCommandsStatsViewersDescriptionVars struct {
}
type KeysCommandsStatsViewersDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsViewersDescription) IsTranslationKey() {
}
func (k KeysCommandsStatsViewersDescription) GetPath() string {
	return "commands.stats.viewers.description"
}
func (k KeysCommandsStatsViewersDescription) GetPathSlice() []string {
	return []string{"commands", "stats", "viewers", "description"}
}
func (k KeysCommandsStatsViewersDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsViewersDescription) SetVars(vars KeysCommandsStatsViewersDescriptionVars) twiri18n.TranslationKey[KeysCommandsStatsViewersDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsViewersErrorsCannotGetViewersVars struct {
}
type KeysCommandsStatsViewersErrorsCannotGetViewers struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsViewersErrorsCannotGetViewers) IsTranslationKey() {
}
func (k KeysCommandsStatsViewersErrorsCannotGetViewers) GetPath() string {
	return "commands.stats.viewers.errors.cannot_get_viewers"
}
func (k KeysCommandsStatsViewersErrorsCannotGetViewers) GetPathSlice() []string {
	return []string{"commands", "stats", "viewers", "errors", "cannot_get_viewers"}
}
func (k KeysCommandsStatsViewersErrorsCannotGetViewers) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsViewersErrorsCannotGetViewers) SetVars(vars KeysCommandsStatsViewersErrorsCannotGetViewersVars) twiri18n.TranslationKey[KeysCommandsStatsViewersErrorsCannotGetViewersVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsViewersErrors struct {
	CannotGetViewers KeysCommandsStatsViewersErrorsCannotGetViewers
}
type KeysCommandsStatsViewersSuccessViewersDisplayVars struct {
	Count any
}
type KeysCommandsStatsViewersSuccessViewersDisplay struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsViewersSuccessViewersDisplay) IsTranslationKey() {
}
func (k KeysCommandsStatsViewersSuccessViewersDisplay) GetPath() string {
	return "commands.stats.viewers.success.viewers_display"
}
func (k KeysCommandsStatsViewersSuccessViewersDisplay) GetPathSlice() []string {
	return []string{"commands", "stats", "viewers", "success", "viewers_display"}
}
func (k KeysCommandsStatsViewersSuccessViewersDisplay) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsViewersSuccessViewersDisplay) SetVars(vars KeysCommandsStatsViewersSuccessViewersDisplayVars) twiri18n.TranslationKey[KeysCommandsStatsViewersSuccessViewersDisplayVars] {
	k.Vars = twiri18n.Vars{"count": vars.Count}
	return k
}

type KeysCommandsStatsViewersSuccess struct {
	ViewersDisplay KeysCommandsStatsViewersSuccessViewersDisplay
}
type KeysCommandsStatsViewers struct {
	Description	KeysCommandsStatsViewersDescription
	Errors		KeysCommandsStatsViewersErrors
	Success		KeysCommandsStatsViewersSuccess
}
type KeysCommandsStats struct {
	Description	KeysCommandsStatsDescription
	Followage	KeysCommandsStatsFollowage
	Uptime		KeysCommandsStatsUptime
	Viewers		KeysCommandsStatsViewers
}
type KeysCommandsSubageDescriptionVars struct {
}
type KeysCommandsSubageDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSubageDescription) IsTranslationKey() {
}
func (k KeysCommandsSubageDescription) GetPath() string {
	return "commands.subage.description"
}
func (k KeysCommandsSubageDescription) GetPathSlice() []string {
	return []string{"commands", "subage", "description"}
}
func (k KeysCommandsSubageDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSubageDescription) SetVars(vars KeysCommandsSubageDescriptionVars) twiri18n.TranslationKey[KeysCommandsSubageDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSubageErrorsNotSubscriberOrHiddenVars struct {
}
type KeysCommandsSubageErrorsNotSubscriberOrHidden struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSubageErrorsNotSubscriberOrHidden) IsTranslationKey() {
}
func (k KeysCommandsSubageErrorsNotSubscriberOrHidden) GetPath() string {
	return "commands.subage.errors.not_subscriber_or_hidden"
}
func (k KeysCommandsSubageErrorsNotSubscriberOrHidden) GetPathSlice() []string {
	return []string{"commands", "subage", "errors", "not_subscriber_or_hidden"}
}
func (k KeysCommandsSubageErrorsNotSubscriberOrHidden) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSubageErrorsNotSubscriberOrHidden) SetVars(vars KeysCommandsSubageErrorsNotSubscriberOrHiddenVars) twiri18n.TranslationKey[KeysCommandsSubageErrorsNotSubscriberOrHiddenVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSubageErrors struct {
	NotSubscriberOrHidden KeysCommandsSubageErrorsNotSubscriberOrHidden
}
type KeysCommandsSubageResponsesNotSubscriberVars struct {
	User any
}
type KeysCommandsSubageResponsesNotSubscriber struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSubageResponsesNotSubscriber) IsTranslationKey() {
}
func (k KeysCommandsSubageResponsesNotSubscriber) GetPath() string {
	return "commands.subage.responses.not_subscriber"
}
func (k KeysCommandsSubageResponsesNotSubscriber) GetPathSlice() []string {
	return []string{"commands", "subage", "responses", "not_subscriber"}
}
func (k KeysCommandsSubageResponsesNotSubscriber) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSubageResponsesNotSubscriber) SetVars(vars KeysCommandsSubageResponsesNotSubscriberVars) twiri18n.TranslationKey[KeysCommandsSubageResponsesNotSubscriberVars] {
	k.Vars = twiri18n.Vars{"user": vars.User}
	return k
}

type KeysCommandsSubageResponsesNotSubscriberButWasVars struct {
	User	any
	Months	any
}
type KeysCommandsSubageResponsesNotSubscriberButWas struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSubageResponsesNotSubscriberButWas) IsTranslationKey() {
}
func (k KeysCommandsSubageResponsesNotSubscriberButWas) GetPath() string {
	return "commands.subage.responses.not_subscriber_but_was"
}
func (k KeysCommandsSubageResponsesNotSubscriberButWas) GetPathSlice() []string {
	return []string{"commands", "subage", "responses", "not_subscriber_but_was"}
}
func (k KeysCommandsSubageResponsesNotSubscriberButWas) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSubageResponsesNotSubscriberButWas) SetVars(vars KeysCommandsSubageResponsesNotSubscriberButWasVars) twiri18n.TranslationKey[KeysCommandsSubageResponsesNotSubscriberButWasVars] {
	k.Vars = twiri18n.Vars{"user": vars.User, "months": vars.Months}
	return k
}

type KeysCommandsSubageResponsesSubscriptionInfoVars struct {
	User	any
	Tier	any
	Channel	any
	Months	any
}
type KeysCommandsSubageResponsesSubscriptionInfo struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSubageResponsesSubscriptionInfo) IsTranslationKey() {
}
func (k KeysCommandsSubageResponsesSubscriptionInfo) GetPath() string {
	return "commands.subage.responses.subscription_info"
}
func (k KeysCommandsSubageResponsesSubscriptionInfo) GetPathSlice() []string {
	return []string{"commands", "subage", "responses", "subscription_info"}
}
func (k KeysCommandsSubageResponsesSubscriptionInfo) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSubageResponsesSubscriptionInfo) SetVars(vars KeysCommandsSubageResponsesSubscriptionInfoVars) twiri18n.TranslationKey[KeysCommandsSubageResponsesSubscriptionInfoVars] {
	k.Vars = twiri18n.Vars{"user": vars.User, "tier": vars.Tier, "channel": vars.Channel, "months": vars.Months}
	return k
}

type KeysCommandsSubageResponsesStreakInfoVars struct {
	Months any
}
type KeysCommandsSubageResponsesStreakInfo struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSubageResponsesStreakInfo) IsTranslationKey() {
}
func (k KeysCommandsSubageResponsesStreakInfo) GetPath() string {
	return "commands.subage.responses.streak_info"
}
func (k KeysCommandsSubageResponsesStreakInfo) GetPathSlice() []string {
	return []string{"commands", "subage", "responses", "streak_info"}
}
func (k KeysCommandsSubageResponsesStreakInfo) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSubageResponsesStreakInfo) SetVars(vars KeysCommandsSubageResponsesStreakInfoVars) twiri18n.TranslationKey[KeysCommandsSubageResponsesStreakInfoVars] {
	k.Vars = twiri18n.Vars{"months": vars.Months}
	return k
}

type KeysCommandsSubageResponsesTimeRemainingVars struct {
	Duration any
}
type KeysCommandsSubageResponsesTimeRemaining struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSubageResponsesTimeRemaining) IsTranslationKey() {
}
func (k KeysCommandsSubageResponsesTimeRemaining) GetPath() string {
	return "commands.subage.responses.time_remaining"
}
func (k KeysCommandsSubageResponsesTimeRemaining) GetPathSlice() []string {
	return []string{"commands", "subage", "responses", "time_remaining"}
}
func (k KeysCommandsSubageResponsesTimeRemaining) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSubageResponsesTimeRemaining) SetVars(vars KeysCommandsSubageResponsesTimeRemainingVars) twiri18n.TranslationKey[KeysCommandsSubageResponsesTimeRemainingVars] {
	k.Vars = twiri18n.Vars{"duration": vars.Duration}
	return k
}

type KeysCommandsSubageResponses struct {
	NotSubscriber		KeysCommandsSubageResponsesNotSubscriber
	NotSubscriberButWas	KeysCommandsSubageResponsesNotSubscriberButWas
	SubscriptionInfo	KeysCommandsSubageResponsesSubscriptionInfo
	StreakInfo		KeysCommandsSubageResponsesStreakInfo
	TimeRemaining		KeysCommandsSubageResponsesTimeRemaining
}
type KeysCommandsSubage struct {
	Description	KeysCommandsSubageDescription
	Errors		KeysCommandsSubageErrors
	Responses	KeysCommandsSubageResponses
}
type KeysCommandsVipsAlreadyHaveRoleVars struct {
}
type KeysCommandsVipsAlreadyHaveRole struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsAlreadyHaveRole) IsTranslationKey() {
}
func (k KeysCommandsVipsAlreadyHaveRole) GetPath() string {
	return "commands.vips.already_have_role"
}
func (k KeysCommandsVipsAlreadyHaveRole) GetPathSlice() []string {
	return []string{"commands", "vips", "already_have_role"}
}
func (k KeysCommandsVipsAlreadyHaveRole) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsAlreadyHaveRole) SetVars(vars KeysCommandsVipsAlreadyHaveRoleVars) twiri18n.TranslationKey[KeysCommandsVipsAlreadyHaveRoleVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsAddedWithRemoveTimeVars struct {
	UserName	any
	EndTime		any
}
type KeysCommandsVipsAddedWithRemoveTime struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsAddedWithRemoveTime) IsTranslationKey() {
}
func (k KeysCommandsVipsAddedWithRemoveTime) GetPath() string {
	return "commands.vips.added_with_remove_time"
}
func (k KeysCommandsVipsAddedWithRemoveTime) GetPathSlice() []string {
	return []string{"commands", "vips", "added_with_remove_time"}
}
func (k KeysCommandsVipsAddedWithRemoveTime) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsAddedWithRemoveTime) SetVars(vars KeysCommandsVipsAddedWithRemoveTimeVars) twiri18n.TranslationKey[KeysCommandsVipsAddedWithRemoveTimeVars] {
	k.Vars = twiri18n.Vars{"userName": vars.UserName, "endTime": vars.EndTime}
	return k
}

type KeysCommandsVipsRemovedVars struct {
	UserName any
}
type KeysCommandsVipsRemoved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsRemoved) IsTranslationKey() {
}
func (k KeysCommandsVipsRemoved) GetPath() string {
	return "commands.vips.removed"
}
func (k KeysCommandsVipsRemoved) GetPathSlice() []string {
	return []string{"commands", "vips", "removed"}
}
func (k KeysCommandsVipsRemoved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsRemoved) SetVars(vars KeysCommandsVipsRemovedVars) twiri18n.TranslationKey[KeysCommandsVipsRemovedVars] {
	k.Vars = twiri18n.Vars{"userName": vars.UserName}
	return k
}

type KeysCommandsVipsCannotUpdateVars struct {
}
type KeysCommandsVipsCannotUpdate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsCannotUpdate) IsTranslationKey() {
}
func (k KeysCommandsVipsCannotUpdate) GetPath() string {
	return "commands.vips.cannot_update"
}
func (k KeysCommandsVipsCannotUpdate) GetPathSlice() []string {
	return []string{"commands", "vips", "cannot_update"}
}
func (k KeysCommandsVipsCannotUpdate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsCannotUpdate) SetVars(vars KeysCommandsVipsCannotUpdateVars) twiri18n.TranslationKey[KeysCommandsVipsCannotUpdateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsUpdatedVars struct {
	UserName	any
	EndTime		any
}
type KeysCommandsVipsUpdated struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsUpdated) IsTranslationKey() {
}
func (k KeysCommandsVipsUpdated) GetPath() string {
	return "commands.vips.updated"
}
func (k KeysCommandsVipsUpdated) GetPathSlice() []string {
	return []string{"commands", "vips", "updated"}
}
func (k KeysCommandsVipsUpdated) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsUpdated) SetVars(vars KeysCommandsVipsUpdatedVars) twiri18n.TranslationKey[KeysCommandsVipsUpdatedVars] {
	k.Vars = twiri18n.Vars{"userName": vars.UserName, "endTime": vars.EndTime}
	return k
}

type KeysCommandsVipsInvalidDurationVars struct {
}
type KeysCommandsVipsInvalidDuration struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsInvalidDuration) IsTranslationKey() {
}
func (k KeysCommandsVipsInvalidDuration) GetPath() string {
	return "commands.vips.invalid_duration"
}
func (k KeysCommandsVipsInvalidDuration) GetPathSlice() []string {
	return []string{"commands", "vips", "invalid_duration"}
}
func (k KeysCommandsVipsInvalidDuration) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsInvalidDuration) SetVars(vars KeysCommandsVipsInvalidDurationVars) twiri18n.TranslationKey[KeysCommandsVipsInvalidDurationVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsCannotCreateScheduledInDbVars struct {
}
type KeysCommandsVipsCannotCreateScheduledInDb struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsCannotCreateScheduledInDb) IsTranslationKey() {
}
func (k KeysCommandsVipsCannotCreateScheduledInDb) GetPath() string {
	return "commands.vips.cannot_create_scheduled_in_db"
}
func (k KeysCommandsVipsCannotCreateScheduledInDb) GetPathSlice() []string {
	return []string{"commands", "vips", "cannot_create_scheduled_in_db"}
}
func (k KeysCommandsVipsCannotCreateScheduledInDb) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsCannotCreateScheduledInDb) SetVars(vars KeysCommandsVipsCannotCreateScheduledInDbVars) twiri18n.TranslationKey[KeysCommandsVipsCannotCreateScheduledInDbVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsAddedVars struct {
	UserName any
}
type KeysCommandsVipsAdded struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsAdded) IsTranslationKey() {
}
func (k KeysCommandsVipsAdded) GetPath() string {
	return "commands.vips.added"
}
func (k KeysCommandsVipsAdded) GetPathSlice() []string {
	return []string{"commands", "vips", "added"}
}
func (k KeysCommandsVipsAdded) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsAdded) SetVars(vars KeysCommandsVipsAddedVars) twiri18n.TranslationKey[KeysCommandsVipsAddedVars] {
	k.Vars = twiri18n.Vars{"userName": vars.UserName}
	return k
}

type KeysCommandsVipsCannotGetListFromDbVars struct {
}
type KeysCommandsVipsCannotGetListFromDb struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsCannotGetListFromDb) IsTranslationKey() {
}
func (k KeysCommandsVipsCannotGetListFromDb) GetPath() string {
	return "commands.vips.cannot_get_list_from_db"
}
func (k KeysCommandsVipsCannotGetListFromDb) GetPathSlice() []string {
	return []string{"commands", "vips", "cannot_get_list_from_db"}
}
func (k KeysCommandsVipsCannotGetListFromDb) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsCannotGetListFromDb) SetVars(vars KeysCommandsVipsCannotGetListFromDbVars) twiri18n.TranslationKey[KeysCommandsVipsCannotGetListFromDbVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsNoScheduledVipsVars struct {
}
type KeysCommandsVipsNoScheduledVips struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsNoScheduledVips) IsTranslationKey() {
}
func (k KeysCommandsVipsNoScheduledVips) GetPath() string {
	return "commands.vips.no_scheduled_vips"
}
func (k KeysCommandsVipsNoScheduledVips) GetPathSlice() []string {
	return []string{"commands", "vips", "no_scheduled_vips"}
}
func (k KeysCommandsVipsNoScheduledVips) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsNoScheduledVips) SetVars(vars KeysCommandsVipsNoScheduledVipsVars) twiri18n.TranslationKey[KeysCommandsVipsNoScheduledVipsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVips struct {
	AlreadyHaveRole			KeysCommandsVipsAlreadyHaveRole
	AddedWithRemoveTime		KeysCommandsVipsAddedWithRemoveTime
	Removed				KeysCommandsVipsRemoved
	CannotUpdate			KeysCommandsVipsCannotUpdate
	Updated				KeysCommandsVipsUpdated
	InvalidDuration			KeysCommandsVipsInvalidDuration
	CannotCreateScheduledInDb	KeysCommandsVipsCannotCreateScheduledInDb
	Added				KeysCommandsVipsAdded
	CannotGetListFromDb		KeysCommandsVipsCannotGetListFromDb
	NoScheduledVips			KeysCommandsVipsNoScheduledVips
}
type KeysCommandsShoutoutDescriptionVars struct {
}
type KeysCommandsShoutoutDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShoutoutDescription) IsTranslationKey() {
}
func (k KeysCommandsShoutoutDescription) GetPath() string {
	return "commands.shoutout.description"
}
func (k KeysCommandsShoutoutDescription) GetPathSlice() []string {
	return []string{"commands", "shoutout", "description"}
}
func (k KeysCommandsShoutoutDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShoutoutDescription) SetVars(vars KeysCommandsShoutoutDescriptionVars) twiri18n.TranslationKey[KeysCommandsShoutoutDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsShoutoutUserSuccessUserPromotedVars struct {
	User	any
	Game	any
	Url	any
}
type KeysCommandsShoutoutUserSuccessUserPromoted struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShoutoutUserSuccessUserPromoted) IsTranslationKey() {
}
func (k KeysCommandsShoutoutUserSuccessUserPromoted) GetPath() string {
	return "commands.shoutout.user.success.user_promoted"
}
func (k KeysCommandsShoutoutUserSuccessUserPromoted) GetPathSlice() []string {
	return []string{"commands", "shoutout", "user", "success", "user_promoted"}
}
func (k KeysCommandsShoutoutUserSuccessUserPromoted) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShoutoutUserSuccessUserPromoted) SetVars(vars KeysCommandsShoutoutUserSuccessUserPromotedVars) twiri18n.TranslationKey[KeysCommandsShoutoutUserSuccessUserPromotedVars] {
	k.Vars = twiri18n.Vars{"user": vars.User, "game": vars.Game, "url": vars.Url}
	return k
}

type KeysCommandsShoutoutUserSuccess struct {
	UserPromoted KeysCommandsShoutoutUserSuccessUserPromoted
}
type KeysCommandsShoutoutUserDescriptionVars struct {
}
type KeysCommandsShoutoutUserDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShoutoutUserDescription) IsTranslationKey() {
}
func (k KeysCommandsShoutoutUserDescription) GetPath() string {
	return "commands.shoutout.user.description"
}
func (k KeysCommandsShoutoutUserDescription) GetPathSlice() []string {
	return []string{"commands", "shoutout", "user", "description"}
}
func (k KeysCommandsShoutoutUserDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShoutoutUserDescription) SetVars(vars KeysCommandsShoutoutUserDescriptionVars) twiri18n.TranslationKey[KeysCommandsShoutoutUserDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsShoutoutUserErrorsCannotShoutoutUserVars struct {
}
type KeysCommandsShoutoutUserErrorsCannotShoutoutUser struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShoutoutUserErrorsCannotShoutoutUser) IsTranslationKey() {
}
func (k KeysCommandsShoutoutUserErrorsCannotShoutoutUser) GetPath() string {
	return "commands.shoutout.user.errors.cannot_shoutout_user"
}
func (k KeysCommandsShoutoutUserErrorsCannotShoutoutUser) GetPathSlice() []string {
	return []string{"commands", "shoutout", "user", "errors", "cannot_shoutout_user"}
}
func (k KeysCommandsShoutoutUserErrorsCannotShoutoutUser) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShoutoutUserErrorsCannotShoutoutUser) SetVars(vars KeysCommandsShoutoutUserErrorsCannotShoutoutUserVars) twiri18n.TranslationKey[KeysCommandsShoutoutUserErrorsCannotShoutoutUserVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsShoutoutUserErrorsUserNotFoundVars struct {
}
type KeysCommandsShoutoutUserErrorsUserNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShoutoutUserErrorsUserNotFound) IsTranslationKey() {
}
func (k KeysCommandsShoutoutUserErrorsUserNotFound) GetPath() string {
	return "commands.shoutout.user.errors.user_not_found"
}
func (k KeysCommandsShoutoutUserErrorsUserNotFound) GetPathSlice() []string {
	return []string{"commands", "shoutout", "user", "errors", "user_not_found"}
}
func (k KeysCommandsShoutoutUserErrorsUserNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShoutoutUserErrorsUserNotFound) SetVars(vars KeysCommandsShoutoutUserErrorsUserNotFoundVars) twiri18n.TranslationKey[KeysCommandsShoutoutUserErrorsUserNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsShoutoutUserErrors struct {
	CannotShoutoutUser	KeysCommandsShoutoutUserErrorsCannotShoutoutUser
	UserNotFound		KeysCommandsShoutoutUserErrorsUserNotFound
}
type KeysCommandsShoutoutUser struct {
	Success		KeysCommandsShoutoutUserSuccess
	Description	KeysCommandsShoutoutUserDescription
	Errors		KeysCommandsShoutoutUserErrors
}
type KeysCommandsShoutoutContentDescriptionVars struct {
}
type KeysCommandsShoutoutContentDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShoutoutContentDescription) IsTranslationKey() {
}
func (k KeysCommandsShoutoutContentDescription) GetPath() string {
	return "commands.shoutout.content.description"
}
func (k KeysCommandsShoutoutContentDescription) GetPathSlice() []string {
	return []string{"commands", "shoutout", "content", "description"}
}
func (k KeysCommandsShoutoutContentDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShoutoutContentDescription) SetVars(vars KeysCommandsShoutoutContentDescriptionVars) twiri18n.TranslationKey[KeysCommandsShoutoutContentDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsShoutoutContentErrorsCannotPromoteContentVars struct {
}
type KeysCommandsShoutoutContentErrorsCannotPromoteContent struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShoutoutContentErrorsCannotPromoteContent) IsTranslationKey() {
}
func (k KeysCommandsShoutoutContentErrorsCannotPromoteContent) GetPath() string {
	return "commands.shoutout.content.errors.cannot_promote_content"
}
func (k KeysCommandsShoutoutContentErrorsCannotPromoteContent) GetPathSlice() []string {
	return []string{"commands", "shoutout", "content", "errors", "cannot_promote_content"}
}
func (k KeysCommandsShoutoutContentErrorsCannotPromoteContent) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShoutoutContentErrorsCannotPromoteContent) SetVars(vars KeysCommandsShoutoutContentErrorsCannotPromoteContentVars) twiri18n.TranslationKey[KeysCommandsShoutoutContentErrorsCannotPromoteContentVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsShoutoutContentErrors struct {
	CannotPromoteContent KeysCommandsShoutoutContentErrorsCannotPromoteContent
}
type KeysCommandsShoutoutContentSuccessContentPromotedVars struct {
	Content any
}
type KeysCommandsShoutoutContentSuccessContentPromoted struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShoutoutContentSuccessContentPromoted) IsTranslationKey() {
}
func (k KeysCommandsShoutoutContentSuccessContentPromoted) GetPath() string {
	return "commands.shoutout.content.success.content_promoted"
}
func (k KeysCommandsShoutoutContentSuccessContentPromoted) GetPathSlice() []string {
	return []string{"commands", "shoutout", "content", "success", "content_promoted"}
}
func (k KeysCommandsShoutoutContentSuccessContentPromoted) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShoutoutContentSuccessContentPromoted) SetVars(vars KeysCommandsShoutoutContentSuccessContentPromotedVars) twiri18n.TranslationKey[KeysCommandsShoutoutContentSuccessContentPromotedVars] {
	k.Vars = twiri18n.Vars{"content": vars.Content}
	return k
}

type KeysCommandsShoutoutContentSuccess struct {
	ContentPromoted KeysCommandsShoutoutContentSuccessContentPromoted
}
type KeysCommandsShoutoutContent struct {
	Description	KeysCommandsShoutoutContentDescription
	Errors		KeysCommandsShoutoutContentErrors
	Success		KeysCommandsShoutoutContentSuccess
}
type KeysCommandsShoutout struct {
	Description	KeysCommandsShoutoutDescription
	User		KeysCommandsShoutoutUser
	Content		KeysCommandsShoutoutContent
}
type KeysCommandsTtsDescriptionVars struct {
}
type KeysCommandsTtsDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsDescription) IsTranslationKey() {
}
func (k KeysCommandsTtsDescription) GetPath() string {
	return "commands.tts.description"
}
func (k KeysCommandsTtsDescription) GetPathSlice() []string {
	return []string{"commands", "tts", "description"}
}
func (k KeysCommandsTtsDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsDescription) SetVars(vars KeysCommandsTtsDescriptionVars) twiri18n.TranslationKey[KeysCommandsTtsDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsSayDescriptionVars struct {
}
type KeysCommandsTtsSayDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsSayDescription) IsTranslationKey() {
}
func (k KeysCommandsTtsSayDescription) GetPath() string {
	return "commands.tts.say.description"
}
func (k KeysCommandsTtsSayDescription) GetPathSlice() []string {
	return []string{"commands", "tts", "say", "description"}
}
func (k KeysCommandsTtsSayDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsSayDescription) SetVars(vars KeysCommandsTtsSayDescriptionVars) twiri18n.TranslationKey[KeysCommandsTtsSayDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsSayErrorsCannotUseTtsVars struct {
}
type KeysCommandsTtsSayErrorsCannotUseTts struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsSayErrorsCannotUseTts) IsTranslationKey() {
}
func (k KeysCommandsTtsSayErrorsCannotUseTts) GetPath() string {
	return "commands.tts.say.errors.cannot_use_tts"
}
func (k KeysCommandsTtsSayErrorsCannotUseTts) GetPathSlice() []string {
	return []string{"commands", "tts", "say", "errors", "cannot_use_tts"}
}
func (k KeysCommandsTtsSayErrorsCannotUseTts) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsSayErrorsCannotUseTts) SetVars(vars KeysCommandsTtsSayErrorsCannotUseTtsVars) twiri18n.TranslationKey[KeysCommandsTtsSayErrorsCannotUseTtsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsSayErrorsMessageTooLongVars struct {
}
type KeysCommandsTtsSayErrorsMessageTooLong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsSayErrorsMessageTooLong) IsTranslationKey() {
}
func (k KeysCommandsTtsSayErrorsMessageTooLong) GetPath() string {
	return "commands.tts.say.errors.message_too_long"
}
func (k KeysCommandsTtsSayErrorsMessageTooLong) GetPathSlice() []string {
	return []string{"commands", "tts", "say", "errors", "message_too_long"}
}
func (k KeysCommandsTtsSayErrorsMessageTooLong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsSayErrorsMessageTooLong) SetVars(vars KeysCommandsTtsSayErrorsMessageTooLongVars) twiri18n.TranslationKey[KeysCommandsTtsSayErrorsMessageTooLongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsSayErrorsTtsDisabledVars struct {
}
type KeysCommandsTtsSayErrorsTtsDisabled struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsSayErrorsTtsDisabled) IsTranslationKey() {
}
func (k KeysCommandsTtsSayErrorsTtsDisabled) GetPath() string {
	return "commands.tts.say.errors.tts_disabled"
}
func (k KeysCommandsTtsSayErrorsTtsDisabled) GetPathSlice() []string {
	return []string{"commands", "tts", "say", "errors", "tts_disabled"}
}
func (k KeysCommandsTtsSayErrorsTtsDisabled) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsSayErrorsTtsDisabled) SetVars(vars KeysCommandsTtsSayErrorsTtsDisabledVars) twiri18n.TranslationKey[KeysCommandsTtsSayErrorsTtsDisabledVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsSayErrors struct {
	CannotUseTts	KeysCommandsTtsSayErrorsCannotUseTts
	MessageTooLong	KeysCommandsTtsSayErrorsMessageTooLong
	TtsDisabled	KeysCommandsTtsSayErrorsTtsDisabled
}
type KeysCommandsTtsSaySuccessTtsQueuedVars struct {
}
type KeysCommandsTtsSaySuccessTtsQueued struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsSaySuccessTtsQueued) IsTranslationKey() {
}
func (k KeysCommandsTtsSaySuccessTtsQueued) GetPath() string {
	return "commands.tts.say.success.tts_queued"
}
func (k KeysCommandsTtsSaySuccessTtsQueued) GetPathSlice() []string {
	return []string{"commands", "tts", "say", "success", "tts_queued"}
}
func (k KeysCommandsTtsSaySuccessTtsQueued) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsSaySuccessTtsQueued) SetVars(vars KeysCommandsTtsSaySuccessTtsQueuedVars) twiri18n.TranslationKey[KeysCommandsTtsSaySuccessTtsQueuedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsSaySuccess struct {
	TtsQueued KeysCommandsTtsSaySuccessTtsQueued
}
type KeysCommandsTtsSay struct {
	Description	KeysCommandsTtsSayDescription
	Errors		KeysCommandsTtsSayErrors
	Success		KeysCommandsTtsSaySuccess
}
type KeysCommandsTtsSkipDescriptionVars struct {
}
type KeysCommandsTtsSkipDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsSkipDescription) IsTranslationKey() {
}
func (k KeysCommandsTtsSkipDescription) GetPath() string {
	return "commands.tts.skip.description"
}
func (k KeysCommandsTtsSkipDescription) GetPathSlice() []string {
	return []string{"commands", "tts", "skip", "description"}
}
func (k KeysCommandsTtsSkipDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsSkipDescription) SetVars(vars KeysCommandsTtsSkipDescriptionVars) twiri18n.TranslationKey[KeysCommandsTtsSkipDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsSkipErrorsCannotSkipTtsVars struct {
}
type KeysCommandsTtsSkipErrorsCannotSkipTts struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsSkipErrorsCannotSkipTts) IsTranslationKey() {
}
func (k KeysCommandsTtsSkipErrorsCannotSkipTts) GetPath() string {
	return "commands.tts.skip.errors.cannot_skip_tts"
}
func (k KeysCommandsTtsSkipErrorsCannotSkipTts) GetPathSlice() []string {
	return []string{"commands", "tts", "skip", "errors", "cannot_skip_tts"}
}
func (k KeysCommandsTtsSkipErrorsCannotSkipTts) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsSkipErrorsCannotSkipTts) SetVars(vars KeysCommandsTtsSkipErrorsCannotSkipTtsVars) twiri18n.TranslationKey[KeysCommandsTtsSkipErrorsCannotSkipTtsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsSkipErrorsNoTtsPlayingVars struct {
}
type KeysCommandsTtsSkipErrorsNoTtsPlaying struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsSkipErrorsNoTtsPlaying) IsTranslationKey() {
}
func (k KeysCommandsTtsSkipErrorsNoTtsPlaying) GetPath() string {
	return "commands.tts.skip.errors.no_tts_playing"
}
func (k KeysCommandsTtsSkipErrorsNoTtsPlaying) GetPathSlice() []string {
	return []string{"commands", "tts", "skip", "errors", "no_tts_playing"}
}
func (k KeysCommandsTtsSkipErrorsNoTtsPlaying) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsSkipErrorsNoTtsPlaying) SetVars(vars KeysCommandsTtsSkipErrorsNoTtsPlayingVars) twiri18n.TranslationKey[KeysCommandsTtsSkipErrorsNoTtsPlayingVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsSkipErrors struct {
	CannotSkipTts	KeysCommandsTtsSkipErrorsCannotSkipTts
	NoTtsPlaying	KeysCommandsTtsSkipErrorsNoTtsPlaying
}
type KeysCommandsTtsSkipSuccessTtsSkippedVars struct {
}
type KeysCommandsTtsSkipSuccessTtsSkipped struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsSkipSuccessTtsSkipped) IsTranslationKey() {
}
func (k KeysCommandsTtsSkipSuccessTtsSkipped) GetPath() string {
	return "commands.tts.skip.success.tts_skipped"
}
func (k KeysCommandsTtsSkipSuccessTtsSkipped) GetPathSlice() []string {
	return []string{"commands", "tts", "skip", "success", "tts_skipped"}
}
func (k KeysCommandsTtsSkipSuccessTtsSkipped) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsSkipSuccessTtsSkipped) SetVars(vars KeysCommandsTtsSkipSuccessTtsSkippedVars) twiri18n.TranslationKey[KeysCommandsTtsSkipSuccessTtsSkippedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsSkipSuccess struct {
	TtsSkipped KeysCommandsTtsSkipSuccessTtsSkipped
}
type KeysCommandsTtsSkip struct {
	Description	KeysCommandsTtsSkipDescription
	Errors		KeysCommandsTtsSkipErrors
	Success		KeysCommandsTtsSkipSuccess
}
type KeysCommandsTts struct {
	Description	KeysCommandsTtsDescription
	Say		KeysCommandsTtsSay
	Skip		KeysCommandsTtsSkip
}
type KeysCommandsChatWallTimeoutDescriptionVars struct {
}
type KeysCommandsChatWallTimeoutDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallTimeoutDescription) IsTranslationKey() {
}
func (k KeysCommandsChatWallTimeoutDescription) GetPath() string {
	return "commands.chat_wall.timeout.description"
}
func (k KeysCommandsChatWallTimeoutDescription) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "timeout", "description"}
}
func (k KeysCommandsChatWallTimeoutDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallTimeoutDescription) SetVars(vars KeysCommandsChatWallTimeoutDescriptionVars) twiri18n.TranslationKey[KeysCommandsChatWallTimeoutDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallTimeoutErrorsCannotTimeoutUserVars struct {
}
type KeysCommandsChatWallTimeoutErrorsCannotTimeoutUser struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallTimeoutErrorsCannotTimeoutUser) IsTranslationKey() {
}
func (k KeysCommandsChatWallTimeoutErrorsCannotTimeoutUser) GetPath() string {
	return "commands.chat_wall.timeout.errors.cannot_timeout_user"
}
func (k KeysCommandsChatWallTimeoutErrorsCannotTimeoutUser) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "timeout", "errors", "cannot_timeout_user"}
}
func (k KeysCommandsChatWallTimeoutErrorsCannotTimeoutUser) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallTimeoutErrorsCannotTimeoutUser) SetVars(vars KeysCommandsChatWallTimeoutErrorsCannotTimeoutUserVars) twiri18n.TranslationKey[KeysCommandsChatWallTimeoutErrorsCannotTimeoutUserVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallTimeoutErrors struct {
	CannotTimeoutUser KeysCommandsChatWallTimeoutErrorsCannotTimeoutUser
}
type KeysCommandsChatWallTimeoutSuccessUserTimedOutVars struct {
}
type KeysCommandsChatWallTimeoutSuccessUserTimedOut struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallTimeoutSuccessUserTimedOut) IsTranslationKey() {
}
func (k KeysCommandsChatWallTimeoutSuccessUserTimedOut) GetPath() string {
	return "commands.chat_wall.timeout.success.user_timed_out"
}
func (k KeysCommandsChatWallTimeoutSuccessUserTimedOut) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "timeout", "success", "user_timed_out"}
}
func (k KeysCommandsChatWallTimeoutSuccessUserTimedOut) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallTimeoutSuccessUserTimedOut) SetVars(vars KeysCommandsChatWallTimeoutSuccessUserTimedOutVars) twiri18n.TranslationKey[KeysCommandsChatWallTimeoutSuccessUserTimedOutVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallTimeoutSuccess struct {
	UserTimedOut KeysCommandsChatWallTimeoutSuccessUserTimedOut
}
type KeysCommandsChatWallTimeout struct {
	Description	KeysCommandsChatWallTimeoutDescription
	Errors		KeysCommandsChatWallTimeoutErrors
	Success		KeysCommandsChatWallTimeoutSuccess
}
type KeysCommandsChatWallStopDescriptionVars struct {
}
type KeysCommandsChatWallStopDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallStopDescription) IsTranslationKey() {
}
func (k KeysCommandsChatWallStopDescription) GetPath() string {
	return "commands.chat_wall.stop.description"
}
func (k KeysCommandsChatWallStopDescription) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "stop", "description"}
}
func (k KeysCommandsChatWallStopDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallStopDescription) SetVars(vars KeysCommandsChatWallStopDescriptionVars) twiri18n.TranslationKey[KeysCommandsChatWallStopDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallStopErrorsCannotStopChatWallVars struct {
}
type KeysCommandsChatWallStopErrorsCannotStopChatWall struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallStopErrorsCannotStopChatWall) IsTranslationKey() {
}
func (k KeysCommandsChatWallStopErrorsCannotStopChatWall) GetPath() string {
	return "commands.chat_wall.stop.errors.cannot_stop_chat_wall"
}
func (k KeysCommandsChatWallStopErrorsCannotStopChatWall) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "stop", "errors", "cannot_stop_chat_wall"}
}
func (k KeysCommandsChatWallStopErrorsCannotStopChatWall) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallStopErrorsCannotStopChatWall) SetVars(vars KeysCommandsChatWallStopErrorsCannotStopChatWallVars) twiri18n.TranslationKey[KeysCommandsChatWallStopErrorsCannotStopChatWallVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallStopErrors struct {
	CannotStopChatWall KeysCommandsChatWallStopErrorsCannotStopChatWall
}
type KeysCommandsChatWallStopSuccessChatWallStoppedVars struct {
}
type KeysCommandsChatWallStopSuccessChatWallStopped struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallStopSuccessChatWallStopped) IsTranslationKey() {
}
func (k KeysCommandsChatWallStopSuccessChatWallStopped) GetPath() string {
	return "commands.chat_wall.stop.success.chat_wall_stopped"
}
func (k KeysCommandsChatWallStopSuccessChatWallStopped) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "stop", "success", "chat_wall_stopped"}
}
func (k KeysCommandsChatWallStopSuccessChatWallStopped) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallStopSuccessChatWallStopped) SetVars(vars KeysCommandsChatWallStopSuccessChatWallStoppedVars) twiri18n.TranslationKey[KeysCommandsChatWallStopSuccessChatWallStoppedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallStopSuccess struct {
	ChatWallStopped KeysCommandsChatWallStopSuccessChatWallStopped
}
type KeysCommandsChatWallStop struct {
	Description	KeysCommandsChatWallStopDescription
	Errors		KeysCommandsChatWallStopErrors
	Success		KeysCommandsChatWallStopSuccess
}
type KeysCommandsChatWallBanErrorsCannotBanUserVars struct {
}
type KeysCommandsChatWallBanErrorsCannotBanUser struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallBanErrorsCannotBanUser) IsTranslationKey() {
}
func (k KeysCommandsChatWallBanErrorsCannotBanUser) GetPath() string {
	return "commands.chat_wall.ban.errors.cannot_ban_user"
}
func (k KeysCommandsChatWallBanErrorsCannotBanUser) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "ban", "errors", "cannot_ban_user"}
}
func (k KeysCommandsChatWallBanErrorsCannotBanUser) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallBanErrorsCannotBanUser) SetVars(vars KeysCommandsChatWallBanErrorsCannotBanUserVars) twiri18n.TranslationKey[KeysCommandsChatWallBanErrorsCannotBanUserVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallBanErrors struct {
	CannotBanUser KeysCommandsChatWallBanErrorsCannotBanUser
}
type KeysCommandsChatWallBanSuccessUserBannedVars struct {
}
type KeysCommandsChatWallBanSuccessUserBanned struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallBanSuccessUserBanned) IsTranslationKey() {
}
func (k KeysCommandsChatWallBanSuccessUserBanned) GetPath() string {
	return "commands.chat_wall.ban.success.user_banned"
}
func (k KeysCommandsChatWallBanSuccessUserBanned) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "ban", "success", "user_banned"}
}
func (k KeysCommandsChatWallBanSuccessUserBanned) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallBanSuccessUserBanned) SetVars(vars KeysCommandsChatWallBanSuccessUserBannedVars) twiri18n.TranslationKey[KeysCommandsChatWallBanSuccessUserBannedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallBanSuccess struct {
	UserBanned KeysCommandsChatWallBanSuccessUserBanned
}
type KeysCommandsChatWallBanDescriptionVars struct {
}
type KeysCommandsChatWallBanDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallBanDescription) IsTranslationKey() {
}
func (k KeysCommandsChatWallBanDescription) GetPath() string {
	return "commands.chat_wall.ban.description"
}
func (k KeysCommandsChatWallBanDescription) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "ban", "description"}
}
func (k KeysCommandsChatWallBanDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallBanDescription) SetVars(vars KeysCommandsChatWallBanDescriptionVars) twiri18n.TranslationKey[KeysCommandsChatWallBanDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallBan struct {
	Errors		KeysCommandsChatWallBanErrors
	Success		KeysCommandsChatWallBanSuccess
	Description	KeysCommandsChatWallBanDescription
}
type KeysCommandsChatWall struct {
	Timeout	KeysCommandsChatWallTimeout
	Stop	KeysCommandsChatWallStop
	Ban	KeysCommandsChatWallBan
}
type KeysCommandsShorturlErrorsCannotCreateShortUrlVars struct {
	Error any
}
type KeysCommandsShorturlErrorsCannotCreateShortUrl struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShorturlErrorsCannotCreateShortUrl) IsTranslationKey() {
}
func (k KeysCommandsShorturlErrorsCannotCreateShortUrl) GetPath() string {
	return "commands.shorturl.errors.cannot_create_short_url"
}
func (k KeysCommandsShorturlErrorsCannotCreateShortUrl) GetPathSlice() []string {
	return []string{"commands", "shorturl", "errors", "cannot_create_short_url"}
}
func (k KeysCommandsShorturlErrorsCannotCreateShortUrl) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShorturlErrorsCannotCreateShortUrl) SetVars(vars KeysCommandsShorturlErrorsCannotCreateShortUrlVars) twiri18n.TranslationKey[KeysCommandsShorturlErrorsCannotCreateShortUrlVars] {
	k.Vars = twiri18n.Vars{"error": vars.Error}
	return k
}

type KeysCommandsShorturlErrors struct {
	CannotCreateShortUrl KeysCommandsShorturlErrorsCannotCreateShortUrl
}
type KeysCommandsShorturlSuccessShortUrlCreatedVars struct {
	Url any
}
type KeysCommandsShorturlSuccessShortUrlCreated struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShorturlSuccessShortUrlCreated) IsTranslationKey() {
}
func (k KeysCommandsShorturlSuccessShortUrlCreated) GetPath() string {
	return "commands.shorturl.success.short_url_created"
}
func (k KeysCommandsShorturlSuccessShortUrlCreated) GetPathSlice() []string {
	return []string{"commands", "shorturl", "success", "short_url_created"}
}
func (k KeysCommandsShorturlSuccessShortUrlCreated) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShorturlSuccessShortUrlCreated) SetVars(vars KeysCommandsShorturlSuccessShortUrlCreatedVars) twiri18n.TranslationKey[KeysCommandsShorturlSuccessShortUrlCreatedVars] {
	k.Vars = twiri18n.Vars{"url": vars.Url}
	return k
}

type KeysCommandsShorturlSuccess struct {
	ShortUrlCreated KeysCommandsShorturlSuccessShortUrlCreated
}
type KeysCommandsShorturlDescriptionVars struct {
}
type KeysCommandsShorturlDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShorturlDescription) IsTranslationKey() {
}
func (k KeysCommandsShorturlDescription) GetPath() string {
	return "commands.shorturl.description"
}
func (k KeysCommandsShorturlDescription) GetPathSlice() []string {
	return []string{"commands", "shorturl", "description"}
}
func (k KeysCommandsShorturlDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShorturlDescription) SetVars(vars KeysCommandsShorturlDescriptionVars) twiri18n.TranslationKey[KeysCommandsShorturlDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsShorturl struct {
	Errors		KeysCommandsShorturlErrors
	Success		KeysCommandsShorturlSuccess
	Description	KeysCommandsShorturlDescription
}
type KeysCommandsCategoriesAliasesDescriptionVars struct {
}
type KeysCommandsCategoriesAliasesDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesDescription) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesDescription) GetPath() string {
	return "commands.categories_aliases.description"
}
func (k KeysCommandsCategoriesAliasesDescription) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "description"}
}
func (k KeysCommandsCategoriesAliasesDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesDescription) SetVars(vars KeysCommandsCategoriesAliasesDescriptionVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesSetErrorsCannotSetCategoryVars struct {
}
type KeysCommandsCategoriesAliasesSetErrorsCannotSetCategory struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesSetErrorsCannotSetCategory) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesSetErrorsCannotSetCategory) GetPath() string {
	return "commands.categories_aliases.set.errors.cannot_set_category"
}
func (k KeysCommandsCategoriesAliasesSetErrorsCannotSetCategory) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "set", "errors", "cannot_set_category"}
}
func (k KeysCommandsCategoriesAliasesSetErrorsCannotSetCategory) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesSetErrorsCannotSetCategory) SetVars(vars KeysCommandsCategoriesAliasesSetErrorsCannotSetCategoryVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesSetErrorsCannotSetCategoryVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesSetErrorsCategoryNotFoundVars struct {
}
type KeysCommandsCategoriesAliasesSetErrorsCategoryNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesSetErrorsCategoryNotFound) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesSetErrorsCategoryNotFound) GetPath() string {
	return "commands.categories_aliases.set.errors.category_not_found"
}
func (k KeysCommandsCategoriesAliasesSetErrorsCategoryNotFound) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "set", "errors", "category_not_found"}
}
func (k KeysCommandsCategoriesAliasesSetErrorsCategoryNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesSetErrorsCategoryNotFound) SetVars(vars KeysCommandsCategoriesAliasesSetErrorsCategoryNotFoundVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesSetErrorsCategoryNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesSetErrors struct {
	CannotSetCategory	KeysCommandsCategoriesAliasesSetErrorsCannotSetCategory
	CategoryNotFound	KeysCommandsCategoriesAliasesSetErrorsCategoryNotFound
}
type KeysCommandsCategoriesAliasesSetSuccessCategorySetVars struct {
	Category any
}
type KeysCommandsCategoriesAliasesSetSuccessCategorySet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesSetSuccessCategorySet) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesSetSuccessCategorySet) GetPath() string {
	return "commands.categories_aliases.set.success.category_set"
}
func (k KeysCommandsCategoriesAliasesSetSuccessCategorySet) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "set", "success", "category_set"}
}
func (k KeysCommandsCategoriesAliasesSetSuccessCategorySet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesSetSuccessCategorySet) SetVars(vars KeysCommandsCategoriesAliasesSetSuccessCategorySetVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesSetSuccessCategorySetVars] {
	k.Vars = twiri18n.Vars{"category": vars.Category}
	return k
}

type KeysCommandsCategoriesAliasesSetSuccess struct {
	CategorySet KeysCommandsCategoriesAliasesSetSuccessCategorySet
}
type KeysCommandsCategoriesAliasesSetDescriptionVars struct {
}
type KeysCommandsCategoriesAliasesSetDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesSetDescription) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesSetDescription) GetPath() string {
	return "commands.categories_aliases.set.description"
}
func (k KeysCommandsCategoriesAliasesSetDescription) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "set", "description"}
}
func (k KeysCommandsCategoriesAliasesSetDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesSetDescription) SetVars(vars KeysCommandsCategoriesAliasesSetDescriptionVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesSetDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesSet struct {
	Errors		KeysCommandsCategoriesAliasesSetErrors
	Success		KeysCommandsCategoriesAliasesSetSuccess
	Description	KeysCommandsCategoriesAliasesSetDescription
}
type KeysCommandsCategoriesAliasesAliasesErrorsAliasNotFoundVars struct {
}
type KeysCommandsCategoriesAliasesAliasesErrorsAliasNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesAliasesErrorsAliasNotFound) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesAliasesErrorsAliasNotFound) GetPath() string {
	return "commands.categories_aliases.aliases.errors.alias_not_found"
}
func (k KeysCommandsCategoriesAliasesAliasesErrorsAliasNotFound) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "aliases", "errors", "alias_not_found"}
}
func (k KeysCommandsCategoriesAliasesAliasesErrorsAliasNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesAliasesErrorsAliasNotFound) SetVars(vars KeysCommandsCategoriesAliasesAliasesErrorsAliasNotFoundVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesAliasesErrorsAliasNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesAliasesErrorsCannotManageAliasesVars struct {
}
type KeysCommandsCategoriesAliasesAliasesErrorsCannotManageAliases struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesAliasesErrorsCannotManageAliases) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesAliasesErrorsCannotManageAliases) GetPath() string {
	return "commands.categories_aliases.aliases.errors.cannot_manage_aliases"
}
func (k KeysCommandsCategoriesAliasesAliasesErrorsCannotManageAliases) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "aliases", "errors", "cannot_manage_aliases"}
}
func (k KeysCommandsCategoriesAliasesAliasesErrorsCannotManageAliases) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesAliasesErrorsCannotManageAliases) SetVars(vars KeysCommandsCategoriesAliasesAliasesErrorsCannotManageAliasesVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesAliasesErrorsCannotManageAliasesVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesAliasesErrors struct {
	AliasNotFound		KeysCommandsCategoriesAliasesAliasesErrorsAliasNotFound
	CannotManageAliases	KeysCommandsCategoriesAliasesAliasesErrorsCannotManageAliases
}
type KeysCommandsCategoriesAliasesAliasesSuccessAliasAddedVars struct {
	Alias any
}
type KeysCommandsCategoriesAliasesAliasesSuccessAliasAdded struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesAliasesSuccessAliasAdded) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesAliasesSuccessAliasAdded) GetPath() string {
	return "commands.categories_aliases.aliases.success.alias_added"
}
func (k KeysCommandsCategoriesAliasesAliasesSuccessAliasAdded) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "aliases", "success", "alias_added"}
}
func (k KeysCommandsCategoriesAliasesAliasesSuccessAliasAdded) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesAliasesSuccessAliasAdded) SetVars(vars KeysCommandsCategoriesAliasesAliasesSuccessAliasAddedVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesAliasesSuccessAliasAddedVars] {
	k.Vars = twiri18n.Vars{"alias": vars.Alias}
	return k
}

type KeysCommandsCategoriesAliasesAliasesSuccessAliasRemovedVars struct {
	Alias any
}
type KeysCommandsCategoriesAliasesAliasesSuccessAliasRemoved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesAliasesSuccessAliasRemoved) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesAliasesSuccessAliasRemoved) GetPath() string {
	return "commands.categories_aliases.aliases.success.alias_removed"
}
func (k KeysCommandsCategoriesAliasesAliasesSuccessAliasRemoved) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "aliases", "success", "alias_removed"}
}
func (k KeysCommandsCategoriesAliasesAliasesSuccessAliasRemoved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesAliasesSuccessAliasRemoved) SetVars(vars KeysCommandsCategoriesAliasesAliasesSuccessAliasRemovedVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesAliasesSuccessAliasRemovedVars] {
	k.Vars = twiri18n.Vars{"alias": vars.Alias}
	return k
}

type KeysCommandsCategoriesAliasesAliasesSuccess struct {
	AliasAdded	KeysCommandsCategoriesAliasesAliasesSuccessAliasAdded
	AliasRemoved	KeysCommandsCategoriesAliasesAliasesSuccessAliasRemoved
}
type KeysCommandsCategoriesAliasesAliasesDescriptionVars struct {
}
type KeysCommandsCategoriesAliasesAliasesDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesAliasesDescription) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesAliasesDescription) GetPath() string {
	return "commands.categories_aliases.aliases.description"
}
func (k KeysCommandsCategoriesAliasesAliasesDescription) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "aliases", "description"}
}
func (k KeysCommandsCategoriesAliasesAliasesDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesAliasesDescription) SetVars(vars KeysCommandsCategoriesAliasesAliasesDescriptionVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesAliasesDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesAliases struct {
	Errors		KeysCommandsCategoriesAliasesAliasesErrors
	Success		KeysCommandsCategoriesAliasesAliasesSuccess
	Description	KeysCommandsCategoriesAliasesAliasesDescription
}
type KeysCommandsCategoriesAliases struct {
	Description	KeysCommandsCategoriesAliasesDescription
	Set		KeysCommandsCategoriesAliasesSet
	Aliases		KeysCommandsCategoriesAliasesAliases
}
type KeysCommandsOverlaysHideDescriptionVars struct {
}
type KeysCommandsOverlaysHideDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsOverlaysHideDescription) IsTranslationKey() {
}
func (k KeysCommandsOverlaysHideDescription) GetPath() string {
	return "commands.overlays.hide.description"
}
func (k KeysCommandsOverlaysHideDescription) GetPathSlice() []string {
	return []string{"commands", "overlays", "hide", "description"}
}
func (k KeysCommandsOverlaysHideDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsOverlaysHideDescription) SetVars(vars KeysCommandsOverlaysHideDescriptionVars) twiri18n.TranslationKey[KeysCommandsOverlaysHideDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsOverlaysHideErrorsCannotHideOverlayVars struct {
}
type KeysCommandsOverlaysHideErrorsCannotHideOverlay struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsOverlaysHideErrorsCannotHideOverlay) IsTranslationKey() {
}
func (k KeysCommandsOverlaysHideErrorsCannotHideOverlay) GetPath() string {
	return "commands.overlays.hide.errors.cannot_hide_overlay"
}
func (k KeysCommandsOverlaysHideErrorsCannotHideOverlay) GetPathSlice() []string {
	return []string{"commands", "overlays", "hide", "errors", "cannot_hide_overlay"}
}
func (k KeysCommandsOverlaysHideErrorsCannotHideOverlay) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsOverlaysHideErrorsCannotHideOverlay) SetVars(vars KeysCommandsOverlaysHideErrorsCannotHideOverlayVars) twiri18n.TranslationKey[KeysCommandsOverlaysHideErrorsCannotHideOverlayVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsOverlaysHideErrors struct {
	CannotHideOverlay KeysCommandsOverlaysHideErrorsCannotHideOverlay
}
type KeysCommandsOverlaysHideSuccessOverlayHiddenVars struct {
	Overlay any
}
type KeysCommandsOverlaysHideSuccessOverlayHidden struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsOverlaysHideSuccessOverlayHidden) IsTranslationKey() {
}
func (k KeysCommandsOverlaysHideSuccessOverlayHidden) GetPath() string {
	return "commands.overlays.hide.success.overlay_hidden"
}
func (k KeysCommandsOverlaysHideSuccessOverlayHidden) GetPathSlice() []string {
	return []string{"commands", "overlays", "hide", "success", "overlay_hidden"}
}
func (k KeysCommandsOverlaysHideSuccessOverlayHidden) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsOverlaysHideSuccessOverlayHidden) SetVars(vars KeysCommandsOverlaysHideSuccessOverlayHiddenVars) twiri18n.TranslationKey[KeysCommandsOverlaysHideSuccessOverlayHiddenVars] {
	k.Vars = twiri18n.Vars{"overlay": vars.Overlay}
	return k
}

type KeysCommandsOverlaysHideSuccess struct {
	OverlayHidden KeysCommandsOverlaysHideSuccessOverlayHidden
}
type KeysCommandsOverlaysHide struct {
	Description	KeysCommandsOverlaysHideDescription
	Errors		KeysCommandsOverlaysHideErrors
	Success		KeysCommandsOverlaysHideSuccess
}
type KeysCommandsOverlaysDescriptionVars struct {
}
type KeysCommandsOverlaysDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsOverlaysDescription) IsTranslationKey() {
}
func (k KeysCommandsOverlaysDescription) GetPath() string {
	return "commands.overlays.description"
}
func (k KeysCommandsOverlaysDescription) GetPathSlice() []string {
	return []string{"commands", "overlays", "description"}
}
func (k KeysCommandsOverlaysDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsOverlaysDescription) SetVars(vars KeysCommandsOverlaysDescriptionVars) twiri18n.TranslationKey[KeysCommandsOverlaysDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsOverlaysShowDescriptionVars struct {
}
type KeysCommandsOverlaysShowDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsOverlaysShowDescription) IsTranslationKey() {
}
func (k KeysCommandsOverlaysShowDescription) GetPath() string {
	return "commands.overlays.show.description"
}
func (k KeysCommandsOverlaysShowDescription) GetPathSlice() []string {
	return []string{"commands", "overlays", "show", "description"}
}
func (k KeysCommandsOverlaysShowDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsOverlaysShowDescription) SetVars(vars KeysCommandsOverlaysShowDescriptionVars) twiri18n.TranslationKey[KeysCommandsOverlaysShowDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsOverlaysShowErrorsOverlayNotFoundVars struct {
}
type KeysCommandsOverlaysShowErrorsOverlayNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsOverlaysShowErrorsOverlayNotFound) IsTranslationKey() {
}
func (k KeysCommandsOverlaysShowErrorsOverlayNotFound) GetPath() string {
	return "commands.overlays.show.errors.overlay_not_found"
}
func (k KeysCommandsOverlaysShowErrorsOverlayNotFound) GetPathSlice() []string {
	return []string{"commands", "overlays", "show", "errors", "overlay_not_found"}
}
func (k KeysCommandsOverlaysShowErrorsOverlayNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsOverlaysShowErrorsOverlayNotFound) SetVars(vars KeysCommandsOverlaysShowErrorsOverlayNotFoundVars) twiri18n.TranslationKey[KeysCommandsOverlaysShowErrorsOverlayNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsOverlaysShowErrorsCannotShowOverlayVars struct {
}
type KeysCommandsOverlaysShowErrorsCannotShowOverlay struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsOverlaysShowErrorsCannotShowOverlay) IsTranslationKey() {
}
func (k KeysCommandsOverlaysShowErrorsCannotShowOverlay) GetPath() string {
	return "commands.overlays.show.errors.cannot_show_overlay"
}
func (k KeysCommandsOverlaysShowErrorsCannotShowOverlay) GetPathSlice() []string {
	return []string{"commands", "overlays", "show", "errors", "cannot_show_overlay"}
}
func (k KeysCommandsOverlaysShowErrorsCannotShowOverlay) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsOverlaysShowErrorsCannotShowOverlay) SetVars(vars KeysCommandsOverlaysShowErrorsCannotShowOverlayVars) twiri18n.TranslationKey[KeysCommandsOverlaysShowErrorsCannotShowOverlayVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsOverlaysShowErrors struct {
	OverlayNotFound		KeysCommandsOverlaysShowErrorsOverlayNotFound
	CannotShowOverlay	KeysCommandsOverlaysShowErrorsCannotShowOverlay
}
type KeysCommandsOverlaysShowSuccessOverlayShownVars struct {
	Overlay any
}
type KeysCommandsOverlaysShowSuccessOverlayShown struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsOverlaysShowSuccessOverlayShown) IsTranslationKey() {
}
func (k KeysCommandsOverlaysShowSuccessOverlayShown) GetPath() string {
	return "commands.overlays.show.success.overlay_shown"
}
func (k KeysCommandsOverlaysShowSuccessOverlayShown) GetPathSlice() []string {
	return []string{"commands", "overlays", "show", "success", "overlay_shown"}
}
func (k KeysCommandsOverlaysShowSuccessOverlayShown) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsOverlaysShowSuccessOverlayShown) SetVars(vars KeysCommandsOverlaysShowSuccessOverlayShownVars) twiri18n.TranslationKey[KeysCommandsOverlaysShowSuccessOverlayShownVars] {
	k.Vars = twiri18n.Vars{"overlay": vars.Overlay}
	return k
}

type KeysCommandsOverlaysShowSuccess struct {
	OverlayShown KeysCommandsOverlaysShowSuccessOverlayShown
}
type KeysCommandsOverlaysShow struct {
	Description	KeysCommandsOverlaysShowDescription
	Errors		KeysCommandsOverlaysShowErrors
	Success		KeysCommandsOverlaysShowSuccess
}
type KeysCommandsOverlays struct {
	Hide		KeysCommandsOverlaysHide
	Description	KeysCommandsOverlaysDescription
	Show		KeysCommandsOverlaysShow
}
type KeysCommandsPermitDescriptionVars struct {
}
type KeysCommandsPermitDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPermitDescription) IsTranslationKey() {
}
func (k KeysCommandsPermitDescription) GetPath() string {
	return "commands.permit.description"
}
func (k KeysCommandsPermitDescription) GetPathSlice() []string {
	return []string{"commands", "permit", "description"}
}
func (k KeysCommandsPermitDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPermitDescription) SetVars(vars KeysCommandsPermitDescriptionVars) twiri18n.TranslationKey[KeysCommandsPermitDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPermitAddDescriptionVars struct {
}
type KeysCommandsPermitAddDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPermitAddDescription) IsTranslationKey() {
}
func (k KeysCommandsPermitAddDescription) GetPath() string {
	return "commands.permit.add.description"
}
func (k KeysCommandsPermitAddDescription) GetPathSlice() []string {
	return []string{"commands", "permit", "add", "description"}
}
func (k KeysCommandsPermitAddDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPermitAddDescription) SetVars(vars KeysCommandsPermitAddDescriptionVars) twiri18n.TranslationKey[KeysCommandsPermitAddDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPermitAddErrorsCannotAddPermitVars struct {
}
type KeysCommandsPermitAddErrorsCannotAddPermit struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPermitAddErrorsCannotAddPermit) IsTranslationKey() {
}
func (k KeysCommandsPermitAddErrorsCannotAddPermit) GetPath() string {
	return "commands.permit.add.errors.cannot_add_permit"
}
func (k KeysCommandsPermitAddErrorsCannotAddPermit) GetPathSlice() []string {
	return []string{"commands", "permit", "add", "errors", "cannot_add_permit"}
}
func (k KeysCommandsPermitAddErrorsCannotAddPermit) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPermitAddErrorsCannotAddPermit) SetVars(vars KeysCommandsPermitAddErrorsCannotAddPermitVars) twiri18n.TranslationKey[KeysCommandsPermitAddErrorsCannotAddPermitVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPermitAddErrorsUserNotFoundVars struct {
}
type KeysCommandsPermitAddErrorsUserNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPermitAddErrorsUserNotFound) IsTranslationKey() {
}
func (k KeysCommandsPermitAddErrorsUserNotFound) GetPath() string {
	return "commands.permit.add.errors.user_not_found"
}
func (k KeysCommandsPermitAddErrorsUserNotFound) GetPathSlice() []string {
	return []string{"commands", "permit", "add", "errors", "user_not_found"}
}
func (k KeysCommandsPermitAddErrorsUserNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPermitAddErrorsUserNotFound) SetVars(vars KeysCommandsPermitAddErrorsUserNotFoundVars) twiri18n.TranslationKey[KeysCommandsPermitAddErrorsUserNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPermitAddErrors struct {
	CannotAddPermit	KeysCommandsPermitAddErrorsCannotAddPermit
	UserNotFound	KeysCommandsPermitAddErrorsUserNotFound
}
type KeysCommandsPermitAddSuccessPermitAddedVars struct {
	User any
}
type KeysCommandsPermitAddSuccessPermitAdded struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPermitAddSuccessPermitAdded) IsTranslationKey() {
}
func (k KeysCommandsPermitAddSuccessPermitAdded) GetPath() string {
	return "commands.permit.add.success.permit_added"
}
func (k KeysCommandsPermitAddSuccessPermitAdded) GetPathSlice() []string {
	return []string{"commands", "permit", "add", "success", "permit_added"}
}
func (k KeysCommandsPermitAddSuccessPermitAdded) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPermitAddSuccessPermitAdded) SetVars(vars KeysCommandsPermitAddSuccessPermitAddedVars) twiri18n.TranslationKey[KeysCommandsPermitAddSuccessPermitAddedVars] {
	k.Vars = twiri18n.Vars{"user": vars.User}
	return k
}

type KeysCommandsPermitAddSuccess struct {
	PermitAdded KeysCommandsPermitAddSuccessPermitAdded
}
type KeysCommandsPermitAdd struct {
	Description	KeysCommandsPermitAddDescription
	Errors		KeysCommandsPermitAddErrors
	Success		KeysCommandsPermitAddSuccess
}
type KeysCommandsPermitRemoveErrorsCannotRemovePermitVars struct {
}
type KeysCommandsPermitRemoveErrorsCannotRemovePermit struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPermitRemoveErrorsCannotRemovePermit) IsTranslationKey() {
}
func (k KeysCommandsPermitRemoveErrorsCannotRemovePermit) GetPath() string {
	return "commands.permit.remove.errors.cannot_remove_permit"
}
func (k KeysCommandsPermitRemoveErrorsCannotRemovePermit) GetPathSlice() []string {
	return []string{"commands", "permit", "remove", "errors", "cannot_remove_permit"}
}
func (k KeysCommandsPermitRemoveErrorsCannotRemovePermit) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPermitRemoveErrorsCannotRemovePermit) SetVars(vars KeysCommandsPermitRemoveErrorsCannotRemovePermitVars) twiri18n.TranslationKey[KeysCommandsPermitRemoveErrorsCannotRemovePermitVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPermitRemoveErrors struct {
	CannotRemovePermit KeysCommandsPermitRemoveErrorsCannotRemovePermit
}
type KeysCommandsPermitRemoveSuccessPermitRemovedVars struct {
	User any
}
type KeysCommandsPermitRemoveSuccessPermitRemoved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPermitRemoveSuccessPermitRemoved) IsTranslationKey() {
}
func (k KeysCommandsPermitRemoveSuccessPermitRemoved) GetPath() string {
	return "commands.permit.remove.success.permit_removed"
}
func (k KeysCommandsPermitRemoveSuccessPermitRemoved) GetPathSlice() []string {
	return []string{"commands", "permit", "remove", "success", "permit_removed"}
}
func (k KeysCommandsPermitRemoveSuccessPermitRemoved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPermitRemoveSuccessPermitRemoved) SetVars(vars KeysCommandsPermitRemoveSuccessPermitRemovedVars) twiri18n.TranslationKey[KeysCommandsPermitRemoveSuccessPermitRemovedVars] {
	k.Vars = twiri18n.Vars{"user": vars.User}
	return k
}

type KeysCommandsPermitRemoveSuccess struct {
	PermitRemoved KeysCommandsPermitRemoveSuccessPermitRemoved
}
type KeysCommandsPermitRemoveDescriptionVars struct {
}
type KeysCommandsPermitRemoveDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPermitRemoveDescription) IsTranslationKey() {
}
func (k KeysCommandsPermitRemoveDescription) GetPath() string {
	return "commands.permit.remove.description"
}
func (k KeysCommandsPermitRemoveDescription) GetPathSlice() []string {
	return []string{"commands", "permit", "remove", "description"}
}
func (k KeysCommandsPermitRemoveDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPermitRemoveDescription) SetVars(vars KeysCommandsPermitRemoveDescriptionVars) twiri18n.TranslationKey[KeysCommandsPermitRemoveDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPermitRemove struct {
	Errors		KeysCommandsPermitRemoveErrors
	Success		KeysCommandsPermitRemoveSuccess
	Description	KeysCommandsPermitRemoveDescription
}
type KeysCommandsPermit struct {
	Description	KeysCommandsPermitDescription
	Add		KeysCommandsPermitAdd
	Remove		KeysCommandsPermitRemove
}
type KeysCommandsGamesDescriptionVars struct {
}
type KeysCommandsGamesDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesDescription) IsTranslationKey() {
}
func (k KeysCommandsGamesDescription) GetPath() string {
	return "commands.games.description"
}
func (k KeysCommandsGamesDescription) GetPathSlice() []string {
	return []string{"commands", "games", "description"}
}
func (k KeysCommandsGamesDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesDescription) SetVars(vars KeysCommandsGamesDescriptionVars) twiri18n.TranslationKey[KeysCommandsGamesDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesCurrentDescriptionVars struct {
}
type KeysCommandsGamesCurrentDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesCurrentDescription) IsTranslationKey() {
}
func (k KeysCommandsGamesCurrentDescription) GetPath() string {
	return "commands.games.current.description"
}
func (k KeysCommandsGamesCurrentDescription) GetPathSlice() []string {
	return []string{"commands", "games", "current", "description"}
}
func (k KeysCommandsGamesCurrentDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesCurrentDescription) SetVars(vars KeysCommandsGamesCurrentDescriptionVars) twiri18n.TranslationKey[KeysCommandsGamesCurrentDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesCurrentErrorsCannotGetGameVars struct {
}
type KeysCommandsGamesCurrentErrorsCannotGetGame struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesCurrentErrorsCannotGetGame) IsTranslationKey() {
}
func (k KeysCommandsGamesCurrentErrorsCannotGetGame) GetPath() string {
	return "commands.games.current.errors.cannot_get_game"
}
func (k KeysCommandsGamesCurrentErrorsCannotGetGame) GetPathSlice() []string {
	return []string{"commands", "games", "current", "errors", "cannot_get_game"}
}
func (k KeysCommandsGamesCurrentErrorsCannotGetGame) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesCurrentErrorsCannotGetGame) SetVars(vars KeysCommandsGamesCurrentErrorsCannotGetGameVars) twiri18n.TranslationKey[KeysCommandsGamesCurrentErrorsCannotGetGameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesCurrentErrorsNoGameSetVars struct {
}
type KeysCommandsGamesCurrentErrorsNoGameSet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesCurrentErrorsNoGameSet) IsTranslationKey() {
}
func (k KeysCommandsGamesCurrentErrorsNoGameSet) GetPath() string {
	return "commands.games.current.errors.no_game_set"
}
func (k KeysCommandsGamesCurrentErrorsNoGameSet) GetPathSlice() []string {
	return []string{"commands", "games", "current", "errors", "no_game_set"}
}
func (k KeysCommandsGamesCurrentErrorsNoGameSet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesCurrentErrorsNoGameSet) SetVars(vars KeysCommandsGamesCurrentErrorsNoGameSetVars) twiri18n.TranslationKey[KeysCommandsGamesCurrentErrorsNoGameSetVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesCurrentErrors struct {
	CannotGetGame	KeysCommandsGamesCurrentErrorsCannotGetGame
	NoGameSet	KeysCommandsGamesCurrentErrorsNoGameSet
}
type KeysCommandsGamesCurrentSuccessGameDisplayVars struct {
	Game any
}
type KeysCommandsGamesCurrentSuccessGameDisplay struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesCurrentSuccessGameDisplay) IsTranslationKey() {
}
func (k KeysCommandsGamesCurrentSuccessGameDisplay) GetPath() string {
	return "commands.games.current.success.game_display"
}
func (k KeysCommandsGamesCurrentSuccessGameDisplay) GetPathSlice() []string {
	return []string{"commands", "games", "current", "success", "game_display"}
}
func (k KeysCommandsGamesCurrentSuccessGameDisplay) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesCurrentSuccessGameDisplay) SetVars(vars KeysCommandsGamesCurrentSuccessGameDisplayVars) twiri18n.TranslationKey[KeysCommandsGamesCurrentSuccessGameDisplayVars] {
	k.Vars = twiri18n.Vars{"game": vars.Game}
	return k
}

type KeysCommandsGamesCurrentSuccess struct {
	GameDisplay KeysCommandsGamesCurrentSuccessGameDisplay
}
type KeysCommandsGamesCurrent struct {
	Description	KeysCommandsGamesCurrentDescription
	Errors		KeysCommandsGamesCurrentErrors
	Success		KeysCommandsGamesCurrentSuccess
}
type KeysCommandsGamesStatsDescriptionVars struct {
}
type KeysCommandsGamesStatsDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesStatsDescription) IsTranslationKey() {
}
func (k KeysCommandsGamesStatsDescription) GetPath() string {
	return "commands.games.stats.description"
}
func (k KeysCommandsGamesStatsDescription) GetPathSlice() []string {
	return []string{"commands", "games", "stats", "description"}
}
func (k KeysCommandsGamesStatsDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesStatsDescription) SetVars(vars KeysCommandsGamesStatsDescriptionVars) twiri18n.TranslationKey[KeysCommandsGamesStatsDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesStatsErrorsCannotGetStatsVars struct {
}
type KeysCommandsGamesStatsErrorsCannotGetStats struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesStatsErrorsCannotGetStats) IsTranslationKey() {
}
func (k KeysCommandsGamesStatsErrorsCannotGetStats) GetPath() string {
	return "commands.games.stats.errors.cannot_get_stats"
}
func (k KeysCommandsGamesStatsErrorsCannotGetStats) GetPathSlice() []string {
	return []string{"commands", "games", "stats", "errors", "cannot_get_stats"}
}
func (k KeysCommandsGamesStatsErrorsCannotGetStats) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesStatsErrorsCannotGetStats) SetVars(vars KeysCommandsGamesStatsErrorsCannotGetStatsVars) twiri18n.TranslationKey[KeysCommandsGamesStatsErrorsCannotGetStatsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesStatsErrors struct {
	CannotGetStats KeysCommandsGamesStatsErrorsCannotGetStats
}
type KeysCommandsGamesStatsSuccessStatsDisplayVars struct {
	Stats any
}
type KeysCommandsGamesStatsSuccessStatsDisplay struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesStatsSuccessStatsDisplay) IsTranslationKey() {
}
func (k KeysCommandsGamesStatsSuccessStatsDisplay) GetPath() string {
	return "commands.games.stats.success.stats_display"
}
func (k KeysCommandsGamesStatsSuccessStatsDisplay) GetPathSlice() []string {
	return []string{"commands", "games", "stats", "success", "stats_display"}
}
func (k KeysCommandsGamesStatsSuccessStatsDisplay) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesStatsSuccessStatsDisplay) SetVars(vars KeysCommandsGamesStatsSuccessStatsDisplayVars) twiri18n.TranslationKey[KeysCommandsGamesStatsSuccessStatsDisplayVars] {
	k.Vars = twiri18n.Vars{"stats": vars.Stats}
	return k
}

type KeysCommandsGamesStatsSuccess struct {
	StatsDisplay KeysCommandsGamesStatsSuccessStatsDisplay
}
type KeysCommandsGamesStats struct {
	Description	KeysCommandsGamesStatsDescription
	Errors		KeysCommandsGamesStatsErrors
	Success		KeysCommandsGamesStatsSuccess
}
type KeysCommandsGames struct {
	Description	KeysCommandsGamesDescription
	Current		KeysCommandsGamesCurrent
	Stats		KeysCommandsGamesStats
}
type KeysCommandsSeventvDescriptionVars struct {
}
type KeysCommandsSeventvDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvDescription) IsTranslationKey() {
}
func (k KeysCommandsSeventvDescription) GetPath() string {
	return "commands.7tv.description"
}
func (k KeysCommandsSeventvDescription) GetPathSlice() []string {
	return []string{"commands", "7tv", "description"}
}
func (k KeysCommandsSeventvDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvDescription) SetVars(vars KeysCommandsSeventvDescriptionVars) twiri18n.TranslationKey[KeysCommandsSeventvDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvAddDescriptionVars struct {
}
type KeysCommandsSeventvAddDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvAddDescription) IsTranslationKey() {
}
func (k KeysCommandsSeventvAddDescription) GetPath() string {
	return "commands.7tv.add.description"
}
func (k KeysCommandsSeventvAddDescription) GetPathSlice() []string {
	return []string{"commands", "7tv", "add", "description"}
}
func (k KeysCommandsSeventvAddDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvAddDescription) SetVars(vars KeysCommandsSeventvAddDescriptionVars) twiri18n.TranslationKey[KeysCommandsSeventvAddDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvAddErrorsCannotAddEmoteVars struct {
}
type KeysCommandsSeventvAddErrorsCannotAddEmote struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvAddErrorsCannotAddEmote) IsTranslationKey() {
}
func (k KeysCommandsSeventvAddErrorsCannotAddEmote) GetPath() string {
	return "commands.7tv.add.errors.cannot_add_emote"
}
func (k KeysCommandsSeventvAddErrorsCannotAddEmote) GetPathSlice() []string {
	return []string{"commands", "7tv", "add", "errors", "cannot_add_emote"}
}
func (k KeysCommandsSeventvAddErrorsCannotAddEmote) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvAddErrorsCannotAddEmote) SetVars(vars KeysCommandsSeventvAddErrorsCannotAddEmoteVars) twiri18n.TranslationKey[KeysCommandsSeventvAddErrorsCannotAddEmoteVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvAddErrorsEmoteNotFoundVars struct {
}
type KeysCommandsSeventvAddErrorsEmoteNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvAddErrorsEmoteNotFound) IsTranslationKey() {
}
func (k KeysCommandsSeventvAddErrorsEmoteNotFound) GetPath() string {
	return "commands.7tv.add.errors.emote_not_found"
}
func (k KeysCommandsSeventvAddErrorsEmoteNotFound) GetPathSlice() []string {
	return []string{"commands", "7tv", "add", "errors", "emote_not_found"}
}
func (k KeysCommandsSeventvAddErrorsEmoteNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvAddErrorsEmoteNotFound) SetVars(vars KeysCommandsSeventvAddErrorsEmoteNotFoundVars) twiri18n.TranslationKey[KeysCommandsSeventvAddErrorsEmoteNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvAddErrorsEmoteAlreadyExistsVars struct {
}
type KeysCommandsSeventvAddErrorsEmoteAlreadyExists struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvAddErrorsEmoteAlreadyExists) IsTranslationKey() {
}
func (k KeysCommandsSeventvAddErrorsEmoteAlreadyExists) GetPath() string {
	return "commands.7tv.add.errors.emote_already_exists"
}
func (k KeysCommandsSeventvAddErrorsEmoteAlreadyExists) GetPathSlice() []string {
	return []string{"commands", "7tv", "add", "errors", "emote_already_exists"}
}
func (k KeysCommandsSeventvAddErrorsEmoteAlreadyExists) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvAddErrorsEmoteAlreadyExists) SetVars(vars KeysCommandsSeventvAddErrorsEmoteAlreadyExistsVars) twiri18n.TranslationKey[KeysCommandsSeventvAddErrorsEmoteAlreadyExistsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvAddErrors struct {
	CannotAddEmote		KeysCommandsSeventvAddErrorsCannotAddEmote
	EmoteNotFound		KeysCommandsSeventvAddErrorsEmoteNotFound
	EmoteAlreadyExists	KeysCommandsSeventvAddErrorsEmoteAlreadyExists
}
type KeysCommandsSeventvAddSuccessEmoteAddedVars struct {
	Emote any
}
type KeysCommandsSeventvAddSuccessEmoteAdded struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvAddSuccessEmoteAdded) IsTranslationKey() {
}
func (k KeysCommandsSeventvAddSuccessEmoteAdded) GetPath() string {
	return "commands.7tv.add.success.emote_added"
}
func (k KeysCommandsSeventvAddSuccessEmoteAdded) GetPathSlice() []string {
	return []string{"commands", "7tv", "add", "success", "emote_added"}
}
func (k KeysCommandsSeventvAddSuccessEmoteAdded) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvAddSuccessEmoteAdded) SetVars(vars KeysCommandsSeventvAddSuccessEmoteAddedVars) twiri18n.TranslationKey[KeysCommandsSeventvAddSuccessEmoteAddedVars] {
	k.Vars = twiri18n.Vars{"emote": vars.Emote}
	return k
}

type KeysCommandsSeventvAddSuccess struct {
	EmoteAdded KeysCommandsSeventvAddSuccessEmoteAdded
}
type KeysCommandsSeventvAdd struct {
	Description	KeysCommandsSeventvAddDescription
	Errors		KeysCommandsSeventvAddErrors
	Success		KeysCommandsSeventvAddSuccess
}
type KeysCommandsSeventvRemoveDescriptionVars struct {
}
type KeysCommandsSeventvRemoveDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvRemoveDescription) IsTranslationKey() {
}
func (k KeysCommandsSeventvRemoveDescription) GetPath() string {
	return "commands.7tv.remove.description"
}
func (k KeysCommandsSeventvRemoveDescription) GetPathSlice() []string {
	return []string{"commands", "7tv", "remove", "description"}
}
func (k KeysCommandsSeventvRemoveDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvRemoveDescription) SetVars(vars KeysCommandsSeventvRemoveDescriptionVars) twiri18n.TranslationKey[KeysCommandsSeventvRemoveDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvRemoveErrorsCannotRemoveEmoteVars struct {
}
type KeysCommandsSeventvRemoveErrorsCannotRemoveEmote struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvRemoveErrorsCannotRemoveEmote) IsTranslationKey() {
}
func (k KeysCommandsSeventvRemoveErrorsCannotRemoveEmote) GetPath() string {
	return "commands.7tv.remove.errors.cannot_remove_emote"
}
func (k KeysCommandsSeventvRemoveErrorsCannotRemoveEmote) GetPathSlice() []string {
	return []string{"commands", "7tv", "remove", "errors", "cannot_remove_emote"}
}
func (k KeysCommandsSeventvRemoveErrorsCannotRemoveEmote) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvRemoveErrorsCannotRemoveEmote) SetVars(vars KeysCommandsSeventvRemoveErrorsCannotRemoveEmoteVars) twiri18n.TranslationKey[KeysCommandsSeventvRemoveErrorsCannotRemoveEmoteVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvRemoveErrorsEmoteNotFoundVars struct {
}
type KeysCommandsSeventvRemoveErrorsEmoteNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvRemoveErrorsEmoteNotFound) IsTranslationKey() {
}
func (k KeysCommandsSeventvRemoveErrorsEmoteNotFound) GetPath() string {
	return "commands.7tv.remove.errors.emote_not_found"
}
func (k KeysCommandsSeventvRemoveErrorsEmoteNotFound) GetPathSlice() []string {
	return []string{"commands", "7tv", "remove", "errors", "emote_not_found"}
}
func (k KeysCommandsSeventvRemoveErrorsEmoteNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvRemoveErrorsEmoteNotFound) SetVars(vars KeysCommandsSeventvRemoveErrorsEmoteNotFoundVars) twiri18n.TranslationKey[KeysCommandsSeventvRemoveErrorsEmoteNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvRemoveErrors struct {
	CannotRemoveEmote	KeysCommandsSeventvRemoveErrorsCannotRemoveEmote
	EmoteNotFound		KeysCommandsSeventvRemoveErrorsEmoteNotFound
}
type KeysCommandsSeventvRemoveSuccessEmoteRemovedVars struct {
	Emote any
}
type KeysCommandsSeventvRemoveSuccessEmoteRemoved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvRemoveSuccessEmoteRemoved) IsTranslationKey() {
}
func (k KeysCommandsSeventvRemoveSuccessEmoteRemoved) GetPath() string {
	return "commands.7tv.remove.success.emote_removed"
}
func (k KeysCommandsSeventvRemoveSuccessEmoteRemoved) GetPathSlice() []string {
	return []string{"commands", "7tv", "remove", "success", "emote_removed"}
}
func (k KeysCommandsSeventvRemoveSuccessEmoteRemoved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvRemoveSuccessEmoteRemoved) SetVars(vars KeysCommandsSeventvRemoveSuccessEmoteRemovedVars) twiri18n.TranslationKey[KeysCommandsSeventvRemoveSuccessEmoteRemovedVars] {
	k.Vars = twiri18n.Vars{"emote": vars.Emote}
	return k
}

type KeysCommandsSeventvRemoveSuccess struct {
	EmoteRemoved KeysCommandsSeventvRemoveSuccessEmoteRemoved
}
type KeysCommandsSeventvRemove struct {
	Description	KeysCommandsSeventvRemoveDescription
	Errors		KeysCommandsSeventvRemoveErrors
	Success		KeysCommandsSeventvRemoveSuccess
}
type KeysCommandsSeventv struct {
	Description	KeysCommandsSeventvDescription
	Add		KeysCommandsSeventvAdd
	Remove		KeysCommandsSeventvRemove
}
type KeysCommandsPrefixDescriptionVars struct {
}
type KeysCommandsPrefixDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixDescription) IsTranslationKey() {
}
func (k KeysCommandsPrefixDescription) GetPath() string {
	return "commands.prefix.description"
}
func (k KeysCommandsPrefixDescription) GetPathSlice() []string {
	return []string{"commands", "prefix", "description"}
}
func (k KeysCommandsPrefixDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixDescription) SetVars(vars KeysCommandsPrefixDescriptionVars) twiri18n.TranslationKey[KeysCommandsPrefixDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixSetDescriptionVars struct {
}
type KeysCommandsPrefixSetDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixSetDescription) IsTranslationKey() {
}
func (k KeysCommandsPrefixSetDescription) GetPath() string {
	return "commands.prefix.set.description"
}
func (k KeysCommandsPrefixSetDescription) GetPathSlice() []string {
	return []string{"commands", "prefix", "set", "description"}
}
func (k KeysCommandsPrefixSetDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixSetDescription) SetVars(vars KeysCommandsPrefixSetDescriptionVars) twiri18n.TranslationKey[KeysCommandsPrefixSetDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixSetErrorsCannotSetPrefixVars struct {
}
type KeysCommandsPrefixSetErrorsCannotSetPrefix struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixSetErrorsCannotSetPrefix) IsTranslationKey() {
}
func (k KeysCommandsPrefixSetErrorsCannotSetPrefix) GetPath() string {
	return "commands.prefix.set.errors.cannot_set_prefix"
}
func (k KeysCommandsPrefixSetErrorsCannotSetPrefix) GetPathSlice() []string {
	return []string{"commands", "prefix", "set", "errors", "cannot_set_prefix"}
}
func (k KeysCommandsPrefixSetErrorsCannotSetPrefix) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixSetErrorsCannotSetPrefix) SetVars(vars KeysCommandsPrefixSetErrorsCannotSetPrefixVars) twiri18n.TranslationKey[KeysCommandsPrefixSetErrorsCannotSetPrefixVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixSetErrorsPrefixTooLongVars struct {
}
type KeysCommandsPrefixSetErrorsPrefixTooLong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixSetErrorsPrefixTooLong) IsTranslationKey() {
}
func (k KeysCommandsPrefixSetErrorsPrefixTooLong) GetPath() string {
	return "commands.prefix.set.errors.prefix_too_long"
}
func (k KeysCommandsPrefixSetErrorsPrefixTooLong) GetPathSlice() []string {
	return []string{"commands", "prefix", "set", "errors", "prefix_too_long"}
}
func (k KeysCommandsPrefixSetErrorsPrefixTooLong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixSetErrorsPrefixTooLong) SetVars(vars KeysCommandsPrefixSetErrorsPrefixTooLongVars) twiri18n.TranslationKey[KeysCommandsPrefixSetErrorsPrefixTooLongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixSetErrorsInvalidPrefixVars struct {
}
type KeysCommandsPrefixSetErrorsInvalidPrefix struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixSetErrorsInvalidPrefix) IsTranslationKey() {
}
func (k KeysCommandsPrefixSetErrorsInvalidPrefix) GetPath() string {
	return "commands.prefix.set.errors.invalid_prefix"
}
func (k KeysCommandsPrefixSetErrorsInvalidPrefix) GetPathSlice() []string {
	return []string{"commands", "prefix", "set", "errors", "invalid_prefix"}
}
func (k KeysCommandsPrefixSetErrorsInvalidPrefix) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixSetErrorsInvalidPrefix) SetVars(vars KeysCommandsPrefixSetErrorsInvalidPrefixVars) twiri18n.TranslationKey[KeysCommandsPrefixSetErrorsInvalidPrefixVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixSetErrors struct {
	CannotSetPrefix	KeysCommandsPrefixSetErrorsCannotSetPrefix
	PrefixTooLong	KeysCommandsPrefixSetErrorsPrefixTooLong
	InvalidPrefix	KeysCommandsPrefixSetErrorsInvalidPrefix
}
type KeysCommandsPrefixSetSuccessPrefixSetVars struct {
	Prefix any
}
type KeysCommandsPrefixSetSuccessPrefixSet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixSetSuccessPrefixSet) IsTranslationKey() {
}
func (k KeysCommandsPrefixSetSuccessPrefixSet) GetPath() string {
	return "commands.prefix.set.success.prefix_set"
}
func (k KeysCommandsPrefixSetSuccessPrefixSet) GetPathSlice() []string {
	return []string{"commands", "prefix", "set", "success", "prefix_set"}
}
func (k KeysCommandsPrefixSetSuccessPrefixSet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixSetSuccessPrefixSet) SetVars(vars KeysCommandsPrefixSetSuccessPrefixSetVars) twiri18n.TranslationKey[KeysCommandsPrefixSetSuccessPrefixSetVars] {
	k.Vars = twiri18n.Vars{"prefix": vars.Prefix}
	return k
}

type KeysCommandsPrefixSetSuccess struct {
	PrefixSet KeysCommandsPrefixSetSuccessPrefixSet
}
type KeysCommandsPrefixSet struct {
	Description	KeysCommandsPrefixSetDescription
	Errors		KeysCommandsPrefixSetErrors
	Success		KeysCommandsPrefixSetSuccess
}
type KeysCommandsPrefixRemoveSuccessPrefixRemovedVars struct {
}
type KeysCommandsPrefixRemoveSuccessPrefixRemoved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixRemoveSuccessPrefixRemoved) IsTranslationKey() {
}
func (k KeysCommandsPrefixRemoveSuccessPrefixRemoved) GetPath() string {
	return "commands.prefix.remove.success.prefix_removed"
}
func (k KeysCommandsPrefixRemoveSuccessPrefixRemoved) GetPathSlice() []string {
	return []string{"commands", "prefix", "remove", "success", "prefix_removed"}
}
func (k KeysCommandsPrefixRemoveSuccessPrefixRemoved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixRemoveSuccessPrefixRemoved) SetVars(vars KeysCommandsPrefixRemoveSuccessPrefixRemovedVars) twiri18n.TranslationKey[KeysCommandsPrefixRemoveSuccessPrefixRemovedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixRemoveSuccess struct {
	PrefixRemoved KeysCommandsPrefixRemoveSuccessPrefixRemoved
}
type KeysCommandsPrefixRemoveDescriptionVars struct {
}
type KeysCommandsPrefixRemoveDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixRemoveDescription) IsTranslationKey() {
}
func (k KeysCommandsPrefixRemoveDescription) GetPath() string {
	return "commands.prefix.remove.description"
}
func (k KeysCommandsPrefixRemoveDescription) GetPathSlice() []string {
	return []string{"commands", "prefix", "remove", "description"}
}
func (k KeysCommandsPrefixRemoveDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixRemoveDescription) SetVars(vars KeysCommandsPrefixRemoveDescriptionVars) twiri18n.TranslationKey[KeysCommandsPrefixRemoveDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixRemoveErrorsCannotRemovePrefixVars struct {
}
type KeysCommandsPrefixRemoveErrorsCannotRemovePrefix struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixRemoveErrorsCannotRemovePrefix) IsTranslationKey() {
}
func (k KeysCommandsPrefixRemoveErrorsCannotRemovePrefix) GetPath() string {
	return "commands.prefix.remove.errors.cannot_remove_prefix"
}
func (k KeysCommandsPrefixRemoveErrorsCannotRemovePrefix) GetPathSlice() []string {
	return []string{"commands", "prefix", "remove", "errors", "cannot_remove_prefix"}
}
func (k KeysCommandsPrefixRemoveErrorsCannotRemovePrefix) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixRemoveErrorsCannotRemovePrefix) SetVars(vars KeysCommandsPrefixRemoveErrorsCannotRemovePrefixVars) twiri18n.TranslationKey[KeysCommandsPrefixRemoveErrorsCannotRemovePrefixVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixRemoveErrors struct {
	CannotRemovePrefix KeysCommandsPrefixRemoveErrorsCannotRemovePrefix
}
type KeysCommandsPrefixRemove struct {
	Success		KeysCommandsPrefixRemoveSuccess
	Description	KeysCommandsPrefixRemoveDescription
	Errors		KeysCommandsPrefixRemoveErrors
}
type KeysCommandsPrefix struct {
	Description	KeysCommandsPrefixDescription
	Set		KeysCommandsPrefixSet
	Remove		KeysCommandsPrefixRemove
}
type KeysCommandsGiveawaysOpenDescriptionVars struct {
}
type KeysCommandsGiveawaysOpenDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysOpenDescription) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysOpenDescription) GetPath() string {
	return "commands.giveaways.open.description"
}
func (k KeysCommandsGiveawaysOpenDescription) GetPathSlice() []string {
	return []string{"commands", "giveaways", "open", "description"}
}
func (k KeysCommandsGiveawaysOpenDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysOpenDescription) SetVars(vars KeysCommandsGiveawaysOpenDescriptionVars) twiri18n.TranslationKey[KeysCommandsGiveawaysOpenDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysOpenSuccessGiveawayOpenedVars struct {
}
type KeysCommandsGiveawaysOpenSuccessGiveawayOpened struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysOpenSuccessGiveawayOpened) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysOpenSuccessGiveawayOpened) GetPath() string {
	return "commands.giveaways.open.success.giveaway_opened"
}
func (k KeysCommandsGiveawaysOpenSuccessGiveawayOpened) GetPathSlice() []string {
	return []string{"commands", "giveaways", "open", "success", "giveaway_opened"}
}
func (k KeysCommandsGiveawaysOpenSuccessGiveawayOpened) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysOpenSuccessGiveawayOpened) SetVars(vars KeysCommandsGiveawaysOpenSuccessGiveawayOpenedVars) twiri18n.TranslationKey[KeysCommandsGiveawaysOpenSuccessGiveawayOpenedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysOpenSuccess struct {
	GiveawayOpened KeysCommandsGiveawaysOpenSuccessGiveawayOpened
}
type KeysCommandsGiveawaysOpenErrorsCannotOpenGiveawayVars struct {
}
type KeysCommandsGiveawaysOpenErrorsCannotOpenGiveaway struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysOpenErrorsCannotOpenGiveaway) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysOpenErrorsCannotOpenGiveaway) GetPath() string {
	return "commands.giveaways.open.errors.cannot_open_giveaway"
}
func (k KeysCommandsGiveawaysOpenErrorsCannotOpenGiveaway) GetPathSlice() []string {
	return []string{"commands", "giveaways", "open", "errors", "cannot_open_giveaway"}
}
func (k KeysCommandsGiveawaysOpenErrorsCannotOpenGiveaway) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysOpenErrorsCannotOpenGiveaway) SetVars(vars KeysCommandsGiveawaysOpenErrorsCannotOpenGiveawayVars) twiri18n.TranslationKey[KeysCommandsGiveawaysOpenErrorsCannotOpenGiveawayVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysOpenErrors struct {
	CannotOpenGiveaway KeysCommandsGiveawaysOpenErrorsCannotOpenGiveaway
}
type KeysCommandsGiveawaysOpen struct {
	Description	KeysCommandsGiveawaysOpenDescription
	Success		KeysCommandsGiveawaysOpenSuccess
	Errors		KeysCommandsGiveawaysOpenErrors
}
type KeysCommandsGiveawaysPauseErrorsCannotPauseGiveawayVars struct {
}
type KeysCommandsGiveawaysPauseErrorsCannotPauseGiveaway struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysPauseErrorsCannotPauseGiveaway) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysPauseErrorsCannotPauseGiveaway) GetPath() string {
	return "commands.giveaways.pause.errors.cannot_pause_giveaway"
}
func (k KeysCommandsGiveawaysPauseErrorsCannotPauseGiveaway) GetPathSlice() []string {
	return []string{"commands", "giveaways", "pause", "errors", "cannot_pause_giveaway"}
}
func (k KeysCommandsGiveawaysPauseErrorsCannotPauseGiveaway) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysPauseErrorsCannotPauseGiveaway) SetVars(vars KeysCommandsGiveawaysPauseErrorsCannotPauseGiveawayVars) twiri18n.TranslationKey[KeysCommandsGiveawaysPauseErrorsCannotPauseGiveawayVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysPauseErrors struct {
	CannotPauseGiveaway KeysCommandsGiveawaysPauseErrorsCannotPauseGiveaway
}
type KeysCommandsGiveawaysPauseDescriptionVars struct {
}
type KeysCommandsGiveawaysPauseDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysPauseDescription) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysPauseDescription) GetPath() string {
	return "commands.giveaways.pause.description"
}
func (k KeysCommandsGiveawaysPauseDescription) GetPathSlice() []string {
	return []string{"commands", "giveaways", "pause", "description"}
}
func (k KeysCommandsGiveawaysPauseDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysPauseDescription) SetVars(vars KeysCommandsGiveawaysPauseDescriptionVars) twiri18n.TranslationKey[KeysCommandsGiveawaysPauseDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysPauseSuccessGiveawayPausedVars struct {
}
type KeysCommandsGiveawaysPauseSuccessGiveawayPaused struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysPauseSuccessGiveawayPaused) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysPauseSuccessGiveawayPaused) GetPath() string {
	return "commands.giveaways.pause.success.giveaway_paused"
}
func (k KeysCommandsGiveawaysPauseSuccessGiveawayPaused) GetPathSlice() []string {
	return []string{"commands", "giveaways", "pause", "success", "giveaway_paused"}
}
func (k KeysCommandsGiveawaysPauseSuccessGiveawayPaused) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysPauseSuccessGiveawayPaused) SetVars(vars KeysCommandsGiveawaysPauseSuccessGiveawayPausedVars) twiri18n.TranslationKey[KeysCommandsGiveawaysPauseSuccessGiveawayPausedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysPauseSuccess struct {
	GiveawayPaused KeysCommandsGiveawaysPauseSuccessGiveawayPaused
}
type KeysCommandsGiveawaysPause struct {
	Errors		KeysCommandsGiveawaysPauseErrors
	Description	KeysCommandsGiveawaysPauseDescription
	Success		KeysCommandsGiveawaysPauseSuccess
}
type KeysCommandsGiveawaysResetErrorsCannotResetGiveawayVars struct {
}
type KeysCommandsGiveawaysResetErrorsCannotResetGiveaway struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysResetErrorsCannotResetGiveaway) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysResetErrorsCannotResetGiveaway) GetPath() string {
	return "commands.giveaways.reset.errors.cannot_reset_giveaway"
}
func (k KeysCommandsGiveawaysResetErrorsCannotResetGiveaway) GetPathSlice() []string {
	return []string{"commands", "giveaways", "reset", "errors", "cannot_reset_giveaway"}
}
func (k KeysCommandsGiveawaysResetErrorsCannotResetGiveaway) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysResetErrorsCannotResetGiveaway) SetVars(vars KeysCommandsGiveawaysResetErrorsCannotResetGiveawayVars) twiri18n.TranslationKey[KeysCommandsGiveawaysResetErrorsCannotResetGiveawayVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysResetErrors struct {
	CannotResetGiveaway KeysCommandsGiveawaysResetErrorsCannotResetGiveaway
}
type KeysCommandsGiveawaysResetDescriptionVars struct {
}
type KeysCommandsGiveawaysResetDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysResetDescription) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysResetDescription) GetPath() string {
	return "commands.giveaways.reset.description"
}
func (k KeysCommandsGiveawaysResetDescription) GetPathSlice() []string {
	return []string{"commands", "giveaways", "reset", "description"}
}
func (k KeysCommandsGiveawaysResetDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysResetDescription) SetVars(vars KeysCommandsGiveawaysResetDescriptionVars) twiri18n.TranslationKey[KeysCommandsGiveawaysResetDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysResetSuccessGiveawayResetVars struct {
}
type KeysCommandsGiveawaysResetSuccessGiveawayReset struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysResetSuccessGiveawayReset) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysResetSuccessGiveawayReset) GetPath() string {
	return "commands.giveaways.reset.success.giveaway_reset"
}
func (k KeysCommandsGiveawaysResetSuccessGiveawayReset) GetPathSlice() []string {
	return []string{"commands", "giveaways", "reset", "success", "giveaway_reset"}
}
func (k KeysCommandsGiveawaysResetSuccessGiveawayReset) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysResetSuccessGiveawayReset) SetVars(vars KeysCommandsGiveawaysResetSuccessGiveawayResetVars) twiri18n.TranslationKey[KeysCommandsGiveawaysResetSuccessGiveawayResetVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysResetSuccess struct {
	GiveawayReset KeysCommandsGiveawaysResetSuccessGiveawayReset
}
type KeysCommandsGiveawaysReset struct {
	Errors		KeysCommandsGiveawaysResetErrors
	Description	KeysCommandsGiveawaysResetDescription
	Success		KeysCommandsGiveawaysResetSuccess
}
type KeysCommandsGiveawaysFinishDescriptionVars struct {
}
type KeysCommandsGiveawaysFinishDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysFinishDescription) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysFinishDescription) GetPath() string {
	return "commands.giveaways.finish.description"
}
func (k KeysCommandsGiveawaysFinishDescription) GetPathSlice() []string {
	return []string{"commands", "giveaways", "finish", "description"}
}
func (k KeysCommandsGiveawaysFinishDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysFinishDescription) SetVars(vars KeysCommandsGiveawaysFinishDescriptionVars) twiri18n.TranslationKey[KeysCommandsGiveawaysFinishDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysFinishSuccessGiveawayFinishedVars struct {
}
type KeysCommandsGiveawaysFinishSuccessGiveawayFinished struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysFinishSuccessGiveawayFinished) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysFinishSuccessGiveawayFinished) GetPath() string {
	return "commands.giveaways.finish.success.giveaway_finished"
}
func (k KeysCommandsGiveawaysFinishSuccessGiveawayFinished) GetPathSlice() []string {
	return []string{"commands", "giveaways", "finish", "success", "giveaway_finished"}
}
func (k KeysCommandsGiveawaysFinishSuccessGiveawayFinished) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysFinishSuccessGiveawayFinished) SetVars(vars KeysCommandsGiveawaysFinishSuccessGiveawayFinishedVars) twiri18n.TranslationKey[KeysCommandsGiveawaysFinishSuccessGiveawayFinishedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysFinishSuccess struct {
	GiveawayFinished KeysCommandsGiveawaysFinishSuccessGiveawayFinished
}
type KeysCommandsGiveawaysFinishErrorsCannotFinishGiveawayVars struct {
}
type KeysCommandsGiveawaysFinishErrorsCannotFinishGiveaway struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysFinishErrorsCannotFinishGiveaway) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysFinishErrorsCannotFinishGiveaway) GetPath() string {
	return "commands.giveaways.finish.errors.cannot_finish_giveaway"
}
func (k KeysCommandsGiveawaysFinishErrorsCannotFinishGiveaway) GetPathSlice() []string {
	return []string{"commands", "giveaways", "finish", "errors", "cannot_finish_giveaway"}
}
func (k KeysCommandsGiveawaysFinishErrorsCannotFinishGiveaway) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysFinishErrorsCannotFinishGiveaway) SetVars(vars KeysCommandsGiveawaysFinishErrorsCannotFinishGiveawayVars) twiri18n.TranslationKey[KeysCommandsGiveawaysFinishErrorsCannotFinishGiveawayVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysFinishErrors struct {
	CannotFinishGiveaway KeysCommandsGiveawaysFinishErrorsCannotFinishGiveaway
}
type KeysCommandsGiveawaysFinish struct {
	Description	KeysCommandsGiveawaysFinishDescription
	Success		KeysCommandsGiveawaysFinishSuccess
	Errors		KeysCommandsGiveawaysFinishErrors
}
type KeysCommandsGiveawaysDeleteDescriptionVars struct {
}
type KeysCommandsGiveawaysDeleteDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysDeleteDescription) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysDeleteDescription) GetPath() string {
	return "commands.giveaways.delete.description"
}
func (k KeysCommandsGiveawaysDeleteDescription) GetPathSlice() []string {
	return []string{"commands", "giveaways", "delete", "description"}
}
func (k KeysCommandsGiveawaysDeleteDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysDeleteDescription) SetVars(vars KeysCommandsGiveawaysDeleteDescriptionVars) twiri18n.TranslationKey[KeysCommandsGiveawaysDeleteDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysDeleteSuccessGiveawayDeletedVars struct {
}
type KeysCommandsGiveawaysDeleteSuccessGiveawayDeleted struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysDeleteSuccessGiveawayDeleted) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysDeleteSuccessGiveawayDeleted) GetPath() string {
	return "commands.giveaways.delete.success.giveaway_deleted"
}
func (k KeysCommandsGiveawaysDeleteSuccessGiveawayDeleted) GetPathSlice() []string {
	return []string{"commands", "giveaways", "delete", "success", "giveaway_deleted"}
}
func (k KeysCommandsGiveawaysDeleteSuccessGiveawayDeleted) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysDeleteSuccessGiveawayDeleted) SetVars(vars KeysCommandsGiveawaysDeleteSuccessGiveawayDeletedVars) twiri18n.TranslationKey[KeysCommandsGiveawaysDeleteSuccessGiveawayDeletedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysDeleteSuccess struct {
	GiveawayDeleted KeysCommandsGiveawaysDeleteSuccessGiveawayDeleted
}
type KeysCommandsGiveawaysDeleteErrorsCannotDeleteGiveawayVars struct {
}
type KeysCommandsGiveawaysDeleteErrorsCannotDeleteGiveaway struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGiveawaysDeleteErrorsCannotDeleteGiveaway) IsTranslationKey() {
}
func (k KeysCommandsGiveawaysDeleteErrorsCannotDeleteGiveaway) GetPath() string {
	return "commands.giveaways.delete.errors.cannot_delete_giveaway"
}
func (k KeysCommandsGiveawaysDeleteErrorsCannotDeleteGiveaway) GetPathSlice() []string {
	return []string{"commands", "giveaways", "delete", "errors", "cannot_delete_giveaway"}
}
func (k KeysCommandsGiveawaysDeleteErrorsCannotDeleteGiveaway) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGiveawaysDeleteErrorsCannotDeleteGiveaway) SetVars(vars KeysCommandsGiveawaysDeleteErrorsCannotDeleteGiveawayVars) twiri18n.TranslationKey[KeysCommandsGiveawaysDeleteErrorsCannotDeleteGiveawayVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGiveawaysDeleteErrors struct {
	CannotDeleteGiveaway KeysCommandsGiveawaysDeleteErrorsCannotDeleteGiveaway
}
type KeysCommandsGiveawaysDelete struct {
	Description	KeysCommandsGiveawaysDeleteDescription
	Success		KeysCommandsGiveawaysDeleteSuccess
	Errors		KeysCommandsGiveawaysDeleteErrors
}
type KeysCommandsGiveaways struct {
	Open	KeysCommandsGiveawaysOpen
	Pause	KeysCommandsGiveawaysPause
	Reset	KeysCommandsGiveawaysReset
	Finish	KeysCommandsGiveawaysFinish
	Delete	KeysCommandsGiveawaysDelete
}
type KeysCommandsSongrequestAddDescriptionVars struct {
}
type KeysCommandsSongrequestAddDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestAddDescription) IsTranslationKey() {
}
func (k KeysCommandsSongrequestAddDescription) GetPath() string {
	return "commands.songrequest.add.description"
}
func (k KeysCommandsSongrequestAddDescription) GetPathSlice() []string {
	return []string{"commands", "songrequest", "add", "description"}
}
func (k KeysCommandsSongrequestAddDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestAddDescription) SetVars(vars KeysCommandsSongrequestAddDescriptionVars) twiri18n.TranslationKey[KeysCommandsSongrequestAddDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestAddErrorsSongAlreadyInQueueVars struct {
}
type KeysCommandsSongrequestAddErrorsSongAlreadyInQueue struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestAddErrorsSongAlreadyInQueue) IsTranslationKey() {
}
func (k KeysCommandsSongrequestAddErrorsSongAlreadyInQueue) GetPath() string {
	return "commands.songrequest.add.errors.song_already_in_queue"
}
func (k KeysCommandsSongrequestAddErrorsSongAlreadyInQueue) GetPathSlice() []string {
	return []string{"commands", "songrequest", "add", "errors", "song_already_in_queue"}
}
func (k KeysCommandsSongrequestAddErrorsSongAlreadyInQueue) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestAddErrorsSongAlreadyInQueue) SetVars(vars KeysCommandsSongrequestAddErrorsSongAlreadyInQueueVars) twiri18n.TranslationKey[KeysCommandsSongrequestAddErrorsSongAlreadyInQueueVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestAddErrorsCannotAddSongVars struct {
}
type KeysCommandsSongrequestAddErrorsCannotAddSong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestAddErrorsCannotAddSong) IsTranslationKey() {
}
func (k KeysCommandsSongrequestAddErrorsCannotAddSong) GetPath() string {
	return "commands.songrequest.add.errors.cannot_add_song"
}
func (k KeysCommandsSongrequestAddErrorsCannotAddSong) GetPathSlice() []string {
	return []string{"commands", "songrequest", "add", "errors", "cannot_add_song"}
}
func (k KeysCommandsSongrequestAddErrorsCannotAddSong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestAddErrorsCannotAddSong) SetVars(vars KeysCommandsSongrequestAddErrorsCannotAddSongVars) twiri18n.TranslationKey[KeysCommandsSongrequestAddErrorsCannotAddSongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestAddErrorsSongNotFoundVars struct {
}
type KeysCommandsSongrequestAddErrorsSongNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestAddErrorsSongNotFound) IsTranslationKey() {
}
func (k KeysCommandsSongrequestAddErrorsSongNotFound) GetPath() string {
	return "commands.songrequest.add.errors.song_not_found"
}
func (k KeysCommandsSongrequestAddErrorsSongNotFound) GetPathSlice() []string {
	return []string{"commands", "songrequest", "add", "errors", "song_not_found"}
}
func (k KeysCommandsSongrequestAddErrorsSongNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestAddErrorsSongNotFound) SetVars(vars KeysCommandsSongrequestAddErrorsSongNotFoundVars) twiri18n.TranslationKey[KeysCommandsSongrequestAddErrorsSongNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestAddErrorsQueueFullVars struct {
}
type KeysCommandsSongrequestAddErrorsQueueFull struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestAddErrorsQueueFull) IsTranslationKey() {
}
func (k KeysCommandsSongrequestAddErrorsQueueFull) GetPath() string {
	return "commands.songrequest.add.errors.queue_full"
}
func (k KeysCommandsSongrequestAddErrorsQueueFull) GetPathSlice() []string {
	return []string{"commands", "songrequest", "add", "errors", "queue_full"}
}
func (k KeysCommandsSongrequestAddErrorsQueueFull) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestAddErrorsQueueFull) SetVars(vars KeysCommandsSongrequestAddErrorsQueueFullVars) twiri18n.TranslationKey[KeysCommandsSongrequestAddErrorsQueueFullVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestAddErrors struct {
	SongAlreadyInQueue	KeysCommandsSongrequestAddErrorsSongAlreadyInQueue
	CannotAddSong		KeysCommandsSongrequestAddErrorsCannotAddSong
	SongNotFound		KeysCommandsSongrequestAddErrorsSongNotFound
	QueueFull		KeysCommandsSongrequestAddErrorsQueueFull
}
type KeysCommandsSongrequestAddSuccessSongAddedVars struct {
	Song	any
	Artist	any
}
type KeysCommandsSongrequestAddSuccessSongAdded struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestAddSuccessSongAdded) IsTranslationKey() {
}
func (k KeysCommandsSongrequestAddSuccessSongAdded) GetPath() string {
	return "commands.songrequest.add.success.song_added"
}
func (k KeysCommandsSongrequestAddSuccessSongAdded) GetPathSlice() []string {
	return []string{"commands", "songrequest", "add", "success", "song_added"}
}
func (k KeysCommandsSongrequestAddSuccessSongAdded) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestAddSuccessSongAdded) SetVars(vars KeysCommandsSongrequestAddSuccessSongAddedVars) twiri18n.TranslationKey[KeysCommandsSongrequestAddSuccessSongAddedVars] {
	k.Vars = twiri18n.Vars{"song": vars.Song, "artist": vars.Artist}
	return k
}

type KeysCommandsSongrequestAddSuccess struct {
	SongAdded KeysCommandsSongrequestAddSuccessSongAdded
}
type KeysCommandsSongrequestAdd struct {
	Description	KeysCommandsSongrequestAddDescription
	Errors		KeysCommandsSongrequestAddErrors
	Success		KeysCommandsSongrequestAddSuccess
}
type KeysCommandsSongrequestRemoveDescriptionVars struct {
}
type KeysCommandsSongrequestRemoveDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestRemoveDescription) IsTranslationKey() {
}
func (k KeysCommandsSongrequestRemoveDescription) GetPath() string {
	return "commands.songrequest.remove.description"
}
func (k KeysCommandsSongrequestRemoveDescription) GetPathSlice() []string {
	return []string{"commands", "songrequest", "remove", "description"}
}
func (k KeysCommandsSongrequestRemoveDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestRemoveDescription) SetVars(vars KeysCommandsSongrequestRemoveDescriptionVars) twiri18n.TranslationKey[KeysCommandsSongrequestRemoveDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestRemoveErrorsCannotRemoveSongVars struct {
}
type KeysCommandsSongrequestRemoveErrorsCannotRemoveSong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestRemoveErrorsCannotRemoveSong) IsTranslationKey() {
}
func (k KeysCommandsSongrequestRemoveErrorsCannotRemoveSong) GetPath() string {
	return "commands.songrequest.remove.errors.cannot_remove_song"
}
func (k KeysCommandsSongrequestRemoveErrorsCannotRemoveSong) GetPathSlice() []string {
	return []string{"commands", "songrequest", "remove", "errors", "cannot_remove_song"}
}
func (k KeysCommandsSongrequestRemoveErrorsCannotRemoveSong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestRemoveErrorsCannotRemoveSong) SetVars(vars KeysCommandsSongrequestRemoveErrorsCannotRemoveSongVars) twiri18n.TranslationKey[KeysCommandsSongrequestRemoveErrorsCannotRemoveSongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestRemoveErrorsSongNotInQueueVars struct {
}
type KeysCommandsSongrequestRemoveErrorsSongNotInQueue struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestRemoveErrorsSongNotInQueue) IsTranslationKey() {
}
func (k KeysCommandsSongrequestRemoveErrorsSongNotInQueue) GetPath() string {
	return "commands.songrequest.remove.errors.song_not_in_queue"
}
func (k KeysCommandsSongrequestRemoveErrorsSongNotInQueue) GetPathSlice() []string {
	return []string{"commands", "songrequest", "remove", "errors", "song_not_in_queue"}
}
func (k KeysCommandsSongrequestRemoveErrorsSongNotInQueue) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestRemoveErrorsSongNotInQueue) SetVars(vars KeysCommandsSongrequestRemoveErrorsSongNotInQueueVars) twiri18n.TranslationKey[KeysCommandsSongrequestRemoveErrorsSongNotInQueueVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestRemoveErrors struct {
	CannotRemoveSong	KeysCommandsSongrequestRemoveErrorsCannotRemoveSong
	SongNotInQueue		KeysCommandsSongrequestRemoveErrorsSongNotInQueue
}
type KeysCommandsSongrequestRemoveSuccessSongRemovedVars struct {
	Song any
}
type KeysCommandsSongrequestRemoveSuccessSongRemoved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestRemoveSuccessSongRemoved) IsTranslationKey() {
}
func (k KeysCommandsSongrequestRemoveSuccessSongRemoved) GetPath() string {
	return "commands.songrequest.remove.success.song_removed"
}
func (k KeysCommandsSongrequestRemoveSuccessSongRemoved) GetPathSlice() []string {
	return []string{"commands", "songrequest", "remove", "success", "song_removed"}
}
func (k KeysCommandsSongrequestRemoveSuccessSongRemoved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestRemoveSuccessSongRemoved) SetVars(vars KeysCommandsSongrequestRemoveSuccessSongRemovedVars) twiri18n.TranslationKey[KeysCommandsSongrequestRemoveSuccessSongRemovedVars] {
	k.Vars = twiri18n.Vars{"song": vars.Song}
	return k
}

type KeysCommandsSongrequestRemoveSuccess struct {
	SongRemoved KeysCommandsSongrequestRemoveSuccessSongRemoved
}
type KeysCommandsSongrequestRemove struct {
	Description	KeysCommandsSongrequestRemoveDescription
	Errors		KeysCommandsSongrequestRemoveErrors
	Success		KeysCommandsSongrequestRemoveSuccess
}
type KeysCommandsSongrequestQueueDescriptionVars struct {
}
type KeysCommandsSongrequestQueueDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestQueueDescription) IsTranslationKey() {
}
func (k KeysCommandsSongrequestQueueDescription) GetPath() string {
	return "commands.songrequest.queue.description"
}
func (k KeysCommandsSongrequestQueueDescription) GetPathSlice() []string {
	return []string{"commands", "songrequest", "queue", "description"}
}
func (k KeysCommandsSongrequestQueueDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestQueueDescription) SetVars(vars KeysCommandsSongrequestQueueDescriptionVars) twiri18n.TranslationKey[KeysCommandsSongrequestQueueDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestQueueErrorsCannotGetQueueVars struct {
}
type KeysCommandsSongrequestQueueErrorsCannotGetQueue struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestQueueErrorsCannotGetQueue) IsTranslationKey() {
}
func (k KeysCommandsSongrequestQueueErrorsCannotGetQueue) GetPath() string {
	return "commands.songrequest.queue.errors.cannot_get_queue"
}
func (k KeysCommandsSongrequestQueueErrorsCannotGetQueue) GetPathSlice() []string {
	return []string{"commands", "songrequest", "queue", "errors", "cannot_get_queue"}
}
func (k KeysCommandsSongrequestQueueErrorsCannotGetQueue) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestQueueErrorsCannotGetQueue) SetVars(vars KeysCommandsSongrequestQueueErrorsCannotGetQueueVars) twiri18n.TranslationKey[KeysCommandsSongrequestQueueErrorsCannotGetQueueVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestQueueErrors struct {
	CannotGetQueue KeysCommandsSongrequestQueueErrorsCannotGetQueue
}
type KeysCommandsSongrequestQueueSuccessQueueDisplayVars struct {
	Songs any
}
type KeysCommandsSongrequestQueueSuccessQueueDisplay struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestQueueSuccessQueueDisplay) IsTranslationKey() {
}
func (k KeysCommandsSongrequestQueueSuccessQueueDisplay) GetPath() string {
	return "commands.songrequest.queue.success.queue_display"
}
func (k KeysCommandsSongrequestQueueSuccessQueueDisplay) GetPathSlice() []string {
	return []string{"commands", "songrequest", "queue", "success", "queue_display"}
}
func (k KeysCommandsSongrequestQueueSuccessQueueDisplay) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestQueueSuccessQueueDisplay) SetVars(vars KeysCommandsSongrequestQueueSuccessQueueDisplayVars) twiri18n.TranslationKey[KeysCommandsSongrequestQueueSuccessQueueDisplayVars] {
	k.Vars = twiri18n.Vars{"songs": vars.Songs}
	return k
}

type KeysCommandsSongrequestQueueSuccessEmptyQueueVars struct {
}
type KeysCommandsSongrequestQueueSuccessEmptyQueue struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestQueueSuccessEmptyQueue) IsTranslationKey() {
}
func (k KeysCommandsSongrequestQueueSuccessEmptyQueue) GetPath() string {
	return "commands.songrequest.queue.success.empty_queue"
}
func (k KeysCommandsSongrequestQueueSuccessEmptyQueue) GetPathSlice() []string {
	return []string{"commands", "songrequest", "queue", "success", "empty_queue"}
}
func (k KeysCommandsSongrequestQueueSuccessEmptyQueue) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestQueueSuccessEmptyQueue) SetVars(vars KeysCommandsSongrequestQueueSuccessEmptyQueueVars) twiri18n.TranslationKey[KeysCommandsSongrequestQueueSuccessEmptyQueueVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestQueueSuccess struct {
	QueueDisplay	KeysCommandsSongrequestQueueSuccessQueueDisplay
	EmptyQueue	KeysCommandsSongrequestQueueSuccessEmptyQueue
}
type KeysCommandsSongrequestQueue struct {
	Description	KeysCommandsSongrequestQueueDescription
	Errors		KeysCommandsSongrequestQueueErrors
	Success		KeysCommandsSongrequestQueueSuccess
}
type KeysCommandsSongrequestDescriptionVars struct {
}
type KeysCommandsSongrequestDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestDescription) IsTranslationKey() {
}
func (k KeysCommandsSongrequestDescription) GetPath() string {
	return "commands.songrequest.description"
}
func (k KeysCommandsSongrequestDescription) GetPathSlice() []string {
	return []string{"commands", "songrequest", "description"}
}
func (k KeysCommandsSongrequestDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestDescription) SetVars(vars KeysCommandsSongrequestDescriptionVars) twiri18n.TranslationKey[KeysCommandsSongrequestDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequest struct {
	Add		KeysCommandsSongrequestAdd
	Remove		KeysCommandsSongrequestRemove
	Queue		KeysCommandsSongrequestQueue
	Description	KeysCommandsSongrequestDescription
}
type KeysCommandsUtilityDescriptionVars struct {
}
type KeysCommandsUtilityDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsUtilityDescription) IsTranslationKey() {
}
func (k KeysCommandsUtilityDescription) GetPath() string {
	return "commands.utility.description"
}
func (k KeysCommandsUtilityDescription) GetPathSlice() []string {
	return []string{"commands", "utility", "description"}
}
func (k KeysCommandsUtilityDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsUtilityDescription) SetVars(vars KeysCommandsUtilityDescriptionVars) twiri18n.TranslationKey[KeysCommandsUtilityDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsUtilityErrorsCannotGetFollowersVars struct {
}
type KeysCommandsUtilityErrorsCannotGetFollowers struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsUtilityErrorsCannotGetFollowers) IsTranslationKey() {
}
func (k KeysCommandsUtilityErrorsCannotGetFollowers) GetPath() string {
	return "commands.utility.errors.cannot_get_followers"
}
func (k KeysCommandsUtilityErrorsCannotGetFollowers) GetPathSlice() []string {
	return []string{"commands", "utility", "errors", "cannot_get_followers"}
}
func (k KeysCommandsUtilityErrorsCannotGetFollowers) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsUtilityErrorsCannotGetFollowers) SetVars(vars KeysCommandsUtilityErrorsCannotGetFollowersVars) twiri18n.TranslationKey[KeysCommandsUtilityErrorsCannotGetFollowersVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsUtilityErrors struct {
	CannotGetFollowers KeysCommandsUtilityErrorsCannotGetFollowers
}
type KeysCommandsUtilitySuccessFirstFollowersVars struct {
	Followers any
}
type KeysCommandsUtilitySuccessFirstFollowers struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsUtilitySuccessFirstFollowers) IsTranslationKey() {
}
func (k KeysCommandsUtilitySuccessFirstFollowers) GetPath() string {
	return "commands.utility.success.first_followers"
}
func (k KeysCommandsUtilitySuccessFirstFollowers) GetPathSlice() []string {
	return []string{"commands", "utility", "success", "first_followers"}
}
func (k KeysCommandsUtilitySuccessFirstFollowers) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsUtilitySuccessFirstFollowers) SetVars(vars KeysCommandsUtilitySuccessFirstFollowersVars) twiri18n.TranslationKey[KeysCommandsUtilitySuccessFirstFollowersVars] {
	k.Vars = twiri18n.Vars{"followers": vars.Followers}
	return k
}

type KeysCommandsUtilitySuccess struct {
	FirstFollowers KeysCommandsUtilitySuccessFirstFollowers
}
type KeysCommandsUtility struct {
	Description	KeysCommandsUtilityDescription
	Errors		KeysCommandsUtilityErrors
	Success		KeysCommandsUtilitySuccess
}
type KeysCommandsPredictionsStartDescriptionVars struct {
}
type KeysCommandsPredictionsStartDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsStartDescription) IsTranslationKey() {
}
func (k KeysCommandsPredictionsStartDescription) GetPath() string {
	return "commands.predictions.start.description"
}
func (k KeysCommandsPredictionsStartDescription) GetPathSlice() []string {
	return []string{"commands", "predictions", "start", "description"}
}
func (k KeysCommandsPredictionsStartDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsStartDescription) SetVars(vars KeysCommandsPredictionsStartDescriptionVars) twiri18n.TranslationKey[KeysCommandsPredictionsStartDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsStartErrorsCannotCreateTwitchClientVars struct {
}
type KeysCommandsPredictionsStartErrorsCannotCreateTwitchClient struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsStartErrorsCannotCreateTwitchClient) IsTranslationKey() {
}
func (k KeysCommandsPredictionsStartErrorsCannotCreateTwitchClient) GetPath() string {
	return "commands.predictions.start.errors.cannot_create_twitch_client"
}
func (k KeysCommandsPredictionsStartErrorsCannotCreateTwitchClient) GetPathSlice() []string {
	return []string{"commands", "predictions", "start", "errors", "cannot_create_twitch_client"}
}
func (k KeysCommandsPredictionsStartErrorsCannotCreateTwitchClient) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsStartErrorsCannotCreateTwitchClient) SetVars(vars KeysCommandsPredictionsStartErrorsCannotCreateTwitchClientVars) twiri18n.TranslationKey[KeysCommandsPredictionsStartErrorsCannotCreateTwitchClientVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsStartErrorsCannotCreatePredictionVars struct {
}
type KeysCommandsPredictionsStartErrorsCannotCreatePrediction struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsStartErrorsCannotCreatePrediction) IsTranslationKey() {
}
func (k KeysCommandsPredictionsStartErrorsCannotCreatePrediction) GetPath() string {
	return "commands.predictions.start.errors.cannot_create_prediction"
}
func (k KeysCommandsPredictionsStartErrorsCannotCreatePrediction) GetPathSlice() []string {
	return []string{"commands", "predictions", "start", "errors", "cannot_create_prediction"}
}
func (k KeysCommandsPredictionsStartErrorsCannotCreatePrediction) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsStartErrorsCannotCreatePrediction) SetVars(vars KeysCommandsPredictionsStartErrorsCannotCreatePredictionVars) twiri18n.TranslationKey[KeysCommandsPredictionsStartErrorsCannotCreatePredictionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsStartErrors struct {
	CannotCreateTwitchClient	KeysCommandsPredictionsStartErrorsCannotCreateTwitchClient
	CannotCreatePrediction		KeysCommandsPredictionsStartErrorsCannotCreatePrediction
}
type KeysCommandsPredictionsStartSuccessPredictionCreatedVars struct {
	Title any
}
type KeysCommandsPredictionsStartSuccessPredictionCreated struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsStartSuccessPredictionCreated) IsTranslationKey() {
}
func (k KeysCommandsPredictionsStartSuccessPredictionCreated) GetPath() string {
	return "commands.predictions.start.success.prediction_created"
}
func (k KeysCommandsPredictionsStartSuccessPredictionCreated) GetPathSlice() []string {
	return []string{"commands", "predictions", "start", "success", "prediction_created"}
}
func (k KeysCommandsPredictionsStartSuccessPredictionCreated) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsStartSuccessPredictionCreated) SetVars(vars KeysCommandsPredictionsStartSuccessPredictionCreatedVars) twiri18n.TranslationKey[KeysCommandsPredictionsStartSuccessPredictionCreatedVars] {
	k.Vars = twiri18n.Vars{"title": vars.Title}
	return k
}

type KeysCommandsPredictionsStartSuccess struct {
	PredictionCreated KeysCommandsPredictionsStartSuccessPredictionCreated
}
type KeysCommandsPredictionsStart struct {
	Description	KeysCommandsPredictionsStartDescription
	Errors		KeysCommandsPredictionsStartErrors
	Success		KeysCommandsPredictionsStartSuccess
}
type KeysCommandsPredictionsResolveErrorsCannotCreateTwitchClientVars struct {
}
type KeysCommandsPredictionsResolveErrorsCannotCreateTwitchClient struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsResolveErrorsCannotCreateTwitchClient) IsTranslationKey() {
}
func (k KeysCommandsPredictionsResolveErrorsCannotCreateTwitchClient) GetPath() string {
	return "commands.predictions.resolve.errors.cannot_create_twitch_client"
}
func (k KeysCommandsPredictionsResolveErrorsCannotCreateTwitchClient) GetPathSlice() []string {
	return []string{"commands", "predictions", "resolve", "errors", "cannot_create_twitch_client"}
}
func (k KeysCommandsPredictionsResolveErrorsCannotCreateTwitchClient) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsResolveErrorsCannotCreateTwitchClient) SetVars(vars KeysCommandsPredictionsResolveErrorsCannotCreateTwitchClientVars) twiri18n.TranslationKey[KeysCommandsPredictionsResolveErrorsCannotCreateTwitchClientVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsResolveErrorsCannotGetPredictionVars struct {
}
type KeysCommandsPredictionsResolveErrorsCannotGetPrediction struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsResolveErrorsCannotGetPrediction) IsTranslationKey() {
}
func (k KeysCommandsPredictionsResolveErrorsCannotGetPrediction) GetPath() string {
	return "commands.predictions.resolve.errors.cannot_get_prediction"
}
func (k KeysCommandsPredictionsResolveErrorsCannotGetPrediction) GetPathSlice() []string {
	return []string{"commands", "predictions", "resolve", "errors", "cannot_get_prediction"}
}
func (k KeysCommandsPredictionsResolveErrorsCannotGetPrediction) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsResolveErrorsCannotGetPrediction) SetVars(vars KeysCommandsPredictionsResolveErrorsCannotGetPredictionVars) twiri18n.TranslationKey[KeysCommandsPredictionsResolveErrorsCannotGetPredictionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsResolveErrorsNoPredictionRunningVars struct {
}
type KeysCommandsPredictionsResolveErrorsNoPredictionRunning struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsResolveErrorsNoPredictionRunning) IsTranslationKey() {
}
func (k KeysCommandsPredictionsResolveErrorsNoPredictionRunning) GetPath() string {
	return "commands.predictions.resolve.errors.no_prediction_running"
}
func (k KeysCommandsPredictionsResolveErrorsNoPredictionRunning) GetPathSlice() []string {
	return []string{"commands", "predictions", "resolve", "errors", "no_prediction_running"}
}
func (k KeysCommandsPredictionsResolveErrorsNoPredictionRunning) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsResolveErrorsNoPredictionRunning) SetVars(vars KeysCommandsPredictionsResolveErrorsNoPredictionRunningVars) twiri18n.TranslationKey[KeysCommandsPredictionsResolveErrorsNoPredictionRunningVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsResolveErrorsNoPredictionVariantVars struct {
}
type KeysCommandsPredictionsResolveErrorsNoPredictionVariant struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsResolveErrorsNoPredictionVariant) IsTranslationKey() {
}
func (k KeysCommandsPredictionsResolveErrorsNoPredictionVariant) GetPath() string {
	return "commands.predictions.resolve.errors.no_prediction_variant"
}
func (k KeysCommandsPredictionsResolveErrorsNoPredictionVariant) GetPathSlice() []string {
	return []string{"commands", "predictions", "resolve", "errors", "no_prediction_variant"}
}
func (k KeysCommandsPredictionsResolveErrorsNoPredictionVariant) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsResolveErrorsNoPredictionVariant) SetVars(vars KeysCommandsPredictionsResolveErrorsNoPredictionVariantVars) twiri18n.TranslationKey[KeysCommandsPredictionsResolveErrorsNoPredictionVariantVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsResolveErrorsCannotResolvePredictionVars struct {
}
type KeysCommandsPredictionsResolveErrorsCannotResolvePrediction struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsResolveErrorsCannotResolvePrediction) IsTranslationKey() {
}
func (k KeysCommandsPredictionsResolveErrorsCannotResolvePrediction) GetPath() string {
	return "commands.predictions.resolve.errors.cannot_resolve_prediction"
}
func (k KeysCommandsPredictionsResolveErrorsCannotResolvePrediction) GetPathSlice() []string {
	return []string{"commands", "predictions", "resolve", "errors", "cannot_resolve_prediction"}
}
func (k KeysCommandsPredictionsResolveErrorsCannotResolvePrediction) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsResolveErrorsCannotResolvePrediction) SetVars(vars KeysCommandsPredictionsResolveErrorsCannotResolvePredictionVars) twiri18n.TranslationKey[KeysCommandsPredictionsResolveErrorsCannotResolvePredictionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsResolveErrors struct {
	CannotCreateTwitchClient	KeysCommandsPredictionsResolveErrorsCannotCreateTwitchClient
	CannotGetPrediction		KeysCommandsPredictionsResolveErrorsCannotGetPrediction
	NoPredictionRunning		KeysCommandsPredictionsResolveErrorsNoPredictionRunning
	NoPredictionVariant		KeysCommandsPredictionsResolveErrorsNoPredictionVariant
	CannotResolvePrediction		KeysCommandsPredictionsResolveErrorsCannotResolvePrediction
}
type KeysCommandsPredictionsResolveSuccessPredictionResolvedVars struct {
	Outcome any
}
type KeysCommandsPredictionsResolveSuccessPredictionResolved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsResolveSuccessPredictionResolved) IsTranslationKey() {
}
func (k KeysCommandsPredictionsResolveSuccessPredictionResolved) GetPath() string {
	return "commands.predictions.resolve.success.prediction_resolved"
}
func (k KeysCommandsPredictionsResolveSuccessPredictionResolved) GetPathSlice() []string {
	return []string{"commands", "predictions", "resolve", "success", "prediction_resolved"}
}
func (k KeysCommandsPredictionsResolveSuccessPredictionResolved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsResolveSuccessPredictionResolved) SetVars(vars KeysCommandsPredictionsResolveSuccessPredictionResolvedVars) twiri18n.TranslationKey[KeysCommandsPredictionsResolveSuccessPredictionResolvedVars] {
	k.Vars = twiri18n.Vars{"outcome": vars.Outcome}
	return k
}

type KeysCommandsPredictionsResolveSuccess struct {
	PredictionResolved KeysCommandsPredictionsResolveSuccessPredictionResolved
}
type KeysCommandsPredictionsResolveDescriptionVars struct {
}
type KeysCommandsPredictionsResolveDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsResolveDescription) IsTranslationKey() {
}
func (k KeysCommandsPredictionsResolveDescription) GetPath() string {
	return "commands.predictions.resolve.description"
}
func (k KeysCommandsPredictionsResolveDescription) GetPathSlice() []string {
	return []string{"commands", "predictions", "resolve", "description"}
}
func (k KeysCommandsPredictionsResolveDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsResolveDescription) SetVars(vars KeysCommandsPredictionsResolveDescriptionVars) twiri18n.TranslationKey[KeysCommandsPredictionsResolveDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsResolve struct {
	Errors		KeysCommandsPredictionsResolveErrors
	Success		KeysCommandsPredictionsResolveSuccess
	Description	KeysCommandsPredictionsResolveDescription
}
type KeysCommandsPredictionsLockErrorsCannotCreateTwitchClientVars struct {
}
type KeysCommandsPredictionsLockErrorsCannotCreateTwitchClient struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsLockErrorsCannotCreateTwitchClient) IsTranslationKey() {
}
func (k KeysCommandsPredictionsLockErrorsCannotCreateTwitchClient) GetPath() string {
	return "commands.predictions.lock.errors.cannot_create_twitch_client"
}
func (k KeysCommandsPredictionsLockErrorsCannotCreateTwitchClient) GetPathSlice() []string {
	return []string{"commands", "predictions", "lock", "errors", "cannot_create_twitch_client"}
}
func (k KeysCommandsPredictionsLockErrorsCannotCreateTwitchClient) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsLockErrorsCannotCreateTwitchClient) SetVars(vars KeysCommandsPredictionsLockErrorsCannotCreateTwitchClientVars) twiri18n.TranslationKey[KeysCommandsPredictionsLockErrorsCannotCreateTwitchClientVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsLockErrorsCannotGetPredictionVars struct {
}
type KeysCommandsPredictionsLockErrorsCannotGetPrediction struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsLockErrorsCannotGetPrediction) IsTranslationKey() {
}
func (k KeysCommandsPredictionsLockErrorsCannotGetPrediction) GetPath() string {
	return "commands.predictions.lock.errors.cannot_get_prediction"
}
func (k KeysCommandsPredictionsLockErrorsCannotGetPrediction) GetPathSlice() []string {
	return []string{"commands", "predictions", "lock", "errors", "cannot_get_prediction"}
}
func (k KeysCommandsPredictionsLockErrorsCannotGetPrediction) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsLockErrorsCannotGetPrediction) SetVars(vars KeysCommandsPredictionsLockErrorsCannotGetPredictionVars) twiri18n.TranslationKey[KeysCommandsPredictionsLockErrorsCannotGetPredictionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsLockErrorsNoPredictionRunningVars struct {
}
type KeysCommandsPredictionsLockErrorsNoPredictionRunning struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsLockErrorsNoPredictionRunning) IsTranslationKey() {
}
func (k KeysCommandsPredictionsLockErrorsNoPredictionRunning) GetPath() string {
	return "commands.predictions.lock.errors.no_prediction_running"
}
func (k KeysCommandsPredictionsLockErrorsNoPredictionRunning) GetPathSlice() []string {
	return []string{"commands", "predictions", "lock", "errors", "no_prediction_running"}
}
func (k KeysCommandsPredictionsLockErrorsNoPredictionRunning) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsLockErrorsNoPredictionRunning) SetVars(vars KeysCommandsPredictionsLockErrorsNoPredictionRunningVars) twiri18n.TranslationKey[KeysCommandsPredictionsLockErrorsNoPredictionRunningVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsLockErrorsCannotLockPredictionVars struct {
}
type KeysCommandsPredictionsLockErrorsCannotLockPrediction struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsLockErrorsCannotLockPrediction) IsTranslationKey() {
}
func (k KeysCommandsPredictionsLockErrorsCannotLockPrediction) GetPath() string {
	return "commands.predictions.lock.errors.cannot_lock_prediction"
}
func (k KeysCommandsPredictionsLockErrorsCannotLockPrediction) GetPathSlice() []string {
	return []string{"commands", "predictions", "lock", "errors", "cannot_lock_prediction"}
}
func (k KeysCommandsPredictionsLockErrorsCannotLockPrediction) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsLockErrorsCannotLockPrediction) SetVars(vars KeysCommandsPredictionsLockErrorsCannotLockPredictionVars) twiri18n.TranslationKey[KeysCommandsPredictionsLockErrorsCannotLockPredictionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsLockErrors struct {
	CannotCreateTwitchClient	KeysCommandsPredictionsLockErrorsCannotCreateTwitchClient
	CannotGetPrediction		KeysCommandsPredictionsLockErrorsCannotGetPrediction
	NoPredictionRunning		KeysCommandsPredictionsLockErrorsNoPredictionRunning
	CannotLockPrediction		KeysCommandsPredictionsLockErrorsCannotLockPrediction
}
type KeysCommandsPredictionsLockSuccessPredictionLockedVars struct {
}
type KeysCommandsPredictionsLockSuccessPredictionLocked struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsLockSuccessPredictionLocked) IsTranslationKey() {
}
func (k KeysCommandsPredictionsLockSuccessPredictionLocked) GetPath() string {
	return "commands.predictions.lock.success.prediction_locked"
}
func (k KeysCommandsPredictionsLockSuccessPredictionLocked) GetPathSlice() []string {
	return []string{"commands", "predictions", "lock", "success", "prediction_locked"}
}
func (k KeysCommandsPredictionsLockSuccessPredictionLocked) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsLockSuccessPredictionLocked) SetVars(vars KeysCommandsPredictionsLockSuccessPredictionLockedVars) twiri18n.TranslationKey[KeysCommandsPredictionsLockSuccessPredictionLockedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsLockSuccess struct {
	PredictionLocked KeysCommandsPredictionsLockSuccessPredictionLocked
}
type KeysCommandsPredictionsLockDescriptionVars struct {
}
type KeysCommandsPredictionsLockDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsLockDescription) IsTranslationKey() {
}
func (k KeysCommandsPredictionsLockDescription) GetPath() string {
	return "commands.predictions.lock.description"
}
func (k KeysCommandsPredictionsLockDescription) GetPathSlice() []string {
	return []string{"commands", "predictions", "lock", "description"}
}
func (k KeysCommandsPredictionsLockDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsLockDescription) SetVars(vars KeysCommandsPredictionsLockDescriptionVars) twiri18n.TranslationKey[KeysCommandsPredictionsLockDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsLock struct {
	Errors		KeysCommandsPredictionsLockErrors
	Success		KeysCommandsPredictionsLockSuccess
	Description	KeysCommandsPredictionsLockDescription
}
type KeysCommandsPredictionsCancelErrorsCannotCreateTwitchClientVars struct {
}
type KeysCommandsPredictionsCancelErrorsCannotCreateTwitchClient struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsCancelErrorsCannotCreateTwitchClient) IsTranslationKey() {
}
func (k KeysCommandsPredictionsCancelErrorsCannotCreateTwitchClient) GetPath() string {
	return "commands.predictions.cancel.errors.cannot_create_twitch_client"
}
func (k KeysCommandsPredictionsCancelErrorsCannotCreateTwitchClient) GetPathSlice() []string {
	return []string{"commands", "predictions", "cancel", "errors", "cannot_create_twitch_client"}
}
func (k KeysCommandsPredictionsCancelErrorsCannotCreateTwitchClient) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsCancelErrorsCannotCreateTwitchClient) SetVars(vars KeysCommandsPredictionsCancelErrorsCannotCreateTwitchClientVars) twiri18n.TranslationKey[KeysCommandsPredictionsCancelErrorsCannotCreateTwitchClientVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsCancelErrorsCannotGetPredictionVars struct {
}
type KeysCommandsPredictionsCancelErrorsCannotGetPrediction struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsCancelErrorsCannotGetPrediction) IsTranslationKey() {
}
func (k KeysCommandsPredictionsCancelErrorsCannotGetPrediction) GetPath() string {
	return "commands.predictions.cancel.errors.cannot_get_prediction"
}
func (k KeysCommandsPredictionsCancelErrorsCannotGetPrediction) GetPathSlice() []string {
	return []string{"commands", "predictions", "cancel", "errors", "cannot_get_prediction"}
}
func (k KeysCommandsPredictionsCancelErrorsCannotGetPrediction) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsCancelErrorsCannotGetPrediction) SetVars(vars KeysCommandsPredictionsCancelErrorsCannotGetPredictionVars) twiri18n.TranslationKey[KeysCommandsPredictionsCancelErrorsCannotGetPredictionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsCancelErrorsNoPredictionRunningVars struct {
}
type KeysCommandsPredictionsCancelErrorsNoPredictionRunning struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsCancelErrorsNoPredictionRunning) IsTranslationKey() {
}
func (k KeysCommandsPredictionsCancelErrorsNoPredictionRunning) GetPath() string {
	return "commands.predictions.cancel.errors.no_prediction_running"
}
func (k KeysCommandsPredictionsCancelErrorsNoPredictionRunning) GetPathSlice() []string {
	return []string{"commands", "predictions", "cancel", "errors", "no_prediction_running"}
}
func (k KeysCommandsPredictionsCancelErrorsNoPredictionRunning) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsCancelErrorsNoPredictionRunning) SetVars(vars KeysCommandsPredictionsCancelErrorsNoPredictionRunningVars) twiri18n.TranslationKey[KeysCommandsPredictionsCancelErrorsNoPredictionRunningVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsCancelErrorsCannotCancelPredictionVars struct {
}
type KeysCommandsPredictionsCancelErrorsCannotCancelPrediction struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsCancelErrorsCannotCancelPrediction) IsTranslationKey() {
}
func (k KeysCommandsPredictionsCancelErrorsCannotCancelPrediction) GetPath() string {
	return "commands.predictions.cancel.errors.cannot_cancel_prediction"
}
func (k KeysCommandsPredictionsCancelErrorsCannotCancelPrediction) GetPathSlice() []string {
	return []string{"commands", "predictions", "cancel", "errors", "cannot_cancel_prediction"}
}
func (k KeysCommandsPredictionsCancelErrorsCannotCancelPrediction) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsCancelErrorsCannotCancelPrediction) SetVars(vars KeysCommandsPredictionsCancelErrorsCannotCancelPredictionVars) twiri18n.TranslationKey[KeysCommandsPredictionsCancelErrorsCannotCancelPredictionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsCancelErrors struct {
	CannotCreateTwitchClient	KeysCommandsPredictionsCancelErrorsCannotCreateTwitchClient
	CannotGetPrediction		KeysCommandsPredictionsCancelErrorsCannotGetPrediction
	NoPredictionRunning		KeysCommandsPredictionsCancelErrorsNoPredictionRunning
	CannotCancelPrediction		KeysCommandsPredictionsCancelErrorsCannotCancelPrediction
}
type KeysCommandsPredictionsCancelSuccessPredictionCancelledVars struct {
}
type KeysCommandsPredictionsCancelSuccessPredictionCancelled struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsCancelSuccessPredictionCancelled) IsTranslationKey() {
}
func (k KeysCommandsPredictionsCancelSuccessPredictionCancelled) GetPath() string {
	return "commands.predictions.cancel.success.prediction_cancelled"
}
func (k KeysCommandsPredictionsCancelSuccessPredictionCancelled) GetPathSlice() []string {
	return []string{"commands", "predictions", "cancel", "success", "prediction_cancelled"}
}
func (k KeysCommandsPredictionsCancelSuccessPredictionCancelled) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsCancelSuccessPredictionCancelled) SetVars(vars KeysCommandsPredictionsCancelSuccessPredictionCancelledVars) twiri18n.TranslationKey[KeysCommandsPredictionsCancelSuccessPredictionCancelledVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsCancelSuccess struct {
	PredictionCancelled KeysCommandsPredictionsCancelSuccessPredictionCancelled
}
type KeysCommandsPredictionsCancelDescriptionVars struct {
}
type KeysCommandsPredictionsCancelDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsCancelDescription) IsTranslationKey() {
}
func (k KeysCommandsPredictionsCancelDescription) GetPath() string {
	return "commands.predictions.cancel.description"
}
func (k KeysCommandsPredictionsCancelDescription) GetPathSlice() []string {
	return []string{"commands", "predictions", "cancel", "description"}
}
func (k KeysCommandsPredictionsCancelDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsCancelDescription) SetVars(vars KeysCommandsPredictionsCancelDescriptionVars) twiri18n.TranslationKey[KeysCommandsPredictionsCancelDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsCancel struct {
	Errors		KeysCommandsPredictionsCancelErrors
	Success		KeysCommandsPredictionsCancelSuccess
	Description	KeysCommandsPredictionsCancelDescription
}
type KeysCommandsPredictions struct {
	Start	KeysCommandsPredictionsStart
	Resolve	KeysCommandsPredictionsResolve
	Lock	KeysCommandsPredictionsLock
	Cancel	KeysCommandsPredictionsCancel
}
type KeysCommandsSpamDescriptionVars struct {
}
type KeysCommandsSpamDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSpamDescription) IsTranslationKey() {
}
func (k KeysCommandsSpamDescription) GetPath() string {
	return "commands.spam.description"
}
func (k KeysCommandsSpamDescription) GetPathSlice() []string {
	return []string{"commands", "spam", "description"}
}
func (k KeysCommandsSpamDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSpamDescription) SetVars(vars KeysCommandsSpamDescriptionVars) twiri18n.TranslationKey[KeysCommandsSpamDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSpamTimeoutDescriptionVars struct {
}
type KeysCommandsSpamTimeoutDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSpamTimeoutDescription) IsTranslationKey() {
}
func (k KeysCommandsSpamTimeoutDescription) GetPath() string {
	return "commands.spam.timeout.description"
}
func (k KeysCommandsSpamTimeoutDescription) GetPathSlice() []string {
	return []string{"commands", "spam", "timeout", "description"}
}
func (k KeysCommandsSpamTimeoutDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSpamTimeoutDescription) SetVars(vars KeysCommandsSpamTimeoutDescriptionVars) twiri18n.TranslationKey[KeysCommandsSpamTimeoutDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSpamTimeoutErrorsCannotTimeoutUserVars struct {
}
type KeysCommandsSpamTimeoutErrorsCannotTimeoutUser struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSpamTimeoutErrorsCannotTimeoutUser) IsTranslationKey() {
}
func (k KeysCommandsSpamTimeoutErrorsCannotTimeoutUser) GetPath() string {
	return "commands.spam.timeout.errors.cannot_timeout_user"
}
func (k KeysCommandsSpamTimeoutErrorsCannotTimeoutUser) GetPathSlice() []string {
	return []string{"commands", "spam", "timeout", "errors", "cannot_timeout_user"}
}
func (k KeysCommandsSpamTimeoutErrorsCannotTimeoutUser) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSpamTimeoutErrorsCannotTimeoutUser) SetVars(vars KeysCommandsSpamTimeoutErrorsCannotTimeoutUserVars) twiri18n.TranslationKey[KeysCommandsSpamTimeoutErrorsCannotTimeoutUserVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSpamTimeoutErrorsUserNotFoundVars struct {
}
type KeysCommandsSpamTimeoutErrorsUserNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSpamTimeoutErrorsUserNotFound) IsTranslationKey() {
}
func (k KeysCommandsSpamTimeoutErrorsUserNotFound) GetPath() string {
	return "commands.spam.timeout.errors.user_not_found"
}
func (k KeysCommandsSpamTimeoutErrorsUserNotFound) GetPathSlice() []string {
	return []string{"commands", "spam", "timeout", "errors", "user_not_found"}
}
func (k KeysCommandsSpamTimeoutErrorsUserNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSpamTimeoutErrorsUserNotFound) SetVars(vars KeysCommandsSpamTimeoutErrorsUserNotFoundVars) twiri18n.TranslationKey[KeysCommandsSpamTimeoutErrorsUserNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSpamTimeoutErrors struct {
	CannotTimeoutUser	KeysCommandsSpamTimeoutErrorsCannotTimeoutUser
	UserNotFound		KeysCommandsSpamTimeoutErrorsUserNotFound
}
type KeysCommandsSpamTimeoutSuccessUserTimedOutVars struct {
	User any
}
type KeysCommandsSpamTimeoutSuccessUserTimedOut struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSpamTimeoutSuccessUserTimedOut) IsTranslationKey() {
}
func (k KeysCommandsSpamTimeoutSuccessUserTimedOut) GetPath() string {
	return "commands.spam.timeout.success.user_timed_out"
}
func (k KeysCommandsSpamTimeoutSuccessUserTimedOut) GetPathSlice() []string {
	return []string{"commands", "spam", "timeout", "success", "user_timed_out"}
}
func (k KeysCommandsSpamTimeoutSuccessUserTimedOut) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSpamTimeoutSuccessUserTimedOut) SetVars(vars KeysCommandsSpamTimeoutSuccessUserTimedOutVars) twiri18n.TranslationKey[KeysCommandsSpamTimeoutSuccessUserTimedOutVars] {
	k.Vars = twiri18n.Vars{"user": vars.User}
	return k
}

type KeysCommandsSpamTimeoutSuccess struct {
	UserTimedOut KeysCommandsSpamTimeoutSuccessUserTimedOut
}
type KeysCommandsSpamTimeout struct {
	Description	KeysCommandsSpamTimeoutDescription
	Errors		KeysCommandsSpamTimeoutErrors
	Success		KeysCommandsSpamTimeoutSuccess
}
type KeysCommandsSpamSettingsDescriptionVars struct {
}
type KeysCommandsSpamSettingsDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSpamSettingsDescription) IsTranslationKey() {
}
func (k KeysCommandsSpamSettingsDescription) GetPath() string {
	return "commands.spam.settings.description"
}
func (k KeysCommandsSpamSettingsDescription) GetPathSlice() []string {
	return []string{"commands", "spam", "settings", "description"}
}
func (k KeysCommandsSpamSettingsDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSpamSettingsDescription) SetVars(vars KeysCommandsSpamSettingsDescriptionVars) twiri18n.TranslationKey[KeysCommandsSpamSettingsDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSpamSettingsErrorsCannotUpdateSettingsVars struct {
}
type KeysCommandsSpamSettingsErrorsCannotUpdateSettings struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSpamSettingsErrorsCannotUpdateSettings) IsTranslationKey() {
}
func (k KeysCommandsSpamSettingsErrorsCannotUpdateSettings) GetPath() string {
	return "commands.spam.settings.errors.cannot_update_settings"
}
func (k KeysCommandsSpamSettingsErrorsCannotUpdateSettings) GetPathSlice() []string {
	return []string{"commands", "spam", "settings", "errors", "cannot_update_settings"}
}
func (k KeysCommandsSpamSettingsErrorsCannotUpdateSettings) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSpamSettingsErrorsCannotUpdateSettings) SetVars(vars KeysCommandsSpamSettingsErrorsCannotUpdateSettingsVars) twiri18n.TranslationKey[KeysCommandsSpamSettingsErrorsCannotUpdateSettingsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSpamSettingsErrorsInvalidSettingsVars struct {
}
type KeysCommandsSpamSettingsErrorsInvalidSettings struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSpamSettingsErrorsInvalidSettings) IsTranslationKey() {
}
func (k KeysCommandsSpamSettingsErrorsInvalidSettings) GetPath() string {
	return "commands.spam.settings.errors.invalid_settings"
}
func (k KeysCommandsSpamSettingsErrorsInvalidSettings) GetPathSlice() []string {
	return []string{"commands", "spam", "settings", "errors", "invalid_settings"}
}
func (k KeysCommandsSpamSettingsErrorsInvalidSettings) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSpamSettingsErrorsInvalidSettings) SetVars(vars KeysCommandsSpamSettingsErrorsInvalidSettingsVars) twiri18n.TranslationKey[KeysCommandsSpamSettingsErrorsInvalidSettingsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSpamSettingsErrors struct {
	CannotUpdateSettings	KeysCommandsSpamSettingsErrorsCannotUpdateSettings
	InvalidSettings		KeysCommandsSpamSettingsErrorsInvalidSettings
}
type KeysCommandsSpamSettingsSuccessSettingsUpdatedVars struct {
}
type KeysCommandsSpamSettingsSuccessSettingsUpdated struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSpamSettingsSuccessSettingsUpdated) IsTranslationKey() {
}
func (k KeysCommandsSpamSettingsSuccessSettingsUpdated) GetPath() string {
	return "commands.spam.settings.success.settings_updated"
}
func (k KeysCommandsSpamSettingsSuccessSettingsUpdated) GetPathSlice() []string {
	return []string{"commands", "spam", "settings", "success", "settings_updated"}
}
func (k KeysCommandsSpamSettingsSuccessSettingsUpdated) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSpamSettingsSuccessSettingsUpdated) SetVars(vars KeysCommandsSpamSettingsSuccessSettingsUpdatedVars) twiri18n.TranslationKey[KeysCommandsSpamSettingsSuccessSettingsUpdatedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSpamSettingsSuccess struct {
	SettingsUpdated KeysCommandsSpamSettingsSuccessSettingsUpdated
}
type KeysCommandsSpamSettings struct {
	Description	KeysCommandsSpamSettingsDescription
	Errors		KeysCommandsSpamSettingsErrors
	Success		KeysCommandsSpamSettingsSuccess
}
type KeysCommandsSpam struct {
	Description	KeysCommandsSpamDescription
	Timeout		KeysCommandsSpamTimeout
	Settings	KeysCommandsSpamSettings
}
type KeysCommandsTracerDisableErrorsCannotDisableTracerVars struct {
}
type KeysCommandsTracerDisableErrorsCannotDisableTracer struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTracerDisableErrorsCannotDisableTracer) IsTranslationKey() {
}
func (k KeysCommandsTracerDisableErrorsCannotDisableTracer) GetPath() string {
	return "commands.tracer.disable.errors.cannot_disable_tracer"
}
func (k KeysCommandsTracerDisableErrorsCannotDisableTracer) GetPathSlice() []string {
	return []string{"commands", "tracer", "disable", "errors", "cannot_disable_tracer"}
}
func (k KeysCommandsTracerDisableErrorsCannotDisableTracer) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTracerDisableErrorsCannotDisableTracer) SetVars(vars KeysCommandsTracerDisableErrorsCannotDisableTracerVars) twiri18n.TranslationKey[KeysCommandsTracerDisableErrorsCannotDisableTracerVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTracerDisableErrors struct {
	CannotDisableTracer KeysCommandsTracerDisableErrorsCannotDisableTracer
}
type KeysCommandsTracerDisableSuccessTracerDisabledVars struct {
}
type KeysCommandsTracerDisableSuccessTracerDisabled struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTracerDisableSuccessTracerDisabled) IsTranslationKey() {
}
func (k KeysCommandsTracerDisableSuccessTracerDisabled) GetPath() string {
	return "commands.tracer.disable.success.tracer_disabled"
}
func (k KeysCommandsTracerDisableSuccessTracerDisabled) GetPathSlice() []string {
	return []string{"commands", "tracer", "disable", "success", "tracer_disabled"}
}
func (k KeysCommandsTracerDisableSuccessTracerDisabled) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTracerDisableSuccessTracerDisabled) SetVars(vars KeysCommandsTracerDisableSuccessTracerDisabledVars) twiri18n.TranslationKey[KeysCommandsTracerDisableSuccessTracerDisabledVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTracerDisableSuccess struct {
	TracerDisabled KeysCommandsTracerDisableSuccessTracerDisabled
}
type KeysCommandsTracerDisableDescriptionVars struct {
}
type KeysCommandsTracerDisableDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTracerDisableDescription) IsTranslationKey() {
}
func (k KeysCommandsTracerDisableDescription) GetPath() string {
	return "commands.tracer.disable.description"
}
func (k KeysCommandsTracerDisableDescription) GetPathSlice() []string {
	return []string{"commands", "tracer", "disable", "description"}
}
func (k KeysCommandsTracerDisableDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTracerDisableDescription) SetVars(vars KeysCommandsTracerDisableDescriptionVars) twiri18n.TranslationKey[KeysCommandsTracerDisableDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTracerDisable struct {
	Errors		KeysCommandsTracerDisableErrors
	Success		KeysCommandsTracerDisableSuccess
	Description	KeysCommandsTracerDisableDescription
}
type KeysCommandsTracerStatusDescriptionVars struct {
}
type KeysCommandsTracerStatusDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTracerStatusDescription) IsTranslationKey() {
}
func (k KeysCommandsTracerStatusDescription) GetPath() string {
	return "commands.tracer.status.description"
}
func (k KeysCommandsTracerStatusDescription) GetPathSlice() []string {
	return []string{"commands", "tracer", "status", "description"}
}
func (k KeysCommandsTracerStatusDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTracerStatusDescription) SetVars(vars KeysCommandsTracerStatusDescriptionVars) twiri18n.TranslationKey[KeysCommandsTracerStatusDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTracerStatusErrorsCannotGetStatusVars struct {
}
type KeysCommandsTracerStatusErrorsCannotGetStatus struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTracerStatusErrorsCannotGetStatus) IsTranslationKey() {
}
func (k KeysCommandsTracerStatusErrorsCannotGetStatus) GetPath() string {
	return "commands.tracer.status.errors.cannot_get_status"
}
func (k KeysCommandsTracerStatusErrorsCannotGetStatus) GetPathSlice() []string {
	return []string{"commands", "tracer", "status", "errors", "cannot_get_status"}
}
func (k KeysCommandsTracerStatusErrorsCannotGetStatus) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTracerStatusErrorsCannotGetStatus) SetVars(vars KeysCommandsTracerStatusErrorsCannotGetStatusVars) twiri18n.TranslationKey[KeysCommandsTracerStatusErrorsCannotGetStatusVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTracerStatusErrors struct {
	CannotGetStatus KeysCommandsTracerStatusErrorsCannotGetStatus
}
type KeysCommandsTracerStatusSuccessStatusDisplayVars struct {
	Status any
}
type KeysCommandsTracerStatusSuccessStatusDisplay struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTracerStatusSuccessStatusDisplay) IsTranslationKey() {
}
func (k KeysCommandsTracerStatusSuccessStatusDisplay) GetPath() string {
	return "commands.tracer.status.success.status_display"
}
func (k KeysCommandsTracerStatusSuccessStatusDisplay) GetPathSlice() []string {
	return []string{"commands", "tracer", "status", "success", "status_display"}
}
func (k KeysCommandsTracerStatusSuccessStatusDisplay) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTracerStatusSuccessStatusDisplay) SetVars(vars KeysCommandsTracerStatusSuccessStatusDisplayVars) twiri18n.TranslationKey[KeysCommandsTracerStatusSuccessStatusDisplayVars] {
	k.Vars = twiri18n.Vars{"status": vars.Status}
	return k
}

type KeysCommandsTracerStatusSuccess struct {
	StatusDisplay KeysCommandsTracerStatusSuccessStatusDisplay
}
type KeysCommandsTracerStatus struct {
	Description	KeysCommandsTracerStatusDescription
	Errors		KeysCommandsTracerStatusErrors
	Success		KeysCommandsTracerStatusSuccess
}
type KeysCommandsTracerDescriptionVars struct {
}
type KeysCommandsTracerDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTracerDescription) IsTranslationKey() {
}
func (k KeysCommandsTracerDescription) GetPath() string {
	return "commands.tracer.description"
}
func (k KeysCommandsTracerDescription) GetPathSlice() []string {
	return []string{"commands", "tracer", "description"}
}
func (k KeysCommandsTracerDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTracerDescription) SetVars(vars KeysCommandsTracerDescriptionVars) twiri18n.TranslationKey[KeysCommandsTracerDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTracerEnableDescriptionVars struct {
}
type KeysCommandsTracerEnableDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTracerEnableDescription) IsTranslationKey() {
}
func (k KeysCommandsTracerEnableDescription) GetPath() string {
	return "commands.tracer.enable.description"
}
func (k KeysCommandsTracerEnableDescription) GetPathSlice() []string {
	return []string{"commands", "tracer", "enable", "description"}
}
func (k KeysCommandsTracerEnableDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTracerEnableDescription) SetVars(vars KeysCommandsTracerEnableDescriptionVars) twiri18n.TranslationKey[KeysCommandsTracerEnableDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTracerEnableErrorsCannotEnableTracerVars struct {
}
type KeysCommandsTracerEnableErrorsCannotEnableTracer struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTracerEnableErrorsCannotEnableTracer) IsTranslationKey() {
}
func (k KeysCommandsTracerEnableErrorsCannotEnableTracer) GetPath() string {
	return "commands.tracer.enable.errors.cannot_enable_tracer"
}
func (k KeysCommandsTracerEnableErrorsCannotEnableTracer) GetPathSlice() []string {
	return []string{"commands", "tracer", "enable", "errors", "cannot_enable_tracer"}
}
func (k KeysCommandsTracerEnableErrorsCannotEnableTracer) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTracerEnableErrorsCannotEnableTracer) SetVars(vars KeysCommandsTracerEnableErrorsCannotEnableTracerVars) twiri18n.TranslationKey[KeysCommandsTracerEnableErrorsCannotEnableTracerVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTracerEnableErrors struct {
	CannotEnableTracer KeysCommandsTracerEnableErrorsCannotEnableTracer
}
type KeysCommandsTracerEnableSuccessTracerEnabledVars struct {
}
type KeysCommandsTracerEnableSuccessTracerEnabled struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTracerEnableSuccessTracerEnabled) IsTranslationKey() {
}
func (k KeysCommandsTracerEnableSuccessTracerEnabled) GetPath() string {
	return "commands.tracer.enable.success.tracer_enabled"
}
func (k KeysCommandsTracerEnableSuccessTracerEnabled) GetPathSlice() []string {
	return []string{"commands", "tracer", "enable", "success", "tracer_enabled"}
}
func (k KeysCommandsTracerEnableSuccessTracerEnabled) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTracerEnableSuccessTracerEnabled) SetVars(vars KeysCommandsTracerEnableSuccessTracerEnabledVars) twiri18n.TranslationKey[KeysCommandsTracerEnableSuccessTracerEnabledVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTracerEnableSuccess struct {
	TracerEnabled KeysCommandsTracerEnableSuccessTracerEnabled
}
type KeysCommandsTracerEnable struct {
	Description	KeysCommandsTracerEnableDescription
	Errors		KeysCommandsTracerEnableErrors
	Success		KeysCommandsTracerEnableSuccess
}
type KeysCommandsTracer struct {
	Disable		KeysCommandsTracerDisable
	Status		KeysCommandsTracerStatus
	Description	KeysCommandsTracerDescription
	Enable		KeysCommandsTracerEnable
}
type KeysCommandsManageDescriptionVars struct {
}
type KeysCommandsManageDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageDescription) IsTranslationKey() {
}
func (k KeysCommandsManageDescription) GetPath() string {
	return "commands.manage.description"
}
func (k KeysCommandsManageDescription) GetPathSlice() []string {
	return []string{"commands", "manage", "description"}
}
func (k KeysCommandsManageDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageDescription) SetVars(vars KeysCommandsManageDescriptionVars) twiri18n.TranslationKey[KeysCommandsManageDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageAddDescriptionVars struct {
}
type KeysCommandsManageAddDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageAddDescription) IsTranslationKey() {
}
func (k KeysCommandsManageAddDescription) GetPath() string {
	return "commands.manage.add.description"
}
func (k KeysCommandsManageAddDescription) GetPathSlice() []string {
	return []string{"commands", "manage", "add", "description"}
}
func (k KeysCommandsManageAddDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageAddDescription) SetVars(vars KeysCommandsManageAddDescriptionVars) twiri18n.TranslationKey[KeysCommandsManageAddDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageAddErrorsUserNotFoundVars struct {
}
type KeysCommandsManageAddErrorsUserNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageAddErrorsUserNotFound) IsTranslationKey() {
}
func (k KeysCommandsManageAddErrorsUserNotFound) GetPath() string {
	return "commands.manage.add.errors.user_not_found"
}
func (k KeysCommandsManageAddErrorsUserNotFound) GetPathSlice() []string {
	return []string{"commands", "manage", "add", "errors", "user_not_found"}
}
func (k KeysCommandsManageAddErrorsUserNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageAddErrorsUserNotFound) SetVars(vars KeysCommandsManageAddErrorsUserNotFoundVars) twiri18n.TranslationKey[KeysCommandsManageAddErrorsUserNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageAddErrorsPermissionAlreadyExistsVars struct {
}
type KeysCommandsManageAddErrorsPermissionAlreadyExists struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageAddErrorsPermissionAlreadyExists) IsTranslationKey() {
}
func (k KeysCommandsManageAddErrorsPermissionAlreadyExists) GetPath() string {
	return "commands.manage.add.errors.permission_already_exists"
}
func (k KeysCommandsManageAddErrorsPermissionAlreadyExists) GetPathSlice() []string {
	return []string{"commands", "manage", "add", "errors", "permission_already_exists"}
}
func (k KeysCommandsManageAddErrorsPermissionAlreadyExists) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageAddErrorsPermissionAlreadyExists) SetVars(vars KeysCommandsManageAddErrorsPermissionAlreadyExistsVars) twiri18n.TranslationKey[KeysCommandsManageAddErrorsPermissionAlreadyExistsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageAddErrorsCannotAddPermissionVars struct {
}
type KeysCommandsManageAddErrorsCannotAddPermission struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageAddErrorsCannotAddPermission) IsTranslationKey() {
}
func (k KeysCommandsManageAddErrorsCannotAddPermission) GetPath() string {
	return "commands.manage.add.errors.cannot_add_permission"
}
func (k KeysCommandsManageAddErrorsCannotAddPermission) GetPathSlice() []string {
	return []string{"commands", "manage", "add", "errors", "cannot_add_permission"}
}
func (k KeysCommandsManageAddErrorsCannotAddPermission) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageAddErrorsCannotAddPermission) SetVars(vars KeysCommandsManageAddErrorsCannotAddPermissionVars) twiri18n.TranslationKey[KeysCommandsManageAddErrorsCannotAddPermissionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageAddErrors struct {
	UserNotFound		KeysCommandsManageAddErrorsUserNotFound
	PermissionAlreadyExists	KeysCommandsManageAddErrorsPermissionAlreadyExists
	CannotAddPermission	KeysCommandsManageAddErrorsCannotAddPermission
}
type KeysCommandsManageAddSuccessPermissionAddedVars struct {
	Permission	any
	User		any
}
type KeysCommandsManageAddSuccessPermissionAdded struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageAddSuccessPermissionAdded) IsTranslationKey() {
}
func (k KeysCommandsManageAddSuccessPermissionAdded) GetPath() string {
	return "commands.manage.add.success.permission_added"
}
func (k KeysCommandsManageAddSuccessPermissionAdded) GetPathSlice() []string {
	return []string{"commands", "manage", "add", "success", "permission_added"}
}
func (k KeysCommandsManageAddSuccessPermissionAdded) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageAddSuccessPermissionAdded) SetVars(vars KeysCommandsManageAddSuccessPermissionAddedVars) twiri18n.TranslationKey[KeysCommandsManageAddSuccessPermissionAddedVars] {
	k.Vars = twiri18n.Vars{"permission": vars.Permission, "user": vars.User}
	return k
}

type KeysCommandsManageAddSuccess struct {
	PermissionAdded KeysCommandsManageAddSuccessPermissionAdded
}
type KeysCommandsManageAdd struct {
	Description	KeysCommandsManageAddDescription
	Errors		KeysCommandsManageAddErrors
	Success		KeysCommandsManageAddSuccess
}
type KeysCommandsManageRemoveDescriptionVars struct {
}
type KeysCommandsManageRemoveDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageRemoveDescription) IsTranslationKey() {
}
func (k KeysCommandsManageRemoveDescription) GetPath() string {
	return "commands.manage.remove.description"
}
func (k KeysCommandsManageRemoveDescription) GetPathSlice() []string {
	return []string{"commands", "manage", "remove", "description"}
}
func (k KeysCommandsManageRemoveDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageRemoveDescription) SetVars(vars KeysCommandsManageRemoveDescriptionVars) twiri18n.TranslationKey[KeysCommandsManageRemoveDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageRemoveErrorsCannotRemovePermissionVars struct {
}
type KeysCommandsManageRemoveErrorsCannotRemovePermission struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageRemoveErrorsCannotRemovePermission) IsTranslationKey() {
}
func (k KeysCommandsManageRemoveErrorsCannotRemovePermission) GetPath() string {
	return "commands.manage.remove.errors.cannot_remove_permission"
}
func (k KeysCommandsManageRemoveErrorsCannotRemovePermission) GetPathSlice() []string {
	return []string{"commands", "manage", "remove", "errors", "cannot_remove_permission"}
}
func (k KeysCommandsManageRemoveErrorsCannotRemovePermission) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageRemoveErrorsCannotRemovePermission) SetVars(vars KeysCommandsManageRemoveErrorsCannotRemovePermissionVars) twiri18n.TranslationKey[KeysCommandsManageRemoveErrorsCannotRemovePermissionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageRemoveErrorsPermissionNotFoundVars struct {
}
type KeysCommandsManageRemoveErrorsPermissionNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageRemoveErrorsPermissionNotFound) IsTranslationKey() {
}
func (k KeysCommandsManageRemoveErrorsPermissionNotFound) GetPath() string {
	return "commands.manage.remove.errors.permission_not_found"
}
func (k KeysCommandsManageRemoveErrorsPermissionNotFound) GetPathSlice() []string {
	return []string{"commands", "manage", "remove", "errors", "permission_not_found"}
}
func (k KeysCommandsManageRemoveErrorsPermissionNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageRemoveErrorsPermissionNotFound) SetVars(vars KeysCommandsManageRemoveErrorsPermissionNotFoundVars) twiri18n.TranslationKey[KeysCommandsManageRemoveErrorsPermissionNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageRemoveErrors struct {
	CannotRemovePermission	KeysCommandsManageRemoveErrorsCannotRemovePermission
	PermissionNotFound	KeysCommandsManageRemoveErrorsPermissionNotFound
}
type KeysCommandsManageRemoveSuccessPermissionRemovedVars struct {
	Permission	any
	User		any
}
type KeysCommandsManageRemoveSuccessPermissionRemoved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageRemoveSuccessPermissionRemoved) IsTranslationKey() {
}
func (k KeysCommandsManageRemoveSuccessPermissionRemoved) GetPath() string {
	return "commands.manage.remove.success.permission_removed"
}
func (k KeysCommandsManageRemoveSuccessPermissionRemoved) GetPathSlice() []string {
	return []string{"commands", "manage", "remove", "success", "permission_removed"}
}
func (k KeysCommandsManageRemoveSuccessPermissionRemoved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageRemoveSuccessPermissionRemoved) SetVars(vars KeysCommandsManageRemoveSuccessPermissionRemovedVars) twiri18n.TranslationKey[KeysCommandsManageRemoveSuccessPermissionRemovedVars] {
	k.Vars = twiri18n.Vars{"permission": vars.Permission, "user": vars.User}
	return k
}

type KeysCommandsManageRemoveSuccess struct {
	PermissionRemoved KeysCommandsManageRemoveSuccessPermissionRemoved
}
type KeysCommandsManageRemove struct {
	Description	KeysCommandsManageRemoveDescription
	Errors		KeysCommandsManageRemoveErrors
	Success		KeysCommandsManageRemoveSuccess
}
type KeysCommandsManageListDescriptionVars struct {
}
type KeysCommandsManageListDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageListDescription) IsTranslationKey() {
}
func (k KeysCommandsManageListDescription) GetPath() string {
	return "commands.manage.list.description"
}
func (k KeysCommandsManageListDescription) GetPathSlice() []string {
	return []string{"commands", "manage", "list", "description"}
}
func (k KeysCommandsManageListDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageListDescription) SetVars(vars KeysCommandsManageListDescriptionVars) twiri18n.TranslationKey[KeysCommandsManageListDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageListErrorsCannotListPermissionsVars struct {
}
type KeysCommandsManageListErrorsCannotListPermissions struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageListErrorsCannotListPermissions) IsTranslationKey() {
}
func (k KeysCommandsManageListErrorsCannotListPermissions) GetPath() string {
	return "commands.manage.list.errors.cannot_list_permissions"
}
func (k KeysCommandsManageListErrorsCannotListPermissions) GetPathSlice() []string {
	return []string{"commands", "manage", "list", "errors", "cannot_list_permissions"}
}
func (k KeysCommandsManageListErrorsCannotListPermissions) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageListErrorsCannotListPermissions) SetVars(vars KeysCommandsManageListErrorsCannotListPermissionsVars) twiri18n.TranslationKey[KeysCommandsManageListErrorsCannotListPermissionsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageListErrors struct {
	CannotListPermissions KeysCommandsManageListErrorsCannotListPermissions
}
type KeysCommandsManageListSuccessPermissionsListVars struct {
	User		any
	Permissions	any
}
type KeysCommandsManageListSuccessPermissionsList struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageListSuccessPermissionsList) IsTranslationKey() {
}
func (k KeysCommandsManageListSuccessPermissionsList) GetPath() string {
	return "commands.manage.list.success.permissions_list"
}
func (k KeysCommandsManageListSuccessPermissionsList) GetPathSlice() []string {
	return []string{"commands", "manage", "list", "success", "permissions_list"}
}
func (k KeysCommandsManageListSuccessPermissionsList) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageListSuccessPermissionsList) SetVars(vars KeysCommandsManageListSuccessPermissionsListVars) twiri18n.TranslationKey[KeysCommandsManageListSuccessPermissionsListVars] {
	k.Vars = twiri18n.Vars{"user": vars.User, "permissions": vars.Permissions}
	return k
}

type KeysCommandsManageListSuccess struct {
	PermissionsList KeysCommandsManageListSuccessPermissionsList
}
type KeysCommandsManageList struct {
	Description	KeysCommandsManageListDescription
	Errors		KeysCommandsManageListErrors
	Success		KeysCommandsManageListSuccess
}
type KeysCommandsManage struct {
	Description	KeysCommandsManageDescription
	Add		KeysCommandsManageAdd
	Remove		KeysCommandsManageRemove
	List		KeysCommandsManageList
}
type KeysCommandsNukeDescriptionVars struct {
}
type KeysCommandsNukeDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsNukeDescription) IsTranslationKey() {
}
func (k KeysCommandsNukeDescription) GetPath() string {
	return "commands.nuke.description"
}
func (k KeysCommandsNukeDescription) GetPathSlice() []string {
	return []string{"commands", "nuke", "description"}
}
func (k KeysCommandsNukeDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsNukeDescription) SetVars(vars KeysCommandsNukeDescriptionVars) twiri18n.TranslationKey[KeysCommandsNukeDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsNukeErrorsCannotCreateBroadcasterClientVars struct {
}
type KeysCommandsNukeErrorsCannotCreateBroadcasterClient struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsNukeErrorsCannotCreateBroadcasterClient) IsTranslationKey() {
}
func (k KeysCommandsNukeErrorsCannotCreateBroadcasterClient) GetPath() string {
	return "commands.nuke.errors.cannot_create_broadcaster_client"
}
func (k KeysCommandsNukeErrorsCannotCreateBroadcasterClient) GetPathSlice() []string {
	return []string{"commands", "nuke", "errors", "cannot_create_broadcaster_client"}
}
func (k KeysCommandsNukeErrorsCannotCreateBroadcasterClient) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsNukeErrorsCannotCreateBroadcasterClient) SetVars(vars KeysCommandsNukeErrorsCannotCreateBroadcasterClientVars) twiri18n.TranslationKey[KeysCommandsNukeErrorsCannotCreateBroadcasterClientVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsNukeErrorsCannotClearChatVars struct {
}
type KeysCommandsNukeErrorsCannotClearChat struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsNukeErrorsCannotClearChat) IsTranslationKey() {
}
func (k KeysCommandsNukeErrorsCannotClearChat) GetPath() string {
	return "commands.nuke.errors.cannot_clear_chat"
}
func (k KeysCommandsNukeErrorsCannotClearChat) GetPathSlice() []string {
	return []string{"commands", "nuke", "errors", "cannot_clear_chat"}
}
func (k KeysCommandsNukeErrorsCannotClearChat) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsNukeErrorsCannotClearChat) SetVars(vars KeysCommandsNukeErrorsCannotClearChatVars) twiri18n.TranslationKey[KeysCommandsNukeErrorsCannotClearChatVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsNukeErrors struct {
	CannotCreateBroadcasterClient	KeysCommandsNukeErrorsCannotCreateBroadcasterClient
	CannotClearChat			KeysCommandsNukeErrorsCannotClearChat
}
type KeysCommandsNukeSuccessChatClearedVars struct {
}
type KeysCommandsNukeSuccessChatCleared struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsNukeSuccessChatCleared) IsTranslationKey() {
}
func (k KeysCommandsNukeSuccessChatCleared) GetPath() string {
	return "commands.nuke.success.chat_cleared"
}
func (k KeysCommandsNukeSuccessChatCleared) GetPathSlice() []string {
	return []string{"commands", "nuke", "success", "chat_cleared"}
}
func (k KeysCommandsNukeSuccessChatCleared) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsNukeSuccessChatCleared) SetVars(vars KeysCommandsNukeSuccessChatClearedVars) twiri18n.TranslationKey[KeysCommandsNukeSuccessChatClearedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsNukeSuccess struct {
	ChatCleared KeysCommandsNukeSuccessChatCleared
}
type KeysCommandsNuke struct {
	Description	KeysCommandsNukeDescription
	Errors		KeysCommandsNukeErrors
	Success		KeysCommandsNukeSuccess
}
type KeysCommands struct {
	Channel			KeysCommandsChannel
	Clip			KeysCommandsClip
	Dota			KeysCommandsDota
	Dudes			KeysCommandsDudes
	Marker			KeysCommandsMarker
	Song			KeysCommandsSong
	Stats			KeysCommandsStats
	Subage			KeysCommandsSubage
	Vips			KeysCommandsVips
	Shoutout		KeysCommandsShoutout
	Tts			KeysCommandsTts
	ChatWall		KeysCommandsChatWall
	Shorturl		KeysCommandsShorturl
	CategoriesAliases	KeysCommandsCategoriesAliases
	Overlays		KeysCommandsOverlays
	Permit			KeysCommandsPermit
	Games			KeysCommandsGames
	Seventv			KeysCommandsSeventv
	Prefix			KeysCommandsPrefix
	Giveaways		KeysCommandsGiveaways
	Songrequest		KeysCommandsSongrequest
	Utility			KeysCommandsUtility
	Predictions		KeysCommandsPredictions
	Spam			KeysCommandsSpam
	Tracer			KeysCommandsTracer
	Manage			KeysCommandsManage
	Nuke			KeysCommandsNuke
}
type KeysErrorsGenericBroadcasterClientVars struct {
}
type KeysErrorsGenericBroadcasterClient struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericBroadcasterClient) IsTranslationKey() {
}
func (k KeysErrorsGenericBroadcasterClient) GetPath() string {
	return "errors.generic.broadcaster_client"
}
func (k KeysErrorsGenericBroadcasterClient) GetPathSlice() []string {
	return []string{"errors", "generic", "broadcaster_client"}
}
func (k KeysErrorsGenericBroadcasterClient) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericBroadcasterClient) SetVars(vars KeysErrorsGenericBroadcasterClientVars) twiri18n.TranslationKey[KeysErrorsGenericBroadcasterClientVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericShouldMentionWithAtVars struct {
}
type KeysErrorsGenericShouldMentionWithAt struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericShouldMentionWithAt) IsTranslationKey() {
}
func (k KeysErrorsGenericShouldMentionWithAt) GetPath() string {
	return "errors.generic.should_mention_with_at"
}
func (k KeysErrorsGenericShouldMentionWithAt) GetPathSlice() []string {
	return []string{"errors", "generic", "should_mention_with_at"}
}
func (k KeysErrorsGenericShouldMentionWithAt) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericShouldMentionWithAt) SetVars(vars KeysErrorsGenericShouldMentionWithAtVars) twiri18n.TranslationKey[KeysErrorsGenericShouldMentionWithAtVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericCannotFindUserDbVars struct {
}
type KeysErrorsGenericCannotFindUserDb struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotFindUserDb) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotFindUserDb) GetPath() string {
	return "errors.generic.cannot_find_user_db"
}
func (k KeysErrorsGenericCannotFindUserDb) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_find_user_db"}
}
func (k KeysErrorsGenericCannotFindUserDb) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotFindUserDb) SetVars(vars KeysErrorsGenericCannotFindUserDbVars) twiri18n.TranslationKey[KeysErrorsGenericCannotFindUserDbVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericCannotFindUserTwitchVars struct {
}
type KeysErrorsGenericCannotFindUserTwitch struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotFindUserTwitch) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotFindUserTwitch) GetPath() string {
	return "errors.generic.cannot_find_user_twitch"
}
func (k KeysErrorsGenericCannotFindUserTwitch) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_find_user_twitch"}
}
func (k KeysErrorsGenericCannotFindUserTwitch) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotFindUserTwitch) SetVars(vars KeysErrorsGenericCannotFindUserTwitchVars) twiri18n.TranslationKey[KeysErrorsGenericCannotFindUserTwitchVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericCannotFindUsersTwitchVars struct {
}
type KeysErrorsGenericCannotFindUsersTwitch struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotFindUsersTwitch) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotFindUsersTwitch) GetPath() string {
	return "errors.generic.cannot_find_users_twitch"
}
func (k KeysErrorsGenericCannotFindUsersTwitch) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_find_users_twitch"}
}
func (k KeysErrorsGenericCannotFindUsersTwitch) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotFindUsersTwitch) SetVars(vars KeysErrorsGenericCannotFindUsersTwitchVars) twiri18n.TranslationKey[KeysErrorsGenericCannotFindUsersTwitchVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGeneric struct {
	BroadcasterClient	KeysErrorsGenericBroadcasterClient
	ShouldMentionWithAt	KeysErrorsGenericShouldMentionWithAt
	CannotFindUserDb	KeysErrorsGenericCannotFindUserDb
	CannotFindUserTwitch	KeysErrorsGenericCannotFindUserTwitch
	CannotFindUsersTwitch	KeysErrorsGenericCannotFindUsersTwitch
}
type KeysErrors struct {
	Generic KeysErrorsGeneric
}
type Keys struct {
	Commands	KeysCommands
	Errors		KeysErrors
}

var Translations = Keys{}
var Store twiri18n.LocalesStore = twiri18n.LocalesStore{"en": map[string]map[string]map[string]string{"commands": map[string]map[string]string{"stats": map[string]string{"description": `Channel statistics commands`, "followage.success.followage_display": `{user} has been following for {duration}`, "followage.description": `Check how long user has been following`, "followage.errors.cannot_get_followage": `Cannot get follow information`, "uptime.description": `Check stream uptime`, "followage.errors.user_not_following": `User is not following this channel`, "uptime.success.uptime_display": `Stream has been live for {duration}`, "uptime.errors.stream_offline": `Stream is currently offline`, "uptime.errors.cannot_get_uptime": `Cannot get stream uptime`, "viewers.errors.cannot_get_viewers": `Cannot get viewer count`, "viewers.success.viewers_display": `Current viewers: {count}`, "viewers.description": `Get current viewer count`}, "utility": map[string]string{"description": `Get first followers of the channel`, "errors.cannot_get_followers": `Cannot get first followers`, "success.first_followers": `First followers: {followers}`}, "categories_aliases": map[string]string{"aliases.errors.cannot_manage_aliases": `Cannot manage category aliases`, "aliases.success.alias_added": `✅ Category alias added: {alias}`, "aliases.success.alias_removed": `✅ Category alias removed: {alias}`, "set.description": `Set stream category`, "aliases.description": `Manage category aliases`, "aliases.errors.alias_not_found": `Alias not found`, "description": `Stream category and aliases management`, "set.errors.category_not_found": `Category not found`, "set.errors.cannot_set_category": `Cannot set stream category`, "set.success.category_set": `✅ Stream category set to: {category}`}, "channel": map[string]string{"title.errors.cannot_set_title": `Cannot set channel title`, "title.errors.title_too_long": `Title is too long`, "title.success.title_set": `✅ Channel title set to: {title}`, "game.errors.cannot_set_game": `Cannot set channel game`, "game.success.game_set": `✅ Channel game set to: {game}`, "description": `Channel management commands`, "title.description": `Set channel title`, "game.description": `Set channel game/category`, "game.errors.game_not_found": `Game not found`}, "chat_wall": map[string]string{"timeout.errors.cannot_timeout_user": `Cannot timeout user from chat wall`, "timeout.success.user_timed_out": `✅ User timed out from chat wall`, "stop.description": `Stop chat wall`, "stop.success.chat_wall_stopped": `✅ Chat wall stopped`, "ban.errors.cannot_ban_user": `Cannot ban user from chat wall`, "ban.success.user_banned": `✅ User banned from chat wall`, "timeout.description": `Timeout user from chat wall`, "stop.errors.cannot_stop_chat_wall": `Cannot stop chat wall`, "ban.description": `Ban user from chat wall`}, "dota": map[string]string{"description": `Dota 2 statistics and information`, "stats.errors.player_not_found": `Dota 2 player not found`, "stats.success.stats_display": `Dota 2 stats for {player}: {stats}`, "stats.description": `Get Dota 2 player statistics`, "stats.errors.cannot_get_stats": `Cannot get Dota 2 statistics`, "matches.success.matches_display": `Recent Dota 2 matches: {matches}`, "matches.description": `Get recent Dota 2 matches`, "matches.errors.cannot_get_matches": `Cannot get recent matches`}, "dudes": map[string]string{"play.description": `Play dudes game`, "play.success.game_started": `✅ Dudes game started`, "stop.errors.cannot_stop_game": `Cannot stop dudes game`, "stop.success.game_stopped": `✅ Dudes game stopped`, "stats.errors.cannot_get_stats": `Cannot get dudes game statistics`, "stats.success.stats_display": `Dudes game stats: {stats}`, "stats.description": `Show dudes game statistics`, "description": `Dudes game related commands`, "play.errors.cannot_start_game": `Cannot start dudes game`, "play.errors.game_already_running": `Dudes game is already running`, "stop.errors.no_game_running": `No dudes game is currently running`, "stop.description": `Stop dudes game`}, "giveaways": map[string]string{"reset.errors.cannot_reset_giveaway": `Cannot reset giveaway`, "delete.success.giveaway_deleted": `✅ Giveaway deleted`, "open.errors.cannot_open_giveaway": `Cannot open giveaway`, "reset.success.giveaway_reset": `✅ Giveaway reset`, "delete.description": `Delete giveaway`, "pause.success.giveaway_paused": `✅ Giveaway paused`, "pause.errors.cannot_pause_giveaway": `Cannot pause giveaway`, "finish.success.giveaway_finished": `✅ Giveaway finished`, "finish.errors.cannot_finish_giveaway": `Cannot finish giveaway`, "open.description": `Open giveaway`, "reset.description": `Reset giveaway`, "finish.description": `Finish giveaway`, "delete.errors.cannot_delete_giveaway": `Cannot delete giveaway`, "open.success.giveaway_opened": `✅ Giveaway opened`, "pause.description": `Pause giveaway`}, "manage": map[string]string{"add.errors.user_not_found": `User not found`, "remove.description": `Remove permission from user`, "remove.errors.cannot_remove_permission": `Cannot remove permission from user`, "remove.success.permission_removed": `✅ Removed permission {permission} from {user}`, "list.errors.cannot_list_permissions": `Cannot list user permissions`, "list.description": `List user permissions`, "add.errors.permission_already_exists": `User already has this permission`, "add.errors.cannot_add_permission": `Cannot add permission to user`, "add.success.permission_added": `✅ Added permission {permission} to {user}`, "remove.errors.permission_not_found": `User does not have this permission`, "list.success.permissions_list": `Permissions for {user}: {permissions}`, "description": `Manage user permissions`, "add.description": `Add permission to user`}, "song": map[string]string{"skip.errors.no_song_to_skip": `No song to skip`, "volume.errors.cannot_set_volume": `Cannot set music volume`, "description": `Music and song management commands`, "current.success.current_song": `Now playing: {song} by {artist}`, "skip.errors.cannot_skip_song": `Cannot skip current song`, "skip.success.song_skipped": `✅ Song skipped`, "volume.errors.invalid_volume": `Invalid volume level`, "volume.success.volume_set": `✅ Volume set to {volume}%`, "volume.description": `Set music volume`, "current.description": `Get current playing song`, "current.errors.no_song_playing": `No song currently playing`, "current.errors.cannot_get_current_song": `Cannot get current song`, "skip.description": `Skip current song`}, "clip": map[string]string{"success.clip_created": `✅ Clip created: {url}`, "description": `Create clip`, "errors.cannot_find_channel": `cannot find channel in db, please contact support`, "errors.cannot_create_broadcaster_client": `cannot create broadcaster twitch client`, "errors.cannot_create_clip": `cannot create clip`, "errors.cannot_get_clip": `cannot get clip`, "errors.empty_clip_url": `empty clip edit url`}, "permit": map[string]string{"remove.success.permit_removed": `✅ Removed link permit from {user}`, "description": `Permit user to post links`, "add.errors.cannot_add_permit": `Cannot add permit for user`, "add.errors.user_not_found": `User not found`, "add.success.permit_added": `✅ {user} can now post links`, "add.description": `Give user permission to post links`, "remove.description": `Remove user's link posting permission`, "remove.errors.cannot_remove_permit": `Cannot remove permit for user`}, "predictions": map[string]string{"cancel.errors.no_prediction_running": `no prediction runed`, "cancel.errors.cannot_create_twitch_client": `cannot create twitch client`, "start.description": `Start prediction. Example usage: !prediction start 100 | Will we win this game? | Yes, win / No, lose`, "start.errors.cannot_create_twitch_client": `cannot create twitch client`, "resolve.success.prediction_resolved": `✅ Prediction resolved with outcome: {outcome}`, "resolve.errors.cannot_get_prediction": `cannot get current prediction`, "lock.errors.cannot_lock_prediction": `cannot lock prediction`, "lock.errors.cannot_get_prediction": `cannot get current prediction`, "cancel.success.prediction_cancelled": `✅ Prediction cancelled`, "cancel.errors.cannot_cancel_prediction": `cannot cancel prediction`, "start.success.prediction_created": `✅ Prediction started: {title}`, "resolve.errors.no_prediction_running": `no prediction runed`, "resolve.description": `Resolve prediction`, "lock.success.prediction_locked": `✅ Prediction locked`, "cancel.description": `Cancel prediction`, "cancel.errors.cannot_get_prediction": `cannot get current prediction`, "start.errors.cannot_create_prediction": `cannot create prediction`, "resolve.errors.no_prediction_variant": `no prediction variant`, "resolve.errors.cannot_resolve_prediction": `cannot cancel prediction`, "resolve.errors.cannot_create_twitch_client": `cannot create twitch client`, "lock.description": `Lock prediction`, "lock.errors.cannot_create_twitch_client": `cannot create twitch client`, "lock.errors.no_prediction_running": `no prediction runed`}, "songrequest": map[string]string{"add.errors.song_already_in_queue": `Song already in queue`, "remove.errors.song_not_in_queue": `Song not in queue`, "remove.success.song_removed": `✅ Removed {song} from queue`, "queue.description": `Show current song request queue`, "queue.success.queue_display": `Song queue: {songs}`, "add.errors.queue_full": `Song request queue is full`, "add.success.song_added": `✅ Added {song} by {artist} to queue`, "add.description": `Add song to request queue`, "remove.description": `Remove song from request queue`, "remove.errors.cannot_remove_song": `Cannot remove song from queue`, "queue.errors.cannot_get_queue": `Cannot get song request queue`, "description": `Song request management commands`, "add.errors.cannot_add_song": `Cannot add song to queue`, "add.errors.song_not_found": `Song not found`, "queue.success.empty_queue": `Song request queue is empty`}, "spam": map[string]string{"timeout.errors.cannot_timeout_user": `Cannot timeout user for spam`, "settings.errors.invalid_settings": `Invalid spam settings`, "settings.success.settings_updated": `✅ Spam detection settings updated`, "description": `Anti-spam and moderation commands`, "timeout.errors.user_not_found": `User not found`, "timeout.success.user_timed_out": `✅ {user} timed out for spam`, "timeout.description": `Timeout user for spam`, "settings.description": `Configure spam detection settings`, "settings.errors.cannot_update_settings": `Cannot update spam settings`}, "tracer": map[string]string{"enable.description": `Enable debug tracing`, "disable.success.tracer_disabled": `✅ Debug tracer disabled`, "disable.description": `Disable debug tracing`, "description": `Debug and tracing commands`, "enable.errors.cannot_enable_tracer": `Cannot enable debug tracer`, "disable.errors.cannot_disable_tracer": `Cannot disable debug tracer`, "status.description": `Check tracer status`, "status.errors.cannot_get_status": `Cannot get tracer status`, "status.success.status_display": `Debug tracer status: {status}`, "enable.success.tracer_enabled": `✅ Debug tracer enabled`}, "vips": map[string]string{"removed": `✅ removed vip from {userName}`, "updated": `✅ updated vip for user {userName} new expriation time {endTime}`, "no_scheduled_vips": `There are no scheduled vips.`, "invalid_duration": `Invalid duration format. Please use formats like "1h", "30m", or "2d".`, "already_have_role": `User already vip or moderator!`, "added": `✅ added vip to {userName}`, "cannot_update": `Cannot update scheduled vip.`, "cannot_get_list_from_db": `Cannot get vip list from database.`, "cannot_create_scheduled_in_db": `Cannot create scheduled vip in database.`, "added_with_remove_time": `✅ added vip to {userName}, will be removed at {endTime}`}, "7tv": map[string]string{"add.errors.cannot_add_emote": `Cannot add 7TV emote`, "remove.errors.emote_not_found": `7TV emote not found`, "remove.success.emote_removed": `✅ Removed 7TV emote: {emote}`, "remove.description": `Remove 7TV emote`, "add.errors.emote_not_found": `7TV emote not found`, "add.errors.emote_already_exists": `Emote already exists`, "remove.errors.cannot_remove_emote": `Cannot remove 7TV emote`, "description": `Manage 7TV emotes`, "add.success.emote_added": `✅ Added 7TV emote: {emote}`, "add.description": `Add 7TV emote`}, "games": map[string]string{"description": `Game statistics and information`, "current.description": `Get current game information`, "current.errors.cannot_get_game": `Cannot get current game information`, "current.errors.no_game_set": `No game currently set`, "current.success.game_display": `Current game: {game}`, "stats.description": `Get game statistics`, "stats.errors.cannot_get_stats": `Cannot get game statistics`, "stats.success.stats_display": `Game stats: {stats}`}, "marker": map[string]string{"errors.cannot_create_broadcaster_client": `cannot create broadcaster twitch client`, "errors.cannot_create_marker": `cannot create marker`, "success.marker_created": `✅ Stream marker created`, "description": `Create a stream marker`, "errors.cannot_find_channel": `cannot find channel in db, please contact support`}, "nuke": map[string]string{"errors.cannot_create_broadcaster_client": `cannot create broadcaster twitch client`, "errors.cannot_clear_chat": `cannot clear chat`, "success.chat_cleared": `✅ Chat cleared`, "description": `Remove all messages from chat`}, "prefix": map[string]string{"remove.description": `Remove custom prefix (use default)`, "remove.errors.cannot_remove_prefix": `Cannot remove command prefix`, "remove.success.prefix_removed": `✅ Command prefix reset to default`, "set.success.prefix_set": `✅ Command prefix set to: {prefix}`, "description": `Set command prefix for the channel`, "set.description": `Set command prefix`, "set.errors.prefix_too_long": `Prefix is too long`, "set.errors.invalid_prefix": `Invalid prefix format`, "set.errors.cannot_set_prefix": `Cannot set command prefix`}, "shoutout": map[string]string{"content.description": `Promote content`, "description": `Promote a user or content`, "user.errors.user_not_found": `User not found`, "user.errors.cannot_shoutout_user": `Cannot shoutout user`, "user.success.user_promoted": `✅ Check out {user}! They were last seen playing {game}: {url}`, "user.description": `Shoutout a user`, "content.errors.cannot_promote_content": `Cannot promote content`, "content.success.content_promoted": `✅ Content promoted: {content}`}, "subage": map[string]string{"responses.time_remaining": `, {duration} remaining`, "responses.not_subscriber": `{user} is not a subscriber.`, "responses.not_subscriber_but_was": `{user} is not a subscriber, but used to be for {months} months.`, "responses.subscription_info": `{user} has a {tier} subscription to {channel} for a total {months} months`, "description": `Displays sub age of user or mentioned user.`, "errors.not_subscriber_or_hidden": `not a subscriber or info hidden`, "responses.streak_info": `, currently on a {months} months streak`}, "tts": map[string]string{"skip.success.tts_skipped": `✅ TTS message skipped`, "skip.description": `Skip current TTS message`, "skip.errors.no_tts_playing": `No TTS message playing`, "say.errors.cannot_use_tts": `Cannot use text-to-speech`, "say.errors.message_too_long": `Message is too long for TTS`, "say.errors.tts_disabled": `Text-to-speech is disabled`, "description": `Text-to-speech commands`, "say.success.tts_queued": `✅ TTS message queued`, "say.description": `Use text-to-speech to say message`, "skip.errors.cannot_skip_tts": `Cannot skip TTS message`}, "overlays": map[string]string{"description": `Overlay management commands`, "show.description": `Show overlay`, "show.errors.overlay_not_found": `Overlay not found`, "show.errors.cannot_show_overlay": `Cannot show overlay`, "show.success.overlay_shown": `✅ Overlay {overlay} shown`, "hide.errors.cannot_hide_overlay": `Cannot hide overlay`, "hide.success.overlay_hidden": `✅ Overlay {overlay} hidden`, "hide.description": `Hide overlay`}, "shorturl": map[string]string{"description": `Create short url`, "errors.cannot_create_short_url": `cannot create short url: {error}`, "success.short_url_created": `✅ {url}`}}, "errors": map[string]map[string]string{"generic": map[string]string{"should_mention_with_at": `you should tag user with @`, "cannot_find_user_db": `Cannot find user in database`, "cannot_find_user_twitch": `Cannot find user on twitch`, "cannot_find_users_twitch": `Cannot find users on twitch`, "broadcaster_client": `Cannot create broadcaster twitch client`}}}}
