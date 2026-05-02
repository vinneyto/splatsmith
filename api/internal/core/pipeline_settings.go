package core

import (
	"encoding/json"
	"fmt"
	"time"
)

type ReconstructionSoftwareName string

const (
	ReconstructionSoftwareNameGlomap ReconstructionSoftwareName = "glomap"
	ReconstructionSoftwareNameColmap ReconstructionSoftwareName = "colmap"
)

func (v ReconstructionSoftwareName) IsValid() bool {
	switch v {
	case ReconstructionSoftwareNameGlomap, ReconstructionSoftwareNameColmap:
		return true
	default:
		return false
	}
}

func (v *ReconstructionSoftwareName) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed := ReconstructionSoftwareName(raw)
	if !parsed.IsValid() {
		return fmt.Errorf("invalid reconstruction.softwareName: %q", raw)
	}
	*v = parsed
	return nil
}

type ReconstructionMatchingMethod string

const (
	ReconstructionMatchingMethodSequential ReconstructionMatchingMethod = "sequential"
	ReconstructionMatchingMethodExhaustive ReconstructionMatchingMethod = "exhaustive"
	ReconstructionMatchingMethodSpatial    ReconstructionMatchingMethod = "spatial"
	ReconstructionMatchingMethodVocabTree  ReconstructionMatchingMethod = "vocab_tree"
	ReconstructionMatchingMethodTransitive ReconstructionMatchingMethod = "transitive"
)

func (v ReconstructionMatchingMethod) IsValid() bool {
	switch v {
	case ReconstructionMatchingMethodSequential,
		ReconstructionMatchingMethodExhaustive,
		ReconstructionMatchingMethodSpatial,
		ReconstructionMatchingMethodVocabTree,
		ReconstructionMatchingMethodTransitive:
		return true
	default:
		return false
	}
}

func (v *ReconstructionMatchingMethod) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed := ReconstructionMatchingMethod(raw)
	if !parsed.IsValid() {
		return fmt.Errorf("invalid reconstruction.matchingMethod: %q", raw)
	}
	*v = parsed
	return nil
}

type TrainingModel string

const (
	TrainingModelSplatfacto TrainingModel = "splatfacto"
)

func (v TrainingModel) IsValid() bool {
	switch v {
	case TrainingModelSplatfacto:
		return true
	default:
		return false
	}
}

func (v *TrainingModel) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed := TrainingModel(raw)
	if !parsed.IsValid() {
		return fmt.Errorf("invalid training.model: %q", raw)
	}
	*v = parsed
	return nil
}

type TrainingThreeDIsp string

const (
	TrainingThreeDIspNone TrainingThreeDIsp = "none"
)

func (v TrainingThreeDIsp) IsValid() bool {
	switch v {
	case TrainingThreeDIspNone:
		return true
	default:
		return false
	}
}

func (v *TrainingThreeDIsp) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed := TrainingThreeDIsp(raw)
	if !parsed.IsValid() {
		return fmt.Errorf("invalid training.3dIsp: %q", raw)
	}
	*v = parsed
	return nil
}

type PostProcessingCropMode string

const (
	PostProcessingCropModeEnvironment PostProcessingCropMode = "environment"
)

func (v PostProcessingCropMode) IsValid() bool {
	switch v {
	case PostProcessingCropModeEnvironment:
		return true
	default:
		return false
	}
}

func (v *PostProcessingCropMode) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed := PostProcessingCropMode(raw)
	if !parsed.IsValid() {
		return fmt.Errorf("invalid postProcessing.cropMode: %q", raw)
	}
	*v = parsed
	return nil
}

type CoordinateSystem string

const (
	CoordinateSystemRhyu CoordinateSystem = "rhyu"
)

func (v CoordinateSystem) IsValid() bool {
	switch v {
	case CoordinateSystemRhyu:
		return true
	default:
		return false
	}
}

func (v *CoordinateSystem) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed := CoordinateSystem(raw)
	if !parsed.IsValid() {
		return fmt.Errorf("invalid coordinate system: %q", raw)
	}
	*v = parsed
	return nil
}

type BackgroundRemovalModel string

const (
	BackgroundRemovalModelU2Net BackgroundRemovalModel = "u2net"
)

