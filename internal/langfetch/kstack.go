package langfetch

func init() { register("kstack", fetchKStack) }

func fetchKStack() (LangInfo, error) {
	return stackFetcher("K‑Stack", []cmdSpec{{"kubectl", "kubectl", []string{"version", "--client", "--short"}}, {"Kustomize", "kustomize", []string{"version"}}, {"Helm", "helm", []string{"version", "--short"}}, {"k3d", "k3d", []string{"version"}}})()
}
