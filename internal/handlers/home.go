package handlers

import (
	"case/internal/models"
	"case/internal/security"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	_ "github.com/lib/pq"
)

func HandlerHome(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {
	fmt.Println("starting home")

	userID, userName := GetUser(c, sl, store)
	role := security.GetRoles(userID, "admin")

	data := NewTemplateData(c, store)
	data.User = userName
	data.Role = role
	sID := fmt.Sprintf("%d", userID)
	metums, _ := models.Metums(c.Context(), db, " user_id = "+sID) //meta_category = 3 AND

	zamenu := ""
	for _, menu := range metums {
		zamenu = zamenu + "<a href='" + menu.MetaLink.String + "' class='menu-item'>" + menu.MetaDescription.String + "</a>"
	}

	data.Form = metums
	data.Menuz = zamenu

	fmt.Println("loading home page")
	return GenerateHTML(c, data, "home")

}

func HandlerLoginForm(c *fiber.Ctx, sl *slog.Logger, store *session.Store, config Config) error {

	sess, err := store.Get(c)

	if err == nil {
		userID := sess.Get("user")
		if userID != nil {
			sl.Info("Session error, No ID set")
			return c.Redirect("/", 302)
		}
	}

	// load page
	data := map[string]string{"Title": "Login Page"}
	return GenerateHTML(c, data, "login")
}

func HandlerLoginSubmit(c *fiber.Ctx, db *sql.DB, sl *slog.Logger, store *session.Store, config Config) error {

	sess, err := store.Get(c)
	if err == nil {
		userID := sess.Get("user")
		if userID != nil {
			fmt.Println("Already logged in")
			return c.Redirect("/", 302)
		}
	}

	// Extract form values
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		fmt.Println("No username and or password provided")
		c.Status(fiber.StatusBadRequest)      // Set HTTP 400 status
		return c.Redirect("/login?error=400") // Redirect to login page
	}

	id, er := models.Authenticate(c.Context(), db, username, password)
	if er != nil {
		fmt.Println("Failed Authentication: ", er.Error())
		return c.Redirect("/login?error=afail")
	}

	if id > 0 {
		// Get session
		sess, err := store.Get(c)
		if err != nil {
			//return c.Status(fiber.StatusInternalServerError).SendString("Session error")
			sl.Info("Session error")
			return c.Redirect("/login?serror")
		}

		// Set session variables
		sess.Set("user", id) // Example: Set user ID
		sess.Set("username", username)
		sess.Set("isAuthenticated", true)

		// Save session
		if err := sess.Save(); err != nil {
			//return c.Status(fiber.StatusInternalServerError).SendString("Failed to save session")
			sl.Info("Failed to save session")
			return c.Redirect("/login?sfail")
		}

		// Redirect to dashboard
		return c.Redirect("/?goodnes=1")
	}

	return nil
}

func HandlerLoginOut(c *fiber.Ctx, sl *slog.Logger, store *session.Store, config Config) error {
	sess, err := store.Get(c)
	if err != nil {
		//return c.Status(fiber.StatusInternalServerError).SendString("Session error")
		sl.Info("Session error")
		return c.Redirect("/login")
	}

	// Destroy session
	sess.Destroy()
	return c.Redirect("/login")
}

func HandlerLoginForgot(c *fiber.Ctx, sl *slog.Logger, store *session.Store, config Config) error {

	sess, err := store.Get(c)
	if err == nil {
		userID := sess.Get("user")
		if userID != nil {
			return c.Redirect("/", 302)
		}
	}

	// load page

	data := map[string]string{"Title": "Forgot Password and/or username"}
	return GenerateHTML(c, data, "forgot")
}

func HandlerHelp(c *fiber.Ctx, sl *slog.Logger, store *session.Store, config Config) error {
	data := map[string]string{"Title": "Help Page"}
	return GenerateHTML(c, data, "help")
}
