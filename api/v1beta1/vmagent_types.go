package v1beta1

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// VMAgentSpec defines the desired state of VMAgent
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version",description="The version of VMAgent"
// +kubebuilder:printcolumn:name="ReplicaCount",type="integer",JSONPath=".spec.replicas",description="The desired replicas number of VMAgent"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type VMAgentSpec struct {
	// PodMetadata configures Labels and Annotations which are propagated to the vmagent pods.
	// +optional
	PodMetadata *EmbeddedObjectMetadata `json:"podMetadata,omitempty"`
	// Image - docker image settings for VMAgent
	// if no specified operator uses default config version
	// +optional
	Image Image `json:"image,omitempty"`
	// ConfigReloadImage - docker image settings for prometheus-config-reloader
	// if no specified operator uses default one
	// +optional
	ConfigReloadImage Image `json:"configReloadImage,omitempty"`
	// ImagePullSecrets An optional list of references to secrets in the same namespace
	// to use for pulling images from registries
	// see http://kubernetes.io/docs/user-guide/images#specifying-imagepullsecrets-on-a-pod
	// +optional
	ImagePullSecrets []v1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// Secrets is a list of Secrets in the same namespace as the vmagent
	// object, which shall be mounted into the vmagent Pods.
	// will be mounted at path /etc/vm/secrets
	// +optional
	Secrets []string `json:"secrets,omitempty"`
	// ConfigMaps is a list of ConfigMaps in the same namespace as the vmagent
	// object, which shall be mounted into the vmagent Pods.
	// will be mounted at path  /etc/vm/configs
	// +optional
	ConfigMaps []string `json:"configMaps,omitempty"`
	// LogLevel for VMAgent to be configured with.
	//INFO, WARN, ERROR, FATAL, PANIC
	// +optional
	// +kubebuilder:validation:Enum=INFO;WARN;ERROR;FATAL;PANIC
	LogLevel string `json:"logLevel,omitempty"`
	// LogFormat for VMAgent to be configured with.
	// +optional
	// +kubebuilder:validation:Enum=default;json
	LogFormat string `json:"logFormat,omitempty"`
	// ReplicaCount is the expected size of the VMAgent cluster. The controller will
	// eventually make the size of the running cluster equal to the expected
	// size.
	// NOTE enable VMSingle deduplication for replica usage
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Number of pods",xDescriptors="urn:alm:descriptor:com.tectonic.ui:podCount,urn:alm:descriptor:io.kubernetes:custom"
	// +optional
	ReplicaCount *int32 `json:"replicaCount,omitempty"`
	// Volumes allows configuration of additional volumes on the output deploy definition.
	// Volumes specified will be appended to other volumes that are generated as a result of
	// StorageSpec objects.
	// +optional
	Volumes []v1.Volume `json:"volumes,omitempty"`
	// VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition.
	// VolumeMounts specified will be appended to other VolumeMounts in the vmagent container,
	// that are generated as a result of StorageSpec objects.
	// +optional
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`
	// Resources container resource request and limits, https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	//if not specified - default setting will be used
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Resources",xDescriptors="urn:alm:descriptor:com.tectonic.ui:resourceRequirements"
	// +optional
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// Affinity If specified, the pod's scheduling constraints.
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Tolerations If specified, the pod's tolerations.
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// SecurityContext holds pod-level security attributes and common container settings.
	// This defaults to the default PodSecurityContext.
	// +optional
	SecurityContext *v1.PodSecurityContext `json:"securityContext,omitempty"`
	// ServiceAccountName is the name of the ServiceAccount to use to run the
	// VMAgent Pods.
	// required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="ServiceAccount name",xDescriptors="urn:alm:descriptor:io.kubernetes:ServiceAccount"
	ServiceAccountName string `json:"serviceAccountName"`
	// Containers property allows to inject additions sidecars. It can be useful for proxies, backup, etc.
	// +optional
	Containers []v1.Container `json:"containers,omitempty"`
	// InitContainers allows adding initContainers to the pod definition. Those can be used to e.g.
	// fetch secrets for injection into the vmagent configuration from external sources. Any
	// errors during the execution of an initContainer will lead to a restart of the Pod. More info: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/
	// Using initContainers for any use case other then secret fetching is entirely outside the scope
	// of what the maintainers will support and by doing so, you accept that this behaviour may break
	// at any time without notice.
	// +optional
	InitContainers []v1.Container `json:"initContainers,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// HostNetwork controls whether the pod may use the node network namespace
	// +optional
	HostNetwork bool `json:"hostNetwork,omitempty"`
	// DNSPolicy set DNS policy for the pod
	// +optional
	DNSPolicy v1.DNSPolicy `json:"dnsPolicy,omitempty"`
	// ScrapeInterval defines how often scrape targets by default
	// +optional
	// +kubebuilder:validation:Pattern:="[0-9]+(ms|s|m|h)"
	ScrapeInterval string `json:"scrapeInterval,omitempty"`

	// APIServerConfig allows specifying a host and auth methods to access apiserver.
	// If left empty, VMAgent is assumed to run inside of the cluster
	// and will discover API servers automatically and use the pod's CA certificate
	// and bearer token file at /var/run/secrets/kubernetes.io/serviceaccount/.
	// +optional
	APIServerConfig *APIServerConfig `json:"aPIServerConfig,omitempty"`
	// OverrideHonorLabels if set to true overrides all user configured honor_labels.
	// If HonorLabels is set in ServiceScrape or PodScrape to true, this overrides honor_labels to false.
	// +optional
	OverrideHonorLabels bool `json:"overrideHonorLabels,omitempty"`
	// OverrideHonorTimestamps allows to globally enforce honoring timestamps in all scrape configs.
	// +optional
	OverrideHonorTimestamps bool `json:"overrideHonorTimestamps,omitempty"`
	// IgnoreNamespaceSelectors if set to true will ignore NamespaceSelector settings from
	// the podscrape and vmservicescrape configs, and they will only discover endpoints
	// within their current namespace.  Defaults to false.
	// +optional
	IgnoreNamespaceSelectors bool `json:"ignoreNamespaceSelectors,omitempty"`
	// EnforcedNamespaceLabel enforces adding a namespace label of origin for each alert
	// and metric that is user created. The label value will always be the namespace of the object that is
	// being created.
	// +optional
	EnforcedNamespaceLabel string `json:"enforcedNamespaceLabel,omitempty"`
	// VMAgentExternalLabelName Name of vmAgent external label used to denote vmAgent instance
	// name. Defaults to the value of `prometheus`. External label will
	// _not_ be added when value is set to empty string (`""`).
	// +optional
	VMAgentExternalLabelName *string `json:"vmAgentExternalLabelName,omitempty"`

	// ExternalLabels The labels to add to any time series or alerts when communicating with
	// external systems (federation, remote storage, etc).
	// +optional
	ExternalLabels map[string]string `json:"externalLabels,omitempty"`
	// RemoteWrite list of victoria metrics /some other remote write system
	// for vm it must looks like: http://victoria-metrics-single:8429/api/v1/write
	// or for cluster different url
	// https://github.com/VictoriaMetrics/VictoriaMetrics/tree/master/app/vmagent#splitting-data-streams-among-multiple-systems
	// +optional
	RemoteWrite []VMAgentRemoteWriteSpec `json:"remoteWrite"`
	// RelabelConfig ConfigMap with global relabel config -remoteWrite.relabelConfig
	// This relabeling is applied to all the collected metrics before sending them to remote storage.
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Key at Configmap with relabelConfig name",xDescriptors="urn:alm:descriptor:io.kubernetes:ConfigMapKeySelector"
	RelabelConfig *v1.ConfigMapKeySelector `json:"relabelConfig,omitempty"`
	// ServiceScrapeSelector defines ServiceScrapes to be selected for target discovery. if
	// neither serviceScrapeNamespaceSelector nor ProbeSelector nor serviceScrapeSelector are specified, configuration is
	// unmanaged.
	// +optional
	ServiceScrapeSelector *metav1.LabelSelector `json:"serviceScrapeSelector,omitempty"`
	// ServiceScrapeNamespaceSelector Namespaces to be selected for ServiceMonitor discovery. If nil, only
	// check own namespace.
	// +optional
	ServiceScrapeNamespaceSelector *metav1.LabelSelector `json:"serviceScrapeNamespaceSelector,omitempty"`
	// PodScrapeSelector defines PodScrapes to be selected for target discovery.
	// if neither PodScrapeNamespaceSelector nor ProbeSelector nor PodScrapeSelector are specified,
	// configuration is unmanaged.
	// +optional
	PodScrapeSelector *metav1.LabelSelector `json:"podScrapeSelector,omitempty"`
	// PodScrapeNamespaceSelector defines Namespaces to be selected for PodMonitor discovery. If nil, only
	// check own namespace.
	// +optional
	PodScrapeNamespaceSelector *metav1.LabelSelector `json:"podScrapeNamespaceSelector,omitempty"`
	// ProbeSelector defines VMProbe to be selected for target probing.
	// if neither PodScrapeNamespaceSelector, nor ProbeSelector nor PodScrapeSelector are specified,
	// configuration is unmanaged.
	// +optional
	ProbeSelector *metav1.LabelSelector `json:"probeSelector,omitempty"`
	// ProbeNamespaceSelector defines Namespaces to be selected for VMProbe discovery. If nil, only
	// check own namespace.
	// +optional
	ProbeNamespaceSelector *metav1.LabelSelector `json:"probeNamespaceSelector,omitempty"`

	// AdditionalScrapeConfigs As scrape configs are appended, the user is responsible to make sure it
	// is valid. Note that using this feature may expose the possibility to
	// break upgrades of VMAgent. It is advised to review VMAgent release
	// notes to ensure that no incompatible scrape configs are going to break
	// VMAgent after the upgrade.
	// +optional
	AdditionalScrapeConfigs *v1.SecretKeySelector `json:"additionalScrapeConfigs,omitempty"`
	// ArbitraryFSAccessThroughSMs configures whether configuration
	// based on a service scrape can access arbitrary files on the file system
	// of the VMAgent container e.g. bearer token files.
	// +optional
	ArbitraryFSAccessThroughSMs ArbitraryFSAccessThroughSMsConfig `json:"arbitraryFSAccessThroughSMs,omitempty"`
	// Port listen address
	// +optional
	Port string `json:"port,omitempty"`
	// ExtraArgs that will be passed to  VMAgent pod
	// for example remoteWrite.tmpDataPath: /tmp
	// it would be converted to flag --remoteWrite.tmpDataPath=/tmp
	// +optional
	ExtraArgs map[string]string `json:"extraArgs,omitempty"`
	// ExtraEnvs that will be added to VMAgent pod
	// +optional
	ExtraEnvs []v1.EnvVar `json:"extraEnvs,omitempty"`
}

