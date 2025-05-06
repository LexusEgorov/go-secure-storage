package grpcserv

import (
	"context"
	"data/internal/models"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/LexusEgorov/go-secure-storage-protos/gen/golang/datapb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type DataProvider interface {
	AddCard(card models.Card, userId int) (string, error)
	AddText(text string, userId int) (string, error)
	AddBinary(binary []byte, userId int) (string, error)
	AddPassword(password models.Password, userId int) (string, error)

	GetCard(filename string) (*models.Card, error)
	GetText(filename string) (string, error)
	GetBinary(filename string) ([]byte, error)
	GetPassword(filename string) (*models.Password, error)
}

type Server struct {
	datapb.UnimplementedDataServer
	l    *logrus.Logger
	s    *grpc.Server
	data DataProvider
}

func (s Server) Add(ctx context.Context, req *datapb.AddRequest) (*datapb.AddResponse, error) {
	dataItem := req.GetData()
	category := req.GetCategory()
	userId := -1

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values, found := md["authorization"]; found {
			var err error
			userId, err = strconv.Atoi(values[0])

			if err != nil {
				return nil, err
			}
		}
	}

	if userId == -1 {
		return nil, errors.New("unauthorized")
	}

	var filename string
	var err error

	switch category {
	case datapb.Category_BINARY:
		filename, err = s.data.AddBinary(dataItem.GetBinary().Binary, userId)
	case datapb.Category_PASSWORD:
		password := models.Password{
			Login:    dataItem.GetPassword().Login,
			Password: dataItem.GetPassword().Password,
		}

		filename, err = s.data.AddPassword(password, userId)
	case datapb.Category_TEXT:
		filename, err = s.data.AddText(dataItem.GetText().Text, userId)
	case datapb.Category_CARD:
		card := models.Card{
			Cvv:    dataItem.GetCard().Cvv,
			Exp:    dataItem.GetCard().Exp,
			Holder: dataItem.GetCard().Holder,
			Number: dataItem.GetCard().Number,
		}

		filename, err = s.data.AddCard(card, userId)
	}

	if err != nil {
		return &datapb.AddResponse{
			Ok: false,
			Response: &datapb.AddResponse_Bad{
				Bad: &datapb.BadResponse{
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &datapb.AddResponse{
		Ok: true,
		Response: &datapb.AddResponse_Success{
			Success: &datapb.SuccessAddResponse{
				Filename: filename,
				Category: req.Category,
			},
		},
	}, nil
}

func (s Server) Get(ctx context.Context, req *datapb.GetRequest) (*datapb.GetResponse, error) {
	category := req.GetCategory()
	filename := req.GetFilename()
	var item datapb.DataItem

	switch category {
	case datapb.Category_BINARY:
		card, err := s.data.GetCard(filename)

		if err != nil {
			return nil, err
		}

		item.Data = &datapb.DataItem_Card{
			Card: &datapb.Card{
				Cvv:    card.Cvv,
				Exp:    card.Exp,
				Holder: card.Holder,
				Number: card.Number,
			},
		}
	case datapb.Category_TEXT:
		text, err := s.data.GetText(filename)

		if err != nil {
			return nil, err
		}

		item.Data = &datapb.DataItem_Text{
			Text: &datapb.Text{
				Text: text,
			},
		}
	case datapb.Category_CARD:
		card, err := s.data.GetCard(filename)

		if err != nil {
			return nil, err
		}

		item.Data = &datapb.DataItem_Card{
			Card: &datapb.Card{
				Cvv:    card.Cvv,
				Exp:    card.Exp,
				Holder: card.Holder,
				Number: card.Number,
			},
		}
	case datapb.Category_PASSWORD:
		password, err := s.data.GetPassword(filename)

		if err != nil {
			return nil, err
		}

		item.Data = &datapb.DataItem_Password{
			Password: &datapb.Password{
				Login:    password.Login,
				Password: password.Password,
			},
		}
	}

	return &datapb.GetResponse{
		Ok: true,
		Response: &datapb.GetResponse_Success{
			Success: &datapb.SuccessGetResponse{
				Data: &item,
			},
		},
	}, nil
}

func NewServer(l *logrus.Logger, data DataProvider) *Server {
	grpcServer := grpc.NewServer()

	server := Server{
		l:    l,
		s:    grpcServer,
		data: data,
	}

	datapb.RegisterDataServer(grpcServer, server)

	return &server
}

func (s Server) RunServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		s.l.Panic(err)
		return err
	}

	s.l.Info("server is running on ", port, " port")

	if err := s.s.Serve(lis); err != nil {
		s.l.Panic(err)
		return err
	}

	return nil
}
