package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"testTaskHezzl/internal/config"
	"testTaskHezzl/internal/good"
	"testTaskHezzl/internal/meta"
	"testTaskHezzl/migrations"
)

const ( //sql-requests
	getListGoodsWithLimitAndOffset = "SELECT * FROM good LIMIT $1 OFFSET $2"
	insertGood                     = "INSERT INTO good(project_id, name, priority, removed, created_at) VALUES ($1, $2, COALESCE((SELECT MAX(priority) FROM good) + 1, 1) , false, NOW()) RETURNING *"
	checkGOOD                      = "SELECT * FROM good WHERE id = $1"
	updateGood                     = "UPDATE goods SET name = $1, description = $2 WHERE id = $3 AND project_id $4 RETURNING *"
	removeGood                     = "UPDATE good SET removed = true WHERE id = $1 RETURNING *"
	getListGoodsFromID             = "SELECT * FROM good WHERE id >= $1"
)

var (
	repriotiizeGoods = [3]string{"CREATE SEQUENCE update_priority START $1", "UPDATE good SET priority = NEXTVAL('update_priority') WHERE id >= $1", "DROP SEQUENCE update_priority	"}
)

type PostgresConn struct {
	conn *pgxpool.Pool
}

func NewConnection(ctx context.Context, cfg config.DBConfig, migrationsPath string) (*PostgresConn, error) {
	dbCFG, err := pgxpool.ParseConfig(cfg.String())
	if err != nil {
		return nil, err
	}
	conn, err := pgxpool.NewWithConfig(ctx, dbCFG)
	if err != nil {
		return nil, err
	}
	if err = migrations.Up(conn, migrationsPath); err != nil {
		return nil, err
	}
	return &PostgresConn{
		conn: conn,
	}, nil
}

func (pc *PostgresConn) CreateGood(ctx context.Context, projectId int, name string) (good.Good, error) {
	row := pc.conn.QueryRow(ctx, insertGood, projectId, name)
	g := good.Good{}
	if err := row.Scan(&g.ID, &g.ProjectID, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt); err != nil {
		return good.Good{}, err
	}
	return g, nil
}

func (pc *PostgresConn) GetListGoods(ctx context.Context, limit, offset int) (meta.Meta, []good.Good, error) {
	rows, err := pc.conn.Query(ctx, getListGoodsWithLimitAndOffset, limit, offset)
	if err != nil {
		return meta.Meta{}, nil, err
	}
	m := meta.Meta{
		Total:   0,
		Removed: 0,
		Limit:   limit,
		Offset:  offset,
	}
	result := make([]good.Good, 0, limit)
	for rows.Next() {
		g := good.Good{}
		if err = rows.Scan(&g.ID, &g.ProjectID, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt); err != nil {
			return meta.Meta{}, nil, fmt.Errorf("can't scan message: %w", err)
		}
		result = append(result, g)
		if g.Removed {
			m.Removed++
		}
		m.Total++
	}
	return m, result, nil
}

func (pc *PostgresConn) UpdateGood(ctx context.Context, id, projectId int, name, desc string) (good.Good, error) {
	row := pc.conn.QueryRow(ctx, updateGood, id, projectId, name, desc)
	g := good.Good{}
	if err := row.Scan(&g.ID, &g.ProjectID, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt); err != nil {
		return good.Good{}, err
	}
	return g, nil
}

func (pc *PostgresConn) RemoveGood(ctx context.Context, id, projectId int) (good.Good, error) {
	row := pc.conn.QueryRow(ctx, removeGood, id, projectId)
	g := good.Good{}
	if err := row.Scan(&g.ID, &g.ProjectID, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt); err != nil {
		return good.Good{}, nil
	}
	return g, nil
}

func (pc *PostgresConn) GoodIsExist(ctx context.Context, id, projectId int) (bool, error) {
	row := pc.conn.QueryRow(ctx, checkGOOD, id)
	g := good.Good{}
	if err := row.Scan(&g.ID, &g.ProjectID, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

func (pc *PostgresConn) ReprioritiizeGood(ctx context.Context, id int, projectID int, newPriority int) ([]good.Good, error) {
	cntRows, err := pc.UpdatePriorityTx(ctx, id, projectID, newPriority)
	if err != nil {
		return nil, fmt.Errorf("error update priority: %w", err)
	}
	rows, err := pc.conn.Query(ctx, getListGoodsFromID, id)
	if err != nil {
		return nil, fmt.Errorf("error select: %w", err)
	}
	result := make([]good.Good, 0, cntRows)
	for rows.Next() {
		g := good.Good{}
		if err = rows.Scan(&g.ID, &g.ProjectID, &g.Name, &g.Description, &g.Priority, &g.Removed, &g.CreatedAt); err != nil {
			return nil, fmt.Errorf("can't scan message: %w", err)
		}
		result = append(result, g)
	}
	return result, nil
}

func (pc *PostgresConn) UpdatePriorityTx(ctx context.Context, id, projectId, newPriority int) (int64, error) {
	tx, err := pc.conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, repriotiizeGoods[0], newPriority)
	ans, err := tx.Exec(ctx, repriotiizeGoods[1], id)
	_, err = tx.Exec(ctx, repriotiizeGoods[2])
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}
	return ans.RowsAffected(), nil
}
