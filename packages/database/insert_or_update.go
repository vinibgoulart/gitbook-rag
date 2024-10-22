package database

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func InsertOrUpdate(ctx *context.Context, db *bun.DB, model interface{}, conflictColumn string, updateFields string) error {
	_, err := db.NewInsert().
		Model(model).
		On("CONFLICT (" + conflictColumn + ") DO UPDATE").
		Set(updateFields).
		Exec(*ctx)

	if err != nil {
		return fmt.Errorf("error inserting or updating: %w", err)
	}

	return nil
}
