package proxy

import (
	"github.com/jkkgbe/open-zcash-pool/policy"
	"github.com/jkkgbe/open-zcash-pool/storage"
)

type Config struct {
	Name                  string     `json:"name"`
	Proxy                 Proxy      `json:"proxy"`
	Upstream              []Upstream `json:"upstream"`
	UpstreamCheckInterval string     `json:"upstreamCheckInterval"`

	Threads int `json:"threads"`

	InstanceId uint32 `json:"instanceId"`

	Coin  string         `json:"coin"`
	Redis storage.Config `json:"redis"`

	NewrelicName    string `json:"newrelicName"`
	NewrelicKey     string `json:"newrelicKey"`
	NewrelicVerbose bool   `json:"newrelicVerbose"`
	NewrelicEnabled bool   `json:"newrelicEnabled"`
}

type Proxy struct {
	Enabled              bool   `json:"enabled"`
	Listen               string `json:"listen"`
	LimitHeadersSize     int    `json:"limitHeadersSize"`
	LimitBodySize        int64  `json:"limitBodySize"`
	BehindReverseProxy   bool   `json:"behindReverseProxy"`
	BlockRefreshInterval string `json:"blockRefreshInterval"`
	Difficulty           int64  `json:"difficulty"`
	StateUpdateInterval  string `json:"stateUpdateInterval"`
	HashrateExpiration   string `json:"hashrateExpiration"`

	Policy policy.Config `json:"policy"`

	MaxFails    int64 `json:"maxFails"`
	HealthCheck bool  `json:"healthCheck"`

	Stratum Stratum `json:"stratum"`
}

type Stratum struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
	Timeout string `json:"timeout"`
	MaxConn int    `json:"maxConn"`
}

type Upstream struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	Timeout string `json:"timeout"`
}
