package langfetch
func init() { register("lnmp", fetchLNMP) }
func fetchLNMP() (LangInfo, error) {
    return stackFetcher("LNMP", []cmdSpec{{"Nginx", "nginx", []string{"-v"}}, {"MySQL", "mysql", []string{"--version"}}, {"PHP", "php", []string{"-v"}},})()
}
