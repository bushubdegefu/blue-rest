package temps

import (
	"os"
	"text/template"
)

func AuthUtilsFrame(appName string) {
	ProjectSettings.CurrentAppName = appName
	ProjectSettings.BackTick = "`"
	// ############################################################
	utils_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(utilsTemplate)
	if err != nil {
		panic(err)
	}

	// create database folder if does not exist
	err = os.MkdirAll("utils", os.ModePerm)
	if err != nil {
		panic(err)
	}

	utils_file, err := os.Create("utils/utils.go")
	if err != nil {
		panic(err)
	}
	defer utils_file.Close()

	err = utils_tmpl.Execute(utils_file, ProjectSettings)
	if err != nil {
		panic(err)
	}

	// ############################################################
	jwt_utils_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(jwtUtilsTemplate)
	if err != nil {
		panic(err)
	}

	jwt_utils_file, err := os.Create("utils/jwt_utils.go")
	if err != nil {
		panic(err)
	}
	defer jwt_utils_file.Close()

	err = jwt_utils_tmpl.Execute(jwt_utils_file, ProjectSettings)
	if err != nil {
		panic(err)
	}
}

func AuthLoginFrame(appName, frame string) {
	ProjectSettings.CurrentAppName = appName
	ProjectSettings.BackTick = "`"
	var utils_tmpl *template.Template
	var err error
	// ############################################################
	if frame == "fiber" {
		utils_tmpl, err = template.New("RenderData").Funcs(FuncMap).Parse(loginTemplateFiber)
		if err != nil {
			panic(err)
		}
	} else {
		utils_tmpl, err = template.New("RenderData").Funcs(FuncMap).Parse(loginTemplate)
		if err != nil {
			panic(err)
		}

	}

	utils_file, err := os.Create("controllers/login.go")
	if err != nil {
		panic(err)
	}
	defer utils_file.Close()

	err = utils_tmpl.Execute(utils_file, ProjectSettings)
	if err != nil {
		panic(err)
	}

}

