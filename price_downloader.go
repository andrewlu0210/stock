package stock

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/andrewlu0210/stock/netutils"
	"github.com/andrewlu0210/stock/parser"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type PriceDownloader struct {
	csvRoot string
	dao     *PriceDAO
}

func (service *PriceDownloader) DownloadStockPrice(dateStr string) (bool, bool) {
	date, err := parser.ParseDate(dateStr)
	if err != nil {
		log.Println(err)
		return false, false
	}
	saved_to_db := false
	is_download := false
	if date.Weekday() == time.Sunday {
		log.Printf("Download date [%s] is Sunday _ (週日不用交易)\n", dateStr)
		return saved_to_db, is_download
	} else if date.Year() >= 2019 && date.Weekday() == time.Saturday {
		//2019年以後週六不交易
		log.Printf("Download date [%s] is Saturday _ (2019年以後週六不交易)\n", dateStr)
		return saved_to_db, is_download
	}
	cnt := service.dao.countByDate(dateStr)
	if cnt > 0 {
		log.Printf("Daily Price[%s - %s] already exists in MongoDB\n", dateStr, date.Weekday())
	} else {
		saved_to_db, is_download = service.downloadCSVAndSave(dateStr, date)
	}
	return saved_to_db, is_download
}

func (service *PriceDownloader) downloadCSVAndSave(dateStr string, date time.Time) (bool, bool) {
	priceRootDir := fmt.Sprintf("%s/%s/%s", service.csvRoot, date.Format("2006"), date.Format("200601"))
	checkMakeDirs(priceRootDir)
	fileName := fmt.Sprintf("%s/%s.csv", priceRootDir, dateStr)
	is_download := false
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Printf("下載[%s]股價超資料...\n", dateStr)
		url := fmt.Sprintf("https://www.twse.com.tw/exchangeReport/MI_INDEX?response=csv&date=%s&type=ALL", dateStr)
		downErr := netutils.DownloadHttpFile(url, fileName)
		checkError(downErr)
		is_download = true
	} else {
		log.Printf("CSV[%s]檔案已存在!\n", fileName)
	}
	//read csv file
	saved_to_db := false
	dailyPrices := service.readCSVFile(fileName, dateStr)
	if len(dailyPrices) > 0 {
		//TODO: add to mongodb
		service.dao.addDailyPrices(dailyPrices)
		saved_to_db = true
		fmt.Printf("[%s]已存入(%d)個股價資訊!\n", dateStr, len(dailyPrices))
	} else {
		fmt.Printf("[%s]0筆股價資訊!\n", dateStr)
		//TODO: 移除該CSV檔案，沒有資料
		//removeFile(fileName)
		is_download = false
	}
	return saved_to_db, is_download
}

func (service *PriceDownloader) readCSVFile(fileName, dateStr string) []*StockPrice {
	csvFile, err := os.Open(fileName)
	checkError(err)
	r := transform.NewReader(csvFile, traditionalchinese.Big5.NewDecoder())
	sc := bufio.NewScanner(r)

	startRead := false
	stockPrices := []*StockPrice{}
	for sc.Scan() {
		line := string(sc.Bytes())
		if startRead {
			//fmt.Println(line)
			sp := service.readCSVLine(line, dateStr)
			if sp != nil {
				stockPrices = append(stockPrices, sp)
			}
		}
		if strings.HasPrefix(line, "\"證券代號") {
			startRead = true
		}
	}
	if err = sc.Err(); err != nil {
		log.Fatal(err)
	}
	if err = csvFile.Close(); err != nil {
		log.Panic(err)
	}
	//return daily prices list
	return stockPrices
}

func (service *PriceDownloader) readCSVLine(line, dateStr string) *StockPrice {
	words := strings.Split(line, "\",\"")
	if len(words) == 16 {
		//0開頭的
		if strings.HasPrefix(words[0], "=") {
			//ETF
			chk := strings.Replace(words[0], "=\"", "", 1)
			if strings.HasPrefix(chk, "00") {
				//只看00開頭
				ret := service.newStockPriceData(dateStr, chk, words)
				return ret
			}
		} else {
			//普通個股
			code := strings.Replace(words[0], "\"", "", -1)
			ret := service.newStockPriceData(dateStr, code, words)
			return ret
		}
	}
	return nil
}

func (service *PriceDownloader) newStockPriceData(dateStr, code string, words []string) *StockPrice {
	vo := &StockPrice{
		DateStr: dateStr,
		Code:    code,
	}
	vo.Name = words[1]
	vo.Qty = parser.StringToInt(strings.Replace(words[2], ",", "", -1))
	vo.Qty2 = parser.StringToInt(strings.Replace(words[3], ",", "", -1))
	vo.Total = parser.StringToLong(strings.Replace(words[4], ",", "", -1))
	vo.StartPrice = parser.StringToFloat(strings.Replace(words[5], ",", "", -1))
	vo.HighPrice = parser.StringToFloat(strings.Replace(words[6], ",", "", -1))
	vo.LowPrice = parser.StringToFloat(strings.Replace(words[7], ",", "", -1))
	vo.EndPrice = parser.StringToFloat(strings.Replace(words[8], ",", "", -1))
	vo.UpDown = strings.Replace(words[9], "\"", "", -1)
	vo.Step = parser.StringToFloat(words[10])
	vo.LastBuyPrice = parser.StringToFloat(strings.Replace(words[11], ",", "", -1))
	vo.LastBuyQty = parser.StringToInt(strings.Replace(words[12], ",", "", -1))
	vo.LastSellPrice = parser.StringToFloat(strings.Replace(words[13], ",", "", -1))
	vo.LastSellQty = parser.StringToInt(strings.Replace(words[14], ",", "", -1))
	vo.PeRatio = parser.StringToFloat(strings.Replace(words[15], "\",", "", -1))

	return vo
}
