package langfetch

func init() { register("lemp", fetchLEMP) }

func fetchLEMP() (LangInfo, error) {
    return stackFetcher("LEMP", []cmdSpec{
        {"Nginx", "nginx", []string{"-v"}},
        {"MySQL", "mysql", []string{"--version"}},
        {"PHP", "php", []string{"-v"}},
    })()
}
