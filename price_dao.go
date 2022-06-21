package stock

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PriceDAO struct {
	db *mongo.Database
}

func (dao *PriceDAO) getLatestDate() string {
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"dateStr": -1}) //sort
	findOptions.SetLimit(1)
	col := dao.db.Collection(STOCK_PRICE_COL)
	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.M{"code": "2317"},
			bson.M{"code": "2330"},
		}},
		{Key: "dateStr", Value: bson.D{
			{Key: "$gte", Value: "20190101"},
		}},
	}
	var results []*StockPrice
	cur, err := col.Find(context.TODO(), filter, findOptions)
	checkError(err)
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var elem StockPrice
		err := cur.Decode(&elem)
		checkError(err)
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	str := ""
	if len(results) > 0 {
		str = results[0].DateStr
	}
	return str
}

func (dao *PriceDAO) getPricesByCode(code, fromDate, toDate string, oldToNew bool) []*StockPrice {
	findOptions := options.Find()
	findOptions.SetLimit(20000)
	dataOrder := 1
	if !oldToNew {
		dataOrder = -1 //NewToOld
	}
	findOptions.SetSort(bson.M{"dateStr": dataOrder}) //sort
	var results []*StockPrice
	filter := bson.D{
		{Key: "code", Value: code},
		{Key: "dateStr", Value: bson.D{
			{Key: "$gte", Value: fromDate},
			{Key: "$lte", Value: toDate},
		}},
	}
	//filter := bson.M{"code": code}
	col := dao.db.Collection(STOCK_PRICE_COL)
	cur, err := col.Find(context.TODO(), filter, findOptions)
	checkError(err)

	for cur.Next(context.TODO()) {
		var elem StockPrice
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Decode error!")
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		fmt.Println("Close Error")
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	return results
}

func (dao *PriceDAO) getPricesByDate(dateStr string) []*StockPrice {
	findOptions := options.Find()
	findOptions.SetLimit(20000)
	//findOptions.SetSort(bson.M{"endPrice": -1}) //sort
	var results []*StockPrice
	filter := bson.M{"dateStr": dateStr}
	col := dao.db.Collection(STOCK_PRICE_COL)
	cur, err := col.Find(context.TODO(), filter, findOptions)
	checkError(err)

	for cur.Next(context.TODO()) {
		var elem StockPrice
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Decode error!")
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		fmt.Println("Close Error")
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	return results
}

func (dao *PriceDAO) getDailyPrice(dateStr, code string) *StockPrice {
	var sp StockPrice
	filter := bson.D{
		{Key: "code", Value: code},
		{Key: "dateStr", Value: dateStr},
	}
	col := dao.db.Collection(STOCK_PRICE_COL)
	err := col.FindOne(context.TODO(), filter).Decode(&sp)
	if err != nil {
		log.Println(err, code)
		return nil
	}
	return &sp
}

func (dao *PriceDAO) addDailyPrices(dailyPrices []*StockPrice) {
	col := dao.db.Collection(STOCK_PRICE_COL)
	docs := make([]interface{}, len(dailyPrices))
	for i, d := range dailyPrices {
		docs[i] = d
	}
	_, err := col.InsertMany(context.TODO(), docs)
	checkError(err)
}

func (dao *PriceDAO) countByDate(dateStr string) int64 {
	col := dao.db.Collection(STOCK_PRICE_COL)
	filter := bson.M{"dateStr": dateStr}
	cnt, err := col.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return 0
	}
	return cnt
}
