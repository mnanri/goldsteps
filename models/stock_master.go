package models

type Stock struct {
	ID                 uint `gorm:"primaryKey"`
	StockCode          int  `gorm:"uniqueIndex"`
	MarketType         string
	StockName          string
	StockExist         bool
	CompanyName        string
	EnglishCompanyName string
	Industry           string
	Representative     string
	SettlementMonth    int
	Capital            int
	Address            string
	Phone              string
	ListingMarket      string
	ListingDate        string
	UnitShares         int
}

type StockDetail struct {
	ID                    uint `gorm:"primaryKey"`
	StockCode             int  `gorm:"uniqueIndex"`
	Feature               string
	Business              string
	EmployeesSolo         *int
	EmployeesConsolidated *int
	AverageAge            *float64
	AverageSalary         *int

	Stock Stock `gorm:"foreignKey:StockCode;references:StockCode;constraint:OnDelete:CASCADE"`
}
