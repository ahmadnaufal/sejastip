package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"sejastip.id/api"
	"sejastip.id/api/entity"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type mysqlDevice struct {
	db *sqlx.DB
}

func NewMysqlDevice(db *sql.DB) api.DeviceRepository {
	newDB := sqlx.NewDb(db, "mysql")
	return &mysqlDevice{newDB}
}

func (m *mysqlDevice) GetUserDevice(ctx context.Context, userID int64) (*entity.Device, error) {
	query := `
		SELECT * FROM user_devices
		WHERE user_id = ?
		ORDER BY updated_at DESC
		LIMIT 1
	`
	result := &entity.Device{}
	err := m.db.GetContext(ctx, result, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrNotFound
		} else {
			return nil, err
		}
	}
	return result, nil
}

func (m *mysqlDevice) InsertUserDevice(ctx context.Context, device *entity.Device) error {
	now := time.Now()
	device.CreatedAt = now
	device.UpdatedAt = now

	query := `INSERT INTO user_devices
		(device_id, platform, user_agent, user_id, created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?, ?)
	`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	// execute query
	res, err := prep.ExecContext(ctx,
		device.DeviceID, device.Platform, device.UserAgent, device.UserID,
		device.CreatedAt, device.UpdatedAt,
	)
	if err != nil {
		return err
	}

	device.ID, err = res.LastInsertId()
	return err
}

func (m *mysqlDevice) UpdateUserDevice(ctx context.Context, ID int64, newDevice *entity.Device) error {
	now := time.Now()
	newDevice.UpdatedAt = now

	query := `
		UPDATE user_devices SET
		device_id = ?, platform = ?, user_agent = ?, updated_at = ?
		WHERE id = ?
	`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := prep.ExecContext(ctx,
		newDevice.DeviceID, newDevice.Platform, newDevice.UserAgent,
		newDevice.UpdatedAt, ID,
	)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		return errors.New(fmt.Sprintf("Unexpected behavior detected when updating device data (total rows affected: %d)", affectedRows))
	}

	return nil
}

func (m *mysqlDevice) RemoveDevice(ctx context.Context, ID int64) error {
	now := time.Now()
	query := `
		UPDATE user_devices SET
		deleted_at = ?
		WHERE id = ?
	`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := prep.ExecContext(ctx, now, ID)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		return errors.New(fmt.Sprintf("Unexpected behavior detected when deleting device (total rows affected: %d)", affectedRows))
	}

	return nil
}
