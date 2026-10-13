package entity

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"

	"github.com/MaminirinaEdwino/turbostackapi/src/config"
)

type FileNode struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"is_dir"`
	Children []FileNode `json:"children,omitempty"`
	Size     int64      `json:"size"`
}

type ProjectManager struct {
	Projects []Project
}

func (mgr *ProjectManager) LoadProjects() error {

	projectFiles, err := os.ReadDir(config.PROJECT_DIR)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, val := range projectFiles {
		fmt.Println(val.Name(), "read file pro")
		if !val.IsDir() {
			var pJson ProjectJSON
			var p Project
			filePath := fmt.Sprintf("%s/%s", config.PROJECT_DIR, val.Name())
			fmt.Println(filePath, "file path ito")
			saveFile, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println(err.Error(), "eraera")
			}

			json.Unmarshal(saveFile, &pJson)
			p = pJson.ToModel()
			mgr.Projects = append(mgr.Projects, p)
		}
	}
	// fmt.Println(mgr.Projects, "prjcts")
	sort.Slice(mgr.Projects, func(i, j int) bool {
		return mgr.Projects[i].updateAt.After(mgr.Projects[j].updateAt)
	})
	return nil
}

func (mgr *ProjectManager) SaveProjects() error {

	file, err := os.OpenFile("project.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	pJson := make([]ProjectJSON, 0, len(mgr.Projects))
	for i := range mgr.Projects {
		pJson = append(pJson, mgr.Projects[i].ToJSON())
	}

	jsonValue, err := json.MarshalIndent(pJson, "", "    ")
	if err != nil {
		return err
	}
	file.Write(jsonValue)
	return nil
}

func (mgr *ProjectManager) LoadProject(name string) ProjectJSON {
	var pJson ProjectJSON
	filePath := filepath.Join(config.PROJECT_DIR, name)
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err, "load lty ee")
		return ProjectJSON{}
	}
	json.Unmarshal(file, &pJson)
	return pJson
}

func (mgr *ProjectManager) SaveProject(project ProjectJSON) error {
	filename := project.Nom
	filepath := fmt.Sprintf("%s/%s.json", config.PROJECT_DIR, filename)
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	pModel := project.ToModel()
	pModel.SetUpdateAt()
	project = pModel.ToJSON()
	jsonData, err := json.MarshalIndent(project, "", "    ")
	if err != nil {
		return nil
	}
	file.Write(jsonData)
	return nil
}

func (mgr *ProjectManager) CheckIfExist(name string) bool {
	for _, project := range mgr.Projects {
		if project.GetNom() == name {
			return true
		}
	}
	return false
}

func (mgr *ProjectManager) GetAll() []Project {
	return mgr.Projects
}

func (mgr *ProjectManager) Create(project Project) {
	var tmp Project
	tmp.SetNom(project.GetNom())
	tmp.SetType(project.GetType())
	tmp.SetCreatedAt()
	tmp.SetUpdateAt()
	mgr.Projects = append(mgr.Projects, tmp)
}

func (mgr *ProjectManager) Update(project Project) {
	for i, val := range mgr.Projects {
		if val.GetNom() == project.GetNom() {
			mgr.Projects[i] = project
			break
		}
	}
}

func (mgr *ProjectManager) Delete(project Project) {
	for i, val := range mgr.Projects {
		if val.GetNom() == project.GetNom() {
			mgr.Projects = slices.Delete(mgr.Projects, i, i+1)
			projectPath := filepath.Join(config.PROJECT_DIR, project.nom)

			os.Remove(projectPath + ".json")
			os.RemoveAll(projectPath)
			break
		}
	}
}

func (mgr *ProjectManager) GetByName(name string) Project {
	for _, val := range mgr.Projects {
		if val.GetNom() == name {
			return val
		}
	}
	return Project{}
}

func (mgr *ProjectManager) GetProjectFiles(projectName string) ([]FileNode, error) {
	rootPath := filepath.Join(config.PROJECT_DIR, projectName)
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return []FileNode{}, nil
	}
	return walkDir(rootPath, "")
}

func (mgr *ProjectManager) GetFileContent(projectName string, filePath string) (string, error) {
	fullPath := filepath.Join(config.PROJECT_DIR, projectName, filePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func walkDir(fullPath, relPath string) ([]FileNode, error) {
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}
	var nodes []FileNode
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		node := FileNode{
			Name:  entry.Name(),
			Path:  filepath.Join(relPath, entry.Name()),
			IsDir: entry.IsDir(),
			Size:  info.Size(),
		}
		if entry.IsDir() {
			children, _ := walkDir(filepath.Join(fullPath, entry.Name()), node.Path)
			node.Children = children
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
