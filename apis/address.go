package apis

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

func bindAddressApi(app core.App, rg *echo.Group) {
	api := addressAPI{app: app}

	subGroup := rg.Group("/address", ActivityLogger(app))
	subGroup.POST("", api.create, RequireAdminAuthOnlyIfAny(app))
}

type addressAPI struct {
	app core.App
}

func (api *addressAPI) create(c echo.Context) error {
	address := &models.Address{}

	form := forms.NewAddressUpsert(api.app, address)

	// load request
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("Failed to load the submitted data due to invalid formatting.", err)
	}

	submitErr := form.Submit()
	if submitErr != nil {
		return NewBadRequestError("Failed to submit the form due to invalid data.", submitErr)
	}

	c.JSON(200, address)
	return nil
}
