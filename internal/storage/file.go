package storage

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SuvorovSergey/gitlab-automerge/internal/config"
	"github.com/SuvorovSergey/gitlab-automerge/pkg/rest"
)

type FileStorage struct {
	Token      string
	HTTPClient rest.BaseClient
}

func NewFileStorage(cfg *config.GitlabConfig) *FileStorage {
	return &FileStorage{
		Token: cfg.Token,
		HTTPClient: rest.BaseClient{
			BaseURL: cfg.BaseUrl,
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
		},
	}
}

func (s *FileStorage) GetAll(projectId int) ([]byte, error) {
	var result []byte
	url := fmt.Sprintf("/api/v4/projects/%d/repository/tree?per_page=100", projectId)

	req, err := s.HTTPClient.CreateRequest(http.MethodGet, url)
	req.Header.Set("PRIVATE-TOKEN", s.Token)

	if err != nil {
		return nil, err
	}

	res, err := s.HTTPClient.SendRequest(req)

	if err != nil {
		return nil, err
	}

	if res.IsOk {
		result, err = res.ReadBody()
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, err
}

func (s *FileStorage) GetOne(projectId int, fileId string) ([]byte, error) {
	var result []byte
	url := fmt.Sprintf("/api/v4/projects/%d/repository/blobs/%s/raw", projectId, fileId)

	req, err := s.HTTPClient.CreateRequest(http.MethodGet, url)
	req.Header.Set("PRIVATE-TOKEN", s.Token)

	if err != nil {
		return nil, err
	}

	res, err := s.HTTPClient.SendRequest(req)

	if err != nil {
		return nil, err
	}

	if res.IsOk {
		result, err = res.ReadBody()
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, err
}
