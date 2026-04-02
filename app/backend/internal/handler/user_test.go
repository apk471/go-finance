package handler

import "testing"

func TestUpdateUserRequestValidateRequiresAtLeastOneField(t *testing.T) {
	req := &UpdateUserRequest{
		ID: "11111111-1111-1111-1111-111111111111",
	}

	if err := req.Validate(); err == nil {
		t.Fatal("expected validation error but got nil")
	}
}
