package models

import (
	"fmt"
	"gade/srv-goldcard/logger"
	"strconv"
	"strings"
	"time"

	"github.com/leekchan/accounting"
)

// PdsNotification is a struct to store PdsNotification data
type PdsNotification struct {
	PhoneNumber             string        `json:"phoneNumber"`
	CIF                     string        `json:"cif"`
	EmailSubject            string        `json:"emailSubject"`
	ContentTitle            string        `json:"contentTitle"`
	ContentDescription      []string      `json:"contentDescription"`
	ContentFooter           []string      `json:"contentFooter"`
	ContentList             []ContentList `json:"contentList"`
	NotificationTitle       string        `json:"notificationTitle"`
	NotificationDescription string        `json:"notificationDescription"`
}

type ContentList struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (pdsNotif *PdsNotification) GcApplication(acc Account, notifType string) {
	pdsNotif.PhoneNumber = acc.PersonalInformation.HandPhoneNumber
	pdsNotif.CIF = acc.CIF
	switch notifType {
	case "failed":
		pdsNotif.EmailSubject = "Pengajuan Kartu Emas Pegadaian Gagal"
		pdsNotif.ContentTitle = "Pengajuan Kartu Emas Pegadaian Gagal"
		pdsNotif.ContentDescription = []string{"Mohon maaf Pengajuan Kartu Emas anda belum berhasil."}
		pdsNotif.ContentFooter = []string{"Silahkan untuk melakukan pengajuan kembali dengan melengkapi data-data yang sesuai."}
		pdsNotif.NotificationDescription = "Pengajuan Gagal"
	case "succeeded":
		pdsNotif.EmailSubject = "Pengajuan Kartu Emas Pegadaian Berhasil"
		pdsNotif.ContentTitle = "Pengajuan Kartu Emas Pegadaian Berhasil"
		pdsNotif.ContentDescription = []string{"Selamat Pengajuan Kartu Emas anda Berhasil"}
		pdsNotif.ContentFooter = []string{"Kartu Emas akan segera diproses, pengiriman kartu maksimal 14 hari kerja."}
		pdsNotif.NotificationDescription = "Pengajuan Berhasil"
	default:
		logger.Make(nil, nil).Fatal("notifType could not be empty")
	}

	ac := accounting.Accounting{Symbol: "Rp. ", Precision: 2}
	pdsNotif.NotificationTitle = "Kartu Emas"
	pdsNotif.ContentList = []ContentList{
		{
			Key:   "Tanggal",
			Value: acc.Application.ApplicationProcessedDate.Format("02/01/2006"),
		},
		{
			Key:   "Waktu",
			Value: acc.Application.ApplicationProcessedDate.Format("15:04"),
		},
		{
			Key:   "Referensi",
			Value: acc.Application.ApplicationNumber,
		},
		{
			Key:   "No Rekening Tabungan Emas",
			Value: acc.Application.SavingAccount,
		},
		{
			Key: "Nama Nasabah",
			Value: strings.Join([]string{acc.PersonalInformation.FirstName,
				acc.PersonalInformation.LastName}, " "),
		},
		{
			Key:   "Harga/0.01gr",
			Value: ac.FormatMoney(acc.Card.StlLimit),
		},
		{
			Key:   "Pengajuan Gram Limit",
			Value: fmt.Sprintf("%f gram", acc.Card.GoldLimit),
		},
		{
			Key:   "Pengajuan Limit",
			Value: ac.FormatMoney(acc.Card.CardLimit),
		},
	}
}

