package models

import (
	"errors"
	"fmt"
)

var (
	// ErrInternalServerError to store internal server error message
	ErrInternalServerError = errors.New("Internal Server Error ")

	// ErrNotFound to store not found error message
	ErrNotFound = errors.New("Item tidak ditemukan")

	// ErrConflict to store conflicted error message
	ErrConflict = errors.New("Item sudah ada")

	// ErrBadParamInput to store bad parameter error message
	ErrBadParamInput = errors.New("Parameter yang diberikan tidak valid")

	// ErrGetCampaignCounter to get campaign counter error message
	ErrGetCampaignCounter = errors.New("Campaign tidak tersedia")

	// ErrRefIDStatus to not found ref_trx error message
	ErrRefIDStatus = errors.New("Transaksi ID ")

	// ErrReqUndefined to store undefine request body error message
	ErrReqUndefined = errors.New("Request data tidak tersedia")

	// ErrFileProdReqsNotFound to errors file production requirements not found
	ErrFileProdReqsNotFound = errors.New(("File Product Requirements Tidak di Temukan"))

	// ErrAddressEmpty to errors address is empty
	ErrAddressEmpty = errors.New(("Alamat tidak boleh kosong"))

	// ErrPostAddressFailed to errors post address failed
	ErrPostAddressFailed = errors.New(("Menambahkan alamat gagal"))
)

// DynamicErr to return parameterize errors
func DynamicErr(message string, args ...string) error {
	return fmt.Errorf(message, args[0])
}
