//go:build windows

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ijry/lyshop/config"
	appcore "github.com/ijry/lyshop/core/app"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type launcher struct {
	configPath string

	mw          *walk.MainWindow
	statusLabel *walk.Label
	logBox      *walk.TextEdit
	startButton *walk.PushButton
	stopButton  *walk.PushButton

	mu          sync.Mutex
	initialized bool
	running     bool
	cancelRun   context.CancelFunc

	webURL   string
	adminURL string
	h5URL    string
}

func main() {
	var (
		configPath string
		autoStart  bool
	)
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")
	flag.BoolVar(&autoStart, "auto-start", true, "start server automatically")
	flag.Parse()

	l := &launcher{
		configPath: resolveConfigPath(configPath),
		webURL:     "http://127.0.0.1:8080/",
		adminURL:   "http://127.0.0.1:8080/admin/",
		h5URL:      "http://127.0.0.1:8080/h5/",
	}
	if err := l.run(autoStart); err != nil {
		walk.MsgBox(nil, "LYShop Launcher", "启动 GUI 失败:\n"+err.Error(), walk.MsgBoxIconError)
	}
}

func (l *launcher) run(autoStart bool) error {
	if err := (declarative.MainWindow{
		AssignTo: &l.mw,
		Title:    "LYShop 本机启动器",
		MinSize:  declarative.Size{Width: 760, Height: 500},
		Layout:   declarative.VBox{Margins: declarative.Margins{Left: 14, Top: 12, Right: 14, Bottom: 14}},
		Children: []declarative.Widget{
			declarative.Label{
				AssignTo: &l.statusLabel,
				Text:     "服务状态：未启动",
			},
			declarative.Composite{
				Layout: declarative.HBox{MarginsZero: true},
				Children: []declarative.Widget{
					declarative.PushButton{
						AssignTo:  &l.startButton,
						Text:      "启动服务",
						OnClicked: l.startService,
					},
					declarative.PushButton{
						AssignTo:  &l.stopButton,
						Text:      "停止服务",
						Enabled:   false,
						OnClicked: l.stopService,
					},
					declarative.PushButton{
						Text: "打开商城",
						OnClicked: func() {
							l.openURL(l.webURL)
						},
					},
					declarative.PushButton{
						Text: "打开后台",
						OnClicked: func() {
							l.openURL(l.adminURL)
						},
					},
					declarative.PushButton{
						Text: "打开 H5",
						OnClicked: func() {
							l.openURL(l.h5URL)
						},
					},
				},
			},
			declarative.GroupBox{
				Title:  "地址（可点击）",
				Layout: declarative.Grid{Columns: 2},
				Children: []declarative.Widget{
					declarative.Label{Text: "商城"},
					declarative.LinkLabel{
						Text: fmt.Sprintf(`<a href="%s">%s</a>`, l.webURL, l.webURL),
						OnLinkActivated: func(link *walk.LinkLabelLink) {
							l.openURL(link.URL())
						},
					},
					declarative.Label{Text: "后台"},
					declarative.LinkLabel{
						Text: fmt.Sprintf(`<a href="%s">%s</a>`, l.adminURL, l.adminURL),
						OnLinkActivated: func(link *walk.LinkLabelLink) {
							l.openURL(link.URL())
						},
					},
					declarative.Label{Text: "H5"},
					declarative.LinkLabel{
						Text: fmt.Sprintf(`<a href="%s">%s</a>`, l.h5URL, l.h5URL),
						OnLinkActivated: func(link *walk.LinkLabelLink) {
							l.openURL(link.URL())
						},
					},
				},
			},
			declarative.GroupBox{
				Title:  "运行日志",
				Layout: declarative.VBox{MarginsZero: true},
				Children: []declarative.Widget{
					declarative.TextEdit{
						AssignTo: &l.logBox,
						ReadOnly: true,
						VScroll:  true,
					},
				},
			},
		},
	}).Create(); err != nil {
		return err
	}

	l.mw.Closing().Attach(func(canceled *bool, _ walk.CloseReason) {
		l.stopService()
	})

	l.appendLog("启动器已加载")
	l.appendLog("配置文件: " + l.configPath)

	if autoStart {
		l.startService()
	}

	l.mw.Run()
	return nil
}

