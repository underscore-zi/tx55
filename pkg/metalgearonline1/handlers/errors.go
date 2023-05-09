package handlers

var ErrUnexpectedArgument = NewError(-1, "unexpected argument")
var ErrNotImplemented = NewError(-2, "not implemented")
var ErrHandlerNotFound = NewError(-3, "handler not found")
var ErrNotFound = NewError(-3, "not found")
var ErrInvalidArguments = NewError(-4, "invalid arguments")
var ErrDatabase = NewError(-5, "database error")
var ErrNotHosting = NewError(-6, "not hosting")
var ErrBanned = NewError(-7, "banned")

type GameError struct {
	Code    int32
	Message string
}

func (e GameError) Error() string {
	return e.Message
}

func NewError(code int32, message string) GameError {
	return GameError{
		Code:    code,
		Message: message,
	}
}
