package genh

import (
	"fmt"
	"log"
)

type testingTB interface {
	Error(args ...any)
	FailNow()
	Helper()
}

func ErrorIf(tb testingTB, err error, args ...any) bool {
	if err == nil {
		return false
	}
	tb.Helper()
	if len(args) > 0 {
		tb.Error(append([]any{err}, args...))
	} else {
		tb.Error(err)
	}
	return true
}

func DieIf(tb testingTB, err error, args ...any) {
	if ErrorIf(tb, err, args...) {
		tb.FailNow()
	}
}

func PanicIf(lg *log.Logger, err error, args ...any) {
	if err == nil {
		return
	}
	if lg == nil {
		lg = log.Default()
	}
	var s string
	if len(args) > 0 {
		s = fmt.Sprint(append([]any{err}, args...)...)
	} else {
		s = err.Error()
	}
	lg.Output(2, s)
	panic(s)
}
