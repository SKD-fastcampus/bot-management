package usecase_test

import (
	"fmt"
	"log"
	"net"

	"github.com/SKD-fastcampus/bot-management/services/geo/internal/adapter/repository"
	"github.com/SKD-fastcampus/bot-management/services/geo/internal/domain/entity"
	"github.com/SKD-fastcampus/bot-management/services/geo/internal/usecase"
)

// 이 파일은 단순한 예제입니다. 실제 데이터베이스 파일 경로를 적절히 변경해야 합니다.

// 테스트에 사용할 IP 주소
var testIP = "211.177.139.105"

func Example_cityLookup() {
	// City DB 리포지토리 생성
	cityRepo, err := repository.NewGeoIP2CityRepository("/Users/k2zoo/Documents/growingup/semo/semo-backend-monorepo/services/geo/data/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatalf("City DB 열기 실패: %v", err)
	}
	defer cityRepo.Close()

	// UseCase 생성
	geoUseCase := usecase.NewGeoUseCase(cityRepo, nil, nil, nil)
	defer geoUseCase.Close()

	// IP 주소 정보 조회
	city, err := geoUseCase.GetCityInfo(testIP)
	if err != nil {
		log.Fatalf("도시 정보 조회 실패: %v", err)
	}

	// 결과 출력 - 모든 정보 표시
	printCityDetails(city, testIP)

	// Output:
	// (실제 출력은 사용하는 데이터베이스에 따라 다를 수 있습니다)
}

func Example_fullLookup() {
	// GeoLite2 통합 리포지토리 생성
	geoRepo, err := repository.NewGeoLite2Repository(
		"/Users/k2zoo/Documents/growingup/semo/semo-backend-monorepo/services/geo/data/GeoLite2-City.mmdb",
		"/Users/k2zoo/Documents/growingup/semo/semo-backend-monorepo/services/geo/data/GeoLite2-Country.mmdb",
		"/Users/k2zoo/Documents/growingup/semo/semo-backend-monorepo/services/geo/data/GeoLite2-ASN.mmdb",
	)
	if err != nil {
		log.Fatalf("GeoLite2 DB 열기 실패: %v", err)
	}
	defer geoRepo.Close()

	// UseCase 생성
	geoUseCase := usecase.NewGeoUseCaseWithGeoLite2(geoRepo)
	defer geoUseCase.Close()

	// 종합적인 지리 정보 조회
	geoData, err := geoUseCase.GetGeoData(testIP)
	if err != nil {
		log.Fatalf("지리 정보 조회 실패: %v", err)
	}

	// 결과 출력 - 모든 정보 표시
	printGeoDataDetails(geoData)

	// 추가로 개별 정보에 대한 상세 출력
	city, err := geoUseCase.GetCityInfo(testIP)
	if err == nil {
		fmt.Println("\n[도시 상세 정보]")
		printCityDetails(city, testIP)
	}

	country, err := geoUseCase.GetCountryInfo(testIP)
	if err == nil {
		fmt.Println("\n[국가 상세 정보]")
		printCountryDetails(country, testIP)
	}

	asn, err := geoUseCase.GetASNInfo(testIP)
	if err == nil {
		fmt.Println("\n[ASN 상세 정보]")
		printASNDetails(asn, testIP)
	}

	// Output:
	// (실제 출력은 사용하는 데이터베이스에 따라 다를 수 있습니다)
}

func Example_anonymousCheck() {
	// Anonymous IP DB 리포지토리 생성
	anonymousRepo, err := repository.NewGeoIP2AnonymousIPRepository("/path/to/GeoIP2-Anonymous-IP.mmdb")
	if err != nil {
		log.Fatalf("Anonymous IP DB 열기 실패: %v", err)
	}
	defer anonymousRepo.Close()

	// UseCase 생성
	geoUseCase := usecase.NewGeoUseCase(nil, nil, nil, anonymousRepo)
	defer geoUseCase.Close()

	// IP 주소가 익명 프록시인지 확인
	isAnonymous, err := geoUseCase.IsAnonymousIP(testIP)
	if err != nil {
		log.Fatalf("익명 IP 확인 실패: %v", err)
	}

	// IP 주소가 Tor 출구 노드인지 확인
	isTorExit, err := geoUseCase.IsTorExitNode(testIP)
	if err != nil {
		log.Fatalf("Tor 출구 노드 확인 실패: %v", err)
	}

	// 결과 출력 - 모든 정보 표시
	fmt.Printf("IP 주소: %s\n", testIP)
	fmt.Printf("익명 프록시 여부: %v\n", isAnonymous)
	fmt.Printf("Tor 출구 노드 여부: %v\n", isTorExit)

	// Anonymous IP DB에서 직접 정보 조회
	ip := net.ParseIP(testIP)
	if ip != nil {
		anonIP, err := anonymousRepo.GetAnonymousIP(ip)
		if err == nil {
			fmt.Printf("\n[익명 IP 상세 정보]\n")
			fmt.Printf("익명 IP 여부: %v\n", anonIP.IsAnonymous)
			fmt.Printf("익명 VPN 여부: %v\n", anonIP.IsAnonymousVPN)
			fmt.Printf("호스팅 제공자 여부: %v\n", anonIP.IsHostingProvider)
			fmt.Printf("공개 프록시 여부: %v\n", anonIP.IsPublicProxy)
			fmt.Printf("주거용 프록시 여부: %v\n", anonIP.IsResidentialProxy)
			fmt.Printf("Tor 출구 노드 여부: %v\n", anonIP.IsTorExitNode)
		}
	}

	// Output:
	// (실제 출력은 사용하는 데이터베이스에 따라 다를 수 있습니다)
}

