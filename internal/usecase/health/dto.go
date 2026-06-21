package health

import "time"

type StatusResponse struct {
	Name        string
	Environment string
	Status      string
	Timestamp   time.Time
}
