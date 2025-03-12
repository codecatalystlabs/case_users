package models

import (
	"context"
	"log"
	"strconv"
)

func (c *Client) SetAsExists() {
	c._exists = true
}

func (c *User) SetAsExists() {
	c._exists = true
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
