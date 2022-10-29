package server_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"pcstore/pb"
	"pcstore/sample"
	"pcstore/serializer"
	"pcstore/server"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	serverAddress   string
	laptopClient    pb.LaptopServiceClient
	laptopStore     server.LaptopStore
	clientGenerator *sync.Once
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
	clientGenerator = &sync.Once{}
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

	clientGenerator.Do(func() {
		laptopStore = server.NewInMemoryLaptopStore()
		serverAddress = startTestLaptopServer(t, laptopStore)
		laptopClient = newTestLaptopServiceClient(t, serverAddress)
	})

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

func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()

	filter := &pb.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.2,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}

	clientGenerator.Do(func() {
		laptopStore = server.NewInMemoryLaptopStore()
		serverAddress = startTestLaptopServer(t, laptopStore)
		laptopClient = newTestLaptopServiceClient(t, serverAddress)
	})

	expectedIDs := make(map[string]bool)

	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()

		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Ram = &pb.Memory{Value: 4096, Unit: pb.Memory_MEGABYTE}
		case 4:
			laptop.PriceUsd = 1999
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.5
			laptop.Cpu.MaxGhz = laptop.Cpu.MinGhz + 2.0
			laptop.Ram = &pb.Memory{Value: 16, Unit: pb.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumberCores = 6
			laptop.Cpu.MinGhz = 2.8
			laptop.Cpu.MaxGhz = laptop.Cpu.MinGhz + 2.0
			laptop.Ram = &pb.Memory{Value: 64, Unit: pb.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		}

		err := laptopStore.Save(laptop)
		require.NoError(t, err)
	}

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)
	found := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		require.Contains(t, expectedIDs, res.GetLaptop().GetId())

		found += 1
	}

	require.Equal(t, len(expectedIDs), found)
}