var jwtUtilsTemplate = `
package utils

import (
	"math/rand"
	"strconv"

	"{{.ProjectName}}/{{.CurrentAppName}}/models"
	"{{.ProjectName}}/configs"
	"{{.ProjectName}}/database"
)

const (
	charsetLen = 62
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateRandomString(length int) (string, error) {
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		index := rand.Intn(charsetLen)
		result[i] = charset[index]
	}

	return string(result), nil
}

func JWTSaltUpdate() {
	dbcon, _ := database.ReturnSession("{{ .CurrentAppName | replaceString}}")

	//  getting salt length from configuration file
	salt_length, _ := strconv.Atoi(configs.AppConfig.Get("JWT_SALT_LENGTH"))

	//  generating for salt A
	salt_a, _ := GenerateRandomString(salt_length)

	//  defining jwt _object to work with
	var jwt_object models.JWTSalt

	// Fethching the JWT object if it exists
	dbcon.Model(&models.JWTSalt{}).Where("id = ?", 1).First(&jwt_object)

	// If it exists updating or creating if it does not exist
	if jwt_object.ID == 1 {
		// salt_b =
		jwt_object.SaltB = jwt_object.SaltA
		jwt_object.SaltA = salt_a

		tx := dbcon.Begin()
		if err := dbcon.Updates(&jwt_object).Error; err != nil {

			tx.Rollback()
		}
		tx.Commit()
	} else {
		// generating salt b and creating new
		salt_b, _ := GenerateRandomString(salt_length)
		jwt_object.SaltA = salt_a
		jwt_object.SaltB = salt_b

		tx := dbcon.Begin()
		// add  data using transaction if values are valid
		if err := tx.Create(&jwt_object).Error; err != nil {

			tx.Rollback()

		}
		tx.Commit()
	}

}

func GetJWTSalt() (salt_a string, salt_b string) {
	dbcon, _ := database.ReturnSession("{{ .CurrentAppName | replaceString}}")
	var jwt_object models.JWTSalt

	// Fethching the JWT object if it exists
	dbcon.Model(&models.JWTSalt{}).Where("id = ?", 1).First(&jwt_object)

	salt_a = jwt_object.SaltA
	salt_b = jwt_object.SaltB
	if salt_a == "" {
		JWTSaltUpdate()
	}
	return salt_a, salt_b

}

`
var utilsTemplate = `
package utils

import (
	"fmt"
	"reflect"
	"time"

	"{{.ProjectName}}/{{.CurrentAppName}}/models"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	jwt.RegisteredClaims
	Email       string   {{.BackTick}}json:"email"{{.BackTick}}
	Permissions []string {{.BackTick}}json:"permissions"{{.BackTick}}
	Groups      []string {{.BackTick}}json:"groups"{{.BackTick}}
	UUID        string   {{.BackTick}}json:"uuid"{{.BackTick}}
	UserID      int      {{.BackTick}}json:"user_id"{{.BackTick}}
}

func PasswordsMatch(hashedPassword, currPassword string) bool {

	var currPasswordHash = models.HashFunc(currPassword)
	return hashedPassword == currPasswordHash
}

// source of this token encode decode functions
// https://github.com/gurleensethi/go-jwt-tutorial/blob/main/main.go
func CreateJWTToken(email string, user_id int, permissions, groups []string, duration int) (string, error) {
	my_claim := UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		Email:            email,
		Permissions:      permissions,
		Groups:           groups,
		UserID:           user_id,
	}

	salt_a, _ := GetJWTSalt()
	exp := time.Now().UTC().Add(time.Duration(duration) * time.Minute)
	my_claim.ExpiresAt = jwt.NewNumericDate(exp)
	my_claim.Issuer = "Blue Admin"
	my_claim.Subject = "UI Authentication Token"
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, my_claim)
	signedString, err := token.SignedString([]byte(salt_a))
	if err != nil {
		return "", fmt.Errorf("error creating signed string: %v", err)
	}

	return signedString, nil
}

func ParseJWTToken(jwtToken string) (UserClaim, error) {
	salt_a, salt_b := GetJWTSalt()
	response_a := UserClaim{}
	response_b := UserClaim{}

	token_a, aerr := jwt.ParseWithClaims(jwtToken, &response_a, func(token *jwt.Token) (interface{}, error) {
		return []byte(salt_a), nil
	})
	token_b, berr := jwt.ParseWithClaims(jwtToken, &response_b, func(token *jwt.Token) (interface{}, error) {
		return []byte(salt_b), nil
	})

	if aerr != nil && berr != nil {
		return UserClaim{}, aerr
	}

	// check token validity, for example token might have been expired
	if !token_a.Valid {
		if !token_b.Valid {
			return UserClaim{}, fmt.Errorf("invalid token with second salt")
		}
		return response_b, nil
	}
	return response_a, nil

}

// Return Unique values in list
func UniqueSlice(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Return Unique values in list
func CheckValueExistsInSlice(slice []string, role_test string) bool {
	for _, role := range slice {
		if role == role_test || role == "superuser" {
			return true
		}
	}
	return false
}

// Struct to Map conversion function
func StructToMap(input interface{}) (map[string]interface{}, error) {
	// Create an empty map
	result := make(map[string]interface{})

	// Get the reflect value of the struct
	val := reflect.ValueOf(input)

	// Ensure that the input is a pointer to a struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Check if the input is a struct
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	// Loop through the struct fields
	for i := 0; i < val.NumField(); i++ {
		// Get the field and its name
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name

		// Insert the field name and value into the map
		result[fieldName] = field.Interface()
	}

	return result, nil
}

`

