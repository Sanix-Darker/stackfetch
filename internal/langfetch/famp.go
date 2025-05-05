package langfetch

func init() { register("famp", fetchFAMP) }

func fetchFAMP() (LangInfo, error) {
    return stackFetcher("FAMP", []cmdSpec{{"Apache", "apachectl", []string{"-v"}}, {"MySQL", "mysql", []string{"--version"}}, {"PHP", "php", []string{"-v"}},})()
}
