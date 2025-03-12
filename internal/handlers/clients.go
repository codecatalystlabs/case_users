package handlers

import (
	"case/internal/models"
	"case/internal/security"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var CtxG *fiber.Ctx
var dbG *sql.DB

func HandlerCasesForm(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	fmt.Println("starting case form")

	userID, userName := GetUser(c, sl, store)
	role := security.GetRoles(userID, "admin")
	id, err := strconv.Atoi(c.Params("i"))
	data := NewTemplateData(c, store)

	//fmt.Println("za id: ", id)

	var client models.Client

	if err != nil || id < 1 {
		client.ID = 0
		data.IsIDPos = false
	} else {
		c, err := models.ClientByID(c.Context(), db, id)
		if err == nil {
			client = *c
		}
		fmt.Println("za id zi: ", c.ID)
		data.IsIDPos = true
	}

	data.User = userName
	data.Role = role
	data.Optionz = Get_Client_Optionz()
	data.Form = client

	fmt.Println("loading case form page")
	return GenerateHTML(c, data, "form_patients")
}
func HandlerCasesSubmit(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {

	id, er := strconv.Atoi(c.FormValue("id"))
	if er != nil {
		id = 0
	}

	client := models.Client{
		ID:               id,
		Firstname:        ParseNullString(c.FormValue("firstname")),
		Lastname:         ParseNullString(c.FormValue("lastname")),
		Othername:        ParseNullString(c.FormValue("othername")),
		Gender:           ParseNullInt(c.FormValue("gender")),
		DateOfBirth:      ParseNullString(c.FormValue("date_of_birth")),
		Age:              ParseNullFloat(c.FormValue("age")),
		Marital:          ParseNullInt(c.FormValue("marital")),
		Nin:              ParseNullString(c.FormValue("nin")),
		Nationality:      ParseNullInt(c.FormValue("nationality")),
		AdmDate:          ParseNullString(c.FormValue("adm_date")),
		AdmFrom:          ParseNullString(c.FormValue("adm_from")),
		LabNo:            ParseNullString(c.FormValue("lab_no")),
		CifNo:            ParseNullString(c.FormValue("cif_no")),
		EtuNo:            ParseNullString(c.FormValue("etu_no")),
		CaseNo:           ParseNullString(c.FormValue("case_no")),
		Occupation:       ParseNullInt(c.FormValue("occupation")),
		OccupationAza:    ParseNullString(c.FormValue("occupation_aza")),
		DateSymptomOnset: ParseNullString(c.FormValue("date_symptom_onset")),
		DateIsolation:    ParseNullString(c.FormValue("date_isolation")),
		Pregnant:         ParseNullInt(c.FormValue("pregnant")),
		AdmWard:          ParseNullString(c.FormValue("adm_ward")),
		Tb:               ParseNullInt(c.FormValue("tb")),
		Asplenia:         ParseNullInt(c.FormValue("asplenia")),
		Hep:              ParseNullInt(c.FormValue("hep")),
		Diabetes:         ParseNullInt(c.FormValue("diabetes")),
		Hiv:              ParseNullInt(c.FormValue("hiv")),
		Liver:            ParseNullInt(c.FormValue("liver")),
		Malignancy:       ParseNullInt(c.FormValue("malignancy")),
		Heart:            ParseNullInt(c.FormValue("heart")),
		Pulmonary:        ParseNullInt(c.FormValue("pulmonary")),
		Kidney:           ParseNullInt(c.FormValue("kidney")),
		Neurologic:       ParseNullInt(c.FormValue("neurologic")),
		Other:            ParseNullString(c.FormValue("other")),
		Transfer:         ParseNullInt(c.FormValue("transfer")),
		Site:             ParseNullInt(c.FormValue("site")),
		Status:           ParseNullString(c.FormValue("status")),

		//Status: ParseNullString(c.FormValue("status")),
	}

	//visID, _ := utilities.GetSequentialVisitID()
	userID := GetCurrentUser(c, store)

	client.EditOn.Valid = true
	client.EditBy.Valid = true

	client.EditBy.Int64 = int64(userID)
	client.EditOn.Time = time.Now()

	if client.ID == 0 {

		client.EnterOn.Valid = true
		client.EnterBy.Valid = true

		client.EnterBy.Int64 = int64(userID)
		client.EnterOn.Time = time.Now()

		client.UUID.Valid = true
		client.UUID.String = models.CreateUUID()

		//appID := models.CreateUUID()
		//client.UUID.String = appID

	}
	fmt.Println("Created 0")
	//return
	fmt.Println(client)
	fmt.Println("Created 1")

	if client.ID == 0 {
		err := client.Insert(c.Context(), db)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		client.SetAsExists()
		err := client.Update(c.Context(), db)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	urlx := "/cases/new/" + strconv.Itoa(client.ID)

	return c.Redirect(urlx)
}
func HandlerCasesList(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	fmt.Println("starting case list")

	userID, userName := GetUser(c, sl, store)
	role := security.GetRoles(userID, "admin")

	data := NewTemplateData(c, store)
	data.User = userName
	data.Role = role

	fmt.Println("loading case list page")
	//

	filter := ""
	clients, err := models.Clients(c.Context(), db, filter)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			fmt.Println("error loading case list: ", err.Error())
		} else {
			fmt.Println("error loading case list: ", err.Error())
		}
	}

	data.Form = clients

	//
	return GenerateHTML(c, data, "list_patients")
}

func HandlerCaseEncounterForm(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	fmt.Println("starting case encounter form")

	CtxG = c
	dbG = db

	userID, userName := GetUser(c, sl, store)
	role := security.GetRoles(userID, "admin")

	id, er := strconv.Atoi(c.Params("i"))
	if er != nil {
		id = 0
	}

	jd, er := strconv.Atoi(c.Params("j"))
	if er != nil {
		jd = 0
	}

	var client *models.Client
	var encounter = &models.Encounter{}
	var clinical = &models.Clinical{}
	var vital = &models.Vital{}

	encounter.ClientID.Valid = true
	encounter.ClientID.Int64 = 0

	clinical.EncounterID.Valid = true
	clinical.EncounterID.Int64 = 0

	vital.EncounterID.Valid = true
	vital.EncounterID.Int64 = 0

	if id > 0 {
		client, _ = models.ClientByID(c.Context(), db, id)
	}

	encounter.ClientID.Int64 = int64(client.ID)

	if jd > 0 {
		encounter, _ = models.EncounterByEncounterID(c.Context(), db, jd)
		clinical.EncounterID.Int64 = int64(encounter.EncounterID)
		vital.EncounterID.Int64 = int64(encounter.EncounterID)

		// add logic to load child forms with encounter ID
	}
	fmt.Println("loading 3")
	data := NewTemplateData(c, store)

	data.User = userName
	data.Role = role
	data.Optionz = Get_Client_Optionz()
	data.Form = encounter
	data.FormRef = client
	data.FormChild1 = clinical
	data.FormChild2 = vital

	fmt.Println("loading case encounter form page")
	return GenerateHTML(c, data, "form_encounters")
}

func HandlerCaseEncounterSubmit(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	return nil
}
