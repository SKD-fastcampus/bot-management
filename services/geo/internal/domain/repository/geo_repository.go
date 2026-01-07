package repository

import (
	"net"

	"github.com/wekeepgrowing/semo-backend-monorepo/services/geo/internal/domain/entity"
)

// BaseGeoRepository는 모든 GeoIP 저장소가 공통으로 가지는 메서드를 정의합니다
type BaseGeoRepository interface {
	// Close는 사용한 리소스를 해제합니다
	Close() error
}

// GeoIP2CityRepository는 GeoIP2/GeoLite2 City 데이터베이스용 인터페이스입니다
type GeoIP2CityRepository interface {
	BaseGeoRepository
	// GetCity는 IP 주소를 받아 도시 정보를 반환합니다
	GetCity(ipAddress net.IP) (entity.City, error)
}

// GeoIP2CountryRepository는 GeoIP2/GeoLite2 Country 데이터베이스용 인터페이스입니다
type GeoIP2CountryRepository interface {
	BaseGeoRepository
	// GetCountry는 IP 주소를 받아 국가 정보를 반환합니다
	GetCountry(ipAddress net.IP) (entity.Country, error)
}

// GeoLite2ASNRepository는 GeoLite2 ASN 데이터베이스용 인터페이스입니다
type GeoLite2ASNRepository interface {
	BaseGeoRepository
	// GetASN은 IP 주소를 받아 ASN 정보를 반환합니다
	GetASN(ipAddress net.IP) (entity.ASN, error)
}

// GeoIP2EnterpriseRepository는 GeoIP2 Enterprise 데이터베이스용 인터페이스입니다
type GeoIP2EnterpriseRepository interface {
	BaseGeoRepository
	// GetEnterprise는 IP 주소를 받아 Enterprise 정보를 반환합니다
	GetEnterprise(ipAddress net.IP) (entity.Enterprise, error)
	// GetCity는 Enterprise DB에서도 도시 정보를 제공합니다
	GetCity(ipAddress net.IP) (entity.City, error)
	// GetCountry는 Enterprise DB에서도 국가 정보를 제공합니다
	GetCountry(ipAddress net.IP) (entity.Country, error)
}

// GeoIP2AnonymousIPRepository는 GeoIP2 Anonymous IP 데이터베이스용 인터페이스입니다
type GeoIP2AnonymousIPRepository interface {
	BaseGeoRepository
	// GetAnonymousIP는 IP 주소를 받아 익명 IP 정보를 반환합니다
	GetAnonymousIP(ipAddress net.IP) (entity.AnonymousIP, error)
}

// GeoIP2ISPRepository는 GeoIP2 ISP 데이터베이스용 인터페이스입니다
type GeoIP2ISPRepository interface {
	BaseGeoRepository
	// GetISP는 IP 주소를 받아 ISP 정보를 반환합니다
	GetISP(ipAddress net.IP) (entity.ISP, error)
	// GetASN은 ISP DB에서도 ASN 정보를 제공합니다
	GetASN(ipAddress net.IP) (entity.ASN, error)
}

// GeoIP2DomainRepository는 GeoIP2 Domain 데이터베이스용 인터페이스입니다
type GeoIP2DomainRepository interface {
	BaseGeoRepository
	// GetDomain은 IP 주소를 받아 도메인 정보를 반환합니다
	GetDomain(ipAddress net.IP) (entity.Domain, error)
}

// GeoIP2ConnectionTypeRepository는 GeoIP2 Connection Type 데이터베이스용 인터페이스입니다
type GeoIP2ConnectionTypeRepository interface {
	BaseGeoRepository
	// GetConnectionType은 IP 주소를 받아 연결 유형 정보를 반환합니다
	GetConnectionType(ipAddress net.IP) (entity.ConnectionType, error)
}

// GeoLite2Repository는 GeoLite2 데이터베이스를 통합해서 사용하는 인터페이스입니다
type GeoLite2Repository interface {
	GeoIP2CityRepository
	GeoIP2CountryRepository
	GeoLite2ASNRepository
}

// GeoIP2FullRepository는 모든 GeoIP2 데이터베이스를 통합해서 사용하는 인터페이스입니다
type GeoIP2FullRepository interface {
	GeoIP2EnterpriseRepository
	GeoIP2AnonymousIPRepository
	GeoIP2ISPRepository
	GeoIP2DomainRepository
	GeoIP2ConnectionTypeRepository
}
