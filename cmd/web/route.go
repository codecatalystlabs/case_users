package main

import (
	"database/sql"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	_ "github.com/lib/pq"

	"case/internal/handlers"
)

func SetRoute(app *fiber.App, db *sql.DB, store *session.Store, sl *slog.Logger, config handlers.Config) {
	RouteHome(app, db, sl, store, config)

	app.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})

	// Main application routes
	appGroup := app.Group("/")
	appGroup.Use(AuthRequired(store)) // Apply middleware for protected routes
	{
		// Home route
		appGroup.Get("/", func(c *fiber.Ctx) error { return handlers.HandlerHome(c, db, sl, store, config) })

		// Add more routes as needed...

		api := app.Group("/api") // Group for all API routes

		enk := api.Group("/encounter")
		sym := api.Group("/sym")
		mob := api.Group("/mob")
		rus := api.Group("/rush")
		lab := api.Group("/lab")
		sta := api.Group("/status")

		emp := app.Group("/employees") // Employees
		usr := app.Group("/users")     // users
		hfs := app.Group("/secure")    // Health facilities
		cse := app.Group("/cases")

		//enc := app.Group("/encounter")
		dis := app.Group("/discharge")

		// Additional routes
		RouteFacilities(hfs, db, sl, config)
		RouteUsers(usr, db, sl, config)
		RouteCases(cse, db, sl, config)
		RouteMorbidity(mob, db, sl, config)
		RouteSymptoms(sym, db, sl, config)
		RouteRush(rus, db, sl, config)
		RouteLab(lab, db, sl, config)

		RouteEmployees(emp, db, sl, config)
		RouteDischarge(dis, db, sl, config)

		//RouteLab(enc, db, sl, config)
		//RouteLab(lab, db, sl, config)

		RouteAPIEncounter(enk, db, sl, config)
		RouteAPIStatus(sta, db, sl, config)
	}
}

func AuthRequired(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return err
		}
		userID := sess.Get("user")
		if userID == nil {
			return c.Redirect("/login", 302)
		}

		// Store user ID in Fiber Locals for later use
		c.Locals("userID", userID)

		return c.Next()
	}
}

func RouteAPIEncounter(v fiber.Router, db *sql.DB, sl *slog.Logger, config handlers.Config) {
	v.Get("/", func(c *fiber.Ctx) error { return handlers.HandlerAPIGetEncounter(c, db, sl, store, config) })
}

func RouteAPIStatus(v fiber.Router, db *sql.DB, sl *slog.Logger, config handlers.Config) {
	v.Get("/list", func(c *fiber.Ctx) error { return handlers.HandlerAPIGetStatuses(c, db, sl, store, config) })
	v.Post("/save", func(c *fiber.Ctx) error { return handlers.HandlerAPIPostStatus(c, db, sl, store, config) })
}

func RouteDischarge(v fiber.Router, db *sql.DB, sl *slog.Logger, config handlers.Config) {
	v.Get("/list", func(c *fiber.Ctx) error { return handlers.GetDischarge(c, db, sl, store, config) })
	v.Get("/certificate", func(c *fiber.Ctx) error { return handlers.Certificate(c, db, sl, store, config) })
	v.Get("/verify", func(c *fiber.Ctx) error { return handlers.VerifyDischarge(c, db, sl, store, config) })
	v.Post("/save", func(c *fiber.Ctx) error { return handlers.Discharge(c, db, sl, store, config) })
}

func RouteHome(app *fiber.App, db *sql.DB, sl *slog.Logger, store *session.Store, config handlers.Config) {
	app.Get("/login", func(c *fiber.Ctx) error { return handlers.HandlerLoginForm(c, sl, store, config) })
	app.Post("/login", func(c *fiber.Ctx) error { return handlers.HandlerLoginSubmit(c, db, sl, store, config) })
	app.Get("/logout", func(c *fiber.Ctx) error { return handlers.HandlerLoginOut(c, sl, store, config) })
	app.Get("/forget", func(c *fiber.Ctx) error { return handlers.HandlerLoginForgot(c, sl, store, config) })
	app.Get("/help", func(c *fiber.Ctx) error { return handlers.HandlerHelp(c, sl, store, config) })
}

