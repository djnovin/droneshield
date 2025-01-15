package handlers

import (
	"droneshield/internal/models"
	"droneshield/internal/services"
	"droneshield/pkg"
	"encoding/json"
	"net/http"
	"strconv"
)

func GeoFenceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetGeofences(w, r)
	case http.MethodPost:
		handleCreateGeofence(w, r)
	case http.MethodPut:
		handleUpdateGeofence(w, r)
	case http.MethodDelete:
		handleDeleteGeofence(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetGeofences(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		id, err := parseID(idStr)
		if err != nil {
			pkg.ErrBadRequest.Send(w)
			return
		}
		geofence, err := services.GetGeofence(id)
		if err != nil {
			pkg.ErrNotFound.Send(w)
			return
		}
		respondWithJSON(w, http.StatusOK, geofence)
		return
	}

	geofences, err := services.GetGeofences()
	if err != nil {
		pkg.ErrInternal.Send(w)
		return
	}
	respondWithJSON(w, http.StatusOK, geofences)
}

// handleCreateGeofence handles POST requests for creating a geofence
func handleCreateGeofence(w http.ResponseWriter, r *http.Request) {
	var geofence models.Geofence
	if err := json.NewDecoder(r.Body).Decode(&geofence); err != nil {
		pkg.ErrBadRequest.Send(w)
		return
	}

	newGeofence, err := services.AddGeofence(geofence.Name, geofence.Polygon)
	if err != nil {
		pkg.ErrBadRequest.Send(w)
		return
	}
	respondWithJSON(w, http.StatusCreated, newGeofence)
}

// handleUpdateGeofence handles PUT requests for updating a geofence
func handleUpdateGeofence(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := parseID(idStr)
	if err != nil {
		pkg.ErrBadRequest.Send(w)
		return
	}

	var geofence models.Geofence
	if err := json.NewDecoder(r.Body).Decode(&geofence); err != nil {
		pkg.ErrBadRequest.Send(w)
		return
	}

	updatedGeofence, err := services.UpdateGeofence(id, geofence.Name, geofence.Polygon)
	if err != nil {
		pkg.ErrInternal.Send(w)
		return
	}
	respondWithJSON(w, http.StatusOK, updatedGeofence)
}

func handleDeleteGeofence(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := parseID(idStr)
	if err != nil {
		pkg.ErrBadRequest.Send(w)
		return
	}

	if err := services.DeleteGeofence(int(id)); err != nil {
		pkg.ErrInternal.Send(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(idStr string) (int64, error) {
	return strconv.ParseInt(idStr, 10, 64)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
