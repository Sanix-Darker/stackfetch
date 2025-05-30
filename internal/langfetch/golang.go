package langfetch

func init() { register("go", fetchGo) }

func fetchGo() (LangInfo, error) {
	ver, err := runFirst("go", "version")
	if err != nil {
		return LangInfo{}, err
	}
	mod := runSilently("go", "env", "GOMOD")
	return LangInfo{"Go", ver, []string{"Module file: " + mod}}, nil
}