func RouteFacilities(v fiber.Router, db *sql.DB, sl *slog.Logger, config handlers.Config) {

	v.Get("/new/:i", func(c *fiber.Ctx) error { return handlers.HandlerFacilityForm(c, db, sl, store, config) })
	v.Post("/save", func(c *fiber.Ctx) error { return handlers.HandlerFacilitySubmit(c, db, sl, store, config) })
	v.Post("/filter", func(c *fiber.Ctx) error { return handlers.HandlerFacilityList(c, db, sl, store, config) })
	v.Get("/list", func(c *fiber.Ctx) error { return handlers.HandlerFacilityList(c, db, sl, store, config) })
	v.Get("/", func(c *fiber.Ctx) error { return handlers.HandlerFacilityList(c, db, sl, store, config) })
}

func RouteUsers(v fiber.Router, db *sql.DB, sl *slog.Logger, config handlers.Config) {

	v.Get("/new/:i", func(c *fiber.Ctx) error { return handlers.HandlerUserForm(c, db, sl, store, config) })
	v.Post("/save", func(c *fiber.Ctx) error { return handlers.HandlerUserSubmit(c, db, sl, store, config) })
	v.Post("/filter", func(c *fiber.Ctx) error { return handlers.HandlerUserList(c, db, sl, store, config) })
	v.Get("/list", func(c *fiber.Ctx) error { return handlers.HandlerUserList(c, db, sl, store, config) })
	v.Get("/", func(c *fiber.Ctx) error { return handlers.HandlerUserList(c, db, sl, store, config) })
}

func RouteEmployees(v fiber.Router, db *sql.DB, sl *slog.Logger, config handlers.Config) {
	v.Get("/new/:i", func(c *fiber.Ctx) error { return handlers.HandlerEmployeeForm(c, db, sl, store, config) })
	v.Post("/save", func(c *fiber.Ctx) error { return handlers.HandlerEmployeeSubmit(c, db, sl, store, config) })
	v.Post("/filter", func(c *fiber.Ctx) error { return handlers.HandlerEmployeeList(c, db, sl, store, config) })
	v.Get("/list", func(c *fiber.Ctx) error { return handlers.HandlerEmployeeList(c, db, sl, store, config) })
	v.Get("/", func(c *fiber.Ctx) error { return handlers.HandlerEmployeeList(c, db, sl, store, config) })
}

func RouteCases(v fiber.Router, db *sql.DB, sl *slog.Logger, config handlers.Config) {

	v.Get("/new/:i", func(c *fiber.Ctx) error { return handlers.HandlerCasesForm(c, db, sl, store, config) })
	v.Post("/save", func(c *fiber.Ctx) error { return handlers.HandlerCasesSubmit(c, db, sl, store, config) })
	v.Post("/filter", func(c *fiber.Ctx) error { return handlers.HandlerCasesList(c, db, sl, store, config) })
	v.Get("/list", func(c *fiber.Ctx) error { return handlers.HandlerCasesList(c, db, sl, store, config) })
	v.Get("/", func(c *fiber.Ctx) error { return handlers.HandlerCasesList(c, db, sl, store, config) })

	///cases/encounters/list/1

	v.Get("/encounters/list/:i", func(c *fiber.Ctx) error { return handlers.HandlerCaseEncounterForm(c, db, sl, store, config) })   //+
	v.Get("/encounters/new/:i/:j", func(c *fiber.Ctx) error { return handlers.HandlerCaseEncounterForm(c, db, sl, store, config) }) //+
	v.Post("/encounters/save", func(c *fiber.Ctx) error { return handlers.HandlerCaseEncounterSubmit(c, db, sl, store, config) })   //+
}

