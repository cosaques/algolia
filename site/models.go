package main

type (
	// CountResponse contains count of distinct queries.
	CountResponse struct {
		Count int `json:"count"`
	}

	// PopularResponse contains list of top popular queris.
	PopularResponse struct {
		Queries []QueryCountResponse `json:"queries"`
	}

	// QueryCountResponse represent a query and number of times it was occured.
	QueryCountResponse struct {
		Query string `json:"query"`
		Count int    `json:"count"`
	}

	// MonitoringMsg is sent to a dashboard to monitor the index progress.
	MonitoringMsg struct {
		Indexed int `json:"indexed"`
	}
)
