package json

import (
	"errors"
	"os"
	"path/filepath"

	scribble "github.com/nanobox-io/golang-scribble"
)

type DBDriver interface {
	Write(string, string, any) error
	Read(string, string, any) error
	ReadAll(string) ([]string, error)
	Delete(string, string) error
}

type dbDriver struct {
	db  *scribble.Driver
	dir string
}

func newDb(p string) (*dbDriver, error) {
	d := &dbDriver{}
	db, err := scribble.New(p, nil)
	if err != nil {
		return d, err
	}
	d.db = db
	d.dir = p
	return d, nil
}

func (d *dbDriver) Write(collection, resource string, v any) error {
	return d.db.Write(collection, resource, v)
}

func (d *dbDriver) Read(collection, resource string, v any) error {
	return d.db.Read(collection, resource, v)
}

func (d *dbDriver) ReadAll(collection string) ([]string, error) {
	var records []string

	if collection == "" {
		return nil, errors.New("collection is missing")
	}

	dir := filepath.Join(d.dir, collection)
	if _, err := os.Stat(d.dir); os.IsNotExist(err) {
		return records, err
	}

	f, err := os.Open(dir)
	if err != nil {
		return records, errors.New("unable to read collection directory")
	}
	defer f.Close()

	items, err := f.ReadDir(0)
	if err != nil {
		return records, err
	}

	for _, item := range items {
		if item.IsDir() || filepath.Ext(item.Name()) != ".json" {
			continue
		}

		b, err := os.ReadFile(filepath.Join(dir, item.Name()))
		if err != nil {
			return records, err
		}

		records = append(records, string(b))
	}

	return records, nil
}

func (d *dbDriver) Delete(collection, resource string) error {
	return d.db.Delete(collection, resource)
}
