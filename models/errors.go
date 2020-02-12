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
	ErrTokenExpired = errors.New("Token Anda telah kadarluarsa")

	// ErrPostSavingAccountFailed to errors post address failed
	ErrPostSavingAccountFailed = errors.New("Menambahkan Rekening Tabungan gagal")

	// ErrCreateApplication to errors create application failed
	ErrCreateApplication = errors.New("Terjadi kesalahan pada proses pembuatan pengajuan")

	// ErrAppNumberNotFound to store error find application number
	ErrAppNumberNotFound = errors.New("Nomor pengajuan tidak ditemukan")

	// ErrAppNumberCompleted to store error application number applied
	ErrAppNumberCompleted = errors.New("Nomor pengajuan telah diajukan")

	// ErrBankNotFound to store error find bank id
	ErrBankNotFound = errors.New("Kode bank tidak ditemukan")

	// ErrEmergecyContactNotFound to store emergency contact not found error message
	ErrEmergecyContactNotFound = errors.New("Emergency contact tidak ditemukan")

	// ErrUpdateRegData to store updating registration data error message
	ErrUpdateRegData = errors.New("Terjadi kesalahan saat update data pengajuan")

	// ErrUpdateOccData to store updating occupation data error message
	ErrUpdateOccData = errors.New("Terjadi kesalahan saat update data pekerjaan")

	// ErrUpdateCardLimit to store updating card limit data error message
	ErrUpdateCardLimit = errors.New("Terjadi kesalahan saat update card limit")

	// ErrZipcodeNotFound to store zip code not found error message
	ErrZipcodeNotFound = errors.New("Kode pos tidak ditemukan")

	// ErrBlacklisted to store blacklisted error message
	ErrBlacklisted = errors.New("Anda tercatat dalam daftar black list kami")

	// ErrUpdateBrixkey to store update brixkey error message
	ErrUpdateBrixkey = errors.New("Terjadi kesalahan saat update brixkey")

	// ErrUpdateBrixkey to store update brixkey error message
	ErrGetAccByBrixkey = errors.New("Brixkey tidak di temukan")

	// ErrGetAccByAccountNumber to store account not found error message
	ErrGetAccByAccountNumber = errors.New("Account Number tidak di temukan")

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

	// ErrGetAppStatus to strore update application status error message
	ErrGetAppStatus = errors.New("Terjadi Kesalahan saat mengambil data status aplikasi")

	// ErrInquiryReg to inquiry registrations to switching/core
	ErrInquiryReg = errors.New("Data Pengajuan sudah pernah terdaftar")

	// ErrSwitchingAPIRequest to store switching api request error message
	ErrSwitchingAPIRequest = "SWITCHING API: RC-%s - %s"

	// ErrAppData to store get application data error message
	ErrAppData = errors.New("Terjadi Kesalahan saat mengambil data aplikasi")

	// ErrPostActivationsFailed to errors post activations failed
	ErrPostActivationsFailed = errors.New("Gagal melakukan aktivasi")

	// ErrAlreadyActivated to errors already activated
	ErrAlreadyActivated = errors.New("Akun ini sudah pernah di aktivasi sebelumnya")

	// ErrStatusActivations to store activations status activation status not "sent" yet
	ErrStatusActivations = errors.New("Status pengajuan tidak sesuai")

	// ErrAppExpired to store application expired error message
	ErrAppExpired = errors.New("PENGAJUAN KADALUARSA")

	// ErrAppExpiredDesc to store the description of application expired error message
	ErrAppExpiredDesc = errors.New("Pengajuan harus dibatalkan karena tidak ada aktivitas selama 12 bulan. Saldo emas akan dikembalikan ke saldo efektif.")

	// ErrGetCurrSTL to store get current STL error message
	ErrGetCurrSTL = errors.New("Terjadi kesalahan ketika mendapatkan harga emas saat ini")

	// ErrGetUserDetail to store get user detail error message
	ErrGetUserDetail = errors.New("Terjadi kesalahan ketika mendapatkan data detail nasabah")

	// ErrGetEffBalance to store get effective gold balance error message
	ErrGetEffBalance = errors.New("Terjadi kesalahan ketika mendapatkan data saldo efektif nasabah")

	// ErrDecreasedSTL to store get decreasing STL error message
	ErrDecreasedSTL = errors.New("HARGA EMAS TURUN")

	// ErrDecreasedSTLDesc to store the description of get decreasing STL error message
	ErrDecreasedSTLDesc = errors.New("Harga emas turun cukup tinggi sejak kamu mengajukan kartu emas. Top Up Tabungan Emas kamu untuk melanjutkan proses aktivasi.")

	// ErrInsertTransactions to store get failed when insert data to table transactions
	ErrInsertTransactions = errors.New("Gagal saat memasukan data transaksi")

	// ErrGetHistoryTransactions to store message failed when get data history account in table transactions
	ErrGetHistoryTransactions = errors.New("Gagal saat mencari data history transaksi")

	// ErrCoreEODStatus to store down time core servive
	ErrCoreEODStatus = errors.New("Mohon maaf, Layanan sedang tidak tersedia")

	// ErrGetCardBalance to store get gold card balance error message
	ErrGetCardBalance = errors.New("Terjadi kesalahan ketika mendapatkan data saldo kartu emas")
)

// DynamicErr to return parameterize errors
func DynamicErr(message string, args []interface{}) error {
	return fmt.Errorf(message, args...)
}
