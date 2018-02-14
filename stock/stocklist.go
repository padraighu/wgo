package stock

func GetStockList() []string {
		lst := make([]string, 10)
		lst = []string{"SPY", "DIA", "IWM", "AAPL", "XOP", "VDE", "XLE", "IVV", "VTI", "UAL"}
        return lst
}