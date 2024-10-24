package auditlog

import (
	"slices"
)

// filter filters audit log and returns whether it should be sent to subscribers.
func filter(auditLog AuditLog, dashboardIDs ...string) bool {
	if len(dashboardIDs) != 0 && auditLog.ChannelID.IsZero() {
		return false
	}

	if !auditLog.ChannelID.IsZero() {
		if !slices.Contains(
			dashboardIDs,
			auditLog.ChannelID.String,
		) {
			return false
		}
	}

	return true
}
