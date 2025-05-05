package langfetch

func init() { register("pern", fetchPERN) }

func fetchPERN() (LangInfo, error) {
    return stackFetcher("PERN", []cmdSpec{
        {"PostgreSQL", "psql", []string{"--version"}},
        {"Node", "node", []string{"-v"}},
        {"npm", "npm", []string{"-v"}},
    })()
}
