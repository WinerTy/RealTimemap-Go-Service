package mark

import (
	"time"

	"github.com/mmcloughlin/geohash"
)

type TimeRangeType string

const (
	DefaultRadius   float64 = 5000.0
	DefaultSRID     int     = 4326
	DefaultDuration int     = 12
	DefaultPersists uint    = 5
)

type Filter struct {
	// Геоданные
	Longitude float64
	Latitude  float64
	Radius    float64
	SRID      int
	Geohash   []string

	// Временные
	ReferenceTime time.Time
	ShowEnded     bool

	// Для временного отрезка
	SearchWindowStart time.Time
	SearchWindowEnd   time.Time
}

func NewFilter(latitude, longitude, radius float64,
	srid, durationHours int, showEnded bool, refTime time.Time) Filter {

	if radius == 0 || radius < 0 {
		radius = DefaultRadius
	}
	if srid == 0 || srid < 1 {
		srid = DefaultSRID
	}
	if durationHours == 0 || durationHours < 0 {
		durationHours = DefaultDuration
	}

	geoHash := geohash.EncodeWithPrecision(latitude, longitude, DefaultPersists)
	neighbors := geohash.Neighbors(geoHash)
	neighbors = append(neighbors, geoHash)

	filter := Filter{
		Longitude:     longitude,
		Latitude:      latitude,
		Radius:        radius,
		SRID:          srid,
		Geohash:       neighbors,
		ShowEnded:     showEnded,
		ReferenceTime: refTime,
	}

	duration := time.Duration(durationHours) * time.Hour
	filter.SearchWindowStart = refTime.Add(-duration)
	filter.SearchWindowEnd = refTime.Add(duration)

	return filter
}
