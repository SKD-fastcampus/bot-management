package entity

// City는 도시 정보를 담는 구조체입니다
type City struct {
	City      CityInfo      `json:"city"`
	Country   CountryInfo   `json:"country"`
	Continent ContinentInfo `json:"continent"`
	Location  LocationInfo  `json:"location"`
}

// Country는 국가 정보를 담는 구조체입니다
type Country struct {
	Country   CountryInfo   `json:"country"`
	Continent ContinentInfo `json:"continent"`
}

// ASN은 자율 시스템 번호 정보를 담는 구조체입니다
type ASN struct {
	AutonomousSystemNumber       uint   `json:"autonomous_system_number"`
	AutonomousSystemOrganization string `json:"autonomous_system_organization"`
}

// CityInfo는 도시 세부 정보를 담는 구조체입니다
type CityInfo struct {
	GeoNameID uint              `json:"geoname_id"`
	Names     map[string]string `json:"names"`
}

// CountryInfo는 국가 세부 정보를 담는 구조체입니다
type CountryInfo struct {
	GeoNameID         uint              `json:"geoname_id"`
	IsInEuropeanUnion bool              `json:"is_in_european_union"`
	IsoCode           string            `json:"iso_code"`
	Names             map[string]string `json:"names"`
}

// ContinentInfo는 대륙 세부 정보를 담는 구조체입니다
type ContinentInfo struct {
	Code      string            `json:"code"`
	GeoNameID uint              `json:"geoname_id"`
	Names     map[string]string `json:"names"`
}

// LocationInfo는 위치 세부 정보를 담는 구조체입니다
type LocationInfo struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	TimeZone  string  `json:"time_zone"`
}

// Enterprise는 Enterprise 수준의 지리 정보를 담는 구조체입니다
type Enterprise struct {
	City               CityInfo               `json:"city"`
	Country            CountryInfo            `json:"country"`
	Continent          ContinentInfo          `json:"continent"`
	Location           LocationInfo           `json:"location"`
	Traits             EnterpriseTraits       `json:"traits"`
	Postal             PostalInfo             `json:"postal"`
	Subdivisions       []SubdivisionInfo      `json:"subdivisions"`
	RegisteredCountry  CountryInfo            `json:"registered_country"`
	RepresentedCountry RepresentedCountryInfo `json:"represented_country"`
}

// SubdivisionInfo는 지역 구분 정보를 담는 구조체입니다
type SubdivisionInfo struct {
	GeoNameID  uint              `json:"geoname_id"`
	IsoCode    string            `json:"iso_code"`
	Names      map[string]string `json:"names"`
	Confidence uint8             `json:"confidence"`
}

// PostalInfo는 우편 정보를 담는 구조체입니다
type PostalInfo struct {
	Code       string `json:"code"`
	Confidence uint8  `json:"confidence"`
}

// RepresentedCountryInfo는 대표 국가 정보를 담는 구조체입니다
type RepresentedCountryInfo struct {
	GeoNameID         uint              `json:"geoname_id"`
	IsInEuropeanUnion bool              `json:"is_in_european_union"`
	IsoCode           string            `json:"iso_code"`
	Names             map[string]string `json:"names"`
	Type              string            `json:"type"`
}

// EnterpriseTraits는 Enterprise 수준의 특성 정보를 담는 구조체입니다
type EnterpriseTraits struct {
	AutonomousSystemNumber       uint    `json:"autonomous_system_number"`
	AutonomousSystemOrganization string  `json:"autonomous_system_organization"`
	ConnectionType               string  `json:"connection_type"`
	Domain                       string  `json:"domain"`
	ISP                          string  `json:"isp"`
	MobileCountryCode            string  `json:"mobile_country_code"`
	MobileNetworkCode            string  `json:"mobile_network_code"`
	Organization                 string  `json:"organization"`
	UserType                     string  `json:"user_type"`
	StaticIPScore                float64 `json:"static_ip_score"`
	IsAnonymousProxy             bool    `json:"is_anonymous_proxy"`
	IsAnycast                    bool    `json:"is_anycast"`
	IsLegitimateProxy            bool    `json:"is_legitimate_proxy"`
	IsSatelliteProvider          bool    `json:"is_satellite_provider"`
}

// AnonymousIP는 익명 IP 정보를 담는 구조체입니다
type AnonymousIP struct {
	IsAnonymous        bool `json:"is_anonymous"`
	IsAnonymousVPN     bool `json:"is_anonymous_vpn"`
	IsHostingProvider  bool `json:"is_hosting_provider"`
	IsPublicProxy      bool `json:"is_public_proxy"`
	IsResidentialProxy bool `json:"is_residential_proxy"`
	IsTorExitNode      bool `json:"is_tor_exit_node"`
}

// ConnectionType은 연결 유형 정보를 담는 구조체입니다
type ConnectionType struct {
	ConnectionType string `json:"connection_type"`
}

// Domain은 도메인 정보를 담는 구조체입니다
type Domain struct {
	Domain string `json:"domain"`
}

// ISP는 ISP 정보를 담는 구조체입니다
type ISP struct {
	AutonomousSystemNumber       uint   `json:"autonomous_system_number"`
	AutonomousSystemOrganization string `json:"autonomous_system_organization"`
	ISP                          string `json:"isp"`
	MobileCountryCode            string `json:"mobile_country_code"`
	MobileNetworkCode            string `json:"mobile_network_code"`
	Organization                 string `json:"organization"`
}
