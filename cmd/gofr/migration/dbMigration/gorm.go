package dbmigration

import (
	"strconv"
	"time"

	"gorm.io/gorm"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type GORM struct {
	database    *gorm.DB
	txn         *gorm.DB
	isTxnActive bool
}

type gofrMigration struct {
	App       string    `gorm:"primary_key"`
	Version   int64     `gorm:"primary_key;auto_increment:false"`
	StartTime time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	EndTime   time.Time `gorm:"default:NULL"`
	Method    string    `gorm:"primary_key"`
}

func NewGorm(d *gorm.DB) *GORM {
	return &GORM{database: d, txn: d.Begin(), isTxnActive: true}
}

func (g *GORM) Run(m Migrator, app, name, methods string, logger log.Logger) error {
	// if DOWN and UP both are required to run at once
	if !g.isTxnActive {
		g.txn = g.database.Begin()
		g.isTxnActive = true
	}

	err := g.preRun(app, methods, name)
	if err != nil {
		if g.txn != nil {
			g.rollBack()
		}

		return err
	}

	ds := &datastore.DataStore{ORM: g.database}

	if methods == UP {
		err = m.Up(ds, logger)
	} else {
		err = m.Down(ds, logger)
	}

	if err != nil {
		g.rollBack()
		return &errors.Response{Reason: "error encountered in running the migration", Detail: err}
	}

	err = g.postRun(app, methods, name)
	if err != nil {
		g.rollBack()
		return err
	}

	return nil
}

func (g *GORM) preRun(app, method, name string) error {
	if !g.database.Migrator().HasTable(&gofrMigration{}) {
		err := g.database.Migrator().CreateTable(&gofrMigration{}).Error
		if err != nil {
			return &errors.Response{Reason: "unable to create gofr_migrations table", Detail: err}
		}
	}

	if g.isDirty(app) {
		return &errors.Response{Reason: "dirty migration check failed"}
	}

	ver, _ := strconv.Atoi(name)

	err := g.txn.Create(&gofrMigration{App: app, Version: int64(ver), StartTime: time.Now(), Method: method}).Error
	if err != nil {
		return &errors.Response{Reason: "unable to insert migration start time", Detail: err}
	}

	return nil
}

func (g *GORM) isDirty(app string) bool {
	var val int64

	err := g.txn.Table("gofr_migrations").Where("app = ? AND end_time is null", app).Count(&val).Error
	if err != nil || val > 0 {
		return true
	}

	return false
}

func (g *GORM) postRun(app, method, name string) error {
	// finish the migration
	err := g.txn.Table("gofr_migrations").Where("app = ? AND version = ? AND method = ?", app, name, method).Update(`end_time`, &gofrMigration{EndTime: time.Now()}).Error
	//err := g.txn.Table("gofr_migrations").Where("app = $1 AND version = $2 AND method = $3", app, name, method).Update(&gofrMigration{EndTime: time.Now()}).Error
	return err
}

func (g *GORM) LastRunVersion(app, method string) (lv int) {
	row := g.database.Table("gofr_migrations").Where("app = ? AND method = ?", app, method).
		Select("COALESCE(MAX(version),0) as version").Row()

	_ = row.Scan(&lv)

	return
}

func (g *GORM) GetAllMigrations(app string) (upMigration, downMigration []int) {
	rows, err := g.database.Table("gofr_migrations").Where("app = ?", app).Select("version, method").Rows()
	if err != nil {
		return nil, nil
	}

	defer rows.Close()

	for rows.Next() {
		var (
			i int
			v string
		)

		_ = rows.Scan(&i, &v)

		if v == UP {
			upMigration = append(upMigration, i)
		} else {
			downMigration = append(downMigration, i)
		}
	}

	return
}

func (g *GORM) FinishMigration() error {
	if g.isTxnActive {
		g.isTxnActive = false
		return g.txn.Commit().Error
	}

	return nil
}

func (g *GORM) rollBack() {
	g.isTxnActive = false
	g.txn.Rollback()
}
