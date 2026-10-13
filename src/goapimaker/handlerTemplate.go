package goapimaker

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/MaminirinaEdwino/turbostackapi/src/utils"
)

func SelectTemplate(HandlerName string, DbCaller string, ResponseType string, Query string, ScanValue string, ResponseWriter string) string {
	return fmt.Sprintf(`
func %s(w http.ResponseWriter, r *http.Request){
	%s
	var res []models.%s

	rows, _ := db.Query("%s")
	
	for rows.Next(){
		var tmp models.%s
		rows.Scan(%s)
		res = append(res, tmp)
	}
	%s
}
	`, HandlerName, DbCaller, ResponseType, Query, ResponseType, ScanValue, ResponseWriter)
}

func DBCallerTemplate() string {
	return `
db := config.DB
defer db.Close()
	`
}
func DBCallerTemplateWebAPp() string {
	return `
db := config.ConnectDB()
defer db.Close()
	`
}

func WriteResponseWriter() string {
	return `
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(res)
	`
}

func SelectBytemplate() *template.Template {
	content := `
func {{ .FuncName }}(w http.ResponseWriter, r *http.Request){
	{{ .Params }} := r.PathValue("{{ .Params }}")
	var res models.{{ .StructName }}
	{{ .DbCallerHandler }}
	rows,_ := db.Query("{{ .Query }}", {{ .Params }})
	rows.Next()
	rows.Scan({{ .ScanParams }})
	{{ .ResponseWriter }}
}
	`
	temp := template.New(content)
	temp.Parse(content)
	return temp
}

func PutTemplate() *template.Template {
	content := `
func {{ .FuncName }}(w http.ResponseWriter, r *http.Request){
	var body models.{{ .EndPointName }}
	var res models.{{ .EndPointName }}
	{{ .Params }} := r.PathValue("{{ .Params }}")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		log.Fatal(err)
	}
	{{ .DbCallerHandler }}
	rows, err := db.Query("{{ .PutQuery }}", {{ .Params }})
	if err != nil {
		log.Fatal(err)
	}
	rows.Next()
	rows.Scan({{ .ScanParams }})
	{{ .ResponseWriter }}
}		
`
	temp := template.New(content)
	temp.Parse(content)
	return temp
}

func PutHandler(modelName, epName, sgbd string, attrs []string, ScanParamsWriter string, params string) string {
	tmp := PutTemplate()

	var tmpBuffer bytes.Buffer
	data := struct {
		EndPointName    string
		DbCallerHandler string
		PutQuery        string
		ScanParams      string
		ResponseWriter  string
		Params          string
		FuncName        string
	}{
		EndPointName:    utils.ToUpperFirstLetter(modelName),
		DbCallerHandler: DBCallerHandler(sgbd),
		PutQuery:        Update(epName, attrs, sgbd),
		ScanParams:      ScanParamsWriter,
		ResponseWriter:  WriteResponseWriter(),
		Params:          params,
		FuncName:        strings.ReplaceAll(epName, " ", "_"),
	}
	err := tmp.Execute(&tmpBuffer, data)
	if err != nil {
		fmt.Println(err)
	}
	return tmpBuffer.String()
}

func DeleteHandler(epName, sgbd, params, tableName string) string {
	return fmt.Sprintf(`
func %s(w http.ResponseWriter, r *http.Request){
	%s := r.PathValue("%s")
	type response struct{
		Message string
	}
	%s
	rows,_ := db.Query("%s", %s)
	rows.Next()
	res := response{
		Message: "users deleted",
	}
	%s
}
	`, epName, params, params, DBCallerHandler(sgbd), Delete(tableName, sgbd), params, WriteResponseWriter())
}

func InsertHandler(epName, sgbd, tableName string, attr []string) string {
	var bodyParam []string
	for _, val := range attr {
		bodyParam = append(bodyParam, fmt.Sprintf("body.%s", utils.ToUpperFirstLetter(val)))
	}
	return fmt.Sprintf(`
	func %s(w http.ResponseWriter, r *http.Request){
		%s
		%s
	res, err := db.Exec("%s", %s)
		%s
		%s
	}`+"\n",
		epName,
		WriteBodyDecodeur(tableName),
		DBCallerHandler(sgbd),
		Insert(tableName, attr),
		strings.Join(bodyParam, ", "),
		utils.WriteErrorCheker("insert error"),
		WriteResponseWriter())
}

func WriteBodyDecodeur(modelName string) string {
	var tmpBuffer bytes.Buffer

	tmp := BodyDecodeurTemplate()

	data := struct {
		EndPointName string
		ErrorChecker string
	}{
		EndPointName: utils.ToUpperFirstLetter(modelName),
		ErrorChecker: utils.WriteErrorCheker("Parsing Error"),
	}

	err := tmp.Execute(&tmpBuffer, data)
	utils.ErrorChecker(err)
	return tmpBuffer.String()
}

func BodyDecodeurTemplate() *template.Template {
	content := `
var body models.{{ .EndPointName }} 
decoder := json.NewDecoder(r.Body) 
err := decoder.Decode(&body)
{{ .ErrorChecker }}
	`
	temp := template.New(content)
	temp.Parse(content)
	return temp
}
