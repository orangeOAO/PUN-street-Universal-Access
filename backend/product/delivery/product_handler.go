package delivery

import (
	"strconv"

	"github.com/PUArallelepiped/PUN-street-Universal-Access/domain"
	"github.com/PUArallelepiped/PUN-street-Universal-Access/swagger"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	ProductUsecase domain.ProductUsecase
}

func NewProductHandler(e *gin.Engine, productUsecase domain.ProductUsecase) {
	handler := &ProductHandler{
		ProductUsecase: productUsecase,
	}

	store := e.Group("/api/v1/store/:storeID")
	{
		store.GET("/products", handler.GetProductById)
		store.POST("/add-product", handler.AddProduct)
		store.PUT("/update-product/:productID", handler.UpdateProduct)
	}
}

func (s *ProductHandler) GetProductById(c *gin.Context) {
	storeID, err := strconv.ParseInt(c.Param("storeID"), 10, 64)
	if err != nil {
		logrus.Error(err)
		c.Status(400)
		return
	}
	products, err := s.ProductUsecase.GetByID(c, storeID)
	if err != nil {
		logrus.Error(err)
		c.Status(500)
		return
	}
	c.JSON(200, products)
}

func (s *ProductHandler) AddProduct(c *gin.Context) {
	storeID, err := strconv.ParseInt(c.Param("storeID"), 10, 64)
	if err != nil {
		logrus.Error(err)
		c.Status(400)
		return
	}
	var product swagger.ProductInfo
	if err := c.BindJSON(&product); err != nil {
		logrus.Error(err)
		c.Status(400)
		return
	}

	err = s.ProductUsecase.AddByStoreId(c, storeID, &product)
	if err != nil {
		logrus.Error(err)
		c.Status(500)
		return
	}

	c.Status(200)
}

func (s *ProductHandler) UpdateProduct(c *gin.Context) {
	storeID, err := strconv.ParseInt(c.Param("storeID"), 10, 64)
	if err != nil {
		logrus.Error(err)
		c.Status(400)
		return
	}
	productID, err := strconv.ParseInt(c.Param("productID"), 10, 64)
	if err != nil {
		logrus.Error(err)
		c.Status(400)
		return
	}

	var product swagger.ProductInfo
	if err := c.BindJSON(&product); err != nil {
		logrus.Error(err)
		c.Status(400)
		return
	}

	err = s.ProductUsecase.UpdateById(c, storeID, productID, &product)
	if err != nil {
		logrus.Error(err)
		c.Status(500)
		return
	}

	c.Status(200)
}
