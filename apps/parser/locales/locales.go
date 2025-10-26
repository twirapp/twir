package locales

import twiri18n "github.com/twirapp/twir/libs/i18n"

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

type KeysErrorsGenericCannotCreateTwitchVars struct {
}
type KeysErrorsGenericCannotCreateTwitch struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotCreateTwitch) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotCreateTwitch) GetPath() string {
	return "errors.generic.cannot_create_twitch"
}
func (k KeysErrorsGenericCannotCreateTwitch) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_create_twitch"}
}
func (k KeysErrorsGenericCannotCreateTwitch) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotCreateTwitch) SetVars(vars KeysErrorsGenericCannotCreateTwitchVars) twiri18n.TranslationKey[KeysErrorsGenericCannotCreateTwitchVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericCannotGetUserVars struct {
	Reason any
}
type KeysErrorsGenericCannotGetUser struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotGetUser) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotGetUser) GetPath() string {
	return "errors.generic.cannot_get_user"
}
func (k KeysErrorsGenericCannotGetUser) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_get_user"}
}
func (k KeysErrorsGenericCannotGetUser) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotGetUser) SetVars(vars KeysErrorsGenericCannotGetUserVars) twiri18n.TranslationKey[KeysErrorsGenericCannotGetUserVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysErrorsGenericCreateSettingsVars struct {
}
type KeysErrorsGenericCreateSettings struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCreateSettings) IsTranslationKey() {
}
func (k KeysErrorsGenericCreateSettings) GetPath() string {
	return "errors.generic.create_settings"
}
func (k KeysErrorsGenericCreateSettings) GetPathSlice() []string {
	return []string{"errors", "generic", "create_settings"}
}
func (k KeysErrorsGenericCreateSettings) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCreateSettings) SetVars(vars KeysErrorsGenericCreateSettingsVars) twiri18n.TranslationKey[KeysErrorsGenericCreateSettingsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericCannotFindChannelDbVars struct {
}
type KeysErrorsGenericCannotFindChannelDb struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotFindChannelDb) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotFindChannelDb) GetPath() string {
	return "errors.generic.cannot_find_channel_db"
}
func (k KeysErrorsGenericCannotFindChannelDb) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_find_channel_db"}
}
func (k KeysErrorsGenericCannotFindChannelDb) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotFindChannelDb) SetVars(vars KeysErrorsGenericCannotFindChannelDbVars) twiri18n.TranslationKey[KeysErrorsGenericCannotFindChannelDbVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericGettingChannelSettingsVars struct {
}
type KeysErrorsGenericGettingChannelSettings struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericGettingChannelSettings) IsTranslationKey() {
}
func (k KeysErrorsGenericGettingChannelSettings) GetPath() string {
	return "errors.generic.getting_channel_settings"
}
func (k KeysErrorsGenericGettingChannelSettings) GetPathSlice() []string {
	return []string{"errors", "generic", "getting_channel_settings"}
}
func (k KeysErrorsGenericGettingChannelSettings) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericGettingChannelSettings) SetVars(vars KeysErrorsGenericGettingChannelSettingsVars) twiri18n.TranslationKey[KeysErrorsGenericGettingChannelSettingsVars] {
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

type KeysErrorsGenericUpdatingSettingsVars struct {
}
type KeysErrorsGenericUpdatingSettings struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericUpdatingSettings) IsTranslationKey() {
}
func (k KeysErrorsGenericUpdatingSettings) GetPath() string {
	return "errors.generic.updating_settings"
}
func (k KeysErrorsGenericUpdatingSettings) GetPathSlice() []string {
	return []string{"errors", "generic", "updating_settings"}
}
func (k KeysErrorsGenericUpdatingSettings) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericUpdatingSettings) SetVars(vars KeysErrorsGenericUpdatingSettingsVars) twiri18n.TranslationKey[KeysErrorsGenericUpdatingSettingsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericCannotGetDbChannelVars struct {
}
type KeysErrorsGenericCannotGetDbChannel struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotGetDbChannel) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotGetDbChannel) GetPath() string {
	return "errors.generic.cannot_get_db_channel"
}
func (k KeysErrorsGenericCannotGetDbChannel) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_get_db_channel"}
}
func (k KeysErrorsGenericCannotGetDbChannel) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotGetDbChannel) SetVars(vars KeysErrorsGenericCannotGetDbChannelVars) twiri18n.TranslationKey[KeysErrorsGenericCannotGetDbChannelVars] {
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

type KeysErrorsGenericCannotGetStreamVars struct {
	Reason any
}
type KeysErrorsGenericCannotGetStream struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotGetStream) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotGetStream) GetPath() string {
	return "errors.generic.cannot_get_stream"
}
func (k KeysErrorsGenericCannotGetStream) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_get_stream"}
}
func (k KeysErrorsGenericCannotGetStream) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotGetStream) SetVars(vars KeysErrorsGenericCannotGetStreamVars) twiri18n.TranslationKey[KeysErrorsGenericCannotGetStreamVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysErrorsGenericCannotGetModeratorsVars struct {
	Reason any
}
type KeysErrorsGenericCannotGetModerators struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotGetModerators) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotGetModerators) GetPath() string {
	return "errors.generic.cannot_get_moderators"
}
func (k KeysErrorsGenericCannotGetModerators) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_get_moderators"}
}
func (k KeysErrorsGenericCannotGetModerators) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotGetModerators) SetVars(vars KeysErrorsGenericCannotGetModeratorsVars) twiri18n.TranslationKey[KeysErrorsGenericCannotGetModeratorsVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysErrorsGenericCannotGetAcceptCommandNameVars struct {
}
type KeysErrorsGenericCannotGetAcceptCommandName struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotGetAcceptCommandName) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotGetAcceptCommandName) GetPath() string {
	return "errors.generic.cannot_get_accept_command_name"
}
func (k KeysErrorsGenericCannotGetAcceptCommandName) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_get_accept_command_name"}
}
func (k KeysErrorsGenericCannotGetAcceptCommandName) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotGetAcceptCommandName) SetVars(vars KeysErrorsGenericCannotGetAcceptCommandNameVars) twiri18n.TranslationKey[KeysErrorsGenericCannotGetAcceptCommandNameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericCannotGetCommandVars struct {
}
type KeysErrorsGenericCannotGetCommand struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotGetCommand) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotGetCommand) GetPath() string {
	return "errors.generic.cannot_get_command"
}
func (k KeysErrorsGenericCannotGetCommand) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_get_command"}
}
func (k KeysErrorsGenericCannotGetCommand) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotGetCommand) SetVars(vars KeysErrorsGenericCannotGetCommandVars) twiri18n.TranslationKey[KeysErrorsGenericCannotGetCommandVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericSomethingWentWrongVars struct {
}
type KeysErrorsGenericSomethingWentWrong struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericSomethingWentWrong) IsTranslationKey() {
}
func (k KeysErrorsGenericSomethingWentWrong) GetPath() string {
	return "errors.generic.something_went_wrong"
}
func (k KeysErrorsGenericSomethingWentWrong) GetPathSlice() []string {
	return []string{"errors", "generic", "something_went_wrong"}
}
func (k KeysErrorsGenericSomethingWentWrong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericSomethingWentWrong) SetVars(vars KeysErrorsGenericSomethingWentWrongVars) twiri18n.TranslationKey[KeysErrorsGenericSomethingWentWrongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericTwirErrorVars struct {
}
type KeysErrorsGenericTwirError struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericTwirError) IsTranslationKey() {
}
func (k KeysErrorsGenericTwirError) GetPath() string {
	return "errors.generic.twir_error"
}
func (k KeysErrorsGenericTwirError) GetPathSlice() []string {
	return []string{"errors", "generic", "twir_error"}
}
func (k KeysErrorsGenericTwirError) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericTwirError) SetVars(vars KeysErrorsGenericTwirErrorVars) twiri18n.TranslationKey[KeysErrorsGenericTwirErrorVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericInternalVars struct {
}
type KeysErrorsGenericInternal struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericInternal) IsTranslationKey() {
}
func (k KeysErrorsGenericInternal) GetPath() string {
	return "errors.generic.internal"
}
func (k KeysErrorsGenericInternal) GetPathSlice() []string {
	return []string{"errors", "generic", "internal"}
}
func (k KeysErrorsGenericInternal) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericInternal) SetVars(vars KeysErrorsGenericInternalVars) twiri18n.TranslationKey[KeysErrorsGenericInternalVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericCannotTimeoutUserVars struct {
}
type KeysErrorsGenericCannotTimeoutUser struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotTimeoutUser) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotTimeoutUser) GetPath() string {
	return "errors.generic.cannot_timeout_user"
}
func (k KeysErrorsGenericCannotTimeoutUser) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_timeout_user"}
}
func (k KeysErrorsGenericCannotTimeoutUser) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotTimeoutUser) SetVars(vars KeysErrorsGenericCannotTimeoutUserVars) twiri18n.TranslationKey[KeysErrorsGenericCannotTimeoutUserVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericCannotBanUserVars struct {
}
type KeysErrorsGenericCannotBanUser struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotBanUser) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotBanUser) GetPath() string {
	return "errors.generic.cannot_ban_user"
}
func (k KeysErrorsGenericCannotBanUser) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_ban_user"}
}
func (k KeysErrorsGenericCannotBanUser) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotBanUser) SetVars(vars KeysErrorsGenericCannotBanUserVars) twiri18n.TranslationKey[KeysErrorsGenericCannotBanUserVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericNotAFollowerVars struct {
}
type KeysErrorsGenericNotAFollower struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericNotAFollower) IsTranslationKey() {
}
func (k KeysErrorsGenericNotAFollower) GetPath() string {
	return "errors.generic.not_a_follower"
}
func (k KeysErrorsGenericNotAFollower) GetPathSlice() []string {
	return []string{"errors", "generic", "not_a_follower"}
}
func (k KeysErrorsGenericNotAFollower) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericNotAFollower) SetVars(vars KeysErrorsGenericNotAFollowerVars) twiri18n.TranslationKey[KeysErrorsGenericNotAFollowerVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericCannotGetMessageVars struct {
}
type KeysErrorsGenericCannotGetMessage struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotGetMessage) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotGetMessage) GetPath() string {
	return "errors.generic.cannot_get_message"
}
func (k KeysErrorsGenericCannotGetMessage) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_get_message"}
}
func (k KeysErrorsGenericCannotGetMessage) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotGetMessage) SetVars(vars KeysErrorsGenericCannotGetMessageVars) twiri18n.TranslationKey[KeysErrorsGenericCannotGetMessageVars] {
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

type KeysErrorsGenericUserNotFoundVars struct {
}
type KeysErrorsGenericUserNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericUserNotFound) IsTranslationKey() {
}
func (k KeysErrorsGenericUserNotFound) GetPath() string {
	return "errors.generic.user_not_found"
}
func (k KeysErrorsGenericUserNotFound) GetPathSlice() []string {
	return []string{"errors", "generic", "user_not_found"}
}
func (k KeysErrorsGenericUserNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericUserNotFound) SetVars(vars KeysErrorsGenericUserNotFoundVars) twiri18n.TranslationKey[KeysErrorsGenericUserNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericGettingUserSettingsVars struct {
}
type KeysErrorsGenericGettingUserSettings struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericGettingUserSettings) IsTranslationKey() {
}
func (k KeysErrorsGenericGettingUserSettings) GetPath() string {
	return "errors.generic.getting_user_settings"
}
func (k KeysErrorsGenericGettingUserSettings) GetPathSlice() []string {
	return []string{"errors", "generic", "getting_user_settings"}
}
func (k KeysErrorsGenericGettingUserSettings) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericGettingUserSettings) SetVars(vars KeysErrorsGenericGettingUserSettingsVars) twiri18n.TranslationKey[KeysErrorsGenericGettingUserSettingsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGenericCannotFindChannelTwitchVars struct {
	Reason any
}
type KeysErrorsGenericCannotFindChannelTwitch struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotFindChannelTwitch) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotFindChannelTwitch) GetPath() string {
	return "errors.generic.cannot_find_channel_twitch"
}
func (k KeysErrorsGenericCannotFindChannelTwitch) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_find_channel_twitch"}
}
func (k KeysErrorsGenericCannotFindChannelTwitch) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotFindChannelTwitch) SetVars(vars KeysErrorsGenericCannotFindChannelTwitchVars) twiri18n.TranslationKey[KeysErrorsGenericCannotFindChannelTwitchVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysErrorsGenericCannotCreateCommandVars struct {
}
type KeysErrorsGenericCannotCreateCommand struct {
	Vars twiri18n.Vars
}

func (k KeysErrorsGenericCannotCreateCommand) IsTranslationKey() {
}
func (k KeysErrorsGenericCannotCreateCommand) GetPath() string {
	return "errors.generic.cannot_create_command"
}
func (k KeysErrorsGenericCannotCreateCommand) GetPathSlice() []string {
	return []string{"errors", "generic", "cannot_create_command"}
}
func (k KeysErrorsGenericCannotCreateCommand) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysErrorsGenericCannotCreateCommand) SetVars(vars KeysErrorsGenericCannotCreateCommandVars) twiri18n.TranslationKey[KeysErrorsGenericCannotCreateCommandVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysErrorsGeneric struct {
	BroadcasterClient		KeysErrorsGenericBroadcasterClient
	CannotCreateTwitch		KeysErrorsGenericCannotCreateTwitch
	CannotGetUser			KeysErrorsGenericCannotGetUser
	CreateSettings			KeysErrorsGenericCreateSettings
	CannotFindChannelDb		KeysErrorsGenericCannotFindChannelDb
	GettingChannelSettings		KeysErrorsGenericGettingChannelSettings
	CannotFindUserDb		KeysErrorsGenericCannotFindUserDb
	UpdatingSettings		KeysErrorsGenericUpdatingSettings
	CannotGetDbChannel		KeysErrorsGenericCannotGetDbChannel
	ShouldMentionWithAt		KeysErrorsGenericShouldMentionWithAt
	CannotGetStream			KeysErrorsGenericCannotGetStream
	CannotGetModerators		KeysErrorsGenericCannotGetModerators
	CannotGetAcceptCommandName	KeysErrorsGenericCannotGetAcceptCommandName
	CannotGetCommand		KeysErrorsGenericCannotGetCommand
	SomethingWentWrong		KeysErrorsGenericSomethingWentWrong
	TwirError			KeysErrorsGenericTwirError
	Internal			KeysErrorsGenericInternal
	CannotTimeoutUser		KeysErrorsGenericCannotTimeoutUser
	CannotBanUser			KeysErrorsGenericCannotBanUser
	NotAFollower			KeysErrorsGenericNotAFollower
	CannotGetMessage		KeysErrorsGenericCannotGetMessage
	CannotFindUserTwitch		KeysErrorsGenericCannotFindUserTwitch
	CannotFindUsersTwitch		KeysErrorsGenericCannotFindUsersTwitch
	UserNotFound			KeysErrorsGenericUserNotFound
	GettingUserSettings		KeysErrorsGenericGettingUserSettings
	CannotFindChannelTwitch		KeysErrorsGenericCannotFindChannelTwitch
	CannotCreateCommand		KeysErrorsGenericCannotCreateCommand
}
type KeysErrors struct {
	Generic KeysErrorsGeneric
}
type KeysServicesChatWallInfoBannedByTwirVars struct {
	BanPhrase any
}
type KeysServicesChatWallInfoBannedByTwir struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallInfoBannedByTwir) IsTranslationKey() {
}
func (k KeysServicesChatWallInfoBannedByTwir) GetPath() string {
	return "services.chat_wall.info.banned_by_twir"
}
func (k KeysServicesChatWallInfoBannedByTwir) GetPathSlice() []string {
	return []string{"services", "chat_wall", "info", "banned_by_twir"}
}
func (k KeysServicesChatWallInfoBannedByTwir) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallInfoBannedByTwir) SetVars(vars KeysServicesChatWallInfoBannedByTwirVars) twiri18n.TranslationKey[KeysServicesChatWallInfoBannedByTwirVars] {
	k.Vars = twiri18n.Vars{"banPhrase": vars.BanPhrase}
	return k
}

type KeysServicesChatWallInfo struct {
	BannedByTwir KeysServicesChatWallInfoBannedByTwir
}
type KeysServicesChatWallErrorsGetChatWallSettingsVars struct {
	Reason any
}
type KeysServicesChatWallErrorsGetChatWallSettings struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsGetChatWallSettings) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsGetChatWallSettings) GetPath() string {
	return "services.chat_wall.errors.get_chat_wall_settings"
}
func (k KeysServicesChatWallErrorsGetChatWallSettings) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "get_chat_wall_settings"}
}
func (k KeysServicesChatWallErrorsGetChatWallSettings) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsGetChatWallSettings) SetVars(vars KeysServicesChatWallErrorsGetChatWallSettingsVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsGetChatWallSettingsVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysServicesChatWallErrorsGetUsersStatsVars struct {
	Reason any
}
type KeysServicesChatWallErrorsGetUsersStats struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsGetUsersStats) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsGetUsersStats) GetPath() string {
	return "services.chat_wall.errors.get_users_stats"
}
func (k KeysServicesChatWallErrorsGetUsersStats) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "get_users_stats"}
}
func (k KeysServicesChatWallErrorsGetUsersStats) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsGetUsersStats) SetVars(vars KeysServicesChatWallErrorsGetUsersStatsVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsGetUsersStatsVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysServicesChatWallErrorsPublishDeletedMessagesVars struct {
	Reason any
}
type KeysServicesChatWallErrorsPublishDeletedMessages struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsPublishDeletedMessages) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsPublishDeletedMessages) GetPath() string {
	return "services.chat_wall.errors.publish_deleted_messages"
}
func (k KeysServicesChatWallErrorsPublishDeletedMessages) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "publish_deleted_messages"}
}
func (k KeysServicesChatWallErrorsPublishDeletedMessages) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsPublishDeletedMessages) SetVars(vars KeysServicesChatWallErrorsPublishDeletedMessagesVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsPublishDeletedMessagesVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysServicesChatWallErrorsPublishBanUsersVars struct {
	Reason any
}
type KeysServicesChatWallErrorsPublishBanUsers struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsPublishBanUsers) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsPublishBanUsers) GetPath() string {
	return "services.chat_wall.errors.publish_ban_users"
}
func (k KeysServicesChatWallErrorsPublishBanUsers) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "publish_ban_users"}
}
func (k KeysServicesChatWallErrorsPublishBanUsers) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsPublishBanUsers) SetVars(vars KeysServicesChatWallErrorsPublishBanUsersVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsPublishBanUsersVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysServicesChatWallErrorsCreateChatWallVars struct {
	Reason any
}
type KeysServicesChatWallErrorsCreateChatWall struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsCreateChatWall) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsCreateChatWall) GetPath() string {
	return "services.chat_wall.errors.create_chat_wall"
}
func (k KeysServicesChatWallErrorsCreateChatWall) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "create_chat_wall"}
}
func (k KeysServicesChatWallErrorsCreateChatWall) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsCreateChatWall) SetVars(vars KeysServicesChatWallErrorsCreateChatWallVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsCreateChatWallVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysServicesChatWallErrorsHandledMessagesToRedisVars struct {
	Reason any
}
type KeysServicesChatWallErrorsHandledMessagesToRedis struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsHandledMessagesToRedis) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsHandledMessagesToRedis) GetPath() string {
	return "services.chat_wall.errors.handled_messages_to_redis"
}
func (k KeysServicesChatWallErrorsHandledMessagesToRedis) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "handled_messages_to_redis"}
}
func (k KeysServicesChatWallErrorsHandledMessagesToRedis) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsHandledMessagesToRedis) SetVars(vars KeysServicesChatWallErrorsHandledMessagesToRedisVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsHandledMessagesToRedisVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysServicesChatWallErrorsChatWallNotFoundVars struct {
}
type KeysServicesChatWallErrorsChatWallNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsChatWallNotFound) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsChatWallNotFound) GetPath() string {
	return "services.chat_wall.errors.chat_wall_not_found"
}
func (k KeysServicesChatWallErrorsChatWallNotFound) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "chat_wall_not_found"}
}
func (k KeysServicesChatWallErrorsChatWallNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsChatWallNotFound) SetVars(vars KeysServicesChatWallErrorsChatWallNotFoundVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsChatWallNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysServicesChatWallErrorsGetChatWallsVars struct {
	Reason any
}
type KeysServicesChatWallErrorsGetChatWalls struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsGetChatWalls) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsGetChatWalls) GetPath() string {
	return "services.chat_wall.errors.get_chat_walls"
}
func (k KeysServicesChatWallErrorsGetChatWalls) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "get_chat_walls"}
}
func (k KeysServicesChatWallErrorsGetChatWalls) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsGetChatWalls) SetVars(vars KeysServicesChatWallErrorsGetChatWallsVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsGetChatWallsVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysServicesChatWallErrorsUpdateChatWallsVars struct {
	Reason any
}
type KeysServicesChatWallErrorsUpdateChatWalls struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsUpdateChatWalls) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsUpdateChatWalls) GetPath() string {
	return "services.chat_wall.errors.update_chat_walls"
}
func (k KeysServicesChatWallErrorsUpdateChatWalls) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "update_chat_walls"}
}
func (k KeysServicesChatWallErrorsUpdateChatWalls) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsUpdateChatWalls) SetVars(vars KeysServicesChatWallErrorsUpdateChatWallsVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsUpdateChatWallsVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysServicesChatWallErrorsGetAlreadyHandledVars struct {
	Reason any
}
type KeysServicesChatWallErrorsGetAlreadyHandled struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsGetAlreadyHandled) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsGetAlreadyHandled) GetPath() string {
	return "services.chat_wall.errors.get_already_handled"
}
func (k KeysServicesChatWallErrorsGetAlreadyHandled) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "get_already_handled"}
}
func (k KeysServicesChatWallErrorsGetAlreadyHandled) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsGetAlreadyHandled) SetVars(vars KeysServicesChatWallErrorsGetAlreadyHandledVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsGetAlreadyHandledVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysServicesChatWallErrorsCreateChatWallWithPhraseVars struct {
}
type KeysServicesChatWallErrorsCreateChatWallWithPhrase struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsCreateChatWallWithPhrase) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsCreateChatWallWithPhrase) GetPath() string {
	return "services.chat_wall.errors.create_chat_wall_with_phrase"
}
func (k KeysServicesChatWallErrorsCreateChatWallWithPhrase) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "create_chat_wall_with_phrase"}
}
func (k KeysServicesChatWallErrorsCreateChatWallWithPhrase) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsCreateChatWallWithPhrase) SetVars(vars KeysServicesChatWallErrorsCreateChatWallWithPhraseVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsCreateChatWallWithPhraseVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysServicesChatWallErrorsCreateChatLogsInDbVars struct {
	Reason any
}
type KeysServicesChatWallErrorsCreateChatLogsInDb struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsCreateChatLogsInDb) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsCreateChatLogsInDb) GetPath() string {
	return "services.chat_wall.errors.create_chat_logs_in_db"
}
func (k KeysServicesChatWallErrorsCreateChatLogsInDb) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "create_chat_logs_in_db"}
}
func (k KeysServicesChatWallErrorsCreateChatLogsInDb) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsCreateChatLogsInDb) SetVars(vars KeysServicesChatWallErrorsCreateChatLogsInDbVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsCreateChatLogsInDbVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysServicesChatWallErrorsGetCurrentChatWallsVars struct {
	Reason any
}
type KeysServicesChatWallErrorsGetCurrentChatWalls struct {
	Vars twiri18n.Vars
}

func (k KeysServicesChatWallErrorsGetCurrentChatWalls) IsTranslationKey() {
}
func (k KeysServicesChatWallErrorsGetCurrentChatWalls) GetPath() string {
	return "services.chat_wall.errors.get_current_chat_walls"
}
func (k KeysServicesChatWallErrorsGetCurrentChatWalls) GetPathSlice() []string {
	return []string{"services", "chat_wall", "errors", "get_current_chat_walls"}
}
func (k KeysServicesChatWallErrorsGetCurrentChatWalls) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesChatWallErrorsGetCurrentChatWalls) SetVars(vars KeysServicesChatWallErrorsGetCurrentChatWallsVars) twiri18n.TranslationKey[KeysServicesChatWallErrorsGetCurrentChatWallsVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysServicesChatWallErrors struct {
	GetChatWallSettings		KeysServicesChatWallErrorsGetChatWallSettings
	GetUsersStats			KeysServicesChatWallErrorsGetUsersStats
	PublishDeletedMessages		KeysServicesChatWallErrorsPublishDeletedMessages
	PublishBanUsers			KeysServicesChatWallErrorsPublishBanUsers
	CreateChatWall			KeysServicesChatWallErrorsCreateChatWall
	HandledMessagesToRedis		KeysServicesChatWallErrorsHandledMessagesToRedis
	ChatWallNotFound		KeysServicesChatWallErrorsChatWallNotFound
	GetChatWalls			KeysServicesChatWallErrorsGetChatWalls
	UpdateChatWalls			KeysServicesChatWallErrorsUpdateChatWalls
	GetAlreadyHandled		KeysServicesChatWallErrorsGetAlreadyHandled
	CreateChatWallWithPhrase	KeysServicesChatWallErrorsCreateChatWallWithPhrase
	CreateChatLogsInDb		KeysServicesChatWallErrorsCreateChatLogsInDb
	GetCurrentChatWalls		KeysServicesChatWallErrorsGetCurrentChatWalls
}
type KeysServicesChatWall struct {
	Info	KeysServicesChatWallInfo
	Errors	KeysServicesChatWallErrors
}
type KeysServicesShortenedurlsErrorsInvalidUrlVars struct {
}
type KeysServicesShortenedurlsErrorsInvalidUrl struct {
	Vars twiri18n.Vars
}

func (k KeysServicesShortenedurlsErrorsInvalidUrl) IsTranslationKey() {
}
func (k KeysServicesShortenedurlsErrorsInvalidUrl) GetPath() string {
	return "services.shortenedurls.errors.invalid_url"
}
func (k KeysServicesShortenedurlsErrorsInvalidUrl) GetPathSlice() []string {
	return []string{"services", "shortenedurls", "errors", "invalid_url"}
}
func (k KeysServicesShortenedurlsErrorsInvalidUrl) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesShortenedurlsErrorsInvalidUrl) SetVars(vars KeysServicesShortenedurlsErrorsInvalidUrlVars) twiri18n.TranslationKey[KeysServicesShortenedurlsErrorsInvalidUrlVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysServicesShortenedurlsErrors struct {
	InvalidUrl KeysServicesShortenedurlsErrorsInvalidUrl
}
type KeysServicesShortenedurls struct {
	Errors KeysServicesShortenedurlsErrors
}
type KeysServicesTtsInfoNotConfiguredVars struct {
}
type KeysServicesTtsInfoNotConfigured struct {
	Vars twiri18n.Vars
}

func (k KeysServicesTtsInfoNotConfigured) IsTranslationKey() {
}
func (k KeysServicesTtsInfoNotConfigured) GetPath() string {
	return "services.tts.info.not_configured"
}
func (k KeysServicesTtsInfoNotConfigured) GetPathSlice() []string {
	return []string{"services", "tts", "info", "not_configured"}
}
func (k KeysServicesTtsInfoNotConfigured) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesTtsInfoNotConfigured) SetVars(vars KeysServicesTtsInfoNotConfiguredVars) twiri18n.TranslationKey[KeysServicesTtsInfoNotConfiguredVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysServicesTtsInfoNoVoicesVars struct {
}
type KeysServicesTtsInfoNoVoices struct {
	Vars twiri18n.Vars
}

func (k KeysServicesTtsInfoNoVoices) IsTranslationKey() {
}
func (k KeysServicesTtsInfoNoVoices) GetPath() string {
	return "services.tts.info.no_voices"
}
func (k KeysServicesTtsInfoNoVoices) GetPathSlice() []string {
	return []string{"services", "tts", "info", "no_voices"}
}
func (k KeysServicesTtsInfoNoVoices) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesTtsInfoNoVoices) SetVars(vars KeysServicesTtsInfoNoVoicesVars) twiri18n.TranslationKey[KeysServicesTtsInfoNoVoicesVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysServicesTtsInfo struct {
	NotConfigured	KeysServicesTtsInfoNotConfigured
	NoVoices	KeysServicesTtsInfoNoVoices
}
type KeysServicesTtsErrorsVoiceDisallowedVars struct {
	UserVoice any
}
type KeysServicesTtsErrorsVoiceDisallowed struct {
	Vars twiri18n.Vars
}

func (k KeysServicesTtsErrorsVoiceDisallowed) IsTranslationKey() {
}
func (k KeysServicesTtsErrorsVoiceDisallowed) GetPath() string {
	return "services.tts.errors.voice_disallowed"
}
func (k KeysServicesTtsErrorsVoiceDisallowed) GetPathSlice() []string {
	return []string{"services", "tts", "errors", "voice_disallowed"}
}
func (k KeysServicesTtsErrorsVoiceDisallowed) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesTtsErrorsVoiceDisallowed) SetVars(vars KeysServicesTtsErrorsVoiceDisallowedVars) twiri18n.TranslationKey[KeysServicesTtsErrorsVoiceDisallowedVars] {
	k.Vars = twiri18n.Vars{"userVoice": vars.UserVoice}
	return k
}

type KeysServicesTtsErrorsNotFoundVars struct {
	UserVoice any
}
type KeysServicesTtsErrorsNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysServicesTtsErrorsNotFound) IsTranslationKey() {
}
func (k KeysServicesTtsErrorsNotFound) GetPath() string {
	return "services.tts.errors.not_found"
}
func (k KeysServicesTtsErrorsNotFound) GetPathSlice() []string {
	return []string{"services", "tts", "errors", "not_found"}
}
func (k KeysServicesTtsErrorsNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysServicesTtsErrorsNotFound) SetVars(vars KeysServicesTtsErrorsNotFoundVars) twiri18n.TranslationKey[KeysServicesTtsErrorsNotFoundVars] {
	k.Vars = twiri18n.Vars{"userVoice": vars.UserVoice}
	return k
}

type KeysServicesTtsErrors struct {
	VoiceDisallowed	KeysServicesTtsErrorsVoiceDisallowed
	NotFound	KeysServicesTtsErrorsNotFound
}
type KeysServicesTts struct {
	Info	KeysServicesTtsInfo
	Errors	KeysServicesTtsErrors
}
type KeysServices struct {
	ChatWall	KeysServicesChatWall
	Shortenedurls	KeysServicesShortenedurls
	Tts		KeysServicesTts
}
type KeysVariablesCommandsInfoNoPassedParamsVars struct {
}
type KeysVariablesCommandsInfoNoPassedParams struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesCommandsInfoNoPassedParams) IsTranslationKey() {
}
func (k KeysVariablesCommandsInfoNoPassedParams) GetPath() string {
	return "variables.commands.info.no_passed_params"
}
func (k KeysVariablesCommandsInfoNoPassedParams) GetPathSlice() []string {
	return []string{"variables", "commands", "info", "no_passed_params"}
}
func (k KeysVariablesCommandsInfoNoPassedParams) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesCommandsInfoNoPassedParams) SetVars(vars KeysVariablesCommandsInfoNoPassedParamsVars) twiri18n.TranslationKey[KeysVariablesCommandsInfoNoPassedParamsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesCommandsInfoCommandWithNameNotFoundVars struct {
	CommandName any
}
type KeysVariablesCommandsInfoCommandWithNameNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesCommandsInfoCommandWithNameNotFound) IsTranslationKey() {
}
func (k KeysVariablesCommandsInfoCommandWithNameNotFound) GetPath() string {
	return "variables.commands.info.command_with_name_not_found"
}
func (k KeysVariablesCommandsInfoCommandWithNameNotFound) GetPathSlice() []string {
	return []string{"variables", "commands", "info", "command_with_name_not_found"}
}
func (k KeysVariablesCommandsInfoCommandWithNameNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesCommandsInfoCommandWithNameNotFound) SetVars(vars KeysVariablesCommandsInfoCommandWithNameNotFoundVars) twiri18n.TranslationKey[KeysVariablesCommandsInfoCommandWithNameNotFoundVars] {
	k.Vars = twiri18n.Vars{"commandName": vars.CommandName}
	return k
}

