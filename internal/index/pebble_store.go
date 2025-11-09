package index

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/pebble"
)

type PebbleStore struct {
	DB     *pebble.DB
	Txt    string
	DbDir  string
	Opened time.Time
}

func EnsurePebble(txtPath string) (*PebbleStore, error) {
	dbDir := txtPath + ".peb"
	opts := &pebble.Options{FormatMajorVersion: pebble.FormatNewest}

	if hasManifest(dbDir) {
		log.Printf("found pebble indexes for %s", txtPath)
		db, err := pebble.Open(dbDir, opts)
		if err != nil {
			return nil, err
		}
		return &PebbleStore{DB: db, Txt: txtPath, DbDir: dbDir, Opened: time.Now()}, nil
	}

	log.Printf("Building pebble indexes for %s", txtPath)
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return nil, err
	}

	db, err := pebble.Open(dbDir, opts)
	if err != nil {
		return nil, err
	}

	if err := bulkLoadTxtIntoPebble(db, txtPath); err != nil {
		db.Close()
	}

	if err := db.Flush(); err != nil {
		db.Close()
	}

	return &PebbleStore{DB: db, Txt: txtPath, DbDir: dbDir, Opened: time.Now()}, nil
}

func (s *PebbleStore) Has(code string) (bool, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return false, nil
	}
	_, closer, err := s.DB.Get([]byte(code))
	if errors.Is(err, pebble.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if closer != nil {
		_ = closer.Close()
	}
	return true, nil
}

func (s *PebbleStore) Close() error { return s.DB.Close() }

func hasManifest(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "MANIFEST-") {
			return true
		}
	}
	return false
}

func bulkLoadTxtIntoPebble(db *pebble.DB, txtPath string) error {
	f, err := os.Open(txtPath)
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1024), 64*1024)

	const rowsPerCommit = 1_000_000
	batch := db.NewBatch()
	defer batch.Close()

	var n int
	for sc.Scan() {
		code := strings.ToUpper(strings.TrimSpace(sc.Text()))
		if code == "" {
			continue
		}
		if err := batch.Set([]byte(code), nil, pebble.NoSync); err != nil {
			return err
		}
		n++
		if n%rowsPerCommit == 0 {
			if err := batch.Commit(pebble.Sync); err != nil {
				return err
			}
			batch = db.NewBatch()
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}
	if err := batch.Commit(pebble.Sync); err != nil {
		return err
	}

	return nil
}
