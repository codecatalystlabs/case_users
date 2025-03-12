package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// Lab represents a row from 'public.lab'.
type Lab struct {
	ID                int             `json:"id"`                  // id
	EncounterID       sql.NullInt64   `json:"encounter_id"`        // encounter_id
	LabDate           sql.NullTime    `json:"lab_date"`            // lab_date
	Test              sql.NullInt64   `json:"test"`                // test
	Results           sql.NullFloat64 `json:"results"`             // results
	Done              sql.NullInt64   `json:"done"`                // done
	Alt               sql.NullFloat64 `json:"alt"`                 // alt
	Ast               sql.NullFloat64 `json:"ast"`                 // ast
	Creatinine        sql.NullFloat64 `json:"creatinine"`          // creatinine
	Potassium         sql.NullFloat64 `json:"potassium"`           // potassium
	Urea              sql.NullFloat64 `json:"urea"`                // urea
	CreatinineKinase  sql.NullFloat64 `json:"creatinine_kinase"`   // creatinine_kinase
	Calcium           sql.NullFloat64 `json:"calcium"`             // calcium
	Sodium            sql.NullFloat64 `json:"sodium"`              // sodium
	Crp               sql.NullFloat64 `json:"crp"`                 // crp
	Glucose           sql.NullFloat64 `json:"glucose"`             // glucose
	Lactate           sql.NullFloat64 `json:"lactate"`             // lactate
	Haemoglobin       sql.NullFloat64 `json:"haemoglobin"`         // haemoglobin
	Total             sql.NullFloat64 `json:"total"`               // total
	Wbc               sql.NullFloat64 `json:"wbc"`                 // wbc
	Platelets         sql.NullFloat64 `json:"platelets"`           // platelets
	Prothrombin       sql.NullFloat64 `json:"prothrombin"`         // prothrombin
	Activated         sql.NullFloat64 `json:"activated"`           // activated
	Other             sql.NullFloat64 `json:"other"`               // other
	OtherSpecifyName  sql.NullFloat64 `json:"other_specify_name"`  // other_specify_name
	OtherSpecifyValue sql.NullFloat64 `json:"other_specify_value"` // other_specify_value
	Alt1              sql.NullInt64   `json:"alt_1"`               // alt_1
	Ast1              sql.NullInt64   `json:"ast_1"`               // ast_1
	Creatinine1       sql.NullInt64   `json:"creatinine_1"`        // creatinine_1
	Potassium1        sql.NullInt64   `json:"potassium_1"`         // potassium_1
	Urea1             sql.NullInt64   `json:"urea_1"`              // urea_1
	CreatinineKinase1 sql.NullInt64   `json:"creatinine_kinase_1"` // creatinine_kinase_1
	Calcium1          sql.NullInt64   `json:"calcium_1"`           // calcium_1
	Sodium1           sql.NullInt64   `json:"sodium_1"`            // sodium_1
	Crp1              sql.NullInt64   `json:"crp_1"`               // crp_1
	Glucose1          sql.NullInt64   `json:"glucose_1"`           // glucose_1
	Lactate1          sql.NullInt64   `json:"lactate_1"`           // lactate_1
	Haemoglobin1      sql.NullInt64   `json:"haemoglobin_1"`       // haemoglobin_1
	Total1            sql.NullInt64   `json:"total_1"`             // total_1
	Wbc1              sql.NullInt64   `json:"wbc_1"`               // wbc_1
	Platelets1        sql.NullInt64   `json:"platelets_1"`         // platelets_1
	Prothrombin1      sql.NullInt64   `json:"prothrombin_1"`       // prothrombin_1
	Activated1        sql.NullInt64   `json:"activated_1"`         // activated_1
	Other1            sql.NullInt64   `json:"other_1"`             // other_1
	EnterBy           sql.NullInt64   `json:"enter_by"`            // enter_by
	EnterOn           sql.NullTime    `json:"enter_on"`            // enter_on
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [Lab] exists in the database.
func (l *Lab) Exists() bool {
	return l._exists
}

// Deleted returns true when the [Lab] has been marked for deletion
// from the database.
func (l *Lab) Deleted() bool {
	return l._deleted
}

// Insert inserts the [Lab] to the database.
func (l *Lab) Insert(ctx context.Context, db DB) error {
	switch {
	case l._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case l._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO public.lab (` +
		`encounter_id, lab_date, test, results, done, alt, ast, creatinine, potassium, urea, creatinine_kinase, calcium, sodium, crp, glucose, lactate, haemoglobin, total, wbc, platelets, prothrombin, activated, other, other_specify_name, other_specify_value, alt_1, ast_1, creatinine_1, potassium_1, urea_1, creatinine_kinase_1, calcium_1, sodium_1, crp_1, glucose_1, lactate_1, haemoglobin_1, total_1, wbc_1, platelets_1, prothrombin_1, activated_1, other_1, enter_by, enter_on` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45` +
		`) RETURNING id`
	// run
	logf(sqlstr, l.EncounterID, l.LabDate, l.Test, l.Results, l.Done, l.Alt, l.Ast, l.Creatinine, l.Potassium, l.Urea, l.CreatinineKinase, l.Calcium, l.Sodium, l.Crp, l.Glucose, l.Lactate, l.Haemoglobin, l.Total, l.Wbc, l.Platelets, l.Prothrombin, l.Activated, l.Other, l.OtherSpecifyName, l.OtherSpecifyValue, l.Alt1, l.Ast1, l.Creatinine1, l.Potassium1, l.Urea1, l.CreatinineKinase1, l.Calcium1, l.Sodium1, l.Crp1, l.Glucose1, l.Lactate1, l.Haemoglobin1, l.Total1, l.Wbc1, l.Platelets1, l.Prothrombin1, l.Activated1, l.Other1, l.EnterBy, l.EnterOn)
	if err := db.QueryRowContext(ctx, sqlstr, l.EncounterID, l.LabDate, l.Test, l.Results, l.Done, l.Alt, l.Ast, l.Creatinine, l.Potassium, l.Urea, l.CreatinineKinase, l.Calcium, l.Sodium, l.Crp, l.Glucose, l.Lactate, l.Haemoglobin, l.Total, l.Wbc, l.Platelets, l.Prothrombin, l.Activated, l.Other, l.OtherSpecifyName, l.OtherSpecifyValue, l.Alt1, l.Ast1, l.Creatinine1, l.Potassium1, l.Urea1, l.CreatinineKinase1, l.Calcium1, l.Sodium1, l.Crp1, l.Glucose1, l.Lactate1, l.Haemoglobin1, l.Total1, l.Wbc1, l.Platelets1, l.Prothrombin1, l.Activated1, l.Other1, l.EnterBy, l.EnterOn).Scan(&l.ID); err != nil {
		return logerror(err)
	}
	// set exists
	l._exists = true
	return nil
}

// Update updates a [Lab] in the database.
func (l *Lab) Update(ctx context.Context, db DB) error {
	switch {
	case !l._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case l._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with composite primary key
	const sqlstr = `UPDATE public.lab SET ` +
		`encounter_id = $1, lab_date = $2, test = $3, results = $4, done = $5, alt = $6, ast = $7, creatinine = $8, potassium = $9, urea = $10, creatinine_kinase = $11, calcium = $12, sodium = $13, crp = $14, glucose = $15, lactate = $16, haemoglobin = $17, total = $18, wbc = $19, platelets = $20, prothrombin = $21, activated = $22, other = $23, other_specify_name = $24, other_specify_value = $25, alt_1 = $26, ast_1 = $27, creatinine_1 = $28, potassium_1 = $29, urea_1 = $30, creatinine_kinase_1 = $31, calcium_1 = $32, sodium_1 = $33, crp_1 = $34, glucose_1 = $35, lactate_1 = $36, haemoglobin_1 = $37, total_1 = $38, wbc_1 = $39, platelets_1 = $40, prothrombin_1 = $41, activated_1 = $42, other_1 = $43, enter_by = $44, enter_on = $45 ` +
		`WHERE id = $46`
	// run
	logf(sqlstr, l.EncounterID, l.LabDate, l.Test, l.Results, l.Done, l.Alt, l.Ast, l.Creatinine, l.Potassium, l.Urea, l.CreatinineKinase, l.Calcium, l.Sodium, l.Crp, l.Glucose, l.Lactate, l.Haemoglobin, l.Total, l.Wbc, l.Platelets, l.Prothrombin, l.Activated, l.Other, l.OtherSpecifyName, l.OtherSpecifyValue, l.Alt1, l.Ast1, l.Creatinine1, l.Potassium1, l.Urea1, l.CreatinineKinase1, l.Calcium1, l.Sodium1, l.Crp1, l.Glucose1, l.Lactate1, l.Haemoglobin1, l.Total1, l.Wbc1, l.Platelets1, l.Prothrombin1, l.Activated1, l.Other1, l.EnterBy, l.EnterOn, l.ID)
	if _, err := db.ExecContext(ctx, sqlstr, l.EncounterID, l.LabDate, l.Test, l.Results, l.Done, l.Alt, l.Ast, l.Creatinine, l.Potassium, l.Urea, l.CreatinineKinase, l.Calcium, l.Sodium, l.Crp, l.Glucose, l.Lactate, l.Haemoglobin, l.Total, l.Wbc, l.Platelets, l.Prothrombin, l.Activated, l.Other, l.OtherSpecifyName, l.OtherSpecifyValue, l.Alt1, l.Ast1, l.Creatinine1, l.Potassium1, l.Urea1, l.CreatinineKinase1, l.Calcium1, l.Sodium1, l.Crp1, l.Glucose1, l.Lactate1, l.Haemoglobin1, l.Total1, l.Wbc1, l.Platelets1, l.Prothrombin1, l.Activated1, l.Other1, l.EnterBy, l.EnterOn, l.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [Lab] to the database.
func (l *Lab) Save(ctx context.Context, db DB) error {
	if l.Exists() {
		return l.Update(ctx, db)
	}
	return l.Insert(ctx, db)
}

// Upsert performs an upsert for [Lab].
func (l *Lab) Upsert(ctx context.Context, db DB) error {
	switch {
	case l._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO public.lab (` +
		`id, encounter_id, lab_date, test, results, done, alt, ast, creatinine, potassium, urea, creatinine_kinase, calcium, sodium, crp, glucose, lactate, haemoglobin, total, wbc, platelets, prothrombin, activated, other, other_specify_name, other_specify_value, alt_1, ast_1, creatinine_1, potassium_1, urea_1, creatinine_kinase_1, calcium_1, sodium_1, crp_1, glucose_1, lactate_1, haemoglobin_1, total_1, wbc_1, platelets_1, prothrombin_1, activated_1, other_1, enter_by, enter_on` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46` +
		`)` +
		` ON CONFLICT (id) DO ` +
		`UPDATE SET ` +
		`encounter_id = EXCLUDED.encounter_id, lab_date = EXCLUDED.lab_date, test = EXCLUDED.test, results = EXCLUDED.results, done = EXCLUDED.done, alt = EXCLUDED.alt, ast = EXCLUDED.ast, creatinine = EXCLUDED.creatinine, potassium = EXCLUDED.potassium, urea = EXCLUDED.urea, creatinine_kinase = EXCLUDED.creatinine_kinase, calcium = EXCLUDED.calcium, sodium = EXCLUDED.sodium, crp = EXCLUDED.crp, glucose = EXCLUDED.glucose, lactate = EXCLUDED.lactate, haemoglobin = EXCLUDED.haemoglobin, total = EXCLUDED.total, wbc = EXCLUDED.wbc, platelets = EXCLUDED.platelets, prothrombin = EXCLUDED.prothrombin, activated = EXCLUDED.activated, other = EXCLUDED.other, other_specify_name = EXCLUDED.other_specify_name, other_specify_value = EXCLUDED.other_specify_value, alt_1 = EXCLUDED.alt_1, ast_1 = EXCLUDED.ast_1, creatinine_1 = EXCLUDED.creatinine_1, potassium_1 = EXCLUDED.potassium_1, urea_1 = EXCLUDED.urea_1, creatinine_kinase_1 = EXCLUDED.creatinine_kinase_1, calcium_1 = EXCLUDED.calcium_1, sodium_1 = EXCLUDED.sodium_1, crp_1 = EXCLUDED.crp_1, glucose_1 = EXCLUDED.glucose_1, lactate_1 = EXCLUDED.lactate_1, haemoglobin_1 = EXCLUDED.haemoglobin_1, total_1 = EXCLUDED.total_1, wbc_1 = EXCLUDED.wbc_1, platelets_1 = EXCLUDED.platelets_1, prothrombin_1 = EXCLUDED.prothrombin_1, activated_1 = EXCLUDED.activated_1, other_1 = EXCLUDED.other_1, enter_by = EXCLUDED.enter_by, enter_on = EXCLUDED.enter_on `
	// run
	logf(sqlstr, l.ID, l.EncounterID, l.LabDate, l.Test, l.Results, l.Done, l.Alt, l.Ast, l.Creatinine, l.Potassium, l.Urea, l.CreatinineKinase, l.Calcium, l.Sodium, l.Crp, l.Glucose, l.Lactate, l.Haemoglobin, l.Total, l.Wbc, l.Platelets, l.Prothrombin, l.Activated, l.Other, l.OtherSpecifyName, l.OtherSpecifyValue, l.Alt1, l.Ast1, l.Creatinine1, l.Potassium1, l.Urea1, l.CreatinineKinase1, l.Calcium1, l.Sodium1, l.Crp1, l.Glucose1, l.Lactate1, l.Haemoglobin1, l.Total1, l.Wbc1, l.Platelets1, l.Prothrombin1, l.Activated1, l.Other1, l.EnterBy, l.EnterOn)
	if _, err := db.ExecContext(ctx, sqlstr, l.ID, l.EncounterID, l.LabDate, l.Test, l.Results, l.Done, l.Alt, l.Ast, l.Creatinine, l.Potassium, l.Urea, l.CreatinineKinase, l.Calcium, l.Sodium, l.Crp, l.Glucose, l.Lactate, l.Haemoglobin, l.Total, l.Wbc, l.Platelets, l.Prothrombin, l.Activated, l.Other, l.OtherSpecifyName, l.OtherSpecifyValue, l.Alt1, l.Ast1, l.Creatinine1, l.Potassium1, l.Urea1, l.CreatinineKinase1, l.Calcium1, l.Sodium1, l.Crp1, l.Glucose1, l.Lactate1, l.Haemoglobin1, l.Total1, l.Wbc1, l.Platelets1, l.Prothrombin1, l.Activated1, l.Other1, l.EnterBy, l.EnterOn); err != nil {
		return logerror(err)
	}
	// set exists
	l._exists = true
	return nil
}

// Delete deletes the [Lab] from the database.
func (l *Lab) Delete(ctx context.Context, db DB) error {
	switch {
	case !l._exists: // doesn't exist
		return nil
	case l._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM public.lab ` +
		`WHERE id = $1`
	// run
	logf(sqlstr, l.ID)
	if _, err := db.ExecContext(ctx, sqlstr, l.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	l._deleted = true
	return nil
}

// LabByID retrieves a row from 'public.lab' as a [Lab].
//
// Generated from index 'lab_pkey'.
func LabByID(ctx context.Context, db DB, id int) (*Lab, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, encounter_id, lab_date, test, results, done, alt, ast, creatinine, potassium, urea, creatinine_kinase, calcium, sodium, crp, glucose, lactate, haemoglobin, total, wbc, platelets, prothrombin, activated, other, other_specify_name, other_specify_value, alt_1, ast_1, creatinine_1, potassium_1, urea_1, creatinine_kinase_1, calcium_1, sodium_1, crp_1, glucose_1, lactate_1, haemoglobin_1, total_1, wbc_1, platelets_1, prothrombin_1, activated_1, other_1, enter_by, enter_on ` +
		`FROM public.lab ` +
		`WHERE id = $1`
	// run
	logf(sqlstr, id)
	l := Lab{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, id).Scan(&l.ID, &l.EncounterID, &l.LabDate, &l.Test, &l.Results, &l.Done, &l.Alt, &l.Ast, &l.Creatinine, &l.Potassium, &l.Urea, &l.CreatinineKinase, &l.Calcium, &l.Sodium, &l.Crp, &l.Glucose, &l.Lactate, &l.Haemoglobin, &l.Total, &l.Wbc, &l.Platelets, &l.Prothrombin, &l.Activated, &l.Other, &l.OtherSpecifyName, &l.OtherSpecifyValue, &l.Alt1, &l.Ast1, &l.Creatinine1, &l.Potassium1, &l.Urea1, &l.CreatinineKinase1, &l.Calcium1, &l.Sodium1, &l.Crp1, &l.Glucose1, &l.Lactate1, &l.Haemoglobin1, &l.Total1, &l.Wbc1, &l.Platelets1, &l.Prothrombin1, &l.Activated1, &l.Other1, &l.EnterBy, &l.EnterOn); err != nil {
		return nil, logerror(err)
	}
	return &l, nil
}
