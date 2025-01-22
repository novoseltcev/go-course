package server

import (
	"encoding/json"
	"os"

	"github.com/spf13/afero"

	"github.com/novoseltcev/go-course/internal/storages"
)

func Backup(fs afero.Fs, path string, storager storages.MetricStorager) error {
	if path == "" {
		return nil
	}

	fd, err := fs.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}
	defer fd.Close()

	return json.NewEncoder(fd).Encode(storager)
}
