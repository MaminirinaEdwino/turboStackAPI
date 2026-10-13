package entity

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/MaminirinaEdwino/turbostackapi/src/config"
	"github.com/MaminirinaEdwino/turbostackapi/src/goapimaker"
	"github.com/MaminirinaEdwino/turbostackapi/src/utils"
)

type webAppMaker struct {
	ProjectName string
	WebApp      WebApp
	Techno      string
	Api         RestApi
	BDD         BDD
}

func (wap *webAppMaker) SetupArch() {
	projectpath := fmt.Sprintf("%s/web-app/", wap.ProjectName)
	folderList := []string{
		"src/views",
		"src/controller",
		"src/models",
		"src/static",
		"src/router",
		"src/config",
		"src/static/css",
		"src/static/assets",
	}
	for _, dir := range folderList {
		config.CheckCreateDir(projectpath + dir)
	}
}

func WebAppSelectTemplate(query, dbCaller, returnType, scanValue, pageName, ModelName string) string {
	return fmt.Sprintf(`func (w http.ResponseWriter, r *http.Request){
	%s
	%s
	rows, err := db.Query("%s")
	if err != nil {
		http.Error(w, "Erreur BDD: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var returnValues []returnType
	for rows.Next() {
		var returnValue returnType
		if err := rows.Scan(%s); err != nil {
			continue
		}
		returnValues = append(returnValues, returnValue)
	}

	renderTemplate(w, "%s.html", map[string]interface{}{
		"%s": returnValues,
	})
}`, returnType, dbCaller, query, scanValue, pageName, ModelName)
}

func WebAppPostViewtemplate(pageName string) string {
	return fmt.Sprintf(`func (w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "%s.html", map[string]interface{}{
		"Title": "Create",
	})
}`, pageName)
}

func WebAppPostActionTemplate(dbCaller, redirectUri, paramsExtraction, paramsChecker, query, paramsExec string) string {
	return fmt.Sprintf(`func (w http.ResponseWriter, r *http.Request) {
	// Parsing de la Form Data
	%s
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		r.ParseForm()
	}

	%s

	if %s {
		http.Error(w, "Les champs 'name' et 'email' sont requis", http.StatusBadRequest)
		return
	}

	// Insertion PostgreSQL via driver pq
	query := "%s"
	_, err := db.Exec(query, %s)
	if err != nil {
		http.Error(w, "Erreur lors de la création: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirection Post-Redirect-Get vers la liste
	http.Redirect(w, r, "%s", http.StatusSeeOther)
}`, dbCaller, paramsExtraction, paramsChecker, query, paramsExec, redirectUri)
}

func WebAppEditTemplate(dbCaller, params, returnType, query, scanValue, pageName, ModelName string) string {
	return fmt.Sprintf(`func HandleUserEdit(w http.ResponseWriter, r *http.Request) {
	%s
	%s := r.PathValue("%s")
	%s
	var returnValue ReturnType
	query := "%s"
	err := db.QueryRow(query, %s).Scan(%s)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
		return
	}

	renderTemplate(w, "%s.html", map[string]interface{}{
		"%s": returnValue,
	})
}`, dbCaller, params, params, returnType, query, params, scanValue, pageName, ModelName)
}

func WebAppEditActionTemplate(dbCaller, params, contentExtraction, query, queryValue, redirectUri string) string {
	return fmt.Sprintf(`func (w http.ResponseWriter, r *http.Request) {
	%s
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		r.ParseForm()
	}

	%s := r.PathValue("%s")
	%s

	query := "%s"
	_, err := db.Exec(query, %s, %s)
	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirection vers la page de détails de l'utilisateur
	http.Redirect(w, r, "%s"+%s, http.StatusSeeOther)
}`, dbCaller, params, params, contentExtraction, query, queryValue, params, redirectUri, params)
}

func WebAppDeleteActionTemplate(dbCaller, params, query, redirectUri string) string {
	return fmt.Sprintf(`func (w http.ResponseWriter, r *http.Request) {
	%s
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		r.ParseForm()
	}

	%s := r.PathValue("%s")

	query := "%s"
	_, err := db.Exec(query, %s)
	if err != nil {
		http.Error(w, "Erreur suppression: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "%s", http.StatusSeeOther)
}`, dbCaller, params, params, query, params, redirectUri)
}

