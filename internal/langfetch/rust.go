package langfetch

func init() { register("rust", fetchRust) }

func fetchRust() (LangInfo, error) {
	ver, err := runFirst("rustc", "--version")
	if err != nil {
		return LangInfo{}, err
	}
	cargo := runSilently("cargo", "--version")
	return LangInfo{"Rust", ver, []string{cargo}}, nil
}
