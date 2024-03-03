package server

import (
	"log"
	"reflect"
	"time"
	"unsafe"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"rsbackend/internal/controller"
	"rsbackend/internal/model"
)

type Server struct {
	staticData  *controller.StaticData
	dynamicData *controller.DynamicData
}

func newServer(staticDataPath string) *Server {
	staticData := controller.NewStaticData(staticDataPath)
	return &Server{
		staticData: staticData,
		dynamicData: controller.NewDynamicData(
			len(staticData.GoodsList),
			len(staticData.CityList),
		),
	}
}

func (s *Server) Run(port string) {
	app := fiber.New()
	app.Get("/", s.web)
	app.Get("/static", s.static)
	app.Get("/dynamic", s.dynamic)
	app.Post("/report_price", s.reportPrice)
	app.Post("/new_city", s.newCity)
	app.Post("/new_goods", s.newGoods)

	app.Use(limiter.New(limiter.Config{
		Max:        2,
		Expiration: 5 * time.Second,
	}))
	log.Fatal(app.Listen(port))
}

func (s *Server) web(c *fiber.Ctx) error {
	return c.SendString("UnderDevelop")
}

func (s *Server) static(c *fiber.Ctx) error {
	err := c.Status(fiber.StatusOK).Send(s.staticData.GetData())
	c.Set("content-type", "application/json; charset=utf-8")
	return err
}

func (s *Server) dynamic(c *fiber.Ctx) error {
	data := s.dynamicData.GetData()
	header := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	header.Len *= 8 // int64 has 8 bytes
	header.Cap *= 8

	byteSlice := *(*[]byte)(unsafe.Pointer(&header))
	return c.Status(fiber.StatusOK).Send(byteSlice)
}

func (s *Server) reportPrice(c *fiber.Ctx) error {
	data := new(model.PriceRecord)
	c.BodyParser(data)
	s.dynamicData.ModifyPrice(data)
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) newCity(c *fiber.Ctx) error {
	data := new(model.City)
	c.BodyParser(data)
	s.staticData.NewCity(data)
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) newGoods(c *fiber.Ctx) error {
	data := new(model.Goods)
	c.BodyParser(data)
	s.staticData.NewGoods(data)
	return c.SendStatus(fiber.StatusOK)
}
