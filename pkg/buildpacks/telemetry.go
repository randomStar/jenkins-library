package buildpacks

import (
	"encoding/json"

	"github.com/SAP/jenkins-library/pkg/cnbutils/privacy"
	"github.com/SAP/jenkins-library/pkg/cnbutils/project"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/pkg/errors"
)

const version = 3

type BuildpacksTelemetry struct {
	customData *telemetry.CustomData
	data       cnbBuildTelemetry
}

func NewBuildpacksTelemetry(customData *telemetry.CustomData) BuildpacksTelemetry {
	return BuildpacksTelemetry{
		customData: customData,
		data: cnbBuildTelemetry{
			Version: version,
		},
	}
}

func (d BuildpacksTelemetry) Export() error {
	d.customData.Custom1Label = "cnbBuildStepData"
	customData, err := json.Marshal(d.data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal custom telemetry data")
	}
	d.customData.Custom1 = string(customData)
	return nil
}

func (d BuildpacksTelemetry) Data() []cnbBuildTelemetryData {
	return d.data.Data
}

func (d BuildpacksTelemetry) Version() int {
	return d.data.Version
}

func (d BuildpacksTelemetry) WithImage(image string) {
	d.data.builder = image
}

func (d BuildpacksTelemetry) AddSegment(segment Segment) {
	segment.data.Builder = d.data.builder
	d.data.Data = append(d.data.Data, segment.data)
}

type cnbBuildTelemetry struct {
	builder string
	Version int                     `json:"version"`
	Data    []cnbBuildTelemetryData `json:"data"`
}

type cnbBuildTelemetryData struct {
	ImageTag          string                                 `json:"imageTag"`
	AdditionalTags    []string                               `json:"additionalTags"`
	BindingKeys       []string                               `json:"bindingKeys"`
	Path              PathEnum                               `json:"path"`
	BuildEnv          cnbBuildTelemetryDataBuildEnv          `json:"buildEnv"`
	Buildpacks        cnbBuildTelemetryDataBuildpacks        `json:"buildpacks"`
	ProjectDescriptor cnbBuildTelemetryDataProjectDescriptor `json:"projectDescriptor"`
	BuildTool         string                                 `json:"buildTool"`
	Builder           string                                 `json:"builder"`
}

type cnbBuildTelemetryDataBuildEnv struct {
	KeysFromConfig            []string               `json:"keysFromConfig"`
	KeysFromProjectDescriptor []string               `json:"keysFromProjectDescriptor"`
	KeysOverall               []string               `json:"keysOverall"`
	JVMVersion                string                 `json:"jvmVersion"`
	KeyValues                 map[string]interface{} `json:"keyValues"`
}

type cnbBuildTelemetryDataBuildpacks struct {
	FromConfig            []string `json:"FromConfig"`
	FromProjectDescriptor []string `json:"FromProjectDescriptor"`
	Overall               []string `json:"overall"`
}

type cnbBuildTelemetryDataProjectDescriptor struct {
	Used        bool `json:"used"`
	IncludeUsed bool `json:"includeUsed"`
	ExcludeUsed bool `json:"excludeUsed"`
}

type Segment struct {
	data cnbBuildTelemetryData
}

func NewSegment() Segment {
	return Segment{
		data: cnbBuildTelemetryData{},
	}
}

func (s Segment) WithBindings(bindings map[string]interface{}) {
	var bindingKeys []string
	for k := range bindings {
		bindingKeys = append(bindingKeys, k)
	}
	s.data.BindingKeys = bindingKeys
}

func (s Segment) WithEnv(env map[string]interface{}) {
	s.data.BuildEnv.KeysFromConfig = []string{}
	s.data.BuildEnv.KeysOverall = []string{}
	for key := range env {
		s.data.BuildEnv.KeysFromConfig = append(s.data.BuildEnv.KeysFromConfig, key)
		s.data.BuildEnv.KeysOverall = append(s.data.BuildEnv.KeysOverall, key)
	}
}

// Merge tags?
func (s Segment) WithTags(tag string, additionalTags []string) {
	s.data.ImageTag = tag
	s.data.AdditionalTags = additionalTags
}

func (s Segment) WithPath(path PathEnum) {
	s.data.Path = path
}

func (s Segment) WithBuildTool(buildTool string) {
	s.data.BuildTool = buildTool
}

func (s Segment) WithBuilder(builder string) {
	s.data.Builder = privacy.FilterBuilder(builder)
}

func (s Segment) WithBuildpacksFromConfig(buildpacks []string) {
	s.data.Buildpacks.FromConfig = privacy.FilterBuildpacks(buildpacks)
}

func (s Segment) WithBuildpacksOverall(buildpacks []string) {
	s.data.Buildpacks.Overall = privacy.FilterBuildpacks(buildpacks)
}

func (s Segment) WithKeyValues(env map[string]interface{}) {
	s.data.BuildEnv.KeyValues = privacy.FilterEnv(env)
}

func (s Segment) WithProjectDescriptor(descriptor *project.Descriptor) {
	descriptorKeys := s.data.BuildEnv.KeysFromProjectDescriptor
	overallKeys := s.data.BuildEnv.KeysOverall
	for key := range descriptor.EnvVars {
		descriptorKeys = append(descriptorKeys, key)
		overallKeys = append(overallKeys, key)
	}
	s.data.BuildEnv.KeysFromProjectDescriptor = descriptorKeys
	s.data.BuildEnv.KeysOverall = overallKeys
	s.data.Buildpacks.FromProjectDescriptor = privacy.FilterBuildpacks(descriptor.Buildpacks)
	s.data.ProjectDescriptor.Used = true
	s.data.ProjectDescriptor.IncludeUsed = descriptor.Include != nil
	s.data.ProjectDescriptor.ExcludeUsed = descriptor.Exclude != nil
}
