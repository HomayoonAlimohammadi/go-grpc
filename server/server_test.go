package server_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"pcbook/pb"
	"pcbook/sample"
	"pcbook/serializer"
	"pcbook/server"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type testCase struct {
	name   string
	laptop *pb.Laptop
	store  server.LaptopStore
	code   codes.Code
}

func setup() {
	fmt.Printf("\033[1;33m%s\033[0m", "Starting server tests...\n")
}

func teardown() {
	fmt.Printf("\033[1;33m%s\033[0m", "Finished server tests.\n")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestServerCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopNoID := sample.NewLaptop()
	laptopNoID.Id = ""

	laptopInvalidID := sample.NewLaptop()
	laptopInvalidID.Id = "invalid-uuid"

	laptopDuplicateID := sample.NewLaptop()
	storeDuplicateID := server.NewInMemoryLaptopStore()
	err := storeDuplicateID.Save(laptopDuplicateID)
	require.Nil(t, err)

	testCases := []testCase{
		{
			name:   "success_with_id",
			laptop: sample.NewLaptop(),
			store:  server.NewInMemoryLaptopStore(),
			code:   codes.OK,
		},
		{
			name:   "success_no_id",
			laptop: laptopNoID,
			store:  server.NewInMemoryLaptopStore(),
			code:   codes.OK,
		},
		{
			name:   "failure_invalid_id",
			laptop: laptopInvalidID,
			store:  server.NewInMemoryLaptopStore(),
			code:   codes.InvalidArgument,
		},
		{
			name:   "failure_duplicate_id",
			laptop: laptopDuplicateID,
			store:  storeDuplicateID,
			code:   codes.AlreadyExists,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := &pb.CreateLaptopRequest{
				Laptop: tc.laptop,
			}
			server := server.NewLaptopServer(tc.store)
			res, err := server.CreateLaptop(context.Background(), req)

			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Id)
				if len(tc.laptop.Id) > 0 {
					require.Equal(t, tc.laptop.Id, res.Id)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, st.Code())
			}

		})
	}
}

func startTestLaptopServer(t *testing.T, laptopStore server.LaptopStore) string {
	laptopServer := server.NewLaptopServer(laptopStore)

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":8000")
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return listener.Addr().String()
}

func newTestLaptopServiceClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := server.NewInMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopStore)
	laptopClient := newTestLaptopServiceClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedID := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	other, err := laptopStore.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)
	requireSameLaptop(t, laptop, other)
}

func requireSameLaptop(t *testing.T, first, second *pb.Laptop) {
	json1, err := serializer.ProtoBufToJSON(first)
	require.NoError(t, err)

	json2, err := serializer.ProtoBufToJSON(second)
	require.NoError(t, err)

	require.Equal(t, json1, json2)
}
