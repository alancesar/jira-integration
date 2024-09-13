package financial

import (
	"errors"
	"strconv"
)

type (
	Volume struct {
		PartnerID   string
		CreatedAt   string
		Offer       string
		Product     string
		PartnerName string
		Volume      float64
		Operations  int
	}
)

func NewVolume(columns []string) (Volume, error) {
	if len(columns) != 7 {
		return Volume{}, errors.New("unexpected number of columns")
	}

	volume, err := strconv.ParseFloat(columns[2], 64)
	if err != nil {
		return Volume{}, err
	}

	operations, err := strconv.Atoi(columns[3])
	if err != nil {
		return Volume{}, err
	}

	return Volume{
		PartnerID:   columns[1],
		CreatedAt:   columns[4],
		Offer:       columns[5],
		Product:     columns[6],
		PartnerName: columns[0],
		Volume:      volume,
		Operations:  operations,
	}, nil
}
