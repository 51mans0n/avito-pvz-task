package logging

import "go.uber.org/zap"

// L  ― «сырое» *zap.Logger (нужен, когда хочется лог‑поля zap.Field).
var L *zap.Logger

// lg ― «сахарный» *zap.SugaredLogger (короткие методы Info / Error / Fatal).
var lg *zap.SugaredLogger

// Init инициализирует логгер.
// prod=false ➜ zap.NewDevelopment, prod=true ➜ zap.NewProduction.
func Init(prod bool) {
	var err error
	if prod {
		L, err = zap.NewProduction()
	} else {
		L, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
	lg = L.Sugar()
}

// Sync вызываем в `main` через defer.
func Sync() { _ = L.Sync() }

// S — получить «сахарный» логгер.
func S() *zap.SugaredLogger { return lg }