type KeysVariablesCommandsInfoGetCountVars struct {
}
type KeysVariablesCommandsInfoGetCount struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesCommandsInfoGetCount) IsTranslationKey() {
}
func (k KeysVariablesCommandsInfoGetCount) GetPath() string {
	return "variables.commands.info.get_count"
}
func (k KeysVariablesCommandsInfoGetCount) GetPathSlice() []string {
	return []string{"variables", "commands", "info", "get_count"}
}
func (k KeysVariablesCommandsInfoGetCount) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesCommandsInfoGetCount) SetVars(vars KeysVariablesCommandsInfoGetCountVars) twiri18n.TranslationKey[KeysVariablesCommandsInfoGetCountVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesCommandsInfo struct {
	NoPassedParams		KeysVariablesCommandsInfoNoPassedParams
	CommandWithNameNotFound	KeysVariablesCommandsInfoCommandWithNameNotFound
	GetCount		KeysVariablesCommandsInfoGetCount
}
type KeysVariablesCommands struct {
	Info KeysVariablesCommandsInfo
}
type KeysVariablesRandomErrorsEmptyPhraseVars struct {
}
type KeysVariablesRandomErrorsEmptyPhrase struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesRandomErrorsEmptyPhrase) IsTranslationKey() {
}
func (k KeysVariablesRandomErrorsEmptyPhrase) GetPath() string {
	return "variables.random.errors.empty_phrase"
}
func (k KeysVariablesRandomErrorsEmptyPhrase) GetPathSlice() []string {
	return []string{"variables", "random", "errors", "empty_phrase"}
}
func (k KeysVariablesRandomErrorsEmptyPhrase) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesRandomErrorsEmptyPhrase) SetVars(vars KeysVariablesRandomErrorsEmptyPhraseVars) twiri18n.TranslationKey[KeysVariablesRandomErrorsEmptyPhraseVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesRandomErrorsGetOnlineUserVars struct {
}
type KeysVariablesRandomErrorsGetOnlineUser struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesRandomErrorsGetOnlineUser) IsTranslationKey() {
}
func (k KeysVariablesRandomErrorsGetOnlineUser) GetPath() string {
	return "variables.random.errors.get_online_user"
}
func (k KeysVariablesRandomErrorsGetOnlineUser) GetPathSlice() []string {
	return []string{"variables", "random", "errors", "get_online_user"}
}
func (k KeysVariablesRandomErrorsGetOnlineUser) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesRandomErrorsGetOnlineUser) SetVars(vars KeysVariablesRandomErrorsGetOnlineUserVars) twiri18n.TranslationKey[KeysVariablesRandomErrorsGetOnlineUserVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesRandomErrorsNotPassedParamsVars struct {
}
type KeysVariablesRandomErrorsNotPassedParams struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesRandomErrorsNotPassedParams) IsTranslationKey() {
}
func (k KeysVariablesRandomErrorsNotPassedParams) GetPath() string {
	return "variables.random.errors.not_passed_params"
}
func (k KeysVariablesRandomErrorsNotPassedParams) GetPathSlice() []string {
	return []string{"variables", "random", "errors", "not_passed_params"}
}
func (k KeysVariablesRandomErrorsNotPassedParams) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesRandomErrorsNotPassedParams) SetVars(vars KeysVariablesRandomErrorsNotPassedParamsVars) twiri18n.TranslationKey[KeysVariablesRandomErrorsNotPassedParamsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesRandomErrorsFirstLargerSecondVars struct {
}
type KeysVariablesRandomErrorsFirstLargerSecond struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesRandomErrorsFirstLargerSecond) IsTranslationKey() {
}
func (k KeysVariablesRandomErrorsFirstLargerSecond) GetPath() string {
	return "variables.random.errors.first_larger_second"
}
func (k KeysVariablesRandomErrorsFirstLargerSecond) GetPathSlice() []string {
	return []string{"variables", "random", "errors", "first_larger_second"}
}
func (k KeysVariablesRandomErrorsFirstLargerSecond) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesRandomErrorsFirstLargerSecond) SetVars(vars KeysVariablesRandomErrorsFirstLargerSecondVars) twiri18n.TranslationKey[KeysVariablesRandomErrorsFirstLargerSecondVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesRandomErrorsParametersNotSpecifiedVars struct {
}
type KeysVariablesRandomErrorsParametersNotSpecified struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesRandomErrorsParametersNotSpecified) IsTranslationKey() {
}
func (k KeysVariablesRandomErrorsParametersNotSpecified) GetPath() string {
	return "variables.random.errors.parameters_not_specified"
}
func (k KeysVariablesRandomErrorsParametersNotSpecified) GetPathSlice() []string {
	return []string{"variables", "random", "errors", "parameters_not_specified"}
}
func (k KeysVariablesRandomErrorsParametersNotSpecified) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesRandomErrorsParametersNotSpecified) SetVars(vars KeysVariablesRandomErrorsParametersNotSpecifiedVars) twiri18n.TranslationKey[KeysVariablesRandomErrorsParametersNotSpecifiedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesRandomErrorsWrongWithParamsVars struct {
}
type KeysVariablesRandomErrorsWrongWithParams struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesRandomErrorsWrongWithParams) IsTranslationKey() {
}
func (k KeysVariablesRandomErrorsWrongWithParams) GetPath() string {
	return "variables.random.errors.wrong_with_params"
}
func (k KeysVariablesRandomErrorsWrongWithParams) GetPathSlice() []string {
	return []string{"variables", "random", "errors", "wrong_with_params"}
}
func (k KeysVariablesRandomErrorsWrongWithParams) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesRandomErrorsWrongWithParams) SetVars(vars KeysVariablesRandomErrorsWrongWithParamsVars) twiri18n.TranslationKey[KeysVariablesRandomErrorsWrongWithParamsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesRandomErrorsWrongNumberVars struct {
}
type KeysVariablesRandomErrorsWrongNumber struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesRandomErrorsWrongNumber) IsTranslationKey() {
}
func (k KeysVariablesRandomErrorsWrongNumber) GetPath() string {
	return "variables.random.errors.wrong_number"
}
func (k KeysVariablesRandomErrorsWrongNumber) GetPathSlice() []string {
	return []string{"variables", "random", "errors", "wrong_number"}
}
func (k KeysVariablesRandomErrorsWrongNumber) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesRandomErrorsWrongNumber) SetVars(vars KeysVariablesRandomErrorsWrongNumberVars) twiri18n.TranslationKey[KeysVariablesRandomErrorsWrongNumberVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesRandomErrorsParseFirstNumberVars struct {
}
type KeysVariablesRandomErrorsParseFirstNumber struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesRandomErrorsParseFirstNumber) IsTranslationKey() {
}
func (k KeysVariablesRandomErrorsParseFirstNumber) GetPath() string {
	return "variables.random.errors.parse_first_number"
}
func (k KeysVariablesRandomErrorsParseFirstNumber) GetPathSlice() []string {
	return []string{"variables", "random", "errors", "parse_first_number"}
}
func (k KeysVariablesRandomErrorsParseFirstNumber) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesRandomErrorsParseFirstNumber) SetVars(vars KeysVariablesRandomErrorsParseFirstNumberVars) twiri18n.TranslationKey[KeysVariablesRandomErrorsParseFirstNumberVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesRandomErrorsParseSecondNumberVars struct {
}
type KeysVariablesRandomErrorsParseSecondNumber struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesRandomErrorsParseSecondNumber) IsTranslationKey() {
}
func (k KeysVariablesRandomErrorsParseSecondNumber) GetPath() string {
	return "variables.random.errors.parse_second_number"
}
func (k KeysVariablesRandomErrorsParseSecondNumber) GetPathSlice() []string {
	return []string{"variables", "random", "errors", "parse_second_number"}
}
func (k KeysVariablesRandomErrorsParseSecondNumber) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesRandomErrorsParseSecondNumber) SetVars(vars KeysVariablesRandomErrorsParseSecondNumberVars) twiri18n.TranslationKey[KeysVariablesRandomErrorsParseSecondNumberVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesRandomErrorsLowerNumbersVars struct {
}
type KeysVariablesRandomErrorsLowerNumbers struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesRandomErrorsLowerNumbers) IsTranslationKey() {
}
func (k KeysVariablesRandomErrorsLowerNumbers) GetPath() string {
	return "variables.random.errors.lower_numbers"
}
func (k KeysVariablesRandomErrorsLowerNumbers) GetPathSlice() []string {
	return []string{"variables", "random", "errors", "lower_numbers"}
}
func (k KeysVariablesRandomErrorsLowerNumbers) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesRandomErrorsLowerNumbers) SetVars(vars KeysVariablesRandomErrorsLowerNumbersVars) twiri18n.TranslationKey[KeysVariablesRandomErrorsLowerNumbersVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesRandomErrors struct {
	EmptyPhrase		KeysVariablesRandomErrorsEmptyPhrase
	GetOnlineUser		KeysVariablesRandomErrorsGetOnlineUser
	NotPassedParams		KeysVariablesRandomErrorsNotPassedParams
	FirstLargerSecond	KeysVariablesRandomErrorsFirstLargerSecond
	ParametersNotSpecified	KeysVariablesRandomErrorsParametersNotSpecified
	WrongWithParams		KeysVariablesRandomErrorsWrongWithParams
	WrongNumber		KeysVariablesRandomErrorsWrongNumber
	ParseFirstNumber	KeysVariablesRandomErrorsParseFirstNumber
	ParseSecondNumber	KeysVariablesRandomErrorsParseSecondNumber
	LowerNumbers		KeysVariablesRandomErrorsLowerNumbers
}
type KeysVariablesRandom struct {
	Errors KeysVariablesRandomErrors
}
type KeysVariablesShorturlErrorsCreateShortUrlVars struct {
	Reason any
}
type KeysVariablesShorturlErrorsCreateShortUrl struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesShorturlErrorsCreateShortUrl) IsTranslationKey() {
}
func (k KeysVariablesShorturlErrorsCreateShortUrl) GetPath() string {
	return "variables.shorturl.errors.create_short_url"
}
func (k KeysVariablesShorturlErrorsCreateShortUrl) GetPathSlice() []string {
	return []string{"variables", "shorturl", "errors", "create_short_url"}
}
func (k KeysVariablesShorturlErrorsCreateShortUrl) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesShorturlErrorsCreateShortUrl) SetVars(vars KeysVariablesShorturlErrorsCreateShortUrlVars) twiri18n.TranslationKey[KeysVariablesShorturlErrorsCreateShortUrlVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysVariablesShorturlErrorsUrlRequiredVars struct {
}
type KeysVariablesShorturlErrorsUrlRequired struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesShorturlErrorsUrlRequired) IsTranslationKey() {
}
func (k KeysVariablesShorturlErrorsUrlRequired) GetPath() string {
	return "variables.shorturl.errors.url_required"
}
func (k KeysVariablesShorturlErrorsUrlRequired) GetPathSlice() []string {
	return []string{"variables", "shorturl", "errors", "url_required"}
}
func (k KeysVariablesShorturlErrorsUrlRequired) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesShorturlErrorsUrlRequired) SetVars(vars KeysVariablesShorturlErrorsUrlRequiredVars) twiri18n.TranslationKey[KeysVariablesShorturlErrorsUrlRequiredVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesShorturlErrors struct {
	CreateShortUrl	KeysVariablesShorturlErrorsCreateShortUrl
	UrlRequired	KeysVariablesShorturlErrorsUrlRequired
}
type KeysVariablesShorturl struct {
	Errors KeysVariablesShorturlErrors
}
type KeysVariablesKeywordsErrorsIdNotProvidedVars struct {
}
type KeysVariablesKeywordsErrorsIdNotProvided struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesKeywordsErrorsIdNotProvided) IsTranslationKey() {
}
func (k KeysVariablesKeywordsErrorsIdNotProvided) GetPath() string {
	return "variables.keywords.errors.id_not_provided"
}
func (k KeysVariablesKeywordsErrorsIdNotProvided) GetPathSlice() []string {
	return []string{"variables", "keywords", "errors", "id_not_provided"}
}
func (k KeysVariablesKeywordsErrorsIdNotProvided) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesKeywordsErrorsIdNotProvided) SetVars(vars KeysVariablesKeywordsErrorsIdNotProvidedVars) twiri18n.TranslationKey[KeysVariablesKeywordsErrorsIdNotProvidedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesKeywordsErrorsNotFoundVars struct {
}
type KeysVariablesKeywordsErrorsNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesKeywordsErrorsNotFound) IsTranslationKey() {
}
func (k KeysVariablesKeywordsErrorsNotFound) GetPath() string {
	return "variables.keywords.errors.not_found"
}
func (k KeysVariablesKeywordsErrorsNotFound) GetPathSlice() []string {
	return []string{"variables", "keywords", "errors", "not_found"}
}
func (k KeysVariablesKeywordsErrorsNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesKeywordsErrorsNotFound) SetVars(vars KeysVariablesKeywordsErrorsNotFoundVars) twiri18n.TranslationKey[KeysVariablesKeywordsErrorsNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesKeywordsErrors struct {
	IdNotProvided	KeysVariablesKeywordsErrorsIdNotProvided
	NotFound	KeysVariablesKeywordsErrorsNotFound
}
type KeysVariablesKeywords struct {
	Errors KeysVariablesKeywordsErrors
}
type KeysVariablesSeventvErrorsNoActiveSetVars struct {
}
type KeysVariablesSeventvErrorsNoActiveSet struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSeventvErrorsNoActiveSet) IsTranslationKey() {
}
func (k KeysVariablesSeventvErrorsNoActiveSet) GetPath() string {
	return "variables.7tv.errors.no_active_set"
}
func (k KeysVariablesSeventvErrorsNoActiveSet) GetPathSlice() []string {
	return []string{"variables", "7tv", "errors", "no_active_set"}
}
func (k KeysVariablesSeventvErrorsNoActiveSet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSeventvErrorsNoActiveSet) SetVars(vars KeysVariablesSeventvErrorsNoActiveSetVars) twiri18n.TranslationKey[KeysVariablesSeventvErrorsNoActiveSetVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesSeventvErrorsProfileNotFoundVars struct {
	Reason any
}
type KeysVariablesSeventvErrorsProfileNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSeventvErrorsProfileNotFound) IsTranslationKey() {
}
func (k KeysVariablesSeventvErrorsProfileNotFound) GetPath() string {
	return "variables.7tv.errors.profile_not_found"
}
func (k KeysVariablesSeventvErrorsProfileNotFound) GetPathSlice() []string {
	return []string{"variables", "7tv", "errors", "profile_not_found"}
}
func (k KeysVariablesSeventvErrorsProfileNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSeventvErrorsProfileNotFound) SetVars(vars KeysVariablesSeventvErrorsProfileNotFoundVars) twiri18n.TranslationKey[KeysVariablesSeventvErrorsProfileNotFoundVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysVariablesSeventvErrorsEmoteNotFoundVars struct {
	Name any
}
type KeysVariablesSeventvErrorsEmoteNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSeventvErrorsEmoteNotFound) IsTranslationKey() {
}
func (k KeysVariablesSeventvErrorsEmoteNotFound) GetPath() string {
	return "variables.7tv.errors.emote_not_found"
}
func (k KeysVariablesSeventvErrorsEmoteNotFound) GetPathSlice() []string {
	return []string{"variables", "7tv", "errors", "emote_not_found"}
}
func (k KeysVariablesSeventvErrorsEmoteNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSeventvErrorsEmoteNotFound) SetVars(vars KeysVariablesSeventvErrorsEmoteNotFoundVars) twiri18n.TranslationKey[KeysVariablesSeventvErrorsEmoteNotFoundVars] {
	k.Vars = twiri18n.Vars{"name": vars.Name}
	return k
}

type KeysVariablesSeventvErrorsNoPaintVars struct {
}
type KeysVariablesSeventvErrorsNoPaint struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSeventvErrorsNoPaint) IsTranslationKey() {
}
func (k KeysVariablesSeventvErrorsNoPaint) GetPath() string {
	return "variables.7tv.errors.no_paint"
}
func (k KeysVariablesSeventvErrorsNoPaint) GetPathSlice() []string {
	return []string{"variables", "7tv", "errors", "no_paint"}
}
func (k KeysVariablesSeventvErrorsNoPaint) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSeventvErrorsNoPaint) SetVars(vars KeysVariablesSeventvErrorsNoPaintVars) twiri18n.TranslationKey[KeysVariablesSeventvErrorsNoPaintVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesSeventvErrorsNoRolesVars struct {
}
type KeysVariablesSeventvErrorsNoRoles struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSeventvErrorsNoRoles) IsTranslationKey() {
}
func (k KeysVariablesSeventvErrorsNoRoles) GetPath() string {
	return "variables.7tv.errors.no_roles"
}
func (k KeysVariablesSeventvErrorsNoRoles) GetPathSlice() []string {
	return []string{"variables", "7tv", "errors", "no_roles"}
}
func (k KeysVariablesSeventvErrorsNoRoles) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSeventvErrorsNoRoles) SetVars(vars KeysVariablesSeventvErrorsNoRolesVars) twiri18n.TranslationKey[KeysVariablesSeventvErrorsNoRolesVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesSeventvErrors struct {
	NoActiveSet	KeysVariablesSeventvErrorsNoActiveSet
	ProfileNotFound	KeysVariablesSeventvErrorsProfileNotFound
	EmoteNotFound	KeysVariablesSeventvErrorsEmoteNotFound
	NoPaint		KeysVariablesSeventvErrorsNoPaint
	NoRoles		KeysVariablesSeventvErrorsNoRoles
}
type KeysVariablesSeventv struct {
	Errors KeysVariablesSeventvErrors
}
type KeysVariablesStreamErrorsErrorVars struct {
}
type KeysVariablesStreamErrorsError struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesStreamErrorsError) IsTranslationKey() {
}
func (k KeysVariablesStreamErrorsError) GetPath() string {
	return "variables.stream.errors.error"
}
func (k KeysVariablesStreamErrorsError) GetPathSlice() []string {
	return []string{"variables", "stream", "errors", "error"}
}
func (k KeysVariablesStreamErrorsError) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesStreamErrorsError) SetVars(vars KeysVariablesStreamErrorsErrorVars) twiri18n.TranslationKey[KeysVariablesStreamErrorsErrorVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesStreamErrorsOfflineVars struct {
}
type KeysVariablesStreamErrorsOffline struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesStreamErrorsOffline) IsTranslationKey() {
}
func (k KeysVariablesStreamErrorsOffline) GetPath() string {
	return "variables.stream.errors.offline"
}
func (k KeysVariablesStreamErrorsOffline) GetPathSlice() []string {
	return []string{"variables", "stream", "errors", "offline"}
}
func (k KeysVariablesStreamErrorsOffline) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesStreamErrorsOffline) SetVars(vars KeysVariablesStreamErrorsOfflineVars) twiri18n.TranslationKey[KeysVariablesStreamErrorsOfflineVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesStreamErrorsCountFollowersVars struct {
}
type KeysVariablesStreamErrorsCountFollowers struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesStreamErrorsCountFollowers) IsTranslationKey() {
}
func (k KeysVariablesStreamErrorsCountFollowers) GetPath() string {
	return "variables.stream.errors.count_followers"
}
func (k KeysVariablesStreamErrorsCountFollowers) GetPathSlice() []string {
	return []string{"variables", "stream", "errors", "count_followers"}
}
func (k KeysVariablesStreamErrorsCountFollowers) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesStreamErrorsCountFollowers) SetVars(vars KeysVariablesStreamErrorsCountFollowersVars) twiri18n.TranslationKey[KeysVariablesStreamErrorsCountFollowersVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesStreamErrorsGetHistoryOfCategoriesVars struct {
}
type KeysVariablesStreamErrorsGetHistoryOfCategories struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesStreamErrorsGetHistoryOfCategories) IsTranslationKey() {
}
func (k KeysVariablesStreamErrorsGetHistoryOfCategories) GetPath() string {
	return "variables.stream.errors.get_history_of_categories"
}
func (k KeysVariablesStreamErrorsGetHistoryOfCategories) GetPathSlice() []string {
	return []string{"variables", "stream", "errors", "get_history_of_categories"}
}
func (k KeysVariablesStreamErrorsGetHistoryOfCategories) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesStreamErrorsGetHistoryOfCategories) SetVars(vars KeysVariablesStreamErrorsGetHistoryOfCategoriesVars) twiri18n.TranslationKey[KeysVariablesStreamErrorsGetHistoryOfCategoriesVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesStreamErrors struct {
	Error			KeysVariablesStreamErrorsError
	Offline			KeysVariablesStreamErrorsOffline
	CountFollowers		KeysVariablesStreamErrorsCountFollowers
	GetHistoryOfCategories	KeysVariablesStreamErrorsGetHistoryOfCategories
}
type KeysVariablesStreamInfoOfflineVars struct {
}
type KeysVariablesStreamInfoOffline struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesStreamInfoOffline) IsTranslationKey() {
}
func (k KeysVariablesStreamInfoOffline) GetPath() string {
	return "variables.stream.info.offline"
}
func (k KeysVariablesStreamInfoOffline) GetPathSlice() []string {
	return []string{"variables", "stream", "info", "offline"}
}
func (k KeysVariablesStreamInfoOffline) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesStreamInfoOffline) SetVars(vars KeysVariablesStreamInfoOfflineVars) twiri18n.TranslationKey[KeysVariablesStreamInfoOfflineVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesStreamInfoNoHistoryVars struct {
}
type KeysVariablesStreamInfoNoHistory struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesStreamInfoNoHistory) IsTranslationKey() {
}
func (k KeysVariablesStreamInfoNoHistory) GetPath() string {
	return "variables.stream.info.no_history"
}
func (k KeysVariablesStreamInfoNoHistory) GetPathSlice() []string {
	return []string{"variables", "stream", "info", "no_history"}
}
func (k KeysVariablesStreamInfoNoHistory) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesStreamInfoNoHistory) SetVars(vars KeysVariablesStreamInfoNoHistoryVars) twiri18n.TranslationKey[KeysVariablesStreamInfoNoHistoryVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesStreamInfo struct {
	Offline		KeysVariablesStreamInfoOffline
	NoHistory	KeysVariablesStreamInfoNoHistory
}
type KeysVariablesStream struct {
	Errors	KeysVariablesStreamErrors
	Info	KeysVariablesStreamInfo
}
type KeysVariablesCountdownErrorsNotPassedParamsVars struct {
}
type KeysVariablesCountdownErrorsNotPassedParams struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesCountdownErrorsNotPassedParams) IsTranslationKey() {
}
func (k KeysVariablesCountdownErrorsNotPassedParams) GetPath() string {
	return "variables.countdown.errors.not_passed_params"
}
func (k KeysVariablesCountdownErrorsNotPassedParams) GetPathSlice() []string {
	return []string{"variables", "countdown", "errors", "not_passed_params"}
}
func (k KeysVariablesCountdownErrorsNotPassedParams) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesCountdownErrorsNotPassedParams) SetVars(vars KeysVariablesCountdownErrorsNotPassedParamsVars) twiri18n.TranslationKey[KeysVariablesCountdownErrorsNotPassedParamsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesCountdownErrorsParseDateVars struct {
}
type KeysVariablesCountdownErrorsParseDate struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesCountdownErrorsParseDate) IsTranslationKey() {
}
func (k KeysVariablesCountdownErrorsParseDate) GetPath() string {
	return "variables.countdown.errors.parse_date"
}
func (k KeysVariablesCountdownErrorsParseDate) GetPathSlice() []string {
	return []string{"variables", "countdown", "errors", "parse_date"}
}
func (k KeysVariablesCountdownErrorsParseDate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesCountdownErrorsParseDate) SetVars(vars KeysVariablesCountdownErrorsParseDateVars) twiri18n.TranslationKey[KeysVariablesCountdownErrorsParseDateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesCountdownErrors struct {
	NotPassedParams	KeysVariablesCountdownErrorsNotPassedParams
	ParseDate	KeysVariablesCountdownErrorsParseDate
}
type KeysVariablesCountdown struct {
	Errors KeysVariablesCountdownErrors
}
type KeysVariablesSongInfoFailedGetSpotifyIntegrationVars struct {
}
type KeysVariablesSongInfoFailedGetSpotifyIntegration struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongInfoFailedGetSpotifyIntegration) IsTranslationKey() {
}
func (k KeysVariablesSongInfoFailedGetSpotifyIntegration) GetPath() string {
	return "variables.song.info.failed_get_spotify_integration"
}
func (k KeysVariablesSongInfoFailedGetSpotifyIntegration) GetPathSlice() []string {
	return []string{"variables", "song", "info", "failed_get_spotify_integration"}
}
func (k KeysVariablesSongInfoFailedGetSpotifyIntegration) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongInfoFailedGetSpotifyIntegration) SetVars(vars KeysVariablesSongInfoFailedGetSpotifyIntegrationVars) twiri18n.TranslationKey[KeysVariablesSongInfoFailedGetSpotifyIntegrationVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesSongInfoGetSpotifyIntegrationVars struct {
	Reason any
}
type KeysVariablesSongInfoGetSpotifyIntegration struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongInfoGetSpotifyIntegration) IsTranslationKey() {
}
func (k KeysVariablesSongInfoGetSpotifyIntegration) GetPath() string {
	return "variables.song.info.get_spotify_integration"
}
func (k KeysVariablesSongInfoGetSpotifyIntegration) GetPathSlice() []string {
	return []string{"variables", "song", "info", "get_spotify_integration"}
}
func (k KeysVariablesSongInfoGetSpotifyIntegration) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongInfoGetSpotifyIntegration) SetVars(vars KeysVariablesSongInfoGetSpotifyIntegrationVars) twiri18n.TranslationKey[KeysVariablesSongInfoGetSpotifyIntegrationVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysVariablesSongInfoHistoryVars struct {
	TrackTitle	any
	TrackArtist	any
	Minutes		any
}
type KeysVariablesSongInfoHistory struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongInfoHistory) IsTranslationKey() {
}
func (k KeysVariablesSongInfoHistory) GetPath() string {
	return "variables.song.info.history"
}
func (k KeysVariablesSongInfoHistory) GetPathSlice() []string {
	return []string{"variables", "song", "info", "history"}
}
func (k KeysVariablesSongInfoHistory) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongInfoHistory) SetVars(vars KeysVariablesSongInfoHistoryVars) twiri18n.TranslationKey[KeysVariablesSongInfoHistoryVars] {
	k.Vars = twiri18n.Vars{"trackTitle": vars.TrackTitle, "trackArtist": vars.TrackArtist, "minutes": vars.Minutes}
	return k
}

type KeysVariablesSongInfoNoNeededScopeVars struct {
}
type KeysVariablesSongInfoNoNeededScope struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongInfoNoNeededScope) IsTranslationKey() {
}
func (k KeysVariablesSongInfoNoNeededScope) GetPath() string {
	return "variables.song.info.no_needed_scope"
}
func (k KeysVariablesSongInfoNoNeededScope) GetPathSlice() []string {
	return []string{"variables", "song", "info", "no_needed_scope"}
}
func (k KeysVariablesSongInfoNoNeededScope) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongInfoNoNeededScope) SetVars(vars KeysVariablesSongInfoNoNeededScopeVars) twiri18n.TranslationKey[KeysVariablesSongInfoNoNeededScopeVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesSongInfoGetSpotifyEntityVars struct {
}
type KeysVariablesSongInfoGetSpotifyEntity struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongInfoGetSpotifyEntity) IsTranslationKey() {
}
func (k KeysVariablesSongInfoGetSpotifyEntity) GetPath() string {
	return "variables.song.info.get_spotify_entity"
}
func (k KeysVariablesSongInfoGetSpotifyEntity) GetPathSlice() []string {
	return []string{"variables", "song", "info", "get_spotify_entity"}
}
func (k KeysVariablesSongInfoGetSpotifyEntity) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongInfoGetSpotifyEntity) SetVars(vars KeysVariablesSongInfoGetSpotifyEntityVars) twiri18n.TranslationKey[KeysVariablesSongInfoGetSpotifyEntityVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesSongInfoNoIntegrationsVars struct {
}
type KeysVariablesSongInfoNoIntegrations struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongInfoNoIntegrations) IsTranslationKey() {
}
func (k KeysVariablesSongInfoNoIntegrations) GetPath() string {
	return "variables.song.info.no_integrations"
}
func (k KeysVariablesSongInfoNoIntegrations) GetPathSlice() []string {
	return []string{"variables", "song", "info", "no_integrations"}
}
func (k KeysVariablesSongInfoNoIntegrations) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongInfoNoIntegrations) SetVars(vars KeysVariablesSongInfoNoIntegrationsVars) twiri18n.TranslationKey[KeysVariablesSongInfoNoIntegrationsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesSongInfoLastfmIntegrationVars struct {
}
type KeysVariablesSongInfoLastfmIntegration struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongInfoLastfmIntegration) IsTranslationKey() {
}
func (k KeysVariablesSongInfoLastfmIntegration) GetPath() string {
	return "variables.song.info.lastfm_integration"
}
func (k KeysVariablesSongInfoLastfmIntegration) GetPathSlice() []string {
	return []string{"variables", "song", "info", "lastfm_integration"}
}
func (k KeysVariablesSongInfoLastfmIntegration) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongInfoLastfmIntegration) SetVars(vars KeysVariablesSongInfoLastfmIntegrationVars) twiri18n.TranslationKey[KeysVariablesSongInfoLastfmIntegrationVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesSongInfoSpotifyNotConnectedVars struct {
}
type KeysVariablesSongInfoSpotifyNotConnected struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongInfoSpotifyNotConnected) IsTranslationKey() {
}
func (k KeysVariablesSongInfoSpotifyNotConnected) GetPath() string {
	return "variables.song.info.spotify_not_connected"
}
func (k KeysVariablesSongInfoSpotifyNotConnected) GetPathSlice() []string {
	return []string{"variables", "song", "info", "spotify_not_connected"}
}
func (k KeysVariablesSongInfoSpotifyNotConnected) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongInfoSpotifyNotConnected) SetVars(vars KeysVariablesSongInfoSpotifyNotConnectedVars) twiri18n.TranslationKey[KeysVariablesSongInfoSpotifyNotConnectedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesSongInfo struct {
	FailedGetSpotifyIntegration	KeysVariablesSongInfoFailedGetSpotifyIntegration
	GetSpotifyIntegration		KeysVariablesSongInfoGetSpotifyIntegration
	History				KeysVariablesSongInfoHistory
	NoNeededScope			KeysVariablesSongInfoNoNeededScope
	GetSpotifyEntity		KeysVariablesSongInfoGetSpotifyEntity
	NoIntegrations			KeysVariablesSongInfoNoIntegrations
	LastfmIntegration		KeysVariablesSongInfoLastfmIntegration
	SpotifyNotConnected		KeysVariablesSongInfoSpotifyNotConnected
}
type KeysVariablesSongErrorsGetRecentTracksVars struct {
	Reason any
}
type KeysVariablesSongErrorsGetRecentTracks struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongErrorsGetRecentTracks) IsTranslationKey() {
}
func (k KeysVariablesSongErrorsGetRecentTracks) GetPath() string {
	return "variables.song.errors.get_recent_tracks"
}
func (k KeysVariablesSongErrorsGetRecentTracks) GetPathSlice() []string {
	return []string{"variables", "song", "errors", "get_recent_tracks"}
}
func (k KeysVariablesSongErrorsGetRecentTracks) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongErrorsGetRecentTracks) SetVars(vars KeysVariablesSongErrorsGetRecentTracksVars) twiri18n.TranslationKey[KeysVariablesSongErrorsGetRecentTracksVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysVariablesSongErrorsCreateLastfmServiceVars struct {
	Reason any
}
type KeysVariablesSongErrorsCreateLastfmService struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongErrorsCreateLastfmService) IsTranslationKey() {
}
func (k KeysVariablesSongErrorsCreateLastfmService) GetPath() string {
	return "variables.song.errors.create_lastfm_service"
}
func (k KeysVariablesSongErrorsCreateLastfmService) GetPathSlice() []string {
	return []string{"variables", "song", "errors", "create_lastfm_service"}
}
func (k KeysVariablesSongErrorsCreateLastfmService) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongErrorsCreateLastfmService) SetVars(vars KeysVariablesSongErrorsCreateLastfmServiceVars) twiri18n.TranslationKey[KeysVariablesSongErrorsCreateLastfmServiceVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysVariablesSongErrorsFetchTracksLastfmVars struct {
	Reason any
}
type KeysVariablesSongErrorsFetchTracksLastfm struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongErrorsFetchTracksLastfm) IsTranslationKey() {
}
func (k KeysVariablesSongErrorsFetchTracksLastfm) GetPath() string {
	return "variables.song.errors.fetch_tracks_lastfm"
}
func (k KeysVariablesSongErrorsFetchTracksLastfm) GetPathSlice() []string {
	return []string{"variables", "song", "errors", "fetch_tracks_lastfm"}
}
func (k KeysVariablesSongErrorsFetchTracksLastfm) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongErrorsFetchTracksLastfm) SetVars(vars KeysVariablesSongErrorsFetchTracksLastfmVars) twiri18n.TranslationKey[KeysVariablesSongErrorsFetchTracksLastfmVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysVariablesSongErrorsFetchTracksSpotifyVars struct {
	Reason any
}
type KeysVariablesSongErrorsFetchTracksSpotify struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongErrorsFetchTracksSpotify) IsTranslationKey() {
}
func (k KeysVariablesSongErrorsFetchTracksSpotify) GetPath() string {
	return "variables.song.errors.fetch_tracks_spotify"
}
func (k KeysVariablesSongErrorsFetchTracksSpotify) GetPathSlice() []string {
	return []string{"variables", "song", "errors", "fetch_tracks_spotify"}
}
func (k KeysVariablesSongErrorsFetchTracksSpotify) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongErrorsFetchTracksSpotify) SetVars(vars KeysVariablesSongErrorsFetchTracksSpotifyVars) twiri18n.TranslationKey[KeysVariablesSongErrorsFetchTracksSpotifyVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysVariablesSongErrorsParsePlayedAtVars struct {
}
type KeysVariablesSongErrorsParsePlayedAt struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSongErrorsParsePlayedAt) IsTranslationKey() {
}
func (k KeysVariablesSongErrorsParsePlayedAt) GetPath() string {
	return "variables.song.errors.parse_played_at"
}
func (k KeysVariablesSongErrorsParsePlayedAt) GetPathSlice() []string {
	return []string{"variables", "song", "errors", "parse_played_at"}
}
func (k KeysVariablesSongErrorsParsePlayedAt) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSongErrorsParsePlayedAt) SetVars(vars KeysVariablesSongErrorsParsePlayedAtVars) twiri18n.TranslationKey[KeysVariablesSongErrorsParsePlayedAtVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesSongErrors struct {
	GetRecentTracks		KeysVariablesSongErrorsGetRecentTracks
	CreateLastfmService	KeysVariablesSongErrorsCreateLastfmService
	FetchTracksLastfm	KeysVariablesSongErrorsFetchTracksLastfm
	FetchTracksSpotify	KeysVariablesSongErrorsFetchTracksSpotify
	ParsePlayedAt		KeysVariablesSongErrorsParsePlayedAt
}
type KeysVariablesSong struct {
	Info	KeysVariablesSongInfo
	Errors	KeysVariablesSongErrors
}
type KeysVariablesChatEvalInfoWrongCodeVars struct {
}
type KeysVariablesChatEvalInfoWrongCode struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesChatEvalInfoWrongCode) IsTranslationKey() {
}
func (k KeysVariablesChatEvalInfoWrongCode) GetPath() string {
	return "variables.chat_eval.info.wrong_code"
}
func (k KeysVariablesChatEvalInfoWrongCode) GetPathSlice() []string {
	return []string{"variables", "chat_eval", "info", "wrong_code"}
}
func (k KeysVariablesChatEvalInfoWrongCode) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesChatEvalInfoWrongCode) SetVars(vars KeysVariablesChatEvalInfoWrongCodeVars) twiri18n.TranslationKey[KeysVariablesChatEvalInfoWrongCodeVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesChatEvalInfo struct {
	WrongCode KeysVariablesChatEvalInfoWrongCode
}
type KeysVariablesChatEval struct {
	Info KeysVariablesChatEvalInfo
}
type KeysVariablesRequestErrorsExecuteRequestVars struct {
}
type KeysVariablesRequestErrorsExecuteRequest struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesRequestErrorsExecuteRequest) IsTranslationKey() {
}
func (k KeysVariablesRequestErrorsExecuteRequest) GetPath() string {
	return "variables.request.errors.execute_request"
}
func (k KeysVariablesRequestErrorsExecuteRequest) GetPathSlice() []string {
	return []string{"variables", "request", "errors", "execute_request"}
}
func (k KeysVariablesRequestErrorsExecuteRequest) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesRequestErrorsExecuteRequest) SetVars(vars KeysVariablesRequestErrorsExecuteRequestVars) twiri18n.TranslationKey[KeysVariablesRequestErrorsExecuteRequestVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesRequestErrors struct {
	ExecuteRequest KeysVariablesRequestErrorsExecuteRequest
}
type KeysVariablesRequest struct {
	Errors KeysVariablesRequestErrors
}
type KeysVariablesValorantInfoMatchesVars struct {
	MatchResult	any
	RoundsWon	any
	RoundsLost	any
	Char		any
	KDA		any
}
type KeysVariablesValorantInfoMatches struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesValorantInfoMatches) IsTranslationKey() {
}
func (k KeysVariablesValorantInfoMatches) GetPath() string {
	return "variables.valorant.info.matches"
}
func (k KeysVariablesValorantInfoMatches) GetPathSlice() []string {
	return []string{"variables", "valorant", "info", "matches"}
}
func (k KeysVariablesValorantInfoMatches) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesValorantInfoMatches) SetVars(vars KeysVariablesValorantInfoMatchesVars) twiri18n.TranslationKey[KeysVariablesValorantInfoMatchesVars] {
	k.Vars = twiri18n.Vars{"matchResult": vars.MatchResult, "roundsWon": vars.RoundsWon, "roundsLost": vars.RoundsLost, "char": vars.Char, "KDA": vars.KDA}
	return k
}

