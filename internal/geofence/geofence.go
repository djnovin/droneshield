package algorithms

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// IsPointInPolygon checks if a given point is inside a polygon
//
// Parameters:
//   - point: the point to check
//   - polygon: the polygon to check against
//
// Returns:
//   - bool: true if the point is inside the polygon, false otherwise
func IsPointInPolygon(point Point, polygon []Point) bool {
	n := len(polygon)
	if n < 3 {
		return false
	}

	inside := false
	j := n - 1
	for i := 0; i < n; i++ {
		if (polygon[i].Latitude > point.Latitude) != (polygon[j].Latitude > point.Latitude) &&
			point.Longitude < (polygon[j].Longitude-polygon[i].Longitude)*(point.Latitude-polygon[i].Latitude)/(polygon[j].Latitude-polygon[i].Latitude)+polygon[i].Longitude {
			inside = !inside
		}
		j = i
	}

	return inside
}
