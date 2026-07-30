package stats

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/twirapp/twir/libs/uptime"
)

const (
	staleThreshold      = time.Minute
	restartThreshold    = 5 * time.Minute
	platformPingTimeout = 5 * time.Second
	maxChatMessageRunes = 450
)

type platformPing struct {
	name      string
	url       string
	available bool
	duration  time.Duration
}

var platformPings = []platformPing{
	{name: "twitch", url: "https://api.twitch.tv/helix/users"},
	{name: "kick", url: "https://kick.com/api/v2/channels/kick"},
	{name: "7tv", url: "https://7tv.io/v3/users/twitch/11148817"},
	{name: "discord", url: "https://discord.com/api/v10/gateway"},
}

func summarizeServices(
	statuses []uptime.InstanceStatus,
	now time.Time,
	unavailable string,
) ([]string, []string, []string) {
	byService := make(map[string][]uptime.InstanceStatus)
	for _, status := range statuses {
		byService[status.Service] = append(byService[status.Service], status)
	}

	serviceNames := make([]string, 0, len(byService))
	for service := range byService {
		serviceNames = append(serviceNames, service)
	}
	sort.Strings(serviceNames)

	services := make([]string, 0, len(serviceNames))
	down := []string{}
	restarted := []string{}
	for _, service := range serviceNames {
		alive := []time.Duration{}
		for _, status := range byService[service] {
			if status.IsStale(now, staleThreshold) {
				down = append(down, service+displayInstance(status.Instance))
				continue
			}

			uptimeDuration := status.Uptime(now)
			alive = append(alive, uptimeDuration)
			if uptimeDuration < restartThreshold {
				restarted = append(restarted, service+displayInstance(status.Instance)+" ("+formatUptime(uptimeDuration)+")")
			}
		}

		services = append(services, formatServiceUptime(service, alive, unavailable))
	}

	return services, down, restarted
}

func formatServiceUptime(service string, alive []time.Duration, unavailable string) string {
	if len(alive) == 0 {
		return service + "×0 " + unavailable
	}

	minimum := alive[0]
	maximum := alive[0]
	for _, duration := range alive[1:] {
		if duration < minimum {
			minimum = duration
		}
		if duration > maximum {
			maximum = duration
		}
	}

	uptimeRange := formatUptime(minimum)
	if len(alive) > 1 {
		uptimeRange += "–" + formatUptime(maximum)
	}

	return service + "×" + strconv.Itoa(len(alive)) + " " + uptimeRange
}

func formatUptime(duration time.Duration) string {
	switch {
	case duration < time.Minute:
		return "<1m"
	case duration < time.Hour:
		return strconv.FormatInt(int64(duration/time.Minute), 10) + "m"
	case duration < 24*time.Hour:
		return strconv.FormatInt(int64(duration/time.Hour), 10) + "h"
	default:
		return strconv.FormatInt(int64(duration/(24*time.Hour)), 10) + "d"
	}
}

func displayInstance(instance string) string {
	if strings.HasPrefix(instance, "slot-") {
		return "#" + strings.TrimPrefix(instance, "slot-")
	}
	if len(instance) == 12 {
		if _, err := hex.DecodeString(instance); err == nil {
			return "#" + instance[7:]
		}
	}

	return "#" + instance
}

func formatPings(ctx context.Context) string {
	results := make([]platformPing, len(platformPings))
	copy(results, platformPings)
	client := &http.Client{Timeout: platformPingTimeout}

	var group sync.WaitGroup
	for index := range results {
		group.Add(1)
		go func(result *platformPing) {
			defer group.Done()

			startedAt := time.Now()
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, result.url, nil)
			if err != nil {
				return
			}
			response, err := client.Do(request)
			if err != nil {
				return
			}
			defer response.Body.Close()

			result.available = true
			result.duration = time.Since(startedAt)
		}(&results[index])
	}
	group.Wait()

	formatted := make([]string, 0, len(results))
	for _, result := range results {
		if !result.available {
			formatted = append(formatted, result.name+" ❌")
			continue
		}
		formatted = append(formatted, fmt.Sprintf("%s %dms", result.name, result.duration.Milliseconds()))
	}

	return strings.Join(formatted, " · ")
}

func trimChatMessage(message string) string {
	if utf8.RuneCountInString(message) <= maxChatMessageRunes {
		return message
	}

	runes := []rune(message)
	return string(runes[:maxChatMessageRunes-1]) + "…"
}
