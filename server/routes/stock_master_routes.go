package routes

import (
	"goldsteps/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

const stockfile = "stock_master_data/stock_fundamental_202502.csv"
const stockDetailfile = "stock_master_data/stock_profile_202502.csv"

func importStockMasterDataFromCSV(c echo.Context) error {

	err := repository.ImportStocks(stockfile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save data in Stock"})
	}

	err = repository.ImportStockDetails(stockDetailfile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save data in StockDetail"})
	}

	return c.JSON(http.StatusCreated, nil)
}

func RegisterImportStockMasterDataFromCSV(e *echo.Group) {
	e.GET("/stock_master", importStockMasterDataFromCSV)
}
