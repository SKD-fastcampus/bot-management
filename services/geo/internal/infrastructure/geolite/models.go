package geolite

// City 구조체는 GeoIP2/GeoLite2 City 데이터베이스의 데이터에 해당합니다.
type City struct {
	City struct {
		Names     map[string]string `maxminddb:"names"`
		GeoNameID uint              `maxminddb:"geoname_id"`
	} `maxminddb:"city"`
	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`
	Continent struct {
		Names     map[string]string `maxminddb:"names"`
		Code      string            `maxminddb:"code"`
		GeoNameID uint              `maxminddb:"geoname_id"`
	} `maxminddb:"continent"`
	Subdivisions []struct {
		Names     map[string]string `maxminddb:"names"`
		IsoCode   string            `maxminddb:"iso_code"`
		GeoNameID uint              `maxminddb:"geoname_id"`
	} `maxminddb:"subdivisions"`
	RepresentedCountry struct {
		Names             map[string]string `maxminddb:"names"`
		IsoCode           string            `maxminddb:"iso_code"`
		Type              string            `maxminddb:"type"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"represented_country"`
	Country struct {
		Names             map[string]string `maxminddb:"names"`
		IsoCode           string            `maxminddb:"iso_code"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"country"`
	RegisteredCountry struct {
		Names             map[string]string `maxminddb:"names"`
		IsoCode           string            `maxminddb:"iso_code"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"registered_country"`
	Location struct {
		TimeZone       string  `maxminddb:"time_zone"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		MetroCode      uint    `maxminddb:"metro_code"`
		AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
	} `maxminddb:"location"`
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy"`
		IsAnycast           bool `maxminddb:"is_anycast"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider"`
	} `maxminddb:"traits"`
}

// Country 구조체는 GeoIP2/GeoLite2 Country 데이터베이스의 데이터에 해당합니다.
type Country struct {
	Continent struct {
		Names     map[string]string `maxminddb:"names"`
		Code      string            `maxminddb:"code"`
		GeoNameID uint              `maxminddb:"geoname_id"`
	} `maxminddb:"continent"`
	Country struct {
		Names             map[string]string `maxminddb:"names"`
		IsoCode           string            `maxminddb:"iso_code"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"country"`
	RegisteredCountry struct {
		Names             map[string]string `maxminddb:"names"`
		IsoCode           string            `maxminddb:"iso_code"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"registered_country"`
	RepresentedCountry struct {
		Names             map[string]string `maxminddb:"names"`
		IsoCode           string            `maxminddb:"iso_code"`
		Type              string            `maxminddb:"type"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"represented_country"`
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy"`
		IsAnycast           bool `maxminddb:"is_anycast"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider"`
	} `maxminddb:"traits"`
}

// ASN 구조체는 GeoLite2 ASN 데이터베이스의 데이터에 해당합니다.
type ASN struct {
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization"`
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number"`
}

// Enterprise 구조체는 GeoIP2 Enterprise 데이터베이스의 데이터에 해당합니다.
type Enterprise struct {
	Continent struct {
		Names     map[string]string `maxminddb:"names"`
		Code      string            `maxminddb:"code"`
		GeoNameID uint              `maxminddb:"geoname_id"`
	} `maxminddb:"continent"`
	City struct {
		Names      map[string]string `maxminddb:"names"`
		GeoNameID  uint              `maxminddb:"geoname_id"`
		Confidence uint8             `maxminddb:"confidence"`
	} `maxminddb:"city"`
	Postal struct {
		Code       string `maxminddb:"code"`
		Confidence uint8  `maxminddb:"confidence"`
	} `maxminddb:"postal"`
	Subdivisions []struct {
		Names      map[string]string `maxminddb:"names"`
		IsoCode    string            `maxminddb:"iso_code"`
		GeoNameID  uint              `maxminddb:"geoname_id"`
		Confidence uint8             `maxminddb:"confidence"`
	} `maxminddb:"subdivisions"`
	RepresentedCountry struct {
		Names             map[string]string `maxminddb:"names"`
		IsoCode           string            `maxminddb:"iso_code"`
		Type              string            `maxminddb:"type"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"represented_country"`
	Country struct {
		Names             map[string]string `maxminddb:"names"`
		IsoCode           string            `maxminddb:"iso_code"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		Confidence        uint8             `maxminddb:"confidence"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"country"`
	RegisteredCountry struct {
		Names             map[string]string `maxminddb:"names"`
		IsoCode           string            `maxminddb:"iso_code"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		Confidence        uint8             `maxminddb:"confidence"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"registered_country"`
	Traits struct {
		AutonomousSystemOrganization string  `maxminddb:"autonomous_system_organization"`
		ConnectionType               string  `maxminddb:"connection_type"`
		Domain                       string  `maxminddb:"domain"`
		ISP                          string  `maxminddb:"isp"`
		MobileCountryCode            string  `maxminddb:"mobile_country_code"`
		MobileNetworkCode            string  `maxminddb:"mobile_network_code"`
		Organization                 string  `maxminddb:"organization"`
		UserType                     string  `maxminddb:"user_type"`
		AutonomousSystemNumber       uint    `maxminddb:"autonomous_system_number"`
		StaticIPScore                float64 `maxminddb:"static_ip_score"`
		IsAnonymousProxy             bool    `maxminddb:"is_anonymous_proxy"`
		IsAnycast                    bool    `maxminddb:"is_anycast"`
		IsLegitimateProxy            bool    `maxminddb:"is_legitimate_proxy"`
		IsSatelliteProvider          bool    `maxminddb:"is_satellite_provider"`
	} `maxminddb:"traits"`
	Location struct {
		TimeZone       string  `maxminddb:"time_zone"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		MetroCode      uint    `maxminddb:"metro_code"`
		AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
	} `maxminddb:"location"`
}

// AnonymousIP 구조체는 GeoIP2 Anonymous IP 데이터베이스의 데이터에 해당합니다.
type AnonymousIP struct {
	IsAnonymous        bool `maxminddb:"is_anonymous"`
	IsAnonymousVPN     bool `maxminddb:"is_anonymous_vpn"`
	IsHostingProvider  bool `maxminddb:"is_hosting_provider"`
	IsPublicProxy      bool `maxminddb:"is_public_proxy"`
	IsResidentialProxy bool `maxminddb:"is_residential_proxy"`
	IsTorExitNode      bool `maxminddb:"is_tor_exit_node"`
}

// ConnectionType 구조체는 GeoIP2 Connection-Type 데이터베이스의 데이터에 해당합니다.
type ConnectionType struct {
	ConnectionType string `maxminddb:"connection_type"`
}

// Domain 구조체는 GeoIP2 Domain 데이터베이스의 데이터에 해당합니다.
type Domain struct {
	Domain string `maxminddb:"domain"`
}

// ISP 구조체는 GeoIP2 ISP 데이터베이스의 데이터에 해당합니다.
type ISP struct {
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization"`
	ISP                          string `maxminddb:"isp"`
	MobileCountryCode            string `maxminddb:"mobile_country_code"`
	MobileNetworkCode            string `maxminddb:"mobile_network_code"`
	Organization                 string `maxminddb:"organization"`
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number"`
}
