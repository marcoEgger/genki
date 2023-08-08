package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"io"
)

// PdfMarshaller Create a custom marshaller for requests with header "Accept" set to application/pdf
type PdfMarshaller struct{}

func (m *PdfMarshaller) Marshal(_ interface{}) ([]byte, error) {
	// We don't want to write anything, because it is already written by response interceptor
	return nil, nil
}

func (m *PdfMarshaller) Unmarshal(_ []byte, _ interface{}) error {
	// We don't unmarshal PDFs
	return nil
}

func (m *PdfMarshaller) NewDecoder(_ io.Reader) runtime.Decoder {
	panic("not implemented")
}

func (m *PdfMarshaller) NewEncoder(_ io.Writer) runtime.Encoder {
	panic("not implemented")
}

func (m *PdfMarshaller) ContentType(_ interface{}) string {
	return "application/pdf"
}
