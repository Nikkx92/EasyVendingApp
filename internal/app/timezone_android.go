//go:build android
// +build android

package app

/*
#include <time.h>
#include <stdlib.h>
*/
import "C"
import (
	"time"
	_ "time/tzdata"
)

func InitAndroidTimezoneProperty() {
	var curtime C.time_t
	var curtm C.struct_tm

	C.time(&curtime)
	C.localtime_r(&curtime, &curtm)

	tzOffset := int(curtm.tm_gmtoff)
	tz := C.GoString(curtm.tm_zone)

	time.Local = time.FixedZone(tz, tzOffset)
}