func DBCallerTemplate() string {
	return `
db := config.DB
defer db.Close()
	`
}

func WebAppSelectByParamsTemplate(query, dbCaller, returnType, scanValue, pageName, uriParams, ModelName string) string {
	return fmt.Sprintf(`
func (w http.ResponseWriter, r *http.Request)  {
	%s
	%s := r.PathValue("%s")
	%s
	var returnValue returnType
	query := "%s"
	err := db.QueryRow(query, %s).Scan(%s)
	if err != nil {
		http.Error(w, "model introuvable", http.StatusNotFound)
		return
	}

	renderTemplate(w, "%s.html", map[string]interface{}{
		"%s": returnValue,
		"%s": %s,
	})
}`, goapimaker.DBCallerTemplateWebAPp(), uriParams, uriParams, returnType, query, uriParams, scanValue, pageName, ModelName, utils.ToUpperFirstLetter(uriParams), uriParams)
}

func (wap *webAppMaker) CreateModelFile() {
	if wap.Techno == "go" {
		modelMaker := GoApiMaker{}
		modelMaker.modelAPIExporter(wap.WebApp.bdd.models, wap.ProjectName)
	}
}

func (wap *webAppMaker) CreateConfigFile() {
	if wap.Techno == "go" {
		wap.configAPIExporter2(wap.ProjectName, wap.BDD.models)
	}
}

