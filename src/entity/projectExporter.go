package entity

import (
	"fmt"
	"os"
	"strings"

	"github.com/MaminirinaEdwino/turbostackapi/src/config"
)

func (mgr *ProjectManager) ExporterDB(Project Project) {
	sqlgenerator := Sqlgenerator{}
	db := Project.GetBDD()
	models := db.GetModels()
	projectName := Project.GetNom()
	sqlgenerator.PostgresSQLExporter(models, projectName)
	sqlgenerator.PostgresCRUDExporter(models, projectName)
	fmt.Printf("Exportation des modèles terminée pour le projet : %s\n", projectName)
}

func (mgr *ProjectManager) ExporterAPI(Project Project) {
	apiMaker := GoApiMaker{}
	projectName := Project.GetNom()
	bdd := Project.GetBDD()
	apiMaker.setupFileArch(projectName)
	models := bdd.GetModels()
	api := Project.GetRestApi()
	apiMaker.modelAPIExporter(models, projectName)
	apiMaker.controllerAPIExporter(api.GetEndpoints(), projectName)
	apiMaker.routesAPIExporter(api.GetEndpoints(), projectName)
	apiMaker.configAPIExporter(projectName, Project.bdd.models)
	apiMaker.middlewareAPIExporter(projectName)
	apiMaker.mainAPIExporter(api.GetEndpoints(), projectName)
	apiMaker.writeModSumFile(projectName)
	fmt.Printf("Exportation de l'API terminée pour le projet : %s\n", projectName)
}

func (mgr *ProjectManager) ExporterStaticSite(Project Project) {
	projectName := Project.GetNom()
	site := Project.GetSiteStatique()
	ssm := Staticsitemaker{}
	ssm.SetupStaticArch(projectName)

	// 1. Génération du CSS global
	cssPath := fmt.Sprintf("%s/%s/static/static/css/style.css", config.PROJECT_DIR, projectName)
	cssFile, _ := os.Create(cssPath)

	defer cssFile.Close()

	// 2. Génération des pages HTML
	for _, page := range site.GetPages() {
		pageName := strings.ToLower(strings.ReplaceAll(page.GetNom(), " ", "_"))
		filePath := fmt.Sprintf("%s/%s/static/views/%s.html", config.PROJECT_DIR, projectName, pageName)
		cssPath := fmt.Sprintf("%s/%s/static/static/css/global_%s.css", config.PROJECT_DIR, projectName, pageName)
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

		sb.WriteString(ssm.RenderBlocksToHTML(content, projectName, pageName))

		sb.WriteString("\n</body>\n</html>")

		file.WriteString(sb.String())
		file.Close()
	}
	techno := "golang"
	if techno == "golang" {
		mainPath := fmt.Sprintf("%s/%s/static/main.go", config.PROJECT_DIR, projectName)
		modPath := fmt.Sprintf("%s/%s/static/go.mod", config.PROJECT_DIR, projectName)
		ssm.SetupGoServerCode(site.GetPages(), mainPath)
		file, _ := os.Create(modPath)
		defer file.Close()
		file.WriteString(ssm.WriteGoMod(projectName))
	}
	fmt.Printf("Exportation du site statique terminée : %s\n", projectName)
}

func (mgr *ProjectManager) ExporterWebApp(Project Project) {
	webAppMaker := webAppMaker{
		ProjectName: Project.GetNom(),
		WebApp:      Project.GetWebApp(),
		Techno:      "go",
		Api:         Project.rest_api,
		BDD:         Project.bdd,
	}
	webAppMaker.WebAppGenerator()
}