type KeysVariablesValorantInfo struct {
	Matches KeysVariablesValorantInfoMatches
}
type KeysVariablesValorant struct {
	Info KeysVariablesValorantInfo
}
type KeysVariablesCustomVarErrorsUpdateCustomVarVars struct {
}
type KeysVariablesCustomVarErrorsUpdateCustomVar struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesCustomVarErrorsUpdateCustomVar) IsTranslationKey() {
}
func (k KeysVariablesCustomVarErrorsUpdateCustomVar) GetPath() string {
	return "variables.custom_var.errors.update_custom_var"
}
func (k KeysVariablesCustomVarErrorsUpdateCustomVar) GetPathSlice() []string {
	return []string{"variables", "custom_var", "errors", "update_custom_var"}
}
func (k KeysVariablesCustomVarErrorsUpdateCustomVar) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesCustomVarErrorsUpdateCustomVar) SetVars(vars KeysVariablesCustomVarErrorsUpdateCustomVarVars) twiri18n.TranslationKey[KeysVariablesCustomVarErrorsUpdateCustomVarVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesCustomVarErrorsEvaluateVariableVars struct {
}
type KeysVariablesCustomVarErrorsEvaluateVariable struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesCustomVarErrorsEvaluateVariable) IsTranslationKey() {
}
func (k KeysVariablesCustomVarErrorsEvaluateVariable) GetPath() string {
	return "variables.custom_var.errors.evaluate_variable"
}
func (k KeysVariablesCustomVarErrorsEvaluateVariable) GetPathSlice() []string {
	return []string{"variables", "custom_var", "errors", "evaluate_variable"}
}
func (k KeysVariablesCustomVarErrorsEvaluateVariable) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesCustomVarErrorsEvaluateVariable) SetVars(vars KeysVariablesCustomVarErrorsEvaluateVariableVars) twiri18n.TranslationKey[KeysVariablesCustomVarErrorsEvaluateVariableVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesCustomVarErrorsWrongNumbersVars struct {
	Reason any
}
type KeysVariablesCustomVarErrorsWrongNumbers struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesCustomVarErrorsWrongNumbers) IsTranslationKey() {
}
func (k KeysVariablesCustomVarErrorsWrongNumbers) GetPath() string {
	return "variables.custom_var.errors.wrong_numbers"
}
func (k KeysVariablesCustomVarErrorsWrongNumbers) GetPathSlice() []string {
	return []string{"variables", "custom_var", "errors", "wrong_numbers"}
}
func (k KeysVariablesCustomVarErrorsWrongNumbers) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesCustomVarErrorsWrongNumbers) SetVars(vars KeysVariablesCustomVarErrorsWrongNumbersVars) twiri18n.TranslationKey[KeysVariablesCustomVarErrorsWrongNumbersVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysVariablesCustomVarErrors struct {
	UpdateCustomVar		KeysVariablesCustomVarErrorsUpdateCustomVar
	EvaluateVariable	KeysVariablesCustomVarErrorsEvaluateVariable
	WrongNumbers		KeysVariablesCustomVarErrorsWrongNumbers
}
type KeysVariablesCustomVar struct {
	Errors KeysVariablesCustomVarErrors
}
type KeysVariablesFollowersErrorsGetFollowersVars struct {
}
type KeysVariablesFollowersErrorsGetFollowers struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesFollowersErrorsGetFollowers) IsTranslationKey() {
}
func (k KeysVariablesFollowersErrorsGetFollowers) GetPath() string {
	return "variables.followers.errors.get_followers"
}
func (k KeysVariablesFollowersErrorsGetFollowers) GetPathSlice() []string {
	return []string{"variables", "followers", "errors", "get_followers"}
}
func (k KeysVariablesFollowersErrorsGetFollowers) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesFollowersErrorsGetFollowers) SetVars(vars KeysVariablesFollowersErrorsGetFollowersVars) twiri18n.TranslationKey[KeysVariablesFollowersErrorsGetFollowersVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesFollowersErrors struct {
	GetFollowers KeysVariablesFollowersErrorsGetFollowers
}
type KeysVariablesFollowers struct {
	Errors KeysVariablesFollowersErrors
}
type KeysVariablesSubscribersErrorsGetSubscribersVars struct {
}
type KeysVariablesSubscribersErrorsGetSubscribers struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesSubscribersErrorsGetSubscribers) IsTranslationKey() {
}
func (k KeysVariablesSubscribersErrorsGetSubscribers) GetPath() string {
	return "variables.subscribers.errors.get_subscribers"
}
func (k KeysVariablesSubscribersErrorsGetSubscribers) GetPathSlice() []string {
	return []string{"variables", "subscribers", "errors", "get_subscribers"}
}
func (k KeysVariablesSubscribersErrorsGetSubscribers) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesSubscribersErrorsGetSubscribers) SetVars(vars KeysVariablesSubscribersErrorsGetSubscribersVars) twiri18n.TranslationKey[KeysVariablesSubscribersErrorsGetSubscribersVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesSubscribersErrors struct {
	GetSubscribers KeysVariablesSubscribersErrorsGetSubscribers
}
type KeysVariablesSubscribers struct {
	Errors KeysVariablesSubscribersErrors
}
type KeysVariablesUserErrorsFindUserOnTwitchVars struct {
}
type KeysVariablesUserErrorsFindUserOnTwitch struct {
	Vars twiri18n.Vars
}

func (k KeysVariablesUserErrorsFindUserOnTwitch) IsTranslationKey() {
}
func (k KeysVariablesUserErrorsFindUserOnTwitch) GetPath() string {
	return "variables.user.errors.find_user_on_twitch"
}
func (k KeysVariablesUserErrorsFindUserOnTwitch) GetPathSlice() []string {
	return []string{"variables", "user", "errors", "find_user_on_twitch"}
}
func (k KeysVariablesUserErrorsFindUserOnTwitch) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysVariablesUserErrorsFindUserOnTwitch) SetVars(vars KeysVariablesUserErrorsFindUserOnTwitchVars) twiri18n.TranslationKey[KeysVariablesUserErrorsFindUserOnTwitchVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysVariablesUserErrors struct {
	FindUserOnTwitch KeysVariablesUserErrorsFindUserOnTwitch
}
type KeysVariablesUser struct {
	Errors KeysVariablesUserErrors
}
type KeysVariables struct {
	Commands	KeysVariablesCommands
	Random		KeysVariablesRandom
	Shorturl	KeysVariablesShorturl
	Keywords	KeysVariablesKeywords
	Seventv		KeysVariablesSeventv
	Stream		KeysVariablesStream
	Countdown	KeysVariablesCountdown
	Song		KeysVariablesSong
	ChatEval	KeysVariablesChatEval
	Request		KeysVariablesRequest
	Valorant	KeysVariablesValorant
	CustomVar	KeysVariablesCustomVar
	Followers	KeysVariablesFollowers
	Subscribers	KeysVariablesSubscribers
	User		KeysVariablesUser
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
type KeysCommandsShorturl struct {
	Errors	KeysCommandsShorturlErrors
	Success	KeysCommandsShorturlSuccess
}
type KeysCommandsMarkerErrorsCannotCreateMarkerVars struct {
	Reason any
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
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsMarkerErrors struct {
	CannotCreateMarker KeysCommandsMarkerErrorsCannotCreateMarker
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
type KeysCommandsMarker struct {
	Errors	KeysCommandsMarkerErrors
	Success	KeysCommandsMarkerSuccess
}
type KeysCommandsNukeHintsNukeTimeArgNameVars struct {
}
type KeysCommandsNukeHintsNukeTimeArgName struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsNukeHintsNukeTimeArgName) IsTranslationKey() {
}
func (k KeysCommandsNukeHintsNukeTimeArgName) GetPath() string {
	return "commands.nuke.hints.nukeTimeArgName"
}
func (k KeysCommandsNukeHintsNukeTimeArgName) GetPathSlice() []string {
	return []string{"commands", "nuke", "hints", "nukeTimeArgName"}
}
func (k KeysCommandsNukeHintsNukeTimeArgName) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsNukeHintsNukeTimeArgName) SetVars(vars KeysCommandsNukeHintsNukeTimeArgNameVars) twiri18n.TranslationKey[KeysCommandsNukeHintsNukeTimeArgNameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsNukeHints struct {
	NukeTimeArgName KeysCommandsNukeHintsNukeTimeArgName
}
type KeysCommandsNukeErrorsCannotGetUsersStatsVars struct {
}
type KeysCommandsNukeErrorsCannotGetUsersStats struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsNukeErrorsCannotGetUsersStats) IsTranslationKey() {
}
func (k KeysCommandsNukeErrorsCannotGetUsersStats) GetPath() string {
	return "commands.nuke.errors.cannot_get_users_stats"
}
func (k KeysCommandsNukeErrorsCannotGetUsersStats) GetPathSlice() []string {
	return []string{"commands", "nuke", "errors", "cannot_get_users_stats"}
}
func (k KeysCommandsNukeErrorsCannotGetUsersStats) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsNukeErrorsCannotGetUsersStats) SetVars(vars KeysCommandsNukeErrorsCannotGetUsersStatsVars) twiri18n.TranslationKey[KeysCommandsNukeErrorsCannotGetUsersStatsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsNukeErrorsCannotGetHandeledMessagesVars struct {
}
type KeysCommandsNukeErrorsCannotGetHandeledMessages struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsNukeErrorsCannotGetHandeledMessages) IsTranslationKey() {
}
func (k KeysCommandsNukeErrorsCannotGetHandeledMessages) GetPath() string {
	return "commands.nuke.errors.cannot_get_handeled_messages"
}
func (k KeysCommandsNukeErrorsCannotGetHandeledMessages) GetPathSlice() []string {
	return []string{"commands", "nuke", "errors", "cannot_get_handeled_messages"}
}
func (k KeysCommandsNukeErrorsCannotGetHandeledMessages) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsNukeErrorsCannotGetHandeledMessages) SetVars(vars KeysCommandsNukeErrorsCannotGetHandeledMessagesVars) twiri18n.TranslationKey[KeysCommandsNukeErrorsCannotGetHandeledMessagesVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsNukeErrorsCannotDeleteMessagesVars struct {
}
type KeysCommandsNukeErrorsCannotDeleteMessages struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsNukeErrorsCannotDeleteMessages) IsTranslationKey() {
}
func (k KeysCommandsNukeErrorsCannotDeleteMessages) GetPath() string {
	return "commands.nuke.errors.cannot_delete_messages"
}
func (k KeysCommandsNukeErrorsCannotDeleteMessages) GetPathSlice() []string {
	return []string{"commands", "nuke", "errors", "cannot_delete_messages"}
}
func (k KeysCommandsNukeErrorsCannotDeleteMessages) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsNukeErrorsCannotDeleteMessages) SetVars(vars KeysCommandsNukeErrorsCannotDeleteMessagesVars) twiri18n.TranslationKey[KeysCommandsNukeErrorsCannotDeleteMessagesVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsNukeErrorsTimeoutDurationVars struct {
}
type KeysCommandsNukeErrorsTimeoutDuration struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsNukeErrorsTimeoutDuration) IsTranslationKey() {
}
func (k KeysCommandsNukeErrorsTimeoutDuration) GetPath() string {
	return "commands.nuke.errors.timeout_duration"
}
func (k KeysCommandsNukeErrorsTimeoutDuration) GetPathSlice() []string {
	return []string{"commands", "nuke", "errors", "timeout_duration"}
}
func (k KeysCommandsNukeErrorsTimeoutDuration) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsNukeErrorsTimeoutDuration) SetVars(vars KeysCommandsNukeErrorsTimeoutDurationVars) twiri18n.TranslationKey[KeysCommandsNukeErrorsTimeoutDurationVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsNukeErrorsParseDurationVars struct {
}
type KeysCommandsNukeErrorsParseDuration struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsNukeErrorsParseDuration) IsTranslationKey() {
}
func (k KeysCommandsNukeErrorsParseDuration) GetPath() string {
	return "commands.nuke.errors.parse_duration"
}
func (k KeysCommandsNukeErrorsParseDuration) GetPathSlice() []string {
	return []string{"commands", "nuke", "errors", "parse_duration"}
}
func (k KeysCommandsNukeErrorsParseDuration) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsNukeErrorsParseDuration) SetVars(vars KeysCommandsNukeErrorsParseDurationVars) twiri18n.TranslationKey[KeysCommandsNukeErrorsParseDurationVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsNukeErrorsInvalidDurationVars struct {
}
type KeysCommandsNukeErrorsInvalidDuration struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsNukeErrorsInvalidDuration) IsTranslationKey() {
}
func (k KeysCommandsNukeErrorsInvalidDuration) GetPath() string {
	return "commands.nuke.errors.invalid_duration"
}
func (k KeysCommandsNukeErrorsInvalidDuration) GetPathSlice() []string {
	return []string{"commands", "nuke", "errors", "invalid_duration"}
}
func (k KeysCommandsNukeErrorsInvalidDuration) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsNukeErrorsInvalidDuration) SetVars(vars KeysCommandsNukeErrorsInvalidDurationVars) twiri18n.TranslationKey[KeysCommandsNukeErrorsInvalidDurationVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsNukeErrors struct {
	CannotGetUsersStats		KeysCommandsNukeErrorsCannotGetUsersStats
	CannotGetHandeledMessages	KeysCommandsNukeErrorsCannotGetHandeledMessages
	CannotDeleteMessages		KeysCommandsNukeErrorsCannotDeleteMessages
	TimeoutDuration			KeysCommandsNukeErrorsTimeoutDuration
	ParseDuration			KeysCommandsNukeErrorsParseDuration
	InvalidDuration			KeysCommandsNukeErrorsInvalidDuration
}
type KeysCommandsNuke struct {
	Hints	KeysCommandsNukeHints
	Errors	KeysCommandsNukeErrors
}
type KeysCommandsStatsInfoWatchingStreamVars struct {
	UserWatching any
}
type KeysCommandsStatsInfoWatchingStream struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsInfoWatchingStream) IsTranslationKey() {
}
func (k KeysCommandsStatsInfoWatchingStream) GetPath() string {
	return "commands.stats.info.watching_stream"
}
func (k KeysCommandsStatsInfoWatchingStream) GetPathSlice() []string {
	return []string{"commands", "stats", "info", "watching_stream"}
}
func (k KeysCommandsStatsInfoWatchingStream) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsInfoWatchingStream) SetVars(vars KeysCommandsStatsInfoWatchingStreamVars) twiri18n.TranslationKey[KeysCommandsStatsInfoWatchingStreamVars] {
	k.Vars = twiri18n.Vars{"userWatching": vars.UserWatching}
	return k
}

type KeysCommandsStatsInfo struct {
	WatchingStream KeysCommandsStatsInfoWatchingStream
}
type KeysCommandsStatsMeMessagesVars struct {
}
type KeysCommandsStatsMeMessages struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsMeMessages) IsTranslationKey() {
}
func (k KeysCommandsStatsMeMessages) GetPath() string {
	return "commands.stats.me.messages"
}
func (k KeysCommandsStatsMeMessages) GetPathSlice() []string {
	return []string{"commands", "stats", "me", "messages"}
}
func (k KeysCommandsStatsMeMessages) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsMeMessages) SetVars(vars KeysCommandsStatsMeMessagesVars) twiri18n.TranslationKey[KeysCommandsStatsMeMessagesVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsMeEmotesVars struct {
}
type KeysCommandsStatsMeEmotes struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsMeEmotes) IsTranslationKey() {
}
func (k KeysCommandsStatsMeEmotes) GetPath() string {
	return "commands.stats.me.emotes"
}
func (k KeysCommandsStatsMeEmotes) GetPathSlice() []string {
	return []string{"commands", "stats", "me", "emotes"}
}
func (k KeysCommandsStatsMeEmotes) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsMeEmotes) SetVars(vars KeysCommandsStatsMeEmotesVars) twiri18n.TranslationKey[KeysCommandsStatsMeEmotesVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsMePointsVars struct {
}
type KeysCommandsStatsMePoints struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsMePoints) IsTranslationKey() {
}
func (k KeysCommandsStatsMePoints) GetPath() string {
	return "commands.stats.me.points"
}
func (k KeysCommandsStatsMePoints) GetPathSlice() []string {
	return []string{"commands", "stats", "me", "points"}
}
func (k KeysCommandsStatsMePoints) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsMePoints) SetVars(vars KeysCommandsStatsMePointsVars) twiri18n.TranslationKey[KeysCommandsStatsMePointsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsMeSongsVars struct {
}
type KeysCommandsStatsMeSongs struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsMeSongs) IsTranslationKey() {
}
func (k KeysCommandsStatsMeSongs) GetPath() string {
	return "commands.stats.me.songs"
}
func (k KeysCommandsStatsMeSongs) GetPathSlice() []string {
	return []string{"commands", "stats", "me", "songs"}
}
func (k KeysCommandsStatsMeSongs) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsMeSongs) SetVars(vars KeysCommandsStatsMeSongsVars) twiri18n.TranslationKey[KeysCommandsStatsMeSongsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsMeWatchedVars struct {
}
type KeysCommandsStatsMeWatched struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsStatsMeWatched) IsTranslationKey() {
}
func (k KeysCommandsStatsMeWatched) GetPath() string {
	return "commands.stats.me.watched"
}
func (k KeysCommandsStatsMeWatched) GetPathSlice() []string {
	return []string{"commands", "stats", "me", "watched"}
}
func (k KeysCommandsStatsMeWatched) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsStatsMeWatched) SetVars(vars KeysCommandsStatsMeWatchedVars) twiri18n.TranslationKey[KeysCommandsStatsMeWatchedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsStatsMe struct {
	Messages	KeysCommandsStatsMeMessages
	Emotes		KeysCommandsStatsMeEmotes
	Points		KeysCommandsStatsMePoints
	Songs		KeysCommandsStatsMeSongs
	Watched		KeysCommandsStatsMeWatched
}
type KeysCommandsStats struct {
	Info	KeysCommandsStatsInfo
	Me	KeysCommandsStatsMe
}
type KeysCommandsVipsErrorsInvalidDurationVars struct {
}
type KeysCommandsVipsErrorsInvalidDuration struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsErrorsInvalidDuration) IsTranslationKey() {
}
func (k KeysCommandsVipsErrorsInvalidDuration) GetPath() string {
	return "commands.vips.errors.invalid_duration"
}
func (k KeysCommandsVipsErrorsInvalidDuration) GetPathSlice() []string {
	return []string{"commands", "vips", "errors", "invalid_duration"}
}
func (k KeysCommandsVipsErrorsInvalidDuration) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsErrorsInvalidDuration) SetVars(vars KeysCommandsVipsErrorsInvalidDurationVars) twiri18n.TranslationKey[KeysCommandsVipsErrorsInvalidDurationVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsErrorsAlreadyHaveRoleVars struct {
}
type KeysCommandsVipsErrorsAlreadyHaveRole struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsErrorsAlreadyHaveRole) IsTranslationKey() {
}
func (k KeysCommandsVipsErrorsAlreadyHaveRole) GetPath() string {
	return "commands.vips.errors.already_have_role"
}
func (k KeysCommandsVipsErrorsAlreadyHaveRole) GetPathSlice() []string {
	return []string{"commands", "vips", "errors", "already_have_role"}
}
func (k KeysCommandsVipsErrorsAlreadyHaveRole) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsErrorsAlreadyHaveRole) SetVars(vars KeysCommandsVipsErrorsAlreadyHaveRoleVars) twiri18n.TranslationKey[KeysCommandsVipsErrorsAlreadyHaveRoleVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsErrorsAddedVars struct {
	UserName any
}
type KeysCommandsVipsErrorsAdded struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsErrorsAdded) IsTranslationKey() {
}
func (k KeysCommandsVipsErrorsAdded) GetPath() string {
	return "commands.vips.errors.added"
}
func (k KeysCommandsVipsErrorsAdded) GetPathSlice() []string {
	return []string{"commands", "vips", "errors", "added"}
}
func (k KeysCommandsVipsErrorsAdded) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsErrorsAdded) SetVars(vars KeysCommandsVipsErrorsAddedVars) twiri18n.TranslationKey[KeysCommandsVipsErrorsAddedVars] {
	k.Vars = twiri18n.Vars{"userName": vars.UserName}
	return k
}

type KeysCommandsVipsErrorsNoScheduledVipsVars struct {
}
type KeysCommandsVipsErrorsNoScheduledVips struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsErrorsNoScheduledVips) IsTranslationKey() {
}
func (k KeysCommandsVipsErrorsNoScheduledVips) GetPath() string {
	return "commands.vips.errors.no_scheduled_vips"
}
func (k KeysCommandsVipsErrorsNoScheduledVips) GetPathSlice() []string {
	return []string{"commands", "vips", "errors", "no_scheduled_vips"}
}
func (k KeysCommandsVipsErrorsNoScheduledVips) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsErrorsNoScheduledVips) SetVars(vars KeysCommandsVipsErrorsNoScheduledVipsVars) twiri18n.TranslationKey[KeysCommandsVipsErrorsNoScheduledVipsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsErrorsCannotCreateScheduledInDbVars struct {
}
type KeysCommandsVipsErrorsCannotCreateScheduledInDb struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsErrorsCannotCreateScheduledInDb) IsTranslationKey() {
}
func (k KeysCommandsVipsErrorsCannotCreateScheduledInDb) GetPath() string {
	return "commands.vips.errors.cannot_create_scheduled_in_db"
}
func (k KeysCommandsVipsErrorsCannotCreateScheduledInDb) GetPathSlice() []string {
	return []string{"commands", "vips", "errors", "cannot_create_scheduled_in_db"}
}
func (k KeysCommandsVipsErrorsCannotCreateScheduledInDb) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsErrorsCannotCreateScheduledInDb) SetVars(vars KeysCommandsVipsErrorsCannotCreateScheduledInDbVars) twiri18n.TranslationKey[KeysCommandsVipsErrorsCannotCreateScheduledInDbVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsErrorsAddedWithRemoveTimeVars struct {
	UserName	any
	EndTime		any
}
type KeysCommandsVipsErrorsAddedWithRemoveTime struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsErrorsAddedWithRemoveTime) IsTranslationKey() {
}
func (k KeysCommandsVipsErrorsAddedWithRemoveTime) GetPath() string {
	return "commands.vips.errors.added_with_remove_time"
}
func (k KeysCommandsVipsErrorsAddedWithRemoveTime) GetPathSlice() []string {
	return []string{"commands", "vips", "errors", "added_with_remove_time"}
}
func (k KeysCommandsVipsErrorsAddedWithRemoveTime) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsErrorsAddedWithRemoveTime) SetVars(vars KeysCommandsVipsErrorsAddedWithRemoveTimeVars) twiri18n.TranslationKey[KeysCommandsVipsErrorsAddedWithRemoveTimeVars] {
	k.Vars = twiri18n.Vars{"userName": vars.UserName, "endTime": vars.EndTime}
	return k
}

type KeysCommandsVipsErrorsCannotGetListFromDbVars struct {
}
type KeysCommandsVipsErrorsCannotGetListFromDb struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsErrorsCannotGetListFromDb) IsTranslationKey() {
}
func (k KeysCommandsVipsErrorsCannotGetListFromDb) GetPath() string {
	return "commands.vips.errors.cannot_get_list_from_db"
}
func (k KeysCommandsVipsErrorsCannotGetListFromDb) GetPathSlice() []string {
	return []string{"commands", "vips", "errors", "cannot_get_list_from_db"}
}
func (k KeysCommandsVipsErrorsCannotGetListFromDb) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsErrorsCannotGetListFromDb) SetVars(vars KeysCommandsVipsErrorsCannotGetListFromDbVars) twiri18n.TranslationKey[KeysCommandsVipsErrorsCannotGetListFromDbVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsErrorsRemovedVars struct {
	UserName any
}
type KeysCommandsVipsErrorsRemoved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsErrorsRemoved) IsTranslationKey() {
}
func (k KeysCommandsVipsErrorsRemoved) GetPath() string {
	return "commands.vips.errors.removed"
}
func (k KeysCommandsVipsErrorsRemoved) GetPathSlice() []string {
	return []string{"commands", "vips", "errors", "removed"}
}
func (k KeysCommandsVipsErrorsRemoved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsErrorsRemoved) SetVars(vars KeysCommandsVipsErrorsRemovedVars) twiri18n.TranslationKey[KeysCommandsVipsErrorsRemovedVars] {
	k.Vars = twiri18n.Vars{"userName": vars.UserName}
	return k
}

type KeysCommandsVipsErrorsCannotUpdateVars struct {
}
type KeysCommandsVipsErrorsCannotUpdate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsErrorsCannotUpdate) IsTranslationKey() {
}
func (k KeysCommandsVipsErrorsCannotUpdate) GetPath() string {
	return "commands.vips.errors.cannot_update"
}
func (k KeysCommandsVipsErrorsCannotUpdate) GetPathSlice() []string {
	return []string{"commands", "vips", "errors", "cannot_update"}
}
func (k KeysCommandsVipsErrorsCannotUpdate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsErrorsCannotUpdate) SetVars(vars KeysCommandsVipsErrorsCannotUpdateVars) twiri18n.TranslationKey[KeysCommandsVipsErrorsCannotUpdateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsErrorsUpdatedVars struct {
	UserName	any
	EndTime		any
}
type KeysCommandsVipsErrorsUpdated struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsErrorsUpdated) IsTranslationKey() {
}
func (k KeysCommandsVipsErrorsUpdated) GetPath() string {
	return "commands.vips.errors.updated"
}
func (k KeysCommandsVipsErrorsUpdated) GetPathSlice() []string {
	return []string{"commands", "vips", "errors", "updated"}
}
func (k KeysCommandsVipsErrorsUpdated) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsErrorsUpdated) SetVars(vars KeysCommandsVipsErrorsUpdatedVars) twiri18n.TranslationKey[KeysCommandsVipsErrorsUpdatedVars] {
	k.Vars = twiri18n.Vars{"userName": vars.UserName, "endTime": vars.EndTime}
	return k
}

type KeysCommandsVipsErrors struct {
	InvalidDuration			KeysCommandsVipsErrorsInvalidDuration
	AlreadyHaveRole			KeysCommandsVipsErrorsAlreadyHaveRole
	Added				KeysCommandsVipsErrorsAdded
	NoScheduledVips			KeysCommandsVipsErrorsNoScheduledVips
	CannotCreateScheduledInDb	KeysCommandsVipsErrorsCannotCreateScheduledInDb
	AddedWithRemoveTime		KeysCommandsVipsErrorsAddedWithRemoveTime
	CannotGetListFromDb		KeysCommandsVipsErrorsCannotGetListFromDb
	Removed				KeysCommandsVipsErrorsRemoved
	CannotUpdate			KeysCommandsVipsErrorsCannotUpdate
	Updated				KeysCommandsVipsErrorsUpdated
}
type KeysCommandsVipsHintsUserVars struct {
}
type KeysCommandsVipsHintsUser struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsHintsUser) IsTranslationKey() {
}
func (k KeysCommandsVipsHintsUser) GetPath() string {
	return "commands.vips.hints.user"
}
func (k KeysCommandsVipsHintsUser) GetPathSlice() []string {
	return []string{"commands", "vips", "hints", "user"}
}
func (k KeysCommandsVipsHintsUser) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsHintsUser) SetVars(vars KeysCommandsVipsHintsUserVars) twiri18n.TranslationKey[KeysCommandsVipsHintsUserVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsHintsUnvipInVars struct {
}
type KeysCommandsVipsHintsUnvipIn struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsVipsHintsUnvipIn) IsTranslationKey() {
}
func (k KeysCommandsVipsHintsUnvipIn) GetPath() string {
	return "commands.vips.hints.unvip_in"
}
func (k KeysCommandsVipsHintsUnvipIn) GetPathSlice() []string {
	return []string{"commands", "vips", "hints", "unvip_in"}
}
func (k KeysCommandsVipsHintsUnvipIn) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsVipsHintsUnvipIn) SetVars(vars KeysCommandsVipsHintsUnvipInVars) twiri18n.TranslationKey[KeysCommandsVipsHintsUnvipInVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsVipsHints struct {
	User	KeysCommandsVipsHintsUser
	UnvipIn	KeysCommandsVipsHintsUnvipIn
}
type KeysCommandsVips struct {
	Errors	KeysCommandsVipsErrors
	Hints	KeysCommandsVipsHints
}
type KeysCommandsChannelHintsGameArgNameVars struct {
}
type KeysCommandsChannelHintsGameArgName struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelHintsGameArgName) IsTranslationKey() {
}
func (k KeysCommandsChannelHintsGameArgName) GetPath() string {
	return "commands.channel.hints.gameArgName"
}
func (k KeysCommandsChannelHintsGameArgName) GetPathSlice() []string {
	return []string{"commands", "channel", "hints", "gameArgName"}
}
func (k KeysCommandsChannelHintsGameArgName) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelHintsGameArgName) SetVars(vars KeysCommandsChannelHintsGameArgNameVars) twiri18n.TranslationKey[KeysCommandsChannelHintsGameArgNameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelHints struct {
	GameArgName KeysCommandsChannelHintsGameArgName
}
type KeysCommandsChannelAddCategoryChangeVars struct {
	CategoryName any
}
type KeysCommandsChannelAddCategoryChange struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelAddCategoryChange) IsTranslationKey() {
}
func (k KeysCommandsChannelAddCategoryChange) GetPath() string {
	return "commands.channel.add.category_change"
}
func (k KeysCommandsChannelAddCategoryChange) GetPathSlice() []string {
	return []string{"commands", "channel", "add", "category_change"}
}
func (k KeysCommandsChannelAddCategoryChange) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelAddCategoryChange) SetVars(vars KeysCommandsChannelAddCategoryChangeVars) twiri18n.TranslationKey[KeysCommandsChannelAddCategoryChangeVars] {
	k.Vars = twiri18n.Vars{"categoryName": vars.CategoryName}
	return k
}

type KeysCommandsChannelAdd struct {
	CategoryChange KeysCommandsChannelAddCategoryChange
}
type KeysCommandsChannelErrorsCategoryCannotChangeErrorVars struct {
	ErrorMessage any
}
type KeysCommandsChannelErrorsCategoryCannotChangeError struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsCategoryCannotChangeError) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsCategoryCannotChangeError) GetPath() string {
	return "commands.channel.errors.category_cannot_change_error"
}
func (k KeysCommandsChannelErrorsCategoryCannotChangeError) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "category_cannot_change_error"}
}
func (k KeysCommandsChannelErrorsCategoryCannotChangeError) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsCategoryCannotChangeError) SetVars(vars KeysCommandsChannelErrorsCategoryCannotChangeErrorVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsCategoryCannotChangeErrorVars] {
	k.Vars = twiri18n.Vars{"errorMessage": vars.ErrorMessage}
	return k
}

type KeysCommandsChannelErrorsBroadcasterTwitchClientCannotCreateVars struct {
}
type KeysCommandsChannelErrorsBroadcasterTwitchClientCannotCreate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsBroadcasterTwitchClientCannotCreate) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsBroadcasterTwitchClientCannotCreate) GetPath() string {
	return "commands.channel.errors.broadcaster_twitch_client_cannot_create"
}
func (k KeysCommandsChannelErrorsBroadcasterTwitchClientCannotCreate) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "broadcaster_twitch_client_cannot_create"}
}
func (k KeysCommandsChannelErrorsBroadcasterTwitchClientCannotCreate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsBroadcasterTwitchClientCannotCreate) SetVars(vars KeysCommandsChannelErrorsBroadcasterTwitchClientCannotCreateVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsBroadcasterTwitchClientCannotCreateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelErrorsChannelNotFoundVars struct {
}
type KeysCommandsChannelErrorsChannelNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsChannelNotFound) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsChannelNotFound) GetPath() string {
	return "commands.channel.errors.channel_not_found"
}
func (k KeysCommandsChannelErrorsChannelNotFound) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "channel_not_found"}
}
func (k KeysCommandsChannelErrorsChannelNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsChannelNotFound) SetVars(vars KeysCommandsChannelErrorsChannelNotFoundVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsChannelNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelErrorsCategoryNotFoundVars struct {
}
type KeysCommandsChannelErrorsCategoryNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsCategoryNotFound) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsCategoryNotFound) GetPath() string {
	return "commands.channel.errors.category_not_found"
}
func (k KeysCommandsChannelErrorsCategoryNotFound) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "category_not_found"}
}
func (k KeysCommandsChannelErrorsCategoryNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsCategoryNotFound) SetVars(vars KeysCommandsChannelErrorsCategoryNotFoundVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsCategoryNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelErrorsGameNotFoundVars struct {
}
type KeysCommandsChannelErrorsGameNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsGameNotFound) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsGameNotFound) GetPath() string {
	return "commands.channel.errors.game_not_found"
}
func (k KeysCommandsChannelErrorsGameNotFound) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "game_not_found"}
}
func (k KeysCommandsChannelErrorsGameNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsGameNotFound) SetVars(vars KeysCommandsChannelErrorsGameNotFoundVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsGameNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelErrorsHistoryGameMessageVars struct {
}
type KeysCommandsChannelErrorsHistoryGameMessage struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsHistoryGameMessage) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsHistoryGameMessage) GetPath() string {
	return "commands.channel.errors.history_game_message"
}
func (k KeysCommandsChannelErrorsHistoryGameMessage) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "history_game_message"}
}
func (k KeysCommandsChannelErrorsHistoryGameMessage) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsHistoryGameMessage) SetVars(vars KeysCommandsChannelErrorsHistoryGameMessageVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsHistoryGameMessageVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelErrorsHistoryTitleMessageVars struct {
	Reason any
}
type KeysCommandsChannelErrorsHistoryTitleMessage struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsHistoryTitleMessage) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsHistoryTitleMessage) GetPath() string {
	return "commands.channel.errors.history_title_message"
}
func (k KeysCommandsChannelErrorsHistoryTitleMessage) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "history_title_message"}
}
func (k KeysCommandsChannelErrorsHistoryTitleMessage) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsHistoryTitleMessage) SetVars(vars KeysCommandsChannelErrorsHistoryTitleMessageVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsHistoryTitleMessageVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsChannelErrorsBroadcasterTwitchApiClientVars struct {
	Reason any
}
type KeysCommandsChannelErrorsBroadcasterTwitchApiClient struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsBroadcasterTwitchApiClient) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsBroadcasterTwitchApiClient) GetPath() string {
	return "commands.channel.errors.broadcaster_twitch_api_client"
}
func (k KeysCommandsChannelErrorsBroadcasterTwitchApiClient) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "broadcaster_twitch_api_client"}
}
func (k KeysCommandsChannelErrorsBroadcasterTwitchApiClient) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsBroadcasterTwitchApiClient) SetVars(vars KeysCommandsChannelErrorsBroadcasterTwitchApiClientVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsBroadcasterTwitchApiClientVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsChannelErrorsChannelCannotGetInformationVars struct {
}
type KeysCommandsChannelErrorsChannelCannotGetInformation struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsChannelCannotGetInformation) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsChannelCannotGetInformation) GetPath() string {
	return "commands.channel.errors.channel_cannot_get_information"
}
func (k KeysCommandsChannelErrorsChannelCannotGetInformation) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "channel_cannot_get_information"}
}
func (k KeysCommandsChannelErrorsChannelCannotGetInformation) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsChannelCannotGetInformation) SetVars(vars KeysCommandsChannelErrorsChannelCannotGetInformationVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsChannelCannotGetInformationVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelErrorsAliasCannotGetCategoryVars struct {
}
type KeysCommandsChannelErrorsAliasCannotGetCategory struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsAliasCannotGetCategory) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsAliasCannotGetCategory) GetPath() string {
	return "commands.channel.errors.alias_cannot_get_category"
}
func (k KeysCommandsChannelErrorsAliasCannotGetCategory) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "alias_cannot_get_category"}
}
func (k KeysCommandsChannelErrorsAliasCannotGetCategory) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsAliasCannotGetCategory) SetVars(vars KeysCommandsChannelErrorsAliasCannotGetCategoryVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsAliasCannotGetCategoryVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelErrorsCategoryCannotGetVars struct {
}
type KeysCommandsChannelErrorsCategoryCannotGet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsCategoryCannotGet) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsCategoryCannotGet) GetPath() string {
	return "commands.channel.errors.category_cannot_get"
}
func (k KeysCommandsChannelErrorsCategoryCannotGet) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "category_cannot_get"}
}
func (k KeysCommandsChannelErrorsCategoryCannotGet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsCategoryCannotGet) SetVars(vars KeysCommandsChannelErrorsCategoryCannotGetVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsCategoryCannotGetVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelErrorsCategoryCannotGetErrorVars struct {
	ErrorMessage any
}
type KeysCommandsChannelErrorsCategoryCannotGetError struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsCategoryCannotGetError) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsCategoryCannotGetError) GetPath() string {
	return "commands.channel.errors.category_cannot_get_error"
}
func (k KeysCommandsChannelErrorsCategoryCannotGetError) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "category_cannot_get_error"}
}
func (k KeysCommandsChannelErrorsCategoryCannotGetError) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsCategoryCannotGetError) SetVars(vars KeysCommandsChannelErrorsCategoryCannotGetErrorVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsCategoryCannotGetErrorVars] {
	k.Vars = twiri18n.Vars{"errorMessage": vars.ErrorMessage}
	return k
}

