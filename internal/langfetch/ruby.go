package langfetch

func init() { register("ruby", fetchRuby) }

func fetchRuby() (LangInfo, error) {
    ver, err := runFirst("ruby", "--version")
    if err != nil { return LangInfo{}, err }
    gem := runSilently("gem", "--version")
    return LangInfo{"Ruby", ver, []string{"gem: " + gem}}, nil
}
