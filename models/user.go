package models

import (
	"errors"
	"log"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	. "telepathy/Utils"
)

//用户表
type User struct {
	Id            int64
	Phone         uint32
	Username      string    `orm:"unique;size(32)" form:"Username"  valid:"Required;MaxSize(20);MinSize(6)"`
	Password      string    `orm:"size(32)" form:"Password" valid:"Required;MaxSize(20);MinSize(6)"`
	Repassword    string    `orm:"-" form:"Repassword" valid:"Required"`
	Nickname      string    `orm:"unique;size(32)" form:"Nickname" valid:"Required;MaxSize(20);MinSize(2)"`
	Email         string    `orm:"size(32)" form:"Email" valid:"Email"`
	Remark        string    `orm:"null;size(200)" form:"Remark" valid:"MaxSize(200)"`
	Status        int       `orm:"default(2)" form:"Status" valid:"Range(1,2)"`
	Lastlogintime time.Time `orm:"null;type(datetime)" form:"-"`
	Createtime    time.Time `orm:"type(datetime);auto_now_add" `
}

func (u *User) TableName() string {
	return beego.AppConfig.String("telepathy_user_table")
}

func (u *User) Valid(v *validation.Validation) {
	if u.Password != u.Repassword {
		v.SetError("Repassword", "两次输入的密码不一样")
	}
}

//验证用户信息
func checkUser(u *User) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&u)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

func init() {
	orm.RegisterModel(new(User))
}

//----------------------------------------------------------

//get user list
func Getuserlist(page int64, page_size int64, sort string) (users []orm.Params, count int64) {
	o := orm.NewOrm()
	user := new(User)
	qs := o.QueryTable(user)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	qs.Limit(page_size, offset).OrderBy(sort).Values(&users)
	count, _ = qs.Count()
	return users, count
}

func DelUserById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&User{Id: Id})
	return status, err
}

func GetUserByPhone(phone uint32) (user User) {
	user = User{Phone: phone}
	o := orm.NewOrm()
	o.Read(&user, "Phone")
	return user
}

//---------------------------------------------
//注册相关
func CanRegister(u *User) (err error) {
	user := GetUserByPhone(u.Phone)
	if user.Id == 0 {
		return errors.New("用户不存在")
	}

	return nil
}

func Register(u *User) error {
	if err := CanRegister(u); err != nil {
		return err
	}

	//TODO:发送验证码至手机

	return nil
}

func CanRegisterEnd(u *User) error {
	//检查用户是否已经存在
	user := GetUserByPhone(u.Phone)
	if user.Id > 0 {
		return errors.New("用户已存在")
	}

	//验证两次密码是否一致
	valid := validation.Validation{}
	b, _ := valid.Valid(&u)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

func RegisterEnd(u *User) (int64, error) {
	if err := CanRegisterEnd(u); err != nil {
		return 0, err
	}

	o := orm.NewOrm()
	user := new(User)
	user.Phone = u.Phone
	user.Username = u.Username
	user.Password = Strtomd5(u.Password)
	user.Nickname = u.Nickname
	user.Email = u.Email
	user.Remark = u.Remark
	user.Status = u.Status

	id, err := o.Insert(user)
	return id, err
}

//---------------------------------------------
//登录login相关

func CanLogin(u *User) error {
	//帐号验证
	user := GetUserByPhone(u.Phone)
	if user.Id == 0 {
		return errors.New("用户不存在")
	}

	//密码验证
	if user.Password != Pwdhash(u.Password) {
		return errors.New("密码错误")
	}

	return nil
}

func Login(u *User) error {
	if err := CanLogin(u); err != nil {
		return err
	}

	return nil
}

//---------------------------------------------
//更新用户profile
func CanUpdateProfile(u *User) error {
	//检查用户是否存在
	user := GetUserByPhone(u.Phone)
	if user.Id == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func UpdateProfile(u *User) (int64, error) {
	if err := CanUpdateProfile(u); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	user := make(orm.Params)
	if len(u.Username) > 0 {
		user["Username"] = u.Username
	}
	if len(u.Nickname) > 0 {
		user["Nickname"] = u.Nickname
	}
	if len(u.Email) > 0 {
		user["Email"] = u.Email
	}
	if len(u.Remark) > 0 {
		user["Remark"] = u.Remark
	}
	if len(u.Password) > 0 {
		user["Password"] = Strtomd5(u.Password)
	}
	if u.Status != 0 {
		user["Status"] = u.Status
	}
	if len(user) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table User
	num, err := o.QueryTable(table).Filter("Phone", u.Phone).Update(user)
	return num, err
}
