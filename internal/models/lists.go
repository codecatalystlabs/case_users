package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

func (c *Client) SetAsExists() {
	c._exists = true
}

func (c *User) SetAsExists() {
	c._exists = true
}

func (c *UserRight) SetAsExists() {
	c._exists = true
}

func (c *Encounter) SetAsExists() {
	c._exists = true
}

func (c *Clinical) SetAsExists() {
	c._exists = true
}

func (c *Vital) SetAsExists() {
	c._exists = true
}

type ClientEncounter struct {
	EncounterID   int
	EncounterType string
	EmployeeFname string
	EmployeeLname string
	EncounterDate string
	EncounterTime string
	ClientID      int
}

func ClientEncounters(ctx context.Context, db DB, flt string) ([]ClientEncounter, error) {
	// query
	sqlstr := ` SELECT 
					encounter_id, meta_name,employee_fname, employee_lname, encounter_date, encounter_time, client_id
				FROM encounter 
				LEFT JOIN meta ON meta.meta_id = encounter.encounter_type
				LEFT JOIN employee on employee.employee_id = encounter.managed_by `
	var args []interface{}
	if flt != "" {
		sqlstr += " WHERE " + flt
	}

	fmt.Println(sqlstr)

	// Log the query
	logf(sqlstr)

	// Execute query
	rows, err := db.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()

	// Slice to hold clients
	var clientencounters []ClientEncounter

	// Iterate through rows
	for rows.Next() {
		var e ClientEncounter
		if err := rows.Scan(
			&e.EncounterID, &e.EncounterType, &e.EmployeeFname, &e.EmployeeLname, &e.EncounterDate, &e.EncounterTime, &e.ClientID,
		); err != nil {
			return nil, logerror(err)
		}
		clientencounters = append(clientencounters, e)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}

	return clientencounters, nil

}

func ClinicalByEncounterID(ctx context.Context, db DB, encounterID int) (*Clinical, error) {
	// query
	const sqlstr = `SELECT ` +
		`clinical_id, encounter_id, fever, fatigue, weakness, malaise, myalgia, anorexia, sore_throat, headache, nausea, chest_pain, joint_pain, hiccups, cough, difficulty_breathing, difficulty_swallowing, abdominal_pain, diarrhoea, vomiting, irritability, dysphagia, unusual_bleeding, dehydration, shock, anuria, disorientation, agitation, seizure, meningitis, confusion, coma, bacteraemia, hyperglycemia, hypoglycemia, other_complications, aza_complications_specif, pharyngeal_erythema, pharyngeal_exudate, conjunctival_injection, oedema_face, tender_abdomen, sunken_eyes, tenting_skin, palpable_liver, palpable_spleen, jaundice, enlarged_lymph_nodes, lower_extremity_oedema, bleeding, bleeding_nose, bleeding_mouth, bleeding_vagina, bleeding_rectum, bleeding_sputum, bleeding_urine, bleeding_iv_site, bleeding_other, bleeding_other_specif ` +
		`FROM public.clinical ` +
		`WHERE encounter_id = $1`
	// run
	logf(sqlstr, encounterID)
	c := Clinical{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, encounterID).Scan(&c.ClinicalID, &c.EncounterID, &c.Fever, &c.Fatigue, &c.Weakness, &c.Malaise, &c.Myalgia, &c.Anorexia, &c.SoreThroat, &c.Headache, &c.Nausea, &c.ChestPain, &c.JointPain, &c.Hiccups, &c.Cough, &c.DifficultyBreathing, &c.DifficultySwallowing, &c.AbdominalPain, &c.Diarrhoea, &c.Vomiting, &c.Irritability, &c.Dysphagia, &c.UnusualBleeding, &c.Dehydration, &c.Shock, &c.Anuria, &c.Disorientation, &c.Agitation, &c.Seizure, &c.Meningitis, &c.Confusion, &c.Coma, &c.Bacteraemia, &c.Hyperglycemia, &c.Hypoglycemia, &c.OtherComplications, &c.AzaComplicationsSpecif, &c.PharyngealErythema, &c.PharyngealExudate, &c.ConjunctivalInjection, &c.OedemaFace, &c.TenderAbdomen, &c.SunkenEyes, &c.TentingSkin, &c.PalpableLiver, &c.PalpableSpleen, &c.Jaundice, &c.EnlargedLymphNodes, &c.LowerExtremityOedema, &c.Bleeding, &c.BleedingNose, &c.BleedingMouth, &c.BleedingVagina, &c.BleedingRectum, &c.BleedingSputum, &c.BleedingUrine, &c.BleedingIvSite, &c.BleedingOther, &c.BleedingOtherSpecif); err != nil {
		return nil, logerror(err)
	}
	return &c, nil
}

