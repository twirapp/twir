package eventsub_framework

import (
	"encoding/json"
	esb "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/mozillazg/go-httpheader"
	"io"
	"net/http"
)

const (
	webhookCallbackVerification = "webhook_callback_verification"
	notificationMessageType     = "notification"
)

type SubHandler struct {
	doSignatureVerification bool
	signatureSecret         []byte

	// Challenge handler function.
	// Returns whether the subscription should be accepted.
	VerifyChallenge func(h *esb.ResponseHeaders, chal *esb.SubscriptionChallenge) bool

	// IDTracker used to deduplicate notifications
	IDTracker               IDTracker
	OnDuplicateNotification func(h *esb.ResponseHeaders)

	HandleChannelUpdate func(h *esb.ResponseHeaders, event *esb.EventChannelUpdate)
	HandleChannelFollow func(h *esb.ResponseHeaders, event *esb.EventChannelFollow)
	HandleUserUpdate    func(h *esb.ResponseHeaders, event *esb.EventUserUpdate)

	HandleChannelSubscribe           func(h *esb.ResponseHeaders, event *esb.EventChannelSubscribe)
	HandleChannelSubscriptionEnd     func(h *esb.ResponseHeaders, event *esb.EventChannelSubscriptionEnd)
	HandleChannelSubscriptionGift    func(h *esb.ResponseHeaders, event *esb.EventChannelSubscriptionGift)
	HandleChannelSubscriptionMessage func(h *esb.ResponseHeaders, event *esb.EventChannelSubscriptionMessage)
	HandleChannelCheer               func(h *esb.ResponseHeaders, event *esb.EventChannelCheer)
	HandleChannelRaid                func(h *esb.ResponseHeaders, event *esb.EventChannelRaid)

	HandleChannelBan             func(h *esb.ResponseHeaders, event *esb.EventChannelBan)
	HandleChannelUnban           func(h *esb.ResponseHeaders, event *esb.EventChannelUnban)
	HandleChannelModeratorAdd    func(h *esb.ResponseHeaders, event *esb.EventChannelModeratorAdd)
	HandleChannelModeratorRemove func(h *esb.ResponseHeaders, event *esb.EventChannelModeratorRemove)

	HandleChannelPointsRewardAdd              func(h *esb.ResponseHeaders, event *esb.EventChannelPointsRewardAdd)
	HandleChannelPointsRewardUpdate           func(h *esb.ResponseHeaders, event *esb.EventChannelPointsRewardUpdate)
	HandleChannelPointsRewardRemove           func(h *esb.ResponseHeaders, event *esb.EventChannelPointsRewardRemove)
	HandleChannelPointsRewardRedemptionAdd    func(h *esb.ResponseHeaders, event *esb.EventChannelPointsRewardRedemptionAdd)
	HandleChannelPointsRewardRedemptionUpdate func(h *esb.ResponseHeaders, event *esb.EventChannelPointsRewardRedemptionUpdate)

	HandleChannelPollBegin    func(h *esb.ResponseHeaders, event *esb.EventChannelPollBegin)
	HandleChannelPollProgress func(h *esb.ResponseHeaders, event *esb.EventChannelPollProgress)
	HandleChannelPollEnd      func(h *esb.ResponseHeaders, event *esb.EventChannelPollEnd)

	HandleChannelPredictionBegin    func(h *esb.ResponseHeaders, event *esb.EventChannelPredictionBegin)
	HandleChannelPredictionProgress func(h *esb.ResponseHeaders, event *esb.EventChannelPredictionProgress)
	HandleChannelPredictionLock     func(h *esb.ResponseHeaders, event *esb.EventChannelPredictionLock)
	HandleChannelPredictionEnd      func(h *esb.ResponseHeaders, event *esb.EventChannelPredictionEnd)

	HandleDropEntitlementGrant           func(h *esb.ResponseHeaders, event *esb.EventDropEntitlementGrant)
	HandleExtensionBitsTransactionCreate func(h *esb.ResponseHeaders, event *esb.EventBitsTransactionCreate)

	HandleGoalBegin    func(h *esb.ResponseHeaders, event *esb.EventGoals)
	HandleGoalProgress func(h *esb.ResponseHeaders, event *esb.EventGoals)
	HandleGoalEnd      func(h *esb.ResponseHeaders, event *esb.EventGoals)

	HandleHypeTrainBegin    func(h *esb.ResponseHeaders, event *esb.EventHypeTrainBegin)
	HandleHypeTrainProgress func(h *esb.ResponseHeaders, event *esb.EventHypeTrainProgress)
	HandleHypeTrainEnd      func(h *esb.ResponseHeaders, event *esb.EventHypeTrainEnd)

	HandleStreamOnline  func(h *esb.ResponseHeaders, event *esb.EventStreamOnline)
	HandleStreamOffline func(h *esb.ResponseHeaders, event *esb.EventStreamOffline)

	HandleUserAuthorizationGrant  func(h *esb.ResponseHeaders, event *esb.EventUserAuthorizationGrant)
	HandleUserAuthorizationRevoke func(h *esb.ResponseHeaders, event *esb.EventUserAuthorizationRevoke)
}