func (mgr *webAppMaker) configAPIExporter(projectName string) {
	filePath := fmt.Sprintf("%s/%s/web-app/src/config/db.go", config.PROJECT_DIR, projectName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating config file %s : %v\n", filePath, err)
		return
	}
	defer file.Close()

	var sb strings.Builder
	sb.WriteString("package config\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"database/sql\"\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"log\"\n")
	sb.WriteString("\t_ \"github.com/lib/pq\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString("var DB *sql.DB\n\n")
	sb.WriteString("var connStr = \"user=postgres password=root dbname=postgres sslmode=disable host=localhost port=5432\"\n")
	sb.WriteString(`func ConnectDB() *sql.DB {
	var err error
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return DB
}
`)
	sb.WriteString("// InitDB initialise la connexion à la base de données PostgreSQL\n")
	sb.WriteString("func InitDB() {\n")
	sb.WriteString("\t// Modifiez cette chaîne de connexion selon votre environnement\n")

	sb.WriteString("\tvar err error\n")
	sb.WriteString("\tDB, err = sql.Open(\"postgres\", connStr)\n")
	sb.WriteString("\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\n")
	sb.WriteString("\tif err = DB.Ping(); err != nil {\n\t\tlog.Fatal(err)\n\t}\n\n")
	sb.WriteString("\tfmt.Println(\"Successfully connected to database\")\n")
	sb.WriteString("}\n")

	file.WriteString(sb.String())
}

func (mgr *webAppMaker) configAPIExporter2(projectName string, model []Model) {
	filePath := fmt.Sprintf("%s/%s/web-app/src/config/db.go", config.PROJECT_DIR, projectName)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating config file %s : %v\n", filePath, err)
		return
	}
	defer file.Close()

	var sb strings.Builder
	var sqlgen Sqlgenerator
	sb.WriteString("package config\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"database/sql\"\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"log\"\n")
	sb.WriteString("\t_ \"github.com/lib/pq\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString("var DB *sql.DB\n\n")
	fmt.Fprintf(&sb, "var connStr = \"user=postgres password=root dbname=%s sslmode=disable host=localhost port=5432\"\n", strings.ReplaceAll(strings.ToLower(projectName), " ", "_"))
	sb.WriteString(`func ConnectDB() *sql.DB {
	var err error
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return DB
}
`)
	sb.WriteString("// InitDB initialise la connexion à la base de données PostgreSQL\n")
	sb.WriteString("func InitDB() {\n")
	sb.WriteString("\t// Modifiez cette chaîne de connexion selon votre environnement\n")

	sb.WriteString("\tvar err error\n")
	sb.WriteString("\tDB, err = sql.Open(\"postgres\", connStr)\n")
	sb.WriteString("\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\n")
	sb.WriteString("\tif err = DB.Ping(); err != nil {\n\t\tlog.Fatal(err)\n\t}\n\n")
	sb.WriteString(fmt.Sprintf("\tquery := `%s`\n", sqlgen.GenerateDBAndTableInit(model, projectName)))
	sb.WriteString("\t_, err = DB.Exec(query)\n\tif err != nil {\n\t\tfmt.Errorf(\"échec lors de l'initialisation des tables: %w\", err)\n\t}\n")
	sb.WriteString("\tfmt.Println(\"Successfully connected to database\")\n")
	sb.WriteString("}\n")

	file.WriteString(sb.String())
}

func (wap *webAppMaker) HandleURIParamsSyntaxeForGo(uri string) string {
	uriTab := strings.Split(uri, "/")
	realUri := []string{}
	for _, val := range uriTab {
		res := val
		if strings.Contains(val, ":") {
			res = fmt.Sprintf("{%s}", strings.ReplaceAll(val, ":", ""))
		}
		realUri = append(realUri, res)
	}
	return strings.Join(realUri, "/")
}

func (wap *webAppMaker) WriteBodyType(endpoint Endpoint) string {
	var strBuilder strings.Builder
	fmt.Fprint(&strBuilder, "type bodyType struct{\n")
	for _, val := range endpoint.model {
		for _, field := range val.attributs {
			fmt.Fprintf(&strBuilder, "%s %s `json:\"%s\"`", utils.ToUpperFirstLetter(field.nom), field.type_champs, field.nom)
		}
	}
	fmt.Fprint(&strBuilder, "}\n")
	return strBuilder.String()
}
func (wap *webAppMaker) WriteReturnType(endpoint Endpoint) string {
	var strBuilder strings.Builder
	fmt.Fprint(&strBuilder, "type returnType struct{\n")
	for _, val := range endpoint.returnContent {
		for _, field := range val.attributs {
			fmt.Fprintf(&strBuilder, "%s %s `json:\"%s\"`\n", utils.ToUpperFirstLetter(field.nom), field.type_champs, field.nom)
		}
	}
	fmt.Fprint(&strBuilder, "}\n")
	return strBuilder.String()
}

func (wap *webAppMaker) WriteParamsGetter(endpoint Endpoint) string {
	var strBuilder strings.Builder
	for _, val := range endpoint.params {
		fmt.Fprintf(&strBuilder, "%s := r.PathValue(\"%s\")\n", val, val)
	}
	return strBuilder.String()
}

func (wap *webAppMaker) WriteScanValue(endpoint Endpoint) string {
	var strBuilder []string
	for _, val := range endpoint.returnContent {
		for _, mod := range val.attributs {
			strBuilder = append(strBuilder, fmt.Sprintf("&returnValue.%s", utils.ToUpperFirstLetter(mod.nom)))
		}
	}
	return strings.Join(strBuilder, ", ")
}

func (wap *webAppMaker) WriteContentExtraction(endpoint Endpoint) string {
	var str []string
	for _, val := range endpoint.model[0].attributs {
		str = append(str, fmt.Sprintf("%s := r.FormValue(\"%s\")\n", val.nom, val.nom))
	}
	return strings.Join(str, "")
}

func (wap *webAppMaker) WriteParamsChecker(endpoint Endpoint) string {
	var str []string
	for _, val := range endpoint.model[0].attributs {
		str = append(str, fmt.Sprintf("%s == \"\"", val.nom))
	}
	return strings.Join(str, " || ")
}

func (wap *webAppMaker) WriteControllerForObjectOrArrayReturn(endpoint Endpoint) string {
	var strBuilder strings.Builder
	var method string
	if endpoint.method == "GET" {
		method = "GET"
	} else {
		method = "POST"
	}

	fmt.Fprintf(&strBuilder, "\nmux.HandleFunc(\"%s %s\",", method, wap.HandleURIParamsSyntaxeForGo(endpoint.uri))
	switch endpoint.method {
	case "GET":
		if len(endpoint.model) > 0 {
			if len(endpoint.params) > 0 {
				var attrTab []string
				for _, val := range endpoint.returnContent[0].attributs {
					attrTab = append(attrTab, val.nom)
				}

				fmt.Fprint(&strBuilder, WebAppSelectByParamsTemplate(goapimaker.SelectByWithAttr(endpoint.model[0].nom, strings.Join(attrTab, ", "), endpoint.params[0]), goapimaker.DbCallerPGWebApp(), wap.WriteReturnType(endpoint), wap.WriteScanValue(endpoint), strings.ToLower(strings.ReplaceAll(endpoint.returnPage, " ", "_")), endpoint.params[0], endpoint.model[0].nom))

			} else {
				var attrTab []string
				for _, val := range endpoint.returnContent[0].attributs {
					attrTab = append(attrTab, val.nom)
				}
				fmt.Fprint(&strBuilder, WebAppSelectTemplate(goapimaker.SelectWithAttr(endpoint.model[0].nom, strings.Join(attrTab, ", ")), goapimaker.DbCallerPGWebApp(), wap.WriteReturnType(endpoint), wap.WriteScanValue(endpoint), strings.ToLower(strings.ReplaceAll(endpoint.returnPage, " ", "_")), endpoint.model[0].nom))

			}
		} else {
			fmt.Fprint(&strBuilder, WebAppPostViewtemplate(strings.ToLower(strings.ReplaceAll(endpoint.returnPage, " ", "_"))))
		}
	case "POST":
		var attr []string
		for _, val := range endpoint.model[0].attributs {
			attr = append(attr, val.nom)
		}
		fmt.Fprint(&strBuilder, WebAppPostActionTemplate(goapimaker.DbCallerPGWebApp(), endpoint.redirectUri, wap.WriteContentExtraction(endpoint), wap.WriteParamsChecker(endpoint), goapimaker.Insert(endpoint.model[0].nom, attr), strings.Join(attr, ", ")))
	case "PUT":
		var attr []string
		var attrForQuery []string
		for i, val := range endpoint.model[0].attributs {
			attr = append(attr, fmt.Sprintf("%s = $%d", val.nom, i+1))
			attrForQuery = append(attrForQuery, val.nom)
		}
		fmt.Fprint(&strBuilder, WebAppEditActionTemplate(goapimaker.DbCallerPGWebApp(), endpoint.params[0], wap.WriteContentExtraction(endpoint), goapimaker.Update(endpoint.model[0].nom, attr, endpoint.params[0]), strings.Join(attrForQuery, ", "), endpoint.redirectUri))
	case "DELETE":
		fmt.Fprint(&strBuilder, WebAppDeleteActionTemplate(goapimaker.DbCallerPGWebApp(), endpoint.params[0], goapimaker.Delete(endpoint.model[0].nom, endpoint.params[0]), endpoint.redirectUri))
	}
	strBuilder.WriteString(")")
	return strBuilder.String()
}

func (wap *webAppMaker) PageRenderer() string {
	return `
func renderTemplate(w http.ResponseWriter, pageName string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := template.ParseFiles("src/views/" + pageName)
	if err != nil {
		http.Error(w, "Erreur de chargement du template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur de rendu du template: "+err.Error(), http.StatusInternalServerError)
	}
}	
	`
}
func (mgr *webAppMaker) mainExporter() {
	filePath := fmt.Sprintf("%s/%s/web-app/main.go", config.PROJECT_DIR, mgr.ProjectName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating main file %s : %v\n", filePath, err)
		return
	}
	defer file.Close()

	var sb strings.Builder
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"log\"\n")
	sb.WriteString("\t\"net/http\"\n")
	sb.WriteString("\t\"")
	sb.WriteString(strings.ReplaceAll(mgr.ProjectName, " ", "_"))
	sb.WriteString("/src/config\"\n")
	sb.WriteString("\t\"")
	sb.WriteString(strings.ReplaceAll(mgr.ProjectName, " ", "_"))
	sb.WriteString("/src/controller\"\n")
	sb.WriteString("\t)\n")
	// sb.WriteString("\t\"")
	// sb.WriteString(strings.ReplaceAlreturnValuesl(mgr.ProjectName, " ", "_"))
	// sb.WriteString("/src/middlewares\"\n")
	// sb.WriteString(")\n\n")

	sb.WriteString("func main() {\n")
	sb.WriteString("\t// 1. Initialisation de la base de données\n")
	sb.WriteString("\tconfig.InitDB()\n\n")
	sb.WriteString("\t// 2. Initialisation du Router (Mux)\n")
	sb.WriteString("\tmux := http.NewServeMux()\n\n")
	sb.WriteString("\tcontroller.RegisterRoutes(mux)\n\n")
	sb.WriteString("\tstaticServer := http.FileServer(http.Dir(\"src/static\"))\n\n")
	sb.WriteString("\tmux.Handle(\"GET /static/\", http.StripPrefix(\"/static/\", staticServer))\n")
	sb.WriteString("\t// 5. Lancement du serveur\n")
	sb.WriteString("\tport := \":8080\"\n")
	sb.WriteString("\tfmt.Printf(\"🚀 TurboStack API running on http://localhost%s\\n\", port)\n")
	sb.WriteString("\tlog.Fatal(http.ListenAndServe(port, mux))\n")
	sb.WriteString("}\n")

	file.WriteString(sb.String())
}

