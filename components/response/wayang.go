package response

import "time"

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type MsgErrorWayang struct {
	Message string `json:"message"`
}

type Function struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Label    string `json:"label"`
	ModuleID int    `json:"module_id"`
}
type Module struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Label    string     `json:"label"`
	Function []Function `json:"function"`
}
type Role struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Status      int      `json:"status"`
	Description string   `json:"description"`
	Module      []Module `json:"module"`
}
type Client struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Label       string `json:"label"`
	Key         string `json:"key"`
	Secret      string `json:"secret"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}
type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Username    string `json:"username"`
	FcmDeviceID string `json:"fcm_device_id"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Role        Role   `json:"role"`
	Photo       string `json:"photo"`
	Client      Client `json:"client"`
}

type Profile struct {
	HasPrevious        bool        `json:"has_previous"`
	HasNext            bool        `json:"has_next"`
	PreviousPageNumber interface{} `json:"previous_page_number"`
	NextPageNumber     interface{} `json:"next_page_number"`
	CurrentPage        int         `json:"current_page"`
	PerPage            int         `json:"per_page"`
	TotalPage          int         `json:"total_page"`
	TotalData          int         `json:"total_data"`
	Data               []Data      `json:"data"`
}
type Data struct {
	Nik            string `json:"nik"`
	CompleteName   string `json:"complete_name"`
	Birthplace     string `json:"birthplace"`
	MotherName     string `json:"mother_name"`
	FatherName     string `json:"father_name"`
	NoKk           string `json:"no_kk"`
	DaerahDomisili string `json:"daerah_domisili"`
}

type ProfileKTP struct {
	Nik            string    `json:"nik"`
	NamaLengkap    string    `json:"nama_lengkap"`
	NomorKk        string    `json:"nomor_kk"`
	StatHubKel     string    `json:"stat_hub_kel"`
	TempatLahir    string    `json:"tempat_lahir"`
	TempatTanggal  time.Time `json:"tempat_tanggal"`
	JenisKelamin   string    `json:"jenis_kelamin"`
	Agama          string    `json:"agama"`
	Alamat         string    `json:"alamat"`
	StatusKawin    string    `json:"status_kawin"`
	JenisPekerjaan string    `json:"jenis_pekerjaan"`
	Rt             int       `json:"rt"`
	Rw             int       `json:"rw"`
	Propinsi       string    `json:"propinsi"`
	Kabupaten      string    `json:"kabupaten"`
	Kecamatan      string    `json:"kecamatan"`
	Keluarahan     string    `json:"keluarahan"`
	NoAktaLahir    string    `json:"no_akta_lahir"`
	FlagStatus     string    `json:"flag_status"`
	TanggalEntri   string    `json:"tanggal_entri"`
	TanggalUbah    string    `json:"tanggal_ubah"`
	Foto           string    `json:"foto"`
}

type ProfileKK []struct {
	AddressDistrict     string        `json:"address_district"`
	FamilyConection     string        `json:"family_conection"`
	BirthPlace          string        `json:"birth_place"`
	AddressSubDistrict  string        `json:"address_sub_district"`
	FathersName         string        `json:"fathers_name"`
	AddressProvince     string        `json:"address_province"`
	MothersName         string        `json:"mothers_name"`
	AddressUrbanVillage string        `json:"address_urban_village"`
	FamilyID            string        `json:"family_id"`
	AddressRt           string        `json:"address_rt"`
	Nik                 string        `json:"nik"`
	Address             string        `json:"address"`
	Telp                string        `json:"telp"`
	Gender              string        `json:"gender"`
	MaritalStatus       string        `json:"marital_status"`
	AddressRw           string        `json:"address_rw"`
	FullName            string        `json:"full_name"`
	PostalCode          string        `json:"postal_code"`
	AddressVillage      string        `json:"address_village"`
	BirthDate           string        `json:"birth_date"`
	Photo               string        `json:"photo"`
	Weton               string        `json:"weton"`
	Zodiac              string        `json:"zodiac"`
	Phone               []interface{} `json:"phone"`
}

type ProfileNikByPhone []struct {
	Nik               string `json:"nik"`
	Msisdn            string `json:"msisdn"`
	Operator          string `json:"operator"`
	Status            string `json:"status"`
	TanggalRegistrasi string `json:"tanggal_registrasi"`
	CreatedTime       string `json:"created_time"`
}

type ProfilePhone []struct {
	Nik               string `json:"nik"`
	Msisdn            string `json:"msisdn"`
	Operator          string `json:"operator"`
	Status            string `json:"status"`
	TanggalRegistrasi string `json:"tanggal_registrasi"`
}

type ProfileSTNK struct {
	NOPOLISI     interface{} `json:"NO_POLISI"`
	PEMILIK      interface{} `json:"PEMILIK"`
	ALAMAT       interface{} `json:"ALAMAT"`
	ALAMAT2      interface{} `json:"ALAMAT2"`
	KODELOKASI   interface{} `json:"KODE_LOKASI"`
	MERK         interface{} `json:"MERK"`
	JENISKEND    interface{} `json:"JENIS_KEND"`
	GOLKEND      interface{} `json:"GOL_KEND"`
	TIPE         interface{} `json:"TIPE"`
	MODEL        interface{} `json:"MODEL"`
	RAKIT        float64     `json:"RAKIT"`
	SILINDER     interface{} `json:"SILINDER"`
	WARNA        interface{} `json:"WARNA"`
	NORANGKA     interface{} `json:"NO_RANGKA"`
	NOMESIN      interface{} `json:"NO_MESIN"`
	TNKB         interface{} `json:"TNKB"`
	BBM          interface{} `json:"BBM"`
	NOPOLLAMA    interface{} `json:"NOPOL_LAMA"`
	JMLAS        interface{} `json:"JML_AS"`
	JMLBEBAN     interface{} `json:"JML_BEBAN"`
	STSPEMILIK   interface{} `json:"STSPEMILIK"`
	NOBPKB       interface{} `json:"NO_BPKB"`
	TGLBPKB      interface{} `json:"TGL_BPKB"`
	NOFAKTUR     interface{} `json:"NO_FAKTUR"`
	TGLFAKTUR    interface{} `json:"TGL_FAKTUR"`
	TGLSTNK      interface{} `json:"TGL_STNK"`
	TGLPAJAK     interface{} `json:"TGL_PAJAK"`
	TGLDAFTAR    interface{} `json:"TGL_DAFTAR"`
	JNSDAFTAR    interface{} `json:"JNS_DAFTAR"`
	KDSAMSAT     interface{} `json:"KD_SAMSAT"`
	KDPOLDA      interface{} `json:"KD_POLDA"`
	NOKTP        interface{} `json:"NO_KTP"`
	NOKK         interface{} `json:"NO_KK"`
	NOHP         interface{} `json:"NO_HP"`
	TGLMATIYAD   interface{} `json:"TGL_MATI_YAD"`
	SAMSATBAYAR  interface{} `json:"SAMSAT_BAYAR"`
	KODEJENIS    interface{} `json:"KODE_JENIS"`
	KODEGOLONGAN interface{} `json:"KODE_GOLONGAN"`
}

type MasterProfile struct {
	KTP   ProfileKTP   `json:"ktp"`
	KK    ProfileKK    `json:"kk"`
	Phone ProfilePhone `json:"phone"`
}
