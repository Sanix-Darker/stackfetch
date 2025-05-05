package langfetch

func init() { register("swift", fetchSwift) }

func fetchSwift() (LangInfo, error) {
    ver, err := runFirst("swift", "--version")
    if err != nil { return LangInfo{}, err }
    return LangInfo{"Swift", ver, nil}, nil
}
