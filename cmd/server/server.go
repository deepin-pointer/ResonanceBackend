package server

import (
	"log"
	"time"
	"unsafe"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"

	"github.com/spf13/viper"

	"rsbackend/internal/controller"
	"rsbackend/internal/model"
)

type Server struct {
	staticData  *controller.StaticData
	dynamicData *controller.DynamicData
	fiberServer *fiber.App
}

func newServer() *Server {
	staticData := controller.NewStaticData(viper.GetString("static_file"))
	return &Server{
		staticData: staticData,
		dynamicData: controller.NewDynamicData(
			len(staticData.GoodsList),
			len(staticData.CityList),
		),
	}
}

func (s *Server) Run(port string) {
	go s.dynamicData.LoggingWorker(viper.GetString("dynamic_file"))

	s.fiberServer = fiber.New()

	s.fiberServer.Use(limiter.New(limiter.Config{
		Max:        2,
		Expiration: 5 * time.Second,
	}))

	s.fiberServer.Static("/", "./web")

	s.fiberServer.Get("/api/static", s.static)
	s.fiberServer.Get("/api/dynamic", s.dynamic)

	s.fiberServer.Post("/api/login", s.login)
	s.fiberServer.Post("/api/report_price", s.reportPrice)

	s.fiberServer.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: viper.GetString("sign_key")},
	}))

	s.fiberServer.Post("/api/new_city", s.newCity)
	s.fiberServer.Post("/api/new_goods", s.newGoods)

	log.Fatal(s.fiberServer.Listen(viper.GetString("bind")))
}

func (s *Server) Shutdown() {
	s.fiberServer.Shutdown()
}

func (s *Server) login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("hash")

	users := viper.GetStringMapString("users")

	if hash, ok := users[user]; !ok || hash != pass {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"name":  user,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(viper.GetString("sign_key")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func (s *Server) static(c *fiber.Ctx) error {
	err := c.Status(fiber.StatusOK).Send(s.staticData.GetData())
	c.Set("content-type", "application/json; charset=utf-8")
	return err
}

func (s *Server) dynamic(c *fiber.Ctx) error {
	data := s.dynamicData.GetData()
	byteSlice := unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(data))), len(data)*8)
	err := c.Status(fiber.StatusOK).Send(byteSlice)
	if err != nil {
		return err
	}
	c.Set("content-type", "application/octet-stream")
	return nil
}

func (s *Server) reportPrice(c *fiber.Ctx) error {
	data := new([]model.PriceRecord)
	c.BodyParser(data)
	s.dynamicData.ModifyPrice(data)
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) newCity(c *fiber.Ctx) error {
	data := new(model.City)
	c.BodyParser(data)
	s.staticData.NewCity(data)
	s.dynamicData.AddCity(1)
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) newGoods(c *fiber.Ctx) error {
	data := new(model.Goods)
	c.BodyParser(data)
	s.staticData.NewGoods(data)
	s.dynamicData.AddGoods(1)
	return c.SendStatus(fiber.StatusOK)
}
