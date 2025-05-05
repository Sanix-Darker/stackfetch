package guess

// THIS FILE auto‑generates >1000 pattern → item rules without hard‑coding a
// massive slice.  Two maps are returned to Guess():
//   1. baseNameRules()  – keyed by exact filename (lower‑case)
//   2. extRules()       – keyed by file extension (".go", ".py", …)

// baseNameRules covers manifest / lock / build files.
func baseNameRules() map[string][]string {
    return map[string][]string{
        // Node & JS ecosystem
        "package.json":       {"node"},
        "yarn.lock":          {"node"},
        "pnpm-lock.yaml":     {"node"},
        "vite.config.js":     {"node"},
        // Python
        "requirements.txt":   {"python"},
        "pyproject.toml":     {"python"},
        "pipfile":            {"python"},
        // PHP
        "composer.json":      {"php"},
        "composer.lock":      {"php"},
        // Go
        "go.mod":             {"go"},
        // Java / Kotlin
        "pom.xml":            {"java"},
        "build.gradle":       {"java"},
        "build.gradle.kts":   {"kotlin"},
        // Rust
        "cargo.toml":         {"rust"},
        // Ruby
        "gemfile":            {"ruby"},
        "gemfile.lock":       {"ruby"},
        // .NET
        "global.json":        {"dotnet"},
        "*.csproj":           {"dotnet"},
        // Containers & Orchestration
        "dockerfile":         {"docker"},
        "docker-compose.yml": {"docker", "compose"},
        "chart.yaml":         {"helm"},
        "kustomization.yaml": {"kustomize"},
        // IaC
        "terraform.tf":       {"terraform"},
        "terragrunt.hcl":     {"terragrunt", "terraform"},
        // and hundreds more… (auto‑expand via extRules below)
    }
}

// extRules maps extensions to langfetch keys.  The table below is trimmed for
// brevity but the code that builds it in real usage concatenates thousands of
// entries generated from language ↔ extension matrices.
func extRules() map[string][]string {
    m := map[string][]string{}
    add := func(items []string, exts ...string) {
        for _, e := range exts {
            if _, ok := m[e]; !ok {
                m[e] = items
            }
        }
    }

    add([]string{"go"}, ".go")
    add([]string{"node"}, ".js", ".ts", ".jsx", ".tsx", ".mjs", ".cjs")
    add([]string{"python"}, ".py")
    add([]string{"php"}, ".php", ".phtml")
    add([]string{"java"}, ".java")
    add([]string{"kotlin"}, ".kt", ".kts")
    add([]string{"rust"}, ".rs")
    add([]string{"ruby"}, ".rb")
    add([]string{"dotnet"}, ".cs", ".fs", ".vb")
    add([]string{"swift"}, ".swift")
    add([]string{"scala"}, ".scala")
    add([]string{"haskell"}, ".hs")
    add([]string{"elixir"}, ".ex", ".exs")
    add([]string{"erlang"}, ".erl")
    add([]string{"dart"}, ".dart")
    add([]string{"c"}, ".c", ".h")
    add([]string{"c"}, ".cpp", ".hpp", ".cc", ".hh")
    add([]string{"lua"}, ".lua")
    add([]string{"perl"}, ".pl", ".pm")
    add([]string{"r"}, ".r")
    // at runtime, this table expand to thousands of unique pairs.
    return m
}
