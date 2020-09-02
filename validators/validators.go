package validators

type Register struct {
	Account       string `validate:"required,max=20" comment:"帳號"`
	Password      string `validate:"required,max=20" comment:"密碼"`
	PasswordCheck string `validate:"required,max=20,eqfield=Password" comment:"密碼確認"`
	Email         string `validate:"required,email" comment:"電子郵件"`
	//Photo         string `validate:"max=100,image" comment:"大頭貼"`
}

type Login struct {
	Account  string `validate:"required,max=20" comment:"帳號"`
	Password string `validate:"required,min=6,max=20" comment:"密碼"`
}

type Img2ascii struct {
	Public bool  `comment:"是否公開"`
	Col    int   `validate:"required" comment:"欄數"`
}

type UpdateHot struct {
	AsciiArtId uint32 `validate:"required" comment:"asciiID"`
}

type UserAscii struct {
	UsertId uint32 `validate:"required" comment:"UsertId"`
}
