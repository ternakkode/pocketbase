package forms

import (
	"log"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/forms/validators"
	"github.com/pocketbase/pocketbase/models"
)

// AddressUpsert is a [models.Address] upsert (create/update) form.
type AddressUpsert struct {
	app     core.App
	dao     *daos.Dao
	address *models.Address

	Id      string `form:"id" json:"id"`
	Street  string `form:"street" json:"street"`
	City    string `form:"city" json:"city"`
	State   string `form:"state" json:"state"`
	ZipCode string `form:"zipCode" json:"zipCode"`
	Country string `form:"country" json:"country"`
}

func NewAddressUpsert(app core.App, address *models.Address) *AddressUpsert {
	form := &AddressUpsert{
		app:     app,
		dao:     app.Dao(),
		address: address,
	}

	// load defaults
	form.Id = address.Id
	form.Street = address.Street
	form.City = address.City
	form.State = address.State
	form.ZipCode = address.ZipCode
	form.Country = address.Country

	return form
}

func (form *AddressUpsert) SetDao(dao *daos.Dao) {
	form.dao = dao
}

func (form *AddressUpsert) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Id,
			validation.When(
				form.address.IsNew(),
				validation.Length(models.DefaultIdLength, models.DefaultIdLength),
				validation.Match(idRegex),
				validation.By(validators.UniqueId(form.dao, form.address.TableName())),
			).Else(validation.In(form.address.Id)),
		),
		validation.Field(
			&form.Street,
			validation.Required,
			validation.Length(1, 255),
		),
		validation.Field(
			&form.City,
			validation.Required,
			validation.Length(1, 255),
		),
		validation.Field(
			&form.State,
			validation.Required,
			validation.Length(1, 255),
		),
		validation.Field(
			&form.ZipCode,
			validation.Required,
			validation.Length(1, 6),
		),
		validation.Field(
			&form.Country,
			validation.Required,
			validation.Length(1, 6),
		),
	)
}

func (form *AddressUpsert) Submit(interceptors ...InterceptorFunc[*models.Address]) error {
	if err := form.Validate(); err != nil {
		return err
	}

	// custom insertion id can be set only on create
	if form.address.IsNew() && form.Id != "" {
		form.address.MarkAsNew()
		form.address.SetId(form.Id)
	}

	form.address.Street = form.Street
	form.address.City = form.City
	form.address.State = form.State
	form.address.ZipCode = form.ZipCode
	form.address.Country = form.Country

	return runInterceptors(form.address, func(address *models.Address) error {
		log.Println("Saving address")
		return form.dao.Save(address)
	}, interceptors...)
}
