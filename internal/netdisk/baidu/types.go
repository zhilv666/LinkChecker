package baidu

import "time"

type baiduVerifyResp struct {
	Errno     int    `json:"errno"`
	ErrMsg    string `json:"err_msg"`
	RequestID int64  `json:"request_id"`
	Randsk    string `json:"randsk"`
}

type baiduDataResp struct {
	Csrf             string    `json:"csrf"`
	Uk               int       `json:"uk"`
	Username         string    `json:"username"`
	Loginstate       int       `json:"loginstate"`
	SkinName         string    `json:"skinName"`
	Bdstoken         string    `json:"bdstoken"`
	Photo            string    `json:"photo"`
	IsVip            int       `json:"is_vip"`
	IsSvip           int       `json:"is_svip"`
	IsEvip           int       `json:"is_evip"`
	VipIdentity      int       `json:"vip_identity"`
	Now              time.Time `json:"now"`
	Xduss            string    `json:"XDUSS"`
	CurrActivityCode int       `json:"curr_activity_code"`
	ShowVipAd        int       `json:"show_vip_ad"`
	SharePhoto       string    `json:"share_photo"`
	ShareUk          string    `json:"share_uk"`
	Shareid          int64     `json:"shareid"`
	HitOgc           bool      `json:"hit_ogc"`
	ExpiredType      int       `json:"expiredType"`
	Public           int       `json:"public"`
	Ctime            int       `json:"ctime"`
	Description      string    `json:"description"`
	FollowFlag       int       `json:"followFlag"`
	AccessListFlag   bool      `json:"access_list_flag"`
	ElinkInfo        struct {
		IsElink      int  `json:"isElink"`
		EflagDisable bool `json:"eflag_disable"`
	} `json:"Elink_info"`
	Sharetype     int      `json:"sharetype"`
	ViewVisited   int      `json:"view_visited"`
	ViewLimit     int      `json:"view_limit"`
	OwnerVipLevel int      `json:"owner_vip_level"`
	OwnerSvip10ID string   `json:"owner_svip10_id"`
	OwnerVipType  int      `json:"owner_vip_type"`
	Linkusername  string   `json:"linkusername"`
	SharePageType string   `json:"share_page_type"`
	TitleImg      []string `json:"title_img"`
	FileList      []struct {
		AppID          string `json:"app_id"`
		BlackTag       string `json:"black_tag"`
		Category       int    `json:"category"`
		DeleteFsID     string `json:"delete_fs_id"`
		ExtentInt3     string `json:"extent_int3"`
		ExtentInt8     string `json:"extent_int8"`
		ExtentTinyint1 string `json:"extent_tinyint1"`
		ExtentTinyint2 string `json:"extent_tinyint2"`
		ExtentTinyint3 string `json:"extent_tinyint3"`
		ExtentTinyint4 string `json:"extent_tinyint4"`
		FileKey        string `json:"file_key"`
		FileTag        string `json:"file_tag"`
		FsID           int64  `json:"fs_id"`
		IsScene        string `json:"is_scene"`
		Isdelete       string `json:"isdelete"`
		Isdir          int    `json:"isdir"`
		LocalCtime     int    `json:"local_ctime"`
		LocalMtime     int    `json:"local_mtime"`
		Md5            string `json:"md5"`
		OperID         string `json:"oper_id"`
		OwnerID        string `json:"owner_id"`
		OwnerType      string `json:"owner_type"`
		ParentPath     string `json:"parent_path"`
		Path           string `json:"path"`
		PathMd5        string `json:"path_md5"`
		Privacy        string `json:"privacy"`
		RealCategory   string `json:"real_category"`
		RootNs         int    `json:"root_ns"`
		ServerAtime    string `json:"server_atime"`
		ServerCtime    int    `json:"server_ctime"`
		ServerFilename string `json:"server_filename"`
		ServerMtime    int    `json:"server_mtime"`
		Share          string `json:"share"`
		Size           int64  `json:"size"`
		Source         string `json:"source"`
		Status         string `json:"status"`
		TkbindID       string `json:"tkbind_id"`
		Videotag       string `json:"videotag"`
		Wpfile         string `json:"wpfile"`
	} `json:"file_list"`
	Errortype int `json:"errortype"`
	Errno     int `json:"errno"`
	UfcTime   int `json:"ufcTime"`
	Error     int `json:"error"`
	Data      struct {
		ExpName string `json:"expName"`
	} `json:"data"`
	Self      int `json:"self"`
	ElinkSelf int `json:"elink_self"`
}