func (v BackgroundRemovalModel) IsValid() bool {
	switch v {
	case BackgroundRemovalModelU2Net:
		return true
	default:
		return false
	}
}

func (v *BackgroundRemovalModel) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed := BackgroundRemovalModel(raw)
	if !parsed.IsValid() {
		return fmt.Errorf("invalid segmentation.backgroundRemoval.model: %q", raw)
	}
	*v = parsed
	return nil
}

type ObjectRemovalAction string

const (
	ObjectRemovalActionErase ObjectRemovalAction = "erase"
)

func (v ObjectRemovalAction) IsValid() bool {
	switch v {
	case ObjectRemovalActionErase:
		return true
	default:
		return false
	}
}

func (v *ObjectRemovalAction) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed := ObjectRemovalAction(raw)
	if !parsed.IsValid() {
		return fmt.Errorf("invalid segmentation.objectRemoval.action: %q", raw)
	}
	*v = parsed
	return nil
}

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
	MaxNumImages       int  `json:"maxNumImages"`
	VideoStartTime     int  `json:"videoStartTime"`
	VideoStopTime      int  `json:"videoStopTime"`
	FilterBlurryImages bool `json:"filterBlurryImages"`
}

type ReconstructionSettings struct {
	Enable                          bool                         `json:"enable"`
	SoftwareName                    ReconstructionSoftwareName   `json:"softwareName"`
	EnableEnhancedFeatureExtraction bool                         `json:"enableEnhancedFeatureExtraction"`
	MatchingMethod                  ReconstructionMatchingMethod `json:"matchingMethod"`
	EnableFlHeuristic               bool                         `json:"enableFlHeuristic"`
	FlHeuristicValue                float64                      `json:"flHeuristicValue"`
	EnableFlMetric                  bool                         `json:"enableFlMetric"`
	FlMetricValue                   float64                      `json:"flMetricValue"`
	PosePriors                      PosePriorsSettings           `json:"posePriors"`
}

type PosePriorsSettings struct {
	UsePosePriorColmapModelFiles bool                           `json:"usePosePriorColmapModelFiles"`
	UsePosePriorTransformJSON    PosePriorTransformJSONSettings `json:"usePosePriorTransformJson"`
}

type PosePriorTransformJSONSettings struct {
	Enable               bool   `json:"enable"`
	SourceCoordinateName string `json:"sourceCoordinateName"`
	PoseIsWorldToCam     bool   `json:"poseIsWorldToCam"`
}

type TrainingSettings struct {
	Enable             bool              `json:"enable"`
	MaxSteps           int               `json:"maxSteps"`
	Model              TrainingModel     `json:"model"`
	ThreeDIsp          TrainingThreeDIsp `json:"3dIsp"`
	PreserveSceneScale bool              `json:"preserveSceneScale"`
	EnableDepthLoss    bool              `json:"enableDepthLoss"`
}

type PostProcessingSettings struct {
	CropOutputBounds  bool                   `json:"cropOutputBounds"`
	CropMode          PostProcessingCropMode `json:"cropMode"`
	CleanSplat        bool                   `json:"cleanSplat"`
	EnableSpz         bool                   `json:"enableSpz"`
	EnableSog         bool                   `json:"enableSog"`
	EnableUsdz        bool                   `json:"enableUsdz"`
	EnableVideoExport bool                   `json:"enableVideoExport"`
	PlyCoords         CoordinateSystem       `json:"plyCoords"`
	SpzCoords         CoordinateSystem       `json:"spzCoords"`
	SogCoords         CoordinateSystem       `json:"sogCoords"`
	UsdzCoords        CoordinateSystem       `json:"usdzCoords"`
}

type SphericalCameraSettings struct {
	Enable                       bool   `json:"enable"`
	CubeFacesToRemove            string `json:"cubeFacesToRemove"`
	OptimizeSequentialFrameOrder bool   `json:"optimizeSequentialFrameOrder"`
}

type SegmentationSettings struct {
	BackgroundRemoval BackgroundRemovalSettings `json:"backgroundRemoval"`
	ObjectRemoval     ObjectRemovalSettings     `json:"objectRemoval"`
}

type BackgroundRemovalSettings struct {
	Enable        bool                   `json:"enable"`
	Model         BackgroundRemovalModel `json:"model"`
	MaskThreshold float64                `json:"maskThreshold"`
}

