package postgres

import (
	"github.com/twirapp/twir/libs/repositories/channels_moderation_settings"
)

func makeCreateOrUpdateMap(input channels_moderation_settings.CreateOrUpdateInput) map[string]any {
	return map[string]any{
		"name":                                 input.Name,
		"type":                                 input.Type,
		"channel_id":                           input.ChannelID,
		"enabled":                              input.Enabled,
		"ban_time":                             input.BanTime,
		"ban_message":                          input.BanMessage,
		"warning_message":                      input.WarningMessage,
		"check_clips":                          input.CheckClips,
		"trigger_length":                       input.TriggerLength,
		"max_percentage":                       input.MaxPercentage,
		"deny_list":                            input.DenyList,
		"denied_chat_languages":                input.DeniedChatLanguages,
		"excluded_roles":                       input.ExcludedRoles,
		"max_warnings":                         input.MaxWarnings,
		"deny_list_regexp_enabled":             input.DenyListRegexpEnabled,
		"deny_list_word_boundary_enabled":      input.DenyListWordBoundaryEnabled,
		"deny_list_sensitivity_enabled":        input.DenyListSensitivityEnabled,
		"one_man_spam_minimum_stored_messages": input.OneManSpamMinimumStoredMessages,
		"one_man_spam_message_memory_seconds":  input.OneManSpamMessageMemorySeconds,
	}
}
