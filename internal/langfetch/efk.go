package langfetch

func init() { register("efk", fetchEFK) }

func fetchEFK() (LangInfo, error) {
	return stackFetcher("EFK", []cmdSpec{{"Elasticsearch", "elasticsearch", []string{"--version"}}, {"Fluentd", "fluentd", []string{"--version"}}, {"Kibana", "kibana", []string{"--version"}}})()
}
