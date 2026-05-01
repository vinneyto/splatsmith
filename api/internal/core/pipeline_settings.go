package core

import "time"

type PipelineSettingsRecordType string

const (
	PipelineSettingsRecordTypePreset   PipelineSettingsRecordType = "preset"
	PipelineSettingsRecordTypeSnapshot PipelineSettingsRecordType = "snapshot"
)

type PipelineSettings struct {
	VideoProcessing VideoProcessingSettings `json:"videoProcessing"`
	Reconstruction  ReconstructionSettings  `json:"reconstruction"`
	Training        TrainingSettings        `json:"training"`
	PostProcessing  PostProcessingSettings  `json:"postProcessing"`
	SphericalCamera SphericalCameraSettings `json:"sphericalCamera"`
	Segmentation    SegmentationSettings    `json:"segmentation"`
}

type VideoProcessingSettings struct {
	MaxNumImages       string `json:"maxNumImages"`
	VideoStartTime     string `json:"videoStartTime"`
	VideoStopTime      string `json:"videoStopTime"`
	FilterBlurryImages string `json:"filterBlurryImages"`
}

type ReconstructionSettings struct {
	Enable                          string             `json:"enable"`
	SoftwareName                    string             `json:"softwareName"`
	EnableEnhancedFeatureExtraction string             `json:"enableEnhancedFeatureExtraction"`
	MatchingMethod                  string             `json:"matchingMethod"`
	EnableFlHeuristic               string             `json:"enableFlHeuristic"`
	FlHeuristicValue                string             `json:"flHeuristicValue"`
	EnableFlMetric                  string             `json:"enableFlMetric"`
	FlMetricValue                   string             `json:"flMetricValue"`
	PosePriors                      PosePriorsSettings `json:"posePriors"`
}

type PosePriorsSettings struct {
	UsePosePriorColmapModelFiles string                         `json:"usePosePriorColmapModelFiles"`
	UsePosePriorTransformJSON    PosePriorTransformJSONSettings `json:"usePosePriorTransformJson"`
}

type PosePriorTransformJSONSettings struct {
	Enable               string `json:"enable"`
	SourceCoordinateName string `json:"sourceCoordinateName"`
	PoseIsWorldToCam     string `json:"poseIsWorldToCam"`
}

type TrainingSettings struct {
	Enable             string `json:"enable"`
	MaxSteps           string `json:"maxSteps"`
	Model              string `json:"model"`
	ThreeDIsp          string `json:"3dIsp"`
	PreserveSceneScale string `json:"preserveSceneScale"`
	EnableDepthLoss    string `json:"enableDepthLoss"`
}

type PostProcessingSettings struct {
	CropOutputBounds  string `json:"cropOutputBounds"`
	CropMode          string `json:"cropMode"`
	CleanSplat        string `json:"cleanSplat"`
	EnableSpz         string `json:"enableSpz"`
	EnableSog         string `json:"enableSog"`
	EnableUsdz        string `json:"enableUsdz"`
	EnableVideoExport string `json:"enableVideoExport"`
	PlyCoords         string `json:"plyCoords"`
	SpzCoords         string `json:"spzCoords"`
	SogCoords         string `json:"sogCoords"`
	UsdzCoords        string `json:"usdzCoords"`
}

type SphericalCameraSettings struct {
	Enable                       string `json:"enable"`
	CubeFacesToRemove            string `json:"cubeFacesToRemove"`
	OptimizeSequentialFrameOrder string `json:"optimizeSequentialFrameOrder"`
}

type SegmentationSettings struct {
	BackgroundRemoval BackgroundRemovalSettings `json:"backgroundRemoval"`
	ObjectRemoval     ObjectRemovalSettings     `json:"objectRemoval"`
}

type BackgroundRemovalSettings struct {
	Enable        string `json:"enable"`
	Model         string `json:"model"`
	MaskThreshold string `json:"maskThreshold"`
}

type ObjectRemovalSettings struct {
	Enable  string `json:"enable"`
	Action  string `json:"action"`
	Objects string `json:"objects"`
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
			ThreeDIsp:          "none",
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
	RecordID   string                     `json:"recordId"`
	UserID     string                     `json:"userId"`
	RecordType PipelineSettingsRecordType `json:"recordType"`
	Name       string                     `json:"name"`
	Settings   PipelineSettings           `json:"settings"`
	CreatedAt  time.Time                  `json:"createdAt"`
	UpdatedAt  time.Time                  `json:"updatedAt"`
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
