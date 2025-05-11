package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"go-nextjs-dashboard/config"
	config2 "go-nextjs-dashboard/internal/config"
	"go-nextjs-dashboard/model"
	"go-nextjs-dashboard/router"
	my_validator "go-nextjs-dashboard/validator"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"

	// "encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// User contains user information
type User struct {
	FirstName      string     `validate:"required"`
	LastName       string     `validate:"required"`
	Age            uint8      `json:"age" validate:"gte=0,lte=130"`
	Email          string     `validate:"required,email"`
	Gender         string     `validate:"oneof=male female prefer_not_to"`
	FavouriteColor string     `validate:"iscolor"`                // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

func main() {
	f, err := os.OpenFile("fiber.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to open/create fiber.log file: %w", err))
	}
	defer f.Close()
	log.SetOutput(f)

	config2.LoadConfig()
	config.InitValidator()

	// Connect to the db
	config.ConnectDatabase()

	// DB.AutoMigrate(&User{}, &Customer{}, &Invoice{}, &Revenue{})
	// name := "test name"
	// log.Fatal("name: ", name[0:0])
	// SeedUsers()
	// SeedCustomers()
	// SeedInvoices()
	// SeedRevenues()

	// err = config.Validate.Var("0337725c-931c-4338-b60a-0b5c6a787c32", "uuid4")
	// fmt.Println(err.Error())

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		// Override default error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			// log.Errorf("%v", e)
			if c.Get(fiber.HeaderAccept) == fiber.MIMEApplicationJSON {
				if errs, ok := err.(validator.ValidationErrors); ok {
					trans, _ := config.Uni.GetTranslator(c.Get("Accept-Language"))

					// errorFields := make(map[string][]string, len(errs))

					// for _, e := range errs {
					// 	// can translate each error one at a time.
					// 	// fmt.Println(e.Translate(trans))
					// 	errorFields[e.Field()] = append(errorFields[e.Field()], e.Translate(trans))
					// }

					// return c.Status(422).JSON(fiber.Map{"message": "Invalid fields.", "error": errorFields})
					return c.Status(422).JSON(fiber.Map{"message": "Invalid fields.", "error": errs.Translate(trans)})
				}

				var errs my_validator.MapValidationErrors
				if errors.As(err, &errs) {
					trans, _ := config.Uni.GetTranslator(c.Get("Accept-Language"))
					return c.Status(422).JSON(fiber.Map{"message": "Invalid fields.", "error": errs.Translate(trans)})
				}

				log.Error(err.Error())

				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
				return c.Status(code).JSON(fiber.Map{"message": err.Error()})
			}

			log.Error(err.Error())

			// Set Content-Type: text/plain; charset=utf-8
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

			// Return status code with error message
			return c.Status(code).SendString(err.Error())
		},
	})

	app.Use(logger.New())
	app.Use(recover.New(recover.Config{
		Next:             nil,
		EnableStackTrace: true,
		StackTraceHandler: func(_ *fiber.Ctx, e interface{}) {
			log.Panicf("%v\n%s\n", e, debug.Stack())
			// _, _ = os.Stderr.WriteString(fmt.Sprintf("panic: %v\n%s\n", e, debug.Stack())) //nolint:errcheck // This will never fail
		},
	}))
	app.Use(cors.New(cors.ConfigDefault))

	router.RegisterRoutes(app)

	log.Fatal(app.Listen(":" + config2.Cfg.ServerPort))
}

func SeedUsers() {
	users := []model.User{
		{Name: "User", Email: "user@nextmail.com", Password: "123456"},
	}

	for _, user := range users {
		if err := config.DB.Create(&user).Error; err != nil {
			log.Errorf("cannot seed user %v: %v", user.Name, err)
		}
	}
}