func (mgr *webAppMaker) writeModSumFile() {
	projectName := strings.ReplaceAll(mgr.ProjectName, " ", "_")
	modfilepath := fmt.Sprintf("%s/%s/web-app/go.mod", config.PROJECT_DIR, projectName)
	sumfilepath := fmt.Sprintf("%s/%s/web-app/go.sum", config.PROJECT_DIR, projectName)
	modfile, _ := os.Create(modfilepath)
	modfile.WriteString(goapimaker.WriteGoMod(strings.ReplaceAll(projectName, " ", "_")))
	sumfile, _ := os.Create(sumfilepath)
	sumfile.WriteString(goapimaker.WriteSum())
}

func (wap *webAppMaker) CreateControllerFile() {
	if wap.Techno == "go" {
		var strBuilder strings.Builder
		filePath := config.PROJECT_DIR + "/" + wap.ProjectName + "/web-app/src/controller/controller.go"
		controllerFile, _ := os.Create(filePath)

		strBuilder.WriteString("package controller\n")
		fmt.Fprintf(&strBuilder, "import (\n\"net/http\"\n\"text/template\"\n\"%s/src/config\"\n)", strings.ReplaceAll(wap.ProjectName, " ", "_"))

		strBuilder.WriteString(wap.PageRenderer())

		strBuilder.WriteString("func RegisterRoutes(mux *http.ServeMux){\n")
		for _, val := range wap.Api.endpoints {
			strBuilder.WriteString(wap.WriteControllerForObjectOrArrayReturn(val))
		}
		strBuilder.WriteString("}\n")
		controllerFile.WriteString(strBuilder.String())
	}
}

