package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// Outcome represents a row from 'public.outcome'.
type Outcome struct {
	ID           int            `json:"id"`            // id
	EnrollmentID sql.NullInt64  `json:"enrollment_id"` // enrollment_id
	MovementDate sql.NullTime   `json:"movement_date"` // movement_date
	MovementType sql.NullInt64  `json:"movement_type"` // movement_type
	Note         sql.NullString `json:"note"`          // note
	EnterBy      sql.NullInt64  `json:"enter_by"`      // enter_by
	EnterOn      sql.NullTime   `json:"enter_on"`      // enter_on
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [Outcome] exists in the database.
func (o *Outcome) Exists() bool {
	return o._exists
}

// Deleted returns true when the [Outcome] has been marked for deletion
// from the database.
func (o *Outcome) Deleted() bool {
	return o._deleted
}

// Insert inserts the [Outcome] to the database.
func (o *Outcome) Insert(ctx context.Context, db DB) error {
	switch {
	case o._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case o._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO public.outcome (` +
		`enrollment_id, movement_date, movement_type, note, enter_by, enter_on` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6` +
		`) RETURNING id`
	// run
	logf(sqlstr, o.EnrollmentID, o.MovementDate, o.MovementType, o.Note, o.EnterBy, o.EnterOn)
	if err := db.QueryRowContext(ctx, sqlstr, o.EnrollmentID, o.MovementDate, o.MovementType, o.Note, o.EnterBy, o.EnterOn).Scan(&o.ID); err != nil {
		return logerror(err)
	}
	// set exists
	o._exists = true
	return nil
}

// Update updates a [Outcome] in the database.
func (o *Outcome) Update(ctx context.Context, db DB) error {
	switch {
	case !o._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case o._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with composite primary key
	const sqlstr = `UPDATE public.outcome SET ` +
		`enrollment_id = $1, movement_date = $2, movement_type = $3, note = $4, enter_by = $5, enter_on = $6 ` +
		`WHERE id = $7`
	// run
	logf(sqlstr, o.EnrollmentID, o.MovementDate, o.MovementType, o.Note, o.EnterBy, o.EnterOn, o.ID)
	if _, err := db.ExecContext(ctx, sqlstr, o.EnrollmentID, o.MovementDate, o.MovementType, o.Note, o.EnterBy, o.EnterOn, o.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [Outcome] to the database.
func (o *Outcome) Save(ctx context.Context, db DB) error {
	if o.Exists() {
		return o.Update(ctx, db)
	}
	return o.Insert(ctx, db)
}

// Upsert performs an upsert for [Outcome].
func (o *Outcome) Upsert(ctx context.Context, db DB) error {
	switch {
	case o._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO public.outcome (` +
		`id, enrollment_id, movement_date, movement_type, note, enter_by, enter_on` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`)` +
		` ON CONFLICT (id) DO ` +
		`UPDATE SET ` +
		`enrollment_id = EXCLUDED.enrollment_id, movement_date = EXCLUDED.movement_date, movement_type = EXCLUDED.movement_type, note = EXCLUDED.note, enter_by = EXCLUDED.enter_by, enter_on = EXCLUDED.enter_on `
	// run
	logf(sqlstr, o.ID, o.EnrollmentID, o.MovementDate, o.MovementType, o.Note, o.EnterBy, o.EnterOn)
	if _, err := db.ExecContext(ctx, sqlstr, o.ID, o.EnrollmentID, o.MovementDate, o.MovementType, o.Note, o.EnterBy, o.EnterOn); err != nil {
		return logerror(err)
	}
	// set exists
	o._exists = true
	return nil
}

// Delete deletes the [Outcome] from the database.
func (o *Outcome) Delete(ctx context.Context, db DB) error {
	switch {
	case !o._exists: // doesn't exist
		return nil
	case o._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM public.outcome ` +
		`WHERE id = $1`
	// run
	logf(sqlstr, o.ID)
	if _, err := db.ExecContext(ctx, sqlstr, o.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	o._deleted = true
	return nil
}

// OutcomeByID retrieves a row from 'public.outcome' as a [Outcome].
//
// Generated from index 'outcome_pkey'.
func OutcomeByID(ctx context.Context, db DB, id int) (*Outcome, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, enrollment_id, movement_date, movement_type, note, enter_by, enter_on ` +
		`FROM public.outcome ` +
		`WHERE id = $1`
	// run
	logf(sqlstr, id)
	o := Outcome{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, id).Scan(&o.ID, &o.EnrollmentID, &o.MovementDate, &o.MovementType, &o.Note, &o.EnterBy, &o.EnterOn); err != nil {
		return nil, logerror(err)
	}
	return &o, nil
}
