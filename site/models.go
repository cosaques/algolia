package main

type (
	CountResponse struct {
		Count int `json:"count"`
	}

	PopularResponse struct {
		Queries []QueryCountResponse `json:"queries"`
	}

	QueryCountResponse struct {
		Query string `json:"query"`
		Count int    `json:"count"`
	}

	MonitoringMsg struct {
		Indexed int `json:"indexed"`
	}
)
