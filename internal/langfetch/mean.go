package langfetch

func init() { register("mean", fetchMEAN) }

func fetchMEAN() (LangInfo, error) {
    return stackFetcher("MEAN", []cmdSpec{
        {"MongoDB", "mongo", []string{"--version"}},
        {"Node", "node", []string{"-v"}},
        {"npm", "npm", []string{"-v"}},
        {"Angular CLI", "ng", []string{"version"}},
    })()
}
