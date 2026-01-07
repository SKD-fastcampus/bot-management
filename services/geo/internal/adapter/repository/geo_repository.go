package repository

import (
	"net"

	"github.com/wekeepgrowing/semo-backend-monorepo/services/geo/internal/domain/entity"
	"github.com/wekeepgrowing/semo-backend-monorepo/services/geo/internal/domain/repository"
	"github.com/wekeepgrowing/semo-backend-monorepo/services/geo/internal/infrastructure/geolite"
)

// 기본 구현체 - 모든 구현체의 기초가 됩니다
type baseGeoRepository struct {
	reader *geolite.Reader
}

// Close는 리더와 관련된 리소스를 해제합니다
func (g *baseGeoRepository) Close() error {
	return g.reader.Close()
}

// GeoIP2City는 GeoIP2/GeoLite2 City 데이터베이스 리포지토리 구현체입니다
type GeoIP2City struct {
	baseGeoRepository
}

// NewGeoIP2CityRepository는 City 데이터베이스 리포지토리를 생성합니다
func NewGeoIP2CityRepository(dbPath string) (repository.GeoIP2CityRepository, error) {
	reader, err := geolite.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &GeoIP2City{baseGeoRepository{reader}}, nil
}

// GetCity는 IP 주소에 해당하는 도시 정보를 반환합니다
func (g *GeoIP2City) GetCity(ipAddress net.IP) (entity.City, error) {
	cityData, err := g.reader.City(ipAddress)
	if err != nil {
		return entity.City{}, err
	}

	// geolite.City에서 entity.City로 매핑
	city := entity.City{
		City: entity.CityInfo{
			GeoNameID: cityData.City.GeoNameID,
			Names:     cityData.City.Names,
		},
		Country: entity.CountryInfo{
			GeoNameID:         cityData.Country.GeoNameID,
			IsInEuropeanUnion: cityData.Country.IsInEuropeanUnion,
			IsoCode:           cityData.Country.IsoCode,
			Names:             cityData.Country.Names,
		},
		Continent: entity.ContinentInfo{
			Code:      cityData.Continent.Code,
			GeoNameID: cityData.Continent.GeoNameID,
			Names:     cityData.Continent.Names,
		},
		Location: entity.LocationInfo{
			Latitude:  cityData.Location.Latitude,
			Longitude: cityData.Location.Longitude,
			TimeZone:  cityData.Location.TimeZone,
		},
	}

	return city, nil
}

// GeoIP2Country는 GeoIP2/GeoLite2 Country 데이터베이스 리포지토리 구현체입니다
type GeoIP2Country struct {
	baseGeoRepository
}

// NewGeoIP2CountryRepository는 Country 데이터베이스 리포지토리를 생성합니다
func NewGeoIP2CountryRepository(dbPath string) (repository.GeoIP2CountryRepository, error) {
	reader, err := geolite.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &GeoIP2Country{baseGeoRepository{reader}}, nil
}

// GetCountry는 IP 주소에 해당하는 국가 정보를 반환합니다
func (g *GeoIP2Country) GetCountry(ipAddress net.IP) (entity.Country, error) {
	countryData, err := g.reader.Country(ipAddress)
	if err != nil {
		return entity.Country{}, err
	}

	// geolite.Country에서 entity.Country로 매핑
	country := entity.Country{
		Country: entity.CountryInfo{
			GeoNameID:         countryData.Country.GeoNameID,
			IsInEuropeanUnion: countryData.Country.IsInEuropeanUnion,
			IsoCode:           countryData.Country.IsoCode,
			Names:             countryData.Country.Names,
		},
		Continent: entity.ContinentInfo{
			Code:      countryData.Continent.Code,
			GeoNameID: countryData.Continent.GeoNameID,
			Names:     countryData.Continent.Names,
		},
	}

	return country, nil
}

// GeoLite2ASN은 GeoLite2 ASN 데이터베이스 리포지토리 구현체입니다
type GeoLite2ASN struct {
	baseGeoRepository
}

// NewGeoLite2ASNRepository는 ASN 데이터베이스 리포지토리를 생성합니다
func NewGeoLite2ASNRepository(dbPath string) (repository.GeoLite2ASNRepository, error) {
	reader, err := geolite.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &GeoLite2ASN{baseGeoRepository{reader}}, nil
}

// GetASN은 IP 주소에 해당하는 ASN 정보를 반환합니다
func (g *GeoLite2ASN) GetASN(ipAddress net.IP) (entity.ASN, error) {
	asnData, err := g.reader.ASN(ipAddress)
	if err != nil {
		return entity.ASN{}, err
	}

	// geolite.ASN에서 entity.ASN으로 매핑
	asn := entity.ASN{
		AutonomousSystemNumber:       asnData.AutonomousSystemNumber,
		AutonomousSystemOrganization: asnData.AutonomousSystemOrganization,
	}

	return asn, nil
}

