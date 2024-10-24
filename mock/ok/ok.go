package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	"time"
)

func main() {
	// 指定窗口标题
	targetWindowTitle := "Async Updating"

	// 查找窗口句柄
	hwnd := robotgo.FindWindow(targetWindowTitle)
	if hwnd == 0 {
		fmt.Println("未找到窗口:", targetWindowTitle)
		return
	}

	// 重复执行最小化和还原操作
	for i := 0; i < 10000; i++ {
		// 最小化窗口
		win.ShowWindow(win.HWND(hwnd), win.SW_MINIMIZE)
		time.Sleep(100 * time.Millisecond)

		// 还原窗口
		win.ShowWindow(win.HWND(hwnd), win.SW_RESTORE)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("操作完成。")
}
