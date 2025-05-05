package langfetch

func init() { register("dcp", fetchDCP) }

func fetchDCP() (LangInfo, error) {
    return stackFetcher("DCP", []cmdSpec{{"Docker", "docker", []string{"--version"}}, {"Docker Compose", "docker", []string{"compose", "version"}}, {"Podman", "podman", []string{"--version"}},})()
}
