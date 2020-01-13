package types

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValidationError struct {
	fieldErrors []validator.FieldError
}

func NewValidationError(fieldErrors []validator.FieldError) error {
	return ValidationError{
		fieldErrors: fieldErrors,
	}
}

func (err ValidationError) Error() string {
	return "struct validation failed"
}

func (err ValidationError) FieldErrors() []validator.FieldError {
	return err.fieldErrors
}

// ValidationError creates an 'InvalidArgument' response and attaches all 'FieldViolation' structs.
func (err ValidationError) GrpcStatus() *status.Status {
	var violations []*errdetails.BadRequest_FieldViolation

	for _, violation := range err.FieldErrors() {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       violation.StructNamespace(),
			Description: strings.ToUpper(violation.Tag()),
		})
	}

	grpcStatus := status.Newf(codes.InvalidArgument, err.Error())
	grpcStatus, statErr := grpcStatus.WithDetails(&errdetails.BadRequest{
		FieldViolations: violations,
	},
	)

	if statErr != nil {
		panic(fmt.Sprintf("failed to marshal error details: %s", statErr))
	}

	return grpcStatus
}