// 도시 정보 상세 출력 헬퍼 함수
func printCityDetails(city entity.City, ipStr string) {
	fmt.Printf("IP 주소: %s\n", ipStr)
	fmt.Printf("\n[도시 정보]\n")
	fmt.Printf("도시 GeoName ID: %d\n", city.City.GeoNameID)

	// 도시 이름 출력 (여러 언어)
	fmt.Println("도시 이름:")
	for lang, name := range city.City.Names {
		fmt.Printf("  - %s: %s\n", lang, name)
	}

	fmt.Printf("\n[국가 정보]\n")
	fmt.Printf("국가 GeoName ID: %d\n", city.Country.GeoNameID)
	fmt.Printf("국가 ISO 코드: %s\n", city.Country.IsoCode)
	fmt.Printf("EU 소속 여부: %v\n", city.Country.IsInEuropeanUnion)

	// 국가 이름 출력 (여러 언어)
	fmt.Println("국가 이름:")
	for lang, name := range city.Country.Names {
		fmt.Printf("  - %s: %s\n", lang, name)
	}

	fmt.Printf("\n[대륙 정보]\n")
	fmt.Printf("대륙 GeoName ID: %d\n", city.Continent.GeoNameID)
	fmt.Printf("대륙 코드: %s\n", city.Continent.Code)

	// 대륙 이름 출력 (여러 언어)
	fmt.Println("대륙 이름:")
	for lang, name := range city.Continent.Names {
		fmt.Printf("  - %s: %s\n", lang, name)
	}

	fmt.Printf("\n[위치 정보]\n")
	fmt.Printf("위도: %.6f\n", city.Location.Latitude)
	fmt.Printf("경도: %.6f\n", city.Location.Longitude)
	fmt.Printf("시간대: %s\n", city.Location.TimeZone)
}

// 국가 정보 상세 출력 헬퍼 함수
func printCountryDetails(country entity.Country, ipStr string) {
	fmt.Printf("IP 주소: %s\n", ipStr)

	fmt.Printf("\n[국가 정보]\n")
	fmt.Printf("국가 GeoName ID: %d\n", country.Country.GeoNameID)
	fmt.Printf("국가 ISO 코드: %s\n", country.Country.IsoCode)
	fmt.Printf("EU 소속 여부: %v\n", country.Country.IsInEuropeanUnion)

	// 국가 이름 출력 (여러 언어)
	fmt.Println("국가 이름:")
	for lang, name := range country.Country.Names {
		fmt.Printf("  - %s: %s\n", lang, name)
	}

	fmt.Printf("\n[대륙 정보]\n")
	fmt.Printf("대륙 GeoName ID: %d\n", country.Continent.GeoNameID)
	fmt.Printf("대륙 코드: %s\n", country.Continent.Code)

	// 대륙 이름 출력 (여러 언어)
	fmt.Println("대륙 이름:")
	for lang, name := range country.Continent.Names {
		fmt.Printf("  - %s: %s\n", lang, name)
	}
}

// ASN 정보 상세 출력 헬퍼 함수
func printASNDetails(asn entity.ASN, ipStr string) {
	fmt.Printf("IP 주소: %s\n", ipStr)
	fmt.Printf("\n[ASN 정보]\n")
	fmt.Printf("자율 시스템 번호(ASN): %d\n", asn.AutonomousSystemNumber)
	fmt.Printf("자율 시스템 조직: %s\n", asn.AutonomousSystemOrganization)
}

// GeoData 종합 정보 상세 출력 헬퍼 함수
func printGeoDataDetails(geoData *usecase.GeoData) {
	fmt.Printf("IP 주소: %s\n", geoData.IPAddress)
	fmt.Printf("유효한 IP 주소 여부: %v\n", geoData.IsValid)

	if geoData.City != "" {
		fmt.Printf("도시: %s\n", geoData.City)
	}

	if geoData.CountryCode != "" {
		fmt.Printf("국가 코드: %s\n", geoData.CountryCode)
	}

	if geoData.CountryName != "" {
		fmt.Printf("국가명: %s\n", geoData.CountryName)
	}

	if geoData.ContinentCode != "" {
		fmt.Printf("대륙 코드: %s\n", geoData.ContinentCode)
	}

	if geoData.Latitude != 0 || geoData.Longitude != 0 {
		fmt.Printf("위치: 위도 %.6f, 경도 %.6f\n", geoData.Latitude, geoData.Longitude)
	}

	if geoData.TimeZone != "" {
		fmt.Printf("시간대: %s\n", geoData.TimeZone)
	}

	if geoData.ASN != 0 {
		fmt.Printf("ASN 번호: %d\n", geoData.ASN)
	}

	if geoData.ISP != "" {
		fmt.Printf("ISP: %s\n", geoData.ISP)
	}

	// 익명 IP 관련 정보
	fmt.Printf("익명 프록시 여부: %v\n", geoData.IsAnonymous)
	fmt.Printf("익명 VPN 여부: %v\n", geoData.IsAnonymousVPN)
	fmt.Printf("Tor 출구 노드 여부: %v\n", geoData.IsTorExitNode)
}