var loginTemplate = `
package controllers

import (
	"net/http"
	"time"
	"fmt"

	"{{.ProjectName}}/{{ .CurrentAppName}}/models"
	"{{.ProjectName}}/{{ .CurrentAppName}}/utils"

	"{{.ProjectName}}/common"
	"{{.ProjectName}}/observe"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Login Request for Endpoint
type LoginPost struct {
	GrantType string {{.BackTick}}json:"grant_type" validate:"required" example:"authorization_code,refresh_token,token_decode"{{.BackTick}}
	Email     string {{.BackTick}}json:"email" validate:"email,min=6,max=32"{{.BackTick}}
	Password  string {{.BackTick}}json:"password"{{.BackTick}}
	Token     string {{.BackTick}}json:"token"{{.BackTick}}
}

// Access token Response
type TokenResponse struct {
	AccessToken  string {{.BackTick}}json:"access_token"{{.BackTick}}
	RefreshToken string {{.BackTick}}json:"refresh_token"{{.BackTick}}
	TokenType    string {{.BackTick}}json:"token_type"{{.BackTick}}
}

// Login is a function to login by EMAIL and ID
// @Summary Auth
// @Description Login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body LoginPost true "Login"
// @Success 200 {object} common.ResponseHTTP{data=TokenResponse{}}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /{{.CurrentAppName | replaceString}}/login [post]
func Login(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  Geting dbsession
	db := contx.Get("db").(*gorm.DB)

	// validator initialization
	validate := validator.New()

	//validating post data
	login_request := new(LoginPost)

	//first parse request data
	if err := contx.Bind(&login_request); err != nil {

		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(login_request); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	switch login_request.GrantType {
	case "authorization_code":
		var user models.User
		var permissions []models.Permission
		var groups []models.Group

		// Getting User
		res := db.WithContext(tracer.Tracer).Model(&models.User{}).Where("email = ? AND is_active = ?", login_request.Email, true).First(&user)
		if res.Error != nil {

			return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
				Success: false,
				Message: res.Error.Error(),
				Data:    nil,
			})
		}
		// Getting User Permissions
		permission_join_string := "INNER JOIN user_permissions up ON permissions.id = up.permission_id"
		res_per := db.WithContext(tracer.Tracer).Model(&models.Permission{}).Joins(permission_join_string).Where("up.user_id = ?", user.ID).Scan(&permissions)
		if res_per.Error != nil {
			return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
				Success: false,
				Message: res.Error.Error(),
				Data:    nil,
			})
		}

		// Getting User Groups
		groups_join_string := "INNER JOIN user_groups ug ON groups.id = ug.group_id"
		res_gr := db.WithContext(tracer.Tracer).Model(&models.Group{}).Joins(groups_join_string).Where("ug.user_id = ?", user.ID).Scan(&groups)
		if res_gr.Error != nil {
			return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
				Success: false,
				Message: res.Error.Error(),
				Data:    nil,
			})
		}

		// Checking Password match
		if utils.PasswordsMatch(user.Password, login_request.Password) {
			// preparing permissions for JWT token hash
			user_permissions := make([]string, 0, 20)
			for _, value := range permissions {
				user_permissions = append(user_permissions, string(value.Codename))
			}
			// preparing groups for JWT token hash
			user_groups := make([]string, 0, 20)
			for _, value := range groups {
				user_groups = append(user_groups, string(value.Name))
			}

			accessString, _ := utils.CreateJWTToken(user.Email, int(user.ID), user_permissions, user_groups, 60)
			refreshString, _ := utils.CreateJWTToken(user.Email, int(user.ID), user_permissions, user_groups, 65)

			data := TokenResponse{
				AccessToken:  accessString,
				RefreshToken: refreshString,
				TokenType:    "Bearer",
			}
			// Update the LastLogin field to the current time
			user.LastLogin = time.Now()

			// Save the updated user record
			if err := db.WithContext(tracer.Tracer).Save(&user).Error; err != nil {
				fmt.Println("Error updating user last login:", err)
			}


			return contx.JSON(http.StatusAccepted, common.ResponseHTTP{
				Success: true,
				Message: "Authorization Granted",
				Data:    data,
			})
		} else {
			return contx.JSON(http.StatusUnauthorized, common.ResponseHTTP{
				Success: false,
				Message: "Make sure You are Providing the Correct Credentials",
				Data:    "Authenthication Failed",
			})
		}
	case "refresh_token":
		claims, err := utils.ParseJWTToken(login_request.Token)
		email := claims.Email
		user_permissions := claims.Permissions
		user_groups := claims.Groups
		user_id := claims.UserID
		if err == nil {
			accessString, _ := utils.CreateJWTToken(email, user_id, user_permissions, user_groups, 60)
			refreshString, _ := utils.CreateJWTToken(email, user_id, user_permissions, user_groups, 65)
			data := TokenResponse{
				AccessToken:  accessString,
				RefreshToken: refreshString,
				TokenType:    "Bearer",
			}
			return contx.JSON(http.StatusAccepted, common.ResponseHTTP{
				Success: true,
				Message: "Authorization Granted",
				Data:    data,
			})
		}

		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: "Request Type Unknown",
			Data:    "Currently Not Implemented",
		})
	case "token_decode":
		claims, err := utils.ParseJWTToken(login_request.Token)

		if err == nil {
			return contx.JSON(http.StatusAccepted, common.ResponseHTTP{
				Success: true,
				Message: "Token decode sucessfull",
				Data:    claims,
			})
		}
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    "Unknown grant type",
		})
	default:
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: "Request Type Unknown",
			Data:    "Unknown grant type",
		})
	}
}
// ###############################################################################
// stats endpoint about schemas
// ###############################################################################
type DBStats struct {
	TotalUsers            int    {{.BackTick}}json:"total_users"{{.BackTick}}
	TotalGroups           int    {{.BackTick}}json:"total_groups"{{.BackTick}}
	TotalPermissions      int    {{.BackTick}}json:"total_permissions"{{.BackTick}}
	TotalGroupPermissions int    {{.BackTick}}json:"total_group_permissions"{{.BackTick}}
	TotalUserPermissions  int    {{.BackTick}}json:"total_user_permissions"{{.BackTick}}
	TotalUserGroups       int    {{.BackTick}}json:"total_user_groups"{{.BackTick}}
	TotalContentTypes     int    {{.BackTick}}json:"total_content_types"{{.BackTick}}
	TotalJwtSalts         int    {{.BackTick}}json:"total_jwt_salts"{{.BackTick}}
	TotalUserTokens       int    {{.BackTick}}json:"total_user_tokens"{{.BackTick}}
	LastUserLogin         string {{.BackTick}}json:"last_user_login"{{.BackTick}}
	LastPermissionCreated string {{.BackTick}}json:"last_permission_created"{{.BackTick}}
	LastGroupCreated      string {{.BackTick}}json:"last_group_created"{{.BackTick}}
	LastSaltGenerated     string {{.BackTick}}json:"last_salt_generated"{{.BackTick}}
	ActiveUsers           int    {{.BackTick}}json:"active_users"{{.BackTick}}
	InactiveUsers         int    {{.BackTick}}json:"inactive_users"{{.BackTick}}
	Superusers            int    {{.BackTick}}json:"superusers"{{.BackTick}}
	StaffUsers            int    {{.BackTick}}json:"staff_users"{{.BackTick}}
}

// Service Stat
// @Summary Auth
// @Description Stat
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseHTTP{data=DBStats{}}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /{{.CurrentAppName | replaceString}}/stats [get]
func DbStatEndpoint(contx echo.Context) error {
	//  Geting tracer
	tracer := contx.Get("tracer").(*observe.RouteTracer)

	//  Geting dbsession
	db := contx.Get("db").(*gorm.DB)

	// Getting Stats
	var stats DBStats
	res := db.WithContext(tracer.Tracer).Raw("SELECT * FROM db_stats").Scan(&stats)
	if res.Error != nil {

		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	return contx.JSON(http.StatusAccepted, common.ResponseHTTP{
		Success: true,
		Message: "Authorization Granted",
		Data:    stats,
	})

}
`

