package services

import (
	"aaimg2ascii/datamodels"
	"aaimg2ascii/models"
	"aaimg2ascii/validators"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os/exec"
	"strconv"
	"strings"
)

type Img2asciiService interface {
	UploadImg(file multipart.File, newfileName string) (string, bool)
	Img2ascii(newFileName string, img2ascii *validators.Img2ascii, userData *datamodels.User) (string, bool)
	UpdateHot(updateHot *validators.UpdateHot, userData *datamodels.User) string
	GetHotTop() []datamodels.AsciiArt
	GetUserAsciiById(userData *validators.UserAscii) datamodels.User
	GetMyAscii(userData *datamodels.User) datamodels.User
}

func NewImg2asciiService() Img2asciiService {
	return &img2asciiService{
		db:          models.DB.Mysql,
		baseService: NewBaseService(),
	}
}

type img2asciiService struct {
	db          *gorm.DB
	baseService BaseService
}

type AsciiData struct {
	Rows int    `json:"rows"`
	Cols int    `json:"cols"`
	Data string `json:"data"`
}

//region setblock
func setblock() map[string]string {
	return map[string]string{
		"0000": "⠀",
		"1000": "⠁",
		"0100": "⠂",
		"1100": "⠃",
		"0010": "⠄",
		"1010": "⠅",
		"0110": "⠆",
		"1110": "⠇",
		"2000": "⠈",
		"3000": "⠉",
		"2100": "⠊",
		"3100": "⠋",
		"2010": "⠌",
		"3010": "⠍",
		"2110": "⠎",
		"3110": "⠏",

		"0200": "⠐",
		"1200": "⠑",
		"0300": "⠒",
		"1300": "⠓",
		"0210": "⠔",
		"1210": "⠕",
		"0310": "⠖",
		"1310": "⠗",
		"2200": "⠘",
		"3200": "⠙",
		"2300": "⠚",
		"3300": "⠛",
		"2210": "⠜",
		"3210": "⠝",
		"2310": "⠞",
		"3310": "⠟",

		"0020": "⠠",
		"1020": "⠡",
		"0120": "⠢",
		"1120": "⠣",
		"0030": "⠤",
		"1030": "⠥",
		"0130": "⠦",
		"1130": "⠧",
		"2020": "⠨",
		"3020": "⠩",
		"2120": "⠪",
		"3120": "⠫",
		"2030": "⠬",
		"3030": "⠭",
		"2130": "⠮",
		"3130": "⠯",

		"0220": "⠰",
		"1220": "⠱",
		"0320": "⠲",
		"1320": "⠳",
		"0230": "⠴",
		"1230": "⠵",
		"0330": "⠶",
		"1330": "⠷",
		"2220": "⠸",
		"3220": "⠹",
		"2320": "⠺",
		"3320": "⠻",
		"2230": "⠼",
		"3230": "⠽",
		"2330": "⠾",
		"3330": "⠿",

		"0001": "⡀",
		"1001": "⡁",
		"0101": "⡂",
		"1101": "⡃",
		"0011": "⡄",
		"1011": "⡅",
		"0111": "⡆",
		"1111": "⡇",
		"2001": "⡈",
		"3001": "⡉",
		"2101": "⡊",
		"3101": "⡋",
		"2011": "⡌",
		"3011": "⡍",
		"2111": "⡎",
		"3111": "⡏",

		"0201": "⡐",
		"1201": "⡑",
		"0301": "⡒",
		"1301": "⡓",
		"0211": "⡔",
		"1211": "⡕",
		"0311": "⡖",
		"1311": "⡗",
		"2201": "⡘",
		"3201": "⡙",
		"2301": "⡚",
		"3301": "⡛",
		"2211": "⡜",
		"3211": "⡝",
		"2311": "⡞",
		"3311": "⡟",

		"0021": "⡠",
		"1021": "⡡",
		"0121": "⡢",
		"1121": "⡣",
		"0031": "⡤",
		"1031": "⡥",
		"0131": "⡦",
		"1131": "⡧",
		"2021": "⡨",
		"3021": "⡩",
		"2121": "⡪",
		"3121": "⡫",
		"2031": "⡬",
		"3031": "⡭",
		"2131": "⡮",
		"3131": "⡯",

		"0221": "⡰",
		"1221": "⡱",
		"0321": "⡲",
		"1321": "⡳",
		"0231": "⡴",
		"1231": "⡵",
		"0331": "⡶",
		"1331": "⡷",
		"2221": "⡸",
		"3221": "⡹",
		"2321": "⡺",
		"3321": "⡻",
		"2231": "⡼",
		"3231": "⡽",
		"2331": "⡾",
		"3331": "⡿",

		"0002": "⢀",
		"1002": "⢁",
		"0102": "⢂",
		"1102": "⢃",
		"0012": "⢄",
		"1012": "⢅",
		"0112": "⢆",
		"1112": "⢇",
		"2002": "⢈",
		"3002": "⢉",
		"2102": "⢊",
		"3102": "⢋",
		"2012": "⢌",
		"3012": "⢍",
		"2112": "⢎",
		"3112": "⢏",

		"0202": "⢐",
		"1202": "⢑",
		"0302": "⢒",
		"1302": "⢓",
		"0212": "⢔",
		"1212": "⢕",
		"0312": "⢖",
		"1312": "⢗",
		"2202": "⢘",
		"3202": "⢙",
		"2302": "⢚",
		"3302": "⢛",
		"2212": "⢜",
		"3212": "⢝",
		"2312": "⢞",
		"3312": "⢟",

		"0022": "⢠",
		"1022": "⢡",
		"0122": "⢢",
		"1122": "⢣",
		"0032": "⢤",
		"1032": "⢥",
		"0132": "⢦",
		"1132": "⢧",
		"2022": "⢨",
		"3022": "⢩",
		"2122": "⢪",
		"3122": "⢫",
		"2032": "⢬",
		"3032": "⢭",
		"2132": "⢮",
		"3132": "⢯",

		"0222": "⢰",
		"1222": "⢱",
		"0322": "⢲",
		"1322": "⢳",
		"0232": "⢴",
		"1232": "⢵",
		"0332": "⢶",
		"1332": "⢷",
		"2222": "⢸",
		"3222": "⢹",
		"2322": "⢺",
		"3322": "⢻",
		"2232": "⢼",
		"3232": "⢽",
		"2332": "⢾",
		"3332": "⢿",

		"0003": "⣀",
		"1003": "⣁",
		"0103": "⣂",
		"1103": "⣃",
		"0013": "⣄",
		"1013": "⣅",
		"0113": "⣆",
		"1113": "⣇",
		"2003": "⣈",
		"3003": "⣉",
		"2103": "⣊",
		"3103": "⣋",
		"2013": "⣌",
		"3013": "⣍",
		"2113": "⣎",
		"3113": "⣏",

		"0203": "⣐",
		"1203": "⣑",
		"0303": "⣒",
		"1303": "⣓",
		"0213": "⣔",
		"1213": "⣕",
		"0313": "⣖",
		"1313": "⣗",
		"2203": "⣘",
		"3203": "⣙",
		"2303": "⣚",
		"3303": "⣛",
		"2213": "⣜",
		"3213": "⣝",
		"2313": "⣞",
		"3313": "⣟",

		"0023": "⣠",
		"1023": "⣡",
		"0123": "⣢",
		"1123": "⣣",
		"0033": "⣤",
		"1033": "⣥",
		"0133": "⣦",
		"1133": "⣧",
		"2023": "⣨",
		"3023": "⣩",
		"2123": "⣪",
		"3123": "⣫",
		"2033": "⣬",
		"3033": "⣭",
		"2133": "⣮",
		"3133": "⣯",

		"0223": "⣰",
		"1223": "⣱",
		"0323": "⣲",
		"1323": "⣳",
		"0233": "⣴",
		"1233": "⣵",
		"0333": "⣶",
		"1333": "⣷",
		"2223": "⣸",
		"3223": "⣹",
		"2323": "⣺",
		"3323": "⣻",
		"2233": "⣼",
		"3233": "⣽",
		"2333": "⣾",
		"3333": "⣿",
	}

}

