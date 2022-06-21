package stock

//StockPrice StockPrice
type StockPrice struct {
	DateStr       string  `json:"dateStr" bson:"dateStr"`
	Code          string  `json:"code" bson:"code"` //stock code
	Name          string  `json:"name" bson:"name"`
	Qty           int     `json:"qty" bson:"qty"`     //成交股數
	Qty2          int     `json:"qty2" bson:"qty2"`   //成交筆數
	Total         int64   `json:"total" bson:"total"` //成交金額
	StartPrice    float64 `json:"startPrice" bson:"startPrice"`
	HighPrice     float64 `json:"highPrice" bson:"highPrice"`
	LowPrice      float64 `json:"lowPrice" bson:"lowPrice"`
	EndPrice      float64 `json:"endPrice" bson:"endPrice"`
	UpDown        string  `json:"upDown" bson:"upDown"` //漲跌(''平, '+'漲, '-'跌, 'X'除權息)
	Step          float64 `json:"step" bson:"step"`     //漲跌價差
	LastBuyPrice  float64 `json:"lastBuyPrice" bson:"lastBuyPrice"`
	LastBuyQty    int     `json:"lastBuyQty" bson:"lastBuyQty"`
	LastSellPrice float64 `json:"lastSellPrice" bson:"lastSellPrice"`
	LastSellQty   int     `json:"lastSellQty" bson:"lastSellQty"`
	PeRatio       float64 `json:"peRatio" bson:"peRatio"` //本益比
}
