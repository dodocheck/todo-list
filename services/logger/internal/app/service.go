package app

import "context"

type Service struct {
	logger Logger
}

func NewService(logger Logger) *Service {
	return &Service{
		logger: logger,
	}
}

func (s *Service) Run(ctx context.Context) error {
	return s.logger.Run(ctx)
}

func (s *Service) Close(ctx context.Context) error {
	return s.logger.Close()
}
