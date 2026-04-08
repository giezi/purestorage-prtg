package pureapi

// APIVersionResponse is returned by GET /api/api_version (no auth required).
type APIVersionResponse struct {
	Version []string `json:"version"`
}

// ArraySpaceResponse is returned by GET /api/2.x/arrays/space.
type ArraySpaceResponse struct {
	Items []ArraySpaceItem `json:"items"`
}

type ArraySpaceItem struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Capacity int64      `json:"capacity"`
	Parity   float64    `json:"parity"`
	Space    ArraySpace `json:"space"`
	Time     int64      `json:"time"`
}

type ArraySpace struct {
	DataReduction    float64 `json:"data_reduction"`
	Shared           int64   `json:"shared"`
	Snapshots        int64   `json:"snapshots"`
	System           int64   `json:"system"`
	ThinProvisioning float64 `json:"thin_provisioning"`
	TotalPhysical    int64   `json:"total_physical"`
	TotalProvisioned int64   `json:"total_provisioned"`
	TotalReduction   float64 `json:"total_reduction"`
	Unique           int64   `json:"unique"`
	Virtual          int64   `json:"virtual"`
	Replication      int64   `json:"replication"`
}

// ArrayPerformanceResponse is returned by GET /api/2.x/arrays/performance.
type ArrayPerformanceResponse struct {
	Items []ArrayPerformanceItem `json:"items"`
}

type ArrayPerformanceItem struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Time              int64   `json:"time"`
	ReadsPerSec       int64   `json:"reads_per_sec"`
	WritesPerSec      int64   `json:"writes_per_sec"`
	OthersPerSec      int64   `json:"others_per_sec"`
	ReadBytesPerSec   int64   `json:"read_bytes_per_sec"`
	WriteBytesPerSec  int64   `json:"write_bytes_per_sec"`
	UsecPerReadOp     float64 `json:"usec_per_read_op"`
	UsecPerWriteOp    float64 `json:"usec_per_write_op"`
	UsecPerOtherOp    float64 `json:"usec_per_other_op"`
	BytesPerOp        int64   `json:"bytes_per_op"`
	BytesPerRead      int64   `json:"bytes_per_read"`
	BytesPerWrite     int64   `json:"bytes_per_write"`
	QueueDepth        int64   `json:"queue_depth"`
	QueueUsecPerReadOp  float64 `json:"queue_usec_per_read_op"`
	QueueUsecPerWriteOp float64 `json:"queue_usec_per_write_op"`
	SanUsecPerReadOp    float64 `json:"san_usec_per_read_op"`
	SanUsecPerWriteOp   float64 `json:"san_usec_per_write_op"`
	ServiceUsecPerReadOp  float64 `json:"service_usec_per_read_op"`
	ServiceUsecPerWriteOp float64 `json:"service_usec_per_write_op"`
}

// HardwareResponse is returned by GET /api/2.x/hardware.
type HardwareResponse struct {
	Items []HardwareItem `json:"items"`
}

type HardwareItem struct {
	Name            string  `json:"name"`
	Status          string  `json:"status"`
	Type            string  `json:"type"`
	Details         string  `json:"details"`
	Temperature     float64 `json:"temperature"`
	Voltage         float64 `json:"voltage"`
	Model           string  `json:"model"`
	Serial          string  `json:"serial"`
	IdentifyEnabled bool    `json:"identify_enabled"`
	Index           int     `json:"index"`
	Slot            int     `json:"slot"`
	Speed           int64   `json:"speed"`
}

// VolumeSpaceResponse is returned by GET /api/2.x/volumes/space.
type VolumeSpaceResponse struct {
	Items []VolumeSpaceItem `json:"items"`
}

type VolumeSpaceItem struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Space VolumeSpace `json:"space"`
	Time  int64       `json:"time"`
}

type VolumeSpace struct {
	DataReduction    float64 `json:"data_reduction"`
	Shared           int64   `json:"shared"`
	Snapshots        int64   `json:"snapshots"`
	System           int64   `json:"system"`
	ThinProvisioning float64 `json:"thin_provisioning"`
	TotalPhysical    int64   `json:"total_physical"`
	TotalProvisioned int64   `json:"total_provisioned"`
	TotalReduction   float64 `json:"total_reduction"`
	Unique           int64   `json:"unique"`
	Virtual          int64   `json:"virtual"`
}
