package plugins

import (
	"context"
	"fmt"

	"saastack/core"
	pb "saastack/interfaces/bookstore/proto"
)

// SimpleBookstorePlugin implements the BookstoreServiceServer interface
type SimpleBookstorePlugin struct {
	pb.UnimplementedBookstoreServiceServer
}

func init() {
	plugin := NewSimpleBookstorePlugin()
	core.GlobalRegistry.RegisterPlugin("bookstore", "simple", plugin)
}

// NewSimpleBookstorePlugin creates a new instance of SimpleBookstorePlugin
func NewSimpleBookstorePlugin() *SimpleBookstorePlugin {
	return &SimpleBookstorePlugin{}
}

func (s *SimpleBookstorePlugin) AddBook(ctx context.Context, req *pb.AddBookRequest) (*pb.GenericResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	fmt.Printf("SimpleBookstorePlugin AddBook: ISBN=%s, Title=%s, Author=%s, Year=%s\n",
		req.Isbn, req.Title, req.Author, req.YearPublished)
	return &pb.GenericResponse{
		Result: fmt.Sprintf("Book added: %s by %s (ISBN: %s)", req.Title, req.Author, req.Isbn),
	}, nil
}

func (s *SimpleBookstorePlugin) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GenericResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	fmt.Printf("SimpleBookstorePlugin GetBook: ISBN=%s\n", req.Isbn)
	return &pb.GenericResponse{
		Result: fmt.Sprintf("Book retrieved: ISBN %s", req.Isbn),
	}, nil
}

func (s *SimpleBookstorePlugin) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	fmt.Println("SimpleBookstorePlugin ListBooks: Listing all books")
	return &pb.ListBooksResponse{
		Books: []string{
			"Book 1: Go Programming (ISBN: 123)",
			"Book 2: Microservices Architecture (ISBN: 456)",
			"Book 3: Database Design (ISBN: 789)",
		},
	}, nil
}

func (s *SimpleBookstorePlugin) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.GenericResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	fmt.Printf("SimpleBookstorePlugin DeleteBook: ISBN=%s\n", req.Isbn)
	return &pb.GenericResponse{
		Result: fmt.Sprintf("Book deleted: ISBN %s", req.Isbn),
	}, nil
}
