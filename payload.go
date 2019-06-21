package dumper

type Payload struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	DriverID  int64   `json:"driver_id"`
}
