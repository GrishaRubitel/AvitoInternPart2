package main

import (
	"encoding/json"
	"errors"
	"fmt"
	models "grisha_rubitel/avito_p2/dbModels"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var log = logrus.New()

func main() {
	envMap, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	conn := envMap["POSTGRES_CONN"]
	fmt.Print(conn)
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		log.Warn("DB connection failed: ", err)
		return
	}

	err = models.TenderMigrate(db)
	if err != nil {
		log.Warn("DB tender migration failed: ", err)
		return
	}

	err = models.BindMigrate(db)
	if err != nil {
		log.Warn("DB bind migration failed: ", err)
		return
	}

	router := gin.Default()
	router.Use(cors.Default())

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/ping", func(c *gin.Context) { //ready
			code, resp, err := CheckServer(conn)
			responseReturner(code, resp, err, c)
		})

		tendersGroup := apiGroup.Group("/tenders")
		{
			tendersGroup.GET("", func(c *gin.Context) { //ready
				bodyData := readQueryParams(c)
				code, resp, err := GetTenders(db, bodyData)
				responseReturner(code, resp, err, c)
			})
			tendersGroup.POST("/new", func(c *gin.Context) { //ready
				bodyData, err := readBodyData(c)
				if err != nil {
					responseReturner(http.StatusBadRequest, "", errors.New("error while reading request's body"), c)
				} else {
					code, resp, err := CreateTender(db, bodyData)
					responseReturner(code, resp, err, c)
				}
			})
			// ВАЖНЫЙ КОММЕНТАРИЙ
			// Так как по заданию не дана ассоциативная таблица tenders-employee, а тендер не хранит информацию о конкретном создателе,
			// то я принял решение, что любой пользователь организации, создавшей тендер, может просматривать информацию о
			// соответствующих тендерах
			tendersGroup.GET("/my", func(c *gin.Context) { //ready
				queryData := readQueryParams(c)
				code, resp, err := GetUserTenders(db, queryData)
				responseReturner(code, resp, err, c)
			})
			tendersGroup.GET("/:tenderId/status", func(c *gin.Context) { //ready
				pathData := readQueryParams(c)
				pathData["tenderid"] = c.Param("tenderId")
				code, resp, err := GetTenderStatus(db, pathData)
				responseReturner(code, resp, err, c)
			})
			tendersGroup.PUT("/:tenderId/status", func(c *gin.Context) { //ready i guess
				pathData := readQueryParams(c)
				pathData["tenderid"] = c.Param("tenderId")
				code, resp, err := UpdateTenderStatus(db, pathData)
				responseReturner(code, resp, err, c)
			})
			tendersGroup.PATCH("/:tenderId/edit", func(c *gin.Context) { //ready
				pathData, err := readBodyData(c)
				if err != nil {
					responseReturner(http.StatusBadRequest, "", err, c)
				} else {
					pathData = mergeMaps(pathData, readQueryParams(c))
					pathData["tenderid"] = c.Param("tenderId")
					code, resp, err := EditTender(db, pathData)
					responseReturner(code, resp, err, c)
				}
			})
			tendersGroup.PUT("/:tenderId/rollback/:version", func(c *gin.Context) { //skip for now
				log.Info("Popa")
			})
		}

		bidsGroup := apiGroup.Group("/bids")
		{
			bidsGroup.POST("/new", func(c *gin.Context) {
				bodyData, err := readBodyData(c)
				if err != nil {
					responseReturner(http.StatusBadRequest, "", err, c)
				} else {
					code, resp, err := CreateBid(db, bodyData)
					responseReturner(code, resp, err, c)
				}
			})
			bidsGroup.GET("/my", func(c *gin.Context) { //ready
				pathData := readQueryParams(c)
				code, resp, err := GetUserBids(db, pathData)
				responseReturner(code, resp, err, c)
			})
			bidsGroup.GET("/:id/list", func(c *gin.Context) { // should be tenderId --- ready
				pathData := readQueryParams(c)
				pathData["tenderid"] = c.Param("id")
				code, resp, err := GetBidsForTender(db, pathData)
				responseReturner(code, resp, err, c)
			})
			bidsGroup.GET("/:id/status", func(c *gin.Context) { // ready
				pathData := readQueryParams(c)
				pathData["bidid"] = c.Param("id")
				code, resp, err := GetBidStatus(db, pathData)
				responseReturner(code, string(resp.Status), err, c)
			})
			bidsGroup.PUT("/:id/status", func(c *gin.Context) { // ready
				pathData := readQueryParams(c)
				pathData["bidid"] = c.Param("id")
				code, resp, err := UpdateBidStatus(db, pathData)
				responseReturner(code, resp, err, c)
			})
			bidsGroup.PATCH("/:id/edit", func(c *gin.Context) { // ready
				pathData, err := readBodyData(c)
				if err != nil {
					responseReturner(http.StatusBadRequest, "", err, c)
				} else {
					pathData = mergeMaps(pathData, readQueryParams(c))
					pathData["bidid"] = c.Param("id")
					code, resp, err := EditBid(db, pathData)
					responseReturner(code, resp, err, c)
				}
			})
			bidsGroup.PUT("/:id/submit_decision", func(c *gin.Context) {
				pathData := readQueryParams(c)
				pathData["bidid"] = c.Param("id")
				code, resp, err := SubmitBidDecision(db, pathData)
				responseReturner(code, resp, err, c)
			})
			bidsGroup.PUT("/:id/feedback", func(c *gin.Context) {
				log.Info("Popa")
			})
			bidsGroup.PUT("/:id/rollback/:version", func(c *gin.Context) { // skip for now
				log.Info("Popa")
			})
			bidsGroup.GET("/:id/reviews", func(c *gin.Context) { // should be tenderId
				log.Info("Popa")
			})
		}
	}

	router.Run(envMap["SERVER_ADDRESS"])
}

func responseReturner(code int, resp string, err error, c *gin.Context) {
	if err != nil {
		log.Warn(err)
		c.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
	}
	c.String(code, resp)
}

func readBodyData(c *gin.Context) (map[string]string, error) {
	var bodyData map[string]interface{}

	body, err := c.GetRawData()
	if err != nil {
		responseReturner(http.StatusBadRequest, "", err, c)
		return nil, err
	}

	if err := json.Unmarshal(body, &bodyData); err != nil {
		responseReturner(http.StatusBadRequest, "", err, c)
		return nil, err
	}

	result := make(map[string]string)

	for key, value := range bodyData {
		switch v := value.(type) {
		case string:
			result[key] = v
		case float64:
			result[key] = strconv.FormatFloat(v, 'f', -1, 64)
		default:
			return nil, errors.New("unsupported value type")
		}
	}

	return result, nil
}

func readQueryParams(c *gin.Context) map[string]string {
	paramsData := make(map[string]string)

	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			paramsData[key] = values[0]
		}
	}

	return paramsData
}

func mergeMaps(m1 map[string]string, m2 map[string]string) map[string]string {
	merged := make(map[string]string)
	for k, v := range m1 {
		merged[k] = v
	}
	for key, value := range m2 {
		merged[key] = value
	}
	return merged
}
