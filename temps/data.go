package temps

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var RenderData Data
var ProjectSettings ProjectSetting

// SetBackTick assigns a backtick to the Data struct.
func (d *Data) SetBackTick() {
	d.BackTick = "`"
}

// CapitalizeFirstLetter capitalizes the first letter of the string.
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}

// Data struct holds information about the project.
type Data struct {
	ProjectName    string   `json:"project_name"`
	FrameName      string   `json:"frame_name"`
	AppName        string   `json:"app_name"`
	PackageAppName string   `json:"package_app_name"`
	BackTick       string   `json:"back_tick"`
	Models         []Model  `json:"models"`
	AppNames       []string `json:"app_names"`
	AuthAppName    string   `json:"auth_app_name"`
}

// ProjectSetting contains project configuration.
type ProjectSetting struct {
	ProjectName    string   `json:"project_name"`
	AppNames       []string `json:"app_names"`
	CurrentAppName string   `json:"current_app_name"`
	BackTick       string   `json:"back_tick"`
	AuthAppName    string   `json:"auth_app_name"`
}

// Contains checks if a string is in the AppNames slice.
func (p *ProjectSetting) Contains(str string) bool {
	return slices.Contains(p.AppNames, str)
}

// AppendAppName appends a new app name to AppNames if not already present.
func (p *ProjectSetting) AppendAppName(appName, authAppName string) error {
	if p.Contains(appName) {
		fmt.Println("App already exists, please use another name.")
		return nil
	}
	p.AppNames = append(p.AppNames, appName)

	if authAppName != "" {
		p.AuthAppName = authAppName
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling struct to JSON: %v", err)
	}

	if err := writeToFile("project.json", data); err != nil {
		return err
	}
	return nil
}

// Model represents a data model.
type Model struct {
	Name         string         `json:"name"`
	LowerName    string         `json:"lower_name"`
	RlnModel     []string       `json:"rln_model"`
	BackTick     string         `json:"back_tick"`
	Fields       []Field        `json:"fields"`
	ProjectName  string         `json:"project_name"`
	AppName      string         `json:"app_name"`
	AuthAppName  string         `json:"auth_app_name"`
	SearchFields []string       `json:"search_fields"`
	Relations    []Relationship `json:"relations"`
}

// Relationship defines model relationships.
type Relationship struct {
	TableName       string  `json:"table_name"`
	AppName         string  `json:"app_name"`
	ParentName      string  `json:"parent_name"`
	LowerParentName string  `json:"lower_parent_name"`
	FieldName       string  `json:"field_name"`
	LowerFieldName  string  `json:"lower_field_name"`
	ParentFields    []Field `json:"parent_fields"`
	ChildFields     []Field `json:"child_fields"`
	MtM             bool    `json:"mtm"`
	OtM             bool    `json:"otm"`
	MtO             bool    `json:"mto"`
	BackTick        string  `json:"back_tick"`
}

// Field represents a model field.
type Field struct {
	NormalModelName  string `json:"normal_model_name"`
	ModelName        string `json:"model_name"`
	Name             string `json:"name"`
	LowerName        string `json:"lower_name"`
	Type             string `json:"type"`
	UpperType        string `json:"upper_type"`
	Annotation       string `json:"annotation"`
	MongoAnnotation  string `json:"mongo_annotation"`
	CurdFlag         string `json:"curd_flag"`
	Get              bool   `json:"get"`
	Post             bool   `json:"post"`
	Patch            bool   `json:"patch"`
	Put              bool   `json:"put"`
	OtM              bool   `json:"otm"`
	MtM              bool   `json:"mtm"`
	ProjectName      string `json:"project_name"`
	AppName          string `json:"app_name"`
	AuthAppName      string `json:"auth_app_name"`
	BackTick         string `json:"back_tick"`
	RandomFieldValue string `json:"random_field_value"`
}

// generateRandomValue generates random data based on the field type.
func generateRandomValue(fieldType string) string {
	switch fieldType {
	case "string":
		return fmt.Sprintf("\"%s\"", randomString(10))
	case "int", "int32", "int64":
		return fmt.Sprintf("%d", rand.Intn(1000))
	case "float64":
		return fmt.Sprintf("%f", rand.Float64()*100.0)
	case "bool":
		return fmt.Sprintf("%t", rand.Intn(2) == 0)
	case "time.Time":
		start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := time.Now()
		duration := end.Sub(start)
		randomDuration := time.Duration(rand.Int63n(int64(duration)))
		randomTime := start.Add(randomDuration)
		return fmt.Sprintf("\"%s\"", randomTime.Format(time.RFC3339))
	case "ID":
		return fmt.Sprintf("%v", rand.Intn(1000))
	case "sql.NullInt64":
		return fmt.Sprintf("%f", rand.Float64()*100.0)
	default:
		return "\"\""
	}
}

