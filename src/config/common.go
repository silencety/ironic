package config

type CommonConfig struct{
	// init from config file
	CorsHeaders          string                   `json:"api-cors-headers,omitempty"`
	EnableCors           bool                     `json:"api-enable-cors,omitempty"`
	EnableSelinuxSupport bool                     `json:"selinux-enabled,omitempty"`
	RemappedRoot         string                   `json:"userns-remap,omitempty"`
	SocketGroup          string                   `json:"group,omitempty"`
	CgroupParent         string                   `json:"cgroup-parent,omitempty"`
	MaxTimeOut	     int	              `json:"MaxTimeOut,omitempty"`
}

type Addr struct {
	Proto string
	Addr  string
}

const SectionDefault  = "default"
const SectionDatabase = "database"
const defaultPidFile = "/var/run/ironic.pid"
