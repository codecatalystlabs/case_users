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

type FormData struct {
	FidID   []string `form:"input[fid_id][]"`
	MetaID  []string `form:"input[meta_id][]"`
	Scope   []string `form:"input[scope][]"`
	View    []string `form:"input[view][]"`
	Add     []string `form:"input[add][]"`
	Edit    []string `form:"input[edit][]"`
	Execute []string `form:"input[execute][]"`
}

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

	} else {
		id = 0
	}
	fmt.Println("Creating")
	// Correct struct definition with semicolons
	type funclist struct {
		FID      sql.NullInt64 `json:"fid"`
		MetaID   sql.NullInt64 `json:"meta_id"`
		MetaName string        `json:"meta_name"`
		FScope   sql.NullInt64 `json:"f_scope"`
		FView    sql.NullInt64 `json:"f_view"`
		FCreate  sql.NullInt64 `json:"f_create"`
		FEdit    sql.NullInt64 `json:"f_edit"`
		FRemove  sql.NullInt64 `json:"f_remove"`
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
		ON m.meta_id = ur.user_rights_function AND ur.user_id = $1
	WHERE m.meta_category = 3`

	// Execute query safely with parameterized input
	rows, err := db.QueryContext(c.Context(), mysql, id)
	if err != nil {
		fmt.Println("Query Error:", err.Error())
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

	id, er := strconv.Atoi(c.FormValue("id"))
	if er != nil {
		id = 0
	}

	user := models.User{
		UserID:       id,
		UserName:     ParseNullString(c.FormValue("user_name")),
		UserEmployee: ParseNullInt(c.FormValue("user_employee")),
	}

	if user.UserID == 0 {

		user.UserPass = sql.NullString{String: models.Encrypt("123456"), Valid: true}
		err := user.Insert(c.Context(), db)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		user.SetAsExists()
		err := user.Update_NoPass(c.Context(), db)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	//================================================================

	mysql := `SELECT m.meta_id, m.meta_name FROM meta m WHERE m.meta_category = 3`

	// Execute query safely with parameterized input
	rows, err := db.QueryContext(c.Context(), mysql)
	if err != nil {
		fmt.Println("Query Error:", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var m_id int64
		var m_nm string
		err := rows.Scan(&m_id, &m_nm)
		if err != nil {
			fmt.Println("Row Scan Error:", err)
			continue
		}

		fid, err := strconv.ParseInt(c.FormValue("input_fid_id_"+m_nm), 10, 64)
		scope, err := strconv.ParseInt(c.FormValue("input_scope_"+m_nm), 10, 64)
		view, err := strconv.ParseInt(c.FormValue("input_view_"+m_nm), 10, 64)
		add, err := strconv.ParseInt(c.FormValue("input_add_"+m_nm), 10, 64)
		edit, err := strconv.ParseInt(c.FormValue("input_edit_"+m_nm), 10, 64)
		exec, err := strconv.ParseInt(c.FormValue("input_execute_"+m_nm), 10, 64)

		//fmt.Println(m_nm, ":", scope, "-", view, "-", add, "-", edit, "-", exec)

		right := models.UserRight{}

		right.UserID.Valid = true
		right.UserRightsFunction.Valid = true
		right.FunctionScope.Valid = true
		right.UserRightsCanView.Valid = true
		right.UserRightsCanCreate.Valid = true
		right.UserRightsCanEdit.Valid = true
		right.UserRightsCanRemove.Valid = true

		right.UserID.Int64 = int64(user.UserID)
		right.UserRightsFunction.Int64 = m_id
		right.FunctionScope.Int64 = scope
		right.UserRightsCanView.Int64 = view
		right.UserRightsCanCreate.Int64 = add
		right.UserRightsCanEdit.Int64 = edit
		right.UserRightsCanRemove.Int64 = exec

		if fid > 0 {
			right.UserRightsID = int(fid)
			right.SetAsExists()
			er := right.Update(c.Context(), db)
			if er != nil {
				fmt.Println(err.Error())
			}
		} else {
			er := right.Insert(c.Context(), db)
			if er != nil {
				fmt.Println(err.Error())
			}
		}

	}

	urlx := "/users/new/" + strconv.Itoa(user.UserID)

	return c.Redirect(urlx)
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