type KeysCommandsChannelErrorsCategoryCannotChangeVars struct {
}
type KeysCommandsChannelErrorsCategoryCannotChange struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChannelErrorsCategoryCannotChange) IsTranslationKey() {
}
func (k KeysCommandsChannelErrorsCategoryCannotChange) GetPath() string {
	return "commands.channel.errors.category_cannot_change"
}
func (k KeysCommandsChannelErrorsCategoryCannotChange) GetPathSlice() []string {
	return []string{"commands", "channel", "errors", "category_cannot_change"}
}
func (k KeysCommandsChannelErrorsCategoryCannotChange) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChannelErrorsCategoryCannotChange) SetVars(vars KeysCommandsChannelErrorsCategoryCannotChangeVars) twiri18n.TranslationKey[KeysCommandsChannelErrorsCategoryCannotChangeVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChannelErrors struct {
	CategoryCannotChangeError		KeysCommandsChannelErrorsCategoryCannotChangeError
	BroadcasterTwitchClientCannotCreate	KeysCommandsChannelErrorsBroadcasterTwitchClientCannotCreate
	ChannelNotFound				KeysCommandsChannelErrorsChannelNotFound
	CategoryNotFound			KeysCommandsChannelErrorsCategoryNotFound
	GameNotFound				KeysCommandsChannelErrorsGameNotFound
	HistoryGameMessage			KeysCommandsChannelErrorsHistoryGameMessage
	HistoryTitleMessage			KeysCommandsChannelErrorsHistoryTitleMessage
	BroadcasterTwitchApiClient		KeysCommandsChannelErrorsBroadcasterTwitchApiClient
	ChannelCannotGetInformation		KeysCommandsChannelErrorsChannelCannotGetInformation
	AliasCannotGetCategory			KeysCommandsChannelErrorsAliasCannotGetCategory
	CategoryCannotGet			KeysCommandsChannelErrorsCategoryCannotGet
	CategoryCannotGetError			KeysCommandsChannelErrorsCategoryCannotGetError
	CategoryCannotChange			KeysCommandsChannelErrorsCategoryCannotChange
}
type KeysCommandsChannel struct {
	Hints	KeysCommandsChannelHints
	Add	KeysCommandsChannelAdd
	Errors	KeysCommandsChannelErrors
}
type KeysCommandsClipEmptyClipUrlVars struct {
}
type KeysCommandsClipEmptyClipUrl struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsClipEmptyClipUrl) IsTranslationKey() {
}
func (k KeysCommandsClipEmptyClipUrl) GetPath() string {
	return "commands.clip.empty_clip_url"
}
func (k KeysCommandsClipEmptyClipUrl) GetPathSlice() []string {
	return []string{"commands", "clip", "empty_clip_url"}
}
func (k KeysCommandsClipEmptyClipUrl) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsClipEmptyClipUrl) SetVars(vars KeysCommandsClipEmptyClipUrlVars) twiri18n.TranslationKey[KeysCommandsClipEmptyClipUrlVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsClipCannotGetClipVars struct {
}
type KeysCommandsClipCannotGetClip struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsClipCannotGetClip) IsTranslationKey() {
}
func (k KeysCommandsClipCannotGetClip) GetPath() string {
	return "commands.clip.cannot_get_clip"
}
func (k KeysCommandsClipCannotGetClip) GetPathSlice() []string {
	return []string{"commands", "clip", "cannot_get_clip"}
}
func (k KeysCommandsClipCannotGetClip) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsClipCannotGetClip) SetVars(vars KeysCommandsClipCannotGetClipVars) twiri18n.TranslationKey[KeysCommandsClipCannotGetClipVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsClipClipCreatedVars struct {
	Url any
}
type KeysCommandsClipClipCreated struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsClipClipCreated) IsTranslationKey() {
}
func (k KeysCommandsClipClipCreated) GetPath() string {
	return "commands.clip.clip_created"
}
func (k KeysCommandsClipClipCreated) GetPathSlice() []string {
	return []string{"commands", "clip", "clip_created"}
}
func (k KeysCommandsClipClipCreated) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsClipClipCreated) SetVars(vars KeysCommandsClipClipCreatedVars) twiri18n.TranslationKey[KeysCommandsClipClipCreatedVars] {
	k.Vars = twiri18n.Vars{"url": vars.Url}
	return k
}

type KeysCommandsClipCannotCreateClipVars struct {
}
type KeysCommandsClipCannotCreateClip struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsClipCannotCreateClip) IsTranslationKey() {
}
func (k KeysCommandsClipCannotCreateClip) GetPath() string {
	return "commands.clip.cannot_create_clip"
}
func (k KeysCommandsClipCannotCreateClip) GetPathSlice() []string {
	return []string{"commands", "clip", "cannot_create_clip"}
}
func (k KeysCommandsClipCannotCreateClip) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsClipCannotCreateClip) SetVars(vars KeysCommandsClipCannotCreateClipVars) twiri18n.TranslationKey[KeysCommandsClipCannotCreateClipVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsClip struct {
	EmptyClipUrl		KeysCommandsClipEmptyClipUrl
	CannotGetClip		KeysCommandsClipCannotGetClip
	ClipCreated		KeysCommandsClipClipCreated
	CannotCreateClip	KeysCommandsClipCannotCreateClip
}
type KeysCommandsGamesErrorsDuelCannotGetSenderVars struct {
}
type KeysCommandsGamesErrorsDuelCannotGetSender struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelCannotGetSender) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelCannotGetSender) GetPath() string {
	return "commands.games.errors.duel_cannot_get_sender"
}
func (k KeysCommandsGamesErrorsDuelCannotGetSender) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_cannot_get_sender"}
}
func (k KeysCommandsGamesErrorsDuelCannotGetSender) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelCannotGetSender) SetVars(vars KeysCommandsGamesErrorsDuelCannotGetSenderVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelCannotGetSenderVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsDuelWithYourselfVars struct {
}
type KeysCommandsGamesErrorsDuelWithYourself struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelWithYourself) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelWithYourself) GetPath() string {
	return "commands.games.errors.duel_with_yourself"
}
func (k KeysCommandsGamesErrorsDuelWithYourself) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_with_yourself"}
}
func (k KeysCommandsGamesErrorsDuelWithYourself) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelWithYourself) SetVars(vars KeysCommandsGamesErrorsDuelWithYourselfVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelWithYourselfVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsDuelCannotSetGlobalCooldownVars struct {
	Reason any
}
type KeysCommandsGamesErrorsDuelCannotSetGlobalCooldown struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelCannotSetGlobalCooldown) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelCannotSetGlobalCooldown) GetPath() string {
	return "commands.games.errors.duel_cannot_set_global_cooldown"
}
func (k KeysCommandsGamesErrorsDuelCannotSetGlobalCooldown) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_cannot_set_global_cooldown"}
}
func (k KeysCommandsGamesErrorsDuelCannotSetGlobalCooldown) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelCannotSetGlobalCooldown) SetVars(vars KeysCommandsGamesErrorsDuelCannotSetGlobalCooldownVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelCannotSetGlobalCooldownVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsGamesErrorsDuelCannotSaveToCacheVars struct {
}
type KeysCommandsGamesErrorsDuelCannotSaveToCache struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelCannotSaveToCache) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelCannotSaveToCache) GetPath() string {
	return "commands.games.errors.duel_cannot_save_to_cache"
}
func (k KeysCommandsGamesErrorsDuelCannotSaveToCache) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_cannot_save_to_cache"}
}
func (k KeysCommandsGamesErrorsDuelCannotSaveToCache) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelCannotSaveToCache) SetVars(vars KeysCommandsGamesErrorsDuelCannotSaveToCacheVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelCannotSaveToCacheVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsRouletteCannotSendInitialMessageVars struct {
}
type KeysCommandsGamesErrorsRouletteCannotSendInitialMessage struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsRouletteCannotSendInitialMessage) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsRouletteCannotSendInitialMessage) GetPath() string {
	return "commands.games.errors.roulette_cannot_send_initial_message"
}
func (k KeysCommandsGamesErrorsRouletteCannotSendInitialMessage) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "roulette_cannot_send_initial_message"}
}
func (k KeysCommandsGamesErrorsRouletteCannotSendInitialMessage) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsRouletteCannotSendInitialMessage) SetVars(vars KeysCommandsGamesErrorsRouletteCannotSendInitialMessageVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsRouletteCannotSendInitialMessageVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsSeppukuCannotFindSettingsVars struct {
}
type KeysCommandsGamesErrorsSeppukuCannotFindSettings struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsSeppukuCannotFindSettings) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsSeppukuCannotFindSettings) GetPath() string {
	return "commands.games.errors.seppuku_cannot_find_settings"
}
func (k KeysCommandsGamesErrorsSeppukuCannotFindSettings) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "seppuku_cannot_find_settings"}
}
func (k KeysCommandsGamesErrorsSeppukuCannotFindSettings) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsSeppukuCannotFindSettings) SetVars(vars KeysCommandsGamesErrorsSeppukuCannotFindSettingsVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsSeppukuCannotFindSettingsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsVotebanCannotCheckProgressVars struct {
}
type KeysCommandsGamesErrorsVotebanCannotCheckProgress struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsVotebanCannotCheckProgress) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsVotebanCannotCheckProgress) GetPath() string {
	return "commands.games.errors.voteban_cannot_check_progress"
}
func (k KeysCommandsGamesErrorsVotebanCannotCheckProgress) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "voteban_cannot_check_progress"}
}
func (k KeysCommandsGamesErrorsVotebanCannotCheckProgress) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsVotebanCannotCheckProgress) SetVars(vars KeysCommandsGamesErrorsVotebanCannotCheckProgressVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsVotebanCannotCheckProgressVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsVotebanCannotSetVoteExpirationVars struct {
}
type KeysCommandsGamesErrorsVotebanCannotSetVoteExpiration struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsVotebanCannotSetVoteExpiration) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsVotebanCannotSetVoteExpiration) GetPath() string {
	return "commands.games.errors.voteban_cannot_set_vote_expiration"
}
func (k KeysCommandsGamesErrorsVotebanCannotSetVoteExpiration) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "voteban_cannot_set_vote_expiration"}
}
func (k KeysCommandsGamesErrorsVotebanCannotSetVoteExpiration) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsVotebanCannotSetVoteExpiration) SetVars(vars KeysCommandsGamesErrorsVotebanCannotSetVoteExpirationVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsVotebanCannotSetVoteExpirationVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsDuelCannotSaveResultVars struct {
	Reason any
}
type KeysCommandsGamesErrorsDuelCannotSaveResult struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelCannotSaveResult) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelCannotSaveResult) GetPath() string {
	return "commands.games.errors.duel_cannot_save_result"
}
func (k KeysCommandsGamesErrorsDuelCannotSaveResult) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_cannot_save_result"}
}
func (k KeysCommandsGamesErrorsDuelCannotSaveResult) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelCannotSaveResult) SetVars(vars KeysCommandsGamesErrorsDuelCannotSaveResultVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelCannotSaveResultVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsGamesErrorsDuelCannotSaveDataVars struct {
	Reason any
}
type KeysCommandsGamesErrorsDuelCannotSaveData struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelCannotSaveData) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelCannotSaveData) GetPath() string {
	return "commands.games.errors.duel_cannot_save_data"
}
func (k KeysCommandsGamesErrorsDuelCannotSaveData) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_cannot_save_data"}
}
func (k KeysCommandsGamesErrorsDuelCannotSaveData) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelCannotSaveData) SetVars(vars KeysCommandsGamesErrorsDuelCannotSaveDataVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelCannotSaveDataVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsGamesErrorsDuelWithStreamerVars struct {
}
type KeysCommandsGamesErrorsDuelWithStreamer struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelWithStreamer) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelWithStreamer) GetPath() string {
	return "commands.games.errors.duel_with_streamer"
}
func (k KeysCommandsGamesErrorsDuelWithStreamer) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_with_streamer"}
}
func (k KeysCommandsGamesErrorsDuelWithStreamer) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelWithStreamer) SetVars(vars KeysCommandsGamesErrorsDuelWithStreamerVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelWithStreamerVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsDuelWithBotVars struct {
}
type KeysCommandsGamesErrorsDuelWithBot struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelWithBot) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelWithBot) GetPath() string {
	return "commands.games.errors.duel_with_bot"
}
func (k KeysCommandsGamesErrorsDuelWithBot) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_with_bot"}
}
func (k KeysCommandsGamesErrorsDuelWithBot) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelWithBot) SetVars(vars KeysCommandsGamesErrorsDuelWithBotVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelWithBotVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsDuelCannotCheckCooldownVars struct {
	Reason any
}
type KeysCommandsGamesErrorsDuelCannotCheckCooldown struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelCannotCheckCooldown) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelCannotCheckCooldown) GetPath() string {
	return "commands.games.errors.duel_cannot_check_cooldown"
}
func (k KeysCommandsGamesErrorsDuelCannotCheckCooldown) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_cannot_check_cooldown"}
}
func (k KeysCommandsGamesErrorsDuelCannotCheckCooldown) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelCannotCheckCooldown) SetVars(vars KeysCommandsGamesErrorsDuelCannotCheckCooldownVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelCannotCheckCooldownVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsGamesErrorsDuelCannotValidateParticipantsVars struct {
}
type KeysCommandsGamesErrorsDuelCannotValidateParticipants struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelCannotValidateParticipants) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelCannotValidateParticipants) GetPath() string {
	return "commands.games.errors.duel_cannot_validate_participants"
}
func (k KeysCommandsGamesErrorsDuelCannotValidateParticipants) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_cannot_validate_participants"}
}
func (k KeysCommandsGamesErrorsDuelCannotValidateParticipants) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelCannotValidateParticipants) SetVars(vars KeysCommandsGamesErrorsDuelCannotValidateParticipantsVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelCannotValidateParticipantsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsRouletteCannotGetWithSettingsVars struct {
}
type KeysCommandsGamesErrorsRouletteCannotGetWithSettings struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsRouletteCannotGetWithSettings) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsRouletteCannotGetWithSettings) GetPath() string {
	return "commands.games.errors.roulette_cannot_get_with_settings"
}
func (k KeysCommandsGamesErrorsRouletteCannotGetWithSettings) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "roulette_cannot_get_with_settings"}
}
func (k KeysCommandsGamesErrorsRouletteCannotGetWithSettings) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsRouletteCannotGetWithSettings) SetVars(vars KeysCommandsGamesErrorsRouletteCannotGetWithSettingsVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsRouletteCannotGetWithSettingsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsVotebanCannotLockVars struct {
}
type KeysCommandsGamesErrorsVotebanCannotLock struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsVotebanCannotLock) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsVotebanCannotLock) GetPath() string {
	return "commands.games.errors.voteban_cannot_lock"
}
func (k KeysCommandsGamesErrorsVotebanCannotLock) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "voteban_cannot_lock"}
}
func (k KeysCommandsGamesErrorsVotebanCannotLock) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsVotebanCannotLock) SetVars(vars KeysCommandsGamesErrorsVotebanCannotLockVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsVotebanCannotLockVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsDuelCannotGetWithSettingsVars struct {
}
type KeysCommandsGamesErrorsDuelCannotGetWithSettings struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelCannotGetWithSettings) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelCannotGetWithSettings) GetPath() string {
	return "commands.games.errors.duel_cannot_get_with_settings"
}
func (k KeysCommandsGamesErrorsDuelCannotGetWithSettings) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_cannot_get_with_settings"}
}
func (k KeysCommandsGamesErrorsDuelCannotGetWithSettings) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelCannotGetWithSettings) SetVars(vars KeysCommandsGamesErrorsDuelCannotGetWithSettingsVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelCannotGetWithSettingsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsDuelCannotCheckUserVars struct {
	Reason any
}
type KeysCommandsGamesErrorsDuelCannotCheckUser struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelCannotCheckUser) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelCannotCheckUser) GetPath() string {
	return "commands.games.errors.duel_cannot_check_user"
}
func (k KeysCommandsGamesErrorsDuelCannotCheckUser) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_cannot_check_user"}
}
func (k KeysCommandsGamesErrorsDuelCannotCheckUser) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelCannotCheckUser) SetVars(vars KeysCommandsGamesErrorsDuelCannotCheckUserVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelCannotCheckUserVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsGamesErrorsDuelCannotSetUserCooldownVars struct {
	Reason any
}
type KeysCommandsGamesErrorsDuelCannotSetUserCooldown struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsDuelCannotSetUserCooldown) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsDuelCannotSetUserCooldown) GetPath() string {
	return "commands.games.errors.duel_cannot_set_user_cooldown"
}
func (k KeysCommandsGamesErrorsDuelCannotSetUserCooldown) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "duel_cannot_set_user_cooldown"}
}
func (k KeysCommandsGamesErrorsDuelCannotSetUserCooldown) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsDuelCannotSetUserCooldown) SetVars(vars KeysCommandsGamesErrorsDuelCannotSetUserCooldownVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsDuelCannotSetUserCooldownVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsGamesErrorsRouletteCannotSendDeathMessageVars struct {
}
type KeysCommandsGamesErrorsRouletteCannotSendDeathMessage struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsRouletteCannotSendDeathMessage) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsRouletteCannotSendDeathMessage) GetPath() string {
	return "commands.games.errors.roulette_cannot_send_death_message"
}
func (k KeysCommandsGamesErrorsRouletteCannotSendDeathMessage) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "roulette_cannot_send_death_message"}
}
func (k KeysCommandsGamesErrorsRouletteCannotSendDeathMessage) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsRouletteCannotSendDeathMessage) SetVars(vars KeysCommandsGamesErrorsRouletteCannotSendDeathMessageVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsRouletteCannotSendDeathMessageVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsVotebanCannotFindSettingsVars struct {
}
type KeysCommandsGamesErrorsVotebanCannotFindSettings struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsVotebanCannotFindSettings) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsVotebanCannotFindSettings) GetPath() string {
	return "commands.games.errors.voteban_cannot_find_settings"
}
func (k KeysCommandsGamesErrorsVotebanCannotFindSettings) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "voteban_cannot_find_settings"}
}
func (k KeysCommandsGamesErrorsVotebanCannotFindSettings) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsVotebanCannotFindSettings) SetVars(vars KeysCommandsGamesErrorsVotebanCannotFindSettingsVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsVotebanCannotFindSettingsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsVotebanCannotFindUserVars struct {
}
type KeysCommandsGamesErrorsVotebanCannotFindUser struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsVotebanCannotFindUser) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsVotebanCannotFindUser) GetPath() string {
	return "commands.games.errors.voteban_cannot_find_user"
}
func (k KeysCommandsGamesErrorsVotebanCannotFindUser) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "voteban_cannot_find_user"}
}
func (k KeysCommandsGamesErrorsVotebanCannotFindUser) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsVotebanCannotFindUser) SetVars(vars KeysCommandsGamesErrorsVotebanCannotFindUserVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsVotebanCannotFindUserVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsVotebanCannotSetVoteVars struct {
}
type KeysCommandsGamesErrorsVotebanCannotSetVote struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsVotebanCannotSetVote) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsVotebanCannotSetVote) GetPath() string {
	return "commands.games.errors.voteban_cannot_set_vote"
}
func (k KeysCommandsGamesErrorsVotebanCannotSetVote) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "voteban_cannot_set_vote"}
}
func (k KeysCommandsGamesErrorsVotebanCannotSetVote) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsVotebanCannotSetVote) SetVars(vars KeysCommandsGamesErrorsVotebanCannotSetVoteVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsVotebanCannotSetVoteVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrorsEightballCannotFindVars struct {
}
type KeysCommandsGamesErrorsEightballCannotFind struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesErrorsEightballCannotFind) IsTranslationKey() {
}
func (k KeysCommandsGamesErrorsEightballCannotFind) GetPath() string {
	return "commands.games.errors.8ball_cannot_find"
}
func (k KeysCommandsGamesErrorsEightballCannotFind) GetPathSlice() []string {
	return []string{"commands", "games", "errors", "8ball_cannot_find"}
}
func (k KeysCommandsGamesErrorsEightballCannotFind) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesErrorsEightballCannotFind) SetVars(vars KeysCommandsGamesErrorsEightballCannotFindVars) twiri18n.TranslationKey[KeysCommandsGamesErrorsEightballCannotFindVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesErrors struct {
	DuelCannotGetSender			KeysCommandsGamesErrorsDuelCannotGetSender
	DuelWithYourself			KeysCommandsGamesErrorsDuelWithYourself
	DuelCannotSetGlobalCooldown		KeysCommandsGamesErrorsDuelCannotSetGlobalCooldown
	DuelCannotSaveToCache			KeysCommandsGamesErrorsDuelCannotSaveToCache
	RouletteCannotSendInitialMessage	KeysCommandsGamesErrorsRouletteCannotSendInitialMessage
	SeppukuCannotFindSettings		KeysCommandsGamesErrorsSeppukuCannotFindSettings
	VotebanCannotCheckProgress		KeysCommandsGamesErrorsVotebanCannotCheckProgress
	VotebanCannotSetVoteExpiration		KeysCommandsGamesErrorsVotebanCannotSetVoteExpiration
	DuelCannotSaveResult			KeysCommandsGamesErrorsDuelCannotSaveResult
	DuelCannotSaveData			KeysCommandsGamesErrorsDuelCannotSaveData
	DuelWithStreamer			KeysCommandsGamesErrorsDuelWithStreamer
	DuelWithBot				KeysCommandsGamesErrorsDuelWithBot
	DuelCannotCheckCooldown			KeysCommandsGamesErrorsDuelCannotCheckCooldown
	DuelCannotValidateParticipants		KeysCommandsGamesErrorsDuelCannotValidateParticipants
	RouletteCannotGetWithSettings		KeysCommandsGamesErrorsRouletteCannotGetWithSettings
	VotebanCannotLock			KeysCommandsGamesErrorsVotebanCannotLock
	DuelCannotGetWithSettings		KeysCommandsGamesErrorsDuelCannotGetWithSettings
	DuelCannotCheckUser			KeysCommandsGamesErrorsDuelCannotCheckUser
	DuelCannotSetUserCooldown		KeysCommandsGamesErrorsDuelCannotSetUserCooldown
	RouletteCannotSendDeathMessage		KeysCommandsGamesErrorsRouletteCannotSendDeathMessage
	VotebanCannotFindSettings		KeysCommandsGamesErrorsVotebanCannotFindSettings
	VotebanCannotFindUser			KeysCommandsGamesErrorsVotebanCannotFindUser
	VotebanCannotSetVote			KeysCommandsGamesErrorsVotebanCannotSetVote
	EightballCannotFind			KeysCommandsGamesErrorsEightballCannotFind
}
type KeysCommandsGamesInfoUserNotParticipateVars struct {
}
type KeysCommandsGamesInfoUserNotParticipate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesInfoUserNotParticipate) IsTranslationKey() {
}
func (k KeysCommandsGamesInfoUserNotParticipate) GetPath() string {
	return "commands.games.info.user_not_participate"
}
func (k KeysCommandsGamesInfoUserNotParticipate) GetPathSlice() []string {
	return []string{"commands", "games", "info", "user_not_participate"}
}
func (k KeysCommandsGamesInfoUserNotParticipate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesInfoUserNotParticipate) SetVars(vars KeysCommandsGamesInfoUserNotParticipateVars) twiri18n.TranslationKey[KeysCommandsGamesInfoUserNotParticipateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesInfoUserAlreadyInDuelVars struct {
}
type KeysCommandsGamesInfoUserAlreadyInDuel struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesInfoUserAlreadyInDuel) IsTranslationKey() {
}
func (k KeysCommandsGamesInfoUserAlreadyInDuel) GetPath() string {
	return "commands.games.info.user_already_in_duel"
}
func (k KeysCommandsGamesInfoUserAlreadyInDuel) GetPathSlice() []string {
	return []string{"commands", "games", "info", "user_already_in_duel"}
}
func (k KeysCommandsGamesInfoUserAlreadyInDuel) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesInfoUserAlreadyInDuel) SetVars(vars KeysCommandsGamesInfoUserAlreadyInDuelVars) twiri18n.TranslationKey[KeysCommandsGamesInfoUserAlreadyInDuelVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesInfoSenderAlreadyInDuelVars struct {
}
type KeysCommandsGamesInfoSenderAlreadyInDuel struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesInfoSenderAlreadyInDuel) IsTranslationKey() {
}
func (k KeysCommandsGamesInfoSenderAlreadyInDuel) GetPath() string {
	return "commands.games.info.sender_already_in_duel"
}
func (k KeysCommandsGamesInfoSenderAlreadyInDuel) GetPathSlice() []string {
	return []string{"commands", "games", "info", "sender_already_in_duel"}
}
func (k KeysCommandsGamesInfoSenderAlreadyInDuel) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesInfoSenderAlreadyInDuel) SetVars(vars KeysCommandsGamesInfoSenderAlreadyInDuelVars) twiri18n.TranslationKey[KeysCommandsGamesInfoSenderAlreadyInDuelVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesInfoDuelStatsVars struct {
	Duels	any
	Wins	any
	Loses	any
	Winrate	any
}
type KeysCommandsGamesInfoDuelStats struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesInfoDuelStats) IsTranslationKey() {
}
func (k KeysCommandsGamesInfoDuelStats) GetPath() string {
	return "commands.games.info.duel_stats"
}
func (k KeysCommandsGamesInfoDuelStats) GetPathSlice() []string {
	return []string{"commands", "games", "info", "duel_stats"}
}
func (k KeysCommandsGamesInfoDuelStats) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesInfoDuelStats) SetVars(vars KeysCommandsGamesInfoDuelStatsVars) twiri18n.TranslationKey[KeysCommandsGamesInfoDuelStatsVars] {
	k.Vars = twiri18n.Vars{"duels": vars.Duels, "wins": vars.Wins, "loses": vars.Loses, "winrate": vars.Winrate}
	return k
}

type KeysCommandsGamesInfoVotebanInProgressVars struct {
}
type KeysCommandsGamesInfoVotebanInProgress struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsGamesInfoVotebanInProgress) IsTranslationKey() {
}
func (k KeysCommandsGamesInfoVotebanInProgress) GetPath() string {
	return "commands.games.info.voteban_in_progress"
}
func (k KeysCommandsGamesInfoVotebanInProgress) GetPathSlice() []string {
	return []string{"commands", "games", "info", "voteban_in_progress"}
}
func (k KeysCommandsGamesInfoVotebanInProgress) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsGamesInfoVotebanInProgress) SetVars(vars KeysCommandsGamesInfoVotebanInProgressVars) twiri18n.TranslationKey[KeysCommandsGamesInfoVotebanInProgressVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsGamesInfo struct {
	UserNotParticipate	KeysCommandsGamesInfoUserNotParticipate
	UserAlreadyInDuel	KeysCommandsGamesInfoUserAlreadyInDuel
	SenderAlreadyInDuel	KeysCommandsGamesInfoSenderAlreadyInDuel
	DuelStats		KeysCommandsGamesInfoDuelStats
	VotebanInProgress	KeysCommandsGamesInfoVotebanInProgress
}
type KeysCommandsGames struct {
	Errors	KeysCommandsGamesErrors
	Info	KeysCommandsGamesInfo
}
type KeysCommandsCategoriesAliasesErrorsCategoryNotFoundVars struct {
}
type KeysCommandsCategoriesAliasesErrorsCategoryNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesErrorsCategoryNotFound) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryNotFound) GetPath() string {
	return "commands.categories_aliases.errors.category_not_found"
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryNotFound) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "errors", "category_not_found"}
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryNotFound) SetVars(vars KeysCommandsCategoriesAliasesErrorsCategoryNotFoundVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesErrorsCategoryNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesErrorsTwitchClientCannotToCreateVars struct {
}
type KeysCommandsCategoriesAliasesErrorsTwitchClientCannotToCreate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesErrorsTwitchClientCannotToCreate) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesErrorsTwitchClientCannotToCreate) GetPath() string {
	return "commands.categories_aliases.errors.twitch_client_cannot_to_create"
}
func (k KeysCommandsCategoriesAliasesErrorsTwitchClientCannotToCreate) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "errors", "twitch_client_cannot_to_create"}
}
func (k KeysCommandsCategoriesAliasesErrorsTwitchClientCannotToCreate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesErrorsTwitchClientCannotToCreate) SetVars(vars KeysCommandsCategoriesAliasesErrorsTwitchClientCannotToCreateVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesErrorsTwitchClientCannotToCreateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesErrorsGameCannotToGetVars struct {
}
type KeysCommandsCategoriesAliasesErrorsGameCannotToGet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesErrorsGameCannotToGet) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesErrorsGameCannotToGet) GetPath() string {
	return "commands.categories_aliases.errors.game_cannot_to_get"
}
func (k KeysCommandsCategoriesAliasesErrorsGameCannotToGet) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "errors", "game_cannot_to_get"}
}
func (k KeysCommandsCategoriesAliasesErrorsGameCannotToGet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesErrorsGameCannotToGet) SetVars(vars KeysCommandsCategoriesAliasesErrorsGameCannotToGetVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesErrorsGameCannotToGetVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesErrorsCategoryCannotToGetVars struct {
}
type KeysCommandsCategoriesAliasesErrorsCategoryCannotToGet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesErrorsCategoryCannotToGet) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryCannotToGet) GetPath() string {
	return "commands.categories_aliases.errors.category_cannot_to_get"
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryCannotToGet) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "errors", "category_cannot_to_get"}
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryCannotToGet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryCannotToGet) SetVars(vars KeysCommandsCategoriesAliasesErrorsCategoryCannotToGetVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesErrorsCategoryCannotToGetVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesErrorsCategoryFailedToCreateVars struct {
}
type KeysCommandsCategoriesAliasesErrorsCategoryFailedToCreate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesErrorsCategoryFailedToCreate) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryFailedToCreate) GetPath() string {
	return "commands.categories_aliases.errors.category_failed_to_create"
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryFailedToCreate) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "errors", "category_failed_to_create"}
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryFailedToCreate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryFailedToCreate) SetVars(vars KeysCommandsCategoriesAliasesErrorsCategoryFailedToCreateVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesErrorsCategoryFailedToCreateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesErrorsCategoryCannotDeleteVars struct {
}
type KeysCommandsCategoriesAliasesErrorsCategoryCannotDelete struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesErrorsCategoryCannotDelete) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryCannotDelete) GetPath() string {
	return "commands.categories_aliases.errors.category_cannot_delete"
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryCannotDelete) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "errors", "category_cannot_delete"}
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryCannotDelete) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryCannotDelete) SetVars(vars KeysCommandsCategoriesAliasesErrorsCategoryCannotDeleteVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesErrorsCategoryCannotDeleteVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesErrorsAliasAlreadyExistsVars struct {
	AliasName any
}
type KeysCommandsCategoriesAliasesErrorsAliasAlreadyExists struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesErrorsAliasAlreadyExists) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesErrorsAliasAlreadyExists) GetPath() string {
	return "commands.categories_aliases.errors.alias_already_exists"
}
func (k KeysCommandsCategoriesAliasesErrorsAliasAlreadyExists) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "errors", "alias_already_exists"}
}
func (k KeysCommandsCategoriesAliasesErrorsAliasAlreadyExists) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesErrorsAliasAlreadyExists) SetVars(vars KeysCommandsCategoriesAliasesErrorsAliasAlreadyExistsVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesErrorsAliasAlreadyExistsVars] {
	k.Vars = twiri18n.Vars{"aliasName": vars.AliasName}
	return k
}

type KeysCommandsCategoriesAliasesErrorsAliasEmptyVars struct {
}
type KeysCommandsCategoriesAliasesErrorsAliasEmpty struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesErrorsAliasEmpty) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesErrorsAliasEmpty) GetPath() string {
	return "commands.categories_aliases.errors.alias_empty"
}
func (k KeysCommandsCategoriesAliasesErrorsAliasEmpty) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "errors", "alias_empty"}
}
func (k KeysCommandsCategoriesAliasesErrorsAliasEmpty) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesErrorsAliasEmpty) SetVars(vars KeysCommandsCategoriesAliasesErrorsAliasEmptyVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesErrorsAliasEmptyVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesErrorsAliasNotFoundVars struct {
}
type KeysCommandsCategoriesAliasesErrorsAliasNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesErrorsAliasNotFound) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesErrorsAliasNotFound) GetPath() string {
	return "commands.categories_aliases.errors.alias_not_found"
}
func (k KeysCommandsCategoriesAliasesErrorsAliasNotFound) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "errors", "alias_not_found"}
}
func (k KeysCommandsCategoriesAliasesErrorsAliasNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesErrorsAliasNotFound) SetVars(vars KeysCommandsCategoriesAliasesErrorsAliasNotFoundVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesErrorsAliasNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesErrorsAliasRemovedVars struct {
	AliasName any
}
type KeysCommandsCategoriesAliasesErrorsAliasRemoved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesErrorsAliasRemoved) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesErrorsAliasRemoved) GetPath() string {
	return "commands.categories_aliases.errors.alias_removed"
}
func (k KeysCommandsCategoriesAliasesErrorsAliasRemoved) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "errors", "alias_removed"}
}
func (k KeysCommandsCategoriesAliasesErrorsAliasRemoved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesErrorsAliasRemoved) SetVars(vars KeysCommandsCategoriesAliasesErrorsAliasRemovedVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesErrorsAliasRemovedVars] {
	k.Vars = twiri18n.Vars{"aliasName": vars.AliasName}
	return k
}

