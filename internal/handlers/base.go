package handlers

import (
	"bytes"
	"case/internal/models"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gorilla/schema"
)

type Config struct {
	Address      string `json:"Port"`
	ReadTimeout  int64  `json:"ReadTimeout"`
	WriteTimeout int64  `json:"WriteTimeout"`
	Static       string `json:"Static"`
	Ux           string `json:"Ux"`
	Px           string `json:"Px"`
	Dx           string `json:"Dx"`
	LogFile      string `json:"LogFile"`
	LogData      string `json:"LogData"`
	Facility     string `json:"Facility"`
}

type TemplateData struct {
	CurrentYear     int
	User            string
	Role            int
	IsIDPos         bool
	Form            any
	FormRef         any
	FormChild1      any
	FormChild2      any
	FormChild3      any
	FormChild4      any
	FormChild5      any
	Ses             any
	Items           []interface{}
	Optionz         map[string]map[string]string
	Flash           string
	Menuz           string
	IsAuthenticated bool
	CSRFToken       string // Add a CSRFToken field.
}

func NewTemplateData(c *fiber.Ctx, store *session.Store) *TemplateData {
	//log.Printf("template data")
	return &TemplateData{
		CurrentYear:     time.Now().Year(),
		IsAuthenticated: IzAuthenticated(c, store),
		//CSRFToken:       c.Locals("csrf").(string), // Add the CSRF token.
	}
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate":            HumanDate,
	"humanDateTime":        HumanDateTime,
	"IsNullStringEmpty":    IsNullStringEmpty,
	"GetClientOptionLabel": GetClientOptionLabel,
	"seq":                  Seq,
	"GetOptionField":       GetOptionField,
	"GetDBOptions":         GetDBOptions,
	"GetDBLabel":           GetDBLabel,
}

