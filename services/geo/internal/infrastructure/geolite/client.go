package geolite

import (
	"fmt"
	"net"

	"github.com/oschwald/maxminddb-golang"
)

type databaseType int

const (
	isAnonymousIP = 1 << iota
	isASN
	isCity
	isConnectionType
	isCountry
	isDomain
	isEnterprise
	isISP
)

// Reader는 maxminddb.Reader 구조체를 보유합니다. Open 및 FromBytes 함수를 사용하여 생성할 수 있습니다.
type Reader struct {
	mmdbReader   *maxminddb.Reader
	databaseType databaseType
}

// InvalidMethodError는 지원하지 않는 데이터베이스에서 조회 메서드가 호출될 때 반환됩니다.
// 예를 들어, City 데이터베이스에서 ISP 메서드를 호출하는 경우입니다.
type InvalidMethodError struct {
	Method       string
	DatabaseType string
}

func (e InvalidMethodError) Error() string {
	return fmt.Sprintf(`geoip2: the %s method does not support the %s database`,
		e.Method, e.DatabaseType)
}

// UnknownDatabaseTypeError는 알 수 없는 데이터베이스 유형이 열릴 때 반환됩니다.
type UnknownDatabaseTypeError struct {
	DatabaseType string
}

func (e UnknownDatabaseTypeError) Error() string {
	return fmt.Sprintf(`geoip2: reader does not support the %q database type`,
		e.DatabaseType)
}

// Open은 파일 경로를 인자로 받아 Reader 구조체 또는 오류를 반환합니다.
// 데이터베이스 파일은 메모리 맵을 사용하여 열립니다. Reader 객체의 Close 메서드를 사용하여
// 시스템에 리소스를 반환합니다.
func Open(file string) (*Reader, error) {
	reader, err := maxminddb.Open(file)
	if err != nil {
		return nil, err
	}
	dbType, err := getDBType(reader)
	return &Reader{reader, dbType}, err
}

// FromBytes는 GeoIP2/GeoLite2 데이터베이스 파일에 해당하는 바이트 슬라이스를 받아
// Reader 구조체 또는 오류를 반환합니다. 바이트 슬라이스는 직접 사용되므로,
// 데이터베이스를 연 후에 변경하면 데이터베이스에서 읽는 동안 오류가 발생할 수 있습니다.
func FromBytes(bytes []byte) (*Reader, error) {
	reader, err := maxminddb.FromBytes(bytes)
	if err != nil {
		return nil, err
	}
	dbType, err := getDBType(reader)
	return &Reader{reader, dbType}, err
}

func getDBType(reader *maxminddb.Reader) (databaseType, error) {
	switch reader.Metadata.DatabaseType {
	case "GeoIP2-Anonymous-IP":
		return isAnonymousIP, nil
	case "DBIP-ASN-Lite (compat=GeoLite2-ASN)",
		"GeoLite2-ASN":
		return isASN, nil
	// 이전 호환성을 위해 Country에서 City 조회를 허용합니다
	case "DBIP-City-Lite",
		"DBIP-Country-Lite",
		"DBIP-Country",
		"DBIP-Location (compat=City)",
		"GeoLite2-City",
		"GeoIP2-City",
		"GeoIP2-City-Africa",
		"GeoIP2-City-Asia-Pacific",
		"GeoIP2-City-Europe",
		"GeoIP2-City-North-America",
		"GeoIP2-City-South-America",
		"GeoIP2-Precision-City",
		"GeoLite2-Country",
		"GeoIP2-Country":
		return isCity | isCountry, nil
	case "GeoIP2-Connection-Type":
		return isConnectionType, nil
	case "GeoIP2-Domain":
		return isDomain, nil
	case "DBIP-ISP (compat=Enterprise)",
		"DBIP-Location-ISP (compat=Enterprise)",
		"GeoIP2-Enterprise":
		return isEnterprise | isCity | isCountry, nil
	case "GeoIP2-ISP", "GeoIP2-Precision-ISP":
		return isISP | isASN, nil
	default:
		return 0, UnknownDatabaseTypeError{reader.Metadata.DatabaseType}
	}
}

// Enterprise는 net.IP 구조체로 IP 주소를 받아 Enterprise 구조체 및/또는 오류를 반환합니다.
// 이는 GeoIP2 Enterprise 데이터베이스와 함께 사용하기 위한 것입니다.
func (r *Reader) Enterprise(ipAddress net.IP) (*Enterprise, error) {
	if isEnterprise&r.databaseType == 0 {
		return nil, InvalidMethodError{"Enterprise", r.Metadata().DatabaseType}
	}
	var enterprise Enterprise
	err := r.mmdbReader.Lookup(ipAddress, &enterprise)
	return &enterprise, err
}

