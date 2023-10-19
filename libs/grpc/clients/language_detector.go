package clients

import (
	"fmt"
	"log"

	"github.com/satont/twir/libs/grpc/constants"
	language_detector "github.com/satont/twir/libs/grpc/generated/language-detector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewLanguageDetector(env string) language_detector.LanguageDetectorClient {
	serverAddress := createClientAddr(
		env,
		"language-detector",
		constants.LANGUAGE_DETECTOR_SERVER_PORT,
	)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(
			fmt.Sprintf(
				`{"loadBalancingConfig": [{"%s":{}}]}`,
				roundrobin.Name,
			),
		),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return language_detector.NewLanguageDetectorClient(conn)
}
