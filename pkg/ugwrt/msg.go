package ugwrt

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type ByteMessage struct {
	buf []byte
}

func ResolveMessage(path string) interface{} {
	return nil
}

func (rt *RtInstance) ResolveProtos(folder string) ([]*desc.FileDescriptor, error) {
	parser := protoparse.Parser{}
	protoFiles := make([]string, 0)
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".proto" {
			return nil
		}
		protoFiles = append(protoFiles, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	descs, err := parser.ParseFiles(protoFiles...)
	if err != nil {
		return nil, err
	}
	for _, val := range descs {
		pkgName := val.GetPackage()
		messages := val.GetMessageTypes()
		rt.Log(AppLog, log.InfoLevel, "ResolveProto: Package=[%s]", pkgName)
		for _, enum := range val.GetEnumTypes() {
			rt.Log(AppLog, log.InfoLevel, "ResolveProto: Enum=[%s]", enum.GetName())
		}
		for _, msg := range messages {
			rt.Log(AppLog, log.InfoLevel, "ResolveProto: MsgName=[%s], Options=[%s]",
				msg.GetName())
			for _, field := range msg.GetFields() {
				rt.Log(AppLog, log.InfoLevel, "ResolveProto: Field=[%s], Options=[%s]", field, field.GetFieldOptions())
			}
		}
	}
	// TODO
	return nil, nil
}
