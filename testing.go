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

func DieIf(tb testingTB, err error, args ...any) {
	if err != nil {
		tb.Helper()
		tb.Error(err)
		if len(args) > 0 {
			tb.Error(fmt.Sprint(append([]any{err}, args...)...))
		} else {
			tb.Error(err)
		}
		tb.FailNow()
	}
}

func PanicIf(lg *log.Logger, err error, args ...any) {
	if err != nil {
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
}
