package model

const (
	TaskTypePDFToWord = iota
	TaskTypeTextToVoice
	TaskTypePictureCluster
)

type Task struct {
	ID        int64  `json:"id"`
	FileInput string `json:"fileInput"`
}
type TaskResult struct {
	ID                 int64  `json:"id"`
	NodeID             string `json:"nodeid"`
	Status             int8   `json:"status"`
	CalType            int    `json:"calType"`
	HandleResult       string `json:"handleResult"`
	HandleBeginDate    string `json:"handleBeginDate"`
	HandleCompleteDate string `json:"handleCompleteDate"`
	HandleRemainDate   string `json:"handleRemainDate"`
}