// GeoIP2Enterprise는 GeoIP2 Enterprise 데이터베이스 리포지토리 구현체입니다
type GeoIP2Enterprise struct {
	baseGeoRepository
}

// NewGeoIP2EnterpriseRepository는 Enterprise 데이터베이스 리포지토리를 생성합니다
func NewGeoIP2EnterpriseRepository(dbPath string) (repository.GeoIP2EnterpriseRepository, error) {
	reader, err := geolite.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &GeoIP2Enterprise{baseGeoRepository{reader}}, nil
}

// GetEnterprise는 IP 주소에 해당하는 Enterprise 정보를 반환합니다
func (g *GeoIP2Enterprise) GetEnterprise(ipAddress net.IP) (entity.Enterprise, error) {
	enterpriseData, err := g.reader.Enterprise(ipAddress)
	if err != nil {
		return entity.Enterprise{}, err
	}

	// 세부 지역 정보 변환
	subdivisions := make([]entity.SubdivisionInfo, len(enterpriseData.Subdivisions))
	for i, subdivision := range enterpriseData.Subdivisions {
		subdivisions[i] = entity.SubdivisionInfo{
			GeoNameID:  subdivision.GeoNameID,
			IsoCode:    subdivision.IsoCode,
			Names:      subdivision.Names,
			Confidence: subdivision.Confidence,
		}
	}

	// geolite.Enterprise에서 entity.Enterprise로 매핑
	enterprise := entity.Enterprise{
		City: entity.CityInfo{
			GeoNameID: enterpriseData.City.GeoNameID,
			Names:     enterpriseData.City.Names,
		},
		Country: entity.CountryInfo{
			GeoNameID:         enterpriseData.Country.GeoNameID,
			IsInEuropeanUnion: enterpriseData.Country.IsInEuropeanUnion,
			IsoCode:           enterpriseData.Country.IsoCode,
			Names:             enterpriseData.Country.Names,
		},
		Continent: entity.ContinentInfo{
			Code:      enterpriseData.Continent.Code,
			GeoNameID: enterpriseData.Continent.GeoNameID,
			Names:     enterpriseData.Continent.Names,
		},
		Location: entity.LocationInfo{
			Latitude:  enterpriseData.Location.Latitude,
			Longitude: enterpriseData.Location.Longitude,
			TimeZone:  enterpriseData.Location.TimeZone,
		},
		Postal: entity.PostalInfo{
			Code:       enterpriseData.Postal.Code,
			Confidence: enterpriseData.Postal.Confidence,
		},
		Subdivisions: subdivisions,
		RegisteredCountry: entity.CountryInfo{
			GeoNameID:         enterpriseData.RegisteredCountry.GeoNameID,
			IsInEuropeanUnion: enterpriseData.RegisteredCountry.IsInEuropeanUnion,
			IsoCode:           enterpriseData.RegisteredCountry.IsoCode,
			Names:             enterpriseData.RegisteredCountry.Names,
		},
		RepresentedCountry: entity.RepresentedCountryInfo{
			GeoNameID:         enterpriseData.RepresentedCountry.GeoNameID,
			IsInEuropeanUnion: enterpriseData.RepresentedCountry.IsInEuropeanUnion,
			IsoCode:           enterpriseData.RepresentedCountry.IsoCode,
			Names:             enterpriseData.RepresentedCountry.Names,
			Type:              enterpriseData.RepresentedCountry.Type,
		},
		Traits: entity.EnterpriseTraits{
			AutonomousSystemNumber:       enterpriseData.Traits.AutonomousSystemNumber,
			AutonomousSystemOrganization: enterpriseData.Traits.AutonomousSystemOrganization,
			ConnectionType:               enterpriseData.Traits.ConnectionType,
			Domain:                       enterpriseData.Traits.Domain,
			ISP:                          enterpriseData.Traits.ISP,
			MobileCountryCode:            enterpriseData.Traits.MobileCountryCode,
			MobileNetworkCode:            enterpriseData.Traits.MobileNetworkCode,
			Organization:                 enterpriseData.Traits.Organization,
			UserType:                     enterpriseData.Traits.UserType,
			StaticIPScore:                enterpriseData.Traits.StaticIPScore,
			IsAnonymousProxy:             enterpriseData.Traits.IsAnonymousProxy,
			IsAnycast:                    enterpriseData.Traits.IsAnycast,
			IsLegitimateProxy:            enterpriseData.Traits.IsLegitimateProxy,
			IsSatelliteProvider:          enterpriseData.Traits.IsSatelliteProvider,
		},
	}

	return enterprise, nil
}

