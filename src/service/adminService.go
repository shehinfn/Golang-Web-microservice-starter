package service

import (
	"errors"
	"fmt"
	"log"
	"sondr-backend/src/models"
	"sondr-backend/src/repository"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	//"gorm.io/gorm"
)

/*******************************CREATING SUB-ADMINS**********************************/
func (c *TestAPIAdmin) CreateSubadmin(insert *models.Admins) (*models.Admins, error) {
	password, _ := bcrypt.GenerateFromPassword([]byte(insert.Password), 4)
	insert.Password = string(password)
	if insert.Role == "" {
		insert.Role = "SubAdmin"
	}
	if err := repository.Repo.Insert(insert); err != nil {
		return nil, err
	}

	return insert, nil
}

/******************************GENARATE PASSWORD MANUALLY***************************/
func (c *TestAPIAdmin) GeneratePassword() (string, error) {
	genaratePassword, err := password.Generate(12, 1, 1, false, false)
	if err != nil {
		log.Fatal(err)
	}
	return genaratePassword, nil
}

/*********************************LOGIN*********************************************/
func (r *TestAPIAdmin) Login(req *models.Request) (string, error) {
	obj := models.Admins{}
	whereQuery := "email = ?"
	if err := repository.Repo.Find(&obj, "admins", "", whereQuery, req.Email); err != nil {
		fmt.Println("Login failed")
		return "", errors.New("Invalid creditional")
	}
	err := bcrypt.CompareHashAndPassword([]byte(obj.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("Invalid Creditional")
	}
	return GenrateToken(&obj)

}

func GenrateToken(admin *models.Admins) (string, error) {
	var mySigningKey = []byte(viper.GetString("secret.Key"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = admin.Email
	claims["role"] = admin.Role
	claims["password"] = admin.Password

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

/*********************************LISTING SUBADMINS*********************************************/

func (r *TestAPIAdmin) ListSubAdmins(pageNo int, pageSize int, search string) (*models.AdminResponse, error) {
	var admin []*models.ListAdmin
	var resp models.AdminResponse
	var whereQuery string
	var searchFilter interface{}
	if pageSize == 0 {
		pageSize = 10
	}

	selectQuery := "admins.id, admins.name, admins.email, admins.password, admins.created_at"

	if search != "" {
		whereQuery = "admins.id like ? or admins.name like ?"
		searchFilter = "%" + search + "%"
	}
	var count int
	count, err := repository.Repo.ListAllWithPagination(&admin, selectQuery, "admins", "", pageNo, pageSize, whereQuery, searchFilter)
	if err != nil {
		return nil, err
	}
	resp.ListAdmin = admin
	resp.Count = count
	return &resp, nil
}

/*func (r *TestAPIAdmin) ListSubAdmins(dt models.AdminsPagination) ([]models.Admins, int, error) {
	var Admins []models.Admins
	var countVal string
	Limit := dt.PageSize
	Offset := (dt.PageID - 1) * dt.PageSize
	sqlDB := database.DB

	sqlDB.Offset(Offset).Limit(Limit).Debug().Find(&Admins)
	result := sqlDB.Raw("SELECT count(*) FROM Admins").Scan(&countVal)
	if result == nil {
		return Admins, 0, nil
	}
	count, _ := strconv.Atoi(countVal)
	return Admins, count, nil

}*/

/*********************************READ SUBADMIN DETAILS*********************************************/
func (c *TestAPIAdmin) ReadSubAdmin(id uint) (*models.ListAdmin, error) {
	//obj := models.Admins{}
	var admin models.ListAdmin
	whereQuery := "id = ?"

	selectQuery := "admins.name,admins.email,admins.id"
	//var password_tes = bcrypt.
	if err := repository.Repo.Find(&admin, "admins", selectQuery, whereQuery, id); err != nil {
		return nil, err
	}
	return &admin, nil
}

/*********************************UPDATE SUBADMIN DETAILS*********************************************/
/*func (c *TestAPIAdmin) UpdateSubAdmin(update *models.Admins) (string, error) {
	//var admin models.Admins
	// admin.Name = update.Name
	// admin.Email = update.Email
	// admin.Password = update.Password
	if err := repository.Repo.Update(&models.Admins{}, update.ID, update); err != nil {
		return "Unable to update", err
	}
	return "Update Sucessfull", nil

}*/

/*********************************UPDATE SUBADMIN DETAILS*********************************************/
func (c *TestAPIAdmin) UpdateSubAdmin(update *models.Admin) (string, error) {
	val := repository.Repo.UpdateSubAdmin(&models.Admins{}, update.ID, update)
	if val.RowsAffected > 0 {
		return "Update successfull", nil
	}
	return "Record already found or unable to update ", nil
	// val := database.DB.Debug().Model(models.Admins{}).Where("id IN (?)", update.ID).Updates(models.Admins{Name: update.Name, Email: update.Email})
	// if val.RowsAffected > 0 {
	// 	fmt.Println("rows affected", val.RowsAffected)
	// 	return val, nil
	// } else {
	// 	return nil, nil
	// }
}

/*********************************Verify Password*********************************************/

// func (c *TestAPIAdmin) VerifyPassword(req *models.Admins) (string, error) {
// 	obj := models.Admins{}
// 	whereQuery := "id = ?"
// 	if err := repository.Repo.Find(&obj, "", whereQuery, req.ID); err != nil {
// 		return "Login Failed. Please enter correct Email and Password", err
// 	}

// 	if req.Password == obj.Password {
// 		return "Verification Successfull", nil
// 	} else if req.Password == "" {
// 		return "Please enter the  current password to verify", nil
// 	} else {
// 		return "Verification failed. Incorrect password", nil
// 	}
// }

// /*********************************CHANGE PASSWORD*********************************************/

// func (c *TestAPIAdmin) ChangePassword(req *models.Admin) (string, error) {
// 	obj := models.Admins{}
// 	whereQuery := "id = ?"
// 	if err := repository.Repo.Find(&obj, "", whereQuery, req.ID); err != nil {
// 		return "Unable to change password", err
// 	}
// 	if req.Password == obj.Password {
// 		return "Current password and new password should not be same", nil
// 	}
// 	if err := repository.Repo.UpdateSubAdmin(&models.Admins{}, req.ID, obj); err != nil {
// 		return "Unable to change password", nil
// 	}
// 	return "Password changed successfully", nil
// }
