package app

import "context"

type Logger interface {
	Run(ctx context.Context) error
	Close() error
}
