package langfetch

func init() { register("perl", fetchPerl) }

func fetchPerl() (LangInfo, error) {
    ver, err := runFirst("perl", "-v")
    if err != nil { return LangInfo{}, err }
    cpan := runSilently("cpan", "-v")
    return LangInfo{"Perl", ver, []string{cpan}}, nil
}
