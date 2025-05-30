package langfetch

func init() { register("mevn", fetchMEVN) }

func fetchMEVN() (LangInfo, error) {
	return stackFetcher("MEVN", []cmdSpec{
		{"MongoDB", "mongo", []string{"--version"}},
		{"Node", "node", []string{"-v"}},
		{"npm", "npm", []string{"-v"}},
		{"Vue CLI", "vue", []string{"--version"}},
	})()
}
