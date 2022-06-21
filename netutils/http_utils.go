package netutils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

//DownloadHTTPFile DownloadHttpFile ßhttp contents
func DownloadHttpFile(url, fileName string) error {
	client := http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; WOW64; rv:41.0) Gecko/20100101 Firefox/41.0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	fmt.Println(resp.Request.URL.String(), "downloaded!")

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
