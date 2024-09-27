package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type TweetModel struct {
	DB *sql.DB
}

type Tweet struct {
	ID        int64
	UserID    int64
	Content   string
	CreatedAt time.Time
}

func (m TweetModel) Insert(tweet *Tweet) (*Tweet, error) {
	query := `
        INSERT INTO tweets (user_id, content) 
        VALUES ($1, $2)
        RETURNING id, created_at`
	args := []interface{}{tweet.UserID, tweet.Content}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&tweet.ID, &tweet.CreatedAt)
	if err != nil {
		return nil, err
	}

	return tweet, nil
}

func (m TweetModel) GetTweetByID(id int64) (*Tweet, error) {
	query := `
        SELECT id, user_id, content, created_at 
        FROM tweets 
        WHERE id = $1`

	var tweet Tweet

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&tweet.ID,
		&tweet.UserID,
		&tweet.Content,
		&tweet.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &tweet, nil
}

func (m TweetModel) Update(tweet *Tweet) (error, *Tweet) {
	query := `
        UPDATE tweets
        SET content = $1
        WHERE id = $2 
        RETURNING created_at`

	args := []interface{}{tweet.Content, tweet.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&tweet.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound, nil
		}
		return err, nil
	}

	return nil, tweet
}

func (m TweetModel) Delete(id int64) error {
	query := `
        DELETE FROM tweets 
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (m TweetModel) GetTweetByUserID(userID int64) ([]*Tweet, error) {
	query := `
        SELECT id, user_id, content, created_at 
        FROM tweets
        WHERE user_id = $1
        ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []*Tweet
	for rows.Next() {
		var tweet Tweet
		err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Content, &tweet.CreatedAt)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tweets, nil
}
