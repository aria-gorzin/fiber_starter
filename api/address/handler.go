package address

import (
	"net/http"
	"strconv"

	db "github.com/aria/app/db/sqlc"
	"github.com/aria/app/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type AddressRouter struct {
	store    db.Store
	validate *validator.Validate
}

func NewRouter(store db.Store, v *validator.Validate) *AddressRouter {
	return &AddressRouter{store: store, validate: v}
}

/* -------- register -------- */

func (r *AddressRouter) Register(app *fiber.App) {
	ad := app.Group("/addresses")

	ad.Post("", r.createAddress)
	ad.Get("/:id", r.getAddressByID)
	ad.Get("", r.listAddresses) // ?client_id=...
	ad.Put("/:id", r.updateAddressByID)
	ad.Delete("/:id", r.deleteAddress)
}

/* -------- CRUD -------- */

// createAddress creates a new address.
// @Summary Create a new address
// @Description Create a new address
// @Tags addresses
// @Accept json
// @Produce json
// @Param request body createAddressRequest true "Address Info"
// @Success 200 {object} addressResponse
// @Failure 400 {object} util.CustomError
// @Failure 422 {object} util.ValidateErrors
// @Router /addresses [post]
// @Security BearerAuth
func (r *AddressRouter) createAddress(c *fiber.Ctx) error {
	var req createAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return util.ErrInvalidBody(err)
	}
	if err := r.validate.Struct(&req); err != nil {
		return util.NewValidationErrors(err.(validator.ValidationErrors))
	}

	addr, err := r.store.CreateAddress(c.Context(), db.CreateAddressParams{
		ClientID: req.ClientID,
		Title:    req.Title,
		City:     req.City,
		Street:   pgtype.Text{String: util.Deref(req.Street), Valid: req.Street != nil},
		Phone:    pgtype.Text{String: util.Deref(req.Phone), Valid: req.Phone != nil},
		Zip:      pgtype.Text{String: util.Deref(req.Zip), Valid: req.Zip != nil},
		Lat:      pgtype.Float4{Float32: util.Deref(req.Lat), Valid: req.Lat != nil},
		Long:     pgtype.Float4{Float32: util.Deref(req.Long), Valid: req.Long != nil},
	})
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(newAddressResponse(addr))
}

// getAddressByID gets an address by ID.
// @Summary Get an address by ID
// @Description Get an address by ID
// @Tags addresses
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Success 200 {object} addressResponse
// @Failure 400 {object} util.CustomError
// @Failure 404 {object} util.CustomError
// @Router /addresses/{id} [get]
// @Security BearerAuth
func (r *AddressRouter) getAddressByID(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	a, err := r.store.GetAddress(c.Context(), id)
	if err != nil {
		return util.ErrItemNotFound("address")
	}
	return c.JSON(newAddressResponse(a))
}

// listAddresses lists addresses by client ID.
// @Summary List addresses by client ID
// @Description List addresses by client ID
// @Tags addresses
// @Accept json
// @Produce json
// @Param client_id query int true "Client ID"
// @Success 200 {array} addressResponse
// @Failure 400 {object} util.CustomError
// @Failure 404 {object} util.CustomError
// @Router /addresses [get]
// @Security BearerAuth
func (r *AddressRouter) listAddresses(c *fiber.Ctx) error {
	idStr := c.Query("client_id")
	if idStr == "" {
		return util.ErrBadRequest("client_id is required")
	}
	clientID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return util.ErrBadRequest("invalid client_id")
	}

	list, err := r.store.ListAddresses(c.Context(), clientID)
	if err != nil {
		return err
	}

	resp := make([]addressResponse, 0, len(list))
	for _, a := range list {
		resp = append(resp, newAddressResponse(a))
	}
	return c.JSON(resp)
}

// updateAddressByID updates an address by ID.
// @Summary Update an address by ID
// @Description Update an address by ID
// @Tags addresses
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Param request body updateAddressRequest true "Address Info"
// @Success 200 {object} addressResponse
// @Failure 400 {object} util.CustomError
// @Failure 404 {object} util.CustomError
// @Failure 422 {object} util.ValidateErrors
// @Router /addresses/{id} [put]
// @Security BearerAuth
func (r *AddressRouter) updateAddressByID(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)

	var req updateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return util.ErrInvalidBody(err)
	}
	if err := r.validate.Struct(&req); err != nil {
		return util.NewValidationErrors(err.(validator.ValidationErrors))
	}

	a, err := r.store.UpdateAddress(c.Context(), db.UpdateAddressParams{
		ID:     id,
		Title:  req.Title,
		City:   req.City,
		Street: pgtype.Text{String: util.Deref(req.Street), Valid: req.Street != nil},
		Phone:  pgtype.Text{String: util.Deref(req.Phone), Valid: req.Phone != nil},
		Zip:    pgtype.Text{String: util.Deref(req.Zip), Valid: req.Zip != nil},
		Lat:    pgtype.Float4{Float32: util.Deref(req.Lat), Valid: req.Lat != nil},
		Long:   pgtype.Float4{Float32: util.Deref(req.Long), Valid: req.Long != nil},
	})
	if err != nil {
		return err
	}
	return c.JSON(newAddressResponse(a))
}

// deleteAddress deletes an address by ID.
// @Summary Delete an address by ID
// @Description Delete an address by ID
// @Tags addresses
// @Accept json
// @Produce json
// @Param id path int true "Address ID"
// @Success 204
// @Failure 400 {object} util.CustomError
// @Failure 404 {object} util.CustomError
// @Router /addresses/{id} [delete]
// @Security BearerAuth
func (r *AddressRouter) deleteAddress(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	if err := r.store.DeleteAddress(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}
