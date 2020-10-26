package handlers

import (
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/techxmind/logserver/config"
	"github.com/techxmind/logserver/logger"
	"github.com/techxmind/logserver/storage"
	"github.com/techxmind/logserver/storage/kafka"
)

var (
	// For test mock
	StorageGet StorageGetter = newStorage
)

type StorageGetter func(*config.StorageConfig) (storage.Storager, error)

func newStorage(cfg *config.StorageConfig) (storage.Storager, error) {
	var group = storage.NewGroup()

	storageTypes := strings.Split(cfg.Types, ",")

	for _, storageType := range storageTypes {
		storageType = strings.TrimSpace(storageType)
		if storageType == "" {
			continue
		}

		switch strings.ToLower(storageType) {
		case "stdout":
			group.Add(storage.New(os.Stdout))
		case "kafka":
			if cfg.Kafka == nil {
				return nil, errors.New("Kafka storage configuration is missing")
			}

			if s, err := kafka.New(cfg.Kafka); err != nil {
				return nil, err
			} else {
				group.Add(s)
			}
		default:
			return nil, errors.Errorf("Unknow storage %s", storageType)
		}
	}

	if group.Size() == 0 {
		// default storage
		logger.Info("No storage specified, use STDOUT as default")
		group.Add(storage.New(os.Stdout))
	}

	return group, nil
}

func getKafkaStorage(cfg *config.KafkaConfig) (s storage.Storager, err error) {
	if cfg == nil {
		err = errors.New("Kafka storage configuration is missing")
		return
	}

	return kafka.New(cfg)
}
