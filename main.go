package main

import (
	"T-S0ph0n/cmd"
	"T-S0ph0n/pkg"
)

func main() {
	err := cmd.Tsophon.Execute()
	pkg.CheckErr(err, "开启解析命令行解析时出错!")
}
