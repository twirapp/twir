package locales

import twiri18n "github.com/twirapp/twir/libs/i18n"

type KeysCommandsSeventvErrorsNoActiveSetVars struct {
}
type KeysCommandsSeventvErrorsNoActiveSet struct {
	Vars twiri18n.Vars
}

func (k KeysCommandsSeventvErrorsNoActiveSet) IsTranslationKey() {
}
func (k KeysCommandsSeventvErrorsNoActiveSet) GetPath() string {
	return "commands.7tv.errors.no_active_set"
}
func (k KeysCommandsSeventvErrorsNoActiveSet) GetPathSlice() []string {
	return []string{"commands", "7tv", "errors", "no_active_set"}
}
func (k KeysCommandsSeventvErrorsNoActiveSet) GetVars() twiri18n.Vars {
	return k.Vars
}
func (k KeysCommandsSeventvErrorsNoActiveSet) SetVars(vars KeysCommandsSeventvErrorsNoActiveSetVars) twiri18n.TranslationKey[KeysCommandsSeventvErrorsNoActiveSetVars] {
	k.Vars = twiri18n.Vars{}
	return k
}

type KeysCommandsSeventvErrorsProfileNotFoundVars struct {
	Reason any
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
	k.Vars = twiri18n.Vars{"reason": vars.Reason}
	return k
}

type KeysCommandsSeventvErrorsEmoteNotFoundVars struct {
	Name any
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
	k.Vars = twiri18n.Vars{"name": vars.Name}
	return k
}

