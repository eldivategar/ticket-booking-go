package event

import (
	"context"
	"fmt"
	"go-war-ticket-service/configs"
	"go-war-ticket-service/internal/domain"
	"go-war-ticket-service/internal/platform/storage"
	"go-war-ticket-service/internal/utils"
	"strings"
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
		log:         log.Named("EventUsecase"),
		minioClient: minioClient,
		cfg:         cfg,
	}
}

func (u *usecase) CreateEvent(ctx context.Context, event domain.Event) (*domain.Event, error) {
	if event.TotalStock <= 0 {
		return nil, domain.ErrInvalidStock
	}

	if event.Price <= 0 {
		return nil, domain.ErrInvalidPrice
	}

	if event.Date.IsZero() {
		return nil, domain.ErrInvalidDate
	}

	// Save image to minio
	eventName := strings.ReplaceAll(event.Name, " ", "-")
	filename := fmt.Sprintf("%s-%s", eventName, utils.GenerateRandomNumberString(6))
	path, err := storage.UploadImageToMinIO(
		u.minioClient,
		u.cfg.MinioBucket,
		event.Image,
		"events",
		filename,
	)
	if err != nil {
		u.log.Error("failed to upload image to minio: ", err)
		return nil, domain.ErrInternal
	}

	event.Image = path

	return u.repo.CreateEvent(ctx, event)
}

func (u *usecase) GetEventByID(ctx context.Context, eventID uuid.UUID) (*domain.Event, error) {
	event, err := u.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, domain.ErrEventNotFound
	}

	presignedURL, err := storage.GetPresignedObject(
		u.minioClient,
		u.cfg.MinioBucket,
		event.Image,
		u.cfg.MinioEndpoint,
		u.cfg.MinioPublicEndpoint,
		time.Minute*15,
	)
	if err != nil {
		u.log.Error("failed to generate presigned url for image: ", err)
		return nil, domain.ErrInternal
	}

	event.Image = presignedURL

	return event, nil
}

func (u *usecase) GetAllEvent(ctx context.Context) ([]domain.Event, error) {
	events, err := u.repo.GetAllEvent(ctx)
	if err != nil {
		return nil, err
	}

	for i := range events {
		presignedURL, err := storage.GetPresignedObject(
			u.minioClient,
			u.cfg.MinioBucket,
			events[i].Image,
			u.cfg.MinioEndpoint,
			u.cfg.MinioPublicEndpoint,
			time.Minute*15,
		)
		if err != nil {
			u.log.Error("failed to generate presigned url for image: ", err)
			return nil, domain.ErrInternal
		}

		events[i].Image = presignedURL
	}

	return events, nil
}

func (u *usecase) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	event, err := u.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event == nil {
		return domain.ErrEventNotFound
	}

	// Delete image from minio
	if err := storage.DeleteObjectFromMinIO(
		u.minioClient,
		u.cfg.MinioBucket,
		event.Image,
	); err != nil {
		u.log.Error("failed to delete image from minio: ", err)
		return domain.ErrInternal
	}

	return u.repo.DeleteEvent(ctx, eventID)
}
