package util

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/ICOMP-UNC/newworld-agustinhernando2/internal/models"
)

func LoadData(path string) []models.Sale {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}

	var data []models.Sale
	for _, record := range records {
		var s models.Sale
		s.Id, err = strconv.Atoi(record[0])
		if err != nil {
			continue
		}
		s.CustomerEmail = record[1]
		s.DeliveryStatus = record[2]
		s.Total, err = strconv.ParseFloat(record[3], 64)
		if err != nil {
			continue
		}

		data = append(data, s)
	}

	return data
}