type KeysCommandsCategoriesAliasesErrorsCategoryRequiredVars struct {
}
type KeysCommandsCategoriesAliasesErrorsCategoryRequired struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesErrorsCategoryRequired) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryRequired) GetPath() string {
	return "commands.categories_aliases.errors.category_required"
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryRequired) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "errors", "category_required"}
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryRequired) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryRequired) SetVars(vars KeysCommandsCategoriesAliasesErrorsCategoryRequiredVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesErrorsCategoryRequiredVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesErrorsCategoryFailedToGetVars struct {
}
type KeysCommandsCategoriesAliasesErrorsCategoryFailedToGet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesErrorsCategoryFailedToGet) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryFailedToGet) GetPath() string {
	return "commands.categories_aliases.errors.category_failed_to_get"
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryFailedToGet) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "errors", "category_failed_to_get"}
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryFailedToGet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesErrorsCategoryFailedToGet) SetVars(vars KeysCommandsCategoriesAliasesErrorsCategoryFailedToGetVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesErrorsCategoryFailedToGetVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsCategoriesAliasesErrors struct {
	CategoryNotFound		KeysCommandsCategoriesAliasesErrorsCategoryNotFound
	TwitchClientCannotToCreate	KeysCommandsCategoriesAliasesErrorsTwitchClientCannotToCreate
	GameCannotToGet			KeysCommandsCategoriesAliasesErrorsGameCannotToGet
	CategoryCannotToGet		KeysCommandsCategoriesAliasesErrorsCategoryCannotToGet
	CategoryFailedToCreate		KeysCommandsCategoriesAliasesErrorsCategoryFailedToCreate
	CategoryCannotDelete		KeysCommandsCategoriesAliasesErrorsCategoryCannotDelete
	AliasAlreadyExists		KeysCommandsCategoriesAliasesErrorsAliasAlreadyExists
	AliasEmpty			KeysCommandsCategoriesAliasesErrorsAliasEmpty
	AliasNotFound			KeysCommandsCategoriesAliasesErrorsAliasNotFound
	AliasRemoved			KeysCommandsCategoriesAliasesErrorsAliasRemoved
	CategoryRequired		KeysCommandsCategoriesAliasesErrorsCategoryRequired
	CategoryFailedToGet		KeysCommandsCategoriesAliasesErrorsCategoryFailedToGet
}
type KeysCommandsCategoriesAliasesAddAliasAddToCategoryVars struct {
	AliasName	any
	CategoryName	any
}
type KeysCommandsCategoriesAliasesAddAliasAddToCategory struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsCategoriesAliasesAddAliasAddToCategory) IsTranslationKey() {
}
func (k KeysCommandsCategoriesAliasesAddAliasAddToCategory) GetPath() string {
	return "commands.categories_aliases.add.alias_add_to_category"
}
func (k KeysCommandsCategoriesAliasesAddAliasAddToCategory) GetPathSlice() []string {
	return []string{"commands", "categories_aliases", "add", "alias_add_to_category"}
}
func (k KeysCommandsCategoriesAliasesAddAliasAddToCategory) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsCategoriesAliasesAddAliasAddToCategory) SetVars(vars KeysCommandsCategoriesAliasesAddAliasAddToCategoryVars) twiri18n.TranslationKey[KeysCommandsCategoriesAliasesAddAliasAddToCategoryVars] {
	k.Vars = twiri18n.Vars{"aliasName": vars.AliasName, "categoryName": vars.CategoryName}
	return k
}

type KeysCommandsCategoriesAliasesAdd struct {
	AliasAddToCategory KeysCommandsCategoriesAliasesAddAliasAddToCategory
}
type KeysCommandsCategoriesAliases struct {
	Errors	KeysCommandsCategoriesAliasesErrors
	Add	KeysCommandsCategoriesAliasesAdd
}
type KeysCommandsPermitErrorsCannotCreateVars struct {
}
type KeysCommandsPermitErrorsCannotCreate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPermitErrorsCannotCreate) IsTranslationKey() {
}
func (k KeysCommandsPermitErrorsCannotCreate) GetPath() string {
	return "commands.permit.errors.cannot_create"
}
func (k KeysCommandsPermitErrorsCannotCreate) GetPathSlice() []string {
	return []string{"commands", "permit", "errors", "cannot_create"}
}
func (k KeysCommandsPermitErrorsCannotCreate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPermitErrorsCannotCreate) SetVars(vars KeysCommandsPermitErrorsCannotCreateVars) twiri18n.TranslationKey[KeysCommandsPermitErrorsCannotCreateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPermitErrors struct {
	CannotCreate KeysCommandsPermitErrorsCannotCreate
}
type KeysCommandsPermitSuccessAddedPermitVars struct {
	CountPermit	any
	UserName	any
}
type KeysCommandsPermitSuccessAddedPermit struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPermitSuccessAddedPermit) IsTranslationKey() {
}
func (k KeysCommandsPermitSuccessAddedPermit) GetPath() string {
	return "commands.permit.success.added_permit"
}
func (k KeysCommandsPermitSuccessAddedPermit) GetPathSlice() []string {
	return []string{"commands", "permit", "success", "added_permit"}
}
func (k KeysCommandsPermitSuccessAddedPermit) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPermitSuccessAddedPermit) SetVars(vars KeysCommandsPermitSuccessAddedPermitVars) twiri18n.TranslationKey[KeysCommandsPermitSuccessAddedPermitVars] {
	k.Vars = twiri18n.Vars{"countPermit": vars.CountPermit, "userName": vars.UserName}
	return k
}

type KeysCommandsPermitSuccess struct {
	AddedPermit KeysCommandsPermitSuccessAddedPermit
}
type KeysCommandsPermit struct {
	Errors	KeysCommandsPermitErrors
	Success	KeysCommandsPermitSuccess
}
type KeysCommandsTtsErrorsWhileGettingVoicesVars struct {
}
type KeysCommandsTtsErrorsWhileGettingVoices struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsErrorsWhileGettingVoices) IsTranslationKey() {
}
func (k KeysCommandsTtsErrorsWhileGettingVoices) GetPath() string {
	return "commands.tts.errors.while_getting_voices"
}
func (k KeysCommandsTtsErrorsWhileGettingVoices) GetPathSlice() []string {
	return []string{"commands", "tts", "errors", "while_getting_voices"}
}
func (k KeysCommandsTtsErrorsWhileGettingVoices) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsErrorsWhileGettingVoices) SetVars(vars KeysCommandsTtsErrorsWhileGettingVoicesVars) twiri18n.TranslationKey[KeysCommandsTtsErrorsWhileGettingVoicesVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsErrorsWhileDisableVars struct {
}
type KeysCommandsTtsErrorsWhileDisable struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsErrorsWhileDisable) IsTranslationKey() {
}
func (k KeysCommandsTtsErrorsWhileDisable) GetPath() string {
	return "commands.tts.errors.while_disable"
}
func (k KeysCommandsTtsErrorsWhileDisable) GetPathSlice() []string {
	return []string{"commands", "tts", "errors", "while_disable"}
}
func (k KeysCommandsTtsErrorsWhileDisable) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsErrorsWhileDisable) SetVars(vars KeysCommandsTtsErrorsWhileDisableVars) twiri18n.TranslationKey[KeysCommandsTtsErrorsWhileDisableVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsErrorsWhileEnableVars struct {
}
type KeysCommandsTtsErrorsWhileEnable struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsErrorsWhileEnable) IsTranslationKey() {
}
func (k KeysCommandsTtsErrorsWhileEnable) GetPath() string {
	return "commands.tts.errors.while_enable"
}
func (k KeysCommandsTtsErrorsWhileEnable) GetPathSlice() []string {
	return []string{"commands", "tts", "errors", "while_enable"}
}
func (k KeysCommandsTtsErrorsWhileEnable) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsErrorsWhileEnable) SetVars(vars KeysCommandsTtsErrorsWhileEnableVars) twiri18n.TranslationKey[KeysCommandsTtsErrorsWhileEnableVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsErrorsNotConfiguredVars struct {
}
type KeysCommandsTtsErrorsNotConfigured struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsErrorsNotConfigured) IsTranslationKey() {
}
func (k KeysCommandsTtsErrorsNotConfigured) GetPath() string {
	return "commands.tts.errors.not_configured"
}
func (k KeysCommandsTtsErrorsNotConfigured) GetPathSlice() []string {
	return []string{"commands", "tts", "errors", "not_configured"}
}
func (k KeysCommandsTtsErrorsNotConfigured) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsErrorsNotConfigured) SetVars(vars KeysCommandsTtsErrorsNotConfiguredVars) twiri18n.TranslationKey[KeysCommandsTtsErrorsNotConfiguredVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsErrorsSendingToTtsVars struct {
}
type KeysCommandsTtsErrorsSendingToTts struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsErrorsSendingToTts) IsTranslationKey() {
}
func (k KeysCommandsTtsErrorsSendingToTts) GetPath() string {
	return "commands.tts.errors.sending_to_tts"
}
func (k KeysCommandsTtsErrorsSendingToTts) GetPathSlice() []string {
	return []string{"commands", "tts", "errors", "sending_to_tts"}
}
func (k KeysCommandsTtsErrorsSendingToTts) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsErrorsSendingToTts) SetVars(vars KeysCommandsTtsErrorsSendingToTtsVars) twiri18n.TranslationKey[KeysCommandsTtsErrorsSendingToTtsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsErrors struct {
	WhileGettingVoices	KeysCommandsTtsErrorsWhileGettingVoices
	WhileDisable		KeysCommandsTtsErrorsWhileDisable
	WhileEnable		KeysCommandsTtsErrorsWhileEnable
	NotConfigured		KeysCommandsTtsErrorsNotConfigured
	SendingToTts		KeysCommandsTtsErrorsSendingToTts
}
type KeysCommandsTtsInfoChangePitchVars struct {
	NewPitch any
}
type KeysCommandsTtsInfoChangePitch struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsInfoChangePitch) IsTranslationKey() {
}
func (k KeysCommandsTtsInfoChangePitch) GetPath() string {
	return "commands.tts.info.change_pitch"
}
func (k KeysCommandsTtsInfoChangePitch) GetPathSlice() []string {
	return []string{"commands", "tts", "info", "change_pitch"}
}
func (k KeysCommandsTtsInfoChangePitch) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsInfoChangePitch) SetVars(vars KeysCommandsTtsInfoChangePitchVars) twiri18n.TranslationKey[KeysCommandsTtsInfoChangePitchVars] {
	k.Vars = twiri18n.Vars{"newPitch": vars.NewPitch}
	return k
}

type KeysCommandsTtsInfoRateVars struct {
	GlobalRate	any
	UserRate	any
}
type KeysCommandsTtsInfoRate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsInfoRate) IsTranslationKey() {
}
func (k KeysCommandsTtsInfoRate) GetPath() string {
	return "commands.tts.info.rate"
}
func (k KeysCommandsTtsInfoRate) GetPathSlice() []string {
	return []string{"commands", "tts", "info", "rate"}
}
func (k KeysCommandsTtsInfoRate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsInfoRate) SetVars(vars KeysCommandsTtsInfoRateVars) twiri18n.TranslationKey[KeysCommandsTtsInfoRateVars] {
	k.Vars = twiri18n.Vars{"globalRate": vars.GlobalRate, "userRate": vars.UserRate}
	return k
}

type KeysCommandsTtsInfoChangeRateVars struct {
	NewRate any
}
type KeysCommandsTtsInfoChangeRate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsInfoChangeRate) IsTranslationKey() {
}
func (k KeysCommandsTtsInfoChangeRate) GetPath() string {
	return "commands.tts.info.change_rate"
}
func (k KeysCommandsTtsInfoChangeRate) GetPathSlice() []string {
	return []string{"commands", "tts", "info", "change_rate"}
}
func (k KeysCommandsTtsInfoChangeRate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsInfoChangeRate) SetVars(vars KeysCommandsTtsInfoChangeRateVars) twiri18n.TranslationKey[KeysCommandsTtsInfoChangeRateVars] {
	k.Vars = twiri18n.Vars{"newRate": vars.NewRate}
	return k
}

type KeysCommandsTtsInfoVoiceVars struct {
	GlobalVoice	any
	UserVoice	any
}
type KeysCommandsTtsInfoVoice struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsInfoVoice) IsTranslationKey() {
}
func (k KeysCommandsTtsInfoVoice) GetPath() string {
	return "commands.tts.info.voice"
}
func (k KeysCommandsTtsInfoVoice) GetPathSlice() []string {
	return []string{"commands", "tts", "info", "voice"}
}
func (k KeysCommandsTtsInfoVoice) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsInfoVoice) SetVars(vars KeysCommandsTtsInfoVoiceVars) twiri18n.TranslationKey[KeysCommandsTtsInfoVoiceVars] {
	k.Vars = twiri18n.Vars{"globalVoice": vars.GlobalVoice, "userVoice": vars.UserVoice}
	return k
}

type KeysCommandsTtsInfoChangeVoiceVars struct {
	NewVoice any
}
type KeysCommandsTtsInfoChangeVoice struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsInfoChangeVoice) IsTranslationKey() {
}
func (k KeysCommandsTtsInfoChangeVoice) GetPath() string {
	return "commands.tts.info.change_voice"
}
func (k KeysCommandsTtsInfoChangeVoice) GetPathSlice() []string {
	return []string{"commands", "tts", "info", "change_voice"}
}
func (k KeysCommandsTtsInfoChangeVoice) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsInfoChangeVoice) SetVars(vars KeysCommandsTtsInfoChangeVoiceVars) twiri18n.TranslationKey[KeysCommandsTtsInfoChangeVoiceVars] {
	k.Vars = twiri18n.Vars{"newVoice": vars.NewVoice}
	return k
}

type KeysCommandsTtsInfoChangeVolumeVars struct {
	UserVolume any
}
type KeysCommandsTtsInfoChangeVolume struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsInfoChangeVolume) IsTranslationKey() {
}
func (k KeysCommandsTtsInfoChangeVolume) GetPath() string {
	return "commands.tts.info.change_volume"
}
func (k KeysCommandsTtsInfoChangeVolume) GetPathSlice() []string {
	return []string{"commands", "tts", "info", "change_volume"}
}
func (k KeysCommandsTtsInfoChangeVolume) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsInfoChangeVolume) SetVars(vars KeysCommandsTtsInfoChangeVolumeVars) twiri18n.TranslationKey[KeysCommandsTtsInfoChangeVolumeVars] {
	k.Vars = twiri18n.Vars{"userVolume": vars.UserVolume}
	return k
}

type KeysCommandsTtsInfoVoiceDisallowedVars struct {
	VoiceName any
}
type KeysCommandsTtsInfoVoiceDisallowed struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsInfoVoiceDisallowed) IsTranslationKey() {
}
func (k KeysCommandsTtsInfoVoiceDisallowed) GetPath() string {
	return "commands.tts.info.voice_disallowed"
}
func (k KeysCommandsTtsInfoVoiceDisallowed) GetPathSlice() []string {
	return []string{"commands", "tts", "info", "voice_disallowed"}
}
func (k KeysCommandsTtsInfoVoiceDisallowed) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsInfoVoiceDisallowed) SetVars(vars KeysCommandsTtsInfoVoiceDisallowedVars) twiri18n.TranslationKey[KeysCommandsTtsInfoVoiceDisallowedVars] {
	k.Vars = twiri18n.Vars{"voiceName": vars.VoiceName}
	return k
}

type KeysCommandsTtsInfoDisabledVars struct {
}
type KeysCommandsTtsInfoDisabled struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsInfoDisabled) IsTranslationKey() {
}
func (k KeysCommandsTtsInfoDisabled) GetPath() string {
	return "commands.tts.info.disabled"
}
func (k KeysCommandsTtsInfoDisabled) GetPathSlice() []string {
	return []string{"commands", "tts", "info", "disabled"}
}
func (k KeysCommandsTtsInfoDisabled) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsInfoDisabled) SetVars(vars KeysCommandsTtsInfoDisabledVars) twiri18n.TranslationKey[KeysCommandsTtsInfoDisabledVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsInfoEnabledVars struct {
}
type KeysCommandsTtsInfoEnabled struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsInfoEnabled) IsTranslationKey() {
}
func (k KeysCommandsTtsInfoEnabled) GetPath() string {
	return "commands.tts.info.enabled"
}
func (k KeysCommandsTtsInfoEnabled) GetPathSlice() []string {
	return []string{"commands", "tts", "info", "enabled"}
}
func (k KeysCommandsTtsInfoEnabled) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsInfoEnabled) SetVars(vars KeysCommandsTtsInfoEnabledVars) twiri18n.TranslationKey[KeysCommandsTtsInfoEnabledVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsInfoPitchVars struct {
	GlobalPitch	any
	UserPitch	any
}
type KeysCommandsTtsInfoPitch struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsInfoPitch) IsTranslationKey() {
}
func (k KeysCommandsTtsInfoPitch) GetPath() string {
	return "commands.tts.info.pitch"
}
func (k KeysCommandsTtsInfoPitch) GetPathSlice() []string {
	return []string{"commands", "tts", "info", "pitch"}
}
func (k KeysCommandsTtsInfoPitch) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsInfoPitch) SetVars(vars KeysCommandsTtsInfoPitchVars) twiri18n.TranslationKey[KeysCommandsTtsInfoPitchVars] {
	k.Vars = twiri18n.Vars{"globalPitch": vars.GlobalPitch, "userPitch": vars.UserPitch}
	return k
}

type KeysCommandsTtsInfoNoVoicesVars struct {
}
type KeysCommandsTtsInfoNoVoices struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsInfoNoVoices) IsTranslationKey() {
}
func (k KeysCommandsTtsInfoNoVoices) GetPath() string {
	return "commands.tts.info.no_voices"
}
func (k KeysCommandsTtsInfoNoVoices) GetPathSlice() []string {
	return []string{"commands", "tts", "info", "no_voices"}
}
func (k KeysCommandsTtsInfoNoVoices) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsInfoNoVoices) SetVars(vars KeysCommandsTtsInfoNoVoicesVars) twiri18n.TranslationKey[KeysCommandsTtsInfoNoVoicesVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsTtsInfoCurrentVolumeVars struct {
	TtsVolume any
}
type KeysCommandsTtsInfoCurrentVolume struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsTtsInfoCurrentVolume) IsTranslationKey() {
}
func (k KeysCommandsTtsInfoCurrentVolume) GetPath() string {
	return "commands.tts.info.current_volume"
}
func (k KeysCommandsTtsInfoCurrentVolume) GetPathSlice() []string {
	return []string{"commands", "tts", "info", "current_volume"}
}
func (k KeysCommandsTtsInfoCurrentVolume) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsTtsInfoCurrentVolume) SetVars(vars KeysCommandsTtsInfoCurrentVolumeVars) twiri18n.TranslationKey[KeysCommandsTtsInfoCurrentVolumeVars] {
	k.Vars = twiri18n.Vars{"ttsVolume": vars.TtsVolume}
	return k
}

type KeysCommandsTtsInfo struct {
	ChangePitch	KeysCommandsTtsInfoChangePitch
	Rate		KeysCommandsTtsInfoRate
	ChangeRate	KeysCommandsTtsInfoChangeRate
	Voice		KeysCommandsTtsInfoVoice
	ChangeVoice	KeysCommandsTtsInfoChangeVoice
	ChangeVolume	KeysCommandsTtsInfoChangeVolume
	VoiceDisallowed	KeysCommandsTtsInfoVoiceDisallowed
	Disabled	KeysCommandsTtsInfoDisabled
	Enabled		KeysCommandsTtsInfoEnabled
	Pitch		KeysCommandsTtsInfoPitch
	NoVoices	KeysCommandsTtsInfoNoVoices
	CurrentVolume	KeysCommandsTtsInfoCurrentVolume
}
type KeysCommandsTts struct {
	Errors	KeysCommandsTtsErrors
	Info	KeysCommandsTtsInfo
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

type KeysCommandsSubageResponses struct {
	StreakInfo		KeysCommandsSubageResponsesStreakInfo
	TimeRemaining		KeysCommandsSubageResponsesTimeRemaining
	NotSubscriber		KeysCommandsSubageResponsesNotSubscriber
	NotSubscriberButWas	KeysCommandsSubageResponsesNotSubscriberButWas
	SubscriptionInfo	KeysCommandsSubageResponsesSubscriptionInfo
}
type KeysCommandsSubage struct {
	Errors		KeysCommandsSubageErrors
	Responses	KeysCommandsSubageResponses
}
type KeysCommandsChatWallHintsBanPhraseArgNameVars struct {
}
type KeysCommandsChatWallHintsBanPhraseArgName struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallHintsBanPhraseArgName) IsTranslationKey() {
}
func (k KeysCommandsChatWallHintsBanPhraseArgName) GetPath() string {
	return "commands.chat_wall.hints.banPhraseArgName"
}
func (k KeysCommandsChatWallHintsBanPhraseArgName) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "hints", "banPhraseArgName"}
}
func (k KeysCommandsChatWallHintsBanPhraseArgName) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallHintsBanPhraseArgName) SetVars(vars KeysCommandsChatWallHintsBanPhraseArgNameVars) twiri18n.TranslationKey[KeysCommandsChatWallHintsBanPhraseArgNameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallHintsDeletePhraseArgNameVars struct {
}
type KeysCommandsChatWallHintsDeletePhraseArgName struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallHintsDeletePhraseArgName) IsTranslationKey() {
}
func (k KeysCommandsChatWallHintsDeletePhraseArgName) GetPath() string {
	return "commands.chat_wall.hints.deletePhraseArgName"
}
func (k KeysCommandsChatWallHintsDeletePhraseArgName) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "hints", "deletePhraseArgName"}
}
func (k KeysCommandsChatWallHintsDeletePhraseArgName) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallHintsDeletePhraseArgName) SetVars(vars KeysCommandsChatWallHintsDeletePhraseArgNameVars) twiri18n.TranslationKey[KeysCommandsChatWallHintsDeletePhraseArgNameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallHintsTimeoutDurationArgNameVars struct {
}
type KeysCommandsChatWallHintsTimeoutDurationArgName struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallHintsTimeoutDurationArgName) IsTranslationKey() {
}
func (k KeysCommandsChatWallHintsTimeoutDurationArgName) GetPath() string {
	return "commands.chat_wall.hints.timeoutDurationArgName"
}
func (k KeysCommandsChatWallHintsTimeoutDurationArgName) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "hints", "timeoutDurationArgName"}
}
func (k KeysCommandsChatWallHintsTimeoutDurationArgName) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallHintsTimeoutDurationArgName) SetVars(vars KeysCommandsChatWallHintsTimeoutDurationArgNameVars) twiri18n.TranslationKey[KeysCommandsChatWallHintsTimeoutDurationArgNameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallHintsTimeoutPhraseArgNameVars struct {
}
type KeysCommandsChatWallHintsTimeoutPhraseArgName struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallHintsTimeoutPhraseArgName) IsTranslationKey() {
}
func (k KeysCommandsChatWallHintsTimeoutPhraseArgName) GetPath() string {
	return "commands.chat_wall.hints.timeoutPhraseArgName"
}
func (k KeysCommandsChatWallHintsTimeoutPhraseArgName) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "hints", "timeoutPhraseArgName"}
}
func (k KeysCommandsChatWallHintsTimeoutPhraseArgName) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallHintsTimeoutPhraseArgName) SetVars(vars KeysCommandsChatWallHintsTimeoutPhraseArgNameVars) twiri18n.TranslationKey[KeysCommandsChatWallHintsTimeoutPhraseArgNameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallHints struct {
	BanPhraseArgName	KeysCommandsChatWallHintsBanPhraseArgName
	DeletePhraseArgName	KeysCommandsChatWallHintsDeletePhraseArgName
	TimeoutDurationArgName	KeysCommandsChatWallHintsTimeoutDurationArgName
	TimeoutPhraseArgName	KeysCommandsChatWallHintsTimeoutPhraseArgName
}
type KeysCommandsChatWallErrorsLongDurationTimeoutVars struct {
}
type KeysCommandsChatWallErrorsLongDurationTimeout struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallErrorsLongDurationTimeout) IsTranslationKey() {
}
func (k KeysCommandsChatWallErrorsLongDurationTimeout) GetPath() string {
	return "commands.chat_wall.errors.long_duration_timeout"
}
func (k KeysCommandsChatWallErrorsLongDurationTimeout) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "errors", "long_duration_timeout"}
}
func (k KeysCommandsChatWallErrorsLongDurationTimeout) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallErrorsLongDurationTimeout) SetVars(vars KeysCommandsChatWallErrorsLongDurationTimeoutVars) twiri18n.TranslationKey[KeysCommandsChatWallErrorsLongDurationTimeoutVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallErrorsDurationCannotParseVars struct {
}
type KeysCommandsChatWallErrorsDurationCannotParse struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallErrorsDurationCannotParse) IsTranslationKey() {
}
func (k KeysCommandsChatWallErrorsDurationCannotParse) GetPath() string {
	return "commands.chat_wall.errors.duration_cannot_parse"
}
func (k KeysCommandsChatWallErrorsDurationCannotParse) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "errors", "duration_cannot_parse"}
}
func (k KeysCommandsChatWallErrorsDurationCannotParse) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallErrorsDurationCannotParse) SetVars(vars KeysCommandsChatWallErrorsDurationCannotParseVars) twiri18n.TranslationKey[KeysCommandsChatWallErrorsDurationCannotParseVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallErrorsChatWallNotFoundVars struct {
	ErrorPhrase any
}
type KeysCommandsChatWallErrorsChatWallNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallErrorsChatWallNotFound) IsTranslationKey() {
}
func (k KeysCommandsChatWallErrorsChatWallNotFound) GetPath() string {
	return "commands.chat_wall.errors.chat_wall_not_found"
}
func (k KeysCommandsChatWallErrorsChatWallNotFound) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "errors", "chat_wall_not_found"}
}
func (k KeysCommandsChatWallErrorsChatWallNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallErrorsChatWallNotFound) SetVars(vars KeysCommandsChatWallErrorsChatWallNotFoundVars) twiri18n.TranslationKey[KeysCommandsChatWallErrorsChatWallNotFoundVars] {
	k.Vars = twiri18n.Vars{"errorPhrase": vars.ErrorPhrase}
	return k
}

type KeysCommandsChatWallErrorsInvalidDurationVars struct {
}
type KeysCommandsChatWallErrorsInvalidDuration struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallErrorsInvalidDuration) IsTranslationKey() {
}
func (k KeysCommandsChatWallErrorsInvalidDuration) GetPath() string {
	return "commands.chat_wall.errors.invalid_duration"
}
func (k KeysCommandsChatWallErrorsInvalidDuration) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "errors", "invalid_duration"}
}
func (k KeysCommandsChatWallErrorsInvalidDuration) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallErrorsInvalidDuration) SetVars(vars KeysCommandsChatWallErrorsInvalidDurationVars) twiri18n.TranslationKey[KeysCommandsChatWallErrorsInvalidDurationVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsChatWallErrors struct {
	LongDurationTimeout	KeysCommandsChatWallErrorsLongDurationTimeout
	DurationCannotParse	KeysCommandsChatWallErrorsDurationCannotParse
	ChatWallNotFound	KeysCommandsChatWallErrorsChatWallNotFound
	InvalidDuration		KeysCommandsChatWallErrorsInvalidDuration
}
type KeysCommandsChatWallStartChatWallStartVars struct {
	ChatWallPhrase any
}
type KeysCommandsChatWallStartChatWallStart struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallStartChatWallStart) IsTranslationKey() {
}
func (k KeysCommandsChatWallStartChatWallStart) GetPath() string {
	return "commands.chat_wall.start.chat_wall_start"
}
func (k KeysCommandsChatWallStartChatWallStart) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "start", "chat_wall_start"}
}
func (k KeysCommandsChatWallStartChatWallStart) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallStartChatWallStart) SetVars(vars KeysCommandsChatWallStartChatWallStartVars) twiri18n.TranslationKey[KeysCommandsChatWallStartChatWallStartVars] {
	k.Vars = twiri18n.Vars{"chatWallPhrase": vars.ChatWallPhrase}
	return k
}

type KeysCommandsChatWallStart struct {
	ChatWallStart KeysCommandsChatWallStartChatWallStart
}
type KeysCommandsChatWallStopChatWalStopVars struct {
	ChatWallPhrase any
}
type KeysCommandsChatWallStopChatWalStop struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsChatWallStopChatWalStop) IsTranslationKey() {
}
func (k KeysCommandsChatWallStopChatWalStop) GetPath() string {
	return "commands.chat_wall.stop.chat_wal_stop"
}
func (k KeysCommandsChatWallStopChatWalStop) GetPathSlice() []string {
	return []string{"commands", "chat_wall", "stop", "chat_wal_stop"}
}
func (k KeysCommandsChatWallStopChatWalStop) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsChatWallStopChatWalStop) SetVars(vars KeysCommandsChatWallStopChatWalStopVars) twiri18n.TranslationKey[KeysCommandsChatWallStopChatWalStopVars] {
	k.Vars = twiri18n.Vars{"chatWallPhrase": vars.ChatWallPhrase}
	return k
}

