package entity

type MergeRequest struct {
	Id             int    `json:"id"`
	Iid            int    `json:"iid"`
	ProjectId      int    `json:"project_id"`
	Title          string `json:"title"`
	SourceBranch   string `json:"source_branch"`
	TargetBranch   string `json:"target_branch"`
	Upvotes        int    `json:"upvotes"`
	Downvotes      int    `json:"downvotes"`
	WorkInProgress bool   `json:"draft"`
}
