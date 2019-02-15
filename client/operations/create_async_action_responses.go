// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/mcastelino/testapi/models"
)

// CreateAsyncActionReader is a Reader for the CreateAsyncAction structure.
type CreateAsyncActionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateAsyncActionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 204:
		result := NewCreateAsyncActionNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewCreateAsyncActionBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		result := NewCreateAsyncActionDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCreateAsyncActionNoContent creates a CreateAsyncActionNoContent with default headers values
func NewCreateAsyncActionNoContent() *CreateAsyncActionNoContent {
	return &CreateAsyncActionNoContent{}
}

/*CreateAsyncActionNoContent handles this case with default header values.

The update was successful
*/
type CreateAsyncActionNoContent struct {
}

func (o *CreateAsyncActionNoContent) Error() string {
	return fmt.Sprintf("[PUT /actions][%d] createAsyncActionNoContent ", 204)
}

func (o *CreateAsyncActionNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCreateAsyncActionBadRequest creates a CreateAsyncActionBadRequest with default headers values
func NewCreateAsyncActionBadRequest() *CreateAsyncActionBadRequest {
	return &CreateAsyncActionBadRequest{}
}

/*CreateAsyncActionBadRequest handles this case with default header values.

The action cannot be executed due to bad input
*/
type CreateAsyncActionBadRequest struct {
	Payload *models.Error
}

func (o *CreateAsyncActionBadRequest) Error() string {
	return fmt.Sprintf("[PUT /actions][%d] createAsyncActionBadRequest  %+v", 400, o.Payload)
}

func (o *CreateAsyncActionBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateAsyncActionDefault creates a CreateAsyncActionDefault with default headers values
func NewCreateAsyncActionDefault(code int) *CreateAsyncActionDefault {
	return &CreateAsyncActionDefault{
		_statusCode: code,
	}
}

/*CreateAsyncActionDefault handles this case with default header values.

Internal Server Error
*/
type CreateAsyncActionDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the create async action default response
func (o *CreateAsyncActionDefault) Code() int {
	return o._statusCode
}

func (o *CreateAsyncActionDefault) Error() string {
	return fmt.Sprintf("[PUT /actions][%d] createAsyncAction default  %+v", o._statusCode, o.Payload)
}

func (o *CreateAsyncActionDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}