package helpers

import "fmt"

// New обертка для создания новых ошибок
func New(msg string, args ...any) error {
	return fmt.Errorf(msg, args...)
}

// Wrap оборачивает ошибки для прокидывания наверх по стеку вызова
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// Wrapf параметризированое оборачивание ошибок
func Wrapf(err error, msg string, args ...any) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(msg, args...), err)
}
