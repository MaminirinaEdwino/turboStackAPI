package api

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/MaminirinaEdwino/turbostackapi/src/config"
	"github.com/MaminirinaEdwino/turbostackapi/src/entity"
	webview "github.com/webview/webview_go"
)

type ProjectService struct {
	Manager *entity.ProjectManager
	WV      webview.WebView
}

func (ps *ProjectService) Bind(w webview.WebView) {
	w.Bind("createProject", ps.CreateProject)
	w.Bind("fetchProjects", ps.FetchProjects)
	w.Bind("fetchByProjectName", ps.FetchProjectByName)
	w.Bind("saveProject", ps.SaveProject)
	w.Bind("saveBdd", ps.SaveProjectBDD)
	w.Bind("exportProject", ps.ExportProject)
	w.Bind("fetchProjectFiles", ps.FeychProjectFiles)
	w.Bind("getFileContent", ps.GetFileContent)
	w.Bind("getFolderForUpload", ps.FetchFolderForUpload)
	w.Bind("getImageAsBase64", ps.GetImageAsBase64)
	w.Bind("uploadAsset", ps.UploadAsset)
	w.Bind("saveScript", ps.SaveScript)
	w.Bind("startProject", ps.StartProject)
	w.Bind("stopProject", ps.Stopproject)
	w.Bind("getStatus", ps.StatusProject)
	w.Bind("loadPageLibrairie", ps.LoadPageLib)
	w.Bind("loadStyleLib", ps.LoadStyleLib)
	w.Bind("loadCompLib", ps.LoadComponentLib)
	w.Bind("getProjectType", ps.GetProjectType)
	w.Bind("deleteProject", ps.DeleteProject)
}

type ProcessManager struct {
	mu        sync.Mutex
	cmd       *exec.Cmd
	isRunning bool
	logs      []string
	pid       int
}

var pm = &ProcessManager{}

func (s *ProjectService) GetProjectType(projectName string) string {
	project := s.FetchProjectByName(projectName)
	return project.Type
}

func (s *ProjectService) DeleteProject(projectName string) {
	project := s.FetchProjectByName(projectName)
	s.Manager.Delete(project.ToModel())
}

func (s *ProjectService) StartProject(projectName string) map[string]interface{} {
	project := s.FetchProjectByName(projectName)
	var projectDir string
	switch project.Type {
	case "static":
		projectDir = fmt.Sprintf("%s/%s/static", config.PROJECT_DIR, projectName)
	case "api":
		projectDir = fmt.Sprintf("%s/%s/api", config.PROJECT_DIR, projectName)
	case "web_app":
		projectDir = fmt.Sprintf("%s/%s/web_app", config.PROJECT_DIR, projectName)
	}

	switch project.Type {
	case "static":
		s.Manager.ExporterStaticSite(project.ToModel())
	case "api":
		s.Manager.ExporterAPI(project.ToModel())
	case "web_app":
		s.Manager.ExporterWebApp(project.ToModel())
	}
	return handleStartProject(projectDir, s.WV)
}

func (s *ProjectService) Stopproject(projectName string) map[string]interface{} {
	return handleStopProject(projectName)
}

func (s *ProjectService) StatusProject() map[string]interface{} {
	return handleGetStatus()
}

func (s *ProjectService) LoadPageLib() []string {
	var res []string
	dirContent, _ := os.ReadDir(config.PAGE_LIB_DIR)
	for _, val := range dirContent {
		fileContent, _ := os.ReadFile(config.PAGE_LIB_DIR + "/" + val.Name())
		res = append(res, string(fileContent))
	}
	return res
}

func (s *ProjectService) LoadComponentLib() []string {
	var res []string
	dirContent, _ := os.ReadDir(config.COMPONENT_LIB_DIR)
	for _, val := range dirContent {
		fileContent, _ := os.ReadFile(config.COMPONENT_LIB_DIR + "/" + val.Name())
		res = append(res, string(fileContent))
	}
	return res
}

func (s *ProjectService) LoadStyleLib() []string {
	var res []string
	dirContent, _ := os.ReadDir(config.STYLE_LIB_DIR)
	for _, val := range dirContent {
		fileContent, _ := os.ReadFile(config.STYLE_LIB_DIR + "/" + val.Name())
		res = append(res, string(fileContent))
	}
	return res
}

