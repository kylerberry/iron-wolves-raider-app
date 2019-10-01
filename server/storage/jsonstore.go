package storage

import "os"

// JSONStore represents a bungo response cache
type JSONStore struct {
	basePath string
	filename string
}

// NewJSONStore creates a store for a path
func NewJSONStore(filename string) JSONStore {
	return JSONStore{"/storage/", filename}
}

// func (s JSONStore) Read() []byte {

// }

func (s JSONStore) Write(data []byte) {
	f, err := os.Create(s.basePath + s.filename)
	if err != nil {
		panic(err)
	}
	f.Write(data)
}
