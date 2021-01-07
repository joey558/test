package model

type GameType struct {
	Id    int `orm:"pk"`
	Code  string
	Title string
	Icon  string
}

type PlatList struct {
	Id    int `orm:"pk"`
	Code  string
	Title string
	Icon  string
}

type GamePlat struct {
	Id        int `orm:"pk"`
	Plat_code string
	Type_code string
	Title     string
	Sort      int
}

type BannerList struct {
	Id       int `orm:"pk"`
	H5_img   string
	Pc_img   string
	Title    string
	Jump_url string
	Sort     int
	Status   int
}

type ActorList struct {
	Id          int `orm:"pk"`
	Name        string
	Nation      string
	Birth       string
	Sex         int
	Cn_code     string
	Initials    string
	Avatar      string
	Is_hot      int
	Height      string
	Cup         string
	Description string
	Total       int
	Bust        string
	Waist       string
	Hips        string
	Status      int
	Real_name   string
	Nick_name   string
	Updatetime  string
}

type LoadingList struct {
	Id           int `orm:"pk"`
	Img_url      string
	Loading_time int
	Jump_url     string
	Status       int
	Sort         int
}

type TagType struct {
	Id      int `orm:"pk"`
	Img_url string
	Title   string
	Is_top  int
	Status  int
	Sort    int
}

type TagList struct {
	Id    int `orm:"pk"`
	Title string
}

type PopList struct {
	Id         int `orm:"pk"`
	Title      string
	Content    string
	Butt_title string
	Jump_url   string
	Status     int
	Sort       int
}

type UserList struct {
	Id         string `orm:"pk"`
	Account    string
	Pwd        string
	Status     int
	Reg_time   string
	Reg_ip     string
	Reg_app    string
	Login_time string
	Login_ip   string
	Is_agent   int
	Is_test    int
	Agent_name string
	Agent_path string
	Link_code  string
	Reg_code   string
	Note       string
	Is_line    int
	Session_id string
}

type UserInfo struct {
	Id          string `orm:"pk"`
	Account     string
	Kok_amount  float64
	Amount      float64
	Total_in    float64
	Total_out   float64
	Total_bonus float64
	In_num      int
	Out_num     int
	Last_in     string
	Last_out    string
	Phone       string
	Is_test     int
	Nick_name   string
	Avatar      string
	Age         int
	Sex         int
	Per_sign    string
	Score       int
	Level       int
	Exp         int
	Day_num     int
}

type LevelList struct {
	Id    int `orm:"pk"`
	Title string
}

type VideoList struct {
	Id        string `orm:"pk"`
	Title     string
	Content   string
	View_num  int
	Uploader  string
	Like_num  int
	Star_num  int
	Status    int
	Pc_img    string
	H5_img    string
	Url       string
	Score     int
	Is_hot    int
	Is_recomm int
}

type VideoTag struct {
	Id       string
	Video_id string
	Tag_id   int
}

type VideoType struct {
	Id       string
	Video_id string
	Type_id  int
}

type VideoLike struct {
	Id          string
	Video_id    string
	Account     string
	Create_time string
}

type VideoStar struct {
	Id          string
	Video_id    string
	Account     string
	Create_time string
}

type VideoHistory struct {
	Id          string
	Video_id    string
	Account     string
	Create_time string
}

type VideoAct struct {
	Id       string
	Video_id string
	Act_id   int
}

type VideoComm struct {
	Id          string
	Video_id    string
	Account     string
	To_account  string
	Content     string
	P_id        string
	Nick_name   string
	To_name     string
	Like_num    int
	Reply_num   int
	Create_time string
}

type OrderList struct {
	Id           string
	Amount       float64
	Account      string
	Status       int
	Order_number string
	Agent_name   string
	Order_type   int
	Is_test      int
	Agent_path   string
	Pay_id       int
	Create_time  string
	Is_auto      int
	Act_id       int
	Note         string
	Score        int
	Creator      string
	Operator     string
	Pay_time     string
}

type AmountList struct {
	Id           string
	Amount       float64
	Account      string
	Order_number string
	Agent_name   string
	Order_type   int
	Is_test      int
	Agent_path   string
	Create_time  string
	Note         string
	Score        int
}

type TransferList struct {
	Id            string
	Amount        float64
	Account       string
	Status        int
	Order_number  string
	Order_type    int
	Is_test       int
	Create_time   string
	Is_auto       int
	Note          string
	Creator       string
	Pay_time      string
	From_plat     string
	To_plat       string
	Before_amount float64
	After_amount  float64
}

type PlatAccount struct {
	Id          string
	Plat_code   string
	Account     string
	User        string
	Pwd         string
	Create_time string
}

