package ugwrt

import "google.golang.org/protobuf/reflect/protoreflect"

type ByteMessage struct {
	buf []byte
}

func ResolveMessage(path string) interface{} {
	_ = protoreflect.Message(nil)
	// TODO
	return nil
}
