package grpcHandler

import (
	"errors"

	"github.com/ernur-eskermes/product-store/pkg/filters"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorFilterResponse(err error) error {
	var validationErrors filters.ValidationErrors

	br := &errdetails.BadRequest{}

	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       e.Field,
				Description: e.Message,
			})
		}
	}

	st, err := status.New(codes.InvalidArgument, err.Error()).WithDetails(br)
	if err != nil {
		return err
	}

	return st.Err()
}
