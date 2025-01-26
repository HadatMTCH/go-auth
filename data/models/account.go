package models

type Account struct {
	NoRekening string `db:"no_rekening"`
	Nama       string `db:"nama"`
	Nik        string `db:"nik"`
	NoHp       string `db:"no_hp"`
	Saldo      int64  `db:"saldo"`
}

type Transaction struct {
	ID         int    `db:"id"`
	NoRekening string `db:"no_rekening"`
	Nominal    int64  `db:"nominal"`
	Jenis      string `db:"jenis"`
	Waktu      string `db:"waktu"`
}

// Request/Response DTOs
type RegisterRequest struct {
	Nama string `json:"nama"`
	Nik  string `json:"nik"`
	NoHp string `json:"no_hp"`
}

type RegisterResponse struct {
	NoRekening string `json:"no_rekening"`
}

type TabungTarikRequest struct {
	NoRekening string `json:"no_rekening"`
	Nominal    int64  `json:"nominal"`
}

type SaldoResponse struct {
	Saldo int64 `json:"saldo"`
}

type ErrorResponse struct {
	Remark string `json:"remark"`
}
