package model

type (
	// City data
	City struct {
		CityID int    `json:"id" db:"CityID"`
		City   string `json:"City" db:"City"`
		Hash   string `json:"Hash" db:"Hash"` // xxhash64 (normalized) of 'City'
	}
)

// SelectCities database query
func (m *Model) SelectCities() []City {
	x := []City{}
	m.DB.Select(&x, "SELECT CityID, City, Hash FROM Cities ORDER BY City")
	return x
}
