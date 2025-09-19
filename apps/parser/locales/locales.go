package locales

import twiri18n "github.com/twirapp/twir/libs/i18n"

type KeysCommandsFollowageDescriptionVars struct {
}
type KeysCommandsFollowageDescription struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsFollowageDescription) IsTranslationKey() {
}
func (k KeysCommandsFollowageDescription) GetPath() string {
	return "commands.followage.description"
}
func (k KeysCommandsFollowageDescription) GetPathSlice() []string {
	return []string{"commands", "followage", "description"}
}
func (k KeysCommandsFollowageDescription) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsFollowageDescription) SetVars(vars KeysCommandsFollowageDescriptionVars) twiri18n.TranslationKey[KeysCommandsFollowageDescriptionVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsFollowageResponseVars struct {
	User     any
	Duration any
}
type KeysCommandsFollowageResponse struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsFollowageResponse) IsTranslationKey() {
}
func (k KeysCommandsFollowageResponse) GetPath() string {
	return "commands.followage.response"
}
func (k KeysCommandsFollowageResponse) GetPathSlice() []string {
	return []string{"commands", "followage", "response"}
}
func (k KeysCommandsFollowageResponse) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsFollowageResponse) SetVars(vars KeysCommandsFollowageResponseVars) twiri18n.TranslationKey[KeysCommandsFollowageResponseVars] {
	k.Vars = twiri18n.Vars{"user": vars.User, "duration": vars.Duration}
	return k
}

type KeysCommandsFollowage struct {
	Description KeysCommandsFollowageDescription
	Response    KeysCommandsFollowageResponse
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

type KeysCommandsVipsAddedWithRemoveTimeVars struct {
	UserName any
	EndTime  any
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

type KeysCommandsVipsUpdatedVars struct {
	UserName any
	EndTime  any
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

type KeysCommandsVips struct {
	Added                     KeysCommandsVipsAdded
	AddedWithRemoveTime       KeysCommandsVipsAddedWithRemoveTime
	AlreadyHaveRole           KeysCommandsVipsAlreadyHaveRole
	CannotCreateScheduledInDb KeysCommandsVipsCannotCreateScheduledInDb
	InvalidDuration           KeysCommandsVipsInvalidDuration
	NoScheduledVips           KeysCommandsVipsNoScheduledVips
	Updated                   KeysCommandsVipsUpdated
	CannotGetListFromDb       KeysCommandsVipsCannotGetListFromDb
	CannotUpdate              KeysCommandsVipsCannotUpdate
	Removed                   KeysCommandsVipsRemoved
}
type KeysCommands struct {
	Followage KeysCommandsFollowage
	Vips      KeysCommandsVips
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

type KeysErrorsGeneric struct {
	BroadcasterClient     KeysErrorsGenericBroadcasterClient
	CannotFindUserDb      KeysErrorsGenericCannotFindUserDb
	CannotFindUserTwitch  KeysErrorsGenericCannotFindUserTwitch
	CannotFindUsersTwitch KeysErrorsGenericCannotFindUsersTwitch
	ShouldMentionWithAt   KeysErrorsGenericShouldMentionWithAt
}
type KeysErrors struct {
	Generic KeysErrorsGeneric
}
type Keys struct {
	Commands KeysCommands
	Errors   KeysErrors
}

var Translations = Keys{}
var Store twiri18n.LocalesStore = twiri18n.LocalesStore{
	"en": map[string]map[string]map[string]string{
		"commands": map[string]map[string]string{
			"followage": map[string]string{
				"response":    `User {user} followed channel for {duration}`,
				"description": `Check how long a user has been following a Twitch channel.`,
			},
			"vips": map[string]string{
				"added":                         `✅ added vip to {userName}`,
				"no_scheduled_vips":             `There are no scheduled vips.`,
				"removed":                       `✅ removed vip from {userName}`,
				"cannot_update":                 `Cannot update scheduled vip.`,
				"invalid_duration":              `Invalid duration format. Please use formats like "1h", "30m", or "2d".`,
				"already_have_role":             `User already vip or moderator!`,
				"cannot_get_list_from_db":       `Cannot get vip list from database.`,
				"updated":                       `✅ updated vip for user {userName} new expriation time {endTime}`,
				"cannot_create_scheduled_in_db": `Cannot create scheduled vip in database.`,
				"added_with_remove_time":        `✅ added vip to {userName}, will be removed at {endTime}`,
			},
		},
		"errors": map[string]map[string]string{
			"generic": map[string]string{
				"broadcaster_client":       `Cannot create broadcaster twitch client`,
				"should_mention_with_at":   `you should tag user with @`,
				"cannot_find_user_db":      `Cannot find user in database`,
				"cannot_find_user_twitch":  `Cannot find user on twitch`,
				"cannot_find_users_twitch": `Cannot find users on twitch`,
			},
		},
	},
	"ru": map[string]map[string]map[string]string{
		"commands": map[string]map[string]string{
			"followage": map[string]string{
				"response":    `Пользователь {user} подписан на канал {duration}`,
				"description": `Показывает, как долго пользователь подписан на канал.`,
			},
			"vips": map[string]string{"invalid_duration": `Неверный формат длительности. Пожалуйста, используйте форматы, такие как "1ч", "30м" или "2д".`},
		},
	},
}
