package storage

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SuvorovSergey/gitlab-automerge/internal/config"
	"github.com/SuvorovSergey/gitlab-automerge/pkg/rest"
)

type ProjectStorage struct {
	Token      string
	HTTPClient rest.BaseClient
}

func NewProjectStorage(cfg *config.GitlabConfig) *ProjectStorage {
	return &ProjectStorage{
		Token: cfg.Token,
		HTTPClient: rest.BaseClient{
			BaseURL: cfg.BaseUrl,
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
		},
	}
}

func (s *ProjectStorage) GetAll() ([]byte, error) {
	var result []byte
	url := "/api/v4/projects/?per_page=100&membership=true&min_access_level=40&archived=false&simple=true"

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

//only WIP == no
func (s *ProjectStorage) GetMergeRequests(pId int) ([]byte, error) {
	var result []byte
	url := "/api/v4/projects/" + strconv.Itoa(pId) + "/merge_requests?state=opened&wip=no"

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

func (s *ProjectStorage) AcceptMergeRequest(pId, mrId int) error {
	url := fmt.Sprintf("/api/v4/projects/%d/merge_requests/%d/merge?should_remove_source_branch=false", pId, mrId)
	req, err := s.HTTPClient.CreateRequest(http.MethodPut, url)
	req.Header.Set("PRIVATE-TOKEN", s.Token)

	if err != nil {
		return err
	}

	res, err := s.HTTPClient.SendRequest(req)

	if err != nil {
		return err
	}

	if res.IsOk {
		return nil
	}

	return errors.New(res.Error.ToString())
}