func (pdsNotif *PdsNotification) GcActivation(acc Account, notifType string) {
	actDate := acc.Card.UpdatedAt // TODO should change this with activation date
	ac := accounting.Accounting{Symbol: "Rp. ", Precision: 2}
	pdsNotif.NotificationTitle = "Kartu Emas"
	pdsNotif.PhoneNumber = acc.PersonalInformation.HandPhoneNumber
	pdsNotif.CIF = acc.CIF

	switch notifType {
	case "failed":
		pdsNotif.EmailSubject = "Aktivasi Kartu Emas Pegadaian Gagal"
		pdsNotif.ContentTitle = pdsNotif.EmailSubject
		pdsNotif.ContentDescription = []string{"Pastikan data pada kartu dan PIN aplikasi Pegadaian anda sesuai dengan yang anda input."}
		pdsNotif.NotificationDescription = "Aktivasi Gagal"
	case "succeeded":
		pdsNotif.EmailSubject = "Aktivasi Kartu Emas Pegadaian Berhasil"
		pdsNotif.ContentTitle = pdsNotif.EmailSubject
		pdsNotif.ContentDescription = []string{"Selamat Aktivasi Kartu Emas Berhasil"}
		pdsNotif.ContentFooter = []string{"Segera lakukan pergantian PIN Kartu Emas Fisik anda dengan cara :",
			"PIN(spasi)KK(spasi)6 digit pertama No. KK BRI#4 Digit Terakhir No. KK BRI#tgl lahir ddmmyyyy",
			"Contoh:", "PIN KK 518828#1234#17081945", "Kirim ke 3300 melalui nomor Handphone yang terdaftar",
			"Info PIN https://bit.ly/xxx"}
		pdsNotif.NotificationDescription = "Aktivasi Berhasil"
		pdsNotif.ContentList = []ContentList{
			{
				Key:   "Tanggal",
				Value: actDate.Format("02/01/2006"),
			},
			{
				Key:   "Waktu",
				Value: actDate.Format("15:04"),
			},
			{
				Key:   "No",
				Value: acc.Card.CardNumber,
			},
			{
				Key: "Nama Nasabah",
				Value: strings.Join([]string{acc.PersonalInformation.FirstName,
					acc.PersonalInformation.LastName}, " "),
			},
			{
				Key:   "Limit Kartu",
				Value: ac.FormatMoney(acc.Card.CardLimit),
			},
		}
	default:
		logger.Make(nil, nil).Fatal("notifType could not be empty")
	}
}

// MaskChar to replace string with x
func (pdsNotif *PdsNotification) MaskChar(txt string, strt int, end int) string {
	var x string
	for i := 0; i < end-strt; i++ {
		x += "X"
	}
	runes := []rune(txt)
	return strings.ReplaceAll(txt, string(runes[strt:end]), x)
}

// TimeParser to parse date string to date and time
func (pdsNotif *PdsNotification) TimeParser(param string) (string, string) {
	trxDate, _ := time.Parse("2006-01-02 15:04:05", param)
	return trxDate.Format("02/01/2006"), trxDate.Format("15:04:05")
}

// GcTransaction to send goldcard transaction notification to pds
func (pdsNotif *PdsNotification) GcTransaction(trx Transaction) {
	ac := accounting.Accounting{Symbol: "Rp ", Thousand: "."}
	trxDate, trxTime := pdsNotif.TimeParser(trx.TrxDate)

	pdsNotif.NotificationTitle = "Kartu Emas"
	pdsNotif.PhoneNumber = trx.Account.PersonalInformation.HandPhoneNumber
	pdsNotif.CIF = trx.Account.CIF
	pdsNotif.ContentTitle = "Transaksi Kartu Emas Berhasil"
	pdsNotif.ContentDescription = []string{"Transaksi Kartu Emas berhasil. Berikut merupakan rincian transaksi kamu."}
	pdsNotif.NotificationDescription = "Transaksi " + ac.FormatMoney(trx.Nominal) + " Berhasil"
	pdsNotif.ContentList = []ContentList{
		{
			Key:   "Tanggal",
			Value: trxDate,
		},
		{
			Key:   "Waktu",
			Value: trxTime,
		},
		{
			Key:   "Id Transaksi",
			Value: pdsNotif.MaskChar(trx.RefTrx, 6, 12),
		},
		{
			Key:   "Jumlah Transaksi",
			Value: ac.FormatMoney(trx.Nominal),
		},
		{
			Key:   "Tempat",
			Value: trx.Description,
		},
	}
}