//endregion

func (s *img2asciiService) UploadImg(file multipart.File, newFileName string) (string, bool) {
	return s.baseService.SaveImg(file, newFileName)
}

func (s *img2asciiService) Img2ascii(newFileName string, img2ascii *validators.Img2ascii, userData *datamodels.User) (string, bool) {
	var asciiArt models.AsciiArt

	newfilepath := "public/img/" + newFileName
	//call python
	cmd := exec.Command("C:/Users/User/AppData/Local/Programs/Python/Python38-32/python.exe", "py/img2ascii2.py ", newfilepath, strconv.Itoa(img2ascii.Col))
	fmt.Println(cmd.Args)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return err.Error(), false
	}

	//read []byte into AsciiData struct
	asciiData := AsciiData{}
	jsondata := []byte(out)
	err = json.Unmarshal(jsondata, &asciiData)
	if err != nil {
		return err.Error(), false
	}

	block := setblock()
	var asciiContent strings.Builder
	data := strings.Split(asciiData.Data, ",")
	for i, s := range data {
		fmt.Fprint(&asciiContent, block[s])
		if i%asciiData.Cols == 0 && i != 0 {
			fmt.Fprint(&asciiContent, "\n")
		}
	}

	asciiArt = models.AsciiArt{AsciiContent: asciiContent.String(), Public: img2ascii.Public,
		Row: asciiData.Rows, Col: asciiData.Cols, Hot: 0, UserId: userData.UserId}
	s.db.Create(&asciiArt)
	return asciiContent.String(), true
}

func (s *img2asciiService) UpdateHot(updateHot *validators.UpdateHot, userData *datamodels.User) string {
	var asciiArt models.AsciiArt

	if err := s.db.Where("ascii_art_id =?", updateHot.AsciiArtId).
		Where("user_id =?", userData.UserId).First(&asciiArt).Error; err != nil {
		return "查無此資料"
	}

	asciiArt.Hot += 1
	if err := s.db.Save(&asciiArt).Error; err != nil {
		return "資料更新失敗"
	}

	return ""
}

func (s *img2asciiService) GetHotTop() []datamodels.AsciiArt {
	var asciiArts []datamodels.AsciiArt
	if err := s.db.Limit(10).Order("hot desc").Preload("User").Find(&asciiArts).Error; err != nil {
		println(err)
	}
	return asciiArts
}

func (s *img2asciiService) GetUserAsciiById(userData *validators.UserAscii) datamodels.User {
	var user datamodels.User
	if err := s.db.Preload("AsciiArt").Find(&user, userData.UsertId).Error; err != nil {
		println(err)
	}
	return user
}

func (s *img2asciiService) GetMyAscii(userData *datamodels.User) datamodels.User {
	var user datamodels.User
	if err := s.db.Preload("AsciiArt").Find(&user, userData.UserId).Error; err != nil {
		println(err)
	}
	return user
}