func SeedCustomers() {
	customers := []model.Customer{
		{Name: "Delba de Oliveira", Email: "delba@oliveira.com", ImageURL: "/customers/delba-de-oliveira.png"},
		{Name: "Lee Robinson", Email: "lee@robinson.com", ImageURL: "/customers/lee-robinson.png"},
		{Name: "Hector Simpson", Email: "hector@simpson.com", ImageURL: "/customers/hector-simpson.png"},
		{Name: "Steven Tey", Email: "steven@tey.com", ImageURL: "/customers/steven-tey.png"},
		{Name: "Steph Dietz", Email: "steph@dietz.com", ImageURL: "/customers/steph-dietz.png"},
		{Name: "Michael Novotny", Email: "michael@novotny.com", ImageURL: "/customers/michael-novotny.png"},
		{Name: "Evil Rabbit", Email: "evil@rabbit.com", ImageURL: "/customers/evil-rabbit.png"},
		{Name: "Emil Kowalski", Email: "emil@kowalski.com", ImageURL: "/customers/emil-kowalski.png"},
		{Name: "Amy Burns", Email: "amy@burns.com", ImageURL: "/customers/amy-burns.png"},
		{Name: "Balazs Orban", Email: "balazs@orban.com", ImageURL: "/customers/balazs-orban.png"},
	}

	for _, customer := range customers {
		if err := config.DB.Create(&customer).Error; err != nil {
			log.Errorf("cannot seed customer %v: %v", customer.Name, err)
		}
	}
}

// func SeedInvoices() {
// 	var customers []model.Customer
// 	config.DB.Find(&customers)

// 	invoices := []model.Invoice{
// 		{CustomerID: customers[0].ID, Amount: 157.95, Status: "pending", Date: utils.ParseDate("2022-12-06")},
// 		{CustomerID: customers[1].ID, Amount: 203.48, Status: "pending", Date: utils.ParseDate("2022-11-14")},
// 		{CustomerID: customers[4].ID, Amount: 30.40, Status: "paid", Date: utils.ParseDate("2022-10-29")},
// 		{CustomerID: customers[3].ID, Amount: 448.00, Status: "paid", Date: utils.ParseDate("2023-09-10")},
// 		{CustomerID: customers[5].ID, Amount: 345.77, Status: "pending", Date: utils.ParseDate("2023-08-05")},
// 		{CustomerID: customers[7].ID, Amount: 542.46, Status: "pending", Date: utils.ParseDate("2023-07-16")},
// 		{CustomerID: customers[6].ID, Amount: 6.66, Status: "pending", Date: utils.ParseDate("2023-06-27")},
// 		{CustomerID: customers[3].ID, Amount: 325.45, Status: "paid", Date: utils.ParseDate("2023-06-09")},
// 		{CustomerID: customers[4].ID, Amount: 12.50, Status: "paid", Date: utils.ParseDate("2023-06-17")},
// 		{CustomerID: customers[5].ID, Amount: 85.46, Status: "paid", Date: utils.ParseDate("2023-06-07")},
// 		{CustomerID: customers[1].ID, Amount: 5.00, Status: "paid", Date: utils.ParseDate("2023-08-19")},
// 		{CustomerID: customers[5].ID, Amount: 89.45, Status: "paid", Date: utils.ParseDate("2023-06-03")},
// 		{CustomerID: customers[2].ID, Amount: 89.45, Status: "paid", Date: utils.ParseDate("2023-06-18")},
// 		{CustomerID: customers[0].ID, Amount: 89.45, Status: "paid", Date: utils.ParseDate("2023-10-04")},
// 		{CustomerID: customers[2].ID, Amount: 10.00, Status: "paid", Date: utils.ParseDate("2022-06-05")},
// 	}

// 	for _, invoice := range invoices {
// 		if err := config.DB.Create(&invoice).Error; err != nil {
// 			log.Printf("cannot seed invoice %v: %v", invoice, err)
// 		}
// 	}
// }

func SeedRevenues() {
	revenue := []model.Revenue{
		{Month: "Jan", Revenue: 20.00},
		{Month: "Feb", Revenue: 18.00},
		{Month: "Mar", Revenue: 22.00},
		{Month: "Apr", Revenue: 25.00},
		{Month: "May", Revenue: 23.00},
		{Month: "Jun", Revenue: 32.00},
		{Month: "Jul", Revenue: 35.00},
		{Month: "Aug", Revenue: 37.00},
		{Month: "Sep", Revenue: 25.00},
		{Month: "Oct", Revenue: 28.00},
		{Month: "Nov", Revenue: 30.00},
		{Month: "Dec", Revenue: 48.00},
	}

	for _, revenue := range revenue {
		if err := config.DB.Create(&revenue).Error; err != nil {
			log.Errorf("cannot seed revenue %v: %v", revenue, err)
		}
	}
}