// GetCity는 Enterprise DB에서도 도시 정보를 제공합니다
func (g *GeoIP2Enterprise) GetCity(ipAddress net.IP) (entity.City, error) {
	cityData, err := g.reader.City(ipAddress)
	if err != nil {
		return entity.City{}, err
	}

	// geolite.City에서 entity.City로 매핑
	city := entity.City{
		City: entity.CityInfo{
			GeoNameID: cityData.City.GeoNameID,
			Names:     cityData.City.Names,
		},
		Country: entity.CountryInfo{
			GeoNameID:         cityData.Country.GeoNameID,
			IsInEuropeanUnion: cityData.Country.IsInEuropeanUnion,
			IsoCode:           cityData.Country.IsoCode,
			Names:             cityData.Country.Names,
		},
		Continent: entity.ContinentInfo{
			Code:      cityData.Continent.Code,
			GeoNameID: cityData.Continent.GeoNameID,
			Names:     cityData.Continent.Names,
		},
		Location: entity.LocationInfo{
			Latitude:  cityData.Location.Latitude,
			Longitude: cityData.Location.Longitude,
			TimeZone:  cityData.Location.TimeZone,
		},
	}

	return city, nil
}

// GetCountry는 Enterprise DB에서도 국가 정보를 제공합니다
func (g *GeoIP2Enterprise) GetCountry(ipAddress net.IP) (entity.Country, error) {
	countryData, err := g.reader.Country(ipAddress)
	if err != nil {
		return entity.Country{}, err
	}

	// geolite.Country에서 entity.Country로 매핑
	country := entity.Country{
		Country: entity.CountryInfo{
			GeoNameID:         countryData.Country.GeoNameID,
			IsInEuropeanUnion: countryData.Country.IsInEuropeanUnion,
			IsoCode:           countryData.Country.IsoCode,
			Names:             countryData.Country.Names,
		},
		Continent: entity.ContinentInfo{
			Code:      countryData.Continent.Code,
			GeoNameID: countryData.Continent.GeoNameID,
			Names:     countryData.Continent.Names,
		},
	}

	return country, nil
}

// GeoIP2AnonymousIP는 GeoIP2 Anonymous IP 데이터베이스 리포지토리 구현체입니다
type GeoIP2AnonymousIP struct {
	baseGeoRepository
}

// NewGeoIP2AnonymousIPRepository는 Anonymous IP 데이터베이스 리포지토리를 생성합니다
func NewGeoIP2AnonymousIPRepository(dbPath string) (repository.GeoIP2AnonymousIPRepository, error) {
	reader, err := geolite.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &GeoIP2AnonymousIP{baseGeoRepository{reader}}, nil
}

// GetAnonymousIP는 IP 주소에 해당하는 익명 IP 정보를 반환합니다
func (g *GeoIP2AnonymousIP) GetAnonymousIP(ipAddress net.IP) (entity.AnonymousIP, error) {
	anonymousIPData, err := g.reader.AnonymousIP(ipAddress)
	if err != nil {
		return entity.AnonymousIP{}, err
	}

	// geolite.AnonymousIP에서 entity.AnonymousIP로 매핑
	anonymousIP := entity.AnonymousIP{
		IsAnonymous:        anonymousIPData.IsAnonymous,
		IsAnonymousVPN:     anonymousIPData.IsAnonymousVPN,
		IsHostingProvider:  anonymousIPData.IsHostingProvider,
		IsPublicProxy:      anonymousIPData.IsPublicProxy,
		IsResidentialProxy: anonymousIPData.IsResidentialProxy,
		IsTorExitNode:      anonymousIPData.IsTorExitNode,
	}

	return anonymousIP, nil
}

// GeoIP2ISP는 GeoIP2 ISP 데이터베이스 리포지토리 구현체입니다
type GeoIP2ISP struct {
	baseGeoRepository
}

// NewGeoIP2ISPRepository는 ISP 데이터베이스 리포지토리를 생성합니다
func NewGeoIP2ISPRepository(dbPath string) (repository.GeoIP2ISPRepository, error) {
	reader, err := geolite.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &GeoIP2ISP{baseGeoRepository{reader}}, nil
}

