"use client";

import { ChangeEvent } from "react";
import type { PipelineSettings } from "@/store/api/splatmakerApi";
import { Input } from "@/components/ui/input";

type Props = {
  value: PipelineSettings;
  onChange: (next: PipelineSettings) => void;
};

const coordChoices = [
  { label: "Right-Hand, Y-Up (playcanvas)", value: "rhyu" },
  { label: "Left-Hand, Y-Up (babylon.js)", value: "lhyu" },
  { label: "Right-Hand, Z-Up (blender)", value: "rhzu" },
  { label: "Left-Hand, Z-Up (unreal)", value: "lhzu" },
] as const;

function boolValue(v: string) {
  return v === "true";
}

function setByPath<T extends keyof PipelineSettings>(obj: PipelineSettings, section: T, next: PipelineSettings[T]): PipelineSettings {
  return { ...obj, [section]: next };
}

export function PipelineSettingsForm({ value, onChange }: Props) {
  const setBool = (path: (draft: PipelineSettings) => string, apply: (draft: PipelineSettings, v: string) => PipelineSettings) =>
    (e: ChangeEvent<HTMLInputElement>) => {
      const next = e.target.checked ? "true" : "false";
      onChange(apply(value, next));
    };

  return (
    <div className="space-y-8">
      <section className="space-y-3">
        <h3 className="text-lg font-semibold">Video Processing</h3>
        <div className="grid gap-3 md:grid-cols-2">
          <Field label="Max Images" value={value.videoProcessing.maxNumImages} onChange={(v) => onChange(setByPath(value, "videoProcessing", { ...value.videoProcessing, maxNumImages: v }))} />
          <Field label="Video Start Time" value={value.videoProcessing.videoStartTime} onChange={(v) => onChange(setByPath(value, "videoProcessing", { ...value.videoProcessing, videoStartTime: v }))} />
          <Field label="Video Stop Time" value={value.videoProcessing.videoStopTime} onChange={(v) => onChange(setByPath(value, "videoProcessing", { ...value.videoProcessing, videoStopTime: v }))} />
          <BoolField label="Filter Blurry Images" checked={boolValue(value.videoProcessing.filterBlurryImages)} onChange={setBool((d) => d.videoProcessing.filterBlurryImages, (d, v) => setByPath(d, "videoProcessing", { ...d.videoProcessing, filterBlurryImages: v }))} />
        </div>
      </section>

      <section className="space-y-3">
        <h3 className="text-lg font-semibold">Reconstruction</h3>
        <div className="grid gap-3 md:grid-cols-2">
          <BoolField label="Enable Reconstruction" checked={boolValue(value.reconstruction.enable)} onChange={setBool((d) => d.reconstruction.enable, (d, v) => setByPath(d, "reconstruction", { ...d.reconstruction, enable: v }))} />
          <SelectField label="Reconstruction Software" value={value.reconstruction.softwareName} onChange={(v) => onChange(setByPath(value, "reconstruction", { ...value.reconstruction, softwareName: v as PipelineSettings["reconstruction"]["softwareName"] }))} choices={["colmap", "glomap", "hloc", "map_anything"]} />
          <BoolField label="Enhanced Feature Extraction" checked={boolValue(value.reconstruction.enableEnhancedFeatureExtraction)} onChange={setBool((d) => d.reconstruction.enableEnhancedFeatureExtraction, (d, v) => setByPath(d, "reconstruction", { ...d.reconstruction, enableEnhancedFeatureExtraction: v }))} />
          <SelectField label="Matching Method" value={value.reconstruction.matchingMethod} onChange={(v) => onChange(setByPath(value, "reconstruction", { ...value.reconstruction, matchingMethod: v as PipelineSettings["reconstruction"]["matchingMethod"] }))} choices={["sequential", "exhaustive", "spatial", "vocab", "vocab_tree", "transitive"]} />
          <BoolField label="Enable FL Heuristic" checked={boolValue(value.reconstruction.enableFlHeuristic)} onChange={setBool((d) => d.reconstruction.enableFlHeuristic, (d, v) => setByPath(d, "reconstruction", { ...d.reconstruction, enableFlHeuristic: v }))} />
          <Field label="FL Heuristic Value" value={value.reconstruction.flHeuristicValue} onChange={(v) => onChange(setByPath(value, "reconstruction", { ...value.reconstruction, flHeuristicValue: v }))} />
          <BoolField label="Enable FL Metric" checked={boolValue(value.reconstruction.enableFlMetric)} onChange={setBool((d) => d.reconstruction.enableFlMetric, (d, v) => setByPath(d, "reconstruction", { ...d.reconstruction, enableFlMetric: v }))} />
          <Field label="FL Metric Value" value={value.reconstruction.flMetricValue} onChange={(v) => onChange(setByPath(value, "reconstruction", { ...value.reconstruction, flMetricValue: v }))} />
          <BoolField label="Use Pose Prior COLMAP Model" checked={boolValue(value.reconstruction.posePriors.usePosePriorColmapModelFiles)} onChange={setBool((d) => d.reconstruction.posePriors.usePosePriorColmapModelFiles, (d, v) => setByPath(d, "reconstruction", { ...d.reconstruction, posePriors: { ...d.reconstruction.posePriors, usePosePriorColmapModelFiles: v } }))} />
          <BoolField label="Use Transform JSON" checked={boolValue(value.reconstruction.posePriors.usePosePriorTransformJson.enable)} onChange={setBool((d) => d.reconstruction.posePriors.usePosePriorTransformJson.enable, (d, v) => setByPath(d, "reconstruction", { ...d.reconstruction, posePriors: { ...d.reconstruction.posePriors, usePosePriorTransformJson: { ...d.reconstruction.posePriors.usePosePriorTransformJson, enable: v } } }))} />
          <Field label="Source Coordinate Name" value={value.reconstruction.posePriors.usePosePriorTransformJson.sourceCoordinateName} onChange={(v) => onChange(setByPath(value, "reconstruction", { ...value.reconstruction, posePriors: { ...value.reconstruction.posePriors, usePosePriorTransformJson: { ...value.reconstruction.posePriors.usePosePriorTransformJson, sourceCoordinateName: v } } }))} />
          <BoolField label="Pose Is World To Cam" checked={boolValue(value.reconstruction.posePriors.usePosePriorTransformJson.poseIsWorldToCam)} onChange={setBool((d) => d.reconstruction.posePriors.usePosePriorTransformJson.poseIsWorldToCam, (d, v) => setByPath(d, "reconstruction", { ...d.reconstruction, posePriors: { ...d.reconstruction.posePriors, usePosePriorTransformJson: { ...d.reconstruction.posePriors.usePosePriorTransformJson, poseIsWorldToCam: v } } }))} />
        </div>
      </section>

      <section className="space-y-3">
        <h3 className="text-lg font-semibold">Training</h3>
        <div className="grid gap-3 md:grid-cols-2">
          <BoolField label="Enable Training" checked={boolValue(value.training.enable)} onChange={setBool((d) => d.training.enable, (d, v) => setByPath(d, "training", { ...d.training, enable: v }))} />
          <Field label="Max Steps" value={value.training.maxSteps} onChange={(v) => onChange(setByPath(value, "training", { ...value.training, maxSteps: v }))} />
          <SelectField label="Training Model" value={value.training.model} onChange={(v) => onChange(setByPath(value, "training", { ...value.training, model: v as PipelineSettings["training"]["model"] }))} choices={["splatfacto", "splatfacto-big", "splatfacto-mcmc", "splatfacto-w-light", "3dgut", "3dgrt", "nerfacto"]} />
          <SelectField label="3D ISP" value={value.training["3dIsp"]} onChange={(v) => onChange(setByPath(value, "training", { ...value.training, "3dIsp": v as PipelineSettings["training"]["3dIsp"] }))} choices={["none", "bilagrid", "ppisp"]} />
          <BoolField label="Preserve Scene Scale" checked={boolValue(value.training.preserveSceneScale)} onChange={setBool((d) => d.training.preserveSceneScale, (d, v) => setByPath(d, "training", { ...d.training, preserveSceneScale: v }))} />
          <BoolField label="Enable Depth Loss" checked={boolValue(value.training.enableDepthLoss)} onChange={setBool((d) => d.training.enableDepthLoss, (d, v) => setByPath(d, "training", { ...d.training, enableDepthLoss: v }))} />
        </div>
      </section>

      <section className="space-y-3">
        <h3 className="text-lg font-semibold">Post Processing</h3>
        <div className="grid gap-3 md:grid-cols-2">
          <BoolField label="Crop Output Bounds" checked={boolValue(value.postProcessing.cropOutputBounds)} onChange={setBool((d) => d.postProcessing.cropOutputBounds, (d, v) => setByPath(d, "postProcessing", { ...d.postProcessing, cropOutputBounds: v }))} />
          <SelectField label="Crop Mode" value={value.postProcessing.cropMode} onChange={(v) => onChange(setByPath(value, "postProcessing", { ...value.postProcessing, cropMode: v as PipelineSettings["postProcessing"]["cropMode"] }))} choices={["environment", "rigid_body"]} />
          <BoolField label="Clean Splat" checked={boolValue(value.postProcessing.cleanSplat)} onChange={setBool((d) => d.postProcessing.cleanSplat, (d, v) => setByPath(d, "postProcessing", { ...d.postProcessing, cleanSplat: v }))} />
          <BoolField label="Enable SPZ" checked={boolValue(value.postProcessing.enableSpz)} onChange={setBool((d) => d.postProcessing.enableSpz, (d, v) => setByPath(d, "postProcessing", { ...d.postProcessing, enableSpz: v }))} />
          <BoolField label="Enable SOG" checked={boolValue(value.postProcessing.enableSog)} onChange={setBool((d) => d.postProcessing.enableSog, (d, v) => setByPath(d, "postProcessing", { ...d.postProcessing, enableSog: v }))} />
          <BoolField label="Enable USDZ" checked={boolValue(value.postProcessing.enableUsdz)} onChange={setBool((d) => d.postProcessing.enableUsdz, (d, v) => setByPath(d, "postProcessing", { ...d.postProcessing, enableUsdz: v }))} />
          <BoolField label="Enable Video Export" checked={boolValue(value.postProcessing.enableVideoExport)} onChange={setBool((d) => d.postProcessing.enableVideoExport, (d, v) => setByPath(d, "postProcessing", { ...d.postProcessing, enableVideoExport: v }))} />
          <SelectField label="PLY Coords" value={value.postProcessing.plyCoords} onChange={(v) => onChange(setByPath(value, "postProcessing", { ...value.postProcessing, plyCoords: v as PipelineSettings["postProcessing"]["plyCoords"] }))} choices={coordChoices.map((x) => x.value)} labels={Object.fromEntries(coordChoices.map((x) => [x.value, x.label]))} />
          <SelectField label="SPZ Coords" value={value.postProcessing.spzCoords} onChange={(v) => onChange(setByPath(value, "postProcessing", { ...value.postProcessing, spzCoords: v as PipelineSettings["postProcessing"]["spzCoords"] }))} choices={coordChoices.map((x) => x.value)} labels={Object.fromEntries(coordChoices.map((x) => [x.value, x.label]))} />
          <SelectField label="SOG Coords" value={value.postProcessing.sogCoords} onChange={(v) => onChange(setByPath(value, "postProcessing", { ...value.postProcessing, sogCoords: v as PipelineSettings["postProcessing"]["sogCoords"] }))} choices={coordChoices.map((x) => x.value)} labels={Object.fromEntries(coordChoices.map((x) => [x.value, x.label]))} />
          <SelectField label="USDZ Coords" value={value.postProcessing.usdzCoords} onChange={(v) => onChange(setByPath(value, "postProcessing", { ...value.postProcessing, usdzCoords: v as PipelineSettings["postProcessing"]["usdzCoords"] }))} choices={coordChoices.map((x) => x.value)} labels={Object.fromEntries(coordChoices.map((x) => [x.value, x.label]))} />
        </div>
      </section>

      <section className="space-y-3">
        <h3 className="text-lg font-semibold">Spherical Camera</h3>
        <div className="grid gap-3 md:grid-cols-2">
          <BoolField label="Enable Spherical Camera" checked={boolValue(value.sphericalCamera.enable)} onChange={setBool((d) => d.sphericalCamera.enable, (d, v) => setByPath(d, "sphericalCamera", { ...d.sphericalCamera, enable: v }))} />
          <Field label="Cube Faces To Remove" value={value.sphericalCamera.cubeFacesToRemove} onChange={(v) => onChange(setByPath(value, "sphericalCamera", { ...value.sphericalCamera, cubeFacesToRemove: v }))} />
          <BoolField label="Optimize Sequential Frame Order" checked={boolValue(value.sphericalCamera.optimizeSequentialFrameOrder)} onChange={setBool((d) => d.sphericalCamera.optimizeSequentialFrameOrder, (d, v) => setByPath(d, "sphericalCamera", { ...d.sphericalCamera, optimizeSequentialFrameOrder: v }))} />
        </div>
      </section>

      <section className="space-y-3">
        <h3 className="text-lg font-semibold">Segmentation</h3>
        <div className="grid gap-3 md:grid-cols-2">
          <BoolField label="Enable Background Removal" checked={boolValue(value.segmentation.backgroundRemoval.enable)} onChange={setBool((d) => d.segmentation.backgroundRemoval.enable, (d, v) => setByPath(d, "segmentation", { ...d.segmentation, backgroundRemoval: { ...d.segmentation.backgroundRemoval, enable: v } }))} />
          <SelectField label="Background Model" value={value.segmentation.backgroundRemoval.model} onChange={(v) => onChange(setByPath(value, "segmentation", { ...value.segmentation, backgroundRemoval: { ...value.segmentation.backgroundRemoval, model: v as PipelineSettings["segmentation"]["backgroundRemoval"]["model"] } }))} choices={["u2net", "sam2"]} />
          <Field label="Mask Threshold" value={value.segmentation.backgroundRemoval.maskThreshold} onChange={(v) => onChange(setByPath(value, "segmentation", { ...value.segmentation, backgroundRemoval: { ...value.segmentation.backgroundRemoval, maskThreshold: v } }))} />
          <BoolField label="Enable Object Removal" checked={boolValue(value.segmentation.objectRemoval.enable)} onChange={setBool((d) => d.segmentation.objectRemoval.enable, (d, v) => setByPath(d, "segmentation", { ...d.segmentation, objectRemoval: { ...d.segmentation.objectRemoval, enable: v } }))} />
          <SelectField label="Object Removal Action" value={value.segmentation.objectRemoval.action} onChange={(v) => onChange(setByPath(value, "segmentation", { ...value.segmentation, objectRemoval: { ...value.segmentation.objectRemoval, action: v as PipelineSettings["segmentation"]["objectRemoval"]["action"] } }))} choices={["erase", "remove"]} />
          <Field label="Objects" value={value.segmentation.objectRemoval.objects} onChange={(v) => onChange(setByPath(value, "segmentation", { ...value.segmentation, objectRemoval: { ...value.segmentation.objectRemoval, objects: v } }))} />
        </div>
      </section>
    </div>
  );
}

function Field({ label, value, onChange }: { label: string; value: string; onChange: (v: string) => void }) {
  return (
    <label className="space-y-1 text-sm">
      <span className="font-medium">{label}</span>
      <Input value={value} onChange={(e) => onChange(e.target.value)} />
    </label>
  );
}

function BoolField({ label, checked, onChange }: { label: string; checked: boolean; onChange: (e: ChangeEvent<HTMLInputElement>) => void }) {
  return (
    <label className="flex items-center gap-2 rounded-md border p-3 text-sm">
      <input type="checkbox" className="h-4 w-4" checked={checked} onChange={onChange} />
      <span>{label}</span>
    </label>
  );
}

function SelectField({
  label,
  value,
  onChange,
  choices,
  labels,
}: {
  label: string;
  value: string;
  onChange: (v: string) => void;
  choices: readonly string[];
  labels?: Record<string, string>;
}) {
  return (
    <label className="space-y-1 text-sm">
      <span className="font-medium">{label}</span>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="flex h-10 w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm"
      >
        {choices.map((choice) => (
          <option key={choice} value={choice}>
            {labels?.[choice] ?? choice}
          </option>
        ))}
      </select>
    </label>
  );
}
