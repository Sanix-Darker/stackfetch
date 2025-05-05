package langfetch

func init() { register("r", fetchR) }

func fetchR() (LangInfo, error) {
    ver, err := runFirst("R", "--version")
    if err != nil { return LangInfo{}, err }
    return LangInfo{"R", ver, nil}, nil
}
