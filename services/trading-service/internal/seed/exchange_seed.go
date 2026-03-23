package seed

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gorm.io/gorm"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

func RunExchangeSeed(db *gorm.DB) error {
	_, filename, _, _ := runtime.Caller(0)
	csvPath := filepath.Join(filepath.Dir(filename), "exchanges.csv")

	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// skip header row
	for _, record := range records[1:] {
		if len(record) < 8 {
			continue
		}

		exchange := model.Exchange{
			Name:           strings.TrimSpace(record[0]),
			Acronym:        strings.TrimSpace(record[1]),
			MicCode:        strings.TrimSpace(record[2]),
			Polity:         strings.TrimSpace(record[3]),
			Currency:       strings.TrimSpace(record[4]),
			TimeZone:       strings.TrimSpace(record[5]),
			OpenTime:       strings.TrimSpace(record[6]),
			CloseTime:      strings.TrimSpace(record[7]),
			TradingEnabled: true,
		}

		var existing model.Exchange
		err := db.Where("mic_code = ?", exchange.MicCode).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&exchange).Error; err != nil {
				log.Printf("failed to create exchange %s: %v", exchange.MicCode, err)
				return err
			}
			log.Printf("created exchange: %s (%s)", exchange.Name, exchange.MicCode)
		} else if err != nil {
			log.Printf("failed to query exchange %s: %v", exchange.MicCode, err)
			return err
		} else {
			log.Printf("exchange already exists: %s (%s)", exchange.Name, exchange.MicCode)
		}
	}

	return nil
}
