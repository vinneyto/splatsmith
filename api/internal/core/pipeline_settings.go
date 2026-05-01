package core

import "time"

type PipelineSettingsRecordType string

const (
	PipelineSettingsRecordTypePreset   PipelineSettingsRecordType = "preset"
	PipelineSettingsRecordTypeSnapshot PipelineSettingsRecordType = "snapshot"
)

type PipelineSettings struct {
	VideoProcessing VideoProcessingSettings
	Reconstruction  ReconstructionSettings
	Training        TrainingSettings
	PostProcessing  PostProcessingSettings
	SphericalCamera SphericalCameraSettings
	Segmentation    SegmentationSettings
}

type VideoProcessingSettings struct {
	MaxNumImages       string
	VideoStartTime     string
	VideoStopTime      string
	FilterBlurryImages string
}

type ReconstructionSettings struct {
	Enable                          string
	SoftwareName                    string
	EnableEnhancedFeatureExtraction string
	MatchingMethod                  string
	EnableFlHeuristic               string
	FlHeuristicValue                string
	EnableFlMetric                  string
	FlMetricValue                   string
	PosePriors                      PosePriorsSettings
}

type PosePriorsSettings struct {
	UsePosePriorColmapModelFiles string
	UsePosePriorTransformJSON    PosePriorTransformJSONSettings
}

type PosePriorTransformJSONSettings struct {
	Enable               string
	SourceCoordinateName string
	PoseIsWorldToCam     string
}

type TrainingSettings struct {
	Enable             string
	MaxSteps           string
	Model              string
	ThreeDISP          string
	PreserveSceneScale string
	EnableDepthLoss    string
}

type PostProcessingSettings struct {
	CropOutputBounds  string
	CropMode          string
	CleanSplat        string
	EnableSpz         string
	EnableSog         string
	EnableUsdz        string
	EnableVideoExport string
	PlyCoords         string
	SpzCoords         string
	SogCoords         string
	UsdzCoords        string
}

type SphericalCameraSettings struct {
	Enable                       string
	CubeFacesToRemove            string
	OptimizeSequentialFrameOrder string
}

type SegmentationSettings struct {
	BackgroundRemoval BackgroundRemovalSettings
	ObjectRemoval     ObjectRemovalSettings
}

type BackgroundRemovalSettings struct {
	Enable        string
	Model         string
	MaskThreshold string
}

type ObjectRemovalSettings struct {
	Enable  string
	Action  string
	Objects string
}

func NewDefaultPipelineSettings() PipelineSettings {
	return PipelineSettings{
		VideoProcessing: VideoProcessingSettings{
			MaxNumImages:       "300",
			VideoStartTime:     "0",
			VideoStopTime:      "None",
			FilterBlurryImages: "true",
		},
		Reconstruction: ReconstructionSettings{
			Enable:                          "true",
			SoftwareName:                    "glomap",
			EnableEnhancedFeatureExtraction: "false",
			MatchingMethod:                  "sequential",
			EnableFlHeuristic:               "false",
			FlHeuristicValue:                "1.2",
			EnableFlMetric:                  "false",
			FlMetricValue:                   "24",
			PosePriors: PosePriorsSettings{
				UsePosePriorColmapModelFiles: "false",
				UsePosePriorTransformJSON: PosePriorTransformJSONSettings{
					Enable:               "false",
					SourceCoordinateName: "arkit",
					PoseIsWorldToCam:     "true",
				},
			},
		},
		Training: TrainingSettings{
			Enable:             "true",
			MaxSteps:           "15000",
			Model:              "splatfacto",
			ThreeDISP:          "none",
			PreserveSceneScale: "false",
			EnableDepthLoss:    "false",
		},
		PostProcessing: PostProcessingSettings{
			CropOutputBounds:  "false",
			CropMode:          "environment",
			CleanSplat:        "false",
			EnableSpz:         "true",
			EnableSog:         "true",
			EnableUsdz:        "true",
			EnableVideoExport: "true",
			PlyCoords:         "rhyu",
			SpzCoords:         "rhyu",
			SogCoords:         "rhyu",
			UsdzCoords:        "rhyu",
		},
		SphericalCamera: SphericalCameraSettings{
			Enable:                       "false",
			CubeFacesToRemove:            "['down', 'up']",
			OptimizeSequentialFrameOrder: "true",
		},
		Segmentation: SegmentationSettings{
			BackgroundRemoval: BackgroundRemovalSettings{
				Enable:        "false",
				Model:         "u2net",
				MaskThreshold: "0.6",
			},
			ObjectRemoval: ObjectRemovalSettings{
				Enable:  "false",
				Action:  "erase",
				Objects: "['human']",
			},
		},
	}
}

type PipelineSettingsRecord struct {
	RecordID   string
	UserID     string
	RecordType PipelineSettingsRecordType
	Name       string
	Settings   PipelineSettings
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type PipelineSettingsListFilter struct {
	UserID     string
	RecordType *PipelineSettingsRecordType
	Name       *string
	Limit      int
	Offset     int
}

type CreatePipelineSettingsInput struct {
	UserID     string
	RecordType PipelineSettingsRecordType
	Name       string
	Settings   PipelineSettings
}

type UpdatePipelineSettingsInput struct {
	RecordID string
	UserID   string
	Name     *string
	Settings *PipelineSettings
}
