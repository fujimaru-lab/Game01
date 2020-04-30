package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/widget"

	"fyne.io/fyne/layout"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Ball and Board is not Bored!!")
	w.Resize(fyne.NewSize(1000, 800))

	// ボールの生成
	circle := canvas.NewCircle(color.White)
	circle.StrokeWidth = 5

	// 板の生成
	board := canvas.NewLine(color.White)
	board.StrokeWidth = 5
	cont := fyne.NewContainer(board)

	// ポイント表示枠の生成
	pointBox := widget.NewLabel("00 pts")
	pboxCon := fyne.NewContainer(pointBox)

	content := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(50, 50)), circle, cont, pboxCon)

	w.SetContent(content)
	log.Println(circle.Position1, circle.Position2)

	go func() {
		board.Position1 = fyne.NewPos(275, 750)
		board.Position2 = fyne.NewPos(375, 750)
		board.Refresh()
		pointBox.Move(fyne.NewPos(0, w.Canvas().Size().Height-50))
		pointBox.Refresh()
		// キーボードからの入力に応じて反射板を移動させる
		go func() {
			for {
				w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
					log.Printf("key=%s\n", string(ev.Name))
					if ev.Name == fyne.KeyRight {
						board.Move(fyne.NewPos(board.Position().X+15, board.Position().Y))
					}
					if ev.Name == fyne.KeyLeft {
						board.Move(fyne.NewPos(board.Position().X-15, board.Position().Y))
					}
				})
				canvas.Refresh(board)
			}
		}()

		diffX := 1
		diffY := 1
		point := 0
		// ボールを 1 count/msで移動する
		for {
			time.Sleep(time.Millisecond)
			x := circle.Position().X + diffX
			y := circle.Position().Y + diffY
			circle.Move(fyne.NewPos(x, y))

			if x < 0 || w.Canvas().Size().Width < x+50 {
				diffX *= -1
			}
			if y < 0 || w.Canvas().Size().Height < y+50 {
				diffY *= -1
			}
			centerX := (circle.Position1.X + circle.Position2.X) / 2
			if diffY > 0 && board.Position1.X <= centerX && centerX <= board.Position2.X && board.Position().Y <= circle.Position2.Y {
				diffY *= -1
				point++
				pointBox.SetText(fmt.Sprintf("%d pts", point))
				canvas.Refresh(pointBox)
				log.Printf("Ball (x1, y1)=(%d, %d), (x2, y2)=(%d, %d)\n", circle.Position1.X, circle.Position1.Y, circle.Position2.X, circle.Position2.Y)
				log.Printf("Ball Center (x, y)=(%d, %d)\n", circle.Position().X+25, circle.Position().Y+25)
				log.Printf("Board (x1, y1)=(%d, %d), (x2, y2)=(%d, %d)\n", board.Position1.X, board.Position1.Y, board.Position2.X, board.Position2.Y)
				log.Printf("Point:%0d\n", point)
			}
			canvas.Refresh(circle)
		}
	}()
	w.ShowAndRun()
}
