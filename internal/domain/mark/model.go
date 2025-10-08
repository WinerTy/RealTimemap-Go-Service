package mark

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"realtimemap-service/internal/domain/category"
	"realtimemap-service/internal/pkg/entity"
	"time"
)

type Mark struct {
	ID             int      // ID айди
	OwnerID        int      // OwnerID айди юзера
	Name           string   // Name название метки
	AdditionalInfo *string  // AdditionalInfo доп описание
	Geom           Geometry // Geom Геометрия, прихоидт из бд в виде JSON
	IsEnded        bool
	DurationHours  int
	EndAt          time.Time
	Photo          []entity.Image
	Category       category.Category
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func (g *Geometry) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	var source []byte

	switch v := src.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	default:
		return fmt.Errorf("unsupported type for Geometry: %T", src)
	}

	if len(source) == 0 {
		return nil
	}

	return json.Unmarshal(source, g)
}

// Value реализует интерфейс driver.Valuer (для записи в БД).
func (g Geometry) Value() (driver.Value, error) {
	if len(g.Coordinates) == 0 {
		return nil, nil
	}
	return json.Marshal(g)
}
