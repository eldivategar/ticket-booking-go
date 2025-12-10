package ticket

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-war-ticket-service/configs"
	"go-war-ticket-service/internal/domain"
	"go-war-ticket-service/internal/features/order"
	"go-war-ticket-service/internal/platform/pdf"
	"go-war-ticket-service/internal/utils"
	"io"
	"strings"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type TicketWorker struct {
	mqConn      *amqp.Connection
	repo        Repository
	orderRepo   order.Repository
	minioClient *minio.Client
	cfg         configs.Config
	pdfGen      pdf.Generator
	log         *zap.SugaredLogger
}

func NewTicketWorker(
	conn *amqp.Connection,
	tr Repository,
	or order.Repository,
	mc *minio.Client,
	cfg configs.Config,
	pg pdf.Generator,
	logger *zap.SugaredLogger,
) *TicketWorker {
	return &TicketWorker{
		mqConn:      conn,
		repo:        tr,
		orderRepo:   or,
		minioClient: mc,
		cfg:         cfg,
		pdfGen:      pg,
		log:         logger,
	}
}

func (w *TicketWorker) Start() {
	ch, _ := w.mqConn.Channel()
	defer ch.Close()

	// Consume Queue
	msgs, _ := ch.Consume(
		utils.QueueTicketGeneration, // Name of queue
		"",                          // consumer name (empty)
		false,                       // Auto-Ack (set FALSE to manually Ack after success)
		false,                       // exclusive
		false,                       // no-local
		false,                       // no-wait
		nil,                         // args
	)

	// Loop forever waiting for messages
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			if err := w.processMessage(d); err != nil {
				w.log.Errorf("Error processing message: %v", err)
				// Option: d.Nack(false, true) if you want to retry
			}
		}
	}()

	<-forever
}

func (w *TicketWorker) processMessage(d amqp.Delivery) error {
	var payload struct {
		BookingID     string `json:"booking_id"`
		PaymentStatus string `json:"payment_status"`
	}

	if err := json.Unmarshal(d.Body, &payload); err != nil {
		return err
	}

	w.log.Infof("Processing PDF with Booking ID: %s\n", payload.BookingID)

	order, err := w.orderRepo.GetOrderByBookingID(context.Background(), payload.BookingID)
	if err != nil {
		return err
	}

	// Get event image
	object, err := w.minioClient.GetObject(
		context.Background(),
		w.cfg.MinioBucket,
		order.Event.Image,
		minio.GetObjectOptions{},
	)
	if err != nil {
		w.log.Errorf("Warning: Failed to get image from minio: %v", err)
	}
	defer object.Close()

	var imageBase64 string

	// Cek apakah object valid
	stat, err := object.Stat()
	if err == nil && stat.Size > 0 {
		imgBytes := make([]byte, stat.Size)
		_, err := io.ReadFull(object, imgBytes)
		if err == nil {
			imageBase64 = base64.StdEncoding.EncodeToString(imgBytes)
		}
	}

	// Get image extension
	var imgExtension consts.Extension = consts.Jpg // Default
	if strings.HasSuffix(strings.ToLower(order.Event.Image), ".png") {
		imgExtension = consts.Png
	}

	// Create new ticket
	for i := 0; i < order.Quantity; i++ {
		ticketNumber := fmt.Sprintf("TIK-%s-%s", order.BookingID, utils.GenerateRandomNumberString(3))

		// Generate PDF with detail order
		pdfData := pdf.TicketData{
			EventName:        order.Event.Name,
			EventLocation:    order.Event.Location,
			EventDate:        order.Event.Date,
			EventImageBase64: imageBase64,
			ImageExtension:   imgExtension,
			OrderID:          order.BookingID,
			TicketCode:       ticketNumber,
		}

		pdfBytes, err := w.pdfGen.GenerateTicket(pdfData)
		if err != nil {
			w.log.Errorf("failed to generate PDF: %v", err)
			return err
		}

		// Save PDF to S3
		filename := fmt.Sprintf("%s.pdf", ticketNumber)
		objectName := fmt.Sprintf("tickets/%s", filename)
		reader := bytes.NewReader(pdfBytes)

		_, err = w.minioClient.PutObject(
			context.Background(),
			w.cfg.MinioBucket,
			objectName,
			reader,
			int64(len(pdfBytes)),
			minio.PutObjectOptions{ContentType: "application/pdf"},
		)
		if err != nil {
			w.log.Errorf("failed to upload PDF to S3: %v", err)
			return err
		}

		path := objectName

		ticket := domain.Ticket{
			OrderID:      order.ID,
			EventID:      order.EventID,
			UserID:       order.UserID,
			TicketNumber: ticketNumber,
			PDFUrl:       path,
		}

		if err := w.repo.CreateTicket(context.Background(), &ticket); err != nil {
			return err
		}
	}

	// Update order status to completed
	if err := w.orderRepo.UpdateOrderStatus(context.Background(), order.BookingID, domain.OrderStatusCompleted); err != nil {
		return err
	}

	d.Ack(false)

	w.log.Infof("PDF generated and saved to S3 for Booking ID: %s\n", payload.BookingID)
	return nil
}
