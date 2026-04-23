package validationerror

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValidationError map[string]string

func (v ValidationError) Merge(other ValidationError) {
	for k, val := range other {
		v[k] = val
	}
}

func (v ValidationError) ToStatus() error {
	st := status.New(codes.InvalidArgument, "validation error")

	var details []*errdetails.BadRequest_FieldViolation
	for field, desc := range v {
		details = append(details, &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: desc,
		})
	}

	br := &errdetails.BadRequest{FieldViolations: details}

	stWithDetails, err := st.WithDetails(br)
	if err != nil {
		return st.Err()
	}

	return stWithDetails.Err()
}
