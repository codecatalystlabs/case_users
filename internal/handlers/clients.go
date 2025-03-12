package handlers

import (
	"case/internal/models"
	"case/internal/security"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
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

	cE, err := models.ClientEncounters(c.Context(), db, "client_id="+strconv.Itoa(id))
	if err != nil {
		fmt.Println("something something: ", err.Error())
	}

	st, err := models.Statuses(c.Context(), db, "client_id="+strconv.Itoa(id))
	if err != nil {
		fmt.Println("something something: ", err.Error())
	}

	data.User = userName
	data.Role = role
	data.Optionz = Get_Client_Optionz()
	data.Form = client
	data.FormChild1 = cE
	data.FormChild2 = st

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
	CtxG = c
	dbG = db

	userID, userName := GetUser(c, sl, store)
	role := security.GetRoles(userID, "admin")

	data := NewTemplateData(c, store)
	data.User = userName
	data.Role = role

	fmt.Println("loading case list page")
	//
	facility := GetCurrentFacility(c, db, sl, store)
	scope := GetDBInt("user_right", "function_scope", "", "user_id= "+strconv.Itoa(userID), 0)

	filter := ""
	if scope == 15 {
		filter = ""
	} else {
		if facility > 0 {
			filter = " site = " + strconv.Itoa(facility)
		}
	}

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
	var lab = &models.Lab{}

	encounter.ClientID.Valid = true
	encounter.ClientID.Int64 = 0

	clinical.EncounterID.Valid = true
	clinical.EncounterID.Int64 = 0

	vital.EncounterID.Valid = true
	vital.EncounterID.Int64 = 0

	lab.EncounterID.Valid = true
	lab.EncounterID.Int64 = 0

	if id > 0 {
		client, _ = models.ClientByID(c.Context(), db, id)
	}

	encounter.ClientID.Int64 = int64(client.ID)

	if jd > 0 {
		encounter, _ = models.EncounterByEncounterID(c.Context(), db, jd)
		clinical.EncounterID.Int64 = int64(encounter.EncounterID)
		vital.EncounterID.Int64 = int64(encounter.EncounterID)
		lab.EncounterID.Int64 = int64(encounter.EncounterID)

		cl, er := models.ClinicalByEncounterID(c.Context(), db, jd)
		if er == nil {
			clinical = cl
		}
		vt, er := models.VitalByEncounterID(c.Context(), db, jd)
		if er == nil {
			vital = vt
		}

		l, er := models.LabByEncounterID(c.Context(), db, jd)
		if er == nil {
			lab = l
		}

		// add logic to load child forms with encounter ID
	}

	data := NewTemplateData(c, store)

	data.User = userName
	data.Role = role
	data.Optionz = Get_Client_Optionz()
	data.Form = encounter
	data.FormRef = client

	data.FormChild1 = clinical
	data.FormChild2 = vital
	data.FormChild3 = lab

	fmt.Println("loading case encounter form page")
	return GenerateHTML(c, data, "form_encounters")
}

