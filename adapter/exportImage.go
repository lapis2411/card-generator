package adapter

import (
	"fmt"
	"path/filepath"

	"github.com/lapis2411/card-generator/domain"
)

const (
	ImageType = "png"
)

type (
	Exporter struct {
		folder       string
		exportDriver ExportDriver
	}

	ExportDriver interface {
		Save(domain.Image, string) error
	}
)

func NewExporter(folder string, e ExportDriver) domain.ExportImageAdapter {
	return Exporter{
		folder:       folder,
		exportDriver: e,
	}
}

func (e Exporter) Save(imgs []domain.Image) error {
	for _, img := range imgs {
		f := img.Name() + "." + ImageType
		p := filepath.Join(e.folder, f)
		if err := e.exportDriver.Save(img, p); err != nil {
			return fmt.Errorf("failed to save image(%s): %w", p, err)
		}
	}
	return nil
}
