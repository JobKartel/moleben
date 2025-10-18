package repository

import (
	"context"
	"errors"
	"moleben/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository interface {
	Create(ctx context.Context, s *domain.Session) error
	Get(ctx context.Context, id string) (*domain.Session, error)
	Update(ctx context.Context, s *domain.Session) error
	InsertMessage(ctx context.Context, m *domain.Message) error
	ListMessages(ctx context.Context, sessionID string, limit, offset int) ([]domain.Message, error)
}

type PostgresRepo struct{ Pool *pgxpool.Pool }

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo { return &PostgresRepo{Pool: pool} }

var ErrNotFound = errors.New("session not found")

func (r *PostgresRepo) Create(ctx context.Context, s *domain.Session) error {
	_, err := r.Pool.Exec(ctx, `INSERT INTO sessions(id, created_at, updated_at, status, tally, final_message) VALUES ($1,$2,$3,$4,$5,$6)`, s.ID, s.CreatedAt, s.UpdatedAt, string(s.Status), s.Tally, s.FinalMessage)
	return err
}

func (r *PostgresRepo) Get(ctx context.Context, id string) (*domain.Session, error) {
	row := r.Pool.QueryRow(ctx, `SELECT id, created_at, updated_at, status, tally, final_message FROM sessions WHERE id=$1`, id)
	s := domain.Session{}
	var status string
	if err := row.Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt, &status, &s.Tally, &s.FinalMessage); err != nil {
		return nil, ErrNotFound
	}
	s.Status = domain.SessionStatus(status)
	return &s, nil
}

func (r *PostgresRepo) Update(ctx context.Context, s *domain.Session) error {
	s.UpdatedAt = time.Now()
	_, err := r.Pool.Exec(ctx, `UPDATE sessions SET updated_at=$2, status=$3, tally=$4, final_message=$5 WHERE id=$1`, s.ID, s.UpdatedAt, string(s.Status), s.Tally, s.FinalMessage)
	return err
}

func (r *PostgresRepo) InsertMessage(ctx context.Context, m *domain.Message) error {
	_, err := r.Pool.Exec(ctx, `INSERT INTO messages(id, session_id, role, content, created_at) VALUES ($1,$2,$3,$4,$5)`, m.ID, m.SessionID, string(m.Role), m.Content, m.CreatedAt)
	return err
}

func (r *PostgresRepo) ListMessages(ctx context.Context, sessionID string, limit, offset int) ([]domain.Message, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.Pool.Query(ctx, `SELECT id, session_id, role, content, created_at FROM messages WHERE session_id=$1 ORDER BY created_at ASC LIMIT $2 OFFSET $3`, sessionID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Message, 0, limit)
	for rows.Next() {
		var m domain.Message
		var role string
		if err := rows.Scan(&m.ID, &m.SessionID, &role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.Role = domain.MessageRole(role)
		items = append(items, m)
	}
	return items, rows.Err()
}