func VitalByEncounterID(ctx context.Context, db DB, encounterID int) (*Vital, error) {
	// query
	const sqlstr = `SELECT ` +
		`vitals_id, encounter_id, heart_rate, bp_systolic, bp_diastolic, capillary_refill, respiratory_rate, saturation, weight, height, temperature, lowest_consciousness, mental_status, muac, bleeding, shock, meningitis, confusion, seizure, coma, bacteraemia, hyperglycemia, hypoglycemia, other ` +
		`FROM public.vitals ` +
		`WHERE encounter_id = $1`
	// run
	logf(sqlstr, encounterID)
	v := Vital{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, encounterID).Scan(&v.VitalsID, &v.EncounterID, &v.HeartRate, &v.BpSystolic, &v.BpDiastolic, &v.CapillaryRefill, &v.RespiratoryRate, &v.Saturation, &v.Weight, &v.Height, &v.Temperature, &v.LowestConsciousness, &v.MentalStatus, &v.Muac, &v.Bleeding, &v.Shock, &v.Meningitis, &v.Confusion, &v.Seizure, &v.Coma, &v.Bacteraemia, &v.Hyperglycemia, &v.Hypoglycemia, &v.Other); err != nil {
		return nil, logerror(err)
	}
	return &v, nil
}

func (u *User) Update_NoPass(ctx context.Context, db DB) error {

	// update with composite primary key
	const sqlstr = `UPDATE public.users SET ` +
		`user_name = $1, user_employee = $2 ` +
		`WHERE user_id = $3`
	// run
	logf(sqlstr, u.UserName, u.UserEmployee, u.UserID)
	if _, err := db.ExecContext(ctx, sqlstr, u.UserName, u.UserEmployee, u.UserID); err != nil {
		return logerror(err)
	}
	return nil
}

func (u *User) Update_Pass(ctx context.Context, db DB) error {

	// update with composite primary key
	const sqlstr = `UPDATE public.users SET ` +
		`user_pass = $1 ` +
		`WHERE user_id = $2`
	// run
	logf(sqlstr, u.UserPass, u.UserID)
	if _, err := db.ExecContext(ctx, sqlstr, u.UserPass, u.UserID); err != nil {
		return logerror(err)
	}
	return nil
}

func GetFields(ctx context.Context, db DB, sql_statement string) (map[int][]string, error) {
	var args []interface{}
	// Log the query
	logf(sql_statement)

	// Execute query
	rows, err := db.QueryContext(ctx, sql_statement, args...)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()

	zaResults := make(map[int][]string)
	var i, id int
	var labs string

	for rows.Next() {
		if err := rows.Scan(&id, &labs); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		// Append to the map
		zaResults[i] = []string{strconv.Itoa(id), labs}
		i++
	}

	return zaResults, nil
}

