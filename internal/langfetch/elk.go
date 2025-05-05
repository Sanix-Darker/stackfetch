package langfetch

func init() { register("elk", fetchELK) }

func fetchELK() (LangInfo, error) {
    return stackFetcher("ELK", []cmdSpec{
        {"Elasticsearch", "elasticsearch", []string{"--version"}},
        {"Logstash", "logstash", []string{"--version"}},
        {"Kibana", "kibana", []string{"--version"}},
    })()
}
