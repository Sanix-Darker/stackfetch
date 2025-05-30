package services

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// DefaultPorts maps service names (or stack dependency keys) to their common TCP ports.
var DefaultPorts = map[string][]int{
	// Web Servers
	"apache":   {80, 443, 8080, 8000, 8443},
	"nginx":    {80, 443, 8080, 8000, 8443},
	"caddy":    {80, 443, 2019},
	"lighttpd": {80, 443},
	"iis":      {80, 443, 8080},
	"tomcat":   {8080, 8005, 8009, 8443},
	"jetty":    {8080, 8443},
	"gunicorn": {8000, 8080},
	"uwsgi":    {8000, 8080},

	// Databases
	"mysql":       {3306, 33060},
	"mariadb":     {3306, 33060},
	"postgresql":  {5432},
	"mongodb":     {27017, 27018, 27019, 28017},
	"redis":       {6379, 26379},
	"memcached":   {11211},
	"cassandra":   {9042, 9160, 7000, 7001},
	"oracle":      {1521, 5500},
	"sqlserver":   {1433, 1434},
	"clickhouse":  {8123, 9000, 9009},
	"couchdb":     {5984, 5986},
	"riak":        {8087, 8098},
	"neo4j":       {7474, 7687},
	"arangodb":    {8529},
	"influxdb":    {8086, 8088},
	"scylladb":    {9042, 9160, 7000, 7001},
	"cockroachdb": {26257, 8080},
	"dynamodb":    {8000},
	"firebird":    {3050},
	"db2":         {50000},

	// Message Brokers
	"rabbitmq":  {5672, 5671, 15672, 25672},
	"kafka":     {9092, 9093, 9094, 2181},
	"activemq":  {61616, 8161, 5672},
	"nats":      {4222, 6222, 8222},
	"zeromq":    {5555, 5556, 5557},
	"pulsar":    {6650, 8080},
	"mosquitto": {1883, 8883, 9001},
	"emqx":      {1883, 8883, 8083, 8084, 18083},

	// Search & Analytics
	"elasticsearch": {9200, 9300},
	"solr":          {8983, 8984},
	"splunk":        {8000, 8088, 8089, 9997},
	"graylog":       {9000, 12201, 12202, 514},
	"logstash":      {5044, 9600},
	"kibana":        {5601},
	"grafana":       {3000},
	"prometheus":    {9090},
	"loki":          {3100},
	"tempo":         {3200, 4317, 4318},
	"jaeger":        {14250, 14268, 16686},
	"zipkin":        {9411},

	// Infrastructure & Monitoring
	"docker":     {2375, 2376},
	"kubernetes": {6443, 10250, 10255, 10256},
	"consul":     {8300, 8301, 8302, 8500, 8600},
	"etcd":       {2379, 2380},
	"vault":      {8200},
	"nomad":      {4646, 4647, 4648},
	"zookeeper":  {2181, 2888, 3888},
	"hadoop":     {8020, 9000, 50070, 50075, 50090},
	"spark":      {4040, 7077, 8080, 8081},
	"flink":      {6123, 8081, 8082},
	"airflow":    {8080},

	// Network Services
	"ssh":      {22},
	"ftp":      {21, 20},
	"sftp":     {22},
	"ftps":     {990, 989},
	"telnet":   {23},
	"dns":      {53, 853},
	"dhcp":     {67, 68},
	"ntp":      {123},
	"ldap":     {389, 636},
	"radius":   {1812, 1813},
	"kerberos": {88, 464, 749},
	"smb":      {139, 445},
	"nfs":      {2049, 111},
	"iscsi":    {3260},
	"snmp":     {161, 162},
	"syslog":   {514},
	"tftp":     {69},
	"sip":      {5060, 5061},
	"rtp":      {5004, 5005},
	"rtsp":     {554, 8554},
	"xmpp":     {5222, 5269},

	// Email
	"smtp":     {25, 465, 587},
	"imap":     {143, 993},
	"pop3":     {110, 995},
	"exchange": {25, 80, 443, 587, 993, 995},
	"postfix":  {25, 465, 587},
	"dovecot":  {110, 143, 993, 995},
	"sendmail": {25, 587},
	"exim":     {25, 465, 587},

	// Version Control
	"git": {9418},
	"svn": {3690},
	"hg":  {8000},

	// CI/CD
	"jenkins":  {8080, 50000},
	"teamcity": {8111},
	"bamboo":   {8085, 54663},
	"gitlab":   {80, 443, 22, 2222},
	"gitea":    {3000, 2222},
	"drone":    {80, 443, 8000},
	"argo":     {2746, 8080},
	"tekton":   {9090},

	// Security
	"vpn":        {1194, 1701, 1723, 4500, 500},
	"openvpn":    {1194},
	"wireguard":  {51820},
	"tor":        {9001, 9030, 9050, 9051},
	"nessus":     {8834},
	"nexpose":    {3780},
	"burp":       {8080, 8443},
	"metasploit": {3790, 55553},
	"zap":        {8080, 8090},

	// Media
	"rtmp":     {1935},
	"icecast":  {8000},
	"plex":     {32400},
	"emby":     {8096, 8920},
	"jellyfin": {8096, 8920},
	"vlc":      {8080, 4212},
	"ffmpeg":   {8090},

	// Game Servers
	"minecraft":     {25565, 25575},
	"steam":         {27015, 27016, 27017, 27018, 27019, 27020},
	"quake":         {26000, 27960},
	"counterstrike": {27015, 27016, 27017, 27018, 27019, 27020},
	"teamspeak":     {9987, 10011, 30033},
	"mumble":        {64738},
	"ventrilo":      {3784, 3785},
	"discord":       {64738},

	// IoT & Smart Devices
	"homeassistant": {8123},
	"openhab":       {8080, 8443},
	"domoticz":      {8080, 6144},
	"hue":           {80},
	"zigbee":        {8080},
	"mqtt":          {1883, 8883},
	"coap":          {5683, 5684},

	// Developer Tools
	"vscode":       {8000, 8080},
	"codeserver":   {8080},
	"jupyter":      {8888},
	"rstudio":      {8787},
	"phpmyadmin":   {80, 443},
	"adminer":      {8080},
	"pgadmin":      {5050},
	"dbeaver":      {8080},
	"redisinsight": {8001},

	// Cloud Services
	"aws":          {80, 443, 3389, 5432, 3306, 27017, 6379},
	"azure":        {80, 443, 3389, 1433, 3306},
	"gcp":          {80, 443, 3389, 5432, 3306},
	"digitalocean": {80, 443, 22, 3306, 5432},
	"heroku":       {80, 443, 5000, 5432},
	"cloudflare":   {80, 443, 2053, 2083, 2087, 2096, 8443},

	// Special Protocols
	"websocket": {80, 443, 8080, 8443},
	"grpc":      {50051, 50052},
	"graphql":   {4000, 4001, 8080, 9000},
	"rest":      {80, 443, 8080, 8443},
	"soap":      {80, 443, 8080, 8443},
	"xmlrpc":    {80, 443, 8080, 8443},
	"jsonrpc":   {80, 443, 8080, 8443},

	// Other
	"phpfpm":            {9000},
	"fluentd":           {24224},
	"telegraf":          {8125, 8186},
	"chronograf":        {8888},
	"kapacitor":         {9092},
	"harbor":            {80, 443, 4443},
	"portainer":         {9000, 8000},
	"traefik":           {80, 443, 8080},
	"envoy":             {9901, 10000},
	"istio":             {15000, 15001, 15004, 15006, 15010, 15012, 15014, 15021},
	"linkerd":           {4140, 4143, 4191, 9990, 9998},
	"kong":              {8000, 8001, 8443, 8444},
	"nginxproxymanager": {80, 81, 443},
	"minio":             {9000, 9001},
	"ceph":              {6789, 6800, 7300},
	"glusterfs":         {24007, 24008, 49152},
	"rethinkdb":         {8080, 28015, 29015},
	"couchbase":         {8091, 8092, 8093, 8094, 8095, 8096, 11210},
	"orientdb":          {2424, 2480},
	"tarantool":         {3301},
	"timescaledb":       {5432},
	"questdb":           {9000, 9009, 8812},
	"ydb":               {2135, 2136, 8765},
}

// PortStatus represents the open/closed state of a single port for a service.
type PortStatus struct {
	Service string
	Port    int
	Open    bool
}

// CheckPorts probes each port for the given services on localhost.
// It returns a slice of PortStatus entries for every configured port.
// The timeout controls the DialTimeout duration per port.
func CheckPorts(svcs []string, timeout time.Duration) []PortStatus {
	var (
		out []PortStatus
		mu  sync.Mutex
		wg  sync.WaitGroup
	)
	for _, svc := range svcs {
		ports, ok := DefaultPorts[svc]
		if !ok {
			continue
		}
		for _, p := range ports {
			wg.Add(1)
			go func(s string, port int) {
				defer wg.Done()
				addr := fmt.Sprintf("127.0.0.1:%d", port)
				conn, err := net.DialTimeout("tcp", addr, timeout)
				open := err == nil
				if conn != nil {
					conn.Close()
				}
				mu.Lock()
				out = append(out, PortStatus{s, port, open})
				mu.Unlock()
			}(svc, p)
		}
	}
	wg.Wait()
	return out
}
