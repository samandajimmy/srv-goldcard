package models

import (
	"fmt"
	"gade/srv-goldcard/logger"
	"strings"

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
			Value: ac.FormatMoney(acc.Card.CurrentSTL),
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