type KeysCommandsChatWallStop struct {
	ChatWalStop KeysCommandsChatWallStopChatWalStop
}
type KeysCommandsChatWall struct {
	Hints	KeysCommandsChatWallHints
	Errors	KeysCommandsChatWallErrors
	Start	KeysCommandsChatWallStart
	Stop	KeysCommandsChatWallStop
}
type KeysCommandsPredictionsHintsStartPredictionDurationVars struct {
}
type KeysCommandsPredictionsHintsStartPredictionDuration struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsHintsStartPredictionDuration) IsTranslationKey() {
}
func (k KeysCommandsPredictionsHintsStartPredictionDuration) GetPath() string {
	return "commands.predictions.hints.startPredictionDuration"
}
func (k KeysCommandsPredictionsHintsStartPredictionDuration) GetPathSlice() []string {
	return []string{"commands", "predictions", "hints", "startPredictionDuration"}
}
func (k KeysCommandsPredictionsHintsStartPredictionDuration) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsHintsStartPredictionDuration) SetVars(vars KeysCommandsPredictionsHintsStartPredictionDurationVars) twiri18n.TranslationKey[KeysCommandsPredictionsHintsStartPredictionDurationVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsHintsStartPredictionArgTitleVars struct {
}
type KeysCommandsPredictionsHintsStartPredictionArgTitle struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsHintsStartPredictionArgTitle) IsTranslationKey() {
}
func (k KeysCommandsPredictionsHintsStartPredictionArgTitle) GetPath() string {
	return "commands.predictions.hints.startPredictionArgTitle"
}
func (k KeysCommandsPredictionsHintsStartPredictionArgTitle) GetPathSlice() []string {
	return []string{"commands", "predictions", "hints", "startPredictionArgTitle"}
}
func (k KeysCommandsPredictionsHintsStartPredictionArgTitle) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsHintsStartPredictionArgTitle) SetVars(vars KeysCommandsPredictionsHintsStartPredictionArgTitleVars) twiri18n.TranslationKey[KeysCommandsPredictionsHintsStartPredictionArgTitleVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsHintsStartPredictionArgVariantsVars struct {
}
type KeysCommandsPredictionsHintsStartPredictionArgVariants struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsHintsStartPredictionArgVariants) IsTranslationKey() {
}
func (k KeysCommandsPredictionsHintsStartPredictionArgVariants) GetPath() string {
	return "commands.predictions.hints.startPredictionArgVariants"
}
func (k KeysCommandsPredictionsHintsStartPredictionArgVariants) GetPathSlice() []string {
	return []string{"commands", "predictions", "hints", "startPredictionArgVariants"}
}
func (k KeysCommandsPredictionsHintsStartPredictionArgVariants) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsHintsStartPredictionArgVariants) SetVars(vars KeysCommandsPredictionsHintsStartPredictionArgVariantsVars) twiri18n.TranslationKey[KeysCommandsPredictionsHintsStartPredictionArgVariantsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsHintsPredictionResolveOutcomeNumVars struct {
}
type KeysCommandsPredictionsHintsPredictionResolveOutcomeNum struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsHintsPredictionResolveOutcomeNum) IsTranslationKey() {
}
func (k KeysCommandsPredictionsHintsPredictionResolveOutcomeNum) GetPath() string {
	return "commands.predictions.hints.predictionResolveOutcomeNum"
}
func (k KeysCommandsPredictionsHintsPredictionResolveOutcomeNum) GetPathSlice() []string {
	return []string{"commands", "predictions", "hints", "predictionResolveOutcomeNum"}
}
func (k KeysCommandsPredictionsHintsPredictionResolveOutcomeNum) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsHintsPredictionResolveOutcomeNum) SetVars(vars KeysCommandsPredictionsHintsPredictionResolveOutcomeNumVars) twiri18n.TranslationKey[KeysCommandsPredictionsHintsPredictionResolveOutcomeNumVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsHints struct {
	StartPredictionDuration		KeysCommandsPredictionsHintsStartPredictionDuration
	StartPredictionArgTitle		KeysCommandsPredictionsHintsStartPredictionArgTitle
	StartPredictionArgVariants	KeysCommandsPredictionsHintsStartPredictionArgVariants
	PredictionResolveOutcomeNum	KeysCommandsPredictionsHintsPredictionResolveOutcomeNum
}
type KeysCommandsPredictionsErrorsCannotGetCurrentVarVars struct {
	Reason any
}
type KeysCommandsPredictionsErrorsCannotGetCurrentVar struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsErrorsCannotGetCurrentVar) IsTranslationKey() {
}
func (k KeysCommandsPredictionsErrorsCannotGetCurrentVar) GetPath() string {
	return "commands.predictions.errors.cannot_get_current_var"
}
func (k KeysCommandsPredictionsErrorsCannotGetCurrentVar) GetPathSlice() []string {
	return []string{"commands", "predictions", "errors", "cannot_get_current_var"}
}
func (k KeysCommandsPredictionsErrorsCannotGetCurrentVar) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsErrorsCannotGetCurrentVar) SetVars(vars KeysCommandsPredictionsErrorsCannotGetCurrentVarVars) twiri18n.TranslationKey[KeysCommandsPredictionsErrorsCannotGetCurrentVarVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsPredictionsErrorsCannotCancelVars struct {
}
type KeysCommandsPredictionsErrorsCannotCancel struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsErrorsCannotCancel) IsTranslationKey() {
}
func (k KeysCommandsPredictionsErrorsCannotCancel) GetPath() string {
	return "commands.predictions.errors.cannot_cancel"
}
func (k KeysCommandsPredictionsErrorsCannotCancel) GetPathSlice() []string {
	return []string{"commands", "predictions", "errors", "cannot_cancel"}
}
func (k KeysCommandsPredictionsErrorsCannotCancel) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsErrorsCannotCancel) SetVars(vars KeysCommandsPredictionsErrorsCannotCancelVars) twiri18n.TranslationKey[KeysCommandsPredictionsErrorsCannotCancelVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsErrorsCannotCancelVarVars struct {
	Reason any
}
type KeysCommandsPredictionsErrorsCannotCancelVar struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsErrorsCannotCancelVar) IsTranslationKey() {
}
func (k KeysCommandsPredictionsErrorsCannotCancelVar) GetPath() string {
	return "commands.predictions.errors.cannot_cancel_var"
}
func (k KeysCommandsPredictionsErrorsCannotCancelVar) GetPathSlice() []string {
	return []string{"commands", "predictions", "errors", "cannot_cancel_var"}
}
func (k KeysCommandsPredictionsErrorsCannotCancelVar) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsErrorsCannotCancelVar) SetVars(vars KeysCommandsPredictionsErrorsCannotCancelVarVars) twiri18n.TranslationKey[KeysCommandsPredictionsErrorsCannotCancelVarVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsPredictionsErrorsNoVariantVars struct {
}
type KeysCommandsPredictionsErrorsNoVariant struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsErrorsNoVariant) IsTranslationKey() {
}
func (k KeysCommandsPredictionsErrorsNoVariant) GetPath() string {
	return "commands.predictions.errors.no_variant"
}
func (k KeysCommandsPredictionsErrorsNoVariant) GetPathSlice() []string {
	return []string{"commands", "predictions", "errors", "no_variant"}
}
func (k KeysCommandsPredictionsErrorsNoVariant) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsErrorsNoVariant) SetVars(vars KeysCommandsPredictionsErrorsNoVariantVars) twiri18n.TranslationKey[KeysCommandsPredictionsErrorsNoVariantVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsErrorsCannotCreateVars struct {
}
type KeysCommandsPredictionsErrorsCannotCreate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsErrorsCannotCreate) IsTranslationKey() {
}
func (k KeysCommandsPredictionsErrorsCannotCreate) GetPath() string {
	return "commands.predictions.errors.cannot_create"
}
func (k KeysCommandsPredictionsErrorsCannotCreate) GetPathSlice() []string {
	return []string{"commands", "predictions", "errors", "cannot_create"}
}
func (k KeysCommandsPredictionsErrorsCannotCreate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsErrorsCannotCreate) SetVars(vars KeysCommandsPredictionsErrorsCannotCreateVars) twiri18n.TranslationKey[KeysCommandsPredictionsErrorsCannotCreateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsErrorsCannotCreateVarVars struct {
	Reason any
}
type KeysCommandsPredictionsErrorsCannotCreateVar struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsErrorsCannotCreateVar) IsTranslationKey() {
}
func (k KeysCommandsPredictionsErrorsCannotCreateVar) GetPath() string {
	return "commands.predictions.errors.cannot_create_var"
}
func (k KeysCommandsPredictionsErrorsCannotCreateVar) GetPathSlice() []string {
	return []string{"commands", "predictions", "errors", "cannot_create_var"}
}
func (k KeysCommandsPredictionsErrorsCannotCreateVar) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsErrorsCannotCreateVar) SetVars(vars KeysCommandsPredictionsErrorsCannotCreateVarVars) twiri18n.TranslationKey[KeysCommandsPredictionsErrorsCannotCreateVarVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsPredictionsErrorsCannotGetCurrentVars struct {
}
type KeysCommandsPredictionsErrorsCannotGetCurrent struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsErrorsCannotGetCurrent) IsTranslationKey() {
}
func (k KeysCommandsPredictionsErrorsCannotGetCurrent) GetPath() string {
	return "commands.predictions.errors.cannot_get_current"
}
func (k KeysCommandsPredictionsErrorsCannotGetCurrent) GetPathSlice() []string {
	return []string{"commands", "predictions", "errors", "cannot_get_current"}
}
func (k KeysCommandsPredictionsErrorsCannotGetCurrent) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsErrorsCannotGetCurrent) SetVars(vars KeysCommandsPredictionsErrorsCannotGetCurrentVars) twiri18n.TranslationKey[KeysCommandsPredictionsErrorsCannotGetCurrentVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsErrors struct {
	CannotGetCurrentVar	KeysCommandsPredictionsErrorsCannotGetCurrentVar
	CannotCancel		KeysCommandsPredictionsErrorsCannotCancel
	CannotCancelVar		KeysCommandsPredictionsErrorsCannotCancelVar
	NoVariant		KeysCommandsPredictionsErrorsNoVariant
	CannotCreate		KeysCommandsPredictionsErrorsCannotCreate
	CannotCreateVar		KeysCommandsPredictionsErrorsCannotCreateVar
	CannotGetCurrent	KeysCommandsPredictionsErrorsCannotGetCurrent
}
type KeysCommandsPredictionsInfoStartedVars struct {
}
type KeysCommandsPredictionsInfoStarted struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsInfoStarted) IsTranslationKey() {
}
func (k KeysCommandsPredictionsInfoStarted) GetPath() string {
	return "commands.predictions.info.started"
}
func (k KeysCommandsPredictionsInfoStarted) GetPathSlice() []string {
	return []string{"commands", "predictions", "info", "started"}
}
func (k KeysCommandsPredictionsInfoStarted) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsInfoStarted) SetVars(vars KeysCommandsPredictionsInfoStartedVars) twiri18n.TranslationKey[KeysCommandsPredictionsInfoStartedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsInfoNoRunedVars struct {
}
type KeysCommandsPredictionsInfoNoRuned struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsInfoNoRuned) IsTranslationKey() {
}
func (k KeysCommandsPredictionsInfoNoRuned) GetPath() string {
	return "commands.predictions.info.no_runed"
}
func (k KeysCommandsPredictionsInfoNoRuned) GetPathSlice() []string {
	return []string{"commands", "predictions", "info", "no_runed"}
}
func (k KeysCommandsPredictionsInfoNoRuned) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsInfoNoRuned) SetVars(vars KeysCommandsPredictionsInfoNoRunedVars) twiri18n.TranslationKey[KeysCommandsPredictionsInfoNoRunedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsInfoCancelVars struct {
}
type KeysCommandsPredictionsInfoCancel struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsInfoCancel) IsTranslationKey() {
}
func (k KeysCommandsPredictionsInfoCancel) GetPath() string {
	return "commands.predictions.info.cancel"
}
func (k KeysCommandsPredictionsInfoCancel) GetPathSlice() []string {
	return []string{"commands", "predictions", "info", "cancel"}
}
func (k KeysCommandsPredictionsInfoCancel) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsInfoCancel) SetVars(vars KeysCommandsPredictionsInfoCancelVars) twiri18n.TranslationKey[KeysCommandsPredictionsInfoCancelVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsInfoLockedVars struct {
}
type KeysCommandsPredictionsInfoLocked struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsInfoLocked) IsTranslationKey() {
}
func (k KeysCommandsPredictionsInfoLocked) GetPath() string {
	return "commands.predictions.info.locked"
}
func (k KeysCommandsPredictionsInfoLocked) GetPathSlice() []string {
	return []string{"commands", "predictions", "info", "locked"}
}
func (k KeysCommandsPredictionsInfoLocked) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsInfoLocked) SetVars(vars KeysCommandsPredictionsInfoLockedVars) twiri18n.TranslationKey[KeysCommandsPredictionsInfoLockedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsInfoResolvedVars struct {
}
type KeysCommandsPredictionsInfoResolved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPredictionsInfoResolved) IsTranslationKey() {
}
func (k KeysCommandsPredictionsInfoResolved) GetPath() string {
	return "commands.predictions.info.resolved"
}
func (k KeysCommandsPredictionsInfoResolved) GetPathSlice() []string {
	return []string{"commands", "predictions", "info", "resolved"}
}
func (k KeysCommandsPredictionsInfoResolved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPredictionsInfoResolved) SetVars(vars KeysCommandsPredictionsInfoResolvedVars) twiri18n.TranslationKey[KeysCommandsPredictionsInfoResolvedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPredictionsInfo struct {
	Started		KeysCommandsPredictionsInfoStarted
	NoRuned		KeysCommandsPredictionsInfoNoRuned
	Cancel		KeysCommandsPredictionsInfoCancel
	Locked		KeysCommandsPredictionsInfoLocked
	Resolved	KeysCommandsPredictionsInfoResolved
}
type KeysCommandsPredictions struct {
	Hints	KeysCommandsPredictionsHints
	Errors	KeysCommandsPredictionsErrors
	Info	KeysCommandsPredictionsInfo
}
type KeysCommandsPrefixErrorsTooLongVars struct {
}
type KeysCommandsPrefixErrorsTooLong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixErrorsTooLong) IsTranslationKey() {
}
func (k KeysCommandsPrefixErrorsTooLong) GetPath() string {
	return "commands.prefix.errors.too_long"
}
func (k KeysCommandsPrefixErrorsTooLong) GetPathSlice() []string {
	return []string{"commands", "prefix", "errors", "too_long"}
}
func (k KeysCommandsPrefixErrorsTooLong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixErrorsTooLong) SetVars(vars KeysCommandsPrefixErrorsTooLongVars) twiri18n.TranslationKey[KeysCommandsPrefixErrorsTooLongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixErrorsCannotGetCurrentVars struct {
}
type KeysCommandsPrefixErrorsCannotGetCurrent struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixErrorsCannotGetCurrent) IsTranslationKey() {
}
func (k KeysCommandsPrefixErrorsCannotGetCurrent) GetPath() string {
	return "commands.prefix.errors.cannot_get_current"
}
func (k KeysCommandsPrefixErrorsCannotGetCurrent) GetPathSlice() []string {
	return []string{"commands", "prefix", "errors", "cannot_get_current"}
}
func (k KeysCommandsPrefixErrorsCannotGetCurrent) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixErrorsCannotGetCurrent) SetVars(vars KeysCommandsPrefixErrorsCannotGetCurrentVars) twiri18n.TranslationKey[KeysCommandsPrefixErrorsCannotGetCurrentVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixErrorsCannotCreateVars struct {
}
type KeysCommandsPrefixErrorsCannotCreate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixErrorsCannotCreate) IsTranslationKey() {
}
func (k KeysCommandsPrefixErrorsCannotCreate) GetPath() string {
	return "commands.prefix.errors.cannot_create"
}
func (k KeysCommandsPrefixErrorsCannotCreate) GetPathSlice() []string {
	return []string{"commands", "prefix", "errors", "cannot_create"}
}
func (k KeysCommandsPrefixErrorsCannotCreate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixErrorsCannotCreate) SetVars(vars KeysCommandsPrefixErrorsCannotCreateVars) twiri18n.TranslationKey[KeysCommandsPrefixErrorsCannotCreateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixErrorsCannotUpdateVars struct {
}
type KeysCommandsPrefixErrorsCannotUpdate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixErrorsCannotUpdate) IsTranslationKey() {
}
func (k KeysCommandsPrefixErrorsCannotUpdate) GetPath() string {
	return "commands.prefix.errors.cannot_update"
}
func (k KeysCommandsPrefixErrorsCannotUpdate) GetPathSlice() []string {
	return []string{"commands", "prefix", "errors", "cannot_update"}
}
func (k KeysCommandsPrefixErrorsCannotUpdate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixErrorsCannotUpdate) SetVars(vars KeysCommandsPrefixErrorsCannotUpdateVars) twiri18n.TranslationKey[KeysCommandsPrefixErrorsCannotUpdateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixErrorsRequiredVars struct {
}
type KeysCommandsPrefixErrorsRequired struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixErrorsRequired) IsTranslationKey() {
}
func (k KeysCommandsPrefixErrorsRequired) GetPath() string {
	return "commands.prefix.errors.required"
}
func (k KeysCommandsPrefixErrorsRequired) GetPathSlice() []string {
	return []string{"commands", "prefix", "errors", "required"}
}
func (k KeysCommandsPrefixErrorsRequired) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixErrorsRequired) SetVars(vars KeysCommandsPrefixErrorsRequiredVars) twiri18n.TranslationKey[KeysCommandsPrefixErrorsRequiredVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixErrors struct {
	TooLong			KeysCommandsPrefixErrorsTooLong
	CannotGetCurrent	KeysCommandsPrefixErrorsCannotGetCurrent
	CannotCreate		KeysCommandsPrefixErrorsCannotCreate
	CannotUpdate		KeysCommandsPrefixErrorsCannotUpdate
	Required		KeysCommandsPrefixErrorsRequired
}
type KeysCommandsPrefixSuccessUpdatedVars struct {
}
type KeysCommandsPrefixSuccessUpdated struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsPrefixSuccessUpdated) IsTranslationKey() {
}
func (k KeysCommandsPrefixSuccessUpdated) GetPath() string {
	return "commands.prefix.success.updated"
}
func (k KeysCommandsPrefixSuccessUpdated) GetPathSlice() []string {
	return []string{"commands", "prefix", "success", "updated"}
}
func (k KeysCommandsPrefixSuccessUpdated) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsPrefixSuccessUpdated) SetVars(vars KeysCommandsPrefixSuccessUpdatedVars) twiri18n.TranslationKey[KeysCommandsPrefixSuccessUpdatedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsPrefixSuccess struct {
	Updated KeysCommandsPrefixSuccessUpdated
}
type KeysCommandsPrefix struct {
	Errors	KeysCommandsPrefixErrors
	Success	KeysCommandsPrefixSuccess
}
type KeysCommandsManageErrorsCommandWithNameCannotFindVars struct {
}
type KeysCommandsManageErrorsCommandWithNameCannotFind struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsCommandWithNameCannotFind) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsCommandWithNameCannotFind) GetPath() string {
	return "commands.manage.errors.command_with_name_cannot_find"
}
func (k KeysCommandsManageErrorsCommandWithNameCannotFind) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "command_with_name_cannot_find"}
}
func (k KeysCommandsManageErrorsCommandWithNameCannotFind) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsCommandWithNameCannotFind) SetVars(vars KeysCommandsManageErrorsCommandWithNameCannotFindVars) twiri18n.TranslationKey[KeysCommandsManageErrorsCommandWithNameCannotFindVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsCommandNotFoundVars struct {
}
type KeysCommandsManageErrorsCommandNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsCommandNotFound) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsCommandNotFound) GetPath() string {
	return "commands.manage.errors.command_not_found"
}
func (k KeysCommandsManageErrorsCommandNotFound) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "command_not_found"}
}
func (k KeysCommandsManageErrorsCommandNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsCommandNotFound) SetVars(vars KeysCommandsManageErrorsCommandNotFoundVars) twiri18n.TranslationKey[KeysCommandsManageErrorsCommandNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsCommandLargeSizeVars struct {
}
type KeysCommandsManageErrorsCommandLargeSize struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsCommandLargeSize) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsCommandLargeSize) GetPath() string {
	return "commands.manage.errors.command_large_size"
}
func (k KeysCommandsManageErrorsCommandLargeSize) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "command_large_size"}
}
func (k KeysCommandsManageErrorsCommandLargeSize) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsCommandLargeSize) SetVars(vars KeysCommandsManageErrorsCommandLargeSizeVars) twiri18n.TranslationKey[KeysCommandsManageErrorsCommandLargeSizeVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsCommandHaveNoAliasesVars struct {
}
type KeysCommandsManageErrorsCommandHaveNoAliases struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsCommandHaveNoAliases) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsCommandHaveNoAliases) GetPath() string {
	return "commands.manage.errors.command_have_no_aliases"
}
func (k KeysCommandsManageErrorsCommandHaveNoAliases) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "command_have_no_aliases"}
}
func (k KeysCommandsManageErrorsCommandHaveNoAliases) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsCommandHaveNoAliases) SetVars(vars KeysCommandsManageErrorsCommandHaveNoAliasesVars) twiri18n.TranslationKey[KeysCommandsManageErrorsCommandHaveNoAliasesVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsCommandCannotDeleteDefaultVars struct {
}
type KeysCommandsManageErrorsCommandCannotDeleteDefault struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsCommandCannotDeleteDefault) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsCommandCannotDeleteDefault) GetPath() string {
	return "commands.manage.errors.command_cannot_delete_default"
}
func (k KeysCommandsManageErrorsCommandCannotDeleteDefault) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "command_cannot_delete_default"}
}
func (k KeysCommandsManageErrorsCommandCannotDeleteDefault) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsCommandCannotDeleteDefault) SetVars(vars KeysCommandsManageErrorsCommandCannotDeleteDefaultVars) twiri18n.TranslationKey[KeysCommandsManageErrorsCommandCannotDeleteDefaultVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsCommandCannotUpdateResponseVars struct {
}
type KeysCommandsManageErrorsCommandCannotUpdateResponse struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsCommandCannotUpdateResponse) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsCommandCannotUpdateResponse) GetPath() string {
	return "commands.manage.errors.command_cannot_update_response"
}
func (k KeysCommandsManageErrorsCommandCannotUpdateResponse) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "command_cannot_update_response"}
}
func (k KeysCommandsManageErrorsCommandCannotUpdateResponse) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsCommandCannotUpdateResponse) SetVars(vars KeysCommandsManageErrorsCommandCannotUpdateResponseVars) twiri18n.TranslationKey[KeysCommandsManageErrorsCommandCannotUpdateResponseVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsCommandCannotUpdateVars struct {
}
type KeysCommandsManageErrorsCommandCannotUpdate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsCommandCannotUpdate) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsCommandCannotUpdate) GetPath() string {
	return "commands.manage.errors.command_cannot_update"
}
func (k KeysCommandsManageErrorsCommandCannotUpdate) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "command_cannot_update"}
}
func (k KeysCommandsManageErrorsCommandCannotUpdate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsCommandCannotUpdate) SetVars(vars KeysCommandsManageErrorsCommandCannotUpdateVars) twiri18n.TranslationKey[KeysCommandsManageErrorsCommandCannotUpdateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsAliasCannotGetExistedCommandsVars struct {
}
type KeysCommandsManageErrorsAliasCannotGetExistedCommands struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsAliasCannotGetExistedCommands) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsAliasCannotGetExistedCommands) GetPath() string {
	return "commands.manage.errors.alias_cannot_get_existed_commands"
}
func (k KeysCommandsManageErrorsAliasCannotGetExistedCommands) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "alias_cannot_get_existed_commands"}
}
func (k KeysCommandsManageErrorsAliasCannotGetExistedCommands) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsAliasCannotGetExistedCommands) SetVars(vars KeysCommandsManageErrorsAliasCannotGetExistedCommandsVars) twiri18n.TranslationKey[KeysCommandsManageErrorsAliasCannotGetExistedCommandsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsAliasCannotUpdateVars struct {
}
type KeysCommandsManageErrorsAliasCannotUpdate struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsAliasCannotUpdate) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsAliasCannotUpdate) GetPath() string {
	return "commands.manage.errors.alias_cannot_update"
}
func (k KeysCommandsManageErrorsAliasCannotUpdate) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "alias_cannot_update"}
}
func (k KeysCommandsManageErrorsAliasCannotUpdate) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsAliasCannotUpdate) SetVars(vars KeysCommandsManageErrorsAliasCannotUpdateVars) twiri18n.TranslationKey[KeysCommandsManageErrorsAliasCannotUpdateVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsAliasNotCommandVars struct {
}
type KeysCommandsManageErrorsAliasNotCommand struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsAliasNotCommand) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsAliasNotCommand) GetPath() string {
	return "commands.manage.errors.alias_not_command"
}
func (k KeysCommandsManageErrorsAliasNotCommand) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "alias_not_command"}
}
func (k KeysCommandsManageErrorsAliasNotCommand) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsAliasNotCommand) SetVars(vars KeysCommandsManageErrorsAliasNotCommandVars) twiri18n.TranslationKey[KeysCommandsManageErrorsAliasNotCommandVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsCommandWithAliasAlreadyExistsVars struct {
}
type KeysCommandsManageErrorsCommandWithAliasAlreadyExists struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsCommandWithAliasAlreadyExists) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsCommandWithAliasAlreadyExists) GetPath() string {
	return "commands.manage.errors.command_with_alias_already_exists"
}
func (k KeysCommandsManageErrorsCommandWithAliasAlreadyExists) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "command_with_alias_already_exists"}
}
func (k KeysCommandsManageErrorsCommandWithAliasAlreadyExists) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsCommandWithAliasAlreadyExists) SetVars(vars KeysCommandsManageErrorsCommandWithAliasAlreadyExistsVars) twiri18n.TranslationKey[KeysCommandsManageErrorsCommandWithAliasAlreadyExistsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsCommandCannotGetVars struct {
}
type KeysCommandsManageErrorsCommandCannotGet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsCommandCannotGet) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsCommandCannotGet) GetPath() string {
	return "commands.manage.errors.command_cannot_get"
}
func (k KeysCommandsManageErrorsCommandCannotGet) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "command_cannot_get"}
}
func (k KeysCommandsManageErrorsCommandCannotGet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsCommandCannotGet) SetVars(vars KeysCommandsManageErrorsCommandCannotGetVars) twiri18n.TranslationKey[KeysCommandsManageErrorsCommandCannotGetVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsCommandCannotSaveVars struct {
}
type KeysCommandsManageErrorsCommandCannotSave struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsCommandCannotSave) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsCommandCannotSave) GetPath() string {
	return "commands.manage.errors.command_cannot_save"
}
func (k KeysCommandsManageErrorsCommandCannotSave) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "command_cannot_save"}
}
func (k KeysCommandsManageErrorsCommandCannotSave) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsCommandCannotSave) SetVars(vars KeysCommandsManageErrorsCommandCannotSaveVars) twiri18n.TranslationKey[KeysCommandsManageErrorsCommandCannotSaveVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageErrorsAliasAlreadyExistVars struct {
	Alias any
}
type KeysCommandsManageErrorsAliasAlreadyExist struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageErrorsAliasAlreadyExist) IsTranslationKey() {
}
func (k KeysCommandsManageErrorsAliasAlreadyExist) GetPath() string {
	return "commands.manage.errors.alias_already_exist"
}
func (k KeysCommandsManageErrorsAliasAlreadyExist) GetPathSlice() []string {
	return []string{"commands", "manage", "errors", "alias_already_exist"}
}
func (k KeysCommandsManageErrorsAliasAlreadyExist) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageErrorsAliasAlreadyExist) SetVars(vars KeysCommandsManageErrorsAliasAlreadyExistVars) twiri18n.TranslationKey[KeysCommandsManageErrorsAliasAlreadyExistVars] {
	k.Vars = twiri18n.Vars{"alias": vars.Alias}
	return k
}

type KeysCommandsManageErrors struct {
	CommandWithNameCannotFind	KeysCommandsManageErrorsCommandWithNameCannotFind
	CommandNotFound			KeysCommandsManageErrorsCommandNotFound
	CommandLargeSize		KeysCommandsManageErrorsCommandLargeSize
	CommandHaveNoAliases		KeysCommandsManageErrorsCommandHaveNoAliases
	CommandCannotDeleteDefault	KeysCommandsManageErrorsCommandCannotDeleteDefault
	CommandCannotUpdateResponse	KeysCommandsManageErrorsCommandCannotUpdateResponse
	CommandCannotUpdate		KeysCommandsManageErrorsCommandCannotUpdate
	AliasCannotGetExistedCommands	KeysCommandsManageErrorsAliasCannotGetExistedCommands
	AliasCannotUpdate		KeysCommandsManageErrorsAliasCannotUpdate
	AliasNotCommand			KeysCommandsManageErrorsAliasNotCommand
	CommandWithAliasAlreadyExists	KeysCommandsManageErrorsCommandWithAliasAlreadyExists
	CommandCannotGet		KeysCommandsManageErrorsCommandCannotGet
	CommandCannotSave		KeysCommandsManageErrorsCommandCannotSave
	AliasAlreadyExist		KeysCommandsManageErrorsAliasAlreadyExist
}
type KeysCommandsManageAddAliasAddVars struct {
}
type KeysCommandsManageAddAliasAdd struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageAddAliasAdd) IsTranslationKey() {
}
func (k KeysCommandsManageAddAliasAdd) GetPath() string {
	return "commands.manage.add.alias_add"
}
func (k KeysCommandsManageAddAliasAdd) GetPathSlice() []string {
	return []string{"commands", "manage", "add", "alias_add"}
}
func (k KeysCommandsManageAddAliasAdd) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageAddAliasAdd) SetVars(vars KeysCommandsManageAddAliasAddVars) twiri18n.TranslationKey[KeysCommandsManageAddAliasAddVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageAddCommandAddVars struct {
}
type KeysCommandsManageAddCommandAdd struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageAddCommandAdd) IsTranslationKey() {
}
func (k KeysCommandsManageAddCommandAdd) GetPath() string {
	return "commands.manage.add.command_add"
}
func (k KeysCommandsManageAddCommandAdd) GetPathSlice() []string {
	return []string{"commands", "manage", "add", "command_add"}
}
func (k KeysCommandsManageAddCommandAdd) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageAddCommandAdd) SetVars(vars KeysCommandsManageAddCommandAddVars) twiri18n.TranslationKey[KeysCommandsManageAddCommandAddVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageAdd struct {
	AliasAdd	KeysCommandsManageAddAliasAdd
	CommandAdd	KeysCommandsManageAddCommandAdd
}
type KeysCommandsManageEditCommandEditedVars struct {
}
type KeysCommandsManageEditCommandEdited struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageEditCommandEdited) IsTranslationKey() {
}
func (k KeysCommandsManageEditCommandEdited) GetPath() string {
	return "commands.manage.edit.command_edited"
}
func (k KeysCommandsManageEditCommandEdited) GetPathSlice() []string {
	return []string{"commands", "manage", "edit", "command_edited"}
}
func (k KeysCommandsManageEditCommandEdited) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageEditCommandEdited) SetVars(vars KeysCommandsManageEditCommandEditedVars) twiri18n.TranslationKey[KeysCommandsManageEditCommandEditedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageEdit struct {
	CommandEdited KeysCommandsManageEditCommandEdited
}
type KeysCommandsManageRemoveAliasRemovedVars struct {
}
type KeysCommandsManageRemoveAliasRemoved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageRemoveAliasRemoved) IsTranslationKey() {
}
func (k KeysCommandsManageRemoveAliasRemoved) GetPath() string {
	return "commands.manage.remove.alias_removed"
}
func (k KeysCommandsManageRemoveAliasRemoved) GetPathSlice() []string {
	return []string{"commands", "manage", "remove", "alias_removed"}
}
func (k KeysCommandsManageRemoveAliasRemoved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageRemoveAliasRemoved) SetVars(vars KeysCommandsManageRemoveAliasRemovedVars) twiri18n.TranslationKey[KeysCommandsManageRemoveAliasRemovedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageRemoveCommandRemovedVars struct {
}
type KeysCommandsManageRemoveCommandRemoved struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsManageRemoveCommandRemoved) IsTranslationKey() {
}
func (k KeysCommandsManageRemoveCommandRemoved) GetPath() string {
	return "commands.manage.remove.command_removed"
}
func (k KeysCommandsManageRemoveCommandRemoved) GetPathSlice() []string {
	return []string{"commands", "manage", "remove", "command_removed"}
}
func (k KeysCommandsManageRemoveCommandRemoved) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsManageRemoveCommandRemoved) SetVars(vars KeysCommandsManageRemoveCommandRemovedVars) twiri18n.TranslationKey[KeysCommandsManageRemoveCommandRemovedVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsManageRemove struct {
	AliasRemoved	KeysCommandsManageRemoveAliasRemoved
	CommandRemoved	KeysCommandsManageRemoveCommandRemoved
}
type KeysCommandsManage struct {
	Errors	KeysCommandsManageErrors
	Add	KeysCommandsManageAdd
	Edit	KeysCommandsManageEdit
	Remove	KeysCommandsManageRemove
}
type KeysCommandsSongrequestErrorsGetCurrentSongVars struct {
}
type KeysCommandsSongrequestErrorsGetCurrentSong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestErrorsGetCurrentSong) IsTranslationKey() {
}
func (k KeysCommandsSongrequestErrorsGetCurrentSong) GetPath() string {
	return "commands.songrequest.errors.get_current_song"
}
func (k KeysCommandsSongrequestErrorsGetCurrentSong) GetPathSlice() []string {
	return []string{"commands", "songrequest", "errors", "get_current_song"}
}
func (k KeysCommandsSongrequestErrorsGetCurrentSong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestErrorsGetCurrentSong) SetVars(vars KeysCommandsSongrequestErrorsGetCurrentSongVars) twiri18n.TranslationKey[KeysCommandsSongrequestErrorsGetCurrentSongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestErrorsGetLatestSongVars struct {
}
type KeysCommandsSongrequestErrorsGetLatestSong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestErrorsGetLatestSong) IsTranslationKey() {
}
func (k KeysCommandsSongrequestErrorsGetLatestSong) GetPath() string {
	return "commands.songrequest.errors.get_latest_song"
}
func (k KeysCommandsSongrequestErrorsGetLatestSong) GetPathSlice() []string {
	return []string{"commands", "songrequest", "errors", "get_latest_song"}
}
func (k KeysCommandsSongrequestErrorsGetLatestSong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestErrorsGetLatestSong) SetVars(vars KeysCommandsSongrequestErrorsGetLatestSongVars) twiri18n.TranslationKey[KeysCommandsSongrequestErrorsGetLatestSongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestErrorsGetSettingsVars struct {
}
type KeysCommandsSongrequestErrorsGetSettings struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestErrorsGetSettings) IsTranslationKey() {
}
func (k KeysCommandsSongrequestErrorsGetSettings) GetPath() string {
	return "commands.songrequest.errors.get_settings"
}
func (k KeysCommandsSongrequestErrorsGetSettings) GetPathSlice() []string {
	return []string{"commands", "songrequest", "errors", "get_settings"}
}
func (k KeysCommandsSongrequestErrorsGetSettings) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestErrorsGetSettings) SetVars(vars KeysCommandsSongrequestErrorsGetSettingsVars) twiri18n.TranslationKey[KeysCommandsSongrequestErrorsGetSettingsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestErrorsGetCurrentVoteVars struct {
}
type KeysCommandsSongrequestErrorsGetCurrentVote struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestErrorsGetCurrentVote) IsTranslationKey() {
}
func (k KeysCommandsSongrequestErrorsGetCurrentVote) GetPath() string {
	return "commands.songrequest.errors.get_current_vote"
}
func (k KeysCommandsSongrequestErrorsGetCurrentVote) GetPathSlice() []string {
	return []string{"commands", "songrequest", "errors", "get_current_vote"}
}
func (k KeysCommandsSongrequestErrorsGetCurrentVote) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestErrorsGetCurrentVote) SetVars(vars KeysCommandsSongrequestErrorsGetCurrentVoteVars) twiri18n.TranslationKey[KeysCommandsSongrequestErrorsGetCurrentVoteVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestErrorsNotFoundVars struct {
}
type KeysCommandsSongrequestErrorsNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestErrorsNotFound) IsTranslationKey() {
}
func (k KeysCommandsSongrequestErrorsNotFound) GetPath() string {
	return "commands.songrequest.errors.not_found"
}
func (k KeysCommandsSongrequestErrorsNotFound) GetPathSlice() []string {
	return []string{"commands", "songrequest", "errors", "not_found"}
}
func (k KeysCommandsSongrequestErrorsNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestErrorsNotFound) SetVars(vars KeysCommandsSongrequestErrorsNotFoundVars) twiri18n.TranslationKey[KeysCommandsSongrequestErrorsNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestErrorsSearchSongVars struct {
}
type KeysCommandsSongrequestErrorsSearchSong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestErrorsSearchSong) IsTranslationKey() {
}
func (k KeysCommandsSongrequestErrorsSearchSong) GetPath() string {
	return "commands.songrequest.errors.search_song"
}
func (k KeysCommandsSongrequestErrorsSearchSong) GetPathSlice() []string {
	return []string{"commands", "songrequest", "errors", "search_song"}
}
func (k KeysCommandsSongrequestErrorsSearchSong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestErrorsSearchSong) SetVars(vars KeysCommandsSongrequestErrorsSearchSongVars) twiri18n.TranslationKey[KeysCommandsSongrequestErrorsSearchSongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestErrorsUpdateSongVars struct {
}
type KeysCommandsSongrequestErrorsUpdateSong struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestErrorsUpdateSong) IsTranslationKey() {
}
func (k KeysCommandsSongrequestErrorsUpdateSong) GetPath() string {
	return "commands.songrequest.errors.update_song"
}
func (k KeysCommandsSongrequestErrorsUpdateSong) GetPathSlice() []string {
	return []string{"commands", "songrequest", "errors", "update_song"}
}
func (k KeysCommandsSongrequestErrorsUpdateSong) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestErrorsUpdateSong) SetVars(vars KeysCommandsSongrequestErrorsUpdateSongVars) twiri18n.TranslationKey[KeysCommandsSongrequestErrorsUpdateSongVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestErrorsGetCurrentQueueCountVars struct {
}
type KeysCommandsSongrequestErrorsGetCurrentQueueCount struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestErrorsGetCurrentQueueCount) IsTranslationKey() {
}
func (k KeysCommandsSongrequestErrorsGetCurrentQueueCount) GetPath() string {
	return "commands.songrequest.errors.get_current_queue_count"
}
func (k KeysCommandsSongrequestErrorsGetCurrentQueueCount) GetPathSlice() []string {
	return []string{"commands", "songrequest", "errors", "get_current_queue_count"}
}
func (k KeysCommandsSongrequestErrorsGetCurrentQueueCount) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestErrorsGetCurrentQueueCount) SetVars(vars KeysCommandsSongrequestErrorsGetCurrentQueueCountVars) twiri18n.TranslationKey[KeysCommandsSongrequestErrorsGetCurrentQueueCountVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestErrorsGetSongsFromQueueVars struct {
}
type KeysCommandsSongrequestErrorsGetSongsFromQueue struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestErrorsGetSongsFromQueue) IsTranslationKey() {
}
func (k KeysCommandsSongrequestErrorsGetSongsFromQueue) GetPath() string {
	return "commands.songrequest.errors.get_songs_from_queue"
}
func (k KeysCommandsSongrequestErrorsGetSongsFromQueue) GetPathSlice() []string {
	return []string{"commands", "songrequest", "errors", "get_songs_from_queue"}
}
func (k KeysCommandsSongrequestErrorsGetSongsFromQueue) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestErrorsGetSongsFromQueue) SetVars(vars KeysCommandsSongrequestErrorsGetSongsFromQueueVars) twiri18n.TranslationKey[KeysCommandsSongrequestErrorsGetSongsFromQueueVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestErrorsRemoveSongFromQueueVars struct {
}
type KeysCommandsSongrequestErrorsRemoveSongFromQueue struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestErrorsRemoveSongFromQueue) IsTranslationKey() {
}
func (k KeysCommandsSongrequestErrorsRemoveSongFromQueue) GetPath() string {
	return "commands.songrequest.errors.remove_song_from_queue"
}
func (k KeysCommandsSongrequestErrorsRemoveSongFromQueue) GetPathSlice() []string {
	return []string{"commands", "songrequest", "errors", "remove_song_from_queue"}
}
func (k KeysCommandsSongrequestErrorsRemoveSongFromQueue) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestErrorsRemoveSongFromQueue) SetVars(vars KeysCommandsSongrequestErrorsRemoveSongFromQueueVars) twiri18n.TranslationKey[KeysCommandsSongrequestErrorsRemoveSongFromQueueVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestErrorsGetUsersCountVars struct {
}
type KeysCommandsSongrequestErrorsGetUsersCount struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestErrorsGetUsersCount) IsTranslationKey() {
}
func (k KeysCommandsSongrequestErrorsGetUsersCount) GetPath() string {
	return "commands.songrequest.errors.get_users_count"
}
func (k KeysCommandsSongrequestErrorsGetUsersCount) GetPathSlice() []string {
	return []string{"commands", "songrequest", "errors", "get_users_count"}
}
func (k KeysCommandsSongrequestErrorsGetUsersCount) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestErrorsGetUsersCount) SetVars(vars KeysCommandsSongrequestErrorsGetUsersCountVars) twiri18n.TranslationKey[KeysCommandsSongrequestErrorsGetUsersCountVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestErrorsGetVotesCountVars struct {
}
type KeysCommandsSongrequestErrorsGetVotesCount struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestErrorsGetVotesCount) IsTranslationKey() {
}
func (k KeysCommandsSongrequestErrorsGetVotesCount) GetPath() string {
	return "commands.songrequest.errors.get_votes_count"
}
func (k KeysCommandsSongrequestErrorsGetVotesCount) GetPathSlice() []string {
	return []string{"commands", "songrequest", "errors", "get_votes_count"}
}
func (k KeysCommandsSongrequestErrorsGetVotesCount) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestErrorsGetVotesCount) SetVars(vars KeysCommandsSongrequestErrorsGetVotesCountVars) twiri18n.TranslationKey[KeysCommandsSongrequestErrorsGetVotesCountVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestErrors struct {
	GetCurrentSong		KeysCommandsSongrequestErrorsGetCurrentSong
	GetLatestSong		KeysCommandsSongrequestErrorsGetLatestSong
	GetSettings		KeysCommandsSongrequestErrorsGetSettings
	GetCurrentVote		KeysCommandsSongrequestErrorsGetCurrentVote
	NotFound		KeysCommandsSongrequestErrorsNotFound
	SearchSong		KeysCommandsSongrequestErrorsSearchSong
	UpdateSong		KeysCommandsSongrequestErrorsUpdateSong
	GetCurrentQueueCount	KeysCommandsSongrequestErrorsGetCurrentQueueCount
	GetSongsFromQueue	KeysCommandsSongrequestErrorsGetSongsFromQueue
	RemoveSongFromQueue	KeysCommandsSongrequestErrorsRemoveSongFromQueue
	GetUsersCount		KeysCommandsSongrequestErrorsGetUsersCount
	GetVotesCount		KeysCommandsSongrequestErrorsGetVotesCount
}
type KeysCommandsSongrequestInfoOnlyCountSongsVars struct {
	SongsCount any
}
type KeysCommandsSongrequestInfoOnlyCountSongs struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestInfoOnlyCountSongs) IsTranslationKey() {
}
func (k KeysCommandsSongrequestInfoOnlyCountSongs) GetPath() string {
	return "commands.songrequest.info.only_count_songs"
}
func (k KeysCommandsSongrequestInfoOnlyCountSongs) GetPathSlice() []string {
	return []string{"commands", "songrequest", "info", "only_count_songs"}
}
func (k KeysCommandsSongrequestInfoOnlyCountSongs) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestInfoOnlyCountSongs) SetVars(vars KeysCommandsSongrequestInfoOnlyCountSongsVars) twiri18n.TranslationKey[KeysCommandsSongrequestInfoOnlyCountSongsVars] {
	k.Vars = twiri18n.Vars{"songsCount": vars.SongsCount}
	return k
}

