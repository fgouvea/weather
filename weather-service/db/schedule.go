package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fgouvea/weather/weather-service/schedule"
	_ "github.com/lib/pq"
)

var (
	ErrConnectDB    = errors.New("failed to connect to database")
	ErrExecuteQuery = errors.New("error executing query")
)

type ScheduleRepository struct {
	DbConnection *sql.DB
}

func NewScheduleRepository(host, port, user, password, database string) (*ScheduleRepository, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)

	dbConnection, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnectDB, err)
	}

	err = dbConnection.Ping()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnectDB, err)
	}

	return &ScheduleRepository{
		DbConnection: dbConnection,
	}, nil
}

func (r *ScheduleRepository) Find(id string) (schedule.Schedule, error) {
	query := `
	SELECT id, user_id, city_name, status, time FROM weather.Schedules
	WHERE id = $1;
	`

	var scheduleID, userID, cityName, status string
	var scheduleTime time.Time

	err := r.DbConnection.QueryRow(query, id).Scan(&scheduleID, &userID, &cityName, &status, &scheduleTime)

	if err == sql.ErrNoRows {
		return schedule.Schedule{}, errors.New("could not find schedule")
	}

	if err != nil {
		return schedule.Schedule{}, fmt.Errorf("%w: %w", ErrExecuteQuery, err)
	}

	return schedule.Schedule{
		ID:       scheduleID,
		UserID:   userID,
		CityName: cityName,
		Status:   status,
		Time:     scheduleTime,
	}, nil
}

func (r *ScheduleRepository) Save(s schedule.Schedule) error {
	query := `
	INSERT INTO weather.Schedules (id, user_id, city_name, status, time)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT(id)
	DO UPDATE SET
		user_id = $2,
		city_name = $3,
		status = $4,
		time = $5;
	`

	_, err := r.DbConnection.Query(query, s.ID, s.UserID, s.CityName, s.Status, s.Time)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrExecuteQuery, err)
	}

	return nil
}

func (r *ScheduleRepository) FindAllBefore(t time.Time) ([]schedule.Schedule, error) {
	query := `
	SELECT id, user_id, city_name, status, time FROM weather.Schedules
	WHERE status = 'active' AND time < $1;
	`

	rows, err := r.DbConnection.Query(query, t)

	if err == sql.ErrNoRows {
		return []schedule.Schedule{}, nil
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []schedule.Schedule

	for rows.Next() {
		var scheduleID, userID, cityName, status string
		var scheduleTime time.Time

		rows.Scan(&scheduleID, &userID, &cityName, &status, &scheduleTime)

		schdl := schedule.Schedule{
			ID:       scheduleID,
			UserID:   userID,
			CityName: cityName,
			Status:   status,
			Time:     scheduleTime,
		}

		result = append(result, schdl)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrExecuteQuery, err)
	}

	return result, nil
}

func (r *ScheduleRepository) Close() {
	r.DbConnection.Close()
}
