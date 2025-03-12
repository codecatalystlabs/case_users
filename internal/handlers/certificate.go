package handlers

import (
	"bytes"
	"case/internal/models"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

func VerifyDischarge(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	return nil
}

func Discharge(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	//=================

	userID := GetCurrentUser(c, store)
	fmt.Println("we here 2")
	// Check if user is logged in
	if userID == 0 {
		fmt.Println("unauthorized")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	//=============================

	var formData map[string]interface{}

	if err := c.BodyParser(&formData); err != nil {
		fmt.Println("JSON parsing failed:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var s models.Discharge

	s.DischargeID = int(ParseNullInt2(formData["discharge_id"]).Int64)
	s.ClientID = ParseNullInt2(formData["client_id"])
	s.DischargeDate = ParseNullString2(formData["discharge_date"])
	s.FinalDiagnosis = ParseNullString2(formData["final_diagnosis"])
	s.FinalDiagnosisOther = ParseNullString2(formData["final_diagnosis_other"])
	s.DischargeOutcome = ParseNullString2(formData["discharge_outcome"])
	s.DischargeSeqHeari = ParseNullInt2(formData["discharge_seq_heari"])
	s.DischargeSeqPregn = ParseNullInt2(formData["discharge_seq_pregn"])
	s.DischargeSeqOcula = ParseNullInt2(formData["discharge_seq_ocula"])
	s.DischargeSeqExtre = ParseNullInt2(formData["discharge_seq_extre"])
	s.DischargeSeqArthr = ParseNullInt2(formData["discharge_seq_arthr"])
	s.DischargeSeqNeuro = ParseNullInt2(formData["discharge_seq_neuro"])
	s.DischargeSeqOthers = ParseNullInt2(formData["discharge_seq_others"])
	s.CounsellingProvided = ParseNullString2(formData["counselling_provided"])
	s.DischargingOfficer = ParseNullString2(formData["discharging_officer"])
	s.DischargeFacility = ParseNullString2(formData["discharge_facility"])
	s.DischargeSeqOthersAza = ParseNullString2(formData["discharge_seq_others_aza"])

	s.UpdatedBy.Valid = true
	s.UpdatedBy.Int64 = int64(userID)

	s.UpdatedOn.Valid = true
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02")
	s.UpdatedOn.String = formattedTime

	fmt.Println(s)
	// Check if it's a new record (StatusID == 0)
	if s.DischargeID > 0 {
		s.EnteredBy.Valid = true
		s.EnteredBy.Int64 = int64(userID)

		s.EnteredOn.Valid = true
		s.EnteredOn.String = formattedTime

		s.SetAsExists()
		err := s.Update(c.Context(), db)
		if err != nil {
			fmt.Println("update fail:", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	} else {

		err := s.Insert(c.Context(), db)
		if err != nil {
			fmt.Println("insert fail:", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	fmt.Println("we good")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func GetDischarge(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	userID := GetCurrentUser(c, store)

	// Check if user is logged in
	if userID == 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	fmt.Println("get discharge")

	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = "0"
	}

	c_id, err := strconv.Atoi(clientID)
	if err != nil {
		c_id = 0
	}

	discharge, er := models.DischargeByClientID(c.Context(), db, c_id)
	if er != nil {
		fmt.Println(er.Error())
		return nil
	}

	return c.JSON(discharge)
}

func Certificate(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {

	clientID := c.Query("who")
	if clientID == "" {
		clientID = "0"
	}

	c_id, er := strconv.Atoi(clientID)
	if er != nil {
		c_id = 0
	}

	discharge, erx := models.DischargeByClientID(c.Context(), db, c_id)
	if erx != nil {
		fmt.Println(erx.Error())
		return nil
	}

	client, erx := models.ClientByID(c.Context(), db, c_id)
	if erx != nil {
		fmt.Println(erx.Error())
		return nil
	}

	facility, erx := models.FacilityByFacilityID(c.Context(), db, int(client.Site.Int64))
	if erx != nil {
		fmt.Println(erx.Error())
		return nil
	}

	// Create a new PDF
	fmt.Println("Starting certificate...")
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Load Fonts
	pdf.SetFont("Arial", "B", 10)

	// Add Ministry of Health Logo
	logoPath := "../../ui/static/img/logo.png"
	pdf.Image(logoPath, 20, 10, 30, 0, false, "", 0, "") // Centered Logo
	//pdf.Ln(35)                                           // Extra spacing after logo

	// Generate QR Code
	qrLink := "response.health.go.ug/discharge/verify" // Replace with the actual verification link
	qrFile := "qrcode.png"
	err := qrcode.WriteFile(qrLink, qrcode.Medium, 256, qrFile)
	if err != nil {
		fmt.Println(err.Error())
		//panic(err)
	}

	// Add QR Code to PDF
	pdf.Image(qrFile, 240, 10, 30, 30, false, "", 0, "")
	_ = os.Remove(qrFile) // Cleanup QR Code file

	// Ministry of Health Title
	pageWidth, _ := pdf.GetPageSize()
	fmt.Println(pageWidth)

	logo_text := "Ministry of Health, Republic of Uganda"
	textWidth := pdf.GetStringWidth(logo_text)
	logo_x := (pageWidth - textWidth) / 2
	fmt.Println(logo_x)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(260, 20, logo_text, "0", 1, "C", false, 0, "")
	//pdf.Ln(5)

	// Certificate Title
	title_text := "EVD Discharge Certificate"
	title_textWidth := pdf.GetStringWidth(title_text)
	title_x := (pageWidth - title_textWidth) / 2
	fmt.Println(title_x)
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(260, 10, title_text, "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Certification Text

	body1_text := "This is to certify that"
	body1_textWidth := pdf.GetStringWidth(body1_text)
	body1_x := (pageWidth - body1_textWidth) / 2
	fmt.Println(body1_x)
	pdf.SetFont("Arial", "i", 14)
	pdf.CellFormat(260, 10, body1_text, "0", 1, "C", false, 0, "")
	pdf.Ln(6)

	// Name (Handwriting Style)
	pdf.SetFont("Courier", "B", 22) // Courier gives a handwritten feel
	pdf.CellFormat(260, 12, client.Firstname.String+" "+client.Lastname.String+" ("+client.EtuNo.String+")", "0", 1, "C", false, 0, "")
	pdf.Ln(1)

	// Draw Line Below Name
	pdf.Line(55, pdf.GetY(), 260, pdf.GetY()) // Draw a horizontal line
	pdf.Ln(5)

	// Rest of the Certification Text
	pdf.SetFont("Arial", "", 12)
	text1 := "At the date (" + discharge.DischargeDate.String + ") of issue of this certificate, does not present a risk of infecting other persons after testing negative for Ebola\nVirus Disease. The current state of health does not constitute a danger to the community and can therefore, return to their\nhousehold and professional environment to continue their normal daily activities"
	pdf.MultiCell(0, 6, text1, "", "C", false)
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 12)
	text2 := `The family, the community, and the authorities are requested to accept them in order to promote their social integration`
	pdf.MultiCell(0, 10, text2, "", "C", false)
	pdf.Ln(5)

	// Draw a box around the last three words

	pdf.SetFont("Arial", "", 14)
	firstwords := `Completed at facility: `

	pdf.SetFont("Arial", "B", 12) // Set bold font
	lastWords := facility.FacilityName.String

	all_words := firstwords + lastWords

	// Draw the box
	pdf.Rect(20, pdf.GetY(), 250, 10, "D") // Adjust the size of the box (190 width, 30 height)

	// Write the text in the box
	pdf.MultiCell(0, 10, all_words, "", "C", false)

	pdf.Ln(10) // Line break after the text

	// Signatures Section
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(260, 10, "_________________________", "0", 0, "C", false, 0, "")
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(260, 10, "Ebola Treatment Unit Manager", "0", 1, "C", false, 0, "")
	pdf.CellFormat(260, 10, facility.FacilityName.String, "0", 1, "C", false, 0, "")
	pdf.Ln(15)

	// Save PDF to Buffer (In-Memory)
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		fmt.Println("Cert error: ", err.Error())
		//return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed, contact admin"})
	}

	// Set HTTP Headers and Return PDF
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "inline; filename=certificate.pdf")
	c.Set("Content-Length", fmt.Sprintf("%d", buf.Len())) // Set content length for the browser
	return c.Send(buf.Bytes())

}
