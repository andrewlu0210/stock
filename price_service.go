package stock

type PriceService struct {
	dao *PriceDAO
}

// Get A Downloader Instance to download stock price CSV File
func (service *PriceService) GetDownloader(csvRoot string) *PriceDownloader {
	return &PriceDownloader{
		csvRoot: csvRoot,
		dao:     service.dao,
	}
}

// Get latest stock price record date in mongodb
func (service *PriceService) GetLatestDate() string {
	return service.dao.getLatestDate()
}

func (service *PriceService) GetPriceByCodeDate(dateStr, code string) *StockPrice {
	return service.dao.getDailyPrice(dateStr, code)
}

func (service *PriceService) GetPricesByCode(code, fromDate, toDate string, oldToNew bool) []*StockPrice {
	return service.dao.getPricesByCode(code, fromDate, toDate, oldToNew)
}

func (service *PriceService) GetPricesByDate(dateStr string) []*StockPrice {
	return service.dao.getPricesByDate(dateStr)
}
