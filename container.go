package boltdb

import (
	"bytes"
	"fmt"
	"reflect"

	"encoding/gob"

	"github.com/golang/protobuf/proto"
)

type container struct {
	T string
	V []byte
}

func cMarshal(m proto.Message) ([]byte, error) {
	typeName := proto.MessageName(m)
	bs, err := proto.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal proto message: %v", err)
	}

	c := container{T: typeName, V: bs}

	b := bytes.Buffer{}
	g := gob.NewEncoder(&b)
	err = g.Encode(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal container: %v", err)
	}

	return b.Bytes(), nil
}

func cUnmarshal(b []byte) (proto.Message, error) {
	r := bytes.NewReader(b)
	g := gob.NewDecoder(r)

	var c container
	err := g.Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal container: %v", err)
	}

	protoType := proto.MessageType(c.T)
	if protoType == nil {
		return nil, fmt.Errorf("unknown proto message type %s", c.T)
	}
	t := protoType.Elem()
	intPtr := reflect.New(t)
	instance := intPtr.Interface().(proto.Message)
	err = proto.Unmarshal(c.V, instance)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal proto message: %v", err)
	}

	return instance, nil
}
