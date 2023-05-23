package godem

import (
	"archive/zip"

	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type SrtmLocalStorage interface {
	// LoadFile loads a file, if not available, then err!=nil and IsNotExists(err) must be true
	FileExists(dem, zip, fn string) (string, error)
	LoadFile(dem, zip, fn string) ([]byte, error)
	IsNotExists(err error) bool
	SaveFile(dem, zip, fn string, bytes []byte) error
	Unzip(dem, fn string, data []byte) error
}

type LocalFileSrtmStorage struct {
	cacheDirectory string
	source         int
}

var _ SrtmLocalStorage = new(LocalFileSrtmStorage)

func makeDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModeDir|0700); err != nil {
			return err
		}
	}
	return nil
}

func NewLocalFileSrtmStorage(source int) (*LocalFileSrtmStorage, error) {
	var cacheDirectory string
	if source == SOURCE_VIEW {
		cacheDirectory = filepath.Join(os.Getenv("HOME"), ".cache", "godem", "viewfinderpanoramas")
	}
	if source == SOURCE_GPXSEE {
		cacheDirectory = filepath.Join(os.Getenv("HOME"), ".local", "share", "gpxsee", "DEM")
	}

	if err := makeDir(cacheDirectory); err != nil {
		return nil, err
	}

	return &LocalFileSrtmStorage{cacheDirectory: cacheDirectory, source: source}, nil
}

func (ds LocalFileSrtmStorage) getName(dem, zip, fn string) string {
	if ds.source == SOURCE_VIEW {
		return filepath.Join(ds.cacheDirectory, dem, zip, fn)
	}
	return filepath.Join(ds.cacheDirectory, fn)
}

func (ds LocalFileSrtmStorage) LoadFile(dem, zip, fn string) ([]byte, error) {
	f, err := os.Open(ds.getName(dem, zip, fn))
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (ds LocalFileSrtmStorage) FileExists(dem, zip, fn string) (string, error) {
	f := path.Join(ds.getName(dem, zip, fn))
	fi, err := os.Stat(f)
	if err != nil {
		return "", err
	}
	if fi.Size() == 0 {
		return "", os.ErrNotExist
	}
	return f, nil
}

func (ds LocalFileSrtmStorage) IsNotExists(err error) bool {
	return os.IsNotExist(err)
}

func (ds LocalFileSrtmStorage) SaveFile(dem, zip, fn string, bytes []byte) error {
	if err := makeDir(ds.getName(dem, zip, "")); err != nil {
		return err
	}

	f, err := os.Create(ds.getName(dem, zip, fn))
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	return err
}

func (ds LocalFileSrtmStorage) Unzip(dem, fn string, data []byte) error {
	r := bytes.NewReader(data)
	zipReader, err := zip.NewReader(r, int64(len(data)))
	if err != nil {
		return err
	}

	for _, f := range zipReader.File {
		var filePath string
		if ds.source == SOURCE_VIEW {
			filePath = ds.getName(dem, f.Name, "")
		} else {
			filePath = ds.getName("", "", f.Name)
		}

		if !strings.HasPrefix(filePath, filepath.Clean(ds.cacheDirectory)+string(os.PathSeparator)) {
			continue
		}
		if f.FileInfo().IsDir() {
			continue
		}

		if err := makeDir(filepath.Dir(filePath)); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {

			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}
