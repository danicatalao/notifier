package forecast

import (
	"net/http"

	l "github.com/danicatalao/notifier/internal/logger"
	"github.com/gin-gonic/gin"
)

// Forecast provider wrapper for testing

type ForecastHandler struct {
	client  ForecastApiClient
	service ForecastService
	log     l.Logger
}

func NewForecastHandler(c ForecastApiClient, s ForecastService, l l.Logger) *ForecastHandler {
	return &ForecastHandler{client: c, service: s, log: l}
}

func (h *ForecastHandler) AddForecastRoutes(r *gin.RouterGroup) {

	forecastRoutes := r.Group("/forecast")
	{
		forecastRoutes.GET("/city", h.GetCityIDByName)
		forecastRoutes.GET("/:id", h.GetForecast)
		forecastRoutes.GET("/:id/wave/:day", h.GetWaves)
		forecastRoutes.GET("/fw/:cityName", h.GetForeCastAndWave)
	}
}

func (h *ForecastHandler) GetCityIDByName(c *gin.Context) {
	cityName := c.Query("name")

	city, err := h.client.GetCity(c, cityName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"CityID": city.Id,
	})
}

func (h *ForecastHandler) GetForecast(c *gin.Context) {
	cityId := c.Param("id")

	forecast, err := h.client.GetCityForecast(c, cityId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, forecast)
}

func (h *ForecastHandler) GetWaves(c *gin.Context) {
	cityId := c.Param("id")
	day := c.Param("day")

	wave, err := h.client.GetWaveForecast(c, cityId, day)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wave)
}

func (h *ForecastHandler) GetForeCastAndWave(c *gin.Context) {
	cityName := c.Param("cityName")

	forecastWave, err := h.service.GetForecastAndWave(c, cityName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, forecastWave)
}
