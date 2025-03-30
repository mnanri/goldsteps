package repository

import (
	"encoding/csv"
	"os"
	"server/db"
	"server/models"
	"strconv"
	"strings"

	"gorm.io/gorm/clause"
)

func ImportStocks(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, row := range records[1:] { // Skip header
		stockCode, _ := strconv.Atoi(row[0])
		settlementMonth, _ := strconv.Atoi(strings.ReplaceAll(row[8], "月", ""))
		capital, _ := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(row[9], "千円", "000"), ",", ""))
		listingDate := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(row[13], "年", "/"), "月", "/"), "日", "")
		unitShares, _ := strconv.Atoi(strings.ReplaceAll(row[14], "株", ""))
		stockExist := true
		if row[3] == "---" {
			stockExist = false
		}

		stock := models.Stock{
			StockCode:          stockCode,
			MarketType:         row[1],
			StockName:          row[2],
			StockExist:         stockExist,
			CompanyName:        row[4],
			EnglishCompanyName: row[5],
			Industry:           row[6],
			Representative:     row[7],
			SettlementMonth:    settlementMonth,
			Capital:            capital,
			Address:            row[10],
			Phone:              row[11],
			ListingMarket:      row[12],
			ListingDate:        listingDate,
			UnitShares:         unitShares,
		}

		db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&stock)
	}

	return nil
}

func ImportStockDetails(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, row := range records[1:] { // Skip header
		stockCode, _ := strconv.Atoi(row[0])
		employeesSolo, _ := strconv.Atoi(row[3])
		employeesConsolidated, _ := strconv.Atoi(strings.Replace(row[4], "---", "0", -1))
		averageAge, _ := strconv.ParseFloat(row[5], 64)
		averageSalary, _ := strconv.Atoi(row[6])

		stockDetail := models.StockDetail{
			StockCode:             stockCode,
			Feature:               row[1],
			Business:              row[2],
			EmployeesSolo:         employeesSolo,
			EmployeesConsolidated: employeesConsolidated,
			AverageAge:            averageAge,
			AverageSalary:         averageSalary,
		}

		db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&stockDetail)
	}

	return nil
}
