package models

// type Stock struct {
// 	ID                 uint `gorm:"primaryKey"`
// 	StockCode          int  `gorm:"uniqueIndex"`
// 	MarketType         string
// 	StockName          string
// 	StockExist         bool
// 	CompanyName        string
// 	EnglishCompanyName string
// 	Industry           string
// 	Representative     string
// 	SettlementMonth    int
// 	Capital            int
// 	Address            string
// 	Phone              string
// 	ListingMarket      string
// 	ListingDate        string
// 	UnitShares         int
// }

// type StockDetail struct {
// 	ID                    uint `gorm:"primaryKey"`
// 	StockCode             int  `gorm:"uniqueIndex"`
// 	Feature               string
// 	Business              string
// 	EmployeesSolo         *int
// 	EmployeesConsolidated *int
// 	AverageAge            *float64
// 	AverageSalary         *int

// 	Stock Stock `gorm:"foreignKey:StockCode;references:StockCode;constraint:OnDelete:CASCADE"`
// }

type Stock struct {
	ID                 uint   `gorm:"primaryKey" json:"id"`
	StockCode          int    `gorm:"uniqueIndex" json:"stock_code"`
	MarketType         string `json:"market_type"`
	StockName          string `json:"stock_name"`
	StockExist         bool   `json:"stock_exist"`
	CompanyName        string `json:"company_name"`
	EnglishCompanyName string `json:"english_company_name"`
	Industry           string `json:"industry"`
	Representative     string `json:"representative"`
	SettlementMonth    int    `json:"settlement_month"`
	Capital            int    `json:"capital"`
	Address            string `json:"address"`
	Phone              string `json:"phone"`
	ListingMarket      string `json:"listing_market"`
	ListingDate        string `json:"listing_date"`
	UnitShares         int    `json:"unit_shares"`
}

type StockDetail struct {
	ID                    uint    `gorm:"primaryKey" json:"id"`
	StockCode             int     `gorm:"uniqueIndex" json:"stock_code"`
	Feature               string  `json:"feature"`
	Business              string  `json:"business"`
	EmployeesSolo         int     `json:"employees_solo"`
	EmployeesConsolidated int     `json:"employees_consolidated"`
	AverageAge            float64 `json:"average_age"`
	AverageSalary         int     `json:"average_salary"`

	Stock Stock `gorm:"foreignKey:StockCode;references:StockCode;constraint:OnDelete:CASCADE" json:"stock"`
}