// GetISP는 IP 주소에 해당하는 ISP 정보를 반환합니다
func (g *GeoIP2ISP) GetISP(ipAddress net.IP) (entity.ISP, error) {
	ispData, err := g.reader.ISP(ipAddress)
	if err != nil {
		return entity.ISP{}, err
	}

	// geolite.ISP에서 entity.ISP로 매핑
	isp := entity.ISP{
		AutonomousSystemNumber:       ispData.AutonomousSystemNumber,
		AutonomousSystemOrganization: ispData.AutonomousSystemOrganization,
		ISP:                          ispData.ISP,
		MobileCountryCode:            ispData.MobileCountryCode,
		MobileNetworkCode:            ispData.MobileNetworkCode,
		Organization:                 ispData.Organization,
	}

	return isp, nil
}

// GetASN은 ISP DB에서도 ASN 정보를 제공합니다
func (g *GeoIP2ISP) GetASN(ipAddress net.IP) (entity.ASN, error) {
	ispData, err := g.reader.ISP(ipAddress)
	if err != nil {
		return entity.ASN{}, err
	}

	// geolite.ISP에서 entity.ASN으로 매핑
	asn := entity.ASN{
		AutonomousSystemNumber:       ispData.AutonomousSystemNumber,
		AutonomousSystemOrganization: ispData.AutonomousSystemOrganization,
	}

	return asn, nil
}

// GeoIP2Domain은 GeoIP2 Domain 데이터베이스 리포지토리 구현체입니다
type GeoIP2Domain struct {
	baseGeoRepository
}

// NewGeoIP2DomainRepository는 Domain 데이터베이스 리포지토리를 생성합니다
func NewGeoIP2DomainRepository(dbPath string) (repository.GeoIP2DomainRepository, error) {
	reader, err := geolite.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &GeoIP2Domain{baseGeoRepository{reader}}, nil
}

// GetDomain은 IP 주소에 해당하는 도메인 정보를 반환합니다
func (g *GeoIP2Domain) GetDomain(ipAddress net.IP) (entity.Domain, error) {
	domainData, err := g.reader.Domain(ipAddress)
	if err != nil {
		return entity.Domain{}, err
	}

	// geolite.Domain에서 entity.Domain으로 매핑
	domain := entity.Domain{
		Domain: domainData.Domain,
	}

	return domain, nil
}

// GeoIP2ConnectionType은 GeoIP2 Connection Type 데이터베이스 리포지토리 구현체입니다
type GeoIP2ConnectionType struct {
	baseGeoRepository
}

// NewGeoIP2ConnectionTypeRepository는 Connection Type 데이터베이스 리포지토리를 생성합니다
func NewGeoIP2ConnectionTypeRepository(dbPath string) (repository.GeoIP2ConnectionTypeRepository, error) {
	reader, err := geolite.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &GeoIP2ConnectionType{baseGeoRepository{reader}}, nil
}

// GetConnectionType은 IP 주소에 해당하는 연결 유형 정보를 반환합니다
func (g *GeoIP2ConnectionType) GetConnectionType(ipAddress net.IP) (entity.ConnectionType, error) {
	connectionTypeData, err := g.reader.ConnectionType(ipAddress)
	if err != nil {
		return entity.ConnectionType{}, err
	}

	// geolite.ConnectionType에서 entity.ConnectionType으로 매핑
	connectionType := entity.ConnectionType{
		ConnectionType: connectionTypeData.ConnectionType,
	}

	return connectionType, nil
}

// GeoLite2All은 GeoLite2 데이터베이스를 모두 사용하는 리포지토리 구현체입니다
type GeoLite2All struct {
	GeoIP2City
	GeoIP2Country
	GeoLite2ASN
}

// NewGeoLite2Repository는 GeoLite2 통합 리포지토리를 생성합니다
func NewGeoLite2Repository(cityDbPath, countryDbPath, asnDbPath string) (repository.GeoLite2Repository, error) {
	cityReader, err := geolite.Open(cityDbPath)
	if err != nil {
		return nil, err
	}

	countryReader, err := geolite.Open(countryDbPath)
	if err != nil {
		cityReader.Close()
		return nil, err
	}

	asnReader, err := geolite.Open(asnDbPath)
	if err != nil {
		cityReader.Close()
		countryReader.Close()
		return nil, err
	}

	return &GeoLite2All{
		GeoIP2City:    GeoIP2City{baseGeoRepository{cityReader}},
		GeoIP2Country: GeoIP2Country{baseGeoRepository{countryReader}},
		GeoLite2ASN:   GeoLite2ASN{baseGeoRepository{asnReader}},
	}, nil
}

// Close는 모든 리더의 리소스를 해제합니다
func (g *GeoLite2All) Close() error {
	g.GeoIP2City.Close()
	g.GeoIP2Country.Close()
	return g.GeoLite2ASN.Close()
}
