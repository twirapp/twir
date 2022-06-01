package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	rdb "github.com/go-redis/redis/v8"

	"mygrpc/streamstatus"

	twitch "github.com/nicklaw5/helix"
	"google.golang.org/grpc"
)

var (
	ctx   = context.Background()
	port  = flag.Int("port", 50052, "The server port")
	redis = rdb.NewClient(&rdb.Options{
		Addr:     "localhost:6379",
		Password: "576294Aa",
		DB:       0,
	})
)

type server struct {
	streamstatus.UnimplementedMainServer
}

func (s *server) CacheStreams(stream streamstatus.Main_CacheStreamsServer) error {
	for {
		data, readErr := stream.Recv()
		if data != nil {
			log.Print(data.GetChannelIds())
		}

		helix, err := twitch.NewClient(&twitch.Options{
			ClientID:     data.ClientId,
			ClientSecret: data.ClientSecret,
		})
		if err != nil {
			return fmt.Errorf("Cannot create helix client")
		}

		streams, err := helix.GetStreams(&twitch.StreamsParams{UserIDs: data.ChannelIds})

		if err != nil {
			return fmt.Errorf("Cannot get streams from twitch")
		}

		for _, UserID := range data.ChannelIds {
			for _, stream := range streams.Data.Streams {
				streamKey := "streams:" + UserID
				cachedStream, err := redis.Get(ctx, streamKey).Result()

				if err == rdb.Nil {
					// stream not exists
				}

				if stream.UserID == UserID {
					//stream user === data user
				}
			}
		}

		if readErr == io.EOF {
			return stream.SendAndClose(&streamstatus.CachedStreamResult{Success: true})
		}

		if readErr != nil {
			return err
		}
	}
}

func main() {
	flag.Parse()
	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	streamstatus.RegisterMainServer(s, &server{})
	log.Printf("server listening at %v", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
