// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
	"os"
	"syscall"
	"unsafe"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

// https://www.cnblogs.com/zpchcbd/p/12296624.html
var (
	moduser32         = syscall.NewLazyDLL("user32.dll")
	windowFromPoint   = moduser32.NewProc("WindowFromPoint")
	getClassNameAddr  = moduser32.NewProc("GetClassNameW")
	getWindowTextAddr = moduser32.NewProc("GetWindowTextW")
)

func main() {
	go func() {
		w := new(app.Window)
		mylog.Check(loop(w))
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	edit1 := new(widget.Editor)
	edit2 := new(widget.Editor)
	edit3 := new(widget.Editor)

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// 布局编辑器
			layoutEditors(gtx, th, edit1, edit2, edit3)

			evt, ok := gtx.Source.Event(
				pointer.Filter{
					Target: edit1,
					Kinds:  pointer.Press | pointer.Release,
				},
				pointer.Filter{
					Target: edit2,
					Kinds:  pointer.Press | pointer.Release,
				},
				pointer.Filter{
					Target: edit3,
					Kinds:  pointer.Press | pointer.Release,
				},
			)
			if ok {
				e, ok := evt.(pointer.Event)
				if ok {
					switch e.Buttons {
					case pointer.ButtonPrimary:
						// 监听鼠标点击事件
						if e.Kind == pointer.Press {
							pt := e.Position
							// 获取窗口句柄
							hwnd := findWindowFromPoint(pt)
							edit1.SetText(fmt.Sprintf("0x%.8x", hwnd))

							var className [256]uint16
							length := getClassName(hwnd, &className[0], 256)
							edit2.SetText(utf16ToString(className[:length]))

							title := make([]uint16, 256)
							getWindowText(hwnd, &title[0], 256)
							edit3.SetText(utf16ToString(title))
						}
					}
				}
			}
			e.Frame(gtx.Ops)
		}
	}
}

// 布局编辑器
func layoutEditors(gtx layout.Context, th *material.Theme, edit1, edit2, edit3 *widget.Editor) {
	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Editor(th, edit1, "窗口句柄").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Editor(th, edit2, "类名").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Editor(th, edit3, "窗口标题").Layout(gtx)
		}),
	)
}

// findWindowFromPoint 获取鼠标位置下的窗口句柄
func findWindowFromPoint(pt f32.Point) syscall.Handle {
	hwnd, _, _ := syscall.Syscall(windowFromPoint.Addr(), 2, uintptr(pt.X), uintptr(pt.Y), 0)
	return syscall.Handle(hwnd)
}

// getClassName 获取窗口类名
func getClassName(hwnd syscall.Handle, className *uint16, length int) int {
	r, _, _ := syscall.Syscall(getClassNameAddr.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(className)), uintptr(length))
	return int(r)
}

// getWindowText 获取窗口标题
func getWindowText(hwnd syscall.Handle, title *uint16, length int) {
	syscall.Syscall(getWindowTextAddr.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(title)), uintptr(length))
}

// utf16ToString 将[]uint16转为string
func utf16ToString(u []uint16) string {
	return syscall.UTF16ToString(u)
}