// VMAgentRemoteWriteSpec defines the remote storage configuration for VmAgent
// +k8s:openapi-gen=true
type VMAgentRemoteWriteSpec struct {
	// URL of the endpoint to send samples to.
	URL string `json:"url"`
	// BasicAuth allow an endpoint to authenticate over basic authentication
	// +optional
	BasicAuth *BasicAuth `json:"basicAuth,omitempty"`
	// Optional bearer auth token to use for -remoteWrite.url
	// +optional
	BearerTokenSecret *v1.SecretKeySelector `json:"bearerTokenSecret,omitempty"`
	// Interval for flushing the data to remote storage. (default 1s)
	// +optional
	// +kubebuilder:validation:Pattern:="[0-9]+(ms|s|m|h)"
	FlushInterval *string `json:"flushInterval,omitempty"`
	// Optional labels in the form 'name=value' to add to all the metrics before sending them
	// +optional
	Labels map[string]string `json:"label,omitempty"`
	// The maximum size in bytes of unpacked request to send to remote storage
	// +optional
	MaxBlockSize *int32 `json:"maxBlockSize,omitempty"`
	// The maximum file-based buffer size in bytes at -remoteWrite.tmpDataPath
	// +optional
	MaxDiskUsagePerURL *int32 `json:"maxDiskUsagePerURL,omitempty"`
	// The number of concurrent queues
	// +optional
	Queues *int32 `json:"queues,omitempty"`
	// ConfigMap with relabeling config which is applied to metrics before sending them to the corresponding -remoteWrite.url
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Key at Configmap with relabelConfig for remoteWrite",xDescriptors="urn:alm:descriptor:io.kubernetes:ConfigMapKeySelector"
	UrlRelabelConfig *v1.ConfigMapKeySelector `json:"urlRelabelConfig,omitempty"`
	// Timeout for sending a single block of data to -remoteWrite.url (default 1m0s)
	// +optional
	// +kubebuilder:validation:Pattern:="[0-9]+(ms|s|m|h)"
	SendTimeout *string `json:"sendTimeout,omitempty"`
	// Whether to show -remoteWrite.url in the exported metrics. It is hidden by default, since it can contain sensistive auth info
	// +optional
	ShowURL *bool `json:"showURL,omitempty"`
	// Path to directory where temporary data for remote write component is stored (default "vmagent-remotewrite-data")
	// +optional
	TmpDataPath *string `json:"tmpDataPath,omitempty"`
	// TLSConfig describes tls configuration for remote write target
	TLSConfig *TLSConfig `json:"tlsConfig,omitempty"`
}

