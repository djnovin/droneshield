package services

import (
	"database/sql"
	"droneshield/internal/models"
	"fmt"
)

var db *sql.DB

func InitDB(dataSourceName string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, db.Ping()
}

func AddGeofence(name, polygon string) (models.Geofence, error) {
	var geofence models.Geofence

	err := db.QueryRow("INSERT INTO geofences (name, polygon) VALUES ($1, $2) RETURNING id, created_at", name, polygon).
		Scan(&geofence.ID, &geofence.CreatedAt)
	if err != nil {
		return geofence, err
	}

	geofence.Name = name
	geofence.Polygon = polygon
	return geofence, nil
}

func GetGeofence(id int64) (models.Geofence, error) {
	var geofence models.Geofence

	err := db.QueryRow("SELECT id, name, polygon, created_at FROM geofences WHERE id = $1", id).
		Scan(&geofence.ID, &geofence.Name, &geofence.Polygon, &geofence.CreatedAt)
	if err != nil {
		return geofence, err
	}

	return geofence, nil
}

func GetGeofences() ([]models.Geofence, error) {
	rows, err := db.Query("SELECT id, name, polygon, created_at FROM geofences")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var geofences []models.Geofence
	for rows.Next() {
		var geofence models.Geofence
		if err := rows.Scan(&geofence.ID, &geofence.Name, &geofence.Polygon, &geofence.CreatedAt); err != nil {
			return nil, err
		}
		geofences = append(geofences, geofence)
	}

	return geofences, nil
}

func UpdateGeofence(id int64, name, polygon string) (models.Geofence, error) {
	var geofence models.Geofence

	_, err := db.Exec("UPDATE geofences SET name = $1, polygon = $2, updated_at = NOW() WHERE id = $3", name, polygon, id)
	if err != nil {
		return geofence, err
	}

	geofence.ID = id
	geofence.Name = name
	geofence.Polygon = polygon
	return geofence, nil
}

func DeleteGeofence(id int) error {
	_, err := db.Exec("DELETE FROM geofences WHERE id = $1", id)
	return err
}
