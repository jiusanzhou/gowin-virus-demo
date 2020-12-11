/*
 * Copyright (c) 2020 wellwell.work, LLC by Zoe
 *
 * Licensed under the Apache License 2.0 (the "License");
 * You may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"github.com/lxn/walk"
	widgets "github.com/lxn/walk/declarative"
)

// setState() is value different, re-render the widget

// App contains all state
type App struct {
	Title string
	_mw   *walk.MainWindow
}

// NewApp ...
func NewApp(title string, ws ...widgets.Widget) (*App, error) {
	app := &App{
		Title: title,
	}

	return app, app.Create(ws...)
}

// MainWindow ...
func (app *App) MainWindow() *walk.MainWindow {
	return app._mw
}

// Create ...
func (app *App) Create(ws ...widgets.Widget) error {
	wmw := widgets.MainWindow{
		AssignTo: &app._mw,
		Name:     "app",
		Title:    widgets.Bind("title()"),
		Layout:   widgets.VBox{MarginsZero: true},
		Size:     widgets.Size{Width: 540, Height: 360},
		Children: ws,
		Functions: map[string]func(args ...interface{}) (interface{}, error){
			"title": func(args ...interface{}) (interface{}, error) {
				return app.Title, nil
			},
		},
	}

	return wmw.Create()
}

// Run ...
func (app *App) Run() error {
	_ = app._mw.Run()
	return nil
}
