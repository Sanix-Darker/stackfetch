package langfetch

func init() { register("dotnet", fetchDotnet) }

func fetchDotnet() (LangInfo, error) {
	ver, err := runFirst("dotnet", "--version")
	if err != nil {
		return LangInfo{}, err
	}
	info := runSilently("dotnet", "--info")
	return LangInfo{".NET", ver, []string{info}}, nil
}