func RouteCaseDischarge(v fiber.Router, db *sql.DB, sl *slog.Logger, config handlers.Config) { //+
	//+
	v.Get("/view/:i/:j", func(c *fiber.Ctx) error { return handlers.HandlerCasesForm(c, db, sl, store, config) })    //+
	v.Get("/new/:i/:j", func(c *fiber.Ctx) error { return handlers.HandlerCasesForm(c, db, sl, store, config) })     //+
	v.Post("/save/:i/:j", func(c *fiber.Ctx) error { return handlers.HandlerCasesSubmit(c, db, sl, store, config) }) //+

}
func RouteSymptoms(v fiber.Router, db *sql.DB, sl *slog.Logger, config handlers.Config) {

	v.Get("/new/:i", func(c *fiber.Ctx) error { return handlers.HandlerSymptomsForm(c, db, sl, store, config) })
	v.Post("/save", func(c *fiber.Ctx) error { return handlers.HandlerSymptomsSubmit(c, db, sl, store, config) })
	v.Post("/filter", func(c *fiber.Ctx) error { return handlers.HandlerSymptomsList(c, db, sl, store, config) })
	v.Get("/list", func(c *fiber.Ctx) error { return handlers.HandlerSymptomsList(c, db, sl, store, config) })
	v.Get("/", func(c *fiber.Ctx) error { return handlers.HandlerSymptomsList(c, db, sl, store, config) })
}

func RouteMorbidity(v fiber.Router, db *sql.DB, sl *slog.Logger, config handlers.Config) {

	v.Get("/new/:i", func(c *fiber.Ctx) error { return handlers.HandlerMorbidityForm(c, db, sl, store, config) })
	v.Post("/save", func(c *fiber.Ctx) error { return handlers.HandlerMorbiditySubmit(c, db, sl, store, config) })
	v.Post("/filter", func(c *fiber.Ctx) error { return handlers.HandlerMorbidityList(c, db, sl, store, config) })
	v.Get("/list", func(c *fiber.Ctx) error { return handlers.HandlerMorbidityList(c, db, sl, store, config) })
	v.Get("/", func(c *fiber.Ctx) error { return handlers.HandlerMorbidityList(c, db, sl, store, config) })
}

func RouteRush(v fiber.Router, db *sql.DB, sl *slog.Logger, config handlers.Config) {

	v.Get("/new/:i", func(c *fiber.Ctx) error { return handlers.HandlerRushForm(c, db, sl, store, config) })
	v.Post("/save", func(c *fiber.Ctx) error { return handlers.HandlerRushSubmit(c, db, sl, store, config) })
	v.Post("/filter", func(c *fiber.Ctx) error { return handlers.HandlerRushList(c, db, sl, store, config) })
	v.Get("/list", func(c *fiber.Ctx) error { return handlers.HandlerRushList(c, db, sl, store, config) })
	v.Get("/", func(c *fiber.Ctx) error { return handlers.HandlerRushList(c, db, sl, store, config) })
}

func RouteLab(v fiber.Router, db *sql.DB, sl *slog.Logger, config handlers.Config) {

	v.Get("/new/:i", func(c *fiber.Ctx) error { return handlers.HandlerLabForm(c, db, sl, store, config) })
	v.Post("/save", func(c *fiber.Ctx) error { return handlers.HandlerLabSubmit(c, db, sl, store, config) })
	v.Post("/filter", func(c *fiber.Ctx) error { return handlers.HandlerLabList(c, db, sl, store, config) })
	v.Get("/list", func(c *fiber.Ctx) error { return handlers.HandlerLabList(c, db, sl, store, config) })
	v.Get("/", func(c *fiber.Ctx) error { return handlers.HandlerLabList(c, db, sl, store, config) })
}
