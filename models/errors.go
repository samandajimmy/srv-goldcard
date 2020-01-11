package models

import (
	"errors"
	"fmt"
)

var (
	// ErrInternalServerError to store internal server error message
	ErrInternalServerError = errors.New("Internal Server Error")

	// ErrExternalAPI to store external api error message
	ErrExternalAPI = errors.New("External API Errors")

	// ErrSetVar to store setting variable error message
	ErrSetVar = errors.New("Setting variable error")

	// ErrNotFound to store not found error message
	ErrNotFound = errors.New("Item tidak ditemukan")

	// ErrConflict to store conflicted error message
	ErrConflict = errors.New("Item sudah ada")

	// ErrMappingData to store mapping data error message
	ErrMappingData = errors.New("Error mapping data")

	// ErrFileProdReqsNotFound to errors file production requirements not found
	ErrFileProdReqsNotFound = errors.New("File Product Requirements Tidak di Temukan")

	// ErrAddressEmpty to errors address is empty
	ErrAddressEmpty = errors.New("Alamat tidak boleh kosong")

	// ErrPostAddressFailed to errors post address failed
	ErrPostAddressFailed = errors.New("Menambahkan alamat gagal")

	// ErrUsername to store username error message
	ErrUsername = errors.New("Username atau Password yang digunakan tidak valid")

	// ErrPassword to store password error message
	ErrPassword = errors.New("Username atau Password yang digunakan tidak valid")

	// ErrTokenExpired to store password error message
	ErrTokenExpired = errors.New("Token Anda telah expire")

	// ErrPostSavingAccountFailed to errors post address failed
	ErrPostSavingAccountFailed = errors.New("Menambahkan Rekening Tabungan gagal")

	// ErrCreateApplication to errors create application failed
	ErrCreateApplication = errors.New("Terjadi kesalahan pada proses pembuatan pengajuan")

	// ErrAppNumberNotFound to store error find application number
	ErrAppNumberNotFound = errors.New("Nomor pengajuan tidak ditemukan")

	// ErrBankNotFound to store error find bank id
	ErrBankNotFound = errors.New("Kode bank tidak ditemukan")

	// ErrEmergecyContactNotFound to store emergency contact not found error message
	ErrEmergecyContactNotFound = errors.New("Emergency contact tidak ditemukan")

	// ErrUpdateRegData to store updating registration data error message
	ErrUpdateRegData = errors.New("Terjadi kesalahan saat update data pengajuan")

	// ErrUpdateCardLimit to store updating card limit data error message
	ErrUpdateCardLimit = errors.New("Terjadi kesalahan saat update card limit")

	// ErrZipcodeNotFound to store zip code not found error message
	ErrZipcodeNotFound = errors.New("Kode pos tidak ditemukan")

	// ErrBlacklisted to store blacklisted error message
	ErrBlacklisted = errors.New("Anda tercatat dalam daftar black list kami")

	// ErrUpdateBrixkey to store update brixkey error message
	ErrUpdateBrixkey = errors.New("Terjadi kesalahan saat update brixkey")

	// ErrUpdateAppDocID to store update application document ID error message
	ErrUpdateAppDocID = errors.New("Terjadi kesalahan saat update document id")

	// ErrDocIDNotFound to store document id not found error message
	ErrDocIDNotFound = errors.New("Document id tidak ditemukan")

	// ErrAppIDNotFound to store application id not found error message
	ErrAppIDNotFound = errors.New("Application id tidak ditemukan")

	// ErrBriAPIRequest to store bri api request error message
	ErrBriAPIRequest = "BRI API: RC-%s - %s"

	// ErrCreateToken to create/update token error message
	ErrCreateToken = errors.New("Terjadi Kesalaahan saat membuat Token")

	// ErrVerifyToken to verify token error message
	ErrVerifyToken = errors.New("Terjadi Kesalaahan saat verifikasi Token")

	// ErrUpdateAppStatus to update application status error message
	ErrUpdateAppStatus = errors.New("Terjadi Kesalahan saat update Status")
)

// DynamicErr to return parameterize errors
func DynamicErr(message string, args []interface{}) error {
	return fmt.Errorf(message, args...)
}