type KeysCommandsSeventvErrors struct {
	NoActiveSet	KeysCommandsSeventvErrorsNoActiveSet
	ProfileNotFound	KeysCommandsSeventvErrorsProfileNotFound
	EmoteNotFound	KeysCommandsSeventvErrorsEmoteNotFound
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
	Errors		KeysCommandsSeventvErrors
	EmoteInfo	KeysCommandsSeventvEmoteInfo
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

type KeysCommandsClip struct {
	ClipCreated		KeysCommandsClipClipCreated
	CannotCreateClip	KeysCommandsClipCannotCreateClip
	EmptyClipUrl		KeysCommandsClipEmptyClipUrl
	CannotGetClip		KeysCommandsClipCannotGetClip
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

type KeysCommandsSubageResponses struct {
	NotSubscriberButWas	KeysCommandsSubageResponsesNotSubscriberButWas
	SubscriptionInfo	KeysCommandsSubageResponsesSubscriptionInfo
	StreakInfo		KeysCommandsSubageResponsesStreakInfo
	TimeRemaining		KeysCommandsSubageResponsesTimeRemaining
	NotSubscriber		KeysCommandsSubageResponsesNotSubscriber
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
type KeysCommandsSubage struct {
	Responses	KeysCommandsSubageResponses
	Errors		KeysCommandsSubageErrors
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
type KeysCommandsMarker struct {
	Success	KeysCommandsMarkerSuccess
	Errors	KeysCommandsMarkerErrors
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

type KeysCommandsSongrequestErrors struct {
	SearchSong		KeysCommandsSongrequestErrorsSearchSong
	GetLatestSong		KeysCommandsSongrequestErrorsGetLatestSong
	GetCurrentQueueCount	KeysCommandsSongrequestErrorsGetCurrentQueueCount
	GetSongsFromQueue	KeysCommandsSongrequestErrorsGetSongsFromQueue
	UpdateSong		KeysCommandsSongrequestErrorsUpdateSong
	RemoveSongFromQueue	KeysCommandsSongrequestErrorsRemoveSongFromQueue
	GetSettings		KeysCommandsSongrequestErrorsGetSettings
}
type KeysCommandsSongrequest struct {
	Errors KeysCommandsSongrequestErrors
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

type KeysCommandsVips struct {
	Added				KeysCommandsVipsAdded
	NoScheduledVips			KeysCommandsVipsNoScheduledVips
	CannotUpdate			KeysCommandsVipsCannotUpdate
	Updated				KeysCommandsVipsUpdated
	CannotCreateScheduledInDb	KeysCommandsVipsCannotCreateScheduledInDb
	AddedWithRemoveTime		KeysCommandsVipsAddedWithRemoveTime
	CannotGetListFromDb		KeysCommandsVipsCannotGetListFromDb
	Removed				KeysCommandsVipsRemoved
	InvalidDuration			KeysCommandsVipsInvalidDuration
	AlreadyHaveRole			KeysCommandsVipsAlreadyHaveRole
}
type KeysCommands struct {
	Seventv		KeysCommandsSeventv
	Clip		KeysCommandsClip
	Prefix		KeysCommandsPrefix
	Shoutout	KeysCommandsShoutout
	Subage		KeysCommandsSubage
	Marker		KeysCommandsMarker
	Shorturl	KeysCommandsShorturl
	Songrequest	KeysCommandsSongrequest
	Vips		KeysCommandsVips
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

type KeysErrorsGeneric struct {
	BroadcasterClient	KeysErrorsGenericBroadcasterClient
	CannotFindUserDb	KeysErrorsGenericCannotFindUserDb
	CannotFindUserTwitch	KeysErrorsGenericCannotFindUserTwitch
	CannotFindChannelDb	KeysErrorsGenericCannotFindChannelDb
	CannotGetStream		KeysErrorsGenericCannotGetStream
	ShouldMentionWithAt	KeysErrorsGenericShouldMentionWithAt
	CannotFindUsersTwitch	KeysErrorsGenericCannotFindUsersTwitch
	CannotFindChannelTwitch	KeysErrorsGenericCannotFindChannelTwitch
}
type KeysErrors struct {
	Generic KeysErrorsGeneric
}
type Keys struct {
	Commands	KeysCommands
	Errors		KeysErrors
}

var Translations = Keys{}
var Store twiri18n.LocalesStore = twiri18n.LocalesStore{"en": map[string]map[string]map[string]string{"errors": map[string]map[string]string{"generic": map[string]string{"cannot_find_user_db": `Cannot find user in database`, "cannot_find_user_twitch": `Cannot find user on twitch`, "cannot_find_users_twitch": `Cannot find users on twitch`, "cannot_find_channel_db": `Cannot find channel in database`, "cannot_find_channel_twitch": `Cannot find channel on twitch. {reason}`, "cannot_get_stream": `Cannot get stream. {reason}`, "broadcaster_client": `Cannot create broadcaster twitch client`, "should_mention_with_at": `you should tag user with @`}}, "commands": map[string]map[string]string{"prefix": map[string]string{"errors.cannot_update": `Cannot update prefix.`, "success.updated": `Prefix updated.`, "errors.required": `Prefix is required.`, "errors.too_long": `Prefix cannot be longer than 10 characters.`, "errors.cannot_get_current": `Cannot get current prefix.`, "errors.cannot_create": `Cannot create prefix.`}, "songrequest": map[string]string{"errors.get_songs_from_queue": `Cannot get songs from queue.`, "errors.update_song": `Cannot update song.`, "errors.remove_song_from_queue": `Cannot remove song from queue.`, "errors.get_settings": `Cannot get song requests settings.`, "errors.search_song": `Cannot search song.`, "errors.get_latest_song": `Cannot get latest song.`, "errors.get_current_queue_count": `Cannot get current queue count.`}, "subage": map[string]string{"responses.not_subscriber_but_was": `{user} is not a subscriber but was subscribed for {months} months.`, "responses.subscription_info": `{user} is a {tier} subscriber of {channel} for {months} months.`, "responses.streak_info": `(streak: {months} months)`, "responses.time_remaining": `(time remaining: {duration})`, "errors.not_subscriber_or_hidden": `User is not a subscriber or their subscription info is hidden.`, "responses.not_subscriber": `{user} is not a subscriber.`}, "7tv": map[string]string{"errors.emote_not_found": `Emote {name} not found.`, "emote_info.response": `{name}: {link} · Added by @{addedByUserName} {addedByTime} ago · Author {emoteAuthor}`, "errors.no_active_set": `You don't have a 7TV set set.`, "errors.profile_not_found": `Failed to get 7tv profile: {reason}`}, "clip": map[string]string{"empty_clip_url": `Clip URL is empty, please try again.`, "cannot_get_clip": `Cannot get created clip.`, "clip_created": `Clip created: {url}`, "cannot_create_clip": `Cannot create clip.`}, "marker": map[string]string{"errors.cannot_create_marker": `Cannot create marker. {reason}`, "success.marker_created": `Marker created.`}, "shorturl": map[string]string{"errors.cannot_create_short_url": `Cannot create short url. {error}`, "success.short_url_created": `Short url: {url}`}, "shoutout": map[string]string{"response_online": `Join {userName} shining in 🎮 {categoryName} 📜 {title} with 👁️ {viewers} viewers!`, "response_offline": `Explore {userName}’s world, last featured in 🎮 {categoryName} 📜 {title}!`, "errors.bot_have_no_permissions": `we have no permissions for shoutout. Streamer must re-authorize to bot dashboard.`}, "vips": map[string]string{"no_scheduled_vips": `There are no scheduled vips.`, "already_have_role": `User already vip or moderator!`, "cannot_create_scheduled_in_db": `Cannot create scheduled vip in database.`, "added": `✅ added vip to {userName}`, "removed": `✅ removed vip from {userName}`, "cannot_update": `Cannot update scheduled vip.`, "invalid_duration": `Invalid duration format. Please use formats like "1h", "30m", or "2d".`, "cannot_get_list_from_db": `Cannot get vip list from database.`, "added_with_remove_time": `✅ added vip to {userName}, will be removed at {endTime}`, "updated": `✅ updated vip for user {userName} new expriation time {endTime}`}}}}