func GetOptionField(table, field, labs, defaultString string, defaultvalue, whole int64) string {
	zaField := ""
	zaDefa1 := ""
	zaDefa2 := ""
	zaDefa3 := ""
	optionz := ""

	if table == "facility" {
		if defaultvalue == 1 {
			zaDefa1 = "selected"
		}

		if defaultvalue == 2 {
			zaDefa2 = "selected"
		}

		optionz = `<option value=""> -- select -- </option>
					<option value="1" ` + zaDefa1 + `>Mulago ETU</option>
					<option value="2" ` + zaDefa2 + `>Mbale ETU</option>`
		zaField = `<select class="form-control-sm patient-input form-select" name="` + field + `" id="` + field + `" aria-label="` + labs + `">
					` + optionz + `
			      </select>`
	}

	if table == "Status" {
		if defaultString == "Suspect" {
			zaDefa1 = "selected"
		}

		if defaultString == "Case" {
			zaDefa2 = "selected"
		}

		if defaultString == "Other" {
			zaDefa3 = "selected"
		}

		optionz = `<option value=""> -- select -- </option>
					<option value="Suspect" ` + zaDefa1 + `>Suspect</option>
					<option value="Case" ` + zaDefa2 + `>Case</option>
					<option value="Other" ` + zaDefa3 + `>Other</option>`
		zaField = `<select class="form-control-sm patient-input form-select" name="` + field + `" id="` + field + `" aria-label="` + labs + `">
					` + optionz + `
			      </select>`
	}

	if table == "pos" {
		if defaultString == "pos" {
			zaDefa1 = "selected"
		}

		if defaultString == "neg" {
			zaDefa2 = "selected"
		}

		if defaultString == "nd" {
			zaDefa3 = "selected"
		}

		optionz = `<option value=""> -- select -- </option>
					<option value="pos" ` + zaDefa1 + `>Pos</option>
					<option value="neg" ` + zaDefa2 + `>Neg</option>
					<option value="nd"  ` + zaDefa3 + `>ND</option>`
		zaField = `<select class="form-control-sm patient-input form-select" name="` + field + `" id="` + field + `" aria-label="` + labs + `">
					` + optionz + `
			      </select>`
	}

	if table == "po" {
		if defaultString == "pos" {
			zaDefa1 = "selected"
		}

		if defaultString == "neg" {
			zaDefa2 = "selected"
		}

		optionz = `<option value=""> -- select -- </option>
					<option value="pos" ` + zaDefa1 + `>Pos</option>
					<option value="neg" ` + zaDefa2 + `>Neg</option>`
		zaField = `<select class="form-control-sm patient-input form-select" name="` + field + `" id="` + field + `" aria-label="` + labs + `">
					` + optionz + `
			      </select>`
	}

	if table == "posx" {
		if defaultString == "pos" {
			zaDefa1 = "selected"
		}

		if defaultString == "neg" {
			zaDefa2 = "selected"
		}

		if defaultString == "indeterminate" {
			zaDefa3 = "selected"
		}

		optionz = `<option value=""> -- select -- </option>
					<option value="pos" ` + zaDefa1 + `>Pos</option>
					<option value="neg" ` + zaDefa2 + `>Neg</option>
					<option value="indeterminate"  ` + zaDefa3 + `>Indeterminate</option>`
		zaField = `<select class="form-control-sm patient-input form-select" name="` + field + `" id="` + field + `" aria-label="` + labs + `">
					` + optionz + `
			      </select>`
	}

	if table == "yn" {
		if defaultvalue == 1 {
			zaDefa1 = "selected"
		}

		if defaultvalue == 2 {
			zaDefa2 = "selected"
		}

		optionz = `<option value=""> -- select -- </option>
					<option value="1" ` + zaDefa1 + `>Yes</option>
					<option value="2" ` + zaDefa2 + `>No</option>`

		zaField = `<select class="form-control-sm patient-input form-select" name="` + field + `" id="` + field + `" aria-label="` + labs + `">
					` + optionz + `
			      </select>`
	}

	if table == "YN" {
		if defaultvalue == 1 {
			zaDefa1 = "selected"
		}

		if defaultvalue == 2 {
			zaDefa2 = "selected"
		}

		if defaultvalue == 3 {
			zaDefa3 = "selected"
		}

		optionz = `<option value=""> -- select -- </option>
					<option value="1" ` + zaDefa1 + `>Yes</option>
					<option value="2" ` + zaDefa2 + `>No</option>
					<option value="3" ` + zaDefa3 + `>Unknown</option>`
		zaField = `<select class="form-control-sm patient-input form-select" name="` + field + `" id="` + field + `" aria-label="` + labs + `">
					` + optionz + `
			      </select>`
	}

	if table == "e_rdt" {
		if defaultString == "Not Done" {
			zaDefa1 = "selected"
		}

		if defaultString == "Oraquick" {
			zaDefa2 = "selected"
		}

		if defaultString == "Others" {
			zaDefa3 = "selected"
		}

		optionz = `<option value=""> -- select -- </option>
					<option value="Not Done" ` + zaDefa1 + `>Not Done</option>
					<option value="Oraquick" ` + zaDefa2 + `>Oraquick</option>
					<option value="Others" ` + zaDefa3 + `>Others</option>`
		zaField = `<select class="form-control-sm patient-input form-select" name="` + field + `" id="` + field + `" aria-label="` + labs + `">
					` + optionz + `
			      </select>`
	}

	if table == "blood" {
		if defaultString == "ND" {
			zaDefa1 = "selected"
		}

		if defaultString == "Arterial" {
			zaDefa2 = "selected"
		}

		if defaultString == "Others" {
			zaDefa3 = "selected"
		}

		optionz = `<option value=""> -- select -- </option>
					<option value="ND" ` + zaDefa1 + `>ND</option>
					<option value="Arterial" ` + zaDefa2 + `>Arterial</option>
					<option value="Venous" ` + zaDefa3 + `>Venous</option>`
		zaField = `<select class="form-control-sm patient-input form-select" name="` + field + `" id="` + field + `" aria-label="` + labs + `">
					` + optionz + `
			      </select>`
	}

	if table == "e_pcr" {
		if defaultString == "Not Done" {
			zaDefa1 = "selected"
		}

		if defaultString == "GeneXpert" {
			zaDefa2 = "selected"
		}

		if defaultString == "Others" {
			zaDefa3 = "selected"
		}

		optionz = `<option value=""> -- select -- </option>
					<option value="Not Done" ` + zaDefa1 + `>Not Done</option>
					<option value="GeneXpert" ` + zaDefa2 + `>GeneXpert</option>
					<option value="Others" ` + zaDefa3 + `>Others</option>`
		zaField = `<select class="form-control-sm patient-input form-select" name="` + field + `" id="` + field + `" aria-label="` + labs + `">
					` + optionz + `
			      </select>`
	}

	if whole == 1 {
		return zaField
	}
	return optionz
}

func GetUser(c *fiber.Ctx, sl *slog.Logger, store *session.Store) (int, string) {
	sess, err := store.Get(c)
	if err != nil {
		sl.Info("Session error")
	}

	userID, ok := sess.Get("user").(int)
	if !ok {
		fmt.Println("Failed to convert session value to int")
		return 0, ""
	}

	username, ok := sess.Get("username").(string)
	if !ok {
		fmt.Println("Failed to convert session value to string")
		return 0, ""
	}

	return userID, username
}

