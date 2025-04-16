package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/51mans0n/avito-pvz-task/internal/db"
	"github.com/51mans0n/avito-pvz-task/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestRepo_CreatePVZ(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	xdb := sqlx.NewDb(sqlDB, "postgres")
	repo := db.NewRepo(xdb)

	mock.ExpectExec("INSERT INTO pvz").
		WithArgs("some-uuid", "Москва", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreatePVZ(context.Background(), &model.PVZ{
		ID:   "some-uuid",
		City: "Москва",
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_CreatePVZ_Error(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	xdb := sqlx.NewDb(sqlDB, "postgres")
	repo := db.NewRepo(xdb)

	mock.ExpectExec("INSERT INTO pvz").WillReturnError(context.DeadlineExceeded)

	err = repo.CreatePVZ(context.Background(), &model.PVZ{
		ID:   "fail-uuid",
		City: "Спб",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "deadline exceeded")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_CreateReception_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	xdb := sqlx.NewDb(sqlDB, "postgres")
	repo := db.NewRepo(xdb)

	mock.ExpectQuery(`SELECT count\(\*\) FROM receptions`).
		WithArgs("82cc7cda-bd24-468f-b7b7-844d66b6693c", "in_progress").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectExec(`INSERT INTO receptions \(id,pvz_id,date_time,status\)`).
		WithArgs("rec-111", "82cc7cda-bd24-468f-b7b7-844d66b6693c", sqlmock.AnyArg(), "in_progress").
		WillReturnResult(sqlmock.NewResult(1, 1))

	rec := &model.Reception{
		ID:       "rec-111",
		PVZID:    "82cc7cda-bd24-468f-b7b7-844d66b6693c",
		DateTime: time.Now(),
		Status:   "in_progress",
	}
	err = repo.CreateReception(context.Background(), rec)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_CreateReception_AlreadyOpen(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	xdb := sqlx.NewDb(sqlDB, "postgres")
	repo := db.NewRepo(xdb)

	mock.ExpectQuery(`SELECT count\(\*\) FROM receptions`).
		WithArgs("82cc7cda-bd24-468f-b7b7-844d66b6693c", "in_progress").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	rec := &model.Reception{
		ID:     "rec-222",
		PVZID:  "82cc7cda-bd24-468f-b7b7-844d66b6693c",
		Status: "in_progress",
	}

	err = repo.CreateReception(context.Background(), rec)
	require.Error(t, err)
	require.Contains(t, err.Error(), "already an open reception")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_CreateProduct_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()
	xdb := sqlx.NewDb(sqlDB, "postgres")
	repo := db.NewRepo(xdb)

	mock.ExpectQuery(`SELECT id, pvz_id, date_time, status FROM receptions`).
		WithArgs("82cc7cda-bd24-468f-b7b7-844d66b6693c", "in_progress").
		WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "date_time", "status"}).
			AddRow("rec-active", "82cc7cda-bd24-468f-b7b7-844d66b6693c", time.Now(), "in_progress"))

	mock.ExpectExec(`INSERT INTO products \(id,reception_id,date_time,type\)`).
		WithArgs("prod-xyz", "rec-active", sqlmock.AnyArg(), "электроника").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateProduct(context.Background(), "82cc7cda-bd24-468f-b7b7-844d66b6693c", &model.Product{
		ID:       "prod-xyz",
		Type:     "электроника",
		DateTime: time.Now(),
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_CreateProduct_NoActiveReception(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()
	xdb := sqlx.NewDb(sqlDB, "postgres")
	repo := db.NewRepo(xdb)

	mock.ExpectQuery(`SELECT id, pvz_id, date_time, status FROM receptions`).
		WithArgs("82cc7cda-bd24-468f-b7b7-844d66b6693c", "in_progress").
		WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "date_time", "status"})) // empty

	err = repo.CreateProduct(context.Background(), "82cc7cda-bd24-468f-b7b7-844d66b6693c", &model.Product{
		ID:       "prod-abc",
		Type:     "обувь",
		DateTime: time.Now(),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no active reception found")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_DeleteLastProduct_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()
	xdb := sqlx.NewDb(sqlDB, "postgres")
	repo := db.NewRepo(xdb)

	mock.ExpectQuery(`SELECT id, pvz_id, date_time, status FROM receptions`).
		WithArgs("82cc7cda-bd24-468f-b7b7-844d66b6693c", "in_progress").
		WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "date_time", "status"}).
			AddRow("rec-xxx", "82cc7cda-bd24-468f-b7b7-844d66b6693c", time.Now(), "in_progress"))

	mock.ExpectQuery(`SELECT id FROM products`).
		WithArgs("rec-xxx").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("prod-latest"))

	mock.ExpectExec(`DELETE FROM products`).
		WithArgs("prod-latest").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteLastProduct(context.Background(), "82cc7cda-bd24-468f-b7b7-844d66b6693c")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_DeleteLastProduct_NoActiveReception(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()
	xdb := sqlx.NewDb(sqlDB, "postgres")
	repo := db.NewRepo(xdb)

	mock.ExpectQuery(`SELECT id, pvz_id, date_time, status FROM receptions`).
		WithArgs("82cc7cda-bd24-468f-b7b7-844d66b6693c", "in_progress").
		WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "date_time", "status"}))

	err = repo.DeleteLastProduct(context.Background(), "82cc7cda-bd24-468f-b7b7-844d66b6693c")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no active reception")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_DeleteLastProduct_NoProducts(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()
	xdb := sqlx.NewDb(sqlDB, "postgres")
	repo := db.NewRepo(xdb)

	mock.ExpectQuery(`SELECT id, pvz_id, date_time, status FROM receptions`).
		WithArgs("82cc7cda-bd24-468f-b7b7-844d66b6693c", "in_progress").
		WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "date_time", "status"}).
			AddRow("rec-abc", "82cc7cda-bd24-468f-b7b7-844d66b6693c", time.Now(), "in_progress"))

	mock.ExpectQuery(`SELECT id FROM products`).
		WithArgs("rec-abc").
		WillReturnRows(sqlmock.NewRows([]string{"id"})) // no rows

	err = repo.DeleteLastProduct(context.Background(), "82cc7cda-bd24-468f-b7b7-844d66b6693c")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no products to delete")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_CloseLastReception_Success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()
	xdb := sqlx.NewDb(sqlDB, "postgres")
	repo := db.NewRepo(xdb)

	mock.ExpectQuery(`SELECT id, pvz_id, date_time, status FROM receptions`).
		WithArgs("82cc7cda-bd24-468f-b7b7-844d66b6693c", "in_progress").
		WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "date_time", "status"}).
			AddRow("rec-xyz", "82cc7cda-bd24-468f-b7b7-844d66b6693c", time.Now(), "in_progress"))

	mock.ExpectExec(`UPDATE receptions SET status = \$1 WHERE id = \$2`).
		WithArgs("close", "rec-xyz").
		WillReturnResult(sqlmock.NewResult(0, 1))

	rec, err := repo.CloseLastReception(context.Background(), "82cc7cda-bd24-468f-b7b7-844d66b6693c")
	require.NoError(t, err)
	require.Equal(t, "close", rec.Status)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_CloseLastReception_NoActive(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()
	xdb := sqlx.NewDb(sqlDB, "postgres")
	repo := db.NewRepo(xdb)

	mock.ExpectQuery(`SELECT id, pvz_id, date_time, status FROM receptions`).
		WithArgs("82cc7cda-bd24-468f-b7b7-844d66b6693c", "in_progress").
		WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "date_time", "status"}))

	rc, err := repo.CloseLastReception(context.Background(), "82cc7cda-bd24-468f-b7b7-844d66b6693c")
	require.Nil(t, rc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no active reception")

	require.NoError(t, mock.ExpectationsWereMet())
}
