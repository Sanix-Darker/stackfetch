package langfetch

func init() { register("lamp", fetchLAMP) }

func fetchLAMP() (LangInfo, error) {
    return stackFetcher("LAMP", []cmdSpec{
        {"Apache", "apache2", []string{"-v"}},
        {"MySQL", "mysql", []string{"--version"}},
        {"PHP", "php", []string{"-v"}},
    })()
}
