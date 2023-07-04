package drawing

import (
	"github.com/zzl/go-gdiplus/gdip"
	"log"
)

func init() {
	status := gdip.GdiplusStartup()
	if status != gdip.Ok {
		log.Fatalln("gdip startup failed")
	}

}
