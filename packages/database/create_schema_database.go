package database

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/content"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/page"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/space"
)

func CreateSchemaDatabase(db *pg.DB) error {
	models := []interface{}{
		(*space.Space)(nil),
		(*content.Content)(nil),
		(*page.Page)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:          false,
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
