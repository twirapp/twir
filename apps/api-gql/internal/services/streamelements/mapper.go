package streamelements

import (
	"strings"

	"github.com/twirapp/twir/apps/api-gql/internal/services/importer"
	streamelementsintegration "github.com/twirapp/twir/libs/integrations/streamelements"
)

func NormalizeCommands(input []streamelementsintegration.Command) ([]importer.Command, []importer.Failure) {
	commands := make([]importer.Command, 0, len(input))
	failures := make([]importer.Failure, 0)
	for _, command := range input {
		role, ok := streamElementsRole(command.AccessLevel)
		if !ok {
			failures = append(failures, importer.Failure{
				Name:   command.Name,
				Reason: importer.FailureUnsupportedRole,
			})
			continue
		}

		isReply, ok := streamElementsResponseType(command.Type)
		if !ok {
			failures = append(failures, importer.Failure{
				Name:   command.Name,
				Reason: importer.FailureUnsupportedResponse,
			})
			continue
		}

		aliases := command.Aliases
		if aliases == nil {
			aliases = []string{}
		}

		response := command.Response
		if command.Type == "action" && !strings.HasPrefix(response, "/me ") {
			response = "/me " + response
		}

		commands = append(commands, importer.Command{
			Name:        command.Name,
			Response:    response,
			Enabled:     command.Enabled && (command.EnabledOnline || command.EnabledOffline),
			Visible:     !command.Hidden,
			IsReply:     isReply,
			Aliases:     aliases,
			Cooldown:    command.Cooldown.Global,
			Role:        role,
			OnlineOnly:  command.EnabledOnline && !command.EnabledOffline,
			OfflineOnly: command.EnabledOffline && !command.EnabledOnline,
		})
	}

	return commands, failures
}

func NormalizeTimers(input []streamelementsintegration.Timer) ([]importer.Timer, []importer.Failure) {
	timers := make([]importer.Timer, 0, len(input))
	failures := make([]importer.Failure, 0)
	for _, timer := range input {
		if timer.Online.Enabled && timer.Offline.Enabled && timer.Online.Interval != timer.Offline.Interval {
			failures = append(failures, importer.Failure{
				Name:   timer.Name,
				Reason: importer.FailureIncompatibleInterval,
			})
			continue
		}

		interval := timer.Online.Interval
		if !timer.Online.Enabled {
			interval = timer.Offline.Interval
		}
		timers = append(timers, importer.Timer{
			Name:            timer.Name,
			Message:         timer.Message,
			Enabled:         timer.Enabled,
			OnlineEnabled:   timer.Online.Enabled,
			OfflineEnabled:  timer.Offline.Enabled,
			TimeInterval:    interval,
			MessageInterval: timer.ChatLines,
		})
	}

	return timers, failures
}

func streamElementsRole(level int) (importer.RoleRequirement, bool) {
	switch {
	case level == 100:
		return importer.RoleEveryone, true
	case level == 250:
		return importer.RoleSubscriber, true
	case level == 400:
		return importer.RoleVip, true
	case level == 500:
		return importer.RoleModerator, true
	case level >= 1000:
		return importer.RoleBroadcaster, true
	default:
		return "", false
	}
}

func streamElementsResponseType(responseType string) (bool, bool) {
	switch responseType {
	case "reply":
		return true, true
	case "say", "action":
		return false, true
	default:
		return false, false
	}
}
