package pkg

import (
	"os"
)

func openFile(Filepath string) *os.File {
	file, err := os.Open(Filepath)
	CheckErr(err, "打开文件失败!")
	return file
}

func createFile(filePath string) *os.File {
	file, err := os.Create(filePath)
	CheckErr(err, "创建文件失败!")
	return file
}

func deleteFile(filePath string) {
	err := os.Remove(filePath)
	CheckErr(err, "删除指定文件失败!")
}
