package eventsub_bindings

import (
	"encoding/json"
)

type Subscription struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Version   string      `json:"version"`
	Status    string      `json:"status"`
	Cost      int         `json:"cost"`
	Condition interface{} `json:"condition"`
	CreatedAt string      `json:"created_at"`
}

func (s *Subscription) ConditionChannelBan() (*ConditionChannelBan, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelBan
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelSubscribe() (*ConditionChannelSubscribe, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelSubscribe
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelSubscriptionEnd() (*ConditionChannelSubscriptionEnd, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelSubscriptionEnd
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelSubscriptionGift() (*ConditionChannelSubscriptionGift, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelSubscriptionGift
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelSubscriptionMessage() (*ConditionChannelSubscriptionMessage, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelSubscriptionMessage
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelCheer() (*ConditionChannelCheer, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelCheer
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelUpdate() (*ConditionChannelUpdate, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelUpdate
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelFollow() (*ConditionChannelFollow, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelFollow
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelUnban() (*ConditionChannelUnban, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelUnban
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelRaid() (*ConditionChannelRaid, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelRaid
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelModeratorAdd() (*ConditionChannelModeratorAdd, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelModeratorAdd
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelModeratorRemove() (*ConditionChannelModeratorRemove, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelModeratorRemove
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelPointsRewardAdd() (*ConditionChannelPointsRewardAdd, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelPointsRewardAdd
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelPointsRewardUpdate() (*ConditionChannelPointsRewardUpdate, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelPointsRewardUpdate
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelPointsRewardRemove() (*ConditionChannelPointsRewardRemove, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelPointsRewardRemove
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelPointsRewardRedemptionAdd() (*ConditionChannelPointsRewardRedemptionAdd, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelPointsRewardRedemptionAdd
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelPointsRewardRedemptionUpdate() (*ConditionChannelPointsRewardRedemptionUpdate, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelPointsRewardRedemptionUpdate
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelPollBegin() (*ConditionChannelPollBegin, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelPollBegin
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelPollProgress() (*ConditionChannelPollProgress, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelPollProgress
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelPollEnd() (*ConditionChannelPollEnd, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelPollEnd
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelPredictionBegin() (*ConditionChannelPredictionBegin, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelPredictionBegin
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelPredictionLock() (*ConditionChannelPredictionLock, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelPredictionLock
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionChannelPredictionEnd() (*ConditionChannelPredictionEnd, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionChannelPredictionEnd
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionDropEntitlementGrant() (*ConditionDropEntitlementGrant, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionDropEntitlementGrant
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionExtensionBitsTransactionCreate() (*ConditionExtensionBitsTransactionCreate, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionExtensionBitsTransactionCreate
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionGoals() (*ConditionGoals, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionGoals
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionHypeTrainBegin() (*ConditionHypeTrainBegin, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionHypeTrainBegin
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionHypeTrainProgress() (*ConditionHypeTrainProgress, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionHypeTrainProgress
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionHypeTrainEnd() (*ConditionHypeTrainEnd, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionHypeTrainEnd
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionStreamOnline() (*ConditionStreamOnline, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionStreamOnline
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionStreamOffline() (*ConditionStreamOffline, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionStreamOffline
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionUserAuthorizationGrant() (*ConditionUserAuthorizationGrant, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionUserAuthorizationGrant
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionUserAuthorizationRevoke() (*ConditionUserAuthorizationRevoke, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionUserAuthorizationRevoke
		return &condition, json.Unmarshal(data, &condition)
	}
}

func (s *Subscription) ConditionUserUpdate() (*ConditionUserUpdate, error) {
	if data, err := json.Marshal(s.Condition); err != nil {
		return nil, err
	} else {
		var condition ConditionUserUpdate
		return &condition, json.Unmarshal(data, &condition)
	}
}
