package handlers

import (
	"case/internal/models"
	"case/internal/security"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func HandlerUserForm(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {

	userID, userName := GetUser(c, sl, store)
	role := security.GetRoles(userID, "admin")

	CtxG = c
	dbG = db

	//id := c.Params("i") // Retrieve the "id" parameter
	id, err := strconv.Atoi(c.Params("i"))
	if err != nil {
		fmt.Println(err.Error())
	}
	//uzer := models.User{}
	var uzer models.User

	uzer.UserEmployee.Valid = true
	uzer.UserEmployee.Int64 = 0

	data := NewTemplateData(c, store)

	if id > 0 {
		u, er := models.UserByUserID(c.Context(), db, id)
		if er == nil {
			uzer = *u
		}

	}

	// Correct struct definition with semicolons
	type funclist struct {
		FID      int    `json:"fid"`
		MetaID   int    `json:"meta_id"`
		MetaName string `json:"meta_name"`
		FScope   int    `json:"f_scope"`
		FView    int    `json:"f_view"`
		FCreate  int    `json:"f_create"`
		FEdit    int    `json:"f_edit"`
		FRemove  int    `json:"f_remove"`
	}

	// Use parameterized query to prevent SQL injection
	mysql := `  
	SELECT 
		ur.user_rights_id, m.meta_id, m.meta_name, function_scope,
		COALESCE(ur.user_rights_can_view, 0), 
		COALESCE(ur.user_rights_can_create, 0), 
		COALESCE(ur.user_rights_can_edit, 0), 
		COALESCE(ur.user_rights_can_remove, 0)
	FROM meta m 
	LEFT JOIN public.user_right ur 
		ON m.meta_id = ur.user_rights_function AND ur.user_id = ?
	WHERE m.meta_category = 3`

	// Execute query safely with parameterized input
	rows, err := db.QueryContext(c.Context(), mysql, id)
	if err != nil {
		fmt.Println("Query Error:", err)
	}
	defer rows.Close()

	// Slice to store results
	var functions []funclist

	// Iterate over query results
	for rows.Next() {
		var f funclist
		err := rows.Scan(
			&f.FID, &f.MetaID, &f.MetaName, &f.FScope,
			&f.FView, &f.FCreate, &f.FEdit, &f.FRemove,
		)
		if err != nil {
			fmt.Println("Row Scan Error:", err)
			continue
		}
		functions = append(functions, f)
	}

	// Check for errors after looping
	if err = rows.Err(); err != nil {
		fmt.Println("Rows Iteration Error:", err)
	}

	data.User = userName
	data.Role = role
	data.Form = uzer
	data.FormChild1 = functions

	return GenerateHTML(c, data, "form_user")
}
func HandlerUserSubmit(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	return nil
}
func HandlerUserList(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	fmt.Println("starting user list")

	userID, userName := GetUser(c, sl, store)
	role := security.GetRoles(userID, "admin")

	data := NewTemplateData(c, store)
	data.User = userName
	data.Role = role

	fmt.Println("loading user list page")
	//
	mysql := "SELECT user_id, user_name, CONCAT(employee_fname, ' ', employee_lname) AS  lab FROM users, employee Where users.user_employee = employee_id"
	var args []interface{}
	type userlist struct {
		UserID   int
		UserName string
		Lab      string
	}

	// Execute query
	rows, err := db.QueryContext(c.Context(), mysql, args...)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	// Slice to hold clients
	var users []userlist

	// Iterate through rows
	for rows.Next() {
		var u userlist

		if err := rows.Scan(
			&u.UserID, &u.UserName, &u.Lab,
		); err != nil {
			fmt.Println(err.Error())
		}

		users = append(users, u)
	}

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			fmt.Println("error loading user list: ", err.Error())
		} else {
			fmt.Println("error loading user list: ", err.Error())
		}
	}

	data.Form = users

	//
	return GenerateHTML(c, data, "list_users")
}
