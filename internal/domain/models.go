package domain

type Dataset struct {
	Data            []int32  `json:"data"`
	BackgroundColor []string `json:"backgroundColor"`
}

type Stats struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}
