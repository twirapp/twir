package model

type ChannelGamesVoteBanRedisStruct struct {
	TargetUserId   string `redis:"target_user_id"`
	TargetUserName string `redis:"target_user_name"`
	TargetIsMod    bool   `redis:"target_is_mod"`
	TotalVotes     int    `redis:"total_votes"`
	PositiveVotes  int    `redis:"positive_votes"`
	NegativeVotes  int    `redis:"negative_votes"`
}
