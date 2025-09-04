package main

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/tsoy/rental-rewards/internal/data"
	"net/http"
	"strings"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string  `json:"email"`
		FullName *string `json:"full_name"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Email == "" {
		app.badRequestResponse(w, r, errors.New("e-mail required!"))
		return
	}
	user := &data.User{
		Email:    input.Email,
		FullName: input.FullName,
	}

	err = app.models.Users.Insert(context.TODO(), user)

	if err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			app.badRequestResponse(w, r, err)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
}
func (app *application) createPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID      string `json:"user_id"`
		AmountCents int64  `json:"amount_cents"`
		Currency    string `json:"currency"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.AmountCents <= 0 {
		app.badRequestResponse(w, r, errors.New("e-mail required!"))
		return
	}

	currency := strings.ToUpper(strings.TrimSpace(input.Currency))
	//input.Currency
	if currency != "" {
		if _, ok := data.Currencies[input.Currency]; !ok {
			app.badRequestResponse(w, r, errors.New("invalid currency"))
			return
		}
	}

	userId, err := uuid.Parse(input.UserID)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("malformed user id"))
		return
	}

	var idkPtr *string
	idk := r.Header.Get("Idempotency-Key")
	if idk != "" {
		idkPtr = &idk
	}

	payment := data.Payment{
		UserID:      userId,
		AmountCents: input.AmountCents,
		Currency:    input.Currency,
		ExternalRef: idkPtr,
	}
	err = app.models.Payments.Insert(context.TODO(), &payment)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidUserId), errors.Is(err, data.ErrDuplicateTransaction):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"payment": payment}, nil)

	//fmt.Fprintln(w, "create a new payment")
}
func (app *application) getPaymentHandler(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	payment, err := app.models.Payments.Get(context.TODO(), id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r, err)
		} else {
			app.serverErrorResponse(w, r, err)

		}
		return
	}
	//p, err := app.models.Payments.Get(
	//
	//	p, err := repo.GetPayment(r.Context(), pid)
	//	if err != nil { http.Error(w, "not found", 404); return }
	//	writeJSON(w, 200, p)
	app.writeJSON(w, http.StatusOK, envelope{"payment": payment}, nil)

	//fmt.Fprintf(w, "payment info %d\n", id)
}