func HandlerCaseEncounterSubmit(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	id, er := strconv.Atoi(c.FormValue("id"))
	if er != nil {
		id = 0
	}

	userID := GetCurrentUser(c, store)

	//encounter
	encounter := models.Encounter{
		EncounterID:   id,
		EncounterType: ParseNullInt(c.FormValue("encounter_type")),
		EncounterTime: ParseNullString(c.FormValue("encounter_time")),
		ClientID:      ParseNullInt(c.FormValue("cid")),
		EncounterDate: ParseNullString(c.FormValue("encounter_date")),
		ManagedBy:     ParseNullInt(c.FormValue("managed_by")),
	}

	if id == 0 {
		encounter.EnterOn.Valid = true
		encounter.EnterBy.Valid = true

		encounter.EnterBy.Int64 = int64(userID)
		encounter.EnterOn.Time = time.Now()
		err := encounter.Insert(c.Context(), db)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		encounter.SetAsExists()
		err := encounter.Update(c.Context(), db)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	//vital

	vitals_id, er := strconv.Atoi(c.FormValue("vitals_id"))
	if er != nil {
		vitals_id = 0
	}

	vital := models.Vital{
		VitalsID:            vitals_id,
		EncounterID:         sql.NullInt64{Int64: int64(encounter.EncounterID), Valid: true},
		HeartRate:           ParseNullFloat(c.FormValue("heart_rate")),
		BpSystolic:          ParseNullFloat(c.FormValue("bp_systolic")),
		BpDiastolic:         ParseNullFloat(c.FormValue("bp_diastolic")),
		CapillaryRefill:     ParseNullInt(c.FormValue("capillary_refill")),
		RespiratoryRate:     ParseNullFloat(c.FormValue("respiratory_rate")),
		Saturation:          ParseNullFloat(c.FormValue("saturation")),
		Weight:              ParseNullFloat(c.FormValue("weight")),
		Height:              ParseNullFloat(c.FormValue("height")),
		Temperature:         ParseNullFloat(c.FormValue("temperature")),
		LowestConsciousness: ParseNullString(c.FormValue("lowest_consciousness")),
		MentalStatus:        ParseNullString(c.FormValue("mental_status")),
		Muac:                ParseNullFloat(c.FormValue("muac")),
		Bleeding:            ParseNullInt(c.FormValue("bleeding")),
		Shock:               ParseNullInt(c.FormValue("shock")),
		Meningitis:          ParseNullInt(c.FormValue("meningitis")),
		Confusion:           ParseNullInt(c.FormValue("confusion")),
		Seizure:             ParseNullInt(c.FormValue("seizure")),
		Coma:                ParseNullInt(c.FormValue("coma")),
		Bacteraemia:         ParseNullInt(c.FormValue("bacteraemia")),
		Hyperglycemia:       ParseNullInt(c.FormValue("hyperglycemia")),
		Hypoglycemia:        ParseNullInt(c.FormValue("hypoglycemia")),
		Other:               ParseNullString(c.FormValue("other")),
	}

	if vitals_id == 0 {
		err := vital.Insert(c.Context(), db)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		vital.SetAsExists()
		err := vital.Update(c.Context(), db)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	//clinical

	clinical_id, er := strconv.Atoi(c.FormValue("clinical_id"))
	if er != nil {
		clinical_id = 0
	}

	clinical := models.Clinical{
		ClinicalID:             vitals_id,
		EncounterID:            sql.NullInt64{Int64: int64(encounter.EncounterID), Valid: true},
		Fever:                  ParseNullInt(c.FormValue("fever")),
		Fatigue:                ParseNullInt(c.FormValue("fatigue")),
		Weakness:               ParseNullInt(c.FormValue("weakness")),
		Malaise:                ParseNullInt(c.FormValue("malaise")),
		Myalgia:                ParseNullInt(c.FormValue("myalgia")),
		Anorexia:               ParseNullInt(c.FormValue("anorexia")),
		SoreThroat:             ParseNullInt(c.FormValue("sore_throat")),
		Headache:               ParseNullInt(c.FormValue("headache")),
		Nausea:                 ParseNullInt(c.FormValue("nausea")),
		ChestPain:              ParseNullInt(c.FormValue("chest_pain")),
		JointPain:              ParseNullInt(c.FormValue("joint_pain")),
		Hiccups:                ParseNullInt(c.FormValue("hiccups")),
		Cough:                  ParseNullInt(c.FormValue("cough")),
		DifficultyBreathing:    ParseNullInt(c.FormValue("difficulty_breathing")),
		DifficultySwallowing:   ParseNullInt(c.FormValue("difficulty_swallowing")),
		AbdominalPain:          ParseNullInt(c.FormValue("abdominal_pain")),
		Diarrhoea:              ParseNullInt(c.FormValue("diarrhoea")),
		Vomiting:               ParseNullInt(c.FormValue("vomiting")),
		Irritability:           ParseNullInt(c.FormValue("irritability")),
		Dysphagia:              ParseNullInt(c.FormValue("dysphagia")),
		UnusualBleeding:        ParseNullInt(c.FormValue("unusual_bleeding")),
		Dehydration:            ParseNullInt(c.FormValue("dehydration")),
		Shock:                  ParseNullInt(c.FormValue("shock")),
		Anuria:                 ParseNullInt(c.FormValue("anuria")),
		Disorientation:         ParseNullInt(c.FormValue("disorientation")),
		Agitation:              ParseNullInt(c.FormValue("agitation")),
		Seizure:                ParseNullInt(c.FormValue("seizure")),
		Meningitis:             ParseNullInt(c.FormValue("meningitis")),
		Confusion:              ParseNullInt(c.FormValue("confusion")),
		Coma:                   ParseNullInt(c.FormValue("coma")),
		Bacteraemia:            ParseNullInt(c.FormValue("bacteraemia")),
		Hyperglycemia:          ParseNullInt(c.FormValue("hyperglycemia")),
		Hypoglycemia:           ParseNullInt(c.FormValue("hypoglycemia")),
		OtherComplications:     ParseNullInt(c.FormValue("other_complications")),
		AzaComplicationsSpecif: ParseNullString(c.FormValue("aza_complications_specif")),
		PharyngealErythema:     ParseNullInt(c.FormValue("pharyngeal_erythema")),
		PharyngealExudate:      ParseNullInt(c.FormValue("pharyngeal_exudate")),
		ConjunctivalInjection:  ParseNullInt(c.FormValue("conjunctival_injection")),
		OedemaFace:             ParseNullInt(c.FormValue("oedema_face")),
		TenderAbdomen:          ParseNullInt(c.FormValue("tender_abdomen")),
		SunkenEyes:             ParseNullInt(c.FormValue("sunken_eyes")),
		TentingSkin:            ParseNullInt(c.FormValue("tenting_skin")),
		PalpableLiver:          ParseNullInt(c.FormValue("palpable_liver")),
		PalpableSpleen:         ParseNullInt(c.FormValue("palpable_spleen")),
		Jaundice:               ParseNullInt(c.FormValue("jaundice")),
		EnlargedLymphNodes:     ParseNullInt(c.FormValue("enlarged_lymph_nodes")),
		LowerExtremityOedema:   ParseNullInt(c.FormValue("lower_extremity_oedema")),
		Bleeding:               ParseNullInt(c.FormValue("clinical_bleeding")),
		BleedingNose:           ParseNullInt(c.FormValue("bleeding_nose")),
		BleedingMouth:          ParseNullInt(c.FormValue("bleeding_mouth")),
		BleedingVagina:         ParseNullInt(c.FormValue("bleeding_vagina")),
		BleedingRectum:         ParseNullInt(c.FormValue("bleeding_rectum")),
		BleedingSputum:         ParseNullInt(c.FormValue("bleeding_sputum")),
		BleedingUrine:          ParseNullInt(c.FormValue("bleeding_urine")),
		BleedingIvSite:         ParseNullInt(c.FormValue("bleeding_iv_site")),
		BleedingOther:          ParseNullInt(c.FormValue("bleeding_other")),
		BleedingOtherSpecif:    ParseNullString(c.FormValue("bleeding_other_specif")),
	}

	if clinical_id == 0 {
		err := clinical.Insert(c.Context(), db)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		clinical.SetAsExists()
		err := clinical.Update(c.Context(), db)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	lab_id, er := strconv.Atoi(c.FormValue("lab_id"))
	if er != nil {
		clinical_id = 0
	}

	lab := models.Lab{
		LabID:                 vitals_id,
		EncounterID:           sql.NullInt64{Int64: int64(encounter.EncounterID), Valid: true},
		Specimen:              ParseNullInt(c.FormValue("specimen")),
		SampleBlood:           ParseNullInt(c.FormValue("sample_blood")),
		SampleUrine:           ParseNullInt(c.FormValue("sample_urine")),
		SampleSwab:            ParseNullInt(c.FormValue("sample_swab")),
		SampleAza:             ParseNullString(c.FormValue("sample_aza")),
		EbolaRdt:              ParseNullInt(c.FormValue("ebola_rdt")),
		EbolaRdtDate:          ParseNullString(c.FormValue("ebola_rdt_date")),
		EbolaRdtResults:       ParseNullInt(c.FormValue("ebola_rdt_results")),
		EbolaPcr:              ParseNullInt(c.FormValue("ebola_pcr")),
		EbolaPcrAza:           ParseNullString(c.FormValue("ebola_pcr_aza")),
		EbolaPcrDate:          ParseNullString(c.FormValue("ebola_pcr_date")),
		EbolaPcrGp:            ParseNullInt(c.FormValue("ebola_pcr_gp")),
		EbolaPcrGpCt:          ParseNullFloat(c.FormValue("ebola_pcr_gp_ct")),
		EbolaPcrNp:            ParseNullInt(c.FormValue("ebola_pcr_np")),
		EbolaPcrNpCt:          ParseNullFloat(c.FormValue("ebola_pcr_np_ct")),
		EbolaPcrIndeterminate: ParseNullInt(c.FormValue("ebola_pcr_indeterminate")),
		MalariaRdtDate:        ParseNullString(c.FormValue("malaria_rdt_date")),
		MalariaRdtResult:      ParseNullInt(c.FormValue("malaria_rdt_result")),
		BloodCultureDate:      ParseNullString(c.FormValue("blood_culture_date")),
		BloodCultureResult:    ParseNullInt(c.FormValue("blood_culture_result")),
		TestPosInfection:      ParseNullInt(c.FormValue("test_pos_infection")),
		TestPosInfectionAza:   ParseNullString(c.FormValue("test_pos_infection_aza")),
		Haemoglobinuria:       ParseNullInt(c.FormValue("haemoglobinuria")),
		Proteinuria:           ParseNullInt(c.FormValue("proteinuria")),
		Hematuria:             ParseNullInt(c.FormValue("hematuria")),
		BloodGas:              ParseNullInt(c.FormValue("blood_gas")),
		Ph:                    ParseNullFloat(c.FormValue("ph")),
		Pco2:                  ParseNullFloat(c.FormValue("pco2")),
		Pao2:                  ParseNullFloat(c.FormValue("pao2")),
		Hco3:                  ParseNullFloat(c.FormValue("hco3")),
		OxygenTherapy:         ParseNullFloat(c.FormValue("oxygen_therapy")),
		AltSgpt:               ParseNullFloat(c.FormValue("alt_sgpt")),
		AstSgo:                ParseNullFloat(c.FormValue("ast_sgo")),
		Creatinine:            ParseNullFloat(c.FormValue("creatinine")),
		Potassium:             ParseNullFloat(c.FormValue("potassium")),
		Urea:                  ParseNullFloat(c.FormValue("urea")),
		CreatinineKinase:      ParseNullFloat(c.FormValue("creatinine_kinase")),
		Calcium:               ParseNullFloat(c.FormValue("calcium")),
		Sodium:                ParseNullFloat(c.FormValue("sodium")),
		AltSgptNd:             ParseNullInt(c.FormValue("alt_sgpt_nd")),
		AstSgoNd:              ParseNullInt(c.FormValue("ast_sgo_nd")),
		CreatinineNd:          ParseNullInt(c.FormValue("creatinine_nd")),
		PotassiumNd:           ParseNullInt(c.FormValue("potassium_nd")),
		UreaNd:                ParseNullInt(c.FormValue("urea_nd")),
		CreatinineKinaseNd:    ParseNullInt(c.FormValue("creatinine_kinase_nd")),
		CalciumNd:             ParseNullInt(c.FormValue("calcium_nd")),
		SodiumNd:              ParseNullInt(c.FormValue("sodium_nd")),
		Glucose:               ParseNullFloat(c.FormValue("glucose")),
		Lactate:               ParseNullFloat(c.FormValue("lactate")),
		Haemoglobin:           ParseNullFloat(c.FormValue("haemoglobin")),
		TotalBilirubin:        ParseNullFloat(c.FormValue("total_bilirubin")),
		WbcCount:              ParseNullFloat(c.FormValue("wbc_count")),
		Platelets:             ParseNullFloat(c.FormValue("platelets")),
		Pt:                    ParseNullFloat(c.FormValue("pt")),
		Aptt:                  ParseNullFloat(c.FormValue("aptt")),
		GlucoseNd:             ParseNullInt(c.FormValue("glucose_nd")),
		LactateNd:             ParseNullInt(c.FormValue("lactate_nd")),
		HaemoglobinNd:         ParseNullInt(c.FormValue("haemoglobin_nd")),
		TotalBilirubinNd:      ParseNullInt(c.FormValue("total_bilirubin_nd")),
		WbcCountNd:            ParseNullInt(c.FormValue("wbc_count_nd")),
		PlateletsNd:           ParseNullInt(c.FormValue("platelets_nd")),
		PtNd:                  ParseNullInt(c.FormValue("pt_nd")),
		ApttNd:                ParseNullInt(c.FormValue("aptt_nd")),
	}

	if lab_id == 0 {
		err := lab.Insert(c.Context(), db)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		lab.SetAsExists()
		err := lab.Update(c.Context(), db)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// re route

	urlx := "/cases/encounters/new/" + strconv.Itoa(int(encounter.ClientID.Int64)) + "/" + strconv.Itoa(encounter.EncounterID)

	return c.Redirect(urlx)
}

func HandlerAPIGetEncounter(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	// Get ID from the query parameter

	id := c.Query("id")

	if id == "" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "",
		})
	}

	encounter_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "",
		})
	}

	//encounter := models.Encounter{}
	//encounter, err := models.EncounterByEncounterID(c.Context(), db, encounter_id)
	/* if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "",
		})
	} */

	var clinical = &models.Clinical{}
	var vital = &models.Vital{}

	clinical, _ = models.ClinicalByEncounterID(c.Context(), db, encounter_id)
	vital, _ = models.VitalByEncounterID(c.Context(), db, encounter_id)

	rtnStr := ` Vitals<br />
				<table class="full-width" border="1">
					<tr>
						<td>Weight: ` + fmt.Sprintf("%.2f", vital.Weight.Float64) + `</td>
						<td>Height: ` + fmt.Sprintf("%.2f", vital.Height.Float64) + `</td>
					</tr>
				</table>
				Symptomms<br/>
				<table class="full-width" border="1">
					<tr>
						<td valign="top">
							Fever: ` + strconv.Itoa(int(clinical.Fever.Int64)) + `<br/>
							Fatigue:<br/>
							Weakness:<br/>
							Malaise:<br/>
							Myalgia:<br/>
							Anorexia:<br/>
							Sore throat
						</td>
						<td valign="top">
							Headache:<br/> 
							Nausea:<br/> 
							Chest pain:<br/> 
							Joint Pain:<br/> 
							Hiccups:<br/>
							Cough:<br/>
						</td>
						<td valign="top">
							Chest pain:<br/>
							Difficulty breathing:<br/>
							Difficulty swallowing:<br/> 
							Abdominal pain:<br/> 
							Diarrhoea:<br/>
							Vomiting:<br/>
							Irritability / Confusion:<br/> 
						</td>
					</tr>
				</table>

				<br/>
				Signs<br/>
				<table class="full-width" border="1">
					<tr>
						<td valign="top">
							Pharyngeal erythema:<br/>  
							Pharyngeal exudate:<br/>  
							Conjunctival injection/bleeding:<br/>  
							Oedema of face/neck:<br/> 
							Tender abdomen:<br/> 
							Sunken eyes or fontanelle:<br/>  
							Tenting on skin pinch:<br/>  
							Palpable liver:<br/> 
							Palpable spleen Rash:<br/> 
							Jaundice:<br/> 

						</td>
						<td valign="top">
							Enlarged lymph nodes:<br/>
							Lower extremity oedema :<br/> 
							Bleeding:<br/> 
						</td>
					</tr>
				</table>
				<br/>
				Specimen <br/>
				<table class="full-width" border="1">
					<tr>
						<td valign="top">
						</td>
					</tr>
				</table>
				<br/>
				Lab Results <br/>
				<table class="full-width" border="1">
					<tr>
						<td valign="top">
						</td>
					</tr>
				</table>
				<br/>
				Medications <br/>
				<table class="full-width" border="1">
					<tr>
						<td valign="top">
						</td>
					</tr>
				</table>`

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": rtnStr,
	})

}

func HandlerAPIGetStatuses(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	userID := GetCurrentUser(c, store)

	// Check if user is logged in
	if userID == 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = "0"
	}

	statuses, er := models.Statusez(c.Context(), db, " client_id = "+clientID)
	if er != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching statuses",
		})
	}

	return c.JSON(statuses)

}

func HandlerAPIPostStatus(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {

	//=================

	userID := GetCurrentUser(c, store)
	fmt.Println("we here")
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

	var s models.Status

	s.StatusID = int(ParseNullInt2(formData["status_id"]).Int64)
	s.ClientID = ParseNullInt2(formData["client_id"])
	s.StatusDate = ParseNullString2(formData["status_date"])
	s.Status = ParseNullString2(formData["status"])
	s.StatusNotes = ParseNullString2(formData["status_notes"])

	s.UpdatedBy.Valid = true
	s.UpdatedBy.Int64 = int64(userID)

	s.UpdatedOn.Valid = true
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02")
	s.UpdatedOn.String = formattedTime

	fmt.Println(s)
	// Check if it's a new record (StatusID == 0)
	if s.StatusID > 0 {
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