func NewSubHandler(doSignatureVerification bool, secret []byte) *SubHandler {
	if doSignatureVerification && secret == nil {
		panic("secret must be set if signature verification is enabled")
	}

	return &SubHandler{
		doSignatureVerification: doSignatureVerification,
		signatureSecret:         secret,
	}
}

func (s *SubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		defer r.Body.Close()
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if s.doSignatureVerification {
			valid, err := VerifyRequestSignature(r, bodyBytes, s.signatureSecret)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if !valid {
				http.Error(w, "Invalid request signature", http.StatusForbidden)
				return
			}
		}

		var h esb.ResponseHeaders
		if err := httpheader.Decode(r.Header, &h); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if s.IDTracker != nil {
			duplicate, err := s.IDTracker.AddAndCheckIfDuplicate(r.Context(), h.MessageID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if duplicate {
				if s.OnDuplicateNotification != nil {
					go s.OnDuplicateNotification(&h)
				}
				w.WriteHeader(http.StatusNoContent) // ignore and return 2XX code
				return
			}
		}

		if h.MessageType == webhookCallbackVerification {
			var data esb.SubscriptionChallenge
			if err := json.Unmarshal(bodyBytes, &data); err != nil {
				http.Error(w, "Invalid JSON body", http.StatusBadRequest)
				return
			}
			if s.VerifyChallenge == nil || s.VerifyChallenge(&h, &data) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(data.Challenge))
			} else {
				http.Error(w, "Invalid subscription", http.StatusBadRequest)
			}
			return
		} else if h.MessageType == notificationMessageType {
			var notification esb.EventNotification
			if err := json.Unmarshal(bodyBytes, &notification); err != nil {
				http.Error(w, "Invalid JSON body", http.StatusBadRequest)
				return
			}
			event := notification.Event

			switch h.SubscriptionType {
			case "channel.update":
				var data esb.EventChannelUpdate
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelUpdate != nil {
					go s.HandleChannelUpdate(&h, &data)
				}
			case "channel.follow":
				var data esb.EventChannelFollow
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelFollow != nil {
					go s.HandleChannelFollow(&h, &data)
				}
			case "channel.subscribe":
				var data esb.EventChannelSubscribe
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelSubscribe != nil {
					go s.HandleChannelSubscribe(&h, &data)
				}
			case "channel.subscription.end":
				var data esb.EventChannelSubscriptionEnd
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelSubscriptionEnd != nil {
					go s.HandleChannelSubscriptionEnd(&h, &data)
				}
			case "channel.subscription.gift":
				var data esb.EventChannelSubscriptionGift
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelSubscriptionGift != nil {
					go s.HandleChannelSubscriptionGift(&h, &data)
				}
			case "channel.subscription.message":
				var data esb.EventChannelSubscriptionMessage
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelSubscriptionMessage != nil {
					go s.HandleChannelSubscriptionMessage(&h, &data)
				}
			case "channel.cheer":
				var data esb.EventChannelCheer
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelCheer != nil {
					go s.HandleChannelCheer(&h, &data)
				}
			case "channel.raid":
				var data esb.EventChannelRaid
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelRaid != nil {
					go s.HandleChannelRaid(&h, &data)
				}
			case "channel.ban":
				var data esb.EventChannelBan
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelBan != nil {
					go s.HandleChannelBan(&h, &data)
				}
			case "channel.unban":
				var data esb.EventChannelUnban
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelUnban != nil {
					go s.HandleChannelUnban(&h, &data)
				}
			case "channel.moderator.add":
				var data esb.EventChannelModeratorAdd
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelModeratorAdd != nil {
					go s.HandleChannelModeratorAdd(&h, &data)
				}
			case "channel.moderator.remove":
				var data esb.EventChannelModeratorRemove
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelModeratorRemove != nil {
					go s.HandleChannelModeratorRemove(&h, &data)
				}
			case "channel.channel_points_custom_reward.add":
				var data esb.EventChannelPointsRewardAdd
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelPointsRewardAdd != nil {
					go s.HandleChannelPointsRewardAdd(&h, &data)
				}
			case "channel.channel_points_custom_reward.update":
				var data esb.EventChannelPointsRewardUpdate
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelPointsRewardUpdate != nil {
					go s.HandleChannelPointsRewardUpdate(&h, &data)
				}
			case "channel.channel_points_custom_reward.remove":
				var data esb.EventChannelPointsRewardRemove
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelPointsRewardRemove != nil {
					go s.HandleChannelPointsRewardRemove(&h, &data)
				}
			case "channel.channel_points_custom_reward_redemption.add":
				var data esb.EventChannelPointsRewardRedemptionAdd
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelPointsRewardRedemptionAdd != nil {
					go s.HandleChannelPointsRewardRedemptionAdd(&h, &data)
				}
			case "channel.channel_points_custom_reward_redemption.update":
				var data esb.EventChannelPointsRewardRedemptionUpdate
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelPointsRewardRedemptionUpdate != nil {
					go s.HandleChannelPointsRewardRedemptionUpdate(&h, &data)
				}
			case "channel.poll.begin":
				var data esb.EventChannelPollBegin
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelPollBegin != nil {
					go s.HandleChannelPollBegin(&h, &data)
				}
			case "channel.poll.progress":
				var data esb.EventChannelPollProgress
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelPollProgress != nil {
					go s.HandleChannelPollProgress(&h, &data)
				}
			case "channel.poll.end":
				var data esb.EventChannelPollEnd
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelPollEnd != nil {
					go s.HandleChannelPollEnd(&h, &data)
				}
			case "channel.prediction.begin":
				var data esb.EventChannelPredictionBegin
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelPredictionBegin != nil {
					go s.HandleChannelPredictionBegin(&h, &data)
				}
			case "channel.prediction.progress":
				var data esb.EventChannelPredictionProgress
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelPredictionProgress != nil {
					go s.HandleChannelPredictionProgress(&h, &data)
				}
			case "channel.prediction.lock":
				var data esb.EventChannelPredictionLock
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelPredictionLock != nil {
					go s.HandleChannelPredictionLock(&h, &data)
				}
			case "channel.prediction.end":
				var data esb.EventChannelPredictionEnd
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleChannelPredictionEnd != nil {
					go s.HandleChannelPredictionEnd(&h, &data)
				}
			case "drop.entitlement.grant":
				var data esb.EventDropEntitlementGrant
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleDropEntitlementGrant != nil {
					go s.HandleDropEntitlementGrant(&h, &data)
				}
			case "extension.bits_transaction.create":
				var data esb.EventBitsTransactionCreate
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleExtensionBitsTransactionCreate != nil {
					go s.HandleExtensionBitsTransactionCreate(&h, &data)
				}
			case "channel.goal.begin":
				var data esb.EventGoals
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleGoalBegin != nil {
					go s.HandleGoalBegin(&h, &data)
				}
			case "channel.goal.progress":
				var data esb.EventGoals
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleGoalProgress != nil {
					go s.HandleGoalProgress(&h, &data)
				}
			case "channel.goal.end":
				var data esb.EventGoals
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleGoalEnd != nil {
					go s.HandleGoalEnd(&h, &data)
				}
			case "channel.hype_train.begin":
				var data esb.EventHypeTrainBegin
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleHypeTrainBegin != nil {
					go s.HandleHypeTrainBegin(&h, &data)
				}
			case "channel.hype_train.progress":
				var data esb.EventHypeTrainProgress
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleHypeTrainProgress != nil {
					go s.HandleHypeTrainProgress(&h, &data)
				}
			case "channel.hype_train.end":
				var data esb.EventHypeTrainEnd
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleHypeTrainEnd != nil {
					go s.HandleHypeTrainEnd(&h, &data)
				}
			case "stream.online":
				var data esb.EventStreamOnline
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleStreamOnline != nil {
					go s.HandleStreamOnline(&h, &data)
				}
			case "stream.offline":
				var data esb.EventStreamOffline
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleStreamOffline != nil {
					go s.HandleStreamOffline(&h, &data)
				}
			case "user.authorization.grant":
				var data esb.EventUserAuthorizationGrant
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleUserAuthorizationGrant != nil {
					go s.HandleUserAuthorizationGrant(&h, &data)
				}
			case "user.authorization.revoke":
				var data esb.EventUserAuthorizationRevoke
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleUserAuthorizationRevoke != nil {
					go s.HandleUserAuthorizationRevoke(&h, &data)
				}
			case "user.update":
				var data esb.EventUserUpdate
				if err := json.Unmarshal(event, &data); err != nil {
					http.Error(w, "Invalid JSON body", http.StatusBadRequest)
					return
				}
				if s.HandleUserUpdate != nil {
					go s.HandleUserUpdate(&h, &data)
				}
			default:
				http.Error(w, "Unknown notification type", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			http.Error(w, "Unknown message type", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