func (wap *webAppMaker) RenderBlocksToHTML(blocks []pageContent, projectName string, pageName string) string {
	cssPath := fmt.Sprintf("%s/%s/web-app/src/static/css/%s.css", config.PROJECT_DIR, projectName, pageName)
	cssFile, _ := os.OpenFile(cssPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)

	var sb strings.Builder

	var desktopSb strings.Builder
	desktopSb.Grow(4096)
	var tabletSb strings.Builder
	tabletSb.Grow(4096)
	defer desktopSb.Reset()
	defer tabletSb.Reset()
	for _, b := range blocks {
		block := b

		tag := fmt.Sprintf("%v", block.tag)
		content := fmt.Sprintf("%v", block.content)
		className := fmt.Sprintf("%v", block.className)
		style := fmt.Sprintf("%v", block.styles)
		id := fmt.Sprintf("%v", block.id)
		inputType := fmt.Sprintf("%v", block.inputType)
		placeholder := fmt.Sprintf("%v", block.placeholder)

		if tag == "img" {
			fmt.Fprintf(&sb, "<img src=\"%s\" class=\"%s\" data-block-id=\"%s\" />", content, className, id)
			continue
		} else if tag == "input" {
			fmt.Fprintf(&sb, "<input type=\"%s\" class=\"%s\" data-block-id=\"%s\" placeholder=\"%s\" name=\"%s\" value=\"%s\"/>", inputType, className, id, placeholder, block.name, block.value)
			continue
		} else if tag == "form" {
			fmt.Fprintf(&sb, "<%s class=\"%s\" data-block-id=\"%s\" action=\"%s\" method=\"POST\" >", tag, className, id, block.formAction)
		} else if tag == "a" {
			fmt.Fprintf(&sb, "<%s class=\"%s\" data-block-id=\"%s\" href=\"%s\"> ", tag, className, id, block.href)
		} else {
			fmt.Fprintf(&sb, "<%s class=\"%s\" data-block-id=\"%s\" > ", tag, className, id)
		}

		var cssVal map[string]map[string]string
		if style != "" {
			er := json.Unmarshal([]byte(style), &cssVal)
			if er != nil {
				fmt.Println("error", er)
				continue
			}
		}

		desktop := cssVal["desktop"]
		tablet := cssVal["tablet"]
		mobile := cssVal["mobile"]
		if len(tablet) > 0 {
			fmt.Fprintf(&tabletSb, "[data-block-id=\"%s\"]{\n", id)
			for key, val := range tablet {
				checkValueSb(&tabletSb, key, val)
			}
			fmt.Fprint(&tabletSb, "}\n")
		}

		if len(mobile) > 0 {
			fmt.Fprintf(&desktopSb, "[data-block-id=\"%s\"]{\n", id)
			for key, val := range mobile {
				checkValueSb(&desktopSb, key, val)
				// fmt.Fprint(&desktopSb, key, val)
			}
			fmt.Fprint(&desktopSb, "}\n")
		}

		if len(desktop) > 0 {
			fmt.Fprintf(cssFile, "[data-block-id=\"%s\"]{\n", id)
			for key, val := range desktop {
				checkValueFile(cssFile, key, val)
			}
			fmt.Fprint(cssFile, "}\n")
		}
		sb.WriteString(content)

		if len(block.children) > 0 {
			sb.WriteString(wap.RenderBlocksToHTML(block.children, projectName, pageName))
		}

		fmt.Fprintf(&sb, "</%s>", tag)
	}
	if desktopSb.Len() > 0 {
		cssFile.WriteString("@media (max-width: 375px) {\n")
		cssFile.WriteString(desktopSb.String())
		cssFile.WriteString("}\n")
	}
	if tabletSb.Len() > 0 {
		cssFile.WriteString("@media (min-width: 376px) and (max-width: 1024px) {\n")
		cssFile.WriteString(tabletSb.String())
		cssFile.WriteString("}\n")
	}
	return sb.String()
}

