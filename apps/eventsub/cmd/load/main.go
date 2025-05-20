package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func main() {
	count := 2_000_000

	messages := make([]eventsub_bindings.EventChannelPointsRewardRedemptionAdd, 0, count)
	for i := 0; i < count; i++ {
		messages = append(
			messages, eventsub_bindings.EventChannelPointsRewardRedemptionAdd{
				ID:                   strconv.Itoa(i),
				BroadcasterUserID:    "870284304",
				BroadcasterUserLogin: "twirdev",
				BroadcasterUserName:  "TwirDev",
				UserID:               "1174832614",
				UserLogin:            "postgresiq",
				UserName:             "postgresiq",
				UserInput:            "",
				Status:               "unfulfilled",
				Reward: eventsub_bindings.Reward{
					ID:     strconv.Itoa(i),
					Title:  "test",
					Cost:   100,
					Prompt: "",
				},
				RedeemedAt: "",
			},
		)
	}

	convertedMessages := make([][]byte, 0, len(messages))
	for _, m := range messages {
		jsonData, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		convertedMessages = append(convertedMessages, jsonData)
	}

	ticker := time.NewTicker(100 * time.Millisecond)

	i := 0

	for {
		select {
		case <-ticker.C:
			var wg sync.WaitGroup

			for _, msg := range convertedMessages[i : i+100] {
				wg.Add(1)
				go func() {
					defer wg.Done()

					body := bytes.NewBuffer(msg)

					req, _ := http.NewRequest("POST", "https://eventsub.twir.app", body)

					_, err := http.DefaultClient.Do(req)
					if err != nil {
						panic(err)
					}
				}()
			}

			wg.Wait()
			log.Println("Sent 100 messages")
			i += 100
		}
	}
}
