package order

import (
	"context"
	"go-war-ticket-service/internal/domain"
	rabbitmq "go-war-ticket-service/internal/platform/message_broker/rabbit_mq"
	"go-war-ticket-service/internal/utils"

	"go.uber.org/zap"
)

type service struct {
	repo Repository
	log  *zap.SugaredLogger
	mq   rabbitmq.Publisher
}

func NewService(
	repo Repository,
	log *zap.SugaredLogger,
	mq rabbitmq.Publisher,
) Service {
	return &service{
		repo: repo,
		log:  log.Named("OrderService"),
		mq:   mq,
	}
}

func (s *service) ProcessPaymentWebhook(ctx context.Context, payload PaymentWebhookRequest) error {
	if payload.PaymentStatus != "PAID" && payload.PaymentStatus != "SETTLEMENT" {
		s.log.Info("Payment status not paid, ignoring...")
		return nil
	}

	order, err := s.repo.GetOrderByBookingID(ctx, payload.BookingID)
	if err != nil {
		return err
	}

	if order.Status == domain.OrderStatusPaid || order.Status == domain.OrderStatusCompleted {
		s.log.Info("Order already paid, ignoring...")
		return nil
	}

	if err := s.repo.UpdateOrderStatus(ctx, order.BookingID, domain.OrderStatusPaid); err != nil {
		return err
	}

	body := map[string]interface{}{
		"booking_id": order.BookingID,
		"status":     order.Status,
	}

	err = s.mq.Publish(ctx, utils.QueueTicketGeneration, body)
	if err != nil {
		s.log.Errorf("failed to publish message: %v", err)
		return err
	}

	return nil
}
