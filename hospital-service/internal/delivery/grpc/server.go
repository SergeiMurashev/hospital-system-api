package grpc

import (
	"context"
	"github.com/sergeimurashev/hospital-system-api/proto"

	"github.com/sergeimurashev/hospital-system-api/hospital-service/internal/domain"
	"github.com/sergeimurashev/hospital-system-api/hospital-service/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	proto.UnimplementedHospitalServiceServer
	hospitalService service.HospitalService
}

func NewServer(hospitalService service.HospitalService) *Server {
	return &Server{
		hospitalService: hospitalService,
	}
}

func (s *Server) CreateHospital(ctx context.Context, req *proto.CreateHospitalRequest) (*proto.Hospital, error) {
	hospital, err := s.hospitalService.Create(ctx, domain.CreateHospitalRequest{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Rooms:   req.Rooms,
	})
	if err != nil {
		return nil, err
	}

	return convertHospitalToProto(hospital), nil
}

func (s *Server) GetHospital(ctx context.Context, req *proto.GetHospitalRequest) (*proto.Hospital, error) {
	hospital, err := s.hospitalService.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return convertHospitalToProto(hospital), nil
}

func (s *Server) UpdateHospital(ctx context.Context, req *proto.UpdateHospitalRequest) (*proto.Hospital, error) {
	hospital, err := s.hospitalService.Update(ctx, domain.UpdateHospitalRequest{
		ID:      req.Id,
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Rooms:   req.Rooms,
	})
	if err != nil {
		return nil, err
	}

	return convertHospitalToProto(hospital), nil
}

func (s *Server) DeleteHospital(ctx context.Context, req *proto.DeleteHospitalRequest) (*proto.DeleteHospitalResponse, error) {
	err := s.hospitalService.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteHospitalResponse{
		Success: true,
	}, nil
}

func (s *Server) ListHospitals(ctx context.Context, req *proto.ListHospitalsRequest) (*proto.ListHospitalsResponse, error) {
	hospitals, total, err := s.hospitalService.List(ctx, int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, err
	}

	protoHospitals := make([]*proto.Hospital, len(hospitals))
	for i, hospital := range hospitals {
		protoHospitals[i] = convertHospitalToProto(hospital)
	}

	return &proto.ListHospitalsResponse{
		Hospitals: protoHospitals,
		Total:     int32(total),
	}, nil
}

func (s *Server) GetRooms(ctx context.Context, req *proto.GetRoomsRequest) (*proto.GetRoomsResponse, error) {
	rooms, err := s.hospitalService.GetRooms(ctx, req.HospitalId)
	if err != nil {
		return nil, err
	}

	protoRooms := make([]*proto.Room, len(rooms))
	for i, room := range rooms {
		protoRooms[i] = convertRoomToProto(room)
	}

	return &proto.GetRoomsResponse{
		Rooms: protoRooms,
	}, nil
}

func convertHospitalToProto(hospital *domain.Hospital) *proto.Hospital {
	if hospital == nil {
		return nil
	}

	protoRooms := make([]*proto.Room, len(hospital.Rooms))
	for i, room := range hospital.Rooms {
		protoRooms[i] = convertRoomToProto(room)
	}

	return &proto.Hospital{
		Id:        hospital.ID,
		Name:      hospital.Name,
		Address:   hospital.Address,
		Phone:     hospital.Phone,
		CreatedAt: timestamppb.New(hospital.CreatedAt),
		UpdatedAt: timestamppb.New(hospital.UpdatedAt),
		Rooms:     protoRooms,
	}
}

func convertRoomToProto(room *domain.Room) *proto.Room {
	if room == nil {
		return nil
	}

	return &proto.Room{
		Id:         room.ID,
		Name:       room.Name,
		HospitalId: room.HospitalID,
		CreatedAt:  timestamppb.New(room.CreatedAt),
		UpdatedAt:  timestamppb.New(room.UpdatedAt),
	}
}