// VmAgentStatus defines the observed state of VmAgent
// +k8s:openapi-gen=true
type VMAgentStatus struct {
	// ReplicaCount Total number of non-terminated pods targeted by this VMAlert
	// cluster (their labels match the selector).
	Replicas int32 `json:"replicas"`
	// UpdatedReplicas Total number of non-terminated pods targeted by this VMAlert
	// cluster that have the desired version spec.
	UpdatedReplicas int32 `json:"updatedReplicas"`
	// AvailableReplicas Total number of available pods (ready for at least minReadySeconds)
	// targeted by this VMAlert cluster.
	AvailableReplicas int32 `json:"availableReplicas"`
	// UnavailableReplicas Total number of unavailable pods targeted by this VMAlert cluster.
	UnavailableReplicas int32 `json:"unavailableReplicas"`
}

// VMAgent - is a tiny but brave agent, which helps you collect metrics from various sources and stores them in VictoriaMetrics
// or any other Prometheus-compatible storage system that supports the remote_write protocol.
// +operator-sdk:gen-csv:customresourcedefinitions.displayName="VMAgent App"
// +operator-sdk:gen-csv:customresourcedefinitions.resources="Deployment,apps"
// +operator-sdk:gen-csv:customresourcedefinitions.resources="Service,v1"
// +operator-sdk:gen-csv:customresourcedefinitions.resources="Secret,v1"
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=vmagents,scope=Namespaced
type VMAgent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VMAgentSpec   `json:"spec,omitempty"`
	Status VMAgentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// VMAgentList contains a list of VMAgent
type VMAgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VMAgent `json:"items"`
}

func (cr *VMAgent) AsOwner() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion:         cr.APIVersion,
			Kind:               cr.Kind,
			Name:               cr.Name,
			UID:                cr.UID,
			Controller:         pointer.BoolPtr(true),
			BlockOwnerDeletion: pointer.BoolPtr(true),
		},
	}
}

func (cr VMAgent) PodAnnotations() map[string]string {
	annotations := map[string]string{}
	if cr.Spec.PodMetadata != nil {
		for annotation, value := range cr.Spec.PodMetadata.Annotations {
			annotations[annotation] = value
		}
	}
	return annotations
}

func (cr VMAgent) Annotations() map[string]string {
	annotations := make(map[string]string)
	for annotation, value := range cr.ObjectMeta.Annotations {
		if !strings.HasPrefix(annotation, "kubectl.kubernetes.io/") {
			annotations[annotation] = value
		}
	}
	return annotations
}

func (cr VMAgent) SelectorLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":      "vmagent",
		"app.kubernetes.io/instance":  cr.Name,
		"app.kubernetes.io/component": "monitoring",
		"managed-by":                  "vm-operator",
	}
}

func (cr VMAgent) PodLabels() map[string]string {
	labels := cr.SelectorLabels()
	if cr.Spec.PodMetadata != nil {
		for label, value := range cr.Spec.PodMetadata.Labels {
			labels[label] = value
		}
	}
	return labels
}

func (cr VMAgent) FinalLabels() map[string]string {
	labels := cr.SelectorLabels()
	if cr.ObjectMeta.Labels != nil {
		for label, value := range cr.ObjectMeta.Labels {
			labels[label] = value
		}
	}
	return labels
}

func (cr VMAgent) PrefixedName() string {
	return fmt.Sprintf("vmagent-%s", cr.Name)
}

func (cr VMAgent) TLSAssetName() string {
	return fmt.Sprintf("tls-assets-vmagent-%s", cr.Name)
}

func (cr VMAgent) HealthPath() string {
	return buildPathWithPrefixFlag(cr.Spec.ExtraArgs, healthPath)
}

func (cr VMAgent) MetricPath() string {
	return buildPathWithPrefixFlag(cr.Spec.ExtraArgs, metricPath)
}

func (cr VMAgent) ReloadPathWithPort(port string) string {
	return fmt.Sprintf("http://localhost:%s%s", port, buildPathWithPrefixFlag(cr.Spec.ExtraArgs, reloadPath))
}

func init() {
	SchemeBuilder.Register(&VMAgent{}, &VMAgentList{})
}
