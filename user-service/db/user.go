package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fgouvea/weather/user-service/user"
	_ "github.com/lib/pq"
)

var (
	ErrConnectDB    = errors.New("failed to connect to database")
	ErrExecuteQuery = errors.New("error executing query")
)

type UserRepository struct {
	DbConnection *sql.DB
}

func NewUserRepository(host, port, user, password, database string) (*UserRepository, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)

	dbConnection, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnectDB, err)
	}

	err = dbConnection.Ping()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnectDB, err)
	}

	return &UserRepository{
		DbConnection: dbConnection,
	}, nil
}

func (r *UserRepository) Find(id string) (*user.User, error) {
	query := `
	SELECT id, name, notification_config FROM weather.Users
	WHERE id = $1;
	`

	var userID, name, rawNotificationConfig string

	err := r.DbConnection.QueryRow(query, id).Scan(&userID, &name, &rawNotificationConfig)

	if err == sql.ErrNoRows {
		return nil, user.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrExecuteQuery, err)
	}

	var notificationConfig user.NotificationConfig
	err = json.Unmarshal([]byte(rawNotificationConfig), &notificationConfig)

	if err != nil {
		return nil, fmt.Errorf("error reading notification config: %w", ErrExecuteQuery, err)
	}

	return &user.User{
		ID:                 userID,
		Name:               name,
		NotificationConfig: notificationConfig,
	}, nil
}

func (r *UserRepository) Save(u *user.User) error {
	query := `
	INSERT INTO weather.Users (id, name, notification_config)
	VALUES ($1, $2, $3)
	ON CONFLICT(id)
	DO UPDATE SET
		name = $2,
		notification_config = $3;
	`

	notificationConfig, err := json.Marshal(u.NotificationConfig)

	if err != nil {
		return err
	}

	_, err = r.DbConnection.Query(query, u.ID, u.Name, notificationConfig)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrExecuteQuery, err)
	}

	return nil
}

func (r *UserRepository) Close() {
	r.DbConnection.Close()
}
