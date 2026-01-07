package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wekeepgrowing/semo-backend-monorepo/services/geo/internal/usecase"
)

// GeoHandler는 지오로케이션 관련 HTTP 핸들러입니다
type GeoHandler struct {
	geoUseCase *usecase.GeoUseCase
}

// NewGeoHandler는 새로운 GeoHandler 인스턴스를 생성합니다
func NewGeoHandler(geoUseCase *usecase.GeoUseCase) *GeoHandler {
	return &GeoHandler{
		geoUseCase: geoUseCase,
	}
}

// RegisterRoutes는 Echo 라우터에 핸들러 경로를 등록합니다
func (h *GeoHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/geo/ip/:ip", h.GetGeoData)
	e.GET("/geo/city/:ip", h.GetCityInfo)
	e.GET("/geo/country/:ip", h.GetCountryInfo)
	e.GET("/geo/asn/:ip", h.GetASNInfo)
	e.GET("/geo/anonymous/:ip", h.CheckAnonymousIP)
}

// GetGeoData는 IP 주소에 대한 종합적인 지리 정보를 반환합니다
// @Summary IP 주소의 지리 정보 조회
// @Description IP 주소에 대한 종합적인 지리 정보(도시, 국가, ASN 등)를 반환합니다
// @Tags geo
// @Accept json
// @Produce json
// @Param ip path string true "IP 주소"
// @Success 200 {object} usecase.GeoData
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /geo/ip/{ip} [get]
func (h *GeoHandler) GetGeoData(c echo.Context) error {
	ipStr := c.Param("ip")
	if ipStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IP 주소가 필요합니다",
		})
	}

	geoData, err := h.geoUseCase.GetGeoData(ipStr)
	if err != nil {
		status := http.StatusInternalServerError
		if err == usecase.ErrInvalidIPAddress {
			status = http.StatusBadRequest
		}
		return c.JSON(status, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, geoData)
}

// GetCityInfo는 IP 주소에 대한 도시 정보를 반환합니다
// @Summary IP 주소의 도시 정보 조회
// @Description IP 주소에 대한 도시 정보를 반환합니다
// @Tags geo
// @Accept json
// @Produce json
// @Param ip path string true "IP 주소"
// @Success 200 {object} entity.City
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /geo/city/{ip} [get]
func (h *GeoHandler) GetCityInfo(c echo.Context) error {
	ipStr := c.Param("ip")
	if ipStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IP 주소가 필요합니다",
		})
	}

	city, err := h.geoUseCase.GetCityInfo(ipStr)
	if err != nil {
		status := http.StatusInternalServerError
		if err == usecase.ErrInvalidIPAddress {
			status = http.StatusBadRequest
		}
		return c.JSON(status, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, city)
}

// GetCountryInfo는 IP 주소에 대한 국가 정보를 반환합니다
// @Summary IP 주소의 국가 정보 조회
// @Description IP 주소에 대한 국가 정보를 반환합니다
// @Tags geo
// @Accept json
// @Produce json
// @Param ip path string true "IP 주소"
// @Success 200 {object} entity.Country
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /geo/country/{ip} [get]
func (h *GeoHandler) GetCountryInfo(c echo.Context) error {
	ipStr := c.Param("ip")
	if ipStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IP 주소가 필요합니다",
		})
	}

	country, err := h.geoUseCase.GetCountryInfo(ipStr)
	if err != nil {
		status := http.StatusInternalServerError
		if err == usecase.ErrInvalidIPAddress {
			status = http.StatusBadRequest
		}
		return c.JSON(status, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, country)
}

// GetASNInfo는 IP 주소에 대한 ASN 정보를 반환합니다
// @Summary IP 주소의 ASN 정보 조회
// @Description IP 주소에 대한 ASN(자율 시스템 번호) 정보를 반환합니다
// @Tags geo
// @Accept json
// @Produce json
// @Param ip path string true "IP 주소"
// @Success 200 {object} entity.ASN
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /geo/asn/{ip} [get]
func (h *GeoHandler) GetASNInfo(c echo.Context) error {
	ipStr := c.Param("ip")
	if ipStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IP 주소가 필요합니다",
		})
	}

	asn, err := h.geoUseCase.GetASNInfo(ipStr)
	if err != nil {
		status := http.StatusInternalServerError
		if err == usecase.ErrInvalidIPAddress {
			status = http.StatusBadRequest
		}
		return c.JSON(status, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, asn)
}

// CheckAnonymousIP는 IP 주소가 익명 프록시인지 확인합니다
// @Summary IP 주소의 익명성 확인
// @Description IP 주소가 익명 프록시, VPN 또는 Tor 출구 노드인지 확인합니다
// @Tags geo
// @Accept json
// @Produce json
// @Param ip path string true "IP 주소"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /geo/anonymous/{ip} [get]
func (h *GeoHandler) CheckAnonymousIP(c echo.Context) error {
	ipStr := c.Param("ip")
	if ipStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "IP 주소가 필요합니다",
		})
	}

	isAnonymous, err := h.geoUseCase.IsAnonymousIP(ipStr)
	if err != nil {
		status := http.StatusInternalServerError
		if err == usecase.ErrInvalidIPAddress {
			status = http.StatusBadRequest
		} else if err == usecase.ErrFeatureNotSupported {
			return c.JSON(http.StatusOK, map[string]bool{
				"is_anonymous":     false,
				"is_tor_exit_node": false,
				"feature_support":  false,
			})
		}
		return c.JSON(status, map[string]string{
			"error": err.Error(),
		})
	}

	isTorExit, _ := h.geoUseCase.IsTorExitNode(ipStr)

	return c.JSON(http.StatusOK, map[string]bool{
		"is_anonymous":     isAnonymous,
		"is_tor_exit_node": isTorExit,
		"feature_support":  true,
	})
}
