package database

import (
	"fmt"

	"github.com/go-pg/pg"
)

func InsertOrUpdate(db *pg.DB, model interface{}, conflictColumn string, updateFields string) error {
	_, err := db.Model(model).
		OnConflict(fmt.Sprintf("(%s) DO UPDATE", conflictColumn)).
		Set(updateFields).
		Insert()

	if err != nil {
		return fmt.Errorf("error inserting or updating: %w", err)
	}

	return nil
}
