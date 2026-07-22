package timers

import (
	"context"
	"testing"
)

func TestBuildWatchedChannelIDsQueryUsesTwitchBinding(t *testing.T) {
	t.Parallel()

	db := newDryRunPostgresDB(t)
	statement := buildWatchedChannelIDsQuery(db, context.Background()).Scan(&[]struct {
		ID string `gorm:"column:id"`
	}{}).Statement
	sql := statement.SQL.String()

	assertQueryContains(t, sql,
		"SELECT DISTINCT cp.channel_id AS id",
		"FROM channels_streams cs",
		"JOIN channel_platforms cp",
		"cp.platform = 'twitch'",
		`cp.platform_channel_id = cs."userId"`,
		"WHERE cs.platform = 'twitch'",
	)
	assertQueryExcludes(t, sql, "c.twitch_user_id", "JOIN users u", "u.platform_id")
	assertSingleBindingSource(t, sql)
}
