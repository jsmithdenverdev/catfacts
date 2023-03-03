package catfacts

import "github.com/chuckpreslar/emission"

func result[TIn, TOk any, TErr error](delegate func(TIn) (TOk, TErr), ok func(TOk), err func(TErr)) func(TIn) {
	return func(in TIn) {
		result, e := delegate(in)
		a, _ := any(e).(error)
		if a != nil {
			err(a.(TErr))
			return
		}
		ok(result)
	}
}

func emitResult[TIn, TOk any, TError error](emitter *emission.Emitter, okTopic string, errTopic string) func(delegate func(TIn) (TOk, TError)) func(TIn) {
	ok := func(ok TOk) {
		emitter.Emit(okTopic, ok)
	}

	err := func(err TError) {
		emitter.Emit(errTopic, err)
	}

	return func(delegate func(TIn) (TOk, TError)) func(TIn) {
		return result[TIn, TOk, TError](delegate, ok, err)
	}
}
