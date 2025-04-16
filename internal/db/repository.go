package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/51mans0n/avito-pvz-task/internal/model"
)

// Repository - интерфейс, который нужен хендлерам (pvz, receptions...)
// чтобы не зависеть от конкретной *Repo
type Repository interface {
	CreatePVZ(ctx context.Context, pvz *model.PVZ) error
	GetPVZListWithFilter(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]model.PVZWithReceptions, error)
	CreateReception(ctx context.Context, rec *model.Reception) error
	CreateProduct(ctx context.Context, pvzID string, prod *model.Product) error
	DeleteLastProduct(ctx context.Context, pvzID string) error
	CloseLastReception(ctx context.Context, pvzID string) (*model.Reception, error)
}

// Убедимся, что *Repo реализует Repository:
var _ Repository = (*Repo)(nil)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreatePVZ(ctx context.Context, pvz *model.PVZ) error {
	query, args, err := sq.Insert("pvz").
		Columns("id", "city", "registration_date").
		Values(pvz.ID, pvz.City, pvz.RegistrationDate).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repo) GetPVZListWithFilter(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]model.PVZWithReceptions, error) {
	q := sq.Select("id", "city", "registration_date").
		From("pvz").
		OrderBy("registration_date DESC").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		PlaceholderFormat(sq.Dollar)

	sqlPVZ, argsPVZ, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var pvzRows []struct {
		ID               string    `db:"id"`
		City             string    `db:"city"`
		RegistrationDate time.Time `db:"registration_date"`
	}
	err = r.db.SelectContext(ctx, &pvzRows, sqlPVZ, argsPVZ...)
	if err != nil {
		return nil, err
	}

	result := make([]model.PVZWithReceptions, 0, len(pvzRows))
	for _, row := range pvzRows {
		item := model.PVZWithReceptions{
			PVZ: &model.PVZResponse{
				ID:               row.ID,
				City:             row.City,
				RegistrationDate: row.RegistrationDate,
			},
			Receptions: []model.ReceptionWithProd{},
		}

		recs, err := r.getReceptions(ctx, row.ID, startDate, endDate)
		if err != nil {
			return nil, err
		}
		rwp := make([]model.ReceptionWithProd, 0, len(recs))
		for _, rc := range recs {
			prods, err := r.getProducts(ctx, rc.ID)
			if err != nil {
				return nil, err
			}
			rwp = append(rwp, model.ReceptionWithProd{
				Reception: &model.ReceptionResponse{
					ID:       rc.ID,
					PVZID:    rc.PVZID,
					DateTime: rc.DateTime,
					Status:   rc.Status,
				},
				Products: convertProducts(prods),
			})
		}
		item.Receptions = rwp

		result = append(result, item)
	}
	return result, nil
}

func (r *Repo) CreateReception(ctx context.Context, rec *model.Reception) error {
	var countOpen int
	qCheck := sq.Select("count(*)").From("receptions").
		Where(sq.Eq{"pvz_id": rec.PVZID, "status": "in_progress"}).
		PlaceholderFormat(sq.Dollar)

	sqlCheck, argsCheck, err := qCheck.ToSql()
	if err != nil {
		return err
	}
	if err := r.db.GetContext(ctx, &countOpen, sqlCheck, argsCheck...); err != nil {
		return err
	}
	if countOpen > 0 {
		return fmt.Errorf("there is already an open reception")
	}

	qIns, argsIns, err := sq.Insert("receptions").
		Columns("id", "pvz_id", "date_time", "status").
		Values(rec.ID, rec.PVZID, rec.DateTime, rec.Status).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, qIns, argsIns...)
	return err
}

func (r *Repo) CreateProduct(ctx context.Context, pvzID string, prod *model.Product) error {
	rec, err := r.getActiveReception(ctx, pvzID)
	if err != nil {
		return err
	}
	if rec == nil {
		return fmt.Errorf("no active reception found for pvz %s", pvzID)
	}

	prod.ReceptionID = rec.ID
	q, args, err := sq.Insert("products").
		Columns("id", "reception_id", "date_time", "type").
		Values(prod.ID, prod.ReceptionID, prod.DateTime, prod.Type).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, q, args...)
	return err
}

func (r *Repo) DeleteLastProduct(ctx context.Context, pvzID string) error {
	rec, err := r.getActiveReception(ctx, pvzID)
	if err != nil {
		return err
	}
	if rec == nil {
		return fmt.Errorf("no active reception found for pvz %s", pvzID)
	}

	qSel, argsSel, err := sq.Select("id").
		From("products").
		Where(sq.Eq{"reception_id": rec.ID}).
		OrderBy("date_time DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	var prodID string
	if err := r.db.GetContext(ctx, &prodID, qSel, argsSel...); err != nil {
		if isNoRowsErr(err) {
			return errors.New("no products to delete")
		}
		return err
	}

	qDel, argsDel, err := sq.Delete("products").
		Where(sq.Eq{"id": prodID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, qDel, argsDel...)
	return err
}

func (r *Repo) CloseLastReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	rec, err := r.getActiveReception(ctx, pvzID)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, fmt.Errorf("no active reception found")
	}

	qUp, argsUp, err := sq.Update("receptions").
		Set("status", "close").
		Where(sq.Eq{"id": rec.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	_, err = r.db.ExecContext(ctx, qUp, argsUp...)
	if err != nil {
		return nil, err
	}
	rec.Status = "close"
	return rec, nil
}

func (r *Repo) getActiveReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	q := sq.Select("id", "pvz_id", "date_time", "status").
		From("receptions").
		Where(sq.Eq{"pvz_id": pvzID, "status": "in_progress"}).
		OrderBy("date_time DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var rec model.Reception
	if err := r.db.GetContext(ctx, &rec, sqlStr, args...); err != nil {
		if isNoRowsErr(err) {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

func (r *Repo) getReceptions(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]*model.Reception, error) {
	q := sq.Select("id", "pvz_id", "date_time", "status").
		From("receptions").
		Where(sq.Eq{"pvz_id": pvzID}).
		PlaceholderFormat(sq.Dollar)

	if startDate != nil {
		q = q.Where(sq.GtOrEq{"date_time": *startDate})
	}
	if endDate != nil {
		q = q.Where(sq.LtOrEq{"date_time": *endDate})
	}
	q = q.OrderBy("date_time DESC")

	sqlRec, argsRec, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	var recs []*model.Reception
	if err := r.db.SelectContext(ctx, &recs, sqlRec, argsRec...); err != nil {
		return nil, err
	}
	return recs, nil
}

func (r *Repo) getProducts(ctx context.Context, receptionID string) ([]*model.Product, error) {
	q := sq.Select("id", "reception_id", "date_time", "type").
		From("products").
		Where(sq.Eq{"reception_id": receptionID}).
		OrderBy("date_time DESC").
		PlaceholderFormat(sq.Dollar)

	sqlProd, argsProd, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	var prods []*model.Product
	if err := r.db.SelectContext(ctx, &prods, sqlProd, argsProd...); err != nil {
		return nil, err
	}
	return prods, nil
}

func convertProducts(ps []*model.Product) []model.ProductResponse {
	result := make([]model.ProductResponse, 0, len(ps))
	for _, p := range ps {
		result = append(result, model.ProductResponse{
			ID:          p.ID,
			DateTime:    p.DateTime,
			Type:        p.Type,
			ReceptionID: p.ReceptionID,
		})
	}
	return result
}

func isNoRowsErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "no rows in result set")
}
