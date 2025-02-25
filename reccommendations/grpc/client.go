package grpc

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Rich-T-kid/musicShare/reccommendations/grpc/protobuff"
)

func GetReccomendations(ctx context.Context, user_uuid string) ([]*pb.SongBody, error) {
	conn, err := grpc.NewClient(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	cc := pb.NewSongServiceClient(conn)
	response, err := cc.GetSong(ctx, &pb.SongRequest{UserId: user_uuid})
	if err != nil {
		log.Println("Error has occured trying to get users song of the day,", user_uuid)
		return nil, err
	}
	fmt.Println("Got the response of this back", response)
	return response.Songs, nil
}
