package entity

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/MaminirinaEdwino/turbostackapi/src/config"
	"github.com/MaminirinaEdwino/turbostackapi/src/goapimaker"
)

type GoApiMaker struct{}

func (mgr *GoApiMaker) setupFileArch(name string) {
	projectPath := fmt.Sprintf("%s/api/", name)
	for _, val := range goapimaker.DirList {
		config.CheckCreateDir(projectPath + val)
	}
}

func (api *GoApiMaker) modelAPIExporter(models []Model, projectName string) {
	for _, model := range models {
		mName := model.GetNom()
		if mName == "" {
			continue
		}
		n := strings.ReplaceAll(mName, " ", "_")
		filePath := fmt.Sprintf("%s/%s/api/src/models/%s.go", config.PROJECT_DIR, projectName, strings.ToLower(n))
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating model file %s : %v\n", filePath, err)
			continue
		}

		var sb strings.Builder
		sb.WriteString("package models\n\n")
		structName := strings.ToUpper(mName[:1]) + mName[1:]
		fmt.Fprintf(&sb, "type %s struct {\n", structName)

		for _, field := range model.GetAttributs() {
			fieldName := strings.ToUpper(field.GetNom()[:1]) + field.GetNom()[1:]
			goType := goapimaker.MapToGoType(field.GetType())
			fmt.Fprintf(&sb, "\t%s %s `json:\"%s\"`\n", fieldName, goType, strings.ToLower(field.GetNom()))
		}
		sb.WriteString("}\n")
		file.WriteString(sb.String())
		file.Close()
	}
}
func (mgr *GoApiMaker) routesAPIExporter(endpoints []Endpoint, projectName string) {
	filePath := fmt.Sprintf("%s/%s/api/src/routes/routes.go", config.PROJECT_DIR, projectName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating routes file %s : %v\n", filePath, err)
		return
	}
	defer file.Close()

	var sb strings.Builder
	sb.WriteString("package routes\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"net/http\"\n")
	sb.WriteString("\t\"")
	sb.WriteString(strings.ReplaceAll(projectName, " ", "_"))
	sb.WriteString("/src/controllers\"\n")
	sb.WriteString(")\n\n")

	// Groupement des endpoints par nom de modèle
	groups := make(map[string][]Endpoint)
	for _, ep := range endpoints {
		mName := "General"
		if len(ep.GetModel()) > 0 {
			mName = ep.GetModel()[0].GetNom()
		}
		groups[mName] = append(groups[mName], ep)
	}

	for groupName, eps := range groups {
		cleanGroupName := strings.ReplaceAll(groupName, " ", "_")
		groupFuncName := strings.ToUpper(cleanGroupName[:1]) + cleanGroupName[1:]
		fmt.Fprintf(&sb, "func RegisterRoutes%s(mux *http.ServeMux) {\n", groupFuncName)
		for _, ep := range eps {
			method := ep.GetMethod()
			uri := ep.GetUri()
			eName := ep.GetNom()
			funcName := strings.ToUpper(eName[:1]) + strings.ReplaceAll(eName[1:], " ", "_")
			fmt.Fprintf(&sb, "\tmux.HandleFunc(\"%s %s\", controllers.%s)\n", method, uri, funcName)
		}
		sb.WriteString("}\n\n")
	}
	file.WriteString(sb.String())
}

func (mgr *GoApiMaker) controllerAPIExporter(endpoints []Endpoint, projectName string) {
	for _, ep := range endpoints {
		eName := ep.GetNom()
		if eName == "" {
			continue
		}

		filePath := fmt.Sprintf("%s/%s/api/src/controllers/%s.go", config.PROJECT_DIR, projectName, strings.ToLower(strings.ReplaceAll(eName, " ", "_")))
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating controller %s : %v\n", filePath, err)
			continue
		}

		var sb strings.Builder
		goapimaker.ControllerImportWriter(&sb, projectName, ep.method)

		funcName := strings.ToUpper(eName[:1]) + eName[1:]

		// Identification du modèle associé
		var activeModel Model
		if len(ep.GetModel()) > 0 {
			activeModel = ep.GetModel()[0]
		}
		modelName := activeModel.GetNom()

		fmt.Fprintf(&sb, "// %s handles the %s request for %s\n", funcName, ep.GetMethod(), ep.GetUri())
		// fmt.Fprintf(&sb, "func %s(w http.ResponseWriter, r *http.Request) {\n", strings.ReplaceAll(funcName, " ", "_"))
		var tmpBuffer bytes.Buffer
		if modelName != "" {
			structName := strings.ToUpper(modelName[:1]) + modelName[1:]
			tableName := strings.ToLower(modelName)

			var columns []string
			var scanTargets []string
			for _, attr := range activeModel.GetAttributs() {
				columns = append(columns, strings.ToLower(attr.GetNom()))
				scanTargets = append(scanTargets, "&tmp."+strings.ToUpper(attr.GetNom()[:1])+attr.GetNom()[1:])
			}

			switch ep.method {
			case "DELETE":
				fmt.Fprint(&sb, goapimaker.DeleteHandler(strings.ReplaceAll(funcName, " ", "_"), "pg", ep.params[0], tableName))
			case "PUT":
				fmt.Fprint(&sb, goapimaker.PutHandler(ep.model[0].nom, strings.ReplaceAll(funcName, " ", "_"), "pg", columns, strings.ReplaceAll(strings.Join(scanTargets, ", "), "tmp", "res"), ep.params[0]))
			case "GET":
				if len(ep.params) > 0 {
					data := struct {
						StructName      string
						DbCallerHandler string
						Query           string
						Params          string
						ScanParams      string
						ResponseWriter  string
						FuncName        string
					}{
						StructName:      structName,
						DbCallerHandler: goapimaker.DBCallerHandler("pg"),
						Query:           goapimaker.SelectBy(tableName, ep.params[0]),
						ScanParams:      strings.ReplaceAll(strings.Join(scanTargets, ", "), "tmp", "res"),
						ResponseWriter:  goapimaker.WriteResponseWriter(),
						Params:          ep.params[0],
						FuncName:        strings.ReplaceAll(funcName, " ", "_"),
					}
					err := goapimaker.SelectBytemplate().Execute(&tmpBuffer, data)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Fprintln(&sb, tmpBuffer.String())
				} else {

					fmt.Fprint(&sb, goapimaker.SelectTemplate(strings.ReplaceAll(funcName, " ", "_"), goapimaker.DBCallerHandler("pg"), structName, goapimaker.Select(tableName), strings.Join(scanTargets, ", "), goapimaker.WriteResponseWriter()))
				}
			case "POST":
				var attrs []string
				for _, val := range ep.model[0].attributs {
					attrs = append(attrs, fmt.Sprintf("%s", val.nom))
				}
				fmt.Fprint(&sb, goapimaker.InsertHandler(strings.ReplaceAll(funcName, " ", "_"), "pg", tableName, attrs))
			}

		} else {
			fmt.Fprintf(&sb, "\tjson.NewEncoder(w).Encode(map[string]string{\"message\": \"%s endpoint reached\"})\n", eName)
		}

		file.WriteString(sb.String())
		file.Close()
	}
}

