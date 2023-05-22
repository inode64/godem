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

func NewLocalFileSrtmStorage(cacheDirectory string) (*LocalFileSrtmStorage, error) {
	if len(cacheDirectory) == 0 {
		cacheDirectory = path.Join(os.Getenv("HOME"), ".cache", "godem")
	}

	if err := makeDir(cacheDirectory); err != nil {
		return nil, err
	}

	return &LocalFileSrtmStorage{cacheDirectory: cacheDirectory}, nil
}

func (ds LocalFileSrtmStorage) LoadFile(dem, zip, fn string) ([]byte, error) {
	f, err := os.Open(path.Join(ds.cacheDirectory, dem, zip, fn))
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
	f := path.Join(ds.cacheDirectory, dem, zip, fn)
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
	if err := makeDir(path.Join(ds.cacheDirectory, dem)); err != nil {
		return err
	}
	if err := makeDir(path.Join(ds.cacheDirectory, dem, zip)); err != nil {
		return err
	}
	f, err := os.Create(path.Join(ds.cacheDirectory, dem, zip, fn))
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
		filePath := filepath.Join(ds.cacheDirectory, dem, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(ds.cacheDirectory)+string(os.PathSeparator)) {
			continue
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				panic(err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {

			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}
