package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	elastic "gopkg.in/olivere/elastic.v5"
)


var (
	ErrNotFound = errors.New("Entity not found")
)

type Respository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product , error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)

}

type elasticRespository struct {
	client *elastic.Client
}

type productDocument struct {
	Name 		string		`json:"name"`
	Description  string 	 `json:"description"`
	Price 		 float64      `json:"price"`
}

func NewElasticRepository(url string) (Respository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)

	if err != nil {
		return nil, err
	}

	return &elasticRespository{client},nil
}

func (r *elasticRespository) Close() {
}

func (r *elasticRespository) PutProduct(ctx context.Context, p Product) error {
	_, err := r.client.Index().
	          Index("catalog").
			  Type("product").
			  Id(p.ID).
			  BodyJson(productDocument{
				Name: p.Name,
				Description: p.Description,
				Price:  p.Price,
			  }).
			  Do(ctx)

			  return err
}


func (r *elasticRespository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	  a, err := r.client.Index().
	            Index("catalog").
				Type("product").
				Id(id).
				Do(ctx)

	  if err != nil {
		return nil, err
	  }

      if !a.Found {
		return nil, ErrNotFound
	  } 

      p := productDocument{}

	  if err = json.Unmarshal(*a.Source, &p); err != nil {
		return nil, err
	  }

	  return &Product {
		ID:   id,
        Name: p.Name,
		Description: p.Description,
		Price: p.Price,
	  }, err
}


func (r *elasticRespository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	a, err := r.client.Search().
	          Index("catalog").
			  Type("product").
			  Query(elastic.NewMatchAllQuery()).
			  From(int(skip)).Size(int(take)).
			  Do(ctx)
    if err != nil {
		log.Println(err)
		return nil, err
	}
	
	prod := []Product{}

	for _, hit := range a.Hits.Hits {
       p := productDocument{}
	   if err = json.Unmarshal(*hit.Source, &p); err == nil {
		  prod = append(prod, Product{
			    ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
		  })
	   }
	}

	return prod, err
}



func (r *elasticRespository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {

	items := []*elastic.MultiGetItem{}

	for _, id := range ids {  // creating list of obj eg {"_index": "catalog","_type": "product","_id": "101"} This objs will be further used to fetch the product from the 
	// elastic search its not doing db.query

		items = append(items, elastic.NewMultiGetItem().
	             Index("Catalog").
				 Type("product"). 
				 Id(id),
	           )

	}
   
	res, err := r.client.MultiGet().  // instead of performing GET ops for each id fetch Add() Addes all the item that to be fetched in a single go
	            Add(items...).
				Do(ctx)

   if err != nil {
	log.Fatal(err)
	return nil, err
   }
   
   products := []Product{}

   for _, doc := range res.Docs {
	p := productDocument{}
	if err = json.Unmarshal(*doc.Source, p); err == nil {
	products = append(products, Product{
        ID:  	doc.Id,
		Name:   p.Name,
		Description:   p.Description,
		Price:      p.Price,
	  })
   }
 }

 return products, nil

}


func (r *elasticRespository) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err = json.Unmarshal(*hit.Source, &p); err == nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, err
}