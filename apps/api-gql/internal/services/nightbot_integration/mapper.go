package nightbot_integration

import (
	"strconv"
	"strings"

	"github.com/twirapp/twir/apps/api-gql/internal/services/importer"
)

func NormalizeCommands(input nightbotCustomCommandsResponse) ([]importer.Command, []importer.Failure) {
	aliasesByCommand := make(map[string][]string)
	for _, command := range input.Commands {
		if command.Alias == nil || *command.Alias == "" {
			continue
		}

		aliasTarget := normalizeCommandName(*command.Alias)
		aliasesByCommand[aliasTarget] = append(aliasesByCommand[aliasTarget], normalizeCommandName(command.Name))
	}

	commands := make([]importer.Command, 0, len(input.Commands))
	failures := make([]importer.Failure, 0)
	for _, command := range input.Commands {
		if command.Alias != nil && *command.Alias != "" {
			continue
		}

		role, ok := nightbotRole(command.UserLevel)
		if !ok {
			failures = append(failures, importer.Failure{
				Name:   command.Name,
				Reason: importer.FailureUnsupportedRole,
			})
			continue
		}

		aliases := aliasesByCommand[normalizeCommandName(command.Name)]
		if aliases == nil {
			aliases = []string{}
		}

		commands = append(commands, importer.Command{
			Name:     normalizeCommandName(command.Name),
			Response: command.Message,
			Enabled:  true,
			Visible:  true,
			IsReply:  true,
			Aliases:  aliases,
			Cooldown: command.CoolDown,
			Role:     role,
		})
	}

	return commands, failures
}

func NormalizeTimers(input nightbotTimersResponse) ([]importer.Timer, []importer.Failure) {
	timers := make([]importer.Timer, 0, len(input.Timers))
	failures := make([]importer.Failure, 0)
	for _, timer := range input.Timers {
		interval, ok := cronIntervalMinutes(timer.Interval)
		if !ok {
			failures = append(failures, importer.Failure{
				Name:   timer.Name,
				Reason: importer.FailureIncompatibleInterval,
			})
			continue
		}

		timers = append(timers, importer.Timer{
			Name:            timer.Name,
			Message:         timer.Message,
			Enabled:         timer.Enabled,
			OnlineEnabled:   true,
			TimeInterval:    interval,
			MessageInterval: timer.Lines,
		})
	}

	return timers, failures
}

func normalizeCommandName(name string) string {
	return strings.ToLower(strings.TrimPrefix(name, "!"))
}

func nightbotRole(level string) (importer.RoleRequirement, bool) {
	switch level {
	case "owner":
		return importer.RoleBroadcaster, true
	case "moderator":
		return importer.RoleModerator, true
	case "subscriber":
		return importer.RoleSubscriber, true
	case "twitch_vip":
		return importer.RoleVip, true
	case "everyone":
		return importer.RoleEveryone, true
	default:
		return "", false
	}
}

func cronIntervalMinutes(value string) (int, bool) {
	fields := strings.Fields(value)
	if len(fields) != 5 || fields[1] != "*" || fields[2] != "*" || fields[3] != "*" || fields[4] != "*" {
		return 0, false
	}

	minute := fields[0]
	if strings.HasPrefix(minute, "*/") {
		interval, err := strconv.Atoi(strings.TrimPrefix(minute, "*/"))
		return interval, err == nil && interval > 0
	}

	if minuteOfHour, err := strconv.Atoi(minute); err == nil && minuteOfHour >= 0 && minuteOfHour < 60 {
		return 60, true
	}

	return 0, false
}
