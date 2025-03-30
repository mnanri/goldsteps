package routes

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"server/repository"

	"github.com/labstack/echo/v4"
)

var (
	stockfile       string
	stockDetailfile string
)

func init() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return
	}

	// Construct the file paths
	stockfile = filepath.Join(wd, "stock_master_data", "stock_fundamental_202502.csv")
	stockDetailfile = filepath.Join(wd, "stock_master_data", "stock_profile_202502.csv")
}

func importStockMasterDataFromCSV(c echo.Context) error {

	err := repository.ImportStocks(stockfile)
	if err != nil {
		fmt.Println("Error importing stocks:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save data in Stock"})
	}

	err = repository.ImportStockDetails(stockDetailfile)
	if err != nil {
		fmt.Println("Error importing stock details:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save data in StockDetail"})
	}

	return c.JSON(http.StatusCreated, nil)
}

func RegisterImportStockMasterDataFromCSV(e *echo.Group) {
	e.GET("/stock_master", importStockMasterDataFromCSV)
}
