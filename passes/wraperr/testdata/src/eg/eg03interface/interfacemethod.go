package eg03interface

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type handler interface {
	Handle() error
}
type App struct {
	handler handler
}

type Message struct {
	text string
}

func (app *App) CallInterfaceMethod() error { // want CallInterfaceMethod:"badFunc"
	eg, _ := errgroup.WithContext(context.Background())
	eg.Go(app.handler.Handle)
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (app *App) CallInterfaceMethod2() error { // want CallInterfaceMethod2:"badFunc"
	eg, _ := errgroup.WithContext(context.Background())
	fn := app.handler.Handle
	eg.Go(fn)
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