func (l *launcher) startService() {
	l.mu.Lock()
	if l.running {
		l.mu.Unlock()
		l.appendLog("服务已在运行")
		return
	}
	l.mu.Unlock()

	if err := l.ensureInitialized(); err != nil {
		l.setStatus("服务状态：初始化失败")
		l.appendLog("初始化失败: " + err.Error())
		walk.MsgBox(l.mw, "LYShop Launcher", "初始化失败:\n"+err.Error(), walk.MsgBoxIconError)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	l.mu.Lock()
	l.cancelRun = cancel
	l.running = true
	l.mu.Unlock()

	l.setStatus("服务状态：启动中...")
	l.updateRunControls(true)
	l.appendLog("正在启动服务...")

	go func() {
		err := appcore.RunWithContext(ctx)
		l.mw.Synchronize(func() {
			l.mu.Lock()
			l.running = false
			l.cancelRun = nil
			l.mu.Unlock()
			l.updateRunControls(false)

			if err != nil {
				l.setStatus("服务状态：异常退出")
				l.appendLog("服务异常退出: " + err.Error())
				return
			}
			l.setStatus("服务状态：已停止")
			l.appendLog("服务已停止")
		})
	}()

	go func() {
		time.Sleep(250 * time.Millisecond)
		l.mw.Synchronize(func() {
			l.mu.Lock()
			running := l.running
			l.mu.Unlock()
			if running {
				l.setStatus("服务状态：运行中")
				l.appendLog("服务运行中")
			}
		})
	}()
}

func (l *launcher) stopService() {
	l.mu.Lock()
	cancel := l.cancelRun
	running := l.running
	l.mu.Unlock()
	if !running || cancel == nil {
		return
	}
	l.setStatus("服务状态：停止中...")
	l.appendLog("正在停止服务...")
	cancel()
}

func (l *launcher) ensureInitialized() error {
	l.mu.Lock()
	if l.initialized {
		l.mu.Unlock()
		return nil
	}
	l.mu.Unlock()

	if err := appcore.Init(l.configPath); err != nil {
		return err
	}

	base := fmt.Sprintf("http://127.0.0.1:%d", config.Global.Server.Port)
	l.webURL = base + "/"
	l.adminURL = base + "/admin/"
	l.h5URL = base + "/h5/"

	l.mu.Lock()
	l.initialized = true
	l.mu.Unlock()

	l.appendLog("初始化完成（SQLite + 嵌入前端资源可直接使用）")
	return nil
}

func (l *launcher) updateRunControls(running bool) {
	if l.startButton != nil {
		l.startButton.SetEnabled(!running)
	}
	if l.stopButton != nil {
		l.stopButton.SetEnabled(running)
	}
}

func (l *launcher) setStatus(text string) {
	if l.statusLabel != nil {
		l.statusLabel.SetText(text)
	}
}

func (l *launcher) appendLog(line string) {
	if l.logBox == nil {
		return
	}
	now := time.Now().Format("15:04:05")
	current := strings.TrimRight(l.logBox.Text(), "\r\n")
	next := fmt.Sprintf("[%s] %s", now, line)
	if current == "" {
		l.logBox.SetText(next)
		return
	}
	l.logBox.SetText(current + "\r\n" + next)
	l.logBox.SetTextSelection(len(l.logBox.Text()), len(l.logBox.Text()))
}

func (l *launcher) openURL(target string) {
	if err := exec.Command("rundll32", "url.dll,FileProtocolHandler", target).Start(); err != nil {
		l.appendLog("打开浏览器失败: " + err.Error())
	}
}

func resolveConfigPath(raw string) string {
	if filepath.IsAbs(raw) {
		return raw
	}
	exePath, err := os.Executable()
	if err != nil {
		return raw
	}
	return filepath.Join(filepath.Dir(exePath), raw)
}