func (mgr *GoApiMaker) middlewareAPIExporter(projectName string) {
	goapimaker.MiddlewareWriter(projectName)
}

func (mgr *GoApiMaker) mainAPIExporter(endpoints []Endpoint, projectName string) {
	filePath := fmt.Sprintf("%s/%s/api/main.go", config.PROJECT_DIR, projectName)
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
	sb.WriteString(strings.ReplaceAll(projectName, " ", "_"))
	sb.WriteString("/src/config\"\n")
	sb.WriteString("\t\"")
	sb.WriteString(strings.ReplaceAll(projectName, " ", "_"))
	sb.WriteString("/src/routes\"\n")
	sb.WriteString("\t\"")
	sb.WriteString(strings.ReplaceAll(projectName, " ", "_"))
	sb.WriteString("/src/middlewares\"\n")
	sb.WriteString(")\n\n")

	sb.WriteString("func main() {\n")
	sb.WriteString("\t// 1. Initialisation de la base de données\n")
	sb.WriteString("\tconfig.InitDB()\n\n")
	sb.WriteString("\t// 2. Initialisation du Router (Mux)\n")
	sb.WriteString("\tmux := http.NewServeMux()\n\n")

	// Récupération des groupes de routes (logique identique à routesAPIExporter)
	groups := make(map[string]bool)
	for _, ep := range endpoints {
		mName := "General"
		if len(ep.GetModel()) > 0 {
			mName = ep.GetModel()[0].GetNom()
		}
		groups[mName] = true
	}

	sb.WriteString("\t// 3. Enregistrement des routes par modèle\n")
	for groupName := range groups {
		cleanGroupName := strings.ReplaceAll(groupName, " ", "_")
		groupFuncName := strings.ToUpper(cleanGroupName[:1]) + cleanGroupName[1:]
		fmt.Fprintf(&sb, "\troutes.RegisterRoutes%s(mux)\n", groupFuncName)
	}

	sb.WriteString("\n\t// 4. Application du middleware CORS\n")
	sb.WriteString("\thandler := middlewares.CorsMiddleware(mux)\n\n")

	sb.WriteString("\t// 5. Lancement du serveur\n")
	sb.WriteString("\tport := \":8080\"\n")
	sb.WriteString("\tfmt.Printf(\"🚀 TurboStack API running on http://localhost%s\\n\", port)\n")
	sb.WriteString("\tlog.Fatal(http.ListenAndServe(port, handler))\n")
	sb.WriteString("}\n")

	file.WriteString(sb.String())
}

func (mgr *GoApiMaker) writeModSumFile(projectName string) {
	modfilepath := fmt.Sprintf("%s/%s/api/go.mod", config.PROJECT_DIR, projectName)
	sumfilepath := fmt.Sprintf("%s/%s/api/go.sum", config.PROJECT_DIR, projectName)
	modfile, _ := os.Create(modfilepath)
	modfile.WriteString(goapimaker.WriteGoMod(strings.ReplaceAll(projectName, " ", "_")))
	sumfile, _ := os.Create(sumfilepath)
	sumfile.WriteString(goapimaker.WriteSum())
}

func (mgr *GoApiMaker) configAPIExporter(projectName string, model []Model) {
	filePath := fmt.Sprintf("%s/%s/api/src/config/db.go", config.PROJECT_DIR, projectName)

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
	sb.WriteString("// InitDB initialise la connexion à la base de données PostgreSQL\n")
	sb.WriteString("func InitDB() {\n")
	sb.WriteString("\t// Modifiez cette chaîne de connexion selon votre environnement\n")
	sb.WriteString(fmt.Sprintf("\tconnStr := \"user=postgres password=root dbname=%s sslmode=disable host=localhost port=5432\"\n", strings.ReplaceAll(strings.ToLower(projectName), " ", "_")))
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
