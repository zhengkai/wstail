package main

import (
	"fmt"
	"pb"
	"strings"
	"unicode/utf8"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

func encodeReturn(u proto.Message) (b []byte) {

	msg, err := ptypes.MarshalAny(u)
	if err != nil {
		fmt.Println(err)
	}

	send := &pb.MsgReturn{
		Msg: []*any.Any{msg},
	}

	b, err = proto.Marshal(send)
	if err != nil {
		fmt.Println(err)
	}

	return
}

func utf8string(input string) string {
	return strings.Map(fixUtf8, input)
}

func fixUtf8(r rune) rune {
	if r == utf8.RuneError {
		return -1
	}
	return r
}
