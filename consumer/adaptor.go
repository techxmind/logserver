package consumer

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/techxmind/rollingfile"
)

var (
	_marshalers  = make(map[string]func(string) (Marshaler, error))
	_sinkTargets = make(map[string]func(string) (io.Writer, error))
)

func init() {
	RegisterMarshaler("json", func(_ string) (Marshaler, error) {
		return MarshalerFunc(JSONMarshaler), nil
	})

	RegisterMarshaler("csv", func(args string) (Marshaler, error) {
		headers := make([]string, 0)
		for _, header := range strings.Split(args, ",") {
			header = strings.TrimSpace(header)
			if header != "" {
				headers = append(headers, header)
			}
		}
		return NewCSVMarshaler(headers)
	})

	RegisterSinkTarget("stdout", func(_ string) (io.Writer, error) {
		return os.Stdout, nil
	})

	RegisterSinkTarget("file", func(args string) (io.Writer, error) {
		if args == "" {
			return nil, fmt.Errorf("SinkTarget[file] args[filename] is missing")
		}
		// filename:rollingFileSize:rollingFileLifeTime
		argsArr := strings.Split(args, ":")
		if len(argsArr) == 1 {
			return os.OpenFile(args, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		} else {
			filename := argsArr[0]
			options := make([]rollingfile.Option, 0, 1)
			maxSize, err := strconv.Atoi(argsArr[1])
			if err != nil {
				return nil, fmt.Errorf("SinkTarget[rollingfile] args invalid")
			}
			options = append(options, rollingfile.MaxSize(maxSize))
			if len(argsArr) > 2 {
				maxAge, err := strconv.Atoi(argsArr[2])
				if err != nil {
					return nil, fmt.Errorf("SinkTarget[rollingfile] args invalid")
				}
				options = append(options, rollingfile.MaxAge(maxAge))
			}
			if len(argsArr) > 3 {
				options = append(options, rollingfile.Suffix(argsArr[3]))
			}
			return rollingfile.New(filename, options...)
		}
	})
}

func RegisterMarshaler(marshaler string, factory func(string) (Marshaler, error)) {
	_marshalers[marshaler] = factory
}

func GetMarshaler(marshaler, args string) (Marshaler, error) {
	factory, ok := _marshalers[marshaler]
	if !ok {
		return nil, fmt.Errorf("Marshaler[%s] is not registered", marshaler)
	}

	return factory(args)
}

func RegisterSinkTarget(target string, factory func(string) (io.Writer, error)) {
	_sinkTargets[target] = factory
}

func GetSinkTarget(target, args string) (io.Writer, error) {
	factory, ok := _sinkTargets[target]
	if !ok {
		return nil, fmt.Errorf("SinkTarget[%s] is not registered", target)
	}

	return factory(args)
}