// City는 net.IP 구조체로 IP 주소를 받아 City 구조체 및/또는 오류를 반환합니다.
// 다른 데이터베이스와 함께 사용할 수 있지만, 이 메서드는 일반적으로 GeoIP2 또는 GeoLite2 City
// 데이터베이스와 함께 사용해야 합니다.
func (r *Reader) City(ipAddress net.IP) (*City, error) {
	if isCity&r.databaseType == 0 {
		return nil, InvalidMethodError{"City", r.Metadata().DatabaseType}
	}
	var city City
	err := r.mmdbReader.Lookup(ipAddress, &city)
	return &city, err
}

// Country는 net.IP 구조체로 IP 주소를 받아 Country 구조체 및/또는 오류를 반환합니다.
// 다른 데이터베이스와 함께 사용할 수 있지만, 이 메서드는 일반적으로 GeoIP2 또는 GeoLite2 Country
// 데이터베이스와 함께 사용해야 합니다.
func (r *Reader) Country(ipAddress net.IP) (*Country, error) {
	if isCountry&r.databaseType == 0 {
		return nil, InvalidMethodError{"Country", r.Metadata().DatabaseType}
	}
	var country Country
	err := r.mmdbReader.Lookup(ipAddress, &country)
	return &country, err
}

// AnonymousIP는 net.IP 구조체로 IP 주소를 받아 AnonymousIP 구조체 및/또는 오류를 반환합니다.
func (r *Reader) AnonymousIP(ipAddress net.IP) (*AnonymousIP, error) {
	if isAnonymousIP&r.databaseType == 0 {
		return nil, InvalidMethodError{"AnonymousIP", r.Metadata().DatabaseType}
	}
	var anonIP AnonymousIP
	err := r.mmdbReader.Lookup(ipAddress, &anonIP)
	return &anonIP, err
}

// ASN은 net.IP 구조체로 IP 주소를 받아 ASN 구조체 및/또는 오류를 반환합니다.
func (r *Reader) ASN(ipAddress net.IP) (*ASN, error) {
	if isASN&r.databaseType == 0 {
		return nil, InvalidMethodError{"ASN", r.Metadata().DatabaseType}
	}
	var val ASN
	err := r.mmdbReader.Lookup(ipAddress, &val)
	return &val, err
}

// ConnectionType은 net.IP 구조체로 IP 주소를 받아 ConnectionType 구조체 및/또는 오류를 반환합니다.
func (r *Reader) ConnectionType(ipAddress net.IP) (*ConnectionType, error) {
	if isConnectionType&r.databaseType == 0 {
		return nil, InvalidMethodError{"ConnectionType", r.Metadata().DatabaseType}
	}
	var val ConnectionType
	err := r.mmdbReader.Lookup(ipAddress, &val)
	return &val, err
}

// Domain은 net.IP 구조체로 IP 주소를 받아 Domain 구조체 및/또는 오류를 반환합니다.
func (r *Reader) Domain(ipAddress net.IP) (*Domain, error) {
	if isDomain&r.databaseType == 0 {
		return nil, InvalidMethodError{"Domain", r.Metadata().DatabaseType}
	}
	var val Domain
	err := r.mmdbReader.Lookup(ipAddress, &val)
	return &val, err
}

// ISP는 net.IP 구조체로 IP 주소를 받아 ISP 구조체 및/또는 오류를 반환합니다.
func (r *Reader) ISP(ipAddress net.IP) (*ISP, error) {
	if isISP&r.databaseType == 0 {
		return nil, InvalidMethodError{"ISP", r.Metadata().DatabaseType}
	}
	var val ISP
	err := r.mmdbReader.Lookup(ipAddress, &val)
	return &val, err
}

// Metadata는 인자를 받지 않고 Reader가 사용하는 MaxMind 데이터베이스에 대한 메타데이터가 포함된
// 구조체를 반환합니다.
func (r *Reader) Metadata() maxminddb.Metadata {
	return r.mmdbReader.Metadata
}

// Close는 데이터베이스 파일을 가상 메모리에서 해제하고 시스템에 리소스를 반환합니다.
func (r *Reader) Close() error {
	return r.mmdbReader.Close()
}