func GetCurrentFacility(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store) int {
	sqlstr := ` SELECT
					facility
				FROM public.users u, public.employee e
				WHERE u.user_employee = e.employee_id AND u.user_id= $1`

	userID, _ := GetUser(c, sl, store)

	var facility int
	if err := db.QueryRowContext(c.Context(), sqlstr, userID).Scan(&facility); err != nil {
		return 0
	}
	return facility
}

func HumanDate(t time.Time) string {
	return t.Format("02 Jan 2006")
}

func HumanDateTime(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func Seq(start, end int) []int {
	s := make([]int, end-start+1)
	for i := range s {
		s[i] = start + i
	}
	return s
}

func IsNullStringEmpty(nullable sql.NullString) bool {
	return !nullable.Valid || nullable.String == ""
}

// GenerateHTML renders an HTML template with given data
func GenerateHTML(c *fiber.Ctx, zdata interface{}, filenames ...string) error {
	var files []string
	for _, file := range filenames {
		files = append(files, GetPath(fmt.Sprintf("../../ui/html/%s.html", file)))
	}

	// Parse templates using the global functions variable
	templates, err := template.New("").Funcs(functions).ParseFiles(files...)
	if err != nil {
		return c.Status(500).SendString("Template parsing error: " + err.Error())
	}

	// Execute template and write output
	c.Set("Content-Type", "text/html")
	if err := templates.ExecuteTemplate(c.Response().BodyWriter(), "layout", zdata); err != nil {
		return c.Status(500).SendString("Template execution error: " + err.Error())
	}

	return nil
}

func GetPath(topFile string) (rtn string) {
	var dirAbsPath string

	ex, err := os.Executable()
	if err == nil {
		dirAbsPath = filepath.Dir(ex)
		rtn = dirAbsPath + "/" + topFile
	} else {
		rtn = topFile
	}

	return
}

func GetParent() (parentPath string) {
	parentPath = filepath.Join("..", "..")
	return
}

type NullableString struct {
	sql.NullString
}

func (ns *NullableString) UnmarshalText(text []byte) error {
	ns.Valid = len(text) > 0
	ns.String = string(text)
	return nil
}

type NullableFloat64 struct {
	sql.NullFloat64
}

type NullableTime struct {
	sql.NullTime
}

func (nf *NullableTime) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		nf.Valid = false
		return nil
	}

	// Try to parse using a standard format (adjust based on your expected input)
	n, err := time.Parse("2006-01-02 15:04:05", string(text))
	if err != nil {
		return err
	}

	nf.Time = n
	nf.Valid = true
	return nil
}

func (nf *NullableFloat64) UnmarshalText(text []byte) error {
	nf.Valid = len(text) > 0
	n, err := strconv.ParseFloat(string(text), 64)
	if err != nil {
		return err
	}
	nf.Float64 = n
	return nil
}

type NullableInt64 struct {
	sql.NullInt64
}

func (ni *NullableInt64) UnmarshalText(text []byte) error {
	ni.Valid = len(text) > 0
	n, err := strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		return err
	}
	ni.Int64 = n
	return nil
}

// DecodeFormData decodes form data into a given struct
func DecodeFormData(c *fiber.Ctx, v interface{}) error {
	// Get all form data and convert it to map[string][]string
	//formData := make(map[string][]string)
	/*
		c.Request().PostArgs().VisitAll(func(key, value []byte) {
			formData[string(key)] = []string{string(value)}
		})
	*/

	postArgs := c.Context().PostArgs()

	// Create a map to store form data
	formData := make(map[string][]string)

	// Visit all arguments and store them in formData map
	postArgs.VisitAll(func(key, value []byte) {
		formData[string(key)] = []string{string(value)}
	})

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true) // Ignore unknown fields

	// Register converters for sql.Null* types
	decoder.RegisterConverter(sql.NullString{}, func(s string) reflect.Value {
		var ns NullableString
		if err := ns.UnmarshalText([]byte(s)); err != nil {
			return reflect.ValueOf(sql.NullString{})
		}
		return reflect.ValueOf(ns.NullString)
	})

	decoder.RegisterConverter(sql.NullFloat64{}, func(s string) reflect.Value {
		var nf NullableFloat64
		if err := nf.UnmarshalText([]byte(s)); err != nil {
			return reflect.ValueOf(sql.NullFloat64{})
		}
		return reflect.ValueOf(nf.NullFloat64)
	})

	decoder.RegisterConverter(sql.NullTime{}, func(s string) reflect.Value {
		var nf NullableTime
		if err := nf.UnmarshalText([]byte(s)); err != nil {
			return reflect.ValueOf(sql.NullTime{})
		}
		return reflect.ValueOf(nf.NullTime)
	})

	decoder.RegisterConverter(sql.NullInt64{}, func(s string) reflect.Value {
		var ni NullableInt64
		if err := ni.UnmarshalText([]byte(s)); err != nil {
			return reflect.ValueOf(sql.NullInt64{})
		}
		return reflect.ValueOf(ni.NullInt64)
	})

	// Decode form data into the struct
	if err := decoder.Decode(v, formData); err != nil {
		fmt.Println("Decoding error:", err)
		return err
	}

	return nil
}