// 1. Démarrer le projet Go
func handleStartProject(projectDir string, w webview.WebView) map[string]interface{} {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.isRunning {
		return map[string]interface{}{
			"isRunning": pm.isRunning,
			"logs":      pm.logs,
			"pid":       pm.pid,
		}
	}

	// Commande pour exécuter votre projet Go (ex: 'go run main.go' dans le dossier cible)
	pm.cmd = exec.Command("go", "run", ".")
	pm.cmd.Dir = projectDir
	// Optionnel : Définir le répertoire de travail du projet
	// pm.cmd.Dir = "./generated_project"

	// Capturer les sorties standard et d'erreur pour le dashboard
	stdout, _ := pm.cmd.StdoutPipe()
	stderr, _ := pm.cmd.StderrPipe()

	if err := pm.cmd.Start(); err != nil {

		return map[string]interface{}{
			"isRunning": pm.isRunning,
			"logs":      pm.logs,
			"pid":       pm.pid,
		}
	}

	pm.isRunning = true
	pm.pid = pm.cmd.Process.Pid
	pm.logs = []string{"Projet started"}

	// Goroutine pour lire les logs en continu sans bloquer l'application
	go func() {
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for scanner.Scan() {
			pm.mu.Lock()
			pm.logs = append(pm.logs, scanner.Text())
			data, _ := json.Marshal(map[string]interface{}{
				"isRunning": pm.isRunning,
				"logs":      pm.logs,
				"pid":       pm.pid,
			})
			Dispatch(w, "get-status-event", string(data))
			pm.mu.Unlock()
		}

		// Attendre la fin du processus
		_ = pm.cmd.Wait()
		pm.mu.Lock()
		pm.isRunning = false
		pm.logs = append(pm.logs, "The  project has stopped")
		pm.mu.Unlock()
	}()

	// return fmt.Sprintf(`{"isRunning": %t, "logs": %q}`, pm.isRunning, pm.logs)
	return map[string]interface{}{
		"isRunning": pm.isRunning,
		"logs":      pm.logs,
		"pid":       pm.pid,
	}
}

func PkillByName(processName string) error {
	// 'pkill' accepte le nom exact ou partiel de l'exécutable
	cmd := exec.Command("pkill", "-9", processName)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("échec de pkill sur '%s' (%v): %s", processName, err, string(output))
	}
	return nil
}

// 2. Arrêter le projet Go
func handleStopProject(projectName string) map[string]interface{} {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if !pm.isRunning || pm.cmd == nil || pm.cmd.Process == nil {

		return map[string]interface{}{
			"isRunning": pm.isRunning,
			"logs":      pm.logs,
			"pid":       pm.pid,
		}
	}

	// err := syscall.Kill(pm.cmd.Process.Pid, syscall.SIGTERM)
	err := PkillByName(strings.ReplaceAll(projectName, " ", "_"))
	if err != nil {
		pm.logs = append(pm.logs, err.Error())
		return map[string]interface{}{
			"isRunning": pm.isRunning,
			"logs":      pm.logs,
			"pid":       pm.pid,
		}
	}

	pm.isRunning = false
	pm.logs = append(pm.logs, "Server Stopped")
	return map[string]interface{}{
		"isRunning": pm.isRunning,
		"logs":      pm.logs,
		"pid":       pm.pid,
	}
}

// 3. Récupérer le statut et les logs
func handleGetStatus() map[string]interface{} {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	return map[string]interface{}{
		"isRunning": pm.isRunning,
		"logs":      pm.logs,
		"pid":       pm.pid,
	}
}

func (s *ProjectService) UploadAsset(projectName, fileName, base64file string) string {
	project := s.FetchProjectByName(projectName)
	project.Assets = append(project.Assets, entity.AssetJSON{
		FileName:    fileName,
		Base64Image: base64file,
	})
	s.Manager.SaveProject(project)
	return "success"
}

func (s *ProjectService) SaveScript(projectName, scriptName, script string) string {
	project := s.FetchProjectByName(projectName)
	project.BDD.AddScript(scriptName, script)
	s.Manager.SaveProject(project)
	return "success"
}

func CopyDir(src string, dst string) error {
	// Récupérer les informations du dossier source (pour les permissions)
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Créer le dossier de destination avec les mêmes permissions
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Lire le contenu du dossier source
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Si c'est un sous-dossier, on appelle la fonction de manière récursive
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Si c'est un fichier, on le copie
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copie un fichier unique de src vers dst
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Récupérer les permissions du fichier source
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Créer le fichier de destination avec les mêmes permissions
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	// Copier efficacement les données
	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) ExportProject(name, typeProject, path string) string {
	var pJson entity.ProjectJSON
	config.CheckCreateDir(name)
	pJson = s.FetchProjectByName(name)
	switch typeProject {
	case "api":
		config.CheckCreateDir(fmt.Sprintf("%s/api", name))
		s.Manager.ExporterAPI(pJson.ToModel())
	case "models":
		config.CheckCreateDir(fmt.Sprintf("%s/models", name))
		s.Manager.ExporterDB(pJson.ToModel())
	case "frontend":
		// config.CheckCreateDir(fmt.Sprintf("%s/frontend", name))
		s.Manager.ExporterStaticSite(pJson.ToModel())
	case "full":
		config.CheckCreateDir(fmt.Sprintf("%s/webapp", name))
		s.Manager.ExporterWebApp(pJson.ToModel())
	}
	CopyDir(filepath.Join(config.PROJECT_DIR, name), filepath.Join(path, name))
	return "Success"
}

