package model

type VoteState struct {
	TargetUserID   string `redis:"target_user_id"`
	TargetUserName string `redis:"target_user_name"`
	TargetIsMod    bool   `redis:"target_is_mod"`
	TotalVotes     int    `redis:"total_votes"`
	PositiveVotes  int    `redis:"positive_votes"`
	NegativeVotes  int    `redis:"negative_votes"`
}

var Nil = VoteState{}
