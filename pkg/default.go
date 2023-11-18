package pkg

import (
	"fmt"
	"os"
)

func CheckErr(err error, notice string) {
	if err != nil {
		fmt.Println(notice)
		fmt.Println(err.Error())
		os.Exit(500)
	}
}

func CheckErr_Df(err error, notice, fileName string) {
	if err != nil {
		deleteFile(fileName)
		fmt.Println(notice)
		os.Exit(500)
	}
}
