package serializer

import (
	"fmt"
	"math/rand"
	"os"
	"pcbook/pb"
	"pcbook/sample"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func setup() {
	fmt.Printf("\033[1;33m%s\033[0m", "Starting serializer tests...\n")
}

func teardown() {
	fmt.Printf("\033[1;33m%s\033[0m", "Finished serializer tests.\n")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int()
	}
}

func TestWriteProtoToBinary(t *testing.T) {
	t.Parallel()

	dir := "../tmp/test1.bin"

	laptop := sample.NewLaptop()
	err := WriteProtoBufToBinaryFile(laptop, dir)
	require.NoError(t, err)
}

func TestReadProtoFromBinary(t *testing.T) {
	t.Parallel()

	dir := "../tmp/test2.bin"

	laptop1 := sample.NewLaptop()
	err := WriteProtoBufToBinaryFile(laptop1, dir)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err = ReadBinaryToProtoFile(dir, laptop2)
	require.NoError(t, err)

	require.True(t, proto.Equal(laptop1, laptop2))
}

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop1 := sample.NewLaptop()

	err := WriteProtoBufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)

	err = WriteProtobufToJSONFile(laptop1, jsonFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err = ReadBinaryToProtoFile(binaryFile, laptop2)
	require.NoError(t, err)

	require.True(t, proto.Equal(laptop1, laptop2))
}
