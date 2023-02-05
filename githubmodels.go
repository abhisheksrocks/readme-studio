package main

import "fmt"

type GithubResultModel[data any] struct {
	Data   data               `json:"data"`
	Errors []GithubErrorModel `json:"errors"`
}

type GithubModel interface {
	makeQuery(any) string
}

// func

// type GithubModelInterface[data any] interface {
// 	makeQuery(...string) string
// 	resultStr
// }

type GithubErrorModel struct {
	Path       []string `json:"path"`
	Extensions struct {
		Code      string `json:"code"`
		TypeName  string `json:"typeName"`
		FieldName string `json:"fieldName"`
	} `json:"extensions"`
	Locations []struct {
		Line   int `json:"line"`
		Column int `json:"column"`
	} `json:"locations"`
	Message string `json:"message"`
}

type GithubRepositoryCardModel struct {
	Repository struct {
		Name        string `json:"name"`
		IsArchived  bool   `json:"isArchived"`
		Description string `json:"description"`
		Parent      struct {
			NameWithOwner string `json:"nameWithOwner"`
		} `json:"parent"`
		Languages struct {
			Nodes []struct {
				Name  string `json:"name"`
				Color string `json:"color"`
			} `json:"nodes"`
		} `json:"languages"`
		StargazerCount int `json:"stargazerCount"`
		ForkCount      int `json:"forkCount"`
	} `json:"repository"`
}

func (*GithubRepositoryCardModel) makeQuery(name string, owner string) string {
	return fmt.Sprintf(
		`{
			repository(name: "%s", owner: "%s") {
				name
				isArchived
				description
				parent {
					nameWithOwner
				}
				languages(first: 1, orderBy: {field: SIZE, direction: DESC}) {
					nodes {
						name
						color
					}
				}
				stargazerCount
				forkCount
			}
			}`, name, owner)
}

// func (*GithubRepositoryCardModel) resultStruct() GithubResultModel[GithubRepositoryCardModel] {
// 	return GithubResultModel[GithubRepositoryCardModel]{
// 		// Data:   new(GithubRepositoryCardModel),
// 		Data:   GithubRepositoryCardModel{},
// 		Errors: []GithubErrorModel{},
// 	}
// }

// func (*GithubModel) resultStruct() GithubResultModel[GithubModel] {
// 	return GithubResultModel[GithubModel]{}
// }
