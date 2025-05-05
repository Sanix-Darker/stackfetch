package langfetch

func init() { register("dart", fetchDart) }

func fetchDart() (LangInfo, error) {
    ver, err := runFirst("dart", "--version")
    if err != nil { return LangInfo{}, err }
    return LangInfo{"Dart", ver, nil}, nil
}
