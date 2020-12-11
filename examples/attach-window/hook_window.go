package main

import (
	"errors"
	"fmt"

	utils "labs.zoe.im/gowin-virus-demo"
)

// WindowAttachToWeChat 窗口吸附在微信旁边
func WindowAttachToWeChat(mw *utils.Window, name string) error {

	if name == "" {
		name = "微信"
	}

	// 找到微信的主窗口
	ws, err := utils.ListWindows() // TODO: filter by process id
	if err != nil {
		return err
	}

	var w *utils.Window
	for _, x := range ws {
		// TODO: try to check this???
		fmt.Println(x.Title())
		if x.Title() == name {
			w = x
		}
	}

	if w == nil {
		fmt.Println("查找微信窗口失败")
		return errors.New("can't find the window")
	}

	fmt.Println("[设置事件] 找到了微信的主窗口")
	// 设置位置跟随

	// 位置和大小变动
	moving := false
	_, err = w.SetWinEventHook(func(evt *utils.Event) error {
		if evt.ObjectID != 0 {
			return nil
		}
		if evt.HWND != w.Handle() {
			// 其他窗口激活则隐藏
			// fmt.Println("其他窗口")
			return nil
		}
		switch evt.Type {
		case utils.EVENT_SYSTEM_MOVESIZEEND:
			moving = false
		case utils.EVENT_SYSTEM_MOVESIZESTART:
			moving = true
		case utils.EVENT_OBJECT_LOCATIONCHANGE:
			if moving {
				rect := w.GetRect()
				mw.SetPosition(int(rect.Right-7), int(rect.Top))
			}
		case utils.EVENT_OBJECT_HIDE, utils.EVENT_SYSTEM_MINIMIZESTART:
			// ,
			// 隐藏
			fmt.Println("微信窗口隐藏了,隐藏本窗口")
		case utils.EVENT_OBJECT_SHOW, utils.EVENT_SYSTEM_FOREGROUND, utils.EVENT_SYSTEM_MINIMIZEEND:
			// ,
			fmt.Println("微信窗口显示了,显示本窗口")
			// 显示
			// 设置和激活当前窗口
			mw.Show()
			mw.SetZ(utils.HWND_TOP)
			rect := w.GetRect()
			mw.SetPosition(int(rect.Right-7), int(rect.Top))
			// TODO: 同时立即将微信设置上来

		default:
			// fmt.Printf("未知事件 HWND: %v Type: %x Event: %v\n", w.hwnd, evt.Type, evt)
		}
		return nil
	})
	if err != nil {
		fmt.Println("设置失败")
		return err
	}

	// 将微信设置为显示并为前景
	w.Show()
	w.SetForeground()

	// 将本窗口设置到微信窗口旁
	rect := w.GetRect()
	// 因为窗口可以调整大小所以去掉8px
	mw.SetPosition(int(rect.Right)-7, int(rect.Top))
	mw.SetZ(utils.HWND_TOP)

	return nil
}