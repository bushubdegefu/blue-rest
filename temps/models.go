package temps

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func ModelDataFrame() {

	// ############################################################
	models_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(gmodelTemplate)
	if err != nil {
		panic(err)
	}

	migration_function_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(migrationFuncTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	// #################################################
	err = os.MkdirAll("models", os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, model := range RenderData.Models {

		folder_path := fmt.Sprintf("models/%v.go", model.Name)
		folder_path = strings.ToLower(folder_path)
		models_file, err := os.Create(folder_path)
		if err != nil {
			panic(err)
		}

		err = models_tmpl.Execute(models_file, model)
		if err != nil {
			panic(err)
		}
		models_file.Close()

	}

	init_file, err := os.Create("models/init.go")
	if err != nil {
		panic(err)
	}

	err = migration_function_tmpl.Execute(init_file, RenderData)
	if err != nil {
		panic(err)
	}
	defer init_file.Close()

	// ############################################################
	helper_function_tmpl, err := template.New("RenderData").Parse(helperFunctionsTemplate)
	if err != nil {
		panic(err)
	}

	helper_file, err := os.Create("models/helper.go")
	if err != nil {
		panic(err)
	}

	err = helper_function_tmpl.Execute(helper_file, RenderData)
	if err != nil {
		panic(err)
	}
	defer init_file.Close()

}

func MigrationInit() {

	// ###################################################
	err := os.MkdirAll("models", os.ModePerm)
	if err != nil {
		panic(err)
	}
	// ############################################################
	migration_function_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(migrationFuncTemplate)
	if err != nil {
		panic(err)
	}

	init_file, err := os.Create("models/init.go")
	if err != nil {
		panic(err)
	}

	err = migration_function_tmpl.Execute(init_file, RenderData)
	if err != nil {
		panic(err)
	}
	defer init_file.Close()

}

func DbConnDataFrame() {
	InitProjectJSON()
	//  creating database connection folder
	// ############################################################
	database_tmpl, err := template.New("RenderData").Funcs(FuncMap).Parse(databaseTemplate)
	if err != nil {
		panic(err)
	}

	// create database folder if does not exist
	err = os.MkdirAll("database", os.ModePerm)
	if err != nil {
		panic(err)
	}

	database_conn_file, err := os.Create("database/database.go")
	if err != nil {
		panic(err)
	}
	defer database_conn_file.Close()

	err = database_tmpl.Execute(database_conn_file, ProjectSettings)
	if err != nil {
		panic(err)
	}

}

var gmodelTemplate = `
package models

import (
	"time"
	"fmt"
	"gorm.io/gorm"
	{{- $break_3 := false }}
	{{- range .Fields}}
	{{- if eq .Name "UUID" }}
	"github.com/google/uuid"
	{{- $break_3 = true }}
	{{- end}}
	{{- end}}
	"database/sql"
	"log"
)

// {{.Name}} Database model info
// @Description App type information
type {{.Name}} struct {
	// The following fields will be ignored by Swagger
   	CreatedAt time.Time {{.BackTick}}json:"created_at,omitempty"{{.BackTick}}
    UpdatedAt time.Time {{.BackTick}}json:"updated_at,omitempty"{{.BackTick}}
    {{range .Fields}} {{.Name}} {{.Type}}  {{.BackTick}}{{.Annotation}}{{.BackTick}}
{{end}}}

{{- $hasUUID := false }}
{{- range .Fields}}
    {{- if eq .Name "UUID" }}
        {{- $hasUUID = true }}
    {{- end }}
{{- end }}

{{- $hasPassword := false }}
{{- range .Fields}}
    {{- if eq .Name "Password" }}
        {{- $hasPassword = true }}
    {{- end }}
{{- end }}

{{- if not $hasUUID }}
func (entity *{{.Name}}) BeforeCreate(tx *gorm.DB) (err error) {
   	{{- if $hasPassword }}
  		entity.Password = HashFunc(entity.Password)
   	{{- end }}
    entity.CreatedAt = time.Now()
    return
}
{{- end }}

{{- $break_4 := false }}
{{- range .Fields}}
{{- if eq .Name "UUID" }}
func (entity *{{.NormalModelName}}) BeforeCreate(tx *gorm.DB) (err error) {
	gen, _ := uuid.NewV7()
	entity.CreatedAt = time.Now();
	id := gen.String()
	entity.UUID = id
	{{- if $hasPassword }}
  		entity.Password = HashFunc(entity.Password)
   	{{- end }}
	return
}
{{- $break_4 = true }}
{{- end}}
{{- end}}

func (entity *{{.Name}}) BeforeUpdate(tx *gorm.DB) (err error) {
	entity.UpdatedAt = time.Now();
	return
}

func (entity *{{.Name}}) Populate(tx *gorm.DB) {
	// Create ContentType for User model
	contentType := ContentType{
		AppLabel: "{{ .AppName | replaceString }}",
		Model:    "{{.LowerName}}",
	}
	if err := tx.Create(&contentType).Error; err != nil {
		log.Fatalf("Failed to create ContentType: %v", err)
	}

	// Create Permissions for User model
	permissions := []Permission{
		{ContentTypeID: sql.NullInt64{Int64: int64(contentType.ID), Valid: true}, Codename: "{{ .AppName | replaceString }}_can_add_{{.LowerName}}"},
		{ContentTypeID: sql.NullInt64{Int64: int64(contentType.ID), Valid: true}, Codename: "{{ .AppName | replaceString }}_can_view_{{.LowerName}}"},
		{ContentTypeID: sql.NullInt64{Int64: int64(contentType.ID), Valid: true}, Codename: "{{ .AppName | replaceString }}_can_change_{{.LowerName}}"},
		{ContentTypeID: sql.NullInt64{Int64: int64(contentType.ID), Valid: true}, Codename: "{{ .AppName | replaceString }}_can_delete_{{.LowerName}}"},
	}

	for _, permission := range permissions {
		if err := tx.Create(&permission).Error; err != nil {
			log.Fatalf("Failed to create Permission: %v", err)
		}
	}

	fmt.Println("Populated ContentType and Permissions for {{.Name}} request actions successfully")
}

// {{.Name}}Post model info
// @Description {{.Name}}Post type information
type {{.Name}}Post struct {
  	{{range .Fields}} {{- if .Post}} {{.Name}} {{.Type}} {{.BackTick}}{{.Annotation}}{{.BackTick}}{{- end}}
{{end}}}

// {{.Name}}Get model info
// @Description {{.Name}}Get type information
type {{.Name}}Get struct {
	{{range .Fields}} {{- if .Get}}	{{.Name}} {{.Type}}  {{.BackTick}}{{.Annotation}}{{.BackTick}} {{- end}}
	{{end}}
	CreatedAt time.Time {{.BackTick}}json:"created_at,omitempty"{{.BackTick}}
	UpdatedAt time.Time {{.BackTick}}json:"updated_at,omitempty"{{.BackTick}}
}

// {{.Name}}Put model info
// @Description {{.Name}}Put type information
type {{.Name}}Put struct {
	{{range .Fields}} {{- if .Put}} {{.Name}} {{.Type}}  {{.BackTick}}{{.Annotation}}{{.BackTick}}{{- end}}
{{end}}}

// {{.Name}}Patch model info
// @Description {{.Name}}Patch type information
type {{.Name}}Patch struct {
	{{range .Fields}}{{- if .Patch}}{{.Name}} {{.Type}}  {{.BackTick}}{{.Annotation}}{{.BackTick}}{{- end}}
{{end}}}
`

var migrationFuncTemplate = `
package models

import (
	"fmt"
	"log"

	"{{.ProjectName}}/database"
	"{{.ProjectName}}/configs"
)

func InitDatabase(test_flag bool) {
	if !test_flag {
		configs.NewEnvFile("./configs")
	}
	database, err  := database.ReturnSession("{{ .AppName | replaceString }}")
	fmt.Println("Connection Opened to Database")
	if err == nil {
		if err := database.AutoMigrate(
			{{- range .Models}}
			&{{.Name}}{},
			{{- end}}
		); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Database Migrated")
	} else {
		panic(err)
	}
}

func CleanDatabase(test_flag bool) {
	if !test_flag {
		configs.NewEnvFile("./configs")
	}
	database, err := database.ReturnSession("{{ .AppName | replaceString }}")
	if err == nil {
		fmt.Println("Connection Opened to Database")
		fmt.Println("Dropping Models if Exist")
		err := database.Migrator().DropTable(
		{{- range .Models}}
			&{{.Name}}{},
		{{- end}}
		)
		if err != nil {
			fmt.Println("Error dropping tables:", err)
		}
		fmt.Println("Database Cleaned")
	} else {
		panic(err)
	}
}

{{ if eq .AuthAppName .AppName }}
func CreateSuperUser() {
	configs.NewEnvFile("./configs")
	db, err := database.ReturnSession("{{.AuthAppName | replaceString }}")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Create the superuser
	user := &User{
		Username:    "superuser",
		Email:       "superuser@mail.com",
		Password:    "default@123",
		IsSuperuser: true,
		IsStaff:     true,
		IsActive:    true,
		FirstName:   "Super",
		LastName:    "Admin",
	}

	// Insert the user into the database
	if err := db.Create(user).Error; err != nil {
		fmt.Printf("failed to create superuser: %v\n", err)
	}
	db.Commit()
	fmt.Println("Superuser created successfully")

}
{{- end}}

func Populate(test_flag bool) {
	if !test_flag {
		configs.NewEnvFile("./configs")
	}
	db, err := database.ReturnSession("{{.AuthAppName | replaceString }}")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	{{- range .Models}}
		(&{{.Name}}{}).Populate(db)
	{{- end}}
}




`

var helperFunctionsTemplate = `
package models

import (
	"crypto/sha512"
	"encoding/hex"

	"{{.ProjectName}}/configs"
)

// Combine password and salt then hash them using the SHA-512
func HashFunc(password string) string {

	// var salt []byte
	// get salt from env variable
	salt := []byte(configs.AppConfig.Get("SECRETE_SALT"))

	// Convert password string to byte slice
	var pwdByte = []byte(password)

	// Create sha-512 hasher
	var sha512 = sha512.New()

	pwdByte = append(pwdByte, salt...)

	sha512.Write(pwdByte)

	// Get the SHA-512 hashed password
	var hashedPassword = sha512.Sum(nil)

	// Convert the hashed to hex string
	var hashedPasswordHex = hex.EncodeToString(hashedPassword)
	return hashedPasswordHex
}

`

var databaseTemplate = `
package database

import (
	"log"
	"os"
	"time"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"{{.ProjectName}}/configs"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	DBConn *gorm.DB
)

func GormLoggerFile(app_name string) (*os.File, error) {
	log_file_name := fmt.Sprintf("%s_gorm.log", app_name)
	gormLogFile, gerr := os.OpenFile(log_file_name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if gerr != nil {
		log.Fatalf("error opening file: %v", gerr)
	}
	return gormLogFile, nil
}

func ReturnSession(app_name string) (*gorm.DB,error) {

	//  setting up database connection based on DB type
	env_name := fmt.Sprintf("%s_DB_TYPE", strings.ToUpper(app_name))
	app_env := configs.AppConfig.Get(env_name)
	//  This is file to output gorm logger on to
	gormlogger,_ := GormLoggerFile(app_name)
	gormFileLogger := log.Logger{}
	gormFileLogger.SetOutput(gormlogger)
	gormFileLogger.Writer()


	gormLogger := log.New(gormFileLogger.Writer(), "\r\n", log.LstdFlags|log.Ldate|log.Ltime|log.Lshortfile)
	newLogger := logger.New(
		gormLogger, // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			Colorful:                  true,        // Enable color
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			// ParameterizedQueries:      true,        // Don't include params in the SQL log

		},
	)

	var DBSession *gorm.DB

	switch app_env {
	case "postgres":
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  configs.AppConfig.Get(fmt.Sprintf("%s_POSTGRES_URI", strings.ToUpper(app_name))),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage,

		}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                 newLogger,
			SkipDefaultTransaction: true,
		})
		if err != nil {
			panic(err)
		}

		sqlDB,err := db.DB()
		if err != nil {
			fmt.Printf("Error during connecting to database: %v\n", err)
			return nil, err
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)

		DBSession = db
	case "sqlite":
		//  this is sqlite connection
		db, _ := gorm.Open(sqlite.Open(configs.AppConfig.Get(fmt.Sprintf("%s_SQLLITE_URI", strings.ToUpper(app_name)))), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                 newLogger,
			SkipDefaultTransaction: true,
		})

		sqlDB,err := db.DB()
		if err != nil {
			fmt.Printf("Error during connecting to database: %v\n", err)
			return nil, err
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)
		DBSession = db
	case "mysql":
		db, _ := gorm.Open(mysql.New(mysql.Config{
			DSN:                       configs.AppConfig.Get(fmt.Sprintf("%s_MYSQL_URI", strings.ToUpper(app_name))), // data source name
			DefaultStringSize:         256,                                // default size for string fields
			DisableDatetimePrecision:  true,                               // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,                               // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,                               //  when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false,                              // auto configure based on currently MySQL version
		}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                 newLogger,
			SkipDefaultTransaction: true,
		})

		sqlDB,err := db.DB()
		if err != nil {
			fmt.Printf("Error during connecting to database: %v\n", err)
			return nil, err
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)
		DBSession = db
	default:
			return nil, fmt.Errorf("database type not supported")

	}

	// Mouting Otel tracer plugin on gorm Session
	err := DBSession.Use(tracing.NewPlugin())
	if err != nil {
		fmt.Printf("Error during connecting to otel plugin: %v\n", err)

	}
	return DBSession,nil

}

`
