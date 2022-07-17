package entity

type Project struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Config *AutomergeConfig
}
