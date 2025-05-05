package langfetch

func init() { register("mern", fetchMERN) }

func fetchMERN() (LangInfo, error) {
    return stackFetcher("MERN", []cmdSpec{
        {"MongoDB", "mongo", []string{"--version"}},
        {"Node", "node", []string{"-v"}},
        {"npm", "npm", []string{"-v"}},
        {"yarn", "yarn", []string{"-v"}},
    })()
}
