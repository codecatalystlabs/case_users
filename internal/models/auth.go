package models

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// Add a new ErrInvalidCredentials error. We'll use this later if a user
	// tries to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// Add a new ErrDuplicateEmail error. We'll use this later if a user
	// tries to signup with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

func Authenticate(ctx context.Context, db DB, email, password string) (id int, err error) {
	var hashedPassword string

	flag, err := ConfirmAvailabilityOfUser(ctx, db)

	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}

	if flag == 0 {
		stmt := "SELECT user_id, user_pass FROM public.users WHERE user_name = $1 "
		err = db.QueryRowContext(ctx, stmt, email).Scan(&id, &hashedPassword)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				fmt.Println("No data: ")
				return 0, ErrInvalidCredentials
			}
			fmt.Println(err.Error())
			return 0, err
		}

		if Encrypt(password) != hashedPassword {
			fmt.Println("fake password: ")
			return 0, ErrInvalidCredentials
		}

		return id, nil
	}

	return 0, nil
}

func ConfirmAvailabilityOfUser(ctx context.Context, db DB) (int, error) {
	var count int

	stml := "SELECT COUNT(user_id) AS C FROM users LIMIT 1"
	rows, err := db.QueryContext(ctx, stml)

	if err != nil {
		return 1, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&count)

		if err != nil {
			fmt.Println("jj: " + err.Error())
			return 1, err
		}

		if count == 0 {

			employee := Employee{}

			employee.EmployeeFname = sql.NullString{String: "Paul", Valid: true}
			employee.EmployeeLname = sql.NullString{String: "Mbaka", Valid: true}

			var emp int64

			err := employee.Insert(ctx, db)

			if err != nil {
				fmt.Println("j1: " + err.Error())
				emp = 0
			} else {
				emp = int64(employee.EmployeeID)
			}

			user := User{}
			user.UserEmployee = sql.NullInt64{Int64: emp, Valid: true}
			user.UserName = sql.NullString{String: "paul.mbaka@gmail.com", Valid: true}
			user.UserPass = sql.NullString{String: Encrypt("123456789=p"), Valid: true}

			err = user.Insert(ctx, db)
			if err != nil {
				fmt.Println("j2: ")
				return 1, err
			}
			fmt.Println("j3: ")
			return 0, nil
		}

		return 0, nil
	}

	return 0, nil
}

func CreateUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		//Danger(err, "Cannot create uuid")
		fmt.Println(err.Error())
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// hash plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
