package index

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"time"

	erwp "github.com/PerumallaGiridhar/oolio/internal/errorwrap"
	"github.com/cockroachdb/pebble"
)

type PebbleStore struct {
	DB     *pebble.DB
	Txt    string
	DbDir  string
	Opened time.Time
}

func EnsurePebble(txtPath string, onErrs ...func()) *PebbleStore {
	dbDir := txtPath + ".peb"
	opts := &pebble.Options{FormatMajorVersion: pebble.FormatNewest}

	if hasManifest(dbDir) {
		db := erwp.MustReturn(erwp.Try(pebble.Open(dbDir, opts)), onErrs...)
		return &PebbleStore{DB: db, Txt: txtPath, DbDir: dbDir, Opened: time.Now()}
	}

	erwp.MustDo(os.MkdirAll(dbDir, 0o755), onErrs...)
	db := erwp.MustReturn(erwp.Try(pebble.Open(dbDir, opts)), onErrs...)
	bulkLoadTxtIntoPebble(db, txtPath, func() { db.Close() })
	erwp.MustDo(db.Flush(), func() { db.Close() })

	return &PebbleStore{DB: db, Txt: txtPath, DbDir: dbDir, Opened: time.Now()}
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

func bulkLoadTxtIntoPebble(db *pebble.DB, txtPath string, onErr ...func()) {
	f := erwp.MustReturn(erwp.Try(os.Open(txtPath)), onErr...)
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
		erwp.MustDo(batch.Set([]byte(code), nil, pebble.NoSync), onErr...)
		n++
		if n%rowsPerCommit == 0 {
			erwp.MustDo(batch.Commit(pebble.Sync), onErr...)
			batch = db.NewBatch()
		}
	}
	erwp.MustDo(sc.Err(), onErr...)
	erwp.MustDo(batch.Commit(pebble.Sync), onErr...)
}
