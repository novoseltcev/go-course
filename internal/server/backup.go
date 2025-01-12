package server

import (
	"context"
	"encoding/json"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/novoseltcev/go-course/internal/storages"
)

func Backup(path string, storager storages.MetricStorager) error {
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}

	defer fd.Close()

	if err := json.NewEncoder(fd).Encode(storager); err != nil {
		return err
	}

	return nil
}

func BackupWorker(ctx context.Context, interval time.Duration, path string, storager storages.MetricStorager) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(interval)
		}

		if err := Backup(path, storager); err != nil {
			log.WithError(err).Error("failed to backup metrics")
		}
	}
}
