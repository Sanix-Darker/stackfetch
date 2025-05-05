package langfetch

// This file registers a bunch of **single‑binary** runtimes & DevOps tools in
// ~20 lines.
import "strings"

func init() {
    simple := []struct{
        Key, Name, Bin string
        Args           []string
    }{
        // ─── Languages ────────────────────────────────────────────────
        {"kotlin", "Kotlin", "kotlinc", []string{"-version"}},
        {"typescript", "TypeScript", "tsc", []string{"--version"}},
        {"scala", "Scala", "scala", []string{"-version"}},
        {"haskell", "Haskell", "ghc", []string{"--version"}},
        {"clojure", "Clojure", "clojure", []string{"-Sdescribe"}},
        {"erlang", "Erlang", "erl", []string{"-eval", "erlang:display(erlang:system_info(otp_release)), halt()."}},
        {"elixir", "Elixir", "elixir", []string{"--version"}},
        {"fsharp", "F#", "dotnet", []string{"fsi", "--version"}},
        {"ocaml", "OCaml", "ocaml", []string{"-version"}},
        {"julia", "Julia", "julia", []string{"--version"}},
        // ─── DevOps / Containers ──────────────────────────────────────
        {"docker", "Docker", "docker", []string{"--version"}},
        {"compose", "Docker‑Compose", "docker", []string{"compose", "version"}},
        {"podman", "Podman", "podman", []string{"--version"}},
        {"kubectl", "kubectl", "kubectl", []string{"version", "--client", "--short"}},
        {"helm", "Helm", "helm", []string{"version", "--short"}},
        {"kustomize", "Kustomize", "kustomize", []string{"version"}},
        {"terraform", "Terraform", "terraform", []string{"version"}},
        {"ansible", "Ansible", "ansible", []string{"--version"}},
        {"vagrant", "Vagrant", "vagrant", []string{"--version"}},
        {"packer", "Packer", "packer", []string{"--version"}},
        {"git", "Git", "git", []string{"--version"}},
        {"make", "GNU Make", "make", []string{"--version"}},
        {"redis", "Redis", "redis‑server", []string{"--version"}},
        {"rabbitmq", "RabbitMQ", "rabbitmqctl", []string{"status"}},
        {"cassandra", "Cassandra", "cassandra", []string{"-v"}},
        {"kafka", "Apache Kafka", "kafka‑topology‑builder", []string{"--version"}},
    }
    for _, s := range simple {
        // just for loop vars
        key, name, bin, args := s.Key, s.Name, s.Bin, s.Args
        register(key, func() (LangInfo, error) {
            ver, err := runFirst(bin, args...)
            if err != nil {
                // degrade gracefully – mark missing but don’t error; allows multi‑fetch runs.
                return LangInfo{}, err
            }
            // Some commands output like "Version X", trim common prefixes.
            ver = strings.TrimSpace(strings.TrimPrefix(ver, name+" version"))
            return LangInfo{Name: name, Version: ver}, nil
        })
    }
}
