package errors

// Base error structure
type BaseError struct {
	Message string
}

func (s *BaseError) Error() string { return s.Message }

// When human being rejects action
type CallRejected BaseError
func (s *CallRejected) Error() string { return s.Message }

// When tofa client rejects action
type CallForbidden BaseError
func (s *CallForbidden) Error() string { return s.Message }

// When server sends malformed request (aka code 400)
type BadCall BaseError
func (s *BadCall) Error() string { return s.Message }

// When client does not reply in mean time
// Nothing which server can fix from its side
type CallTimedOut BaseError
func (s *CallTimedOut) Error() string { return s.Message }

// Used in preflight/ping requests
// When connection doesn't establish in mean time
type ConnectTimedOut BaseError
func (s *ConnectTimedOut) Error() string { return s.Message }

// Fault's on client side
// When there is a conflict between client GUI and Daemon.
// Nothing which server can fix from its side
type ClDaConflict BaseError
func (s *ClDaConflict) Error() string { return s.Message }

// Fired amid unexpected response codes from client
type UnsupportedResponseCode BaseError
func (s *UnsupportedResponseCode) Error() string { return s.Message }

// Fired amid bad URI
type BadURI BaseError
func (s *BadURI) Error() string { return s.Message }

// Amid unsuported URI
type UnsupportedURI BaseError
func (s *UnsupportedURI) Error() string { return s.Message }

// When request failed due to different causes
type RequestFailed BaseError
func (s *RequestFailed) Error() string { return s.Message }

// When client is busy processing another call to the same app
type CallConflicts BaseError
func (s *CallConflicts) Error() string { return s.Message }

// Retrieves error based on code
func GetErrorByCode(code int) error {
	switch code {
	case 403:
		return &CallForbidden{"call forbidden"}
	case 400:
		return &BadCall{"bad call"}
	case 408:
		return &CallTimedOut{"call timed out"}
	case 409:
		return &CallConflicts{"call conflicts"}
	case 570:
		return &CallRejected{"call rejected"}
	case 571:
		return &ClDaConflict{"client gui<->daemon conflict"}
	default:
		return &UnsupportedResponseCode{"unsupported response code"}
	}
}
