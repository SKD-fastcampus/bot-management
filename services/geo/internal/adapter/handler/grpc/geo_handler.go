package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "github.com/wekeepgrowing/semo-backend-monorepo/proto/geo/v1"
	"github.com/wekeepgrowing/semo-backend-monorepo/services/geo/internal/usecase"
)

// GeoHandler는 지오로케이션 관련 gRPC 핸들러입니다
type GeoHandler struct {
	proto.UnimplementedGeoServiceServer
	geoUseCase *usecase.GeoUseCase
}

// NewGeoHandler는 새로운 GeoHandler 인스턴스를 생성합니다
func NewGeoHandler(geoUseCase *usecase.GeoUseCase) *GeoHandler {
	return &GeoHandler{
		geoUseCase: geoUseCase,
	}
}

// GetGeoData는 IP 주소에 대한 종합적인 지리 정보를 반환합니다
func (h *GeoHandler) GetGeoData(ctx context.Context, req *proto.IpRequest) (*proto.GeoDataResponse, error) {
	if req.Ip == "" {
		return nil, status.Error(codes.InvalidArgument, "IP 주소가 필요합니다")
	}

	geoData, err := h.geoUseCase.GetGeoData(req.Ip)
	if err != nil {
		if err == usecase.ErrInvalidIPAddress {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &proto.GeoDataResponse{
		IpAddress:      geoData.IPAddress,
		City:           geoData.City,
		CountryCode:    geoData.CountryCode,
		CountryName:    geoData.CountryName,
		ContinentCode:  geoData.ContinentCode,
		Latitude:       geoData.Latitude,
		Longitude:      geoData.Longitude,
		TimeZone:       geoData.TimeZone,
		Asn:            uint32(geoData.ASN),
		Isp:            geoData.ISP,
		IsValid:        geoData.IsValid,
		IsAnonymous:    geoData.IsAnonymous,
		IsAnonymousVpn: geoData.IsAnonymousVPN,
		IsTorExitNode:  geoData.IsTorExitNode,
	}

	return response, nil
}

// GetCityInfo는 IP 주소에 대한 도시 정보를 반환합니다
func (h *GeoHandler) GetCityInfo(ctx context.Context, req *proto.IpRequest) (*proto.CityResponse, error) {
	if req.Ip == "" {
		return nil, status.Error(codes.InvalidArgument, "IP 주소가 필요합니다")
	}

	city, err := h.geoUseCase.GetCityInfo(req.Ip)
	if err != nil {
		if err == usecase.ErrInvalidIPAddress {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &proto.CityResponse{
		City: &proto.CityInfo{
			GeonameId: uint32(city.City.GeoNameID),
			Names:     city.City.Names,
		},
		Country: &proto.CountryInfo{
			GeonameId:         uint32(city.Country.GeoNameID),
			IsInEuropeanUnion: city.Country.IsInEuropeanUnion,
			IsoCode:           city.Country.IsoCode,
			Names:             city.Country.Names,
		},
		Continent: &proto.ContinentInfo{
			Code:      city.Continent.Code,
			GeonameId: uint32(city.Continent.GeoNameID),
			Names:     city.Continent.Names,
		},
		Location: &proto.LocationInfo{
			Latitude:  city.Location.Latitude,
			Longitude: city.Location.Longitude,
			TimeZone:  city.Location.TimeZone,
		},
	}

	return response, nil
}

// GetCountryInfo는 IP 주소에 대한 국가 정보를 반환합니다
func (h *GeoHandler) GetCountryInfo(ctx context.Context, req *proto.IpRequest) (*proto.CountryResponse, error) {
	if req.Ip == "" {
		return nil, status.Error(codes.InvalidArgument, "IP 주소가 필요합니다")
	}

	country, err := h.geoUseCase.GetCountryInfo(req.Ip)
	if err != nil {
		if err == usecase.ErrInvalidIPAddress {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &proto.CountryResponse{
		Country: &proto.CountryInfo{
			GeonameId:         uint32(country.Country.GeoNameID),
			IsInEuropeanUnion: country.Country.IsInEuropeanUnion,
			IsoCode:           country.Country.IsoCode,
			Names:             country.Country.Names,
		},
		Continent: &proto.ContinentInfo{
			Code:      country.Continent.Code,
			GeonameId: uint32(country.Continent.GeoNameID),
			Names:     country.Continent.Names,
		},
	}

	return response, nil
}

// GetASNInfo는 IP 주소에 대한 ASN 정보를 반환합니다
func (h *GeoHandler) GetASNInfo(ctx context.Context, req *proto.IpRequest) (*proto.ASNResponse, error) {
	if req.Ip == "" {
		return nil, status.Error(codes.InvalidArgument, "IP 주소가 필요합니다")
	}

	asn, err := h.geoUseCase.GetASNInfo(req.Ip)
	if err != nil {
		if err == usecase.ErrInvalidIPAddress {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &proto.ASNResponse{
		AutonomousSystemNumber:       uint32(asn.AutonomousSystemNumber),
		AutonomousSystemOrganization: asn.AutonomousSystemOrganization,
	}

	return response, nil
}

// CheckAnonymousIP는 IP 주소가 익명 프록시인지 확인합니다
func (h *GeoHandler) CheckAnonymousIP(ctx context.Context, req *proto.IpRequest) (*proto.AnonymousResponse, error) {
	if req.Ip == "" {
		return nil, status.Error(codes.InvalidArgument, "IP 주소가 필요합니다")
	}

	isAnonymous, err := h.geoUseCase.IsAnonymousIP(req.Ip)
	if err != nil {
		if err == usecase.ErrInvalidIPAddress {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else if err == usecase.ErrFeatureNotSupported {
			return &proto.AnonymousResponse{
				IsAnonymous:    false,
				IsTorExitNode:  false,
				FeatureSupport: false,
			}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	isTorExit, _ := h.geoUseCase.IsTorExitNode(req.Ip)

	response := &proto.AnonymousResponse{
		IsAnonymous:    isAnonymous,
		IsTorExitNode:  isTorExit,
		FeatureSupport: true,
	}

	return response, nil
}