type ActiveType struct {
	Id    int
	Title string
	Sort  int
}

type ActiveBonus struct {
	Id    int
	Title string
}

type ActiveList struct {
	Id         string
	Title      string
	Sub_title  string
	Pc_img     string
	H5_img     string
	Img        string
	Status     int
	Start_time string
	End_time   string
	S_date     string
	E_date     string
	Content    string
	Jump_url   string
	Note       string
	Type_id    int
	Act_id     int
	Sort       int
}

type SysBank struct {
	Id    int
	Code  string
	Title string
	Icon  string
}

type SysMsg struct {
	Id          string
	Sub_title   string
	Title       string
	Content     string
	Create_time string
	To_user     string //接受的用户,所有用户=all
	Msg_type    int    //消息通知类型:0=系统消息,1=评论回复,2=评论点赞
}

type MsgView struct {
	Id          string
	Msg_id      string
	Account     string
	Create_time string
	Status      int
}

type AppSetup struct {
	Id           string
	App_type     string
	Phone_info   string
	Phone_os     string
	Ip           string
	Uid          string
	Country_code string
	Country_name string
	Province     string
	Last_req     string
	Phone_number string
	Create_time  string
}

type AppLog struct {
	Id           string
	App_type     string
	Phone_info   string
	Phone_os     string
	Ip           string
	Uid          string
	Country_code string
	Country_name string
	Province     string
	Last_req     string
	Phone_number string
	Start_time   string
	Close_time   string
	Account      string
	Day_date     string
	App_minu     int
}

type UserLog struct {
	Id          string
	Ip          string
	Title       string
	Content     string
	Url         string
	Host        string
	Param       string
	Res         string
	Note        string
	Create_time string
	Account     string
	Status      int
}

type PayConfig struct {
	Id          int
	Title       string
	Class_code  string
	Bank_code   string
	Min_amount  float64
	Max_amount  float64
	Status      int
	Is_fix      int
	Fix_amount  string
	Conf_id     int
	Note        string
	Create_time string
}

type ConfList struct {
	Id          int
	Title       string
	Status      int
	Conf        string
	Note        string
	Create_time string
}

type CaptchaList struct {
	Id           string
	Phone_number string
	Code         string
	Account      string
	Create_time  string
}

type TaskList struct {
	Id           string
	Title        string
	Content      string
	Score        int
	Icon         string
	Act_id       int
	Status       int
	Tast_type    int
	Create_time  string
	App_url_code string
}

type TaskLog struct {
	Id          string
	Task_id     string
	Note        string
	Score       int
	Account     string
	Create_time string
}

type ContactList struct {
	Id          int
	Code        string
	Content     string
	Note        string
	Status      int
	Icon        string
	Operator    string
	Create_time string
}

type DomainList struct {
	Id       int
	Domain   string
	App_type string
	Status   int
}

type CommLike struct {
	Id          string
	Comm_id     string
	Account     string
	Create_time string
}

type WebConf struct {
	Id    int
	Code  string
	Title string
	Val   string
}

type AppVersion struct {
	Id      int
	version string
	note    string
	Val     string
}

type AppUrl struct {
	Id       int
	Url      string
	App_type string
	Code     string
	Title    string
}

type PayClass struct {
	Id    int `orm:"pk"`
	Title string
	Code  string
	Icon  string
	Note  string
	Sort  int
}

type ThemeVideo struct {
	Id          string
	Title       string
	Content     string
	View_num    int
	Uploader    string
	Like_num    int
	Status      int
	H5_img      string
	Url         string
	Is_hot      int
	Is_recomm   int
	Theme_id    int
	Sort        int
	Create_time string
}

type ThemeList struct {
	Id          int
	Title       string
	Content     string
	Icon        string
	Create_time string
}

type ThemeLike struct {
	Id          string
	Comm_id     string
	Account     string
	Create_time string
}

type BlogListLike struct {
	Id          string
	BlogListId  string
	Account     string
	Create_time string
}

type BlogCommLike struct {
	Id          string
	Comm_id     string
	Account     string
	Create_time string
}

type BlogList struct {
	Id          string
	Content     string
	Img_url     string
	View_num    int
	Account     string
	Nick_name   string
	Like_num    int
	Avatar      string
	Is_top      int
	Is_good     int
	Comm_num    int
	Create_time string
	Status      int
}

type BlogComm struct {
	Id          string
	Blog_id     string
	Account     string
	To_account  string
	Content     string
	P_id        string
	Nick_name   string
	To_name     string
	Like_num    int
	Reply_num   int
	Create_time string
	Status      int
	Is_up       int
}
