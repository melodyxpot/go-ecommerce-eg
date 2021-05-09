package customer

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/postgresql"
	"time"
)

// Repository db repository
type Repository interface {
	// get customer
	Get(c echo.Context, id int64) (cus customer, err error)

	// create new customer
	Create(c echo.Context, email, firstName, lastName string) (id int64, err error)

	// update customer
	Update(c echo.Context, id int64, email, firstName, lastName string) error

	// delete customer
	Delete(c echo.Context, id int64) error
}

type repository struct {
	pdb postgresql.DB
}

// NewRepository returns a new repostory
func NewRepository(pdb postgresql.DB) Repository {
	return &repository{pdb}
}

func (r repository) Get(c echo.Context, id int64) (cus customer, err error) {
	sql := `select email, first_name, last_name from customers where id = $1`
	row := r.pdb.QueryRow(c, sql, id)
	err = row.Scan(&cus.Email, &cus.FirstName, &cus.LastName)
	if err != nil {
		return cus, errors.Wrapf(errorsd.New("error getting customer from db"), "error getting customer from db: %s", err.Error())
	}
	return cus, err
}

func (r repository) Create(c echo.Context, email, firstName, lastName string) (id int64, err error) {
	now := time.Now().UTC()
	sql := `insert into customers (email, first_name, last_name, created_at, updated_at) values
		($1, $2, $3, $4, $5) returning (id)`
	err = r.pdb.QueryRow(c, sql, email, firstName, lastName, now, now).Scan(&id)
	return id, nil
}

func (r repository) Update(c echo.Context, id int64, email, firstName, lastName string) error {
	sql := `update customers set email = $1, first_name = $2, last_name = $3 where id = $4`
	tag, err := r.pdb.Exec(c, sql, email, firstName, lastName, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("no updated record changed")
	}
	return nil
}

func (r repository) Delete(c echo.Context, id int64) error {
	sql := `delete from customers where id = $1`
	tag, err := r.pdb.Exec(c, sql, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("no delete record changed")
	}
	return nil
}
