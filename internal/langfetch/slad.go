package langfetch

func init() { register("slad", fetchSLAD) }
func fetchSLAD() (LangInfo, error) {
	return stackFetcher("SLAD", []cmdSpec{{"AWS CLI", "aws", []string{"--version"}}, {"SAM CLI", "sam", []string{"--version"}}})()
}
