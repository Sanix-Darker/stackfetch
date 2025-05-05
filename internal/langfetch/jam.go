package langfetch

func init() { register("jam", fetchJAM) }

func fetchJAM() (LangInfo, error) {
    return stackFetcher("JAMstack", []cmdSpec{
        {"Node", "node", []string{"-v"}},
        {"npm", "npm", []string{"-v"}},
        {"Static Site Gen (Gatsby)", "gatsby", []string{"--version"}},
    })()
}
