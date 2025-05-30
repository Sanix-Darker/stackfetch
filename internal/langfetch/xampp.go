package langfetch

func init() { register("xampp", fetchXAMPP) }

func fetchXAMPP() (LangInfo, error) {
	return stackFetcher("XAMPP", []cmdSpec{
		{"Apache", "apache2", []string{"-v"}},
		{"MariaDB", "mysql", []string{"--version"}},
		{"PHP", "php", []string{"-v"}},
		{"Perl", "perl", []string{"-v"}},
	})()
}
