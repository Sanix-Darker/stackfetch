package langfetch

func init() { register("wamp", fetchWAMP) }

func fetchWAMP() (LangInfo, error) {
    return stackFetcher("WAMP", []cmdSpec{
        {"Apache", "httpd", []string{"-v"}},
        {"MySQL", "mysql", []string{"--version"}},
        {"PHP", "php", []string{"-v"}},
    })()
}
