package stock

// Get Price Service Instance
func GetPriceService() *PriceService {
	return &PriceService{
		dao: &PriceDAO{
			db: client.Database(db_name),
		},
	}
}
