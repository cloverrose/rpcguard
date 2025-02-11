package a06parameter

type App struct{}

type Message struct {
	text string
}

// returnParamErr returns err its source is parameter.
func (app *App) returnParamErr(err error) error { // want returnParamErr:"badFunc"
	return err
}

// returnParamFunc returns func its source is parameter.
func (app *App) returnParamFunc(fn func() error) func() error { // want returnParamFunc:"badFunc"
	return fn
}
