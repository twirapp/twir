package predictions

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/stretchr/testify/require"
)

func newReplica(f *fixture, client predictionClient) *Predictions {
	return newPredictions(
		f.settings,
		f.channels,
		&fakeClientFactory{client: client},
		f.store,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		nil,
		nil,
	)
}

func newReplicaClient() *fakePredictionClient {
	return &fakePredictionClient{endResponse: &helix.PredictionsResponse{}}
}

func createPendingPrediction(t *testing.T, f *fixture, matchID int64) pendingPredictionIntent {
	t.Helper()
	message := startMessage(f, matchID)
	f.store.commitErr = errors.New("redis temporarily unavailable")
	_, err := f.predictions.handleMatchStarted(context.Background(), message)
	require.ErrorIs(t, err, f.store.commitErr)
	require.Len(t, f.client.CreateCalls(), 1)
	f.store.commitErr = nil

	intent, err := f.store.GetPending(context.Background(), predictionKey(f.channelID, matchID))
	require.NoError(t, err)
	return intent
}

func activePredictionForIntent(intent pendingPredictionIntent, id string) helix.Prediction {
	return helix.Prediction{
		ID:        id,
		Status:    "ACTIVE",
		Title:     intent.Title,
		CreatedAt: helix.Time{Time: intent.ReservedAt.Add(time.Second)},
		Outcomes: []helix.Outcomes{
			{ID: "yes-outcome", Title: intent.YesOutcomeTitle},
			{ID: "no-outcome", Title: intent.NoOutcomeTitle},
		},
	}
}

func predictionsResponse(predictions ...helix.Prediction) *helix.PredictionsResponse {
	return &helix.PredictionsResponse{
		Data: helix.ManyPredictions{Predictions: predictions},
	}
}

func TestAnotherReplicaRecoversPendingPredictionOnStartReplay(t *testing.T) {
	f := newFixture(t)
	message := startMessage(f, 2_001)
	intent := createPendingPrediction(t, f, message.MatchID)
	replicaClient := newReplicaClient()
	replica := newReplica(f, replicaClient)
	replicaClient.getResponse = predictionsResponse(activePredictionForIntent(intent, "prediction-1"))

	_, err := replica.handleMatchStarted(context.Background(), message)

	require.NoError(t, err)
	require.Len(t, f.client.CreateCalls(), 1)
	require.Empty(t, replicaClient.CreateCalls())
	record, exists := f.store.Record(predictionKey(f.channelID, message.MatchID))
	require.True(t, exists)
	require.Equal(t, storedPrediction{
		PredictionID: "prediction-1",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	}, record)
	getCalls := replicaClient.GetCalls()
	require.Len(t, getCalls, 1)
	require.Empty(t, getCalls[0].ID)
}

func TestAnotherReplicaRecoversPendingPredictionForTerminalEvents(t *testing.T) {
	for _, tt := range []struct {
		name   string
		handle func(*Predictions, *fixture, int64) error
		status string
	}{
		{
			name: "resolve",
			handle: func(replica *Predictions, f *fixture, matchID int64) error {
				_, err := replica.handleMatchEnded(context.Background(), endMessage(f, matchID, true))
				return err
			},
			status: "RESOLVED",
		},
		{
			name: "cancel",
			handle: func(replica *Predictions, f *fixture, matchID int64) error {
				_, err := replica.handleMatchAbandoned(context.Background(), abandonedMessage(f, matchID))
				return err
			},
			status: "CANCELED",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			matchID := int64(2_002)
			intent := createPendingPrediction(t, f, matchID)
			replicaClient := newReplicaClient()
			replica := newReplica(f, replicaClient)
			prediction := activePredictionForIntent(intent, "prediction-1")
			if tt.status == "CANCELED" {
				prediction.Status = "LOCKED"
			}
			replicaClient.getResponse = predictionsResponse(prediction)

			require.NoError(t, tt.handle(replica, f, matchID))
			require.Len(t, f.client.CreateCalls(), 1)
			require.Empty(t, replicaClient.CreateCalls())
			endCalls := replicaClient.EndCalls()
			require.Len(t, endCalls, 1)
			require.Equal(t, "prediction-1", endCalls[0].ID)
			require.Equal(t, tt.status, endCalls[0].Status)
			_, exists := f.store.Record(predictionKey(f.channelID, matchID))
			require.False(t, exists)
			_, err := f.store.GetPending(context.Background(), predictionKey(f.channelID, matchID))
			require.ErrorIs(t, err, errPredictionIntentNotFound)
		})
	}
}

func TestAnotherReplicaRetainsPendingIntentWhenCandidatesAreUnsafe(t *testing.T) {
	for _, tt := range []struct {
		name       string
		candidates func(pendingPredictionIntent) []helix.Prediction
	}{
		{
			name: "mismatched title",
			candidates: func(intent pendingPredictionIntent) []helix.Prediction {
				prediction := activePredictionForIntent(intent, "manual-prediction")
				prediction.Title = "A manual prediction"
				return []helix.Prediction{prediction}
			},
		},
		{
			name: "mismatched outcomes",
			candidates: func(intent pendingPredictionIntent) []helix.Prediction {
				prediction := activePredictionForIntent(intent, "manual-prediction")
				prediction.Outcomes[1].Title = "Maybe"
				return []helix.Prediction{prediction}
			},
		},
		{
			name: "mismatched creation timestamp",
			candidates: func(intent pendingPredictionIntent) []helix.Prediction {
				prediction := activePredictionForIntent(intent, "manual-prediction")
				prediction.CreatedAt = helix.Time{Time: intent.ReservedAt.Add(pendingPredictionCorrelationWindow + time.Second)}
				return []helix.Prediction{prediction}
			},
		},
		{
			name: "inactive candidate",
			candidates: func(intent pendingPredictionIntent) []helix.Prediction {
				prediction := activePredictionForIntent(intent, "manual-prediction")
				prediction.Status = "RESOLVED"
				return []helix.Prediction{prediction}
			},
		},
		{
			name: "ambiguous candidates",
			candidates: func(intent pendingPredictionIntent) []helix.Prediction {
				return []helix.Prediction{
					activePredictionForIntent(intent, "prediction-1"),
					activePredictionForIntent(intent, "prediction-2"),
				}
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			matchID := int64(2_003)
			intent := createPendingPrediction(t, f, matchID)
			replicaClient := newReplicaClient()
			replica := newReplica(f, replicaClient)
			replicaClient.getResponse = predictionsResponse(tt.candidates(intent)...)

			_, err := replica.handleMatchEnded(context.Background(), endMessage(f, matchID, true))

			require.Error(t, err)
			require.Empty(t, replicaClient.EndCalls())
			require.Empty(t, f.store.deleteCalls)
			require.True(t, f.store.HasReservation(predictionKey(f.channelID, matchID)))
			pending, pendingErr := f.store.GetPending(context.Background(), predictionKey(f.channelID, matchID))
			require.NoError(t, pendingErr)
			require.Equal(t, intent, pending)
		})
	}
}
