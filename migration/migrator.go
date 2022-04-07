package migration

import (
	"context"
	"gorm.io/gorm"
	"sort"
	"sync"
)

var m = migrator{scripts: make(map[string][]Script)}

type migrator struct {
	sync.Mutex
	db      *gorm.DB
	scripts map[string][]Script
}

func Init(db *gorm.DB) {
	m.db = db
}

func (m *migrator) register(scripts ...Script) {
	m.Lock()
	defer m.Unlock()
	for _, script := range scripts {
		m.scripts[script.Owner()] = append(m.scripts[script.Owner()], script)
	}
}

func (m *migrator) bookKeep(script Script) error {
	record := &MigrationHistory{
		ScriptOwner:   script.Owner(),
		ScriptVersion: script.Version(),
		ScriptComment: script.Comment(),
	}
	return m.db.Create(record).Error
}

func (m *migrator) execute(ctx context.Context) error {
	versions, err := m.getLastVersion()
	if err != nil {
		return err
	}
	for _, scriptSlice := range m.scripts {
		sort.Slice(scriptSlice, func(i, j int) bool {
			return scriptSlice[i].Version() < scriptSlice[j].Version()
		})
		for _, script := range scriptSlice {
			if script.Version() > versions[script.Owner()] {
				err = script.Up(ctx, m.db)
				if err != nil {
					return err
				}
				err = m.bookKeep(script)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func (m *migrator) getLastVersion() (map[string]uint64, error) {
	var err error
	versions := make(map[string]uint64)
	if !m.db.Migrator().HasTable(tableName) {
		return versions, nil
	}
	var records []MigrationHistory
	err = m.db.Find(&records).Error
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		if record.ScriptVersion > versions[record.ScriptOwner] {
			versions[record.ScriptOwner] = record.ScriptVersion
		}
	}
	return versions, nil
}

func Register(script ...Script) {
	m.register(script...)
}

func Execute(ctx context.Context) error {
	return m.execute(ctx)
}
