package quark

type quarkTokenResp struct {
	Status    int    `json:"status"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
	Data      struct {
		Subscribed bool   `json:"subscribed"`
		Stoken     string `json:"stoken"`
		ShareType  int    `json:"share_type"`
		Author     struct {
			MemberType string `json:"member_type"`
			AvatarURL  string `json:"avatar_url"`
			NickName   string `json:"nick_name"`
		} `json:"author"`
		URLType     int    `json:"url_type"`
		ExpiredType int    `json:"expired_type"`
		ExpiredAt   int64  `json:"expired_at"`
		Title       string `json:"title"`
		FileNum     int    `json:"file_num"`
	} `json:"data"`
	Metadata struct {
		TGroup string `json:"_t_group"`
		GGroup string `json:"_g_group"`
	} `json:"metadata"`
}

type quarkDetailResp struct {
	Status    int    `json:"status"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
	Data      struct {
		IsOwner int `json:"is_owner"`
		Share   struct {
			Title       string `json:"title"`
			ShareType   int    `json:"share_type"`
			ShareID     string `json:"share_id"`
			PwdID       string `json:"pwd_id"`
			ShareURL    string `json:"share_url"`
			URLType     int    `json:"url_type"`
			FirstFid    string `json:"first_fid"`
			ExpiredType int    `json:"expired_type"`
			FileNum     int    `json:"file_num"`
			CreatedAt   int64  `json:"created_at"`
			UpdatedAt   int64  `json:"updated_at"`
			ExpiredAt   int64  `json:"expired_at"`
			ExpiredLeft int64  `json:"expired_left"`
			AuditStatus int    `json:"audit_status"`
			Status      int    `json:"status"`
			ClickPv     int    `json:"click_pv"`
			SavePv      int    `json:"save_pv"`
			DownloadPv  int    `json:"download_pv"`
			FirstFile   struct {
				Fid                     string        `json:"fid"`
				Category                int           `json:"category"`
				FileType                int           `json:"file_type"`
				Size                    int           `json:"size"`
				FormatType              string        `json:"format_type"`
				NameSpace               int           `json:"name_space"`
				SeriesDir               bool          `json:"series_dir"`
				AlbumDir                bool          `json:"album_dir"`
				MoreThanOneLayer        bool          `json:"more_than_one_layer"`
				UploadCameraRootDir     bool          `json:"upload_camera_root_dir"`
				Fps                     float64       `json:"fps"`
				Like                    int           `json:"like"`
				RiskType                int           `json:"risk_type"`
				TagList                 []interface{} `json:"tag_list"`
				FileNameHlStart         int           `json:"file_name_hl_start"`
				FileNameHlEnd           int           `json:"file_name_hl_end"`
				Duration                int           `json:"duration"`
				ScrapeStatus            int           `json:"scrape_status"`
				CurVersionOrDefault     int           `json:"cur_version_or_default"`
				OwnerDriveTypeOrDefault int           `json:"owner_drive_type_or_default"`
				SaveAsSource            bool          `json:"save_as_source"`
				BackupSource            bool          `json:"backup_source"`
				OfflineSource           bool          `json:"offline_source"`
				EnsureValidSaveAsLayer  int           `json:"ensure_valid_save_as_layer"`
				Ban                     bool          `json:"ban"`
				Dir                     bool          `json:"dir"`
				File                    bool          `json:"file"`
				Extra                   struct {
				} `json:"_extra"`
			} `json:"first_file"`
			PathInfo                 string `json:"path_info"`
			PartialViolation         bool   `json:"partial_violation"`
			Size                     int64  `json:"size"`
			FirstLayerFileCategories []int  `json:"first_layer_file_categories"`
			PicTotal                 int    `json:"pic_total"`
			VideoTotal               int    `json:"video_total"`
			IsAllImageFile           int    `json:"is_all_image_file"`
			IsOwner                  int    `json:"is_owner"`
			MemberType               string `json:"member_type"`
			FileOnlyNum              int    `json:"file_only_num"`
			AllFileNum               int    `json:"all_file_num"`
			DownloadPvlimited        bool   `json:"download_pvlimited"`
		} `json:"share"`
		List []struct {
			Fid                 string        `json:"fid"`
			FileName            string        `json:"file_name"`
			PdirFid             string        `json:"pdir_fid"`
			Category            int           `json:"category"`
			FileType            int           `json:"file_type"`
			Size                int           `json:"size"`
			FormatType          string        `json:"format_type"`
			Status              int           `json:"status"`
			Tags                string        `json:"tags"`
			LCreatedAt          int64         `json:"l_created_at"`
			LUpdatedAt          int64         `json:"l_updated_at"`
			Extra               string        `json:"extra"`
			Source              string        `json:"source"`
			FileSource          string        `json:"file_source"`
			NameSpace           int           `json:"name_space"`
			LShotAt             int64         `json:"l_shot_at"`
			SourceDisplay       string        `json:"source_display"`
			IncludeItems        int           `json:"include_items"`
			SeriesDir           bool          `json:"series_dir"`
			AlbumDir            bool          `json:"album_dir"`
			MoreThanOneLayer    bool          `json:"more_than_one_layer"`
			UploadCameraRootDir bool          `json:"upload_camera_root_dir"`
			Fps                 float64       `json:"fps"`
			Like                int           `json:"like"`
			OperatedAt          int64         `json:"operated_at"`
			RiskType            int           `json:"risk_type"`
			TagList             []interface{} `json:"tag_list"`
			BackupSign          int           `json:"backup_sign"`
			FileNameHlStart     int           `json:"file_name_hl_start"`
			FileNameHlEnd       int           `json:"file_name_hl_end"`
			FileStruct          struct {
				FirSource      string `json:"fir_source"`
				SecSource      string `json:"sec_source"`
				ThiSource      string `json:"thi_source"`
				PlatformSource string `json:"platform_source"`
			} `json:"file_struct"`
			Duration   int `json:"duration"`
			EventExtra struct {
				RecentCreatedAt int64 `json:"recent_created_at"`
			} `json:"event_extra"`
			ScrapeStatus            int    `json:"scrape_status"`
			UpdateViewAt            int64  `json:"update_view_at"`
			LastUpdateAt            int64  `json:"last_update_at"`
			ShareFidToken           string `json:"share_fid_token"`
			CurVersionOrDefault     int    `json:"cur_version_or_default"`
			RawNameSpace            int    `json:"raw_name_space"`
			OwnerDriveTypeOrDefault int    `json:"owner_drive_type_or_default"`
			SaveAsSource            bool   `json:"save_as_source"`
			BackupSource            bool   `json:"backup_source"`
			OfflineSource           bool   `json:"offline_source"`
			EnsureValidSaveAsLayer  int    `json:"ensure_valid_save_as_layer"`
			Ban                     bool   `json:"ban"`
			Dir                     bool   `json:"dir"`
			File                    bool   `json:"file"`
			CreatedAt               int64  `json:"created_at"`
			UpdatedAt               int64  `json:"updated_at"`
			Extra0                  struct {
			} `json:"_extra"`
		} `json:"list"`
	} `json:"data"`
	Metadata struct {
		Size          int    `json:"_size"`
		Page          int    `json:"_page"`
		VideoTotal    int    `json:"video_total"`
		Count         int    `json:"_count"`
		Total         int    `json:"_total"`
		CheckFidToken int    `json:"check_fid_token"`
		GGroup        string `json:"_g_group"`
		PicTotal      int    `json:"pic_total"`
		TGroup        string `json:"_t_group"`
	} `json:"metadata"`
}
