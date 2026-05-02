package core

import (
	"encoding/json"
	"fmt"
	"time"
)

type ReconstructionSoftwareName string

const (
	ReconstructionSoftwareNameGlomap      ReconstructionSoftwareName = "glomap"
	ReconstructionSoftwareNameColmap      ReconstructionSoftwareName = "colmap"
	ReconstructionSoftwareNameHloc        ReconstructionSoftwareName = "hloc"
	ReconstructionSoftwareNameMapAnything ReconstructionSoftwareName = "map_anything"
)

func (v ReconstructionSoftwareName) IsValid() bool {
	switch v {
	case ReconstructionSoftwareNameGlomap, ReconstructionSoftwareNameColmap, ReconstructionSoftwareNameHloc, ReconstructionSoftwareNameMapAnything:
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
	ReconstructionMatchingMethodVocab      ReconstructionMatchingMethod = "vocab"
	ReconstructionMatchingMethodVocabTree  ReconstructionMatchingMethod = "vocab_tree"
	ReconstructionMatchingMethodTransitive ReconstructionMatchingMethod = "transitive"
)

func (v ReconstructionMatchingMethod) IsValid() bool {
	switch v {
	case ReconstructionMatchingMethodSequential,
		ReconstructionMatchingMethodExhaustive,
		ReconstructionMatchingMethodSpatial,
		ReconstructionMatchingMethodVocab,
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
	TrainingModelSplatfacto       TrainingModel = "splatfacto"
	TrainingModelSplatfactoBig    TrainingModel = "splatfacto-big"
	TrainingModelSplatfactoMCMC   TrainingModel = "splatfacto-mcmc"
	TrainingModelSplatfactoWLight TrainingModel = "splatfacto-w-light"
	TrainingModel3DGUT            TrainingModel = "3dgut"
	TrainingModel3DGRT            TrainingModel = "3dgrt"
	TrainingModelNerfacto         TrainingModel = "nerfacto"
)

func (v TrainingModel) IsValid() bool {
	switch v {
	case TrainingModelSplatfacto,
		TrainingModelSplatfactoBig,
		TrainingModelSplatfactoMCMC,
		TrainingModelSplatfactoWLight,
		TrainingModel3DGUT,
		TrainingModel3DGRT,
		TrainingModelNerfacto:
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
	TrainingThreeDIspNone     TrainingThreeDIsp = "none"
	TrainingThreeDIspBilagrid TrainingThreeDIsp = "bilagrid"
	TrainingThreeDIspPpisp    TrainingThreeDIsp = "ppisp"
)

func (v TrainingThreeDIsp) IsValid() bool {
	switch v {
	case TrainingThreeDIspNone, TrainingThreeDIspBilagrid, TrainingThreeDIspPpisp:
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
	PostProcessingCropModeRigidBody   PostProcessingCropMode = "rigid_body"
)

func (v PostProcessingCropMode) IsValid() bool {
	switch v {
	case PostProcessingCropModeEnvironment, PostProcessingCropModeRigidBody:
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
	CoordinateSystemLhyu CoordinateSystem = "lhyu"
	CoordinateSystemRhzu CoordinateSystem = "rhzu"
	CoordinateSystemLhzu CoordinateSystem = "lhzu"
)

func (v CoordinateSystem) IsValid() bool {
	switch v {
	case CoordinateSystemRhyu, CoordinateSystemLhyu, CoordinateSystemRhzu, CoordinateSystemLhzu:
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
	BackgroundRemovalModelSAM2  BackgroundRemovalModel = "sam2"
)

func (v BackgroundRemovalModel) IsValid() bool {
	switch v {
	case BackgroundRemovalModelU2Net, BackgroundRemovalModelSAM2:
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
	ObjectRemovalActionErase  ObjectRemovalAction = "erase"
	ObjectRemovalActionRemove ObjectRemovalAction = "remove"
)

func (v ObjectRemovalAction) IsValid() bool {
	switch v {
	case ObjectRemovalActionErase, ObjectRemovalActionRemove:
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

	PipelineSettingsPresetNameStandard = "standard"
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
	Enable                          string                       `json:"enable"`
	SoftwareName                    ReconstructionSoftwareName   `json:"softwareName"`
	EnableEnhancedFeatureExtraction string                       `json:"enableEnhancedFeatureExtraction"`
	MatchingMethod                  ReconstructionMatchingMethod `json:"matchingMethod"`
	EnableFlHeuristic               string                       `json:"enableFlHeuristic"`
	FlHeuristicValue                string                       `json:"flHeuristicValue"`
	EnableFlMetric                  string                       `json:"enableFlMetric"`
	FlMetricValue                   string                       `json:"flMetricValue"`
	PosePriors                      PosePriorsSettings           `json:"posePriors"`
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
	Enable             string            `json:"enable"`
	MaxSteps           string            `json:"maxSteps"`
	Model              TrainingModel     `json:"model"`
	ThreeDIsp          TrainingThreeDIsp `json:"3dIsp"`
	PreserveSceneScale string            `json:"preserveSceneScale"`
	EnableDepthLoss    string            `json:"enableDepthLoss"`
}

type PostProcessingSettings struct {
	CropOutputBounds  string                 `json:"cropOutputBounds"`
	CropMode          PostProcessingCropMode `json:"cropMode"`
	CleanSplat        string                 `json:"cleanSplat"`
	EnableSpz         string                 `json:"enableSpz"`
	EnableSog         string                 `json:"enableSog"`
	EnableUsdz        string                 `json:"enableUsdz"`
	EnableVideoExport string                 `json:"enableVideoExport"`
	PlyCoords         CoordinateSystem       `json:"plyCoords"`
	SpzCoords         CoordinateSystem       `json:"spzCoords"`
	SogCoords         CoordinateSystem       `json:"sogCoords"`
	UsdzCoords        CoordinateSystem       `json:"usdzCoords"`
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
	Enable        string                 `json:"enable"`
	Model         BackgroundRemovalModel `json:"model"`
	MaskThreshold string                 `json:"maskThreshold"`
}

type ObjectRemovalSettings struct {
	Enable  string              `json:"enable"`
	Action  ObjectRemovalAction `json:"action"`
	Objects string              `json:"objects"`
}

func NewDefaultPipelineSettings() PipelineSettings {
	return PipelineSettings{
		VideoProcessing: VideoProcessingSettings{
			MaxNumImages:       "300",
			VideoStartTime:     "0",
			VideoStopTime:      "-1",
			FilterBlurryImages: "true",
		},
		Reconstruction: ReconstructionSettings{
			Enable:                          "true",
			SoftwareName:                    ReconstructionSoftwareNameGlomap,
			EnableEnhancedFeatureExtraction: "false",
			MatchingMethod:                  ReconstructionMatchingMethodSequential,
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
			Model:              TrainingModelSplatfacto,
			ThreeDIsp:          TrainingThreeDIspNone,
			PreserveSceneScale: "false",
			EnableDepthLoss:    "false",
		},
		PostProcessing: PostProcessingSettings{
			CropOutputBounds:  "false",
			CropMode:          PostProcessingCropModeEnvironment,
			CleanSplat:        "false",
			EnableSpz:         "true",
			EnableSog:         "true",
			EnableUsdz:        "true",
			EnableVideoExport: "true",
			PlyCoords:         CoordinateSystemRhyu,
			SpzCoords:         CoordinateSystemRhyu,
			SogCoords:         CoordinateSystemRhyu,
			UsdzCoords:        CoordinateSystemRhyu,
		},
		SphericalCamera: SphericalCameraSettings{
			Enable:                       "false",
			CubeFacesToRemove:            "['down', 'up']",
			OptimizeSequentialFrameOrder: "true",
		},
		Segmentation: SegmentationSettings{
			BackgroundRemoval: BackgroundRemovalSettings{
				Enable:        "false",
				Model:         BackgroundRemovalModelU2Net,
				MaskThreshold: "0.6",
			},
			ObjectRemoval: ObjectRemovalSettings{
				Enable:  "false",
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