type KeysCommandsSongrequestInfoDeleteVars struct {
	SongTitle any
}
type KeysCommandsSongrequestInfoDelete struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestInfoDelete) IsTranslationKey() {
}
func (k KeysCommandsSongrequestInfoDelete) GetPath() string {
	return "commands.songrequest.info.delete"
}
func (k KeysCommandsSongrequestInfoDelete) GetPathSlice() []string {
	return []string{"commands", "songrequest", "info", "delete"}
}
func (k KeysCommandsSongrequestInfoDelete) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestInfoDelete) SetVars(vars KeysCommandsSongrequestInfoDeleteVars) twiri18n.TranslationKey[KeysCommandsSongrequestInfoDeleteVars] {
	k.Vars = twiri18n.Vars{"songTitle": vars.SongTitle}
	return k
}

type KeysCommandsSongrequestInfoSongSkippedVars struct {
	SongTitle any
}
type KeysCommandsSongrequestInfoSongSkipped struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestInfoSongSkipped) IsTranslationKey() {
}
func (k KeysCommandsSongrequestInfoSongSkipped) GetPath() string {
	return "commands.songrequest.info.song_skipped"
}
func (k KeysCommandsSongrequestInfoSongSkipped) GetPathSlice() []string {
	return []string{"commands", "songrequest", "info", "song_skipped"}
}
func (k KeysCommandsSongrequestInfoSongSkipped) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestInfoSongSkipped) SetVars(vars KeysCommandsSongrequestInfoSongSkippedVars) twiri18n.TranslationKey[KeysCommandsSongrequestInfoSongSkippedVars] {
	k.Vars = twiri18n.Vars{"songTitle": vars.SongTitle}
	return k
}

type KeysCommandsSongrequestInfoNoRequestedSongsVars struct {
}
type KeysCommandsSongrequestInfoNoRequestedSongs struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestInfoNoRequestedSongs) IsTranslationKey() {
}
func (k KeysCommandsSongrequestInfoNoRequestedSongs) GetPath() string {
	return "commands.songrequest.info.no_requested_songs"
}
func (k KeysCommandsSongrequestInfoNoRequestedSongs) GetPathSlice() []string {
	return []string{"commands", "songrequest", "info", "no_requested_songs"}
}
func (k KeysCommandsSongrequestInfoNoRequestedSongs) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestInfoNoRequestedSongs) SetVars(vars KeysCommandsSongrequestInfoNoRequestedSongsVars) twiri18n.TranslationKey[KeysCommandsSongrequestInfoNoRequestedSongsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestInfo struct {
	OnlyCountSongs		KeysCommandsSongrequestInfoOnlyCountSongs
	Delete			KeysCommandsSongrequestInfoDelete
	SongSkipped		KeysCommandsSongrequestInfoSongSkipped
	NoRequestedSongs	KeysCommandsSongrequestInfoNoRequestedSongs
}
type KeysCommandsSongrequestValidateErrorsRestrictionsOnUserVars struct {
}
type KeysCommandsSongrequestValidateErrorsRestrictionsOnUser struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestValidateErrorsRestrictionsOnUser) IsTranslationKey() {
}
func (k KeysCommandsSongrequestValidateErrorsRestrictionsOnUser) GetPath() string {
	return "commands.songrequest.validate.errors.restrictions_on_user"
}
func (k KeysCommandsSongrequestValidateErrorsRestrictionsOnUser) GetPathSlice() []string {
	return []string{"commands", "songrequest", "validate", "errors", "restrictions_on_user"}
}
func (k KeysCommandsSongrequestValidateErrorsRestrictionsOnUser) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestValidateErrorsRestrictionsOnUser) SetVars(vars KeysCommandsSongrequestValidateErrorsRestrictionsOnUserVars) twiri18n.TranslationKey[KeysCommandsSongrequestValidateErrorsRestrictionsOnUserVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestValidateErrorsInternalErrorVars struct {
}
type KeysCommandsSongrequestValidateErrorsInternalError struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestValidateErrorsInternalError) IsTranslationKey() {
}
func (k KeysCommandsSongrequestValidateErrorsInternalError) GetPath() string {
	return "commands.songrequest.validate.errors.internal_error"
}
func (k KeysCommandsSongrequestValidateErrorsInternalError) GetPathSlice() []string {
	return []string{"commands", "songrequest", "validate", "errors", "internal_error"}
}
func (k KeysCommandsSongrequestValidateErrorsInternalError) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestValidateErrorsInternalError) SetVars(vars KeysCommandsSongrequestValidateErrorsInternalErrorVars) twiri18n.TranslationKey[KeysCommandsSongrequestValidateErrorsInternalErrorVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestValidateErrorsNeedFollowVars struct {
}
type KeysCommandsSongrequestValidateErrorsNeedFollow struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSongrequestValidateErrorsNeedFollow) IsTranslationKey() {
}
func (k KeysCommandsSongrequestValidateErrorsNeedFollow) GetPath() string {
	return "commands.songrequest.validate.errors.need_follow"
}
func (k KeysCommandsSongrequestValidateErrorsNeedFollow) GetPathSlice() []string {
	return []string{"commands", "songrequest", "validate", "errors", "need_follow"}
}
func (k KeysCommandsSongrequestValidateErrorsNeedFollow) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSongrequestValidateErrorsNeedFollow) SetVars(vars KeysCommandsSongrequestValidateErrorsNeedFollowVars) twiri18n.TranslationKey[KeysCommandsSongrequestValidateErrorsNeedFollowVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSongrequestValidateErrors struct {
	RestrictionsOnUser	KeysCommandsSongrequestValidateErrorsRestrictionsOnUser
	InternalError		KeysCommandsSongrequestValidateErrorsInternalError
	NeedFollow		KeysCommandsSongrequestValidateErrorsNeedFollow
}
type KeysCommandsSongrequestValidate struct {
	Errors KeysCommandsSongrequestValidateErrors
}
type KeysCommandsSongrequest struct {
	Errors		KeysCommandsSongrequestErrors
	Info		KeysCommandsSongrequestInfo
	Validate	KeysCommandsSongrequestValidate
}
type KeysCommandsDudesErrorsSpriteInvalidVars struct {
	AvailableSprites any
}
type KeysCommandsDudesErrorsSpriteInvalid struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesErrorsSpriteInvalid) IsTranslationKey() {
}
func (k KeysCommandsDudesErrorsSpriteInvalid) GetPath() string {
	return "commands.dudes.errors.sprite_invalid"
}
func (k KeysCommandsDudesErrorsSpriteInvalid) GetPathSlice() []string {
	return []string{"commands", "dudes", "errors", "sprite_invalid"}
}
func (k KeysCommandsDudesErrorsSpriteInvalid) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesErrorsSpriteInvalid) SetVars(vars KeysCommandsDudesErrorsSpriteInvalidVars) twiri18n.TranslationKey[KeysCommandsDudesErrorsSpriteInvalidVars] {
	k.Vars = twiri18n.Vars{"availableSprites": vars.AvailableSprites}
	return k
}

type KeysCommandsDudesErrorsColorCannotTriggerVars struct {
}
type KeysCommandsDudesErrorsColorCannotTrigger struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesErrorsColorCannotTrigger) IsTranslationKey() {
}
func (k KeysCommandsDudesErrorsColorCannotTrigger) GetPath() string {
	return "commands.dudes.errors.color_cannot_trigger"
}
func (k KeysCommandsDudesErrorsColorCannotTrigger) GetPathSlice() []string {
	return []string{"commands", "dudes", "errors", "color_cannot_trigger"}
}
func (k KeysCommandsDudesErrorsColorCannotTrigger) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesErrorsColorCannotTrigger) SetVars(vars KeysCommandsDudesErrorsColorCannotTriggerVars) twiri18n.TranslationKey[KeysCommandsDudesErrorsColorCannotTriggerVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesErrorsGrowCannotTriggerVars struct {
}
type KeysCommandsDudesErrorsGrowCannotTrigger struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesErrorsGrowCannotTrigger) IsTranslationKey() {
}
func (k KeysCommandsDudesErrorsGrowCannotTrigger) GetPath() string {
	return "commands.dudes.errors.grow_cannot_trigger"
}
func (k KeysCommandsDudesErrorsGrowCannotTrigger) GetPathSlice() []string {
	return []string{"commands", "dudes", "errors", "grow_cannot_trigger"}
}
func (k KeysCommandsDudesErrorsGrowCannotTrigger) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesErrorsGrowCannotTrigger) SetVars(vars KeysCommandsDudesErrorsGrowCannotTriggerVars) twiri18n.TranslationKey[KeysCommandsDudesErrorsGrowCannotTriggerVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesErrorsJumpCannotTriggerVars struct {
}
type KeysCommandsDudesErrorsJumpCannotTrigger struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesErrorsJumpCannotTrigger) IsTranslationKey() {
}
func (k KeysCommandsDudesErrorsJumpCannotTrigger) GetPath() string {
	return "commands.dudes.errors.jump_cannot_trigger"
}
func (k KeysCommandsDudesErrorsJumpCannotTrigger) GetPathSlice() []string {
	return []string{"commands", "dudes", "errors", "jump_cannot_trigger"}
}
func (k KeysCommandsDudesErrorsJumpCannotTrigger) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesErrorsJumpCannotTrigger) SetVars(vars KeysCommandsDudesErrorsJumpCannotTriggerVars) twiri18n.TranslationKey[KeysCommandsDudesErrorsJumpCannotTriggerVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesErrorsLeaveCannotTriggerVars struct {
}
type KeysCommandsDudesErrorsLeaveCannotTrigger struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesErrorsLeaveCannotTrigger) IsTranslationKey() {
}
func (k KeysCommandsDudesErrorsLeaveCannotTrigger) GetPath() string {
	return "commands.dudes.errors.leave_cannot_trigger"
}
func (k KeysCommandsDudesErrorsLeaveCannotTrigger) GetPathSlice() []string {
	return []string{"commands", "dudes", "errors", "leave_cannot_trigger"}
}
func (k KeysCommandsDudesErrorsLeaveCannotTrigger) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesErrorsLeaveCannotTrigger) SetVars(vars KeysCommandsDudesErrorsLeaveCannotTriggerVars) twiri18n.TranslationKey[KeysCommandsDudesErrorsLeaveCannotTriggerVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesErrorsSpriteCannotTriggerVars struct {
}
type KeysCommandsDudesErrorsSpriteCannotTrigger struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesErrorsSpriteCannotTrigger) IsTranslationKey() {
}
func (k KeysCommandsDudesErrorsSpriteCannotTrigger) GetPath() string {
	return "commands.dudes.errors.sprite_cannot_trigger"
}
func (k KeysCommandsDudesErrorsSpriteCannotTrigger) GetPathSlice() []string {
	return []string{"commands", "dudes", "errors", "sprite_cannot_trigger"}
}
func (k KeysCommandsDudesErrorsSpriteCannotTrigger) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesErrorsSpriteCannotTrigger) SetVars(vars KeysCommandsDudesErrorsSpriteCannotTriggerVars) twiri18n.TranslationKey[KeysCommandsDudesErrorsSpriteCannotTriggerVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesErrorsColorInvalidVars struct {
}
type KeysCommandsDudesErrorsColorInvalid struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesErrorsColorInvalid) IsTranslationKey() {
}
func (k KeysCommandsDudesErrorsColorInvalid) GetPath() string {
	return "commands.dudes.errors.color_invalid"
}
func (k KeysCommandsDudesErrorsColorInvalid) GetPathSlice() []string {
	return []string{"commands", "dudes", "errors", "color_invalid"}
}
func (k KeysCommandsDudesErrorsColorInvalid) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesErrorsColorInvalid) SetVars(vars KeysCommandsDudesErrorsColorInvalidVars) twiri18n.TranslationKey[KeysCommandsDudesErrorsColorInvalidVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesErrors struct {
	SpriteInvalid		KeysCommandsDudesErrorsSpriteInvalid
	ColorCannotTrigger	KeysCommandsDudesErrorsColorCannotTrigger
	GrowCannotTrigger	KeysCommandsDudesErrorsGrowCannotTrigger
	JumpCannotTrigger	KeysCommandsDudesErrorsJumpCannotTrigger
	LeaveCannotTrigger	KeysCommandsDudesErrorsLeaveCannotTrigger
	SpriteCannotTrigger	KeysCommandsDudesErrorsSpriteCannotTrigger
	ColorInvalid		KeysCommandsDudesErrorsColorInvalid
}
type KeysCommandsDudesInfoColorChangedVars struct {
	DudeColor any
}
type KeysCommandsDudesInfoColorChanged struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesInfoColorChanged) IsTranslationKey() {
}
func (k KeysCommandsDudesInfoColorChanged) GetPath() string {
	return "commands.dudes.info.color_changed"
}
func (k KeysCommandsDudesInfoColorChanged) GetPathSlice() []string {
	return []string{"commands", "dudes", "info", "color_changed"}
}
func (k KeysCommandsDudesInfoColorChanged) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesInfoColorChanged) SetVars(vars KeysCommandsDudesInfoColorChangedVars) twiri18n.TranslationKey[KeysCommandsDudesInfoColorChangedVars] {
	k.Vars = twiri18n.Vars{"dudeColor": vars.DudeColor}
	return k
}

type KeysCommandsDudesInfoSpriteChangedVars struct {
	DudeSprite any
}
type KeysCommandsDudesInfoSpriteChanged struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesInfoSpriteChanged) IsTranslationKey() {
}
func (k KeysCommandsDudesInfoSpriteChanged) GetPath() string {
	return "commands.dudes.info.sprite_changed"
}
func (k KeysCommandsDudesInfoSpriteChanged) GetPathSlice() []string {
	return []string{"commands", "dudes", "info", "sprite_changed"}
}
func (k KeysCommandsDudesInfoSpriteChanged) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesInfoSpriteChanged) SetVars(vars KeysCommandsDudesInfoSpriteChangedVars) twiri18n.TranslationKey[KeysCommandsDudesInfoSpriteChangedVars] {
	k.Vars = twiri18n.Vars{"dudeSprite": vars.DudeSprite}
	return k
}

type KeysCommandsDudesInfoColorVars struct {
	DudeColor any
}
type KeysCommandsDudesInfoColor struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesInfoColor) IsTranslationKey() {
}
func (k KeysCommandsDudesInfoColor) GetPath() string {
	return "commands.dudes.info.color"
}
func (k KeysCommandsDudesInfoColor) GetPathSlice() []string {
	return []string{"commands", "dudes", "info", "color"}
}
func (k KeysCommandsDudesInfoColor) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesInfoColor) SetVars(vars KeysCommandsDudesInfoColorVars) twiri18n.TranslationKey[KeysCommandsDudesInfoColorVars] {
	k.Vars = twiri18n.Vars{"dudeColor": vars.DudeColor}
	return k
}

type KeysCommandsDudesInfoSpriteVars struct {
	DudeSprite any
}
type KeysCommandsDudesInfoSprite struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesInfoSprite) IsTranslationKey() {
}
func (k KeysCommandsDudesInfoSprite) GetPath() string {
	return "commands.dudes.info.sprite"
}
func (k KeysCommandsDudesInfoSprite) GetPathSlice() []string {
	return []string{"commands", "dudes", "info", "sprite"}
}
func (k KeysCommandsDudesInfoSprite) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesInfoSprite) SetVars(vars KeysCommandsDudesInfoSpriteVars) twiri18n.TranslationKey[KeysCommandsDudesInfoSpriteVars] {
	k.Vars = twiri18n.Vars{"dudeSprite": vars.DudeSprite}
	return k
}

type KeysCommandsDudesInfoColorRequiredVars struct {
}
type KeysCommandsDudesInfoColorRequired struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesInfoColorRequired) IsTranslationKey() {
}
func (k KeysCommandsDudesInfoColorRequired) GetPath() string {
	return "commands.dudes.info.color_required"
}
func (k KeysCommandsDudesInfoColorRequired) GetPathSlice() []string {
	return []string{"commands", "dudes", "info", "color_required"}
}
func (k KeysCommandsDudesInfoColorRequired) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesInfoColorRequired) SetVars(vars KeysCommandsDudesInfoColorRequiredVars) twiri18n.TranslationKey[KeysCommandsDudesInfoColorRequiredVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesInfoSpriteRequiredVars struct {
	AvailableSprites any
}
type KeysCommandsDudesInfoSpriteRequired struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesInfoSpriteRequired) IsTranslationKey() {
}
func (k KeysCommandsDudesInfoSpriteRequired) GetPath() string {
	return "commands.dudes.info.sprite_required"
}
func (k KeysCommandsDudesInfoSpriteRequired) GetPathSlice() []string {
	return []string{"commands", "dudes", "info", "sprite_required"}
}
func (k KeysCommandsDudesInfoSpriteRequired) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesInfoSpriteRequired) SetVars(vars KeysCommandsDudesInfoSpriteRequiredVars) twiri18n.TranslationKey[KeysCommandsDudesInfoSpriteRequiredVars] {
	k.Vars = twiri18n.Vars{"availableSprites": vars.AvailableSprites}
	return k
}

type KeysCommandsDudesInfoColorResetVars struct {
}
type KeysCommandsDudesInfoColorReset struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsDudesInfoColorReset) IsTranslationKey() {
}
func (k KeysCommandsDudesInfoColorReset) GetPath() string {
	return "commands.dudes.info.color_reset"
}
func (k KeysCommandsDudesInfoColorReset) GetPathSlice() []string {
	return []string{"commands", "dudes", "info", "color_reset"}
}
func (k KeysCommandsDudesInfoColorReset) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsDudesInfoColorReset) SetVars(vars KeysCommandsDudesInfoColorResetVars) twiri18n.TranslationKey[KeysCommandsDudesInfoColorResetVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsDudesInfo struct {
	ColorChanged	KeysCommandsDudesInfoColorChanged
	SpriteChanged	KeysCommandsDudesInfoSpriteChanged
	Color		KeysCommandsDudesInfoColor
	Sprite		KeysCommandsDudesInfoSprite
	ColorRequired	KeysCommandsDudesInfoColorRequired
	SpriteRequired	KeysCommandsDudesInfoSpriteRequired
	ColorReset	KeysCommandsDudesInfoColorReset
}
type KeysCommandsDudes struct {
	Errors	KeysCommandsDudesErrors
	Info	KeysCommandsDudesInfo
}
type KeysCommandsShoutoutErrorsBotHaveNoPermissionsVars struct {
}
type KeysCommandsShoutoutErrorsBotHaveNoPermissions struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShoutoutErrorsBotHaveNoPermissions) IsTranslationKey() {
}
func (k KeysCommandsShoutoutErrorsBotHaveNoPermissions) GetPath() string {
	return "commands.shoutout.errors.bot_have_no_permissions"
}
func (k KeysCommandsShoutoutErrorsBotHaveNoPermissions) GetPathSlice() []string {
	return []string{"commands", "shoutout", "errors", "bot_have_no_permissions"}
}
func (k KeysCommandsShoutoutErrorsBotHaveNoPermissions) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShoutoutErrorsBotHaveNoPermissions) SetVars(vars KeysCommandsShoutoutErrorsBotHaveNoPermissionsVars) twiri18n.TranslationKey[KeysCommandsShoutoutErrorsBotHaveNoPermissionsVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsShoutoutErrors struct {
	BotHaveNoPermissions KeysCommandsShoutoutErrorsBotHaveNoPermissions
}
type KeysCommandsShoutoutResponseOnlineVars struct {
	UserName	any
	CategoryName	any
	Title		any
	Viewers		any
}
type KeysCommandsShoutoutResponseOnline struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShoutoutResponseOnline) IsTranslationKey() {
}
func (k KeysCommandsShoutoutResponseOnline) GetPath() string {
	return "commands.shoutout.response_online"
}
func (k KeysCommandsShoutoutResponseOnline) GetPathSlice() []string {
	return []string{"commands", "shoutout", "response_online"}
}
func (k KeysCommandsShoutoutResponseOnline) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShoutoutResponseOnline) SetVars(vars KeysCommandsShoutoutResponseOnlineVars) twiri18n.TranslationKey[KeysCommandsShoutoutResponseOnlineVars] {
	k.Vars = twiri18n.Vars{"userName": vars.UserName, "categoryName": vars.CategoryName, "title": vars.Title, "viewers": vars.Viewers}
	return k
}

type KeysCommandsShoutoutResponseOfflineVars struct {
	UserName	any
	CategoryName	any
	Title		any
}
type KeysCommandsShoutoutResponseOffline struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsShoutoutResponseOffline) IsTranslationKey() {
}
func (k KeysCommandsShoutoutResponseOffline) GetPath() string {
	return "commands.shoutout.response_offline"
}
func (k KeysCommandsShoutoutResponseOffline) GetPathSlice() []string {
	return []string{"commands", "shoutout", "response_offline"}
}
func (k KeysCommandsShoutoutResponseOffline) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsShoutoutResponseOffline) SetVars(vars KeysCommandsShoutoutResponseOfflineVars) twiri18n.TranslationKey[KeysCommandsShoutoutResponseOfflineVars] {
	k.Vars = twiri18n.Vars{"userName": vars.UserName, "categoryName": vars.CategoryName, "title": vars.Title}
	return k
}

type KeysCommandsShoutout struct {
	Errors		KeysCommandsShoutoutErrors
	ResponseOnline	KeysCommandsShoutoutResponseOnline
	ResponseOffline	KeysCommandsShoutoutResponseOffline
}
type KeysCommandsSeventvProfileInfoResponseVars struct {
	ProfileName		any
	PaintName		any
	UnlockedPaints		any
	Roles			any
	EditorCount		any
	EmoteSetName		any
	EmoteSetCount		any
	EmoteSetCapacity	any
	ProfileCreatedAt	any
}
type KeysCommandsSeventvProfileInfoResponse struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvProfileInfoResponse) IsTranslationKey() {
}
func (k KeysCommandsSeventvProfileInfoResponse) GetPath() string {
	return "commands.7tv.profile_info.response"
}
func (k KeysCommandsSeventvProfileInfoResponse) GetPathSlice() []string {
	return []string{"commands", "7tv", "profile_info", "response"}
}
func (k KeysCommandsSeventvProfileInfoResponse) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvProfileInfoResponse) SetVars(vars KeysCommandsSeventvProfileInfoResponseVars) twiri18n.TranslationKey[KeysCommandsSeventvProfileInfoResponseVars] {
	k.Vars = twiri18n.Vars{"profileName": vars.ProfileName, "paintName": vars.PaintName, "unlockedPaints": vars.UnlockedPaints, "roles": vars.Roles, "editorCount": vars.EditorCount, "emoteSetName": vars.EmoteSetName, "emoteSetCount": vars.EmoteSetCount, "emoteSetCapacity": vars.EmoteSetCapacity, "profileCreatedAt": vars.ProfileCreatedAt}
	return k
}

type KeysCommandsSeventvProfileInfo struct {
	Response KeysCommandsSeventvProfileInfoResponse
}
type KeysCommandsSeventvHintsEmoteForAddArgAliasVars struct {
}
type KeysCommandsSeventvHintsEmoteForAddArgAlias struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvHintsEmoteForAddArgAlias) IsTranslationKey() {
}
func (k KeysCommandsSeventvHintsEmoteForAddArgAlias) GetPath() string {
	return "commands.7tv.hints.emoteForAddArgAlias"
}
func (k KeysCommandsSeventvHintsEmoteForAddArgAlias) GetPathSlice() []string {
	return []string{"commands", "7tv", "hints", "emoteForAddArgAlias"}
}
func (k KeysCommandsSeventvHintsEmoteForAddArgAlias) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvHintsEmoteForAddArgAlias) SetVars(vars KeysCommandsSeventvHintsEmoteForAddArgAliasVars) twiri18n.TranslationKey[KeysCommandsSeventvHintsEmoteForAddArgAliasVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvHintsEmoteForCopyAliasVars struct {
}
type KeysCommandsSeventvHintsEmoteForCopyAlias struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvHintsEmoteForCopyAlias) IsTranslationKey() {
}
func (k KeysCommandsSeventvHintsEmoteForCopyAlias) GetPath() string {
	return "commands.7tv.hints.emoteForCopyAlias"
}
func (k KeysCommandsSeventvHintsEmoteForCopyAlias) GetPathSlice() []string {
	return []string{"commands", "7tv", "hints", "emoteForCopyAlias"}
}
func (k KeysCommandsSeventvHintsEmoteForCopyAlias) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvHintsEmoteForCopyAlias) SetVars(vars KeysCommandsSeventvHintsEmoteForCopyAliasVars) twiri18n.TranslationKey[KeysCommandsSeventvHintsEmoteForCopyAliasVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvHintsEmoteForCopyArgNameVars struct {
}
type KeysCommandsSeventvHintsEmoteForCopyArgName struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvHintsEmoteForCopyArgName) IsTranslationKey() {
}
func (k KeysCommandsSeventvHintsEmoteForCopyArgName) GetPath() string {
	return "commands.7tv.hints.emoteForCopyArgName"
}
func (k KeysCommandsSeventvHintsEmoteForCopyArgName) GetPathSlice() []string {
	return []string{"commands", "7tv", "hints", "emoteForCopyArgName"}
}
func (k KeysCommandsSeventvHintsEmoteForCopyArgName) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvHintsEmoteForCopyArgName) SetVars(vars KeysCommandsSeventvHintsEmoteForCopyArgNameVars) twiri18n.TranslationKey[KeysCommandsSeventvHintsEmoteForCopyArgNameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvHintsCopySetChannelNameVars struct {
}
type KeysCommandsSeventvHintsCopySetChannelName struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvHintsCopySetChannelName) IsTranslationKey() {
}
func (k KeysCommandsSeventvHintsCopySetChannelName) GetPath() string {
	return "commands.7tv.hints.copySetChannelName"
}
func (k KeysCommandsSeventvHintsCopySetChannelName) GetPathSlice() []string {
	return []string{"commands", "7tv", "hints", "copySetChannelName"}
}
func (k KeysCommandsSeventvHintsCopySetChannelName) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvHintsCopySetChannelName) SetVars(vars KeysCommandsSeventvHintsCopySetChannelNameVars) twiri18n.TranslationKey[KeysCommandsSeventvHintsCopySetChannelNameVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvHintsCopySetNameOfSetVars struct {
}
type KeysCommandsSeventvHintsCopySetNameOfSet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvHintsCopySetNameOfSet) IsTranslationKey() {
}
func (k KeysCommandsSeventvHintsCopySetNameOfSet) GetPath() string {
	return "commands.7tv.hints.copySetNameOfSet"
}
func (k KeysCommandsSeventvHintsCopySetNameOfSet) GetPathSlice() []string {
	return []string{"commands", "7tv", "hints", "copySetNameOfSet"}
}
func (k KeysCommandsSeventvHintsCopySetNameOfSet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvHintsCopySetNameOfSet) SetVars(vars KeysCommandsSeventvHintsCopySetNameOfSetVars) twiri18n.TranslationKey[KeysCommandsSeventvHintsCopySetNameOfSetVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvHintsEmoteForAddArgLinkVars struct {
}
type KeysCommandsSeventvHintsEmoteForAddArgLink struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvHintsEmoteForAddArgLink) IsTranslationKey() {
}
func (k KeysCommandsSeventvHintsEmoteForAddArgLink) GetPath() string {
	return "commands.7tv.hints.emoteForAddArgLink"
}
func (k KeysCommandsSeventvHintsEmoteForAddArgLink) GetPathSlice() []string {
	return []string{"commands", "7tv", "hints", "emoteForAddArgLink"}
}
func (k KeysCommandsSeventvHintsEmoteForAddArgLink) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvHintsEmoteForAddArgLink) SetVars(vars KeysCommandsSeventvHintsEmoteForAddArgLinkVars) twiri18n.TranslationKey[KeysCommandsSeventvHintsEmoteForAddArgLinkVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvHints struct {
	EmoteForAddArgAlias	KeysCommandsSeventvHintsEmoteForAddArgAlias
	EmoteForCopyAlias	KeysCommandsSeventvHintsEmoteForCopyAlias
	EmoteForCopyArgName	KeysCommandsSeventvHintsEmoteForCopyArgName
	CopySetChannelName	KeysCommandsSeventvHintsCopySetChannelName
	CopySetNameOfSet	KeysCommandsSeventvHintsCopySetNameOfSet
	EmoteForAddArgLink	KeysCommandsSeventvHintsEmoteForAddArgLink
}
type KeysCommandsSeventvErrorsEmotesetNotActiveVars struct {
}
type KeysCommandsSeventvErrorsEmotesetNotActive struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsEmotesetNotActive) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsEmotesetNotActive) GetPath() string {
	return "commands.7tv.errors.emoteset_not_active"
}
func (k KeysCommandsSeventvErrorsEmotesetNotActive) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "emoteset_not_active"}
}
func (k KeysCommandsSeventvErrorsEmotesetNotActive) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsEmotesetNotActive) SetVars(vars KeysCommandsSeventvErrorsEmotesetNotActiveVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsEmotesetNotActiveVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvErrorsEmotesetBroadcasterNotActiveVars struct {
}
type KeysCommandsSeventvErrorsEmotesetBroadcasterNotActive struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsEmotesetBroadcasterNotActive) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsEmotesetBroadcasterNotActive) GetPath() string {
	return "commands.7tv.errors.emoteset_broadcaster_not_active"
}
func (k KeysCommandsSeventvErrorsEmotesetBroadcasterNotActive) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "emoteset_broadcaster_not_active"}
}
func (k KeysCommandsSeventvErrorsEmotesetBroadcasterNotActive) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsEmotesetBroadcasterNotActive) SetVars(vars KeysCommandsSeventvErrorsEmotesetBroadcasterNotActiveVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsEmotesetBroadcasterNotActiveVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvErrorsProfileNotFoundVars struct {
}
type KeysCommandsSeventvErrorsProfileNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsProfileNotFound) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsProfileNotFound) GetPath() string {
	return "commands.7tv.errors.profile_not_found"
}
func (k KeysCommandsSeventvErrorsProfileNotFound) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "profile_not_found"}
}
func (k KeysCommandsSeventvErrorsProfileNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsProfileNotFound) SetVars(vars KeysCommandsSeventvErrorsProfileNotFoundVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsProfileNotFoundVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvErrorsEmoteFailedToRemoveVars struct {
	Reason any
}
type KeysCommandsSeventvErrorsEmoteFailedToRemove struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsEmoteFailedToRemove) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsEmoteFailedToRemove) GetPath() string {
	return "commands.7tv.errors.emote_failed_to_remove"
}
func (k KeysCommandsSeventvErrorsEmoteFailedToRemove) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "emote_failed_to_remove"}
}
func (k KeysCommandsSeventvErrorsEmoteFailedToRemove) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsEmoteFailedToRemove) SetVars(vars KeysCommandsSeventvErrorsEmoteFailedToRemoveVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsEmoteFailedToRemoveVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsSeventvErrorsEmoteFailedToRenameVars struct {
	Reason any
}
type KeysCommandsSeventvErrorsEmoteFailedToRename struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsEmoteFailedToRename) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsEmoteFailedToRename) GetPath() string {
	return "commands.7tv.errors.emote_failed_to_rename"
}
func (k KeysCommandsSeventvErrorsEmoteFailedToRename) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "emote_failed_to_rename"}
}
func (k KeysCommandsSeventvErrorsEmoteFailedToRename) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsEmoteFailedToRename) SetVars(vars KeysCommandsSeventvErrorsEmoteFailedToRenameVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsEmoteFailedToRenameVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsSeventvErrorsEmoteNotFoundVars struct {
	EmoteName any
}
type KeysCommandsSeventvErrorsEmoteNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsEmoteNotFound) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsEmoteNotFound) GetPath() string {
	return "commands.7tv.errors.emote_not_found"
}
func (k KeysCommandsSeventvErrorsEmoteNotFound) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "emote_not_found"}
}
func (k KeysCommandsSeventvErrorsEmoteNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsEmoteNotFound) SetVars(vars KeysCommandsSeventvErrorsEmoteNotFoundVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsEmoteNotFoundVars] {
	k.Vars = twiri18n.Vars{"emoteName": vars.EmoteName}
	return k
}

