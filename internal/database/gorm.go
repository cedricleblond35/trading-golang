package database

import (
	"context"
	"os"

	"trading/internal/config"
	"trading/internal/model"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GORM implements Client.
type GORM struct {
	db *gorm.DB
}

func GormOpen(ctx context.Context, debugSQL bool) (*GORM, error) {
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PWD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + config.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Enables SQL queries logging
	if debugSQL && os.Getenv("FORCE_DISABLE_SQL_DEBUG") == "" {
		db = db.Debug()
	}

	return &GORM{
		db: db,
	}, nil
}

// Create implements Client.
func (c *GORM) Create(m model.Model) error {
	return c.db.Create(m).Error
}

// Create implements Client if it isn't exist
func  (c *GORM) FirstOrCreate(m model.Model) error {
	return c.db.FirstOrCreate(m).Error
}

// Delete implements Client.
func (c *GORM) Delete(m model.Model) error {
	return c.db.Delete(m).Error
}

// IsNotFound implements Client.
func (*GORM) IsNotFound(err error) bool {
	return errors.Cause(err) == gorm.ErrRecordNotFound
}

// Load implements Client.
func (c *GORM) Load(m model.Model, query string, args ...any) error {
	return c.db.Where(query, args...).Take(m).Error
}

// Load implements Client.
func (c *GORM) LoadLast(m model.Model, query string, args ...any) error {
	return c.db.Where(query, args...).Last(m).Error
}

// Loads implements Client.
func (c *GORM) Loads(m any, query string, args ...any) error {
	return c.db.Where(query, args...).Find(m).Error
}

// Ping implements Client.
func (c *GORM) Ping() error {
	if c.db == nil {
		return errors.New("no active database connection")
	}
	db, err := c.db.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

// Save implements Client.
func (c *GORM) Save(m model.Model) error {
	return c.db.Save(m).Error
}

// Update implements Client.
func (c *GORM) Update(m model.Model, fields map[string]any) error {
	err := c.db.Model(m).Updates(fields).Error
	return err
}

func (c *GORM) DB() *gorm.DB {
	return c.db
}
