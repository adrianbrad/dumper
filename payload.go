package dumper

type Payload struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	DriverID  int64   `json:"driver_id"`
}

func (p Payload) Valid() bool {
	if p.Latitude == 0 {
		return false
	}
	if p.Longitude == 0 {
		return false
	}
	if p.DriverID == 0 {
		return false
	}

	return true
}
