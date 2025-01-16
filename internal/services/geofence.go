package services

import (
	"database/sql"
	"droneshield/internal/models"
	"fmt"
)

// var db *sql.DB

type GeofenceService struct{
  DB *sql.DB
}

func NewGeofenceService(db *sql.DB) *GeofenceService {
  return &GeofenceService{DB: db}
}

func InitDB(dataSourceName string) (*sql.DB, error) {
	var err error

	db, err = sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, db.Ping()
}

func (s *GeofenceService) AddGeofence(name, polygon string) (models.Geofence, error) {
  query := `INSERT INTO geofences (name, polygon) VALUES (:name, :polygon) RETURNING id, created_at`
  params:= map[string]interface{}{
    "name": name,
    "polygon": polygon,
  }

  var geofence models.Geofence

  err := s.DB.NamedQueryRowx(query, params).StructScan(&geofence)
  if err != nil {
    return geofence, fmt.Errorf("failed to insert geofence: %w", err)
  }
  return geofence, nil

	return geofence, nil
}

func (s *GeofenceService) GetGeofence(id int64) (models.Geofence, error) {
  var geofence models.Geofence

  err := s.DB.Get(&geofence, "SELECT id, name, polygon, created_at FROM geofences WHERE id = $1", id)

  if err != nil {
    if error.Is(err, sql.ErrNoRows) {
      return geofence, models.ErrNotFound
    }
    return geofence, fmt.Errorf("failed to get geofence: %w", err)
  }

  return geofence, nil
}

func (s *GeofenceService) GetGeofences() ([]models.Geofence, error) {
  var geofences []models.Geofence

  err := s.DB.Select(&geofences, "SELECT id, name, polygon, created_at FROM geofences")
	if err != nil {
		return nil, fmt.Errorf("failed to get geofences: %w", err)
	}

	return geofences, nil
}


func (s *GeofenceService) UpdateGeofence(id int64, name, polygon string) (models.Geofence, error) {
	query := `
		UPDATE geofences
		SET name = :name, polygon = :polygon, updated_at = NOW()
		WHERE id = :id
		RETURNING id, name, polygon, created_at, updated_at
	`

	params := map[string]interface{}{
		"id":      id,
		"name":    name,
		"polygon": polygon,
	}

	var geofence models.Geofence

	// Execute the query with named parameters
	rows, err := s.DB.NamedQuery(query, params)
	if err != nil {
		return geofence, fmt.Errorf("failed to update geofence with ID %d: %w", id, err)
	}
	defer rows.Close()

	// Fetch the updated row
	if rows.Next() {
		err = rows.StructScan(&geofence)
		if err != nil {
			return geofence, fmt.Errorf("failed to scan updated geofence with ID %d: %w", id, err)
		}
	} else {
		return geofence, fmt.Errorf("no geofence found with ID %d", id)
	}

	return geofence, nil
}

func (s *GeofenceService) DeleteGeofence(id int64) error {
  query := "DELETE FROM geofences WHERE id = $1"

  result, err := s.DB.Exec(query, id)
  if err != nil {
    return fmt.Errorf("failed to delete geofence with ID %d: %w", id, err)
  }

  rowsAffected, err := result.RowsAffected()
  if err != nil {
    return fmt.Errorf("failed to get rows affected after deleting geofence with ID %d: %w", id, err)
  }

  if rowsAffected == 0 {
    return fmt.Errorf("no geofence found with ID %d", id)
  }

  return nil
}
