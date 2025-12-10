package user

import (
	"context"
	"go-war-ticket-service/configs"
	"go-war-ticket-service/internal/domain"
	"go-war-ticket-service/internal/platform/storage"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type usecase struct {
	repo        Repository
	log         *zap.SugaredLogger
	minioClient *minio.Client
	cfg         configs.Config
}

func NewUsecase(
	r Repository,
	log *zap.SugaredLogger,
	minioClient *minio.Client,
	cfg configs.Config,
) Usecase {
	return &usecase{
		repo:        r,
		log:         log.Named("UserUsecase"),
		minioClient: minioClient,
		cfg:         cfg,
	}
}

func (u *usecase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		u.log.Errorf("failed to get user by ID: %v", err)
		return nil, domain.ErrInternal
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	if user.Avatar != "" {
		presignedUrl, err := storage.GetPresignedObject(
			u.minioClient,
			u.cfg.MinioBucket,
			user.Avatar,
			u.cfg.MinioEndpoint,
			u.cfg.MinioPublicEndpoint,
			time.Minute*15,
		)
		if err != nil {
			u.log.Error("failed to generate presigned url for avatar: ", err)
			return nil, domain.ErrInternal
		}

		user.Avatar = presignedUrl
	}

	return user, nil
}
