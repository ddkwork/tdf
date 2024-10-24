package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	// 窗口标题
	windowTitle := "Untitled - 网络控制台[D]【0.3.6.42/2019-07-28 00:03:56】【0.3.6.1503/2019-07-28 00:03:53】"

	// 查找窗口
	hwnd := robotgo.FindWindow(windowTitle)
	if hwnd == 0 {
		fmt.Println("未找到窗口:", windowTitle)
		return
	}

	// 激活窗口
	robotgo.SetActiveWindow(hwnd)

	// 等待一段时间以确保窗口已激活
	time.Sleep(2 * time.Second)

	// 假设表格有10行c
	numRows := 10

	for i := 0; i < numRows; i++ {
		// 每按一次键，获取对应行的 Hexdump
		robotgo.KeyTap("down")      // 按下“下”键以获取当前行的 hexdump
		time.Sleep(1 * time.Second) // 等待 Hexdump 显示

		// 复制当前显示的 hexdump 内容
		robotgo.KeyTap("c", "control") // Ctrl+C 复制
		time.Sleep(1 * time.Second)    // 等待复制完成

		// 打印处理信息
		fmt.Printf("已处理第 %d 行\n", i+1)
	}

	fmt.Println("所有行已处理完成。")
}
