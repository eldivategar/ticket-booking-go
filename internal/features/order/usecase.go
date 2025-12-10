package order

import (
	"context"
	"fmt"
	"go-war-ticket-service/configs"
	"go-war-ticket-service/internal/domain"
	"go-war-ticket-service/internal/platform/storage"
	"go-war-ticket-service/internal/utils"
	"go-war-ticket-service/internal/utils/contextutil"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type usecase struct {
	repo        Repository
	log         *zap.SugaredLogger
	minioClient *minio.Client
	cfg         configs.Config
	cache       *redis.Client
}

func NewUsecase(
	r Repository,
	log *zap.SugaredLogger,
	minioClient *minio.Client,
	cfg configs.Config,
	cache *redis.Client,
) Usecase {
	return &usecase{
		repo:        r,
		log:         log.Named("OrderUsecase"),
		minioClient: minioClient,
		cfg:         cfg,
		cache:       cache,
	}
}

func (u *usecase) CreateOrder(ctx context.Context, order domain.Order) (*domain.Order, error) {
	if err := u.decreaseStockInRedis(ctx, order.EventID, order.Quantity); err != nil {
		return nil, err
	}

	currentUserID, _ := contextutil.GetUserID(ctx)
	takenUserID := strings.Split(currentUserID.String(), "-")[0]
	bookingID := fmt.Sprintf("WT-%s-%s", takenUserID, utils.GenerateRandomString(6))

	newOrder := domain.Order{
		BookingID: strings.ToUpper(bookingID),
		UserID:    currentUserID,
		EventID:   order.EventID,
		Quantity:  order.Quantity,
		Status:    domain.OrderStatusPending,
	}

	// Create new order in DB
	if err := u.repo.CreateOrder(ctx, &newOrder); err != nil {
		u.log.Errorf("failed to create order: %v", err)
		u.rollbackStock(ctx, order.EventID, order.Quantity)
		return nil, err
	}

	presignedUrl, _ := storage.GetPresignedObject(
		u.minioClient,
		u.cfg.MinioBucket,
		newOrder.Event.Image,
		u.cfg.MinioEndpoint,
		u.cfg.MinioPublicEndpoint,
		time.Minute*15,
	)

	newOrder.Event.Image = presignedUrl

	return &newOrder, nil
}

func (u *usecase) GetOrderByBookingID(ctx context.Context, bookingID string) (*domain.Order, error) {
	order, err := u.repo.GetOrderByBookingID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	// presigned url
	presignedUrl, _ := storage.GetPresignedObject(
		u.minioClient,
		u.cfg.MinioBucket,
		order.Event.Image,
		u.cfg.MinioEndpoint,
		u.cfg.MinioPublicEndpoint,
		time.Minute*15,
	)
	order.Event.Image = presignedUrl

	// presigned url ticket
	for i, ticket := range order.Ticket {
		order.Ticket[i].PDFUrl, err = storage.GetPresignedObject(
			u.minioClient,
			u.cfg.MinioBucket,
			ticket.PDFUrl,
			u.cfg.MinioEndpoint,
			u.cfg.MinioPublicEndpoint,
			time.Minute*15,
		)
		if err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (u *usecase) GetOrderList(ctx context.Context) ([]domain.Order, error) {
	currentUserID, _ := contextutil.GetUserID(ctx)

	orders, err := u.repo.GetOrderList(ctx, currentUserID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Order, len(orders))
	for i, order := range orders {
		// presigned url
		presignedUrl, _ := storage.GetPresignedObject(
			u.minioClient,
			u.cfg.MinioBucket,
			order.Event.Image,
			u.cfg.MinioEndpoint,
			u.cfg.MinioPublicEndpoint,
			time.Minute*15,
		)
		order.Event.Image = presignedUrl
		result[i] = order
	}

	return result, nil
}

// Decrease stock in Redis & Handle Cache Miss
func (u *usecase) decreaseStockInRedis(ctx context.Context, eventID uuid.UUID, qty int) error {
	redisKey := fmt.Sprintf(utils.EventStockKey, eventID.String())

	// Check if key exists
	exists, err := u.cache.Exists(ctx, redisKey).Result()
	if err != nil {
		return fmt.Errorf("redis error: %w", err)
	}

	// CACHE MISS: If key doesn't exist, load from DB and set to Redis
	if exists == 0 {
		u.log.Info("Cache miss detected, loading from DB...")
		event, err := u.repo.GetEventByID(ctx, eventID)
		if err != nil {
			return err
		}
		if event == nil {
			return domain.ErrEventNotFound
		}

		// SetNX (Set if Not Exists) to avoid race condition during re-hydrate
		// Set stock based on what's in the DB
		updated := u.cache.SetNX(ctx, redisKey, event.AvailableStock, 1*time.Hour).Val()
		if !updated {
			u.log.Warn("Race condition on cache re-hydration")
		}
	}

	// EXECUTE DECREMENT (Atomic Operation)
	remainingStock, err := u.cache.DecrBy(ctx, redisKey, int64(qty)).Result()
	if err != nil {
		return fmt.Errorf("failed to decrement redis: %w", err)
	}

	// Validate stock
	if remainingStock < 0 {
		// Stock is empty! Return the number (Increment back)
		u.cache.IncrBy(ctx, redisKey, int64(qty))
		return domain.ErrNotEnoughStock
	}

	return nil
}

func (u *usecase) rollbackStock(ctx context.Context, eventID uuid.UUID, qty int) {
	redisKey := fmt.Sprintf(utils.EventStockKey, eventID.String())
	u.cache.IncrBy(ctx, redisKey, int64(qty))
}