// GetCurrentUser retrieves the current user ID from the session
func GetCurrentUser(c *fiber.Ctx, store *session.Store) int {
	sess, err := store.Get(c)
	if err != nil {
		fmt.Println("Error retrieving session:", err)
		return 0
	}

	userID := sess.Get("user")
	userInt, ok := userID.(int)
	if !ok {
		fmt.Println("User not found or not an int")
		return 0
	}
	return userInt
}

// IsAuthenticated middleware checks if a user is authenticated
func IsAuthenticated(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			fmt.Println("Error retrieving session:", err)
			return c.Redirect("/login", fiber.StatusFound)
		}

		if sess.Get("isAuthenticated") == nil {
			fmt.Println("No active session")
			return c.Redirect("/login", fiber.StatusFound)
		}

		return c.Next()
	}
}

func IzAuthenticated(c *fiber.Ctx, store *session.Store) bool {

	sess, err := store.Get(c)
	if err != nil {
		fmt.Println("Error retrieving session:", err)
		return false
	}

	if sess.Get("isAuthenticated") == nil {
		fmt.Println("No active session")
		return false
	}

	return true
}

func Get_Client_Optionz() (opt map[string]map[string]string) {
	opt = make(map[string]map[string]string)
	// Add data to the map of maps
	opt["sex"] = map[string]string{"": " -- ", "1": "Male", "2": "Female"}
	opt["sex2"] = map[string]string{"": " -- ", "Male": "Male", "Female": "Female"}
	opt["occup"] = map[string]string{"": " -- ", "1": "Healthcare worker", "2": "Non-Healthcare worker"}
	opt["yn"] = map[string]string{"": " -- ", "1": "Yes", "2": "No"}
	opt["yn_extra"] = map[string]string{"": " -- ", "1": "Yes", "2": "No", "3": "Unknown"}
	opt["marital"] = map[string]string{"": " -- ",
		"1": "Married",
		"2": "Cohabiting",
		"3": "Widowed",
		"4": "Separated",
		"5": "Divorced",
		"6": "Single",
	}
	opt["nationality"] = map[string]string{"": " -- ", "1": "Ugandan", "2": "EAC", "3": "Other"}
	opt["ethnicity"] = map[string]string{"": " -- ", "1": "Black", "2": "Other"}
	opt["mental"] = map[string]string{"": " -- ", "a": "A", "v": "V", "p": "P", "u": "U"}

	opt["preg"] = map[string]string{"": " -- ", "1": "Yes", "2": "No", "3": "ND"}
	opt["ward"] = map[string]string{"": " -- ", "1": "Ward", "2": "ICU"}
	opt["result1"] = map[string]string{"": " -- ", "1": "Pos", "2": "Neg", "3": "indeterminate"}
	opt["result2"] = map[string]string{"": " -- ", "1": "Pos", "2": "Neg", "3": "ND"}
	return
}

// {{ GetDBOptions "site" "" "" "employee" "Employee" .FormChild2.Hyperglycemia.Int64}}

