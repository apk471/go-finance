package validation

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apk471/go-boilerplate/internal/errs"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type strictJSONPayload struct {
	Name string `json:"name" validate:"required,min=2"`
}

func (p *strictJSONPayload) Validate() error {
	return validator.New().Struct(p)
}

type queryPayload struct {
	Limit int `query:"limit" validate:"min=1"`
}

func (p *queryPayload) Validate() error {
	return validator.New().Struct(p)
}

func TestBindAndValidateRejectsUnknownJSONFields(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Ayush","extra":"nope"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := BindAndValidate(c, &strictJSONPayload{})
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	httpErr, ok := err.(*errs.HTTPError)
	if !ok {
		t.Fatalf("expected *errs.HTTPError, got %T", err)
	}
	if httpErr.Status != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, httpErr.Status)
	}
	if len(httpErr.Errors) != 1 || httpErr.Errors[0].Field != "extra" {
		t.Fatalf("expected unknown field error for extra, got %#v", httpErr.Errors)
	}
}

func TestBindAndValidateRejectsMalformedJSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Ayush"`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := BindAndValidate(c, &strictJSONPayload{})
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	httpErr, ok := err.(*errs.HTTPError)
	if !ok {
		t.Fatalf("expected *errs.HTTPError, got %T", err)
	}
	if len(httpErr.Errors) != 1 || httpErr.Errors[0].Field != "body" {
		t.Fatalf("expected body error, got %#v", httpErr.Errors)
	}
}

func TestBindAndValidateRejectsInvalidJSONTypes(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":123}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := BindAndValidate(c, &strictJSONPayload{})
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	httpErr, ok := err.(*errs.HTTPError)
	if !ok {
		t.Fatalf("expected *errs.HTTPError, got %T", err)
	}
	if len(httpErr.Errors) != 1 || httpErr.Errors[0].Field != "name" {
		t.Fatalf("expected field error for name, got %#v", httpErr.Errors)
	}
}

func TestBindAndValidateValidatesQueryParams(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?limit=0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := BindAndValidate(c, &queryPayload{})
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	httpErr, ok := err.(*errs.HTTPError)
	if !ok {
		t.Fatalf("expected *errs.HTTPError, got %T", err)
	}
	if len(httpErr.Errors) != 1 || httpErr.Errors[0].Field != "limit" {
		t.Fatalf("expected field error for limit, got %#v", httpErr.Errors)
	}
}
