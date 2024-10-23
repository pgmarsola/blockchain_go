package apis

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func Client() {

	http.HandleFunc("/allwallets", handleAllWallets)
	http.HandleFunc("/createwallet", handleCreateWallet)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Interface() {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	app.Static("/", "./static")

	app.Get("/", RenderHome)
	app.Listen(":8081")
}

func RenderHome(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{})
}
