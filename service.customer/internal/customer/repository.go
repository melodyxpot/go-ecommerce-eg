package customer

import (
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
	"github.com/labstack/echo"
)

// Repository db repository
type Repository interface {
	// get customer
	Get(c echo.Context, id int64) (cus customer, err error)

	// create new customer
	Create(c echo.Context, email string) (id int64, err error)

	// update customer
	Update(c echo.Context, id int64, email string) error

	// delete customer
	Delete(c echo.Context, id int64) error
}

type repository struct {
	db dbcontext.DB
}

// NewRepository returns a new repostory
func NewRepository(db dbcontext.DB) Repository {
	return &repository{db}
}

func (r repository) Get(c echo.Context, id int64) (cus customer, err error) {
	return cus, err
}

func (r repository) Create(c echo.Context, email string) (id int64, err error) {
	return id, nil
}

func (r repository) Update(c echo.Context, id int64, email string) error {
	return nil
}

func (r repository) Delete(c echo.Context, id int64) error {
	return nil
}