// writeToFile writes data to a file.
func writeToFile(fileName string, data []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	return nil
}

// LoadData loads data from a JSON configuration file.
func LoadData(fileName string) error {
	if fileName == "" {
		fileName = "config.json"
	}

	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening JSON file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&RenderData); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	initializeRenderData()

	return nil
}

// initializeRenderData initializes RenderData with default values and relationships.
func initializeRenderData() {

	// Ensure RenderData.BackTick is set correctly.
	RenderData.BackTick = "`"
	RenderData.AuthAppName = ProjectSettings.AuthAppName
	// Loop through each model in RenderData.Models
	for i := range RenderData.Models { // Use range with index to modify the model directly
		// Directly modify the model in the slice using its index.
		model := &RenderData.Models[i] // Get a reference to the model

		// Set the fields on the model
		model.LowerName = strings.ToLower(model.Name)
		model.AppName = RenderData.AppName         // Set AppName from RenderData
		model.AuthAppName = RenderData.AuthAppName // Set AppName from RenderData
		model.ProjectName = RenderData.ProjectName // Set ProjectName from RenderData
		model.BackTick = "`"                       // Set backtick for the model

		// Initialize fields and relationships
		initializeFieldsAndRelationships(model)
	}
}

// initializeFieldsAndRelationships initializes fields and relationships for a model.
func initializeFieldsAndRelationships(model *Model) {
	for i := range model.Fields {
		field := &model.Fields[i]
		// "Get$Post$Patch$Put$OtM$MtM"
		cf := strings.Split(field.CurdFlag, "$")
		field.BackTick = "`"
		field.AuthAppName = RenderData.AuthAppName
		field.ModelName = strings.ToLower(model.Name)
		field.NormalModelName = model.Name
		field.LowerName = strings.ToLower(field.Name)
		field.UpperType = CapitalizeFirstLetter(field.Type)
		field.Get, _ = strconv.ParseBool(cf[0])
		field.Post, _ = strconv.ParseBool(cf[1])
		field.Patch, _ = strconv.ParseBool(cf[2])
		field.Put, _ = strconv.ParseBool(cf[3])
		field.AppName = RenderData.AppName
		field.ProjectName = RenderData.ProjectName
		field.RandomFieldValue = generateRandomValue(field.Type)
	}

	model.Relations = initializeRelations(model)
}

// initializeRelations initializes the relationships for a model.
func initializeRelations(model *Model) []Relationship {
	relations := make([]Relationship, 0)
	for _, relation := range model.RlnModel {
		rlf := strings.Split(relation, "$")
		curRelation := Relationship{
			AppName:         RenderData.AppName,
			ParentName:      model.Name,
			LowerParentName: model.LowerName,
			FieldName:       rlf[0],
			LowerFieldName:  strings.ToLower(rlf[0]),
			MtM:             rlf[1] == "mtm",
			OtM:             rlf[1] == "otm",
			MtO:             rlf[1] == "mto",
			BackTick:        "`",
		}
		if len(rlf) > 2 {
			curRelation.TableName = rlf[2]
		}
		curRelation.ParentFields = model.Fields
		relations = append(relations, curRelation)
	}
	return relations
}

// CommonProjectName saves the project name into the project JSON file.
func CommonProjectName(projectName string, authAppName string) {
	if authAppName == "" {
		authAppName = "django_auth"
	}
	data, _ := json.MarshalIndent(&ProjectSetting{ProjectName: projectName, AuthAppName: authAppName}, "", "  ")
	if err := writeToFile("project.json", data); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

// GetProjectName retrieves the project name from the project JSON file.
func GetProjectName() string {
	file, err := os.Open("project.json")
	if err != nil {
		panic("Project not initialized: open project.json: no such file or directory")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&ProjectSettings); err != nil {
		panic("Error decoding project.json: " + err.Error())
	}
	return ProjectSettings.ProjectName
}

// InitProjectJSON initializes the project JSON file if not already present.
func InitProjectJSON() {
	file, err := os.Open("project.json")
	if err != nil {
		fmt.Println("project.json not found, please initialize it.")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&ProjectSettings); err != nil {
		fmt.Println("project.json not found, please initialize it.")
	}
}
