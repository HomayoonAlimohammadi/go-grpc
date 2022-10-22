package serializer

import (
	"fmt"
	"io/ioutil"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func WriteProtoBufToBinaryFile(message proto.Message, filename string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("can not marshal the proto message: %w", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("can not write data to binary file: %w", err)
	}

	return nil
}

func ReadBinaryToProtoFile(filename string, message proto.Message) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("can not read binary file from path: %w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("can not unmarshal data to proto message: %w", err)
	}

	return nil
}

func ProtoBufToJSON(message proto.Message) (string, error) {
	bytesMarshal, err := protojson.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("can not initialize proto marshaler: %w", err)
	}

	return string(bytesMarshal), nil
}

func WriteProtobufToJSONFile(message proto.Message, filename string) error {
	data, err := ProtoBufToJSON(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to JSON: %w", err)
	}

	err = ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("cannot write JSON data to file: %w", err)
	}

	return nil
}
