package pdf

import (
	"fmt"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type Generator interface {
	GenerateTicket(data TicketData) ([]byte, error)
}
type TicketData struct {
	EventName        string
	EventLocation    string
	EventDate        time.Time
	EventImageBase64 string
	ImageExtension   consts.Extension
	OrderID          string
	TicketCode       string
}

type marotoGenerator struct{}

func NewMarotoGenerator() Generator {
	return &marotoGenerator{}
}

func (m *marotoGenerator) GenerateTicket(data TicketData) ([]byte, error) {
	// Setup: Portrait, A4, Margin 10
	p := pdf.NewMaroto(consts.Portrait, consts.A4)
	p.SetPageMargins(10, 10, 10)

	// --- COLORS ---
	darkGray := color.Color{Red: 50, Green: 50, Blue: 50}
	lightGray := color.Color{Red: 150, Green: 150, Blue: 150}

	// --- 1. HEADER IMAGE ---
	if data.EventImageBase64 != "" {
		p.Row(50, func() {
			p.Col(12, func() {
				_ = p.Base64Image(data.EventImageBase64, data.ImageExtension, props.Rect{
					Center:  true,
					Percent: 100,
				})
			})
		})
	}

	// Spacer
	p.Row(5, func() {})

	// --- 2. EVENT TITLE & LOCATION ---
	p.Row(10, func() {
		p.Col(12, func() {
			p.Text(data.EventName, props.Text{
				Size:  14,
				Style: consts.Bold,
				Align: consts.Left,
				Color: darkGray,
			})
		})
	})

	p.Row(15, func() {
		p.Col(12, func() {
			p.Text("Lokasi", props.Text{Size: 8, Color: lightGray, Top: 1})
			p.Text("📍 "+data.EventLocation, props.Text{
				Size:  10,
				Color: darkGray,
				Top:   5,
			})
		})
	})

	// Dashed Line Separator
	p.Line(1.0, props.Line{Color: lightGray, Style: consts.Dashed})
	p.Row(5, func() {}) // Spacer

	// --- 3. GRID INFO (Order ID, Code, Date, Time) ---

	// Row 1: Order ID & Ticket Code
	p.Row(15, func() {
		p.Col(6, func() {
			p.Text("Order ID", props.Text{Size: 8, Color: lightGray})
			p.Text(data.OrderID, props.Text{Size: 11, Style: consts.Bold, Color: darkGray, Top: 4})
		})
		p.Col(6, func() {
			p.Text("Kode Tiket", props.Text{Size: 8, Color: lightGray})
			p.Text(data.TicketCode, props.Text{Size: 11, Style: consts.Bold, Color: darkGray, Top: 4})
		})
	})

	// Format Date & Time
	dateStr := data.EventDate.Format("02 Jan 2006")
	timeStr := data.EventDate.Format("15:04") + " WIB"

	// Row 2: Event Date & Time
	p.Row(15, func() {
		p.Col(6, func() {
			p.Text("Tanggal Event", props.Text{Size: 8, Color: lightGray})
			p.Text(dateStr, props.Text{Size: 11, Style: consts.Bold, Color: darkGray, Top: 4})
		})
		p.Col(6, func() {
			p.Text("Waktu", props.Text{Size: 8, Color: lightGray})
			p.Text(timeStr, props.Text{Size: 11, Style: consts.Bold, Color: darkGray, Top: 4})
		})
	})

	p.Row(5, func() {}) // Spacer
	p.Line(1.0, props.Line{Color: lightGray, Style: consts.Dashed})
	p.Row(5, func() {}) // Spacer

	// --- 4. FOOTER (Disclaimer & QR) ---
	p.Row(35, func() {
		// Left: Disclaimer Text
		p.Col(8, func() {
			p.Text("Informasi Tiket", props.Text{Size: 9, Style: consts.Bold})
			p.Text("• Tunjukkan e-Tiket ini kepada panitia di lokasi.", props.Text{Size: 7, Color: darkGray, Top: 6})
			p.Text("• Wajib membawa kartu identitas yang berlaku.", props.Text{Size: 7, Color: darkGray, Top: 10})
			p.Text("• Dilarang membawa senjata tajam/obat terlarang.", props.Text{Size: 7, Color: darkGray, Top: 14})
		})

		// Right: QR Code
		p.Col(4, func() {
			p.QrCode(data.TicketCode, props.Rect{
				Center:  true,
				Percent: 100,
			})
			p.Text("Scan Code", props.Text{
				Size:  7,
				Align: consts.Center,
				Top:   30,
			})
		})
	})

	pdfBytes, err := p.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfBytes.Bytes(), nil
}
