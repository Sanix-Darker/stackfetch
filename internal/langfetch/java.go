package langfetch

func init() { register("java", fetchJava) }

func fetchJava() (LangInfo, error) {
	ver, err := runFirst("java", "-version")
	if err != nil {
		return LangInfo{}, err
	}
	mvn := runSilently("mvn", "-v")
	gradle := runSilently("gradle", "-v")
	return LangInfo{"Java", ver, []string{mvn, gradle}}, nil
}
