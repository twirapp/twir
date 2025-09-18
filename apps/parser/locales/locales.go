package locales

import twiri18n "github.com/twirapp/twir/libs/i18n"

type KeysCommandsVipsInvalid_duration struct {
	twiri18n.TranslationKey
	Value	string
}

func (k KeysCommandsVipsInvalid_duration) IsTranslationKey() {
}
func (k KeysCommandsVipsInvalid_duration) GetPath() string {
	return "commands.vips.invalid_duration"
}
func (k KeysCommandsVipsInvalid_duration) GetPathSlice() []string {
	return []string{"commands", "vips", "invalid_duration"}
}

type KeysCommandsVips struct {
	twiri18n.TranslationKey
	Invalid_duration	KeysCommandsVipsInvalid_duration
}

func (k KeysCommandsVips) IsTranslationKey() {
}
func (k KeysCommandsVips) GetPath() string {
	return "commands.vips"
}
func (k KeysCommandsVips) GetPathSlice() []string {
	return []string{"commands", "vips"}
}

type KeysCommandsFollowageDescription struct {
	twiri18n.TranslationKey
	Value	string
}

func (k KeysCommandsFollowageDescription) IsTranslationKey() {
}
func (k KeysCommandsFollowageDescription) GetPath() string {
	return "commands.followage.description"
}
func (k KeysCommandsFollowageDescription) GetPathSlice() []string {
	return []string{"commands", "followage", "description"}
}

type KeysCommandsFollowageResponse struct {
	twiri18n.TranslationKey
	Value	string
}

func (k KeysCommandsFollowageResponse) IsTranslationKey() {
}
func (k KeysCommandsFollowageResponse) GetPath() string {
	return "commands.followage.response"
}
func (k KeysCommandsFollowageResponse) GetPathSlice() []string {
	return []string{"commands", "followage", "response"}
}

type KeysCommandsFollowage struct {
	twiri18n.TranslationKey
	Description	KeysCommandsFollowageDescription
	Response	KeysCommandsFollowageResponse
}

func (k KeysCommandsFollowage) IsTranslationKey() {
}
func (k KeysCommandsFollowage) GetPath() string {
	return "commands.followage"
}
func (k KeysCommandsFollowage) GetPathSlice() []string {
	return []string{"commands", "followage"}
}

type KeysCommands struct {
	twiri18n.TranslationKey
	Vips		KeysCommandsVips
	Followage	KeysCommandsFollowage
}

func (k KeysCommands) IsTranslationKey() {
}
func (k KeysCommands) GetPath() string {
	return "commands"
}
func (k KeysCommands) GetPathSlice() []string {
	return []string{"commands"}
}

type Keys struct {
	twiri18n.TranslationKey
	Commands	KeysCommands
}

func (k Keys) IsTranslationKey() {
}
func (k Keys) GetPath() string {
	return ""
}
func (k Keys) GetPathSlice() []string {
	return []string{""}
}

var Translations = Keys{}
var Store twiri18n.LocalesStore = twiri18n.LocalesStore{"en": map[string]map[string]map[string]string{"commands": map[string]map[string]string{"followage": map[string]string{"description": "Check how long a user has been following a Twitch channel.", "response": "User {user} followed channel for {duration}"}, "vips": map[string]string{"invalid_duration": "Invalid duration format. Please use formats like "1h", "30m", or "2d"."}}}, "ru": map[string]map[string]map[string]string{"commands": map[string]map[string]string{"followage": map[string]string{"description": "Показывает, как долго пользователь подписан на канал.", "response": "Пользователь {user} подписан на канал {duration}"}}}}
