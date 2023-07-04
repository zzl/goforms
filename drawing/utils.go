package drawing

import (
	"github.com/zzl/go-gdiplus/gdip"
	"log"
)

func toRectF(rect Rect) RectF {
	return RectF{
		X:      float32(rect.X),
		Y:      float32(rect.Y),
		Width:  float32(rect.Width),
		Height: float32(rect.Height),
	}
}

func ensureOk(status gdip.Status) {
	if status != gdip.Ok {
		log.Fatalln("gdip status not ok:", status)
	}
}

func checkStatus(status gdip.Status) {
	if status != gdip.Ok {
		println("gdip status not ok:", status)
	}
}