func (s *ProjectService) CreateProject(w http.ResponseWriter, r http.Request) string {

	var name, description, projectType string
	
	pJson := entity.ProjectJSON{
		Nom:         name,
		Description: description,
		Type:        projectType,
	}
	if s.Manager.CheckIfExist(name) {
		return "Project already exists"
	}
	s.Manager.Create(pJson.ToModel())
	s.Manager.SaveProject(pJson)
	switch pJson.Type {
	case "bdd":
		s.Manager.ExporterDB(pJson.ToModel())
	case "static":
		s.Manager.ExporterStaticSite(pJson.ToModel())
	case "api":
		s.Manager.ExporterAPI(pJson.ToModel())
	case "webapp":
		s.Manager.ExporterWebApp(pJson.ToModel())
	}
	fmt.Println("create")
	return "Success"
}

func (s *ProjectService) OpenInNavigator(uri string) {}

func (s *ProjectService) SaveProject(name, project string) string {
	var pJson entity.ProjectJSON
	json.Unmarshal([]byte(project), &pJson)
	s.Manager.SaveProject(pJson)
	switch pJson.Type {
	case "bdd":
		s.Manager.ExporterDB(pJson.ToModel())
	case "static":
		s.Manager.ExporterStaticSite(pJson.ToModel())
	case "api":
		s.Manager.ExporterAPI(pJson.ToModel())
	case "webapp":
		s.Manager.ExporterWebApp(pJson.ToModel())
	}
	return "Success"
}

func (s *ProjectService) SaveProjectBDD(name string, bddJson string) string {
	var bdd entity.BDDJSON
	pJson := s.Manager.LoadProject(name)
	json.Unmarshal([]byte(bddJson), &bdd)
	pJson.BDD = bdd
	pJson.Nom = name
	s.Manager.SaveProject(pJson)
	return "Success"
}

func (mgr *ProjectService) FeychProjectFiles(name string) []entity.FileNode {
	file, _ := mgr.Manager.GetProjectFiles(name)
	return file
}

func (s *ProjectService) GetFileContent(projectName, path string) string {
	content, err := s.Manager.GetFileContent(projectName, path)
	if err != nil {
		return fmt.Sprintf("Erreur lors de la lecture du fichier : %v", err)
	}
	return content
}

func (s *ProjectService) FetchProjects() []entity.ProjectJSON {
	var projectList []entity.ProjectJSON
	for _, val := range s.Manager.Projects {
		projectList = append(projectList, val.ToJSON())
	}
	return projectList
}

func (s *ProjectService) FetchProjectByName(name string) entity.ProjectJSON {
	var pJson entity.ProjectJSON
	filePath := filepath.Join(config.PROJECT_DIR, name+".json")
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err, "by name")
	}
	json.Unmarshal(file, &pJson)
	return pJson
}

func (s *ProjectService) FetchFolderForUpload(folder string) []entity.FileNode {
	var dirPath string
	var nodes []entity.FileNode
	if folder == "" {
		dirPath, _ = os.UserHomeDir()
	} else {
		dirPath = folder
	}
	dir, err := os.ReadDir(dirPath)
	if os.IsNotExist(err) {
		return []entity.FileNode{}
	}
	for _, entry := range dir {
		info, err := entry.Info()

		if err != nil {
			continue
		}

		if entry.IsDir() && strings.Split(entry.Name(), "")[0] == "." && entry.Name() != ".turbostack" {
			continue
		}

		if !entry.IsDir() && strings.Split(entry.Name(), "")[0] == "." {
			continue
		}

		if !entry.IsDir() {
			if strings.Contains(entry.Name(), ".pdf") || strings.Contains(entry.Name(), ".rar") || strings.Contains(entry.Name(), ".docx") || strings.Contains(entry.Name(), ".doc") || strings.Contains(entry.Name(), ".txt") || strings.Contains(entry.Name(), ".ase") || strings.Contains(entry.Name(), ".tmp") || strings.Contains(entry.Name(), ".pixi") || strings.Contains(entry.Name(), ".odt") || strings.Contains(entry.Name(), ".sh") || strings.Contains(entry.Name(), ".md") || strings.Contains(entry.Name(), ".db") || strings.Contains(entry.Name(), ".amgp") || strings.Contains(entry.Name(), ".zip") || strings.Contains(entry.Name(), ".json") {
				continue
			}
		}
		node := entity.FileNode{
			Name:  entry.Name(),
			Path:  filepath.Join(dirPath, entry.Name()),
			IsDir: entry.IsDir(),
			Size:  info.Size(),
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func (s *ProjectService) GetImageAsBase64(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var mimeType string
	switch filepath.Ext(path) {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".mp4":
		mimeType = "video/mp4"
	default:
		mimeType = "application/octet-stream"
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}
