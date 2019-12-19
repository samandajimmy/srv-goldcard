package models

// Requirements a struct to store all payload for a list response
type Requirements struct {
	AktivasiFinansial string  `json:"aktivasi_finansial"`
	KYC               string  `json:"kyc"`
	LimitPengajuanMax int64   `json:"limit_pengajuan_max"`
	LimitPengajuanMin int64   `json:"limit_pengajuan_min"`
	OpenTE            string  `json:"open_te"`
	RegistrasiGTE     string  `json:"registrasi_gte"`
	SaldoMinEfektif   float64 `json:"saldo_min_efektif"`
	SaldoTabunganEmas string  `json:"saldo_tabungan_emas"`
	Umur              int64   `json:"umur"`
}

// ValList represent value of Requirements struct
var ValList = Requirements{
	AktivasiFinansial: "true",
	KYC:               "true",
	LimitPengajuanMax: 999000000,
	LimitPengajuanMin: 3000000,
	OpenTE:            "true",
	RegistrasiGTE:     "true",
	SaldoMinEfektif:   0.1,
	SaldoTabunganEmas: "true",
	Umur:              21,
}