func GetDBOptions(table, cat, deflt, fld_name, fld_lab string, deflt_int int64) string {
	sql := ""
	rtn := ""
	switch table {
	case "Employee":
		sql = "SELECT employee_id as code, CONCAT(employee_fname, ' ', employee_lname) AS  lab FROM public.employee"

	case "function":
		sql = "SELECT function_id as code, function_name as lab FROM public.function"
	case "site":
		sql = "SELECT facility_id as code, facility_name as lab FROM public.facility"
	case "test":
	case "meta":
		sql = "Select meta_id as code, meta_name as lab from meta, meta_category WHERE meta.meta_category=meta_category.meta_category_id AND meta_category_name='" + cat + "'"
	}

	res, er := models.GetFields(CtxG.Context(), dbG, sql)
	if er != nil {
		fmt.Println("Error getting fields:", er)
		return ""
	}
	i := 0
	for _, values := range res {

		if deflt == "" {
			zvalue, _ := strconv.ParseInt(values[0], 10, 64)
			if zvalue == deflt_int {
				rtn = rtn + fmt.Sprintf(`<option value="%s" selected>%s</option>`, values[0], values[1])
			} else {
				rtn = rtn + fmt.Sprintf(`<option value="%s">%s</option>`, values[0], values[1])
			}
		} else {
			if values[0] == deflt {
				rtn = rtn + fmt.Sprintf(`<option value="%s" selected>%s</option>`, values[0], values[1])
			} else {
				rtn = rtn + fmt.Sprintf(`<option value="%s">%s</option>`, values[0], values[1])
			}
		}
		i++

	}
	addString := ""
	if i > 1 {
		addString = `<option value=""> -- </option>`
	}

	return `<select class="form-control-sm patient-input form-select" name="` + fld_name + `" id="` + fld_name + `" aria-label="` + fld_lab + `">
			` + addString + rtn + `
			</select>`
}

func GetClientOptionLabel(arrayKey, mapKey string) string {
	options := Get_Client_Optionz()

	if subMap, exists := options[arrayKey]; exists {
		if label, found := subMap[mapKey]; found {
			return label
		}
	}
	return "" // Return an empty string if the keys are not found
}

func GetDBLabel(table, namesFld, indexFld string, indexID int64) string {
	// Use parameterized query to prevent SQL injection
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", namesFld, table, indexFld)
	var label string
	label = ""
	err := dbG.QueryRowContext(CtxG.Context(), query, indexID).Scan(&label)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found for ID:", indexID)
			return ""
		}
		fmt.Println("Error executing query:", err)
		return ""
	}
	return label
}

func GetDBInt(table, namesFld, indexFld, whereString string, indexID int64) int64 {
	// Use parameterized query to prevent SQL injection
	query := ""
	var label int
	label = 0

	var err error

	if whereString == "" {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", namesFld, table, indexFld)
		err = dbG.QueryRowContext(CtxG.Context(), query, indexID).Scan(&label)
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s", namesFld, table, whereString)
		err = dbG.QueryRowContext(CtxG.Context(), query).Scan(&label)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found for ID:", indexID)
			return 0
		}
		fmt.Println("Error executing query:", err)
		return 0
	}
	return int64(label)

}

func ParseNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: value, Valid: true}
}

func ParseNullInt(value string) sql.NullInt64 {
	if value == "" {
		return sql.NullInt64{Valid: false}
	}
	var i int64
	_, err := fmt.Sscanf(value, "%d", &i)
	if err != nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: i, Valid: true}
}

func ParseNullFloat(value string) sql.NullFloat64 {
	if value == "" {
		return sql.NullFloat64{Valid: false}
	}
	var f float64
	_, err := fmt.Sscanf(value, "%f", &f)
	if err != nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: f, Valid: true}
}

func ParseNullTime(value string) sql.NullTime {
	if value == "" {
		return sql.NullTime{Valid: false}
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

// Convert interface{} to sql.NullInt64
func ParseNullInt2(value interface{}) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{Valid: false}
	}

	switch v := value.(type) {
	case float64: // JSON numbers are decoded as float64
		return sql.NullInt64{Int64: int64(v), Valid: true}
	case string:
		if num, err := strconv.ParseInt(v, 10, 64); err == nil {
			return sql.NullInt64{Int64: num, Valid: true}
		}
	}

	return sql.NullInt64{Valid: false}
}

// Convert interface{} to sql.NullString
func ParseNullString2(value interface{}) sql.NullString {
	if value == nil {
		return sql.NullString{Valid: false}
	}

	if str, ok := value.(string); ok && str != "" {
		return sql.NullString{String: str, Valid: true}
	}

	return sql.NullString{Valid: false}
}

// ConvertFiberToGin converts a Fiber context to a Gin context
func ConvertFiberToGin(fctx *fiber.Ctx) (*gin.Context, error) {
	// Create a new HTTP request using Fiber's request data
	req := fctx.Request()

	// Convert Fiber request to standard *http.Request
	httpReq, err := http.NewRequest(
		string(req.Header.Method()),
		fctx.OriginalURL(),
		bytes.NewReader(req.Body()),
	)
	if err != nil {
		return nil, err
	}

	// Copy headers from Fiber to the new request
	req.Header.VisitAll(func(key, value []byte) {
		httpReq.Header.Set(string(key), string(value))
	})

	// Create a new Gin response recorder
	w := httptest.NewRecorder()

	// Create a new Gin context
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = httpReq

	return ginCtx, nil
}
