package stock

import (
	"fmt"
	"log"
	"os"
)

func RemoveFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		log.Println(err)
		log.Printf("Remove [%s] Fail\n", fileName)
	} else {
		log.Printf("Remove [%s] Done\n", fileName)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//CheckMakeDirs - make dir if not exists
func checkMakeDirs(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		pathErr := os.MkdirAll(dir, 0777)
		if pathErr != nil {
			log.Fatal(pathErr)
		}
		fmt.Println(dir, "created")
	}
}