var loginTemplateFiber = `
package controllers

import (
	"net/http"
	"time"
	"fmt"

	"{{.ProjectName}}/{{ .CurrentAppName}}/models"
	"{{.ProjectName}}/{{ .CurrentAppName}}/utils"

	"{{.ProjectName}}/common"
	"{{.ProjectName}}/observe"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Login Request for Endpoint
type LoginPost struct {
	GrantType string {{.BackTick}}json:"grant_type" validate:"required" example:"authorization_code,refresh_token,token_decode"{{.BackTick}}
	Email     string {{.BackTick}}json:"email" validate:"email,min=6,max=32"{{.BackTick}}
	Password  string {{.BackTick}}json:"password"{{.BackTick}}
	Token     string {{.BackTick}}json:"token"{{.BackTick}}
}

// Access token Response
type TokenResponse struct {
	AccessToken  string {{.BackTick}}json:"access_token"{{.BackTick}}
	RefreshToken string {{.BackTick}}json:"refresh_token"{{.BackTick}}
	TokenType    string {{.BackTick}}json:"token_type"{{.BackTick}}
}

// Login is a function to login by EMAIL and ID
// @Summary Auth
// @Description Login
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body LoginPost true "Login"
// @Success 200 {object} common.ResponseHTTP{data=TokenResponse{}}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /blue_auth/login [post]
func Login(contx *fiber.Ctx) error {
	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// validator initialization
	validate := validator.New()

	//validating post data
	login_request := new(LoginPost)

	//first parse request data
	if err := contx.BodyParser(&login_request); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(login_request); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	switch login_request.GrantType {
	case "authorization_code":
		var user models.User
		var permissions []models.Permission
		var groups []models.Group

		// Getting User
		res := db.WithContext(tracer.Tracer).Model(&models.User{}).Where("email = ? AND is_active = ?", login_request.Email, true).First(&user)
		if res.Error != nil {

			return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
				Success: false,
				Message: res.Error.Error(),
				Data:    nil,
			})
		}
		// Getting User Permissions
		permission_join_string := "INNER JOIN user_permissions up ON permissions.id = up.permission_id"
		res_per := db.WithContext(tracer.Tracer).Model(&models.Permission{}).Joins(permission_join_string).Where("up.user_id = ?", user.ID).Scan(&permissions)
		if res_per.Error != nil {
			return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
				Success: false,
				Message: res.Error.Error(),
				Data:    nil,
			})
		}

		// Getting User Groups
		groups_join_string := "INNER JOIN user_groups ug ON groups.id = ug.group_id"
		res_gr := db.WithContext(tracer.Tracer).Model(&models.Group{}).Joins(groups_join_string).Where("ug.user_id = ?", user.ID).Scan(&groups)
		if res_gr.Error != nil {
			return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
				Success: false,
				Message: res.Error.Error(),
				Data:    nil,
			})
		}

		// Checking Password match
		if utils.PasswordsMatch(user.Password, login_request.Password) {
			// preparing permissions for JWT token hash
			user_permissions := make([]string, 0, 20)
			for _, value := range permissions {
				user_permissions = append(user_permissions, string(value.Codename))
			}
			// preparing groups for JWT token hash
			user_groups := make([]string, 0, 20)
			for _, value := range groups {
				user_groups = append(user_groups, string(value.Name))
			}

			accessString, _ := utils.CreateJWTToken(user.Email, int(user.ID), user_permissions, user_groups, 60)
			refreshString, _ := utils.CreateJWTToken(user.Email, int(user.ID), user_permissions, user_groups, 65)

			data := TokenResponse{
				AccessToken:  accessString,
				RefreshToken: refreshString,
				TokenType:    "Bearer",
			}
			// Update the LastLogin field to the current time
			user.LastLogin = time.Now()

			// Save the updated user record
			if err := db.WithContext(tracer.Tracer).Save(&user).Error; err != nil {
				fmt.Println("Error updating user last login:", err)
			}

			return contx.Status(http.StatusAccepted).JSON(common.ResponseHTTP{
				Success: true,
				Message: "Authorization Granted",
				Data:    data,
			})
		} else {
			return contx.Status(http.StatusUnauthorized).JSON(common.ResponseHTTP{
				Success: false,
				Message: "Make sure You are Providing the Correct Credentials",
				Data:    "Authenthication Failed",
			})
		}
	case "refresh_token":
		claims, err := utils.ParseJWTToken(login_request.Token)
		email := claims.Email
		user_permissions := claims.Permissions
		user_groups := claims.Groups
		user_id := claims.UserID
		if err == nil {
			accessString, _ := utils.CreateJWTToken(email, user_id, user_permissions, user_groups, 60)
			refreshString, _ := utils.CreateJWTToken(email, user_id, user_permissions, user_groups, 65)
			data := TokenResponse{
				AccessToken:  accessString,
				RefreshToken: refreshString,
				TokenType:    "Bearer",
			}
			return contx.Status(http.StatusAccepted).JSON(common.ResponseHTTP{
				Success: true,
				Message: "Authorization Granted",
				Data:    data,
			})
		}

		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Request Type Unknown",
			Data:    "Currently Not Implemented",
		})
	case "token_decode":
		claims, err := utils.ParseJWTToken(login_request.Token)

		if err == nil {
			return contx.Status(http.StatusAccepted).JSON(common.ResponseHTTP{
				Success: true,
				Message: "Token decode sucessfull",
				Data:    claims,
			})
		}
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    "Unknown grant type",
		})
	default:
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Request Type Unknown",
			Data:    "Unknown grant type",
		})
	}
}

// ###############################################################################
// stats endpoint about schemas
// ###############################################################################
type DBStats struct {
	TotalUsers            int    {{.BackTick}}json:"total_users"{{.BackTick}}
	TotalGroups           int    {{.BackTick}}json:"total_groups"{{.BackTick}}
	TotalPermissions      int    {{.BackTick}}json:"total_permissions"{{.BackTick}}
	TotalGroupPermissions int    {{.BackTick}}json:"total_group_permissions"{{.BackTick}}
	TotalUserPermissions  int    {{.BackTick}}json:"total_user_permissions"{{.BackTick}}
	TotalUserGroups       int    {{.BackTick}}json:"total_user_groups"{{.BackTick}}
	TotalContentTypes     int    {{.BackTick}}json:"total_content_types"{{.BackTick}}
	TotalJwtSalts         int    {{.BackTick}}json:"total_jwt_salts"{{.BackTick}}
	TotalUserTokens       int    {{.BackTick}}json:"total_user_tokens"{{.BackTick}}
	LastUserLogin         string {{.BackTick}}json:"last_user_login"{{.BackTick}}
	LastPermissionCreated string {{.BackTick}}json:"last_permission_created"{{.BackTick}}
	LastGroupCreated      string {{.BackTick}}json:"last_group_created"{{.BackTick}}
	LastSaltGenerated     string {{.BackTick}}json:"last_salt_generated"{{.BackTick}}
	ActiveUsers           int    {{.BackTick}}json:"active_users"{{.BackTick}}
	InactiveUsers         int    {{.BackTick}}json:"inactive_users"{{.BackTick}}
	Superusers            int    {{.BackTick}}json:"superusers"{{.BackTick}}
	StaffUsers            int    {{.BackTick}}json:"staff_users"{{.BackTick}}
}

// Service Stat
// @Summary Auth
// @Description Stat
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} common.ResponseHTTP{data=DBStats{}}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /blue_auth/stats [get]
func DbStatEndpoint(contx *fiber.Ctx) error {
	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)

	//  Getting Database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// Getting Stats
	var stats DBStats
	res := db.WithContext(tracer.Tracer).Raw("SELECT * FROM db_stats").Scan(&stats)
	if res.Error != nil {

		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	return contx.Status(http.StatusAccepted).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Authorization Granted",
		Data:    stats,
	})

}
`
