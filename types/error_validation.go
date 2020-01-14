package types

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValidationError struct {
	fieldErrors validator.ValidationErrors
	validator   *validator.Validate
	translation ut.Translator
}

func NewValidationError(errors validator.ValidationErrors, validator *validator.Validate) error {

	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	if err := en_translations.RegisterDefaultTranslations(validator, trans); err != nil {
		return err
	}

	return ValidationError{
		fieldErrors: errors,
		validator:   validator,
		translation: trans,
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
		var message string

		message = violation.Translate(err.translation)

		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       violation.StructNamespace(),
			Description: message,
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
