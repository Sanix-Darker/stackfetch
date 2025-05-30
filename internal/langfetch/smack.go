package langfetch

func init() { register("smack", fetchSMACK) }

func fetchSMACK() (LangInfo, error) {
	return stackFetcher("SMACK", []cmdSpec{{"Spark", "spark‑submit", []string{"--version"}}, {"Mesos", "mesos‑master", []string{"--version"}}, {"Akka (sbt)", "sbt", []string{"akkaVersion"}}, {"Cassandra", "cassandra", []string{"-v"}}, {"Kafka", "kafka‑topology‑builder", []string{"--version"}}})()
}
