package langfetch

func init() { register("python", fetchPython); register("py", fetchPython) }

func fetchPython() (LangInfo, error) {
	ver, err := runFirst("python3", "--version")
	if err != nil {
		return LangInfo{}, err
	}
	pip := runSilently("pip3", "--version")
	return LangInfo{"Python", ver, []string{pip}}, nil
}
