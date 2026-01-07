package usecase

import (
	"net"

	"github.com/wekeepgrowing/semo-backend-monorepo/services/geo/internal/domain/entity"
	"github.com/wekeepgrowing/semo-backend-monorepo/services/geo/internal/domain/repository"
)

// GeoUseCase는 지오로케이션 관련 유스케이스를 담당합니다
type GeoUseCase struct {
	cityRepo      repository.GeoIP2CityRepository
	countryRepo   repository.GeoIP2CountryRepository
	asnRepo       repository.GeoLite2ASNRepository
	anonymousRepo repository.GeoIP2AnonymousIPRepository
}

// NewGeoUseCase는 새로운 GeoUseCase 인스턴스를 생성합니다
func NewGeoUseCase(
	cityRepo repository.GeoIP2CityRepository,
	countryRepo repository.GeoIP2CountryRepository,
	asnRepo repository.GeoLite2ASNRepository,
	anonymousRepo repository.GeoIP2AnonymousIPRepository,
) *GeoUseCase {
	return &GeoUseCase{
		cityRepo:      cityRepo,
		countryRepo:   countryRepo,
		asnRepo:       asnRepo,
		anonymousRepo: anonymousRepo,
	}
}

// NewGeoUseCaseWithGeoLite2 은 GeoLite2 통합 리포지토리를 사용하는 GeoUseCase 인스턴스를 생성합니다
func NewGeoUseCaseWithGeoLite2(repo repository.GeoLite2Repository) *GeoUseCase {
	return &GeoUseCase{
		cityRepo:      repo,
		countryRepo:   repo,
		asnRepo:       repo,
		anonymousRepo: nil, // GeoLite2에는 Anonymous IP 데이터가 없습니다
	}
}

// GetCityInfo는 IP 주소에 대한 도시 정보를 조회합니다
func (uc *GeoUseCase) GetCityInfo(ipStr string) (entity.City, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return entity.City{}, ErrInvalidIPAddress
	}

	return uc.cityRepo.GetCity(ip)
}

// GetCountryInfo는 IP 주소에 대한 국가 정보를 조회합니다
func (uc *GeoUseCase) GetCountryInfo(ipStr string) (entity.Country, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return entity.Country{}, ErrInvalidIPAddress
	}

	return uc.countryRepo.GetCountry(ip)
}

// GetASNInfo는 IP 주소에 대한 ASN 정보를 조회합니다
func (uc *GeoUseCase) GetASNInfo(ipStr string) (entity.ASN, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return entity.ASN{}, ErrInvalidIPAddress
	}

	return uc.asnRepo.GetASN(ip)
}

// IsAnonymousIP는 IP 주소가 익명 프록시, VPN 등을 사용하는지 확인합니다
func (uc *GeoUseCase) IsAnonymousIP(ipStr string) (bool, error) {
	if uc.anonymousRepo == nil {
		return false, ErrFeatureNotSupported
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, ErrInvalidIPAddress
	}

	anonIP, err := uc.anonymousRepo.GetAnonymousIP(ip)
	if err != nil {
		return false, err
	}

	return anonIP.IsAnonymous, nil
}

// IsTorExitNode는 IP 주소가 Tor 출구 노드인지 확인합니다
func (uc *GeoUseCase) IsTorExitNode(ipStr string) (bool, error) {
	if uc.anonymousRepo == nil {
		return false, ErrFeatureNotSupported
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, ErrInvalidIPAddress
	}

	anonIP, err := uc.anonymousRepo.GetAnonymousIP(ip)
	if err != nil {
		return false, err
	}

	return anonIP.IsTorExitNode, nil
}

// GetGeoData는 IP 주소에 대한 종합적인 지리 정보를 조회합니다
func (uc *GeoUseCase) GetGeoData(ipStr string) (*GeoData, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, ErrInvalidIPAddress
	}

	city, cityErr := uc.cityRepo.GetCity(ip)
	country, countryErr := uc.countryRepo.GetCountry(ip)
	asn, asnErr := uc.asnRepo.GetASN(ip)

	// 모든 조회가 실패하면 오류를 반환합니다
	if cityErr != nil && countryErr != nil && asnErr != nil {
		return nil, ErrGeoLookupFailed
	}

	geoData := &GeoData{
		IPAddress: ipStr,
		IsValid:   true,
	}

	// 도시 정보가 있으면 설정합니다
	if cityErr == nil {
		geoData.City = city.City.Names["en"]
		geoData.Latitude = city.Location.Latitude
		geoData.Longitude = city.Location.Longitude
		geoData.TimeZone = city.Location.TimeZone
	}

	// 국가 정보가 있으면 설정합니다
	if countryErr == nil {
		geoData.CountryCode = country.Country.IsoCode
		geoData.CountryName = country.Country.Names["en"]
		geoData.ContinentCode = country.Continent.Code
	} else if cityErr == nil {
		// 도시 정보에서 국가 정보를 가져올 수 있습니다
		geoData.CountryCode = city.Country.IsoCode
		geoData.CountryName = city.Country.Names["en"]
		geoData.ContinentCode = city.Continent.Code
	}

	// ASN 정보가 있으면 설정합니다
	if asnErr == nil {
		geoData.ASN = asn.AutonomousSystemNumber
		geoData.ISP = asn.AutonomousSystemOrganization
	}

	// 익명 IP 정보가 있으면 확인합니다
	if uc.anonymousRepo != nil {
		anonIP, err := uc.anonymousRepo.GetAnonymousIP(ip)
		if err == nil {
			geoData.IsAnonymous = anonIP.IsAnonymous
			geoData.IsAnonymousVPN = anonIP.IsAnonymousVPN
			geoData.IsTorExitNode = anonIP.IsTorExitNode
		}
	}

	return geoData, nil
}

// Close는 사용된 리소스를 해제합니다
func (uc *GeoUseCase) Close() error {
	var cityErr, countryErr, asnErr, anonErr error

	if uc.cityRepo != nil {
		cityErr = uc.cityRepo.Close()
	}

	if uc.countryRepo != nil {
		countryErr = uc.countryRepo.Close()
	}

	if uc.asnRepo != nil {
		asnErr = uc.asnRepo.Close()
	}

	if uc.anonymousRepo != nil {
		anonErr = uc.anonymousRepo.Close()
	}

	// 여러 오류가 발생할 경우 첫 번째 발생한 오류를 반환합니다
	if cityErr != nil {
		return cityErr
	}
	if countryErr != nil {
		return countryErr
	}
	if asnErr != nil {
		return asnErr
	}
	if anonErr != nil {
		return anonErr
	}

	return nil
}

// GeoData는 IP 주소에 대한 종합적인 지리 정보를 담는 구조체입니다
type GeoData struct {
	IPAddress      string  `json:"ip_address"`
	City           string  `json:"city,omitempty"`
	CountryCode    string  `json:"country_code,omitempty"`
	CountryName    string  `json:"country_name,omitempty"`
	ContinentCode  string  `json:"continent_code,omitempty"`
	Latitude       float64 `json:"latitude,omitempty"`
	Longitude      float64 `json:"longitude,omitempty"`
	TimeZone       string  `json:"time_zone,omitempty"`
	ASN            uint    `json:"asn,omitempty"`
	ISP            string  `json:"isp,omitempty"`
	IsValid        bool    `json:"is_valid"`
	IsAnonymous    bool    `json:"is_anonymous,omitempty"`
	IsAnonymousVPN bool    `json:"is_anonymous_vpn,omitempty"`
	IsTorExitNode  bool    `json:"is_tor_exit_node,omitempty"`
}
