package langfetch

func init() { register("node", fetchNode) }

func fetchNode() (LangInfo, error) {
	ver, err := runFirst("node", "-v")
	if err != nil {
		return LangInfo{}, err
	}
	npm := runSilently("npm", "-v")
	pkgs := runSilently("npm", "ls", "-g", "--depth=0")
	return LangInfo{"Node.js", ver, []string{"npm: " + npm, "Global packages:\n" + pkgs}}, nil
}
