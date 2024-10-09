package vm

import (
	"fishy/pkg/log"
	"fmt"
	"strings"
)

type ErrorCode int

func (e ErrorCode) String() string {
	if message, ok := ErrorMap[e]; ok {
		return message
	}
	return fmt.Sprintf("unknown error %d", int(e))
}

var ErrorMap = map[ErrorCode]string{
	UNKNOWN_SYSCAlL:  "unknown syscall",
	EPERM:            "operation not permitted",
	ENOENT:           "no such file or directory",
	ESRCH:            "no such process",
	EBADF:            "bad file descriptor",
	ECHILD:           "no child processes",
	EAGAIN:           "resource temporarily unavailable",
	EWOULDBLOCK:      "resource temporarily unavailable",
	EACCES:           "permission denied",
	EFAULT:           "bad address",
	EEXIST:           "file exists",
	ENODEV:           "no such device",
	ENOTDIR:          "not a directory",
	EISDIR:           "is a directory",
	EINVAL:           "invalid argument",
	ESPIPE:           "illegal seek",
	EMLINK:           "too many links",
	EPIPE:            "broken pipe",
	EDEADLK:          "resource deadlock avoided",
	EDEADLOCK:        "resource deadlock avoided",
	ENAMETOOLONG:     "file name too long",
	ENOSYS:           "function not implemented",
	ENOTEMPTY:        "directory not empty",
	ELOOP:            "too many levels of symbolic links",
	EBADMSG:          "bad message",
	EOVERFLOW:        "value too large for defined data type",
	ENOTUNIQ:         "name not unique on network",
	EBADFD:           "file descriptor in bad state",
	EREMCHG:          "remote address changed",
	ENOTSOCK:         "socket operation on non-socket",
	EDESTADDRREQ:     "destination address required",
	EMSGSIZE:         "message too long",
	ESOCKTNOSUPPORT:  "socket type not supported",
	EPFNOSUPPORT:     "protocol family not supported",
	EAFNOSUPPORT:     "address family not supported by protocol",
	EADDRINUSE:       "address already in use",
	EADDRNOTAVAIL:    "cannot assign requested address",
	ENETDOWN:         "network is down",
	ENETUNREACH:      "network is unreachable",
	ENETRESET:        "network dropped connection on reset",
	ECONNABORTED:     "software caused connection abort",
	ECONNRESET:       "connection reset by peer",
	ETIMEDOUT:        "connection timed out",
	ECONNREFUSED:     "connection refused",
	EHOSTDOWN:        "host is down",
	EHOSTUNREACH:     "no route to host",
	EFAILEDCREATE:    "failed to create",
	EADDROUTOFBOUNDS: "address out of bounds",
	EINVALIDLENGTH:   "invalid length",
	EBADHOSTADDRESS:  "bad host address",
}

const (
	UNKNOWN_SYSCAlL ErrorCode = iota + 1
	EPERM
	ENOENT
	ESRCH
	EBADF
	ECHILD
	EAGAIN
	EWOULDBLOCK
	EACCES
	EFAULT
	EEXIST
	ENODEV
	ENOTDIR
	EISDIR
	EINVAL
	ESPIPE
	EMLINK
	EPIPE
	EDEADLK
	EDEADLOCK
	ENAMETOOLONG
	ENOSYS
	ENOTEMPTY
	ELOOP
	EBADMSG
	EOVERFLOW
	ENOTUNIQ
	EBADFD
	EREMCHG
	ENOTSOCK
	EDESTADDRREQ
	EMSGSIZE
	ESOCKTNOSUPPORT
	EPFNOSUPPORT
	EAFNOSUPPORT
	EADDRINUSE
	EADDRNOTAVAIL
	ENETDOWN
	ENETUNREACH
	ENETRESET
	ECONNABORTED
	ECONNRESET
	ETIMEDOUT
	ECONNREFUSED
	EHOSTDOWN
	EHOSTUNREACH
	EFAILEDCREATE
	EADDROUTOFBOUNDS
	EINVALIDLENGTH
	EBADHOSTADDRESS
)

func MatchString(value string) ErrorCode {
	for code, message := range ErrorMap {
		if strings.Contains(strings.ToLower(message), strings.ToLower(value)) {
			return code
		}
	}
	log.Debug("MatchString", "value", value)
	return -1
}
