package temps

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func ModelDataFrame() {

	// Open the JSON file
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer file.Close() // Defer closing the file until the function returns

	// ############################################################
	models_tmpl, err := template.New("RenderData").Parse(gmodelTemplate)
	if err != nil {
		panic(err)
	}

	migration_function_tmpl, err := template.New("RenderData").Parse(migrationFuncTemplate)
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

	//  creating database connection folder
	// ############################################################
	database_tmpl, err := template.New("RenderData").Parse(databaseTemplate)
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

	err = database_tmpl.Execute(database_conn_file, RenderData)
	if err != nil {
		panic(err)
	}
}

func DbConnDataFrame() {

	//  creating database connection folder
	// ############################################################
	database_tmpl, err := template.New("RenderData").Parse(databaseTemplate)
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

	err = database_tmpl.Execute(database_conn_file, RenderData)
	if err != nil {
		panic(err)
	}

}

var gmodelTemplate = `
package models

import (
	"time"
	"gorm.io/gorm"
	{{- $break_3 := false }}
	{{- range .Fields}}
	{{- if eq .Name "UUID" }}
	"github.com/google/uuid"
	{{- $break_3 = true }}
	{{- end}}
	{{- end}}
	{{- $break_2 := false }}
	{{- range .Fields}}
	{{- if eq .Type "sql.NullInt64" }}
	"database/sql"
	{{- $break_2 = true }}
	{{- end}}
	{{- end}}
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

{{- if not $hasUUID }}
    func (entity *{{.Name}}) BeforeCreate(tx *gorm.DB) (err error) {
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
	return
}
{{- $break_4 = true }}
{{- end}}
{{- end}}

func (entity *{{.Name}}) BeforeUpdate(tx *gorm.DB) (err error) {
	entity.UpdatedAt = time.Now();
	return
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
	database, err  := database.ReturnSession()
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
	database, err := database.ReturnSession()
	if err == nil {
		fmt.Println("Connection Opened to Database")
		fmt.Println("Dropping Models if Exist")
		database.Migrator().DropTable(
		{{range .Models}}
			&{{.Name}}{},
		{{end}}
		)

		fmt.Println("Database Cleaned")
	} else {
		panic(err)
	}
}


`

var databaseTemplate = `
package database

import (
	"log"
	"os"
	"time"
	"fmt"

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

func GormLoggerFile() (*os.File,error) {

	gormLogFile, gerr := os.OpenFile("blue-gorm.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if gerr != nil {
		log.Fatalf("error opening file: %v", gerr)
	}
	return gormLogFile,nil
}

func ReturnSession() (*gorm.DB,error) {

	//  setting up database connection based on DB type

	app_env := configs.AppConfig.Get("DB_TYPE")
	//  This is file to output gorm logger on to
	gormlogger,_ := GormLoggerFile()
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
			DSN:                  configs.AppConfig.Get("POSTGRES_URI"),
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
		db, _ := gorm.Open(sqlite.Open(configs.AppConfig.Get("SQLLITE_URI")), &gorm.Config{
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
			DSN:                       configs.AppConfig.Get("MYSQL_URI"), // data source name
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
	case "":
		//  this is sqlite connection
		db, _ := gorm.Open(sqlite.Open("goframe-2.db"), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                 newLogger,
			SkipDefaultTransaction: true,
		})

		sqlDB, err:= db.DB()
		if err != nil {
			fmt.Printf("Error during connecting to database: %v\n", err)
			return nil, err
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)
		DBSession = db
	default:
		//  this is sqlite connection
		db, _ := gorm.Open(sqlite.Open(configs.AppConfig.Get("SQLITE_URI")), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                 newLogger,
			SkipDefaultTransaction: true,
		})

		sqlDB, err:= db.DB()
		if err != nil {
			fmt.Printf("Error during connecting to database: %v\n", err)
			return nil, err
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(5 * time.Second)
		DBSession = db

	}

	DBSession.Use(tracing.NewPlugin())
	return DBSession,nil

}

`