func Clients(ctx context.Context, db DB, flt string) ([]Client, error) {
	// Base SQL query
	sqlstr := `SELECT 
		id, uuid, firstname, lastname, othername, gender, date_of_birth, age, marital, nin, nationality, adm_date, adm_from, lab_no, cif_no, etu_no, case_no, occupation, occupation_aza, date_symptom_onset, date_isolation, pregnant, adm_ward, tb, asplenia, hep, diabetes, hiv, liver, malignancy, heart, pulmonary, kidney, neurologic, other, status, enter_on, enter_by, edit_on, edit_by, transfer, site 
	FROM public.clients`

	// Add filter condition if `flt` is not empty
	var args []interface{}
	if flt != "" {
		sqlstr += " WHERE " + flt
	}

	// Log the query
	logf(sqlstr)

	// Execute query
	rows, err := db.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()

	// Slice to hold clients
	var clients []Client

	// Iterate through rows
	for rows.Next() {
		var c Client
		c._exists = true
		if err := rows.Scan(
			&c.ID, &c.UUID, &c.Firstname, &c.Lastname, &c.Othername, &c.Gender, &c.DateOfBirth, &c.Age, &c.Marital, &c.Nin, &c.Nationality, &c.AdmDate, &c.AdmFrom, &c.LabNo, &c.CifNo, &c.EtuNo, &c.CaseNo, &c.Occupation, &c.OccupationAza, &c.DateSymptomOnset, &c.DateIsolation, &c.Pregnant, &c.AdmWard, &c.Tb, &c.Asplenia, &c.Hep, &c.Diabetes, &c.Hiv, &c.Liver, &c.Malignancy, &c.Heart, &c.Pulmonary, &c.Kidney, &c.Neurologic, &c.Other, &c.Status, &c.EnterOn, &c.EnterBy, &c.EditOn, &c.EditBy, &c.Transfer, &c.Site,
		); err != nil {
			return nil, logerror(err)
		}
		clients = append(clients, c)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}

	return clients, nil
}

func Users(ctx context.Context, db DB, flt string) ([]User, error) {
	// Base SQL query
	sqlstr := `SELECT user_id, user_name, user_pass, user_employee FROM public.users`

	// Add filter condition if `flt` is not empty
	var args []interface{}
	if flt != "" {
		sqlstr += " WHERE " + flt
	}

	// Log the query
	logf(sqlstr)

	// Execute query
	rows, err := db.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()

	// Slice to hold clients
	var users []User

	// Iterate through rows
	for rows.Next() {
		var u User
		u._exists = true
		if err := rows.Scan(
			&u.UserID, &u.UserName, &u.UserPass, &u.UserEmployee,
		); err != nil {
			return nil, logerror(err)
		}

		users = append(users, u)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}

	return users, nil
}

type Metumx struct {
	MetaID          int            `json:"meta_id"`          // meta_id
	MetaCategory    sql.NullInt64  `json:"meta_category"`    // meta_category
	MetaName        sql.NullString `json:"meta_name"`        // meta_name
	MetaOrder       sql.NullInt64  `json:"meta_order"`       // meta_order
	MetaDescription sql.NullString `json:"meta_description"` // meta_description
	MetaLink        sql.NullString `json:"meta_link"`        // meta_link
	Scope           sql.NullInt64  `json:"function_scope"`   // meta_order
	UserID          sql.NullInt64  `json:"user_id"`          // meta_order
	// xo fields
	_exists, _deleted bool
}

func Metums(ctx context.Context, db DB, flt string) ([]Metumx, error) {
	// Base SQL query
	sqlstr := `SELECT meta_id, meta_category, meta_name, meta_order, meta_description, meta_link, f.function_scope, f.user_id 
	           FROM public.meta,
				( 
					Select user_id,function_scope , user_rights_function, 
						user_rights_can_create+user_rights_can_view+user_rights_can_edit+user_rights_can_remove as func
					From public.user_right
					WHERE user_rights_can_create+user_rights_can_view+user_rights_can_edit+user_rights_can_remove > 0 
				) f
				Where f.user_rights_function = meta_id`

	// Add filter condition if `flt` is not empty
	var args []interface{}
	if flt != "" {
		sqlstr += " AND " + flt
	}

	// Log the query
	logf(sqlstr)

	// Execute query
	rows, err := db.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()

	// Slice to hold clients
	var metums []Metumx

	// Iterate through rows
	for rows.Next() {
		var m Metumx
		m._exists = true
		if err := rows.Scan(
			&m.MetaID, &m.MetaCategory, &m.MetaName, &m.MetaOrder, &m.MetaDescription, &m.MetaLink, &m.Scope, &m.UserID,
		); err != nil {
			return nil, logerror(err)
		}

		metums = append(metums, m)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}

	return metums, nil
}
