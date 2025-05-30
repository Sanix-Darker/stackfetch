package langfetch

func init() { register("tall", fetchTALL) }
func fetchTALL() (LangInfo, error) {
	return stackFetcher("TALL", []cmdSpec{{"Laravel", "php", []string{"artisan", "--version"}}, {"Tailwind CLI", "tailwindcss", []string{"--version"}}, {"Alpine.js", "npm", []string{"list", "alpinejs", "--depth=0"}}})()
}
