package exec

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type RetryableExecutor[T any] interface {
	ExecWithFallback(execFn func() (T, error), fallbackFn func() error) (T, error)
}

type retryableExecutor[T any] struct {
	opts *Opts
}

func NewRetryableExecutor[T any](opts ...OptsFn) RetryableExecutor[T] {
	execOpts := NewDefaultExecOpts()

	for _, opt := range opts {
		opt(execOpts)
	}

	return &retryableExecutor[T]{
		execOpts,
	}
}

func (r *retryableExecutor[T]) ExecWithFallback(execFn func() (T, error), fallbackFn func() error) (res T, err error) {
	if execFn == nil {
		return res, ErrMissingExecFn
	}

	if fallbackFn == nil {
		return res, ErrMissingFallbackFn
	}

	execErr := &Error{}

	retryCount := uint(0)

	for {
		res, err = execFn()

		if err == nil {
			return res, nil
		}

		execErr.AddErr(fmt.Errorf("exec function error: %w", err))

		if retryCount == r.opts.maxRetryCount {
			return res, execErr
		}

		if err = fallbackFn(); err != nil && !r.opts.retryOnFallbackError {
			execErr.AddErr(fmt.Errorf("fallback function error: %w", err))

			return res, execErr
		}

		retryCount++

		time.Sleep(r.opts.errTimeout)
	}
}

var (
	ErrMissingExecFn     = errors.New("no exec function provided")
	ErrMissingFallbackFn = errors.New("no fallback function provided")
)

const (
	defaultMaxRetryCount        = 3
	defaultErrTimeout           = 0 * time.Second
	defaultRetryOnFallbackError = true
)

type Opts struct {
	maxRetryCount        uint
	errTimeout           time.Duration
	retryOnFallbackError bool
}

func NewDefaultExecOpts() *Opts {
	return &Opts{
		maxRetryCount:        defaultMaxRetryCount,
		errTimeout:           defaultErrTimeout,
		retryOnFallbackError: defaultRetryOnFallbackError,
	}
}

type OptsFn func(opts *Opts)

func WithMaxRetryCount(maxCount uint) OptsFn {
	return func(opts *Opts) {
		if maxCount == 0 {
			maxCount = defaultMaxRetryCount
		}

		opts.maxRetryCount = maxCount
	}
}

func WithErrTimeout(errTimeout time.Duration) OptsFn {
	return func(opts *Opts) {
		opts.errTimeout = errTimeout
	}
}

func WithRetryOnFallBackError(retryOnFallbackError bool) OptsFn {
	return func(opts *Opts) {
		opts.retryOnFallbackError = retryOnFallbackError
	}
}

type Error struct {
	errs []error
}

func (e *Error) AddErr(err error) {
	e.errs = append(e.errs, err)
}

func (e *Error) Error() string {
	sb := strings.Builder{}

	for i, err := range e.errs {
		sb.WriteString(fmt.Sprintf("error %d: %s\n", i, err))
	}

	return sb.String()
}
