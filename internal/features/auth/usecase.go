package auth

import "go.uber.org/zap"

type usecase struct {
	repo Repository
	log  *zap.SugaredLogger
}

func NewUsecase(
	r Repository,
	log *zap.SugaredLogger,
) Usecase {
	return &usecase{
		repo: r,
		log:  log.Named("AuthUsecase"),
	}
}