func (wap *webAppMaker) GenerateView() {
	projectName := wap.ProjectName
	site := wap.WebApp.pages
	cssPath := fmt.Sprintf("%s/%s/web-app/src/static/css/style.css", config.PROJECT_DIR, projectName)
	cssFile, _ := os.Create(cssPath)

	defer cssFile.Close()
	for _, page := range site {
		pageName := strings.ToLower(strings.ReplaceAll(page.GetNom(), " ", "_"))
		filePath := fmt.Sprintf("%s/%s/web-app/src/views/%s.html", config.PROJECT_DIR, projectName, pageName)
		cssPath := fmt.Sprintf("%s/%s/web-app/src/views/css/global_%s.css", config.PROJECT_DIR, projectName, pageName)
		cssFile, _ := os.OpenFile(cssPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		styleWriter(page, cssFile)
		file, err := os.Create(filePath)
		if err != nil {
			continue
		}

		var sb strings.Builder
		sb.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
		sb.WriteString("\t<meta charset=\"UTF-8\">\n")
		fmt.Fprintf(&sb, "\t<title>%s</title>\n", page.GetNom())
		fmt.Fprint(&sb, "\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\n")
		fmt.Fprintf(&sb, "\t<link rel=\"stylesheet\" href=\"/static/css/%s.css\">\n", pageName)
		fmt.Fprintf(&sb, "\t<link rel=\"stylesheet\" href=\"/static/css/global_%s.css\">\n", pageName)
		sb.WriteString("</head>\n<body>\n")

		// Conversion du contenu JSON en HTML
		content := page.GetContent()

		sb.WriteString(wap.RenderBlocksToHTML(content, projectName, pageName))

		sb.WriteString("\n</body>\n</html>")

		file.WriteString(sb.String())
		file.Close()
	}
}

func (wap *webAppMaker) WebAppGenerator() {
	wap.SetupArch()
	wap.writeModSumFile()
	wap.CreateConfigFile()
	// wap.CreateModelFile()
	wap.CreateControllerFile()
	wap.GenerateView()
	wap.mainExporter()
}