type ObjectRemovalSettings struct {
	Enable  bool                `json:"enable"`
	Action  ObjectRemovalAction `json:"action"`
	Objects string              `json:"objects"`
}

func NewDefaultPipelineSettings() PipelineSettings {
	return PipelineSettings{
		VideoProcessing: VideoProcessingSettings{
			MaxNumImages:       300,
			VideoStartTime:     0,
			VideoStopTime:      -1,
			FilterBlurryImages: true,
		},
		Reconstruction: ReconstructionSettings{
			Enable:                          true,
			SoftwareName:                    ReconstructionSoftwareNameGlomap,
			EnableEnhancedFeatureExtraction: false,
			MatchingMethod:                  ReconstructionMatchingMethodSequential,
			EnableFlHeuristic:               false,
			FlHeuristicValue:                1.2,
			EnableFlMetric:                  false,
			FlMetricValue:                   24,
			PosePriors: PosePriorsSettings{
				UsePosePriorColmapModelFiles: false,
				UsePosePriorTransformJSON: PosePriorTransformJSONSettings{
					Enable:               false,
					SourceCoordinateName: "arkit",
					PoseIsWorldToCam:     true,
				},
			},
		},
		Training: TrainingSettings{
			Enable:             true,
			MaxSteps:           15000,
			Model:              TrainingModelSplatfacto,
			ThreeDIsp:          TrainingThreeDIspNone,
			PreserveSceneScale: false,
			EnableDepthLoss:    false,
		},
		PostProcessing: PostProcessingSettings{
			CropOutputBounds:  false,
			CropMode:          PostProcessingCropModeEnvironment,
			CleanSplat:        false,
			EnableSpz:         true,
			EnableSog:         true,
			EnableUsdz:        true,
			EnableVideoExport: true,
			PlyCoords:         CoordinateSystemRhyu,
			SpzCoords:         CoordinateSystemRhyu,
			SogCoords:         CoordinateSystemRhyu,
			UsdzCoords:        CoordinateSystemRhyu,
		},
		SphericalCamera: SphericalCameraSettings{
			Enable:                       false,
			CubeFacesToRemove:            "['down', 'up']",
			OptimizeSequentialFrameOrder: true,
		},
		Segmentation: SegmentationSettings{
			BackgroundRemoval: BackgroundRemovalSettings{
				Enable:        false,
				Model:         BackgroundRemovalModelU2Net,
				MaskThreshold: 0.6,
			},
			ObjectRemoval: ObjectRemovalSettings{
				Enable:  false,
				Action:  ObjectRemovalActionErase,
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

func (s PipelineSettings) Validate() error {
	if !s.Reconstruction.SoftwareName.IsValid() {
		return fmt.Errorf("invalid reconstruction.softwareName: %q", s.Reconstruction.SoftwareName)
	}
	if !s.Reconstruction.MatchingMethod.IsValid() {
		return fmt.Errorf("invalid reconstruction.matchingMethod: %q", s.Reconstruction.MatchingMethod)
	}
	if !s.Training.Model.IsValid() {
		return fmt.Errorf("invalid training.model: %q", s.Training.Model)
	}
	if !s.Training.ThreeDIsp.IsValid() {
		return fmt.Errorf("invalid training.3dIsp: %q", s.Training.ThreeDIsp)
	}
	if !s.PostProcessing.CropMode.IsValid() {
		return fmt.Errorf("invalid postProcessing.cropMode: %q", s.PostProcessing.CropMode)
	}
	if !s.PostProcessing.PlyCoords.IsValid() || !s.PostProcessing.SpzCoords.IsValid() || !s.PostProcessing.SogCoords.IsValid() || !s.PostProcessing.UsdzCoords.IsValid() {
		return fmt.Errorf("invalid postProcessing.*Coords value")
	}
	if !s.Segmentation.BackgroundRemoval.Model.IsValid() {
		return fmt.Errorf("invalid segmentation.backgroundRemoval.model: %q", s.Segmentation.BackgroundRemoval.Model)
	}
	if !s.Segmentation.ObjectRemoval.Action.IsValid() {
		return fmt.Errorf("invalid segmentation.objectRemoval.action: %q", s.Segmentation.ObjectRemoval.Action)
	}
	return nil
}