type KeysCommandsSeventvErrorsEmoteNotFoundInEmotesetVars struct {
	EmoteName	any
	EmoteSet	any
}
type KeysCommandsSeventvErrorsEmoteNotFoundInEmoteset struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsEmoteNotFoundInEmoteset) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsEmoteNotFoundInEmoteset) GetPath() string {
	return "commands.7tv.errors.emote_not_found_in_emoteset"
}
func (k KeysCommandsSeventvErrorsEmoteNotFoundInEmoteset) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "emote_not_found_in_emoteset"}
}
func (k KeysCommandsSeventvErrorsEmoteNotFoundInEmoteset) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsEmoteNotFoundInEmoteset) SetVars(vars KeysCommandsSeventvErrorsEmoteNotFoundInEmotesetVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsEmoteNotFoundInEmotesetVars] {
	k.Vars = twiri18n.Vars{"emoteName": vars.EmoteName, "emoteSet": vars.EmoteSet}
	return k
}

type KeysCommandsSeventvErrorsEmotesetNotFoundVars struct {
	EmoteName any
}
type KeysCommandsSeventvErrorsEmotesetNotFound struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsEmotesetNotFound) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsEmotesetNotFound) GetPath() string {
	return "commands.7tv.errors.emoteset_not_found"
}
func (k KeysCommandsSeventvErrorsEmotesetNotFound) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "emoteset_not_found"}
}
func (k KeysCommandsSeventvErrorsEmotesetNotFound) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsEmotesetNotFound) SetVars(vars KeysCommandsSeventvErrorsEmotesetNotFoundVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsEmotesetNotFoundVars] {
	k.Vars = twiri18n.Vars{"emoteName": vars.EmoteName}
	return k
}

type KeysCommandsSeventvErrorsProfileFailedToGetVars struct {
	Reason any
}
type KeysCommandsSeventvErrorsProfileFailedToGet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsProfileFailedToGet) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsProfileFailedToGet) GetPath() string {
	return "commands.7tv.errors.profile_failed_to_get"
}
func (k KeysCommandsSeventvErrorsProfileFailedToGet) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "profile_failed_to_get"}
}
func (k KeysCommandsSeventvErrorsProfileFailedToGet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsProfileFailedToGet) SetVars(vars KeysCommandsSeventvErrorsProfileFailedToGetVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsProfileFailedToGetVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsSeventvErrorsEmoteFailedToAddVars struct {
	Reason any
}
type KeysCommandsSeventvErrorsEmoteFailedToAdd struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsEmoteFailedToAdd) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsEmoteFailedToAdd) GetPath() string {
	return "commands.7tv.errors.emote_failed_to_add"
}
func (k KeysCommandsSeventvErrorsEmoteFailedToAdd) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "emote_failed_to_add"}
}
func (k KeysCommandsSeventvErrorsEmoteFailedToAdd) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsEmoteFailedToAdd) SetVars(vars KeysCommandsSeventvErrorsEmoteFailedToAddVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsEmoteFailedToAddVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsSeventvErrorsEmoteFailedToFetchVars struct {
	Reason any
}
type KeysCommandsSeventvErrorsEmoteFailedToFetch struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsEmoteFailedToFetch) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsEmoteFailedToFetch) GetPath() string {
	return "commands.7tv.errors.emote_failed_to_fetch"
}
func (k KeysCommandsSeventvErrorsEmoteFailedToFetch) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "emote_failed_to_fetch"}
}
func (k KeysCommandsSeventvErrorsEmoteFailedToFetch) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsEmoteFailedToFetch) SetVars(vars KeysCommandsSeventvErrorsEmoteFailedToFetchVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsEmoteFailedToFetchVars] {
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsSeventvErrorsEmoteNotFoundInChannelVars struct {
	EmoteSearch any
}
type KeysCommandsSeventvErrorsEmoteNotFoundInChannel struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsEmoteNotFoundInChannel) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsEmoteNotFoundInChannel) GetPath() string {
	return "commands.7tv.errors.emote_not_found_in_channel"
}
func (k KeysCommandsSeventvErrorsEmoteNotFoundInChannel) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "emote_not_found_in_channel"}
}
func (k KeysCommandsSeventvErrorsEmoteNotFoundInChannel) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsEmoteNotFoundInChannel) SetVars(vars KeysCommandsSeventvErrorsEmoteNotFoundInChannelVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsEmoteNotFoundInChannelVars] {
	k.Vars = twiri18n.Vars{"emoteSearch": vars.EmoteSearch}
	return k
}

type KeysCommandsSeventvErrorsEmoteAlreadyExistInChannelVars struct {
	EmoteName any
}
type KeysCommandsSeventvErrorsEmoteAlreadyExistInChannel struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsEmoteAlreadyExistInChannel) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsEmoteAlreadyExistInChannel) GetPath() string {
	return "commands.7tv.errors.emote_already_exist_in_channel"
}
func (k KeysCommandsSeventvErrorsEmoteAlreadyExistInChannel) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "emote_already_exist_in_channel"}
}
func (k KeysCommandsSeventvErrorsEmoteAlreadyExistInChannel) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsEmoteAlreadyExistInChannel) SetVars(vars KeysCommandsSeventvErrorsEmoteAlreadyExistInChannelVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsEmoteAlreadyExistInChannelVars] {
	k.Vars = twiri18n.Vars{"emoteName": vars.EmoteName}
	return k
}

type KeysCommandsSeventvErrors struct {
	EmotesetNotActive		KeysCommandsSeventvErrorsEmotesetNotActive
	EmotesetBroadcasterNotActive	KeysCommandsSeventvErrorsEmotesetBroadcasterNotActive
	ProfileNotFound			KeysCommandsSeventvErrorsProfileNotFound
	EmoteFailedToRemove		KeysCommandsSeventvErrorsEmoteFailedToRemove
	EmoteFailedToRename		KeysCommandsSeventvErrorsEmoteFailedToRename
	EmoteNotFound			KeysCommandsSeventvErrorsEmoteNotFound
	EmoteNotFoundInEmoteset		KeysCommandsSeventvErrorsEmoteNotFoundInEmoteset
	EmotesetNotFound		KeysCommandsSeventvErrorsEmotesetNotFound
	ProfileFailedToGet		KeysCommandsSeventvErrorsProfileFailedToGet
	EmoteFailedToAdd		KeysCommandsSeventvErrorsEmoteFailedToAdd
	EmoteFailedToFetch		KeysCommandsSeventvErrorsEmoteFailedToFetch
	EmoteNotFoundInChannel		KeysCommandsSeventvErrorsEmoteNotFoundInChannel
	EmoteAlreadyExistInChannel	KeysCommandsSeventvErrorsEmoteAlreadyExistInChannel
}
type KeysCommandsSeventvAddEmoteAddVars struct {
}
type KeysCommandsSeventvAddEmoteAdd struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvAddEmoteAdd) IsTranslationKey() {
}
func (k KeysCommandsSeventvAddEmoteAdd) GetPath() string {
	return "commands.7tv.add.emote_add"
}
func (k KeysCommandsSeventvAddEmoteAdd) GetPathSlice() []string {
	return []string{"commands", "7tv", "add", "emote_add"}
}
func (k KeysCommandsSeventvAddEmoteAdd) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvAddEmoteAdd) SetVars(vars KeysCommandsSeventvAddEmoteAddVars) twiri18n.TranslationKey[KeysCommandsSeventvAddEmoteAddVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvAdd struct {
	EmoteAdd KeysCommandsSeventvAddEmoteAdd
}
type KeysCommandsSeventvRemoveEmoteRemoveVars struct {
	EmoteName any
}
type KeysCommandsSeventvRemoveEmoteRemove struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvRemoveEmoteRemove) IsTranslationKey() {
}
func (k KeysCommandsSeventvRemoveEmoteRemove) GetPath() string {
	return "commands.7tv.remove.emote_remove"
}
func (k KeysCommandsSeventvRemoveEmoteRemove) GetPathSlice() []string {
	return []string{"commands", "7tv", "remove", "emote_remove"}
}
func (k KeysCommandsSeventvRemoveEmoteRemove) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvRemoveEmoteRemove) SetVars(vars KeysCommandsSeventvRemoveEmoteRemoveVars) twiri18n.TranslationKey[KeysCommandsSeventvRemoveEmoteRemoveVars] {
	k.Vars = twiri18n.Vars{"emoteName": vars.EmoteName}
	return k
}

type KeysCommandsSeventvRemove struct {
	EmoteRemove KeysCommandsSeventvRemoveEmoteRemove
}
type KeysCommandsSeventvRenameEmoteRenameVars struct {
	OldEmoteName	any
	NewEmoteName	any
}
type KeysCommandsSeventvRenameEmoteRename struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvRenameEmoteRename) IsTranslationKey() {
}
func (k KeysCommandsSeventvRenameEmoteRename) GetPath() string {
	return "commands.7tv.rename.emote_rename"
}
func (k KeysCommandsSeventvRenameEmoteRename) GetPathSlice() []string {
	return []string{"commands", "7tv", "rename", "emote_rename"}
}
func (k KeysCommandsSeventvRenameEmoteRename) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvRenameEmoteRename) SetVars(vars KeysCommandsSeventvRenameEmoteRenameVars) twiri18n.TranslationKey[KeysCommandsSeventvRenameEmoteRenameVars] {
	k.Vars = twiri18n.Vars{"oldEmoteName": vars.OldEmoteName, "newEmoteName": vars.NewEmoteName}
	return k
}

type KeysCommandsSeventvRename struct {
	EmoteRename KeysCommandsSeventvRenameEmoteRename
}
type KeysCommandsSeventvEmoteInfoResponseVars struct {
	Name		any
	Link		any
	AddedByUserName	any
	AddedByTime	any
	EmoteAuthor	any
}
type KeysCommandsSeventvEmoteInfoResponse struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvEmoteInfoResponse) IsTranslationKey() {
}
func (k KeysCommandsSeventvEmoteInfoResponse) GetPath() string {
	return "commands.7tv.emote_info.response"
}
func (k KeysCommandsSeventvEmoteInfoResponse) GetPathSlice() []string {
	return []string{"commands", "7tv", "emote_info", "response"}
}
func (k KeysCommandsSeventvEmoteInfoResponse) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvEmoteInfoResponse) SetVars(vars KeysCommandsSeventvEmoteInfoResponseVars) twiri18n.TranslationKey[KeysCommandsSeventvEmoteInfoResponseVars] {
	k.Vars = twiri18n.Vars{"name": vars.Name, "link": vars.Link, "addedByUserName": vars.AddedByUserName, "addedByTime": vars.AddedByTime, "emoteAuthor": vars.EmoteAuthor}
	return k
}

type KeysCommandsSeventvEmoteInfo struct {
	Response KeysCommandsSeventvEmoteInfoResponse
}
type KeysCommandsSeventv struct {
	ProfileInfo	KeysCommandsSeventvProfileInfo
	Hints		KeysCommandsSeventvHints
	Errors		KeysCommandsSeventvErrors
	Add		KeysCommandsSeventvAdd
	Remove		KeysCommandsSeventvRemove
	Rename		KeysCommandsSeventvRename
	EmoteInfo	KeysCommandsSeventvEmoteInfo
}
type KeysCommands struct {
	Shorturl		KeysCommandsShorturl
	Marker			KeysCommandsMarker
	Nuke			KeysCommandsNuke
	Stats			KeysCommandsStats
	Vips			KeysCommandsVips
	Channel			KeysCommandsChannel
	Clip			KeysCommandsClip
	Games			KeysCommandsGames
	CategoriesAliases	KeysCommandsCategoriesAliases
	Permit			KeysCommandsPermit
	Tts			KeysCommandsTts
	Subage			KeysCommandsSubage
	ChatWall		KeysCommandsChatWall
	Predictions		KeysCommandsPredictions
	Prefix			KeysCommandsPrefix
	Manage			KeysCommandsManage
	Songrequest		KeysCommandsSongrequest
	Dudes			KeysCommandsDudes
	Shoutout		KeysCommandsShoutout
	Seventv			KeysCommandsSeventv
}
type Keys struct {
	Errors		KeysErrors
	Services	KeysServices
	Variables	KeysVariables
	Commands	KeysCommands
}

var Translations = Keys{}
var Store twiri18n.LocalesStore = twiri18n.LocalesStore{"en": map[string]map[string]map[string]string{"errors": map[string]map[string]string{"generic": map[string]string{"cannot_find_users_twitch": `Cannot find users on twitch`, "cannot_create_command": `Cannot create command`, "cannot_get_command": `Cannot get command`, "cannot_get_db_channel": `Cannot get database channel`, "cannot_get_message": `Cannot get messages`, "cannot_find_channel_db": `Cannot find channel in database`, "twir_error": `[Twir error]: Unknown error happened. Please contact developers.`, "cannot_get_user": `Cannot get user: {reason}`, "not_a_follower": `Not a follower`, "cannot_get_accept_command_name": `Cannot get accept command name`, "internal": `❌ Internal error`, "getting_channel_settings": `Error while getting channel settings`, "cannot_timeout_user": `Cannot timeout user`, "user_not_found": `User not found`, "cannot_create_twitch": `Cannot create twitch client`, "create_settings": `Error while creating settings`, "cannot_get_moderators": `Cannot get moderators: {reason}`, "cannot_get_stream": `Cannot get stream. {reason}`, "getting_user_settings": `Error while getting user settings`, "broadcaster_client": `Cannot create broadcaster twitch client`, "updating_settings": `Error while updating settings`, "cannot_find_channel_twitch": `Cannot find channel on twitch. {reason}`, "should_mention_with_at": `You should tag user with @`, "something_went_wrong": `Something went wrong, please try again later`, "cannot_ban_user": `Cannot ban user`, "cannot_find_user_db": `Cannot find user in database`, "cannot_find_user_twitch": `Cannot find user on twitch`}}, "services": map[string]map[string]string{"chat_wall": map[string]string{"errors.get_users_stats": `Cannot get users stats: {reason}`, "errors.publish_ban_users": `Cannot publish ban users: {reason}`, "errors.create_chat_wall_with_phrase": `Cannot create chat wall with phrase that already exists`, "errors.handled_messages_to_redis": `Cannot add handled messages to redis: {reason}`, "errors.create_chat_wall": `Cannot create chat wall: {reason}`, "errors.update_chat_walls": `Cannot update chat wall: {reason}`, "errors.get_current_chat_walls": `Cannot get current chat walls: {reason}`, "info.banned_by_twir": `Banned by twir for chat wall phrase: {banPhrase}`, "errors.get_chat_walls": `Cannot get chat walls: {reason}`, "errors.publish_deleted_messages": `Cannot publish deleted messages: {reason}`, "errors.get_already_handled": `Cannot get already handled messages: {reason}`, "errors.chat_wall_not_found": `Chat wall not found`, "errors.create_chat_logs_in_db": `Cannot create chat wall logs in db: {reason}`, "errors.get_chat_wall_settings": `Cannot get chat wall settings: {reason}`}, "shortenedurls": map[string]string{"errors.invalid_url": `Invalid url`}, "tts": map[string]string{"errors.voice_disallowed": `Voice {userVoice} is disallowed for usage`, "info.not_configured": `TTS not configured`, "info.no_voices": `No voices available`, "errors.not_found": `Voice {userVoice} not found`}}, "variables": map[string]map[string]string{"custom_var": map[string]string{"errors.evaluate_variable": `Cannot evaluate variable. This is internal error, please report this bug`, "errors.wrong_numbers": `Wrong number: {reason}`, "errors.update_custom_var": `Cannot update custom variables`}, "shorturl": map[string]string{"errors.url_required": `URL is required`, "errors.create_short_url": `Cannot create short url: {reason}`}, "subscribers": map[string]string{"errors.get_subscribers": `Cannot get subscribers`}, "chat_eval": map[string]string{"info.wrong_code": `Probably you're doing some suspicious things or wrote wrong code.`}, "commands": map[string]string{"info.command_with_name_not_found": `Command with name {commandName} not found`, "info.get_count": `Cannot get count`, "info.no_passed_params": `Have not passed params to variable`}, "followers": map[string]string{"errors.get_followers": `Cannot get followers`}, "request": map[string]string{"errors.execute_request": `Cannot execute request`}, "7tv": map[string]string{"errors.profile_not_found": `Failed to get 7tv profile: {reason}`, "errors.emote_not_found": `Emote {name} not found`, "errors.no_paint": `No paint`, "errors.no_roles": `No roles`, "errors.no_active_set": `You don't have a 7TV set set`}, "countdown": map[string]string{"errors.not_passed_params": `Have not passed params to variable`, "errors.parse_date": `Cannot parse date`}, "random": map[string]string{"errors.lower_numbers": `Numbers cannot be lower then 0`, "errors.parse_first_number": `Cannot parse first number from arguments`, "errors.first_larger_second": `First number cannot be larger then second`, "errors.parse_second_number": `Cannot parse second number from arguments`, "errors.empty_phrase": `Your phrases contains empty phrase, check you writed commas correctly`, "errors.not_passed_params": `Have not passed params to random variable`, "errors.wrong_with_params": `Something is wrong with your params`, "errors.parameters_not_specified": `Parameters are not specified`, "errors.get_online_user": `Cannot get online user`, "errors.wrong_number": `Wrong number of arguments passed to random`}, "stream": map[string]string{"errors.error": `Error`, "errors.offline": `offline`, "errors.count_followers": `Cannot count followers`, "errors.get_history_of_categories": `Cannot get history of categories`, "info.offline": `Offline or error on getting category`, "info.no_history": `No history recorded`}, "user": map[string]string{"errors.find_user_on_twitch": `Cannot find user on twitch`}, "valorant": map[string]string{"info.matches": `{matchResult}({roundsWon}/{roundsLost}}) — {char} {KDA}`}, "keywords": map[string]string{"errors.id_not_provided": `ID is not provided`, "errors.not_found": `Keyword not found`}, "song": map[string]string{"info.failed_get_spotify_integration": `Failed to get spotify integration`, "info.history": `{trackTitle} - {trackArtist} (~{minutes}m ago)`, "errors.fetch_tracks_lastfm": `Cannot fetch tracks from lastfm: {reason}`, "errors.get_recent_tracks": `Cannot get recent tracks: {reason}`, "info.spotify_not_connected": `Spotify not connected`, "info.get_spotify_entity": `Failed to get spotify entity`, "info.get_spotify_integration": `Cannot get spotify integration: {reason}`, "info.no_needed_scope": `No needed scope, reconnect spotify in dashboard`, "info.no_integrations": `No integrations connected`, "errors.create_lastfm_service": `Cannot create lastfm service: {reason}`, "errors.fetch_tracks_spotify": `Cannot fetch tracks from spotify: {reason}`, "errors.parse_played_at": `Cannot parse played at`, "info.lastfm_integration": `LastFM integration not enabled`}}, "commands": map[string]map[string]string{"clip": map[string]string{"cannot_get_clip": `Cannot get created clip`, "clip_created": `Clip created: {url}`, "cannot_create_clip": `Cannot create clip`, "empty_clip_url": `Clip URL is empty, please try again`}, "manage": map[string]string{"add.alias_add": `✅ Alias added`, "remove.alias_removed": `✅ Alias removed`, "errors.alias_already_exist": `Command with {alias} name or alias already exists`, "errors.command_cannot_get": `Cannot get command`, "errors.command_have_no_aliases": `Command have no aliases`, "errors.alias_not_command": `That alias not in the command`, "add.command_add": `✅ Command added`, "edit.command_edited": `✅ Command edited`, "errors.command_not_found": `Command not found`, "errors.alias_cannot_update": `Cannot update command aliases`, "errors.command_with_name_cannot_find": `Command with that name not found`, "errors.command_with_alias_already_exists": `Command with that name or alias already exists`, "errors.command_cannot_delete_default": `Cannot delete default command`, "remove.command_removed": `✅ Command removed`, "errors.command_cannot_update_response": `Cannot update response because you have more than 1 response in command. Please use UI.`, "errors.command_cannot_update": `Cannot update command`, "errors.command_large_size": `Command name cannot be greatest then 20.`, "errors.command_cannot_save": `Cannot save command`, "errors.alias_cannot_get_existed_commands": `Cannot get existed commands`}, "tts": map[string]string{"info.change_rate": `Rate changed to {newRate}`, "errors.while_getting_voices": `Error while getting voices`, "errors.while_disable": `Error while disabling tts`, "info.voice": `Global voice: {globalVoice} | Your voice: {userVoice}`, "info.current_volume": `Current volume: {ttsVolume}`, "info.enabled": `TTS enabled`, "info.rate": `Global rate: {globalRate} | Your rate: {userRate}`, "info.no_voices": `No voices available`, "info.change_volume": `TTS volume changed to {userVolume}`, "errors.while_enable": `Error while enabling tts`, "errors.not_configured": `TTS is not configured for this channel`, "info.change_pitch": `Pitch changed to {newPitch}`, "info.pitch": `Global pitch: {globalPitch} | Your pitch: {userPitch}`, "info.voice_disallowed": `Voice {voiceName} is disallowed for usage`, "errors.sending_to_tts": `Error while sending message to tts service`, "info.change_voice": `Voice changed to {newVoice}`, "info.disabled": `TTS disabled`}, "7tv": map[string]string{"hints.copySetNameOfSet": `Name of set to copy`, "hints.emoteForAddArgLink": `Link or name`, "hints.emoteForAddArgAlias": `Optional alias`, "errors.profile_failed_to_get": `Failed to get 7tv profile: {reason}`, "errors.emoteset_not_found": `Emote set {emoteName} not found`, "errors.emoteset_broadcaster_not_active": `❌ No active emote set for broadcaster`, "errors.emote_already_exist_in_channel": `❌ Emote {emoteName} already exists in this channel`, "errors.emote_failed_to_fetch": `Failed to fetch 7tv emote: {reason}`, "remove.emote_remove": `✅ Emote {emoteName} removed`, "rename.emote_rename": `✅ Emote {oldEmoteName} renamed to {newEmoteName}"`, "hints.copySetChannelName": `@channelName`, "errors.emote_not_found_in_emoteset": `Emote {emoteName} not found in set {emoteSet}`, "errors.emote_not_found": `Emote {emoteName} not found.`, "hints.emoteForCopyArgName": `Name of emote to copy`, "hints.emoteForCopyAlias": `Alias for emote`, "errors.emote_failed_to_add": `Failed to add 7tv emote: {reason}`, "errors.emote_not_found_in_channel": `❌ Emote {emoteSearch} not found in target channel`, "add.emote_add": `✅ Emote added`, "emote_info.response": `{name}: {link} · Added by @{addedByUserName} {addedByTime} ago · Author {emoteAuthor}`, "profile_info.response": `{profileName} · Paint: {paintName} ({unlockedPaints} unlocked) · Roles: {roles} · Editor for {editorCount} · Set: {emoteSetName} ({emoteSetCount}/{emoteSetCapacity}) · Created: {profileCreatedAt}`, "errors.profile_not_found": `7tv profile not found`, "errors.emote_failed_to_remove": `Failed to remove 7tv emote: {reason}`, "errors.emote_failed_to_rename": `Failed to rename 7tv emote: {reason}`, "errors.emoteset_not_active": `You don't have an active 7TV emote set`}, "games": map[string]string{"errors.duel_cannot_check_user": `Cannot check user in duel: {reason}`, "errors.duel_with_yourself": `You cannot duel with yourself`, "errors.duel_cannot_save_data": `Cannot save duel data: {reason}`, "errors.duel_cannot_validate_participants": `Cannot validate participants`, "errors.duel_with_bot": `You cannot duel with bot`, "info.voteban_in_progress": `Another voteban in progress`, "info.user_not_participate": `You are not participate in any duel`, "errors.roulette_cannot_get_with_settings": `Cannot get roulette settings from database`, "errors.duel_with_streamer": `You cannot duel with streamer`, "errors.8ball_cannot_find": `Cannot find 8ball settings`, "errors.voteban_cannot_find_user": `Cannot find target user`, "errors.voteban_cannot_set_vote": `Cannot set vote`, "errors.voteban_cannot_check_progress": `Cannot check if vote in progress`, "errors.voteban_cannot_set_vote_expiration": `Cannot set vote expiration`, "errors.duel_cannot_save_to_cache": `Cannot save duel data to cache`, "info.duel_stats": `You have shoot {duels} times · {wins} W – {loses} L ({winrate} WR)`, "errors.duel_cannot_set_user_cooldown": `Cannot set user cooldown: {reason}`, "errors.seppuku_cannot_find_settings": `Cannot find seppuku settings`, "errors.voteban_cannot_find_settings": `Cannot find voteban settings`, "errors.duel_cannot_check_cooldown": `Cannot check duel cooldown: {reason}`, "info.sender_already_in_duel": `You already in duel`, "errors.duel_cannot_get_sender": `Cannot get sender current duel`, "errors.roulette_cannot_send_death_message": `Cannot send death message`, "errors.duel_cannot_set_global_cooldown": `Cannot set global cooldown: {reason}`, "errors.voteban_cannot_lock": `Cannot lock voteban`, "info.user_already_in_duel": `Target user already in duel`, "errors.duel_cannot_get_with_settings": `Cannot get duel channel settings`, "errors.roulette_cannot_send_initial_message": `Cannot send initial message`, "errors.duel_cannot_save_result": `Cannot save duel result: {reason}`}, "permit": map[string]string{"errors.cannot_create": `Cannot create permit`, "success.added_permit": `✅ Added {countPermit} permits to {userName}`}, "predictions": map[string]string{"errors.cannot_create": `Cannot create prediction`, "info.started": `✅ Prediction started`, "errors.cannot_get_current": `Cannot get current prediction`, "info.cancel": `✅ Prediction canceled`, "hints.startPredictionDuration": `120`, "hints.startPredictionArgTitle": `Will we win this game?`, "info.resolved": `✅ Prediction resolved`, "hints.predictionResolveOutcomeNum": `Variant number, for example: 1,2,3,4,5`, "errors.cannot_create_var": `Cannot create prediction: {reason}`, "errors.cannot_get_current_var": `Cannot get current prediction: {reason}`, "errors.cannot_cancel": `Cannot cancel prediction`, "errors.cannot_cancel_var": `Cannot cancel prediction: {reason}`, "errors.no_variant": `No prediction variant`, "info.no_runed": `No prediction runed`, "info.locked": `✅ Prediction locked`, "hints.startPredictionArgVariants": `Yes, win / No, lose`}, "prefix": map[string]string{"errors.cannot_update": `Cannot update prefix`, "errors.required": `Prefix is required`, "errors.too_long": `Prefix cannot be longer than 10 characters`, "errors.cannot_get_current": `Cannot get current prefix`, "errors.cannot_create": `Cannot create prefix`, "success.updated": `Prefix updated`}, "shoutout": map[string]string{"errors.bot_have_no_permissions": `We have no permissions for shoutout. Streamer must re-authorize to bot dashboard.`, "response_online": `Join https://twitch.tv/{userName} shining in 🎮 {categoryName} 📜 {title} with 👁️ {viewers} viewers!`, "response_offline": `Explore https://twitch.tv/{userName}’s world, last featured in 🎮 {categoryName} 📜 {title}!`}, "stats": map[string]string{"info.watching_stream": `You watching stream for {userWatching}`, "me.songs": `songs requests`, "me.watched": `watched`, "me.messages": `messages`, "me.emotes": `used emotes`, "me.points": `used points`}, "vips": map[string]string{"errors.no_scheduled_vips": `There are no scheduled vips`, "errors.cannot_create_scheduled_in_db": `Cannot create scheduled vip in database`, "errors.invalid_duration": `Invalid duration format. Please use formats like <1h>, <30m>, or <2d>`, "errors.added": `✅ Added vip to {userName}`, "errors.cannot_get_list_from_db": `Cannot get vip list from database`, "hints.user": `@username`, "errors.removed": `✅ Removed vip from {userName}`, "errors.updated": `✅ Updated vip for user {userName} new expriation time {endTime}`, "errors.cannot_update": `Cannot update scheduled vip`, "hints.unvip_in": `Time in, example: 1w5d1m5s`, "errors.already_have_role": `User already vip or moderator!`, "errors.added_with_remove_time": `✅ Added vip to {userName}, will be removed at {endTime}`}, "categories_aliases": map[string]string{"errors.category_not_found": `Category not found`, "add.alias_add_to_category": `Category alias {aliasName} added with category {categoryName}`, "errors.alias_removed": `Category alias {aliasName} removed`, "errors.twitch_client_cannot_to_create": `Cannot create twitch client`, "errors.category_cannot_to_get": `Cannot get categories`, "errors.alias_already_exists": `Alias {aliasName} already exists`, "errors.game_cannot_to_get": `Cannot get games`, "errors.category_required": `Alias and category are required`, "errors.category_failed_to_create": `Failed to create category`, "errors.category_failed_to_get": `Failed to get caterogies`, "errors.alias_empty": `No categories aliases created`, "errors.category_cannot_delete": `Cannot delete category`, "errors.alias_not_found": `Category alias not found`}, "chat_wall": map[string]string{"errors.chat_wall_not_found": `Chat wall {errorPhrase} not found or already stopped`, "errors.invalid_duration": `Invalid duration. Cannot be longer 2w Examples: 10m, 10, 1h5m`, "stop.chat_wal_stop": `✅ Chat wall {chatWallPhrase} stopped`, "hints.timeoutPhraseArgName": `Phrase to ban`, "hints.banPhraseArgName": `Phrase to ban`, "hints.deletePhraseArgName": `Phrase to delete`, "errors.long_duration_timeout": `Duration of timeout cannot be longer than 2 weeks`, "errors.duration_cannot_parse": `Cannot parse duration`, "start.chat_wall_start": `✅ Chat wall started for 10 minutes, you can stop it with !chat wall stop {chatWallPhrase}`, "hints.timeoutDurationArgName": `Time. Examples: 10m, 10, 1h5m`}, "dudes": map[string]string{"errors.sprite_invalid": `Invalid sprite, available: {availableSprites}`, "errors.color_cannot_trigger": `Cannot trigger dudes color`, "info.sprite": `Your sprite is {dudeSprite}`, "info.color_reset": `Color reset to default`, "info.color_changed": `Color changed to {dudeColor}`, "errors.grow_cannot_trigger": `Cannot trigger dudes grow`, "errors.sprite_cannot_trigger": `Cannot trigger dudes sprite`, "errors.color_invalid": `Invalid color`, "info.color_required": `Color is required`, "info.sprite_required": `Sprite is required, available: {availableSprites}`, "info.sprite_changed": `Sprite changed to {dudeSprite}`, "info.color": `Your color is {dudeColor}`, "errors.jump_cannot_trigger": `Cannot trigger dudes jump`, "errors.leave_cannot_trigger": `Cannot trigger dudes leave`}, "nuke": map[string]string{"hints.nukeTimeArgName": `time, examples: 10m, 10, 1h5m`, "errors.parse_duration": `Cannot parse duration`, "errors.invalid_duration": `Invalid duration. Examples: !nuke 10m phrase, !nuke 10 phrase, !nuke 1h5m phrase`, "errors.cannot_get_users_stats": `Cannot get users stats`, "errors.cannot_get_handeled_messages": `Cannot get handled messages`, "errors.cannot_delete_messages": `Cannot delete messages`, "errors.timeout_duration": `Duration of timeout cannot be longer than 2 weeks`}, "songrequest": map[string]string{"errors.get_votes_count": `Cannot get votes count`, "errors.get_current_song": `Cannot get current song`, "info.delete": `Song {songTitle} deleted from queue`, "info.no_requested_songs": `You haven't requested any song`, "info.only_count_songs": `There is only {songsCount} songs`, "validate.errors.restrictions_on_user": `There are restrictions on user, but i cannot find you in db, sorry :(`, "errors.get_current_queue_count": `Cannot get current queue count`, "errors.not_found": `Current song not found`, "validate.errors.internal_error": `Internal error when checking follow`, "errors.get_latest_song": `Cannot get latest song`, "errors.remove_song_from_queue": `Cannot remove song from queue`, "errors.get_settings": `Cannot get song requests settings`, "validate.errors.need_follow": `For request song you need to be a followed`, "errors.get_current_vote": `Cannot get current vote`, "info.song_skipped": `Song {songTitle} skipped`, "errors.search_song": `Cannot search song`, "errors.update_song": `Cannot update song`, "errors.get_songs_from_queue": `Cannot get songs from queue`, "errors.get_users_count": `Cannot get online users count`}, "subage": map[string]string{"responses.streak_info": `(streak: {months} months)`, "responses.time_remaining": `(time remaining: {duration})`, "responses.not_subscriber": `{user} is not a subscriber`, "responses.not_subscriber_but_was": `{user} is not a subscriber but was subscribed for {months} months`, "errors.not_subscriber_or_hidden": `User is not a subscriber or their subscription info is hidden`, "responses.subscription_info": `{user} is a {tier} subscriber of {channel} for {months} months`}, "channel": map[string]string{"errors.channel_cannot_get_information": `Cannot get channel information`, "errors.category_cannot_get": `Cannot get category`, "errors.category_cannot_get_error": `Cannot get category: {errorMessage}`, "errors.category_cannot_change": `Cannot change category`, "errors.category_cannot_change_error": `Cannot change category: {errorMessage}`, "errors.broadcaster_twitch_client_cannot_create": `Cannot create broadcaster twitch client`, "errors.alias_cannot_get_category": `Cannot get category aliases`, "add.category_change": `✅ {categoryName}`, "errors.history_game_message": `Cannot find used games in database`, "errors.category_not_found": `Category not found`, "errors.channel_not_found": `Channel not found`, "errors.broadcaster_twitch_api_client": `Cannot create broadcaster twitch api client: {reason}`, "errors.history_title_message": `Cannot get history of titles from database: {reason}`, "errors.game_not_found": `Game not found on twitch`, "hints.gameArgName": `Category name or created category alias`}, "marker": map[string]string{"errors.cannot_create_marker": `Cannot create marker. {reason}`, "success.marker_created": `Marker created`}, "shorturl": map[string]string{"errors.cannot_create_short_url": `Cannot create short url. {error}`, "success.short_url_created": `Short url: {url}`}}}}
