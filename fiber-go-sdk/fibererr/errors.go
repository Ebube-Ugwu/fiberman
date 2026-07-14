package fibererr

import "fmt"

type FiberError struct {
	Message string
	Cause   error
}

func (e *FiberError) Error() string {
	return e.Message
}

func (e *FiberError) Unwrap() error {
	return e.Cause
}

type HTTPError struct {
	FiberError
	StatusCode   int
	ResponseBody string
}

type RPCError struct {
	FiberError
	Code int
	Data any
}

type SerializationError struct {
	FiberError
}

type TransportError struct {
	FiberError
}

type TimeoutError struct {
	TransportError
}

func NewHTTPError(message string, statusCode int, responseBody string) *HTTPError {
	return &HTTPError{
		FiberError:   FiberError{Message: message},
		StatusCode:   statusCode,
		ResponseBody: responseBody,
	}
}

func NewRPCError(message string, code int, data any) *RPCError {
	return &RPCError{
		FiberError: FiberError{Message: message},
		Code:       code,
		Data:       data,
	}
}

func NewSerializationError(message string, cause error) *SerializationError {
	return &SerializationError{
		FiberError: FiberError{Message: message, Cause: cause},
	}
}

func NewTransportError(message string, cause error) *TransportError {
	return &TransportError{
		FiberError: FiberError{Message: message, Cause: cause},
	}
}

func NewTimeoutError(message string, cause error) *TimeoutError {
	return &TimeoutError{
		TransportError: TransportError{
			FiberError: FiberError{Message: message, Cause: cause},
		},
	}
}

func FormatHTTPMessage(method string, statusCode int) string {
	return fmt.Sprintf("Fiber RPC call failed for method '%s' with HTTP status %d", method, statusCode)
}

func FormatRPCMessage(method string, code int, rpcMessage string) string {
	return fmt.Sprintf("Fiber RPC error for method '%s': [%d] %s", method, code, rpcMessage)
}
