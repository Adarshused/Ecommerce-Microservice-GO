package main



import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/Adarshused/Ecommerce-Microservice-GO/account"
	"github.com/Adarshused/Ecommerce-Microservice-GO/catalog"
	"github.com/Adarshused/Ecommerce-Microservice-GO/order"
)

type Server struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient *order.Client
} 



func NewGraphQLServer(accountURL, catalogURL, orderURL string) (*Server, error) {

	// connect to account service
	accountClient, err := account.NewClient(accountURL)

	if err != nil {
		return nil, err;
	}

	// connect to catalog servie
	catalogClient, err := catalog.NewClient(catalogURL)

	if err != nil {
		accountClient.Close()
		return nil, err;
	}

	// connect to order service
	orderClient, err := order.NewClient(orderURL)

	if err != nil {
		accountClient.Close()
        catalogClient.Close()
		return nil, err
	}

	return &Server{
		accountClient,
		catalogClient,
		orderClient,
	}, nil

}


// func (s *Server) Mutation() MutationResolver {

//        return mutationResolver {
// 		server: s,
// 	   }
// }

func (s *Server) Query() queryResolver {

	return  queryResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}



