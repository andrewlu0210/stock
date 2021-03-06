# Stock Project
- [Installation](#Installation)
- [MongoDB Setup](#mongo)
- [快速下載股票資訊](#download)


### Installation <span id="Installation"></soan>
- Install

```Shell
$ go get -u github.com/andrewlu0210/stock
```

- Import

```Go
import "github.com/andrewlu0210/stock"
```

### MongoDB Setup <span id="mongo"></span>
- Install MongoDB and get root account and passwordf
```Shell
# assume the following codes in setup.go file
$ go run setup.go
```

```Go
package main

import (
	"github.com/andrewlu0210/stock"
)

var (
	dbHost     = "127.0.0.1"
	dbName     = "stockDB"
	dbAccount  = "account"
	dbPassword = "password"
)

func main() {
	root := "root"
	rootPasswd := "password"

	stock.SetMongo(dbHost, dbName, dbAccount, dbPassword)
	stock.ResetDb(root, rootPasswd)

}
```


### 快速下載股票資訊 <span id="download"></span>
```Shell
# assume the following codes in sample.go file
$ go run sample.go
```

```Go
package main

import (
	"github.com/andrewlu0210/stock"
)

func main() {
	downloadDate := "20220615"
	dbHost := "127.0.0.1"
	dbName := "stockDB"
	dbAccount, dbPasswd := "account", "password"
	csvRoot := "/stock_csv/stock_price"
	stock.SetMongo(dbHost, dbName, dbAccount, dbPasswd)
	stock.ConnectDb()
	defer stock.DisconnectDb()

	downloader := stock.GetPriceService().GetDownloader(csvRoot)
	downloader.DownloadStockPrice(downloadDate)
    // file (20220615.csv) will save to /stock_csv/stock_price/2022/202206

}
```
