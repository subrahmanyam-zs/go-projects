package dbmigration

import (
	"io"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr/config"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

func initTests() *GORM {
	c := config.NewGoDotEnvProvider(log.NewMockLogger(io.Discard), "../../../../configs")

	database, _ := datastore.NewORM(&datastore.DBConfig{
		HostName: c.Get("DB_HOST"),
		Username: c.Get("DB_USER"),
		Password: c.Get("DB_PASSWORD"),
		Database: c.Get("DB_NAME"),
		Port:     c.Get("DB_PORT"),
		Dialect:  c.Get("DB_DIALECT"),
	})

	return &GORM{db: database.DB}
}

func Test_Run(t *testing.T) {
	g := initTests()

	defer func() {
		_ = g.db.Migrator().DropTable("gofr_migrations")
	}()

	tests := []struct {
		desc   string
		method string
		err    error
	}{
		{"success case", "UP", nil},
		{"failure case", "DOWN", &errors.Response{Reason: "error encountered in running the migration", Detail: errors.Error("test error")}},
	}

	for i, tc := range tests {
		err := g.Run(K20180324120906{}, "gofr-app", "20180324120906", tc.method, log.NewMockLogger(io.Discard))

		assert.Equal(t, tc.err, err, "TEST[%d], failed.\n%s", i, tc.desc)
	}
}

func Test_preRun(t *testing.T) {
	g := initTests()

	now := time.Now().UTC()
	g.txn = g.db.Begin()

	createTable(t, g.db)

	defer func() {
		_ = g.txn.Migrator().DropTable("gofr_migrations")
	}()

	insertMigration(t, g.txn, &gofrMigration{App: "gofr-app", Version: int64(20180324120906), StartTime: now, EndTime: now, Method: "UP"})

	expErr := &errors.Response{
		Reason: "unable to insert migration start time",
		Detail: &mysql.MySQLError{
			Number:  uint16(1062),
			Message: "Duplicate entry 'gofr-app-20180324120906-UP' for key 'gofr_migrations.PRIMARY'",
		},
	}

	err := g.preRun("gofr-app", "UP", "20180324120906")

	assert.Equal(t, expErr, err, "TEST failed.\n%s", "failure in starttime insertion (gofr_migrations table)")
}

func Test_postRun(t *testing.T) {
	g := initTests()
	g.txn = g.db.Begin()

	err := g.postRun("gofr-app", "UP", "20180324120906")

	assert.Error(t, err, "TEST failed.\n%s", "failure in endtime update (gofr_migrations table)")
}

func Test_isDirty(t *testing.T) {
	g := initTests()

	expErr := &errors.Response{Reason: "dirty migration check failed"}

	createTable(t, g.db)

	defer func() {
		_ = g.db.Migrator().DropTable("gofr_migrations")
	}()

	insertMigration(t, g.db, &gofrMigration{App: "gofr-app", Version: int64(20180324120906), StartTime: time.Now().UTC(), Method: "UP"})

	err := g.Run(K20180324120906{}, "gofr-app", "20180324120906", "UP", log.NewMockLogger(io.Discard))

	assert.Equal(t, expErr, err, "TEST failed.\n%s", "dirty migration check failure case")
}

func Test_LastRunVersion(t *testing.T) {
	g := initTests()

	now := time.Now().UTC()
	expLastVersion := 20180324120906

	createTable(t, g.db)

	defer func() {
		_ = g.db.Migrator().DropTable("gofr_migrations")
	}()

	insertMigration(t, g.db, &gofrMigration{App: "gofr-app", Version: int64(20180324120906), StartTime: now, EndTime: now, Method: "UP"})

	lastVersion := g.LastRunVersion("gofr-app", "UP")

	assert.Equal(t, expLastVersion, lastVersion, "TEST failed.\n%s", "last version check")
}

func Test_GetAllMigrations(t *testing.T) {
	g := initTests()

	now := time.Now().UTC()
	desc := "get all migrations"

	createTable(t, g.db)

	defer func() {
		_ = g.db.Migrator().DropTable("gofr_migrations")
	}()

	insertMigration(t, g.db, &gofrMigration{App: "gofr-app", Version: int64(20180324120906), StartTime: now, EndTime: now, Method: "UP"})
	insertMigration(t, g.db, &gofrMigration{App: "gofr-app", Version: int64(20180324120906), StartTime: now, EndTime: now, Method: "DOWN"})

	expOut := []int{20180324120906}

	up, down := g.GetAllMigrations("gofr-app")

	assert.Equal(t, expOut, up, "TEST failed.\n%s", desc)

	assert.Equal(t, expOut, down, "TEST failed.\n%s", desc)
}

func Test_GetAllMigrationsError(t *testing.T) {
	g := initTests()

	createMockTable(t, g.db)

	defer func() {
		_ = g.db.Migrator().DropTable("gofr_migrations")
	}()

	up, down := g.GetAllMigrations("gofr-app")

	assert.Nil(t, up, "TEST failed.\n%s", "get all migrations")

	assert.Nil(t, down, "TEST failed.\n%s", "get all migrations")
}

type K20180324120906 struct{}

func (k K20180324120906) Up(d *datastore.DataStore, l log.Logger) error {
	return nil
}

func (k K20180324120906) Down(d *datastore.DataStore, l log.Logger) error {
	return errors.Error("test error")
}

func insertMigration(t *testing.T, g *gorm.DB, mig *gofrMigration) {
	err := g.Create(mig).Error
	if err != nil {
		t.Error(err)
	}
}

func createTable(t *testing.T, g *gorm.DB) {
	err := g.Migrator().CreateTable(&gofrMigration{})
	if err != nil {
		t.Error(err)
	}
}

func createMockTable(t *testing.T, g *gorm.DB) {
	type gofrMigration struct {
		App string `gorm:"primary_key"`
	}

	err := g.Migrator().CreateTable(&gofrMigration{})
	if err != nil {
		t.Error(err)
	}
}
