package catalog

import (
	"context"
	"log"

	"github.com/segmentio/ksuid"
)


type Service interface {
	PostProduct (ctx context.Context, name, description string, price float64) (*Product, error)
	GetProduct (ctx context.Context, id string) (*Product, error)
	GetProducts (ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductsByIDs (ctx context.Context, ids []string) ([]Product, error)
	SearchProducts (ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
    
}

type Product struct {
   
	ID  string 	`json:"id"`
	Name string  `json:"name"`
	Description  string   `json:"description"`
	Price    float64 	`json:"price"`
}

type catalogService struct {
     repository	Respository
}

func NewService(r Respository) Service {
	return &catalogService{r}
}

func (s *catalogService) PostProduct (ctx context.Context, name , description string , price float64) (*Product, error) {

     p := &Product{
          ID: ksuid.New().String(),
		  Name: name,
		  Description:  description,
		  Price:   price,
	    }

	 if err := s.repository.PutProduct(ctx, *p); err != nil {
		log.Fatal("Error while putting the product")
        return nil, err
	 }

	 return p, nil

}

func (s *catalogService) GetProduct (ctx context.Context, id string) (*Product, error) {

	  a, err := s.repository.GetProductByID(ctx, id)

	  if err != nil {
		log.Fatal("Error while fetching single product")
	  }

	  return a, nil;
}

func (s *catalogService) GetProducts (ctx context.Context, skip uint64, take uint64) ([]Product, error) {
     
	 a,err := s.repository.ListProducts(ctx, skip, take);
     
	 if err != nil {
		return nil, err
	 }

	 return a, nil
}

func (s *catalogService) GetProductsByIDs (ctx context.Context, ids []string) ([]Product, error) {
     
	a, err := s.repository.ListProductsWithIDs(ctx, ids);

	if err != nil {
		return nil, err
	}

	return a, nil

}


func (s *catalogService) SearchProducts (ctx context.Context, query string , skip uint64, take uint64) ([]Product, error) {

    if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	
	return s.repository.SearchProducts(ctx, query, skip, take)
}