// GcPayment to send goldcard payment notification to pds
func (pdsNotif *PdsNotification) GcPayment(trx Transaction, bill Billing, pind PaymentInquiryNotificationData) {
	ac := accounting.Accounting{Symbol: "Rp ", Thousand: "."}
	time := time.Now().UTC()

	administration, _ := strconv.ParseInt(pind.Administration, 10, 64)
	dateTime := time.Format("02/01/2006") + " " + time.Format("15:04:05")
	totalPayment := trx.Nominal + administration
	pdsNotif.NotificationTitle = "Kartu Emas"
	pdsNotif.EmailSubject = "Transaksi Pembayaran Kartu Emas Sukses"
	pdsNotif.PhoneNumber = trx.Account.PersonalInformation.HandPhoneNumber
	pdsNotif.CIF = trx.Account.CIF
	pdsNotif.ContentTitle = "Transaksi Pembayaran Kartu Emas Sukses"
	pdsNotif.ContentDescription = []string{"Terima kasih telah melakukan pembayaran.", "Transaksi Kamu Berhasil:"}
	pdsNotif.NotificationDescription = "Terima kasih telah melakukan pembayaran."
	pdsNotif.ContentList = []ContentList{
		{
			Key:   "Metode Pembayaran",
			Value: "",
		},
		{
			Key:   "Referensi",
			Value: pind.ReffSwitching,
		},
		{
			Key:   "Jenis Transaksi",
			Value: "Pembayaran Kartu Emas",
		},
		{
			Key:   "Waktu Transaksi",
			Value: dateTime,
		},
		{
			Key:   "Biaya Channel",
			Value: ac.FormatMoney(administration),
		},
		{
			Key:   "Total Pembayaran",
			Value: ac.FormatMoney(totalPayment),
		},
		{
			Key:   "Sisa Tagihan",
			Value: ac.FormatMoney(bill.DebtAmount),
		},
	}
}

// GcDecreasedSTL to send new card limit notification in cause of decreased STL
func (pdsNotif *PdsNotification) GcDecreasedSTL(acc Account, oldCard Card, refTrx string) {
	ac := accounting.Accounting{Symbol: "Rp ", Thousand: "."}
	trxDate, trxTime := pdsNotif.TimeParser(time.Now().Format(DateTimeFormat))

	pdsNotif.NotificationTitle = "Kartu Emas"
	pdsNotif.PhoneNumber = acc.PersonalInformation.HandPhoneNumber
	pdsNotif.CIF = acc.CIF
	pdsNotif.ContentTitle = "Informasi Limit Kartu Emas"
	pdsNotif.ContentDescription = []string{"Karena ada penurunan harga emas yang sangat signifikan, berikut informasi perubahan limit Kartu Emas kamu saat ini"}
	pdsNotif.NotificationDescription = "Informasi Limit Kartu Emas"
	pdsNotif.ContentList = []ContentList{
		{
			Key:   "Tanggal",
			Value: trxDate,
		},
		{
			Key:   "Waktu",
			Value: trxTime[:5],
		},
		{
			Key:   "Referensi",
			Value: refTrx,
		},
		{
			Key:   "Limit Lama",
			Value: ac.FormatMoney(oldCard.CardLimit),
		},
		{
			Key:   "Saldo Gram Lama",
			Value: fmt.Sprintf("%.4f", oldCard.GoldLimit) + " gram",
		},
		{
			Key:   "Limit Baru",
			Value: ac.FormatMoney(acc.Card.CardLimit),
		},
		{
			Key:   "Harga/0.01gr",
			Value: ac.FormatMoney(acc.Card.StlLimit / 100),
		},
		{
			Key:   "Pengajuan Gram Limit",
			Value: fmt.Sprintf("%.4f", acc.Card.GoldLimit) + " gram",
		},
	}
}

// GcSla2Days to send  notification about update limit's sla
func (pdsNotif *PdsNotification) GcSla2Days(acc Account) {
	pdsNotif.NotificationTitle = "KARTU EMAS"
	pdsNotif.PhoneNumber = acc.PersonalInformation.HandPhoneNumber
	pdsNotif.CIF = acc.CIF
	pdsNotif.EmailSubject = "Pengajuan Limit Diproses"
	pdsNotif.ContentTitle = "Pengajuan Limit Diproses"
	pdsNotif.ContentDescription = []string{"Proses pengajuan limit dapat berlangsung hingga 2 hari kerja."}
	pdsNotif.NotificationDescription = "Pengajuan Limit Diproses"
}
