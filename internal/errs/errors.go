package errs

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
)

var (
	NoConnect          = errors.New("отсутствует соединение с сервером. Возможно, проблемы с сетью")
	Technic            = errors.New("сервис недоступен. Попробуйте позже")
	FnsTokenErr        = errors.New("законлся токен ФНС")
	FnsRefreshTokenErr = errors.New("законлся refresh-токен ФНС")
)

type BusinessError struct {
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return err
	}
	file = filepath.Base(file)
	return fmt.Errorf("(%s:%d): %w", file, line, err)
}
