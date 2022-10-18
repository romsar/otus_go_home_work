package closer

import "github.com/rs/zerolog/log"

var globalCloser = New()

type CloseFunc func() error

// Closer характеризует структуру с функциями для закрытия соединений.
type Closer struct {
	funcs []CloseFunc
}

// New создает Closer.
func New() *Closer {
	return &Closer{
		funcs: make([]CloseFunc, 0),
	}
}

// Add добавляет функцию закрытия глобальному Closer.
func Add(f ...CloseFunc) {
	globalCloser.Add(f...)
}

// CloseAll закрывает все соединения у глобального Closer.
func CloseAll() {
	globalCloser.CloseAll()
}

// Add добавляет функцию закрытия.
func (c *Closer) Add(f ...CloseFunc) {
	c.funcs = append(c.funcs, f...)
}

// CloseAll закрывает все соединения Closer.
func (c *Closer) CloseAll() {
	for _, f := range c.funcs {
		if err := f(); err != nil {
			log.Err(err).Msg("error close")
		}
	}

	c.funcs = make([]CloseFunc, 0)
}
