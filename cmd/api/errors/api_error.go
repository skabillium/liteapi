package errors

type ApiError struct {
	// HttpCode int
	Message string
	Error   string
}

type SupplierError struct {
	ApiError
	Supplier SupplierData `json:"supplier"`
}

type SupplierData struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}
