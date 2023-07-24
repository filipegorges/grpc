package service

import (
	"context"
	"io"
	"log"

	"github.com/filipegorges/fullcycleGRPC/internal/database"
	"github.com/filipegorges/fullcycleGRPC/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(db database.Category) *CategoryService {
	return &CategoryService{CategoryDB: db}
}

func (s *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := s.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		log.Printf("failed to create category %v", err)
		return nil, err
	}

	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (s *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {
	categories, err := s.CategoryDB.FindAll()
	if err != nil {
		log.Printf("failed to list categories %v", err)
		return nil, err
	}

	var pbCategories []*pb.Category
	for _, category := range categories {
		pbCategories = append(pbCategories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return &pb.CategoryList{
		Categories: pbCategories,
	}, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
	category, err := s.CategoryDB.Find(in.Id)
	if err != nil {
		log.Printf("failed to get category %v", err)
		return nil, err
	}

	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (s *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			log.Printf("failed to receive category %v", err)
			return err
		}
		categoryResult, err := s.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			log.Printf("failed to create category %v", err)
			return err
		}
		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
	}
}

func (s *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("failed to receive category %v", err)
			return err
		}
		categoryResult, err := s.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			log.Printf("failed to create category %v", err)
			return err
		}
		err = stream.Send(&pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
		if err != nil {
			log.Printf("failed to send category %v", err)
			return err
		}
	}
}
