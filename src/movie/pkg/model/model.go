package model

import "movieexample.com/src/metadata/pkg/model"

// Moviedetails includes movie metadata it's aggregated ratig
type MovieDetails struct {
	Rating   *float64       `json:"rating,omitEmpty"`
	Metadata model.Metadata `json:"metadata"`
}
