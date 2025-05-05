package langfetch

func init() { register("lapp", fetchLAPP) }

func fetchLAPP() (LangInfo, error) {
    return stackFetcher("LAPP", []cmdSpec{{"Apache", "apache2", []string{"-v"}}, {"PostgreSQL", "psql", []string{"--version"}}, {"PHP", "php", []string{"-v"}},})()
}
