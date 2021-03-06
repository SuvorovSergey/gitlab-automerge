package gitlab

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/SuvorovSergey/gitlab-automerge/internal/config"
	"github.com/SuvorovSergey/gitlab-automerge/internal/entity"
	"github.com/SuvorovSergey/gitlab-automerge/internal/storage"
)

type Gitlab struct {
	configFile     string
	ProjectStorage *storage.ProjectStorage
	FileStorage    *storage.FileStorage
}

func NewGitlab(cfg *config.GitlabConfig) *Gitlab {
	return &Gitlab{
		ProjectStorage: storage.NewProjectStorage(cfg),
		FileStorage:    storage.NewFileStorage(cfg),
		configFile:     cfg.ConfigFile,
	}
}

//fetch user projects
func (g *Gitlab) Projects() []entity.Project {
	var projects []entity.Project
	r, err := g.ProjectStorage.GetAll()
	if err != nil {
		log.Println(err)
	}

	if err = json.Unmarshal(r, &projects); err != nil {
		log.Println("failed to decode projects due to error %w", err)
	}

	return projects
}

//fetch all user projects and fetch configuration for every project
func (g *Gitlab) ProjectsWithConfig() []entity.Project {
	var wg sync.WaitGroup

	log.Println("fetching projects...")
	projects := g.Projects()

	for i := range projects {
		wg.Add(1)
		go func(p *entity.Project) {
			defer wg.Done()
			config, err := g.ProjectConfig(p.Id)
			if err != nil {
				log.Println("cant fetch configuration for project: ", p.Name)
			}
			//save config to project struct
			if config != nil {
				p.Config = config
			}
		}(&projects[i])
	}
	wg.Wait()

	var withConfig []entity.Project
	//remove projects without automerge configuration
	for i := 0; i < len(projects); i++ {
		if projects[i].Config != nil {
			log.Printf("found project: %s (id: %d, config:%+v)", projects[i].Name, projects[i].Id, projects[i].Config)
			withConfig = append(withConfig, projects[i])
		}
	}

	return withConfig
}

//fetch automerge config for project
func (g *Gitlab) ProjectConfig(projectId int) (*entity.AutomergeConfig, error) {
	var files []entity.File
	var config *entity.AutomergeConfig

	r, err := g.FileStorage.GetAll(projectId)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(r, &files); err != nil {
		log.Println("failed to decode files due to error %w", err)
		return nil, err
	}

	for _, file := range files {
		if file.Name == g.configFile {
			r, err := g.FileStorage.GetOne(projectId, file.Id)
			if err != nil {
				log.Println("failed to fetch project config %w", err)
				return nil, err

			}

			if err = json.Unmarshal(r, &config); err != nil {
				log.Println("failed to decode project config due to error %w", err)
				return nil, err
			}

			return config, nil
		}
	}

	return nil, nil
}

//fetch project merge requests
func (g *Gitlab) MergeRequests(p *entity.Project) []entity.MergeRequest {
	var mr []entity.MergeRequest

	r, err := g.ProjectStorage.GetMergeRequests(p.Id)

	if err != nil {
		log.Println(err)
	}

	if err = json.Unmarshal(r, &mr); err != nil {
		log.Println("failed to decode merge requests due to error %w", err)
	}

	return mr
}

//accept merge request
func (g *Gitlab) AcceptMergeRequest(pId, mId int) error {
	return g.ProjectStorage.AcceptMergeRequest(pId, mId)
}

//fetch all merge requests for every project in slice and accept
func (g *Gitlab) AcceptAllMergeRequests(projects []entity.Project) {
	for i := range projects {
		go func(project *entity.Project) {
			log.Printf("fetching merge requests for project \"%s\"...", project.Name)
			mergeRequests := g.MergeRequests(project)

			for _, mr := range mergeRequests {
				log.Printf("found in \"%s\": %s from %s to %s, id:%d",
					project.Name, mr.Title, mr.SourceBranch, mr.TargetBranch, mr.Iid)

				if mr.Upvotes-mr.Downvotes >= project.Config.UpvotesThreshold {
					log.Printf("trying to merge %s", mr.Title)
					if err := g.AcceptMergeRequest(project.Id, mr.Iid); err != nil {
						log.Printf("ERROR during merging Merge request %s: %+v", mr.Title, err)
					} else {
						log.Printf("SUCCESS merge %s from %s to %s", mr.Title, mr.SourceBranch, mr.TargetBranch)
					}
				}
			}
		}(&projects[i])
	}
}
