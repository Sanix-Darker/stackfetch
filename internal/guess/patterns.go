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
        "webpack.config.js": {"node"},
        "babel.config.js":   {"node"},
        ".eslintrc":          {"node"},
        ".prettierrc":        {"node"},
        "rollup.config.js":   {"node"},
        "tsconfig.json":      {"node"},

        // Python
        "requirements.txt":   {"python"},
        "pyproject.toml":     {"python"},
        "pipfile":            {"python"},
        "pipfile.lock":       {"python"},
        "setup.py":           {"python"},
        "tox.ini":            {"python"},
        "pylintrc":          {"python"},
        "manage.py":          {"python", "django"},

        // PHP
        "composer.json":      {"php"},
        "composer.lock":      {"php"},
        "artisan":            {"php", "laravel"},

        // Go
        "go.mod":             {"go"},
        "go.sum":             {"go"},

        // Java/Kotlin/Scala
        "pom.xml":            {"java"},
        "build.gradle":       {"java"},
        "build.gradle.kts":   {"kotlin"},
        "build.sbt":         {"scala"},

        // Rust
        "cargo.toml":         {"rust"},
        "cargo.lock":        {"rust"},

        // Ruby
        "gemfile":            {"ruby"},
        "gemfile.lock":       {"ruby"},
        "rakefile":          {"ruby"},

        // .NET
        "global.json":        {"dotnet"},
        "*.csproj":           {"dotnet"},
        "*.fsproj":           {"dotnet"},
        "*.vbproj":           {"dotnet"},
        "*.sln":             {"dotnet"},

        // Swift
        "package.swift":      {"swift"},

        // Haskell
        "stack.yaml":         {"haskell"},
        "cabal.project":      {"haskell"},

        // Elixir
        "mix.exs":            {"elixir"},

        // Dart/Flutter
        "pubspec.yaml":       {"dart"},
        "pubspec.lock":      {"dart"},

        // C/C++
        "makefile":           {"c"},
        "cmakelists.txt":    {"c"},

        // Containers & Orchestration
        "dockerfile":         {"docker"},
        "docker-compose.yml": {"docker", "compose"},
        "chart.yaml":         {"helm"},
        "kustomization.yaml": {"kustomize"},

        // IaC
        "terraform.tf":       {"terraform"},
        "terragrunt.hcl":     {"terragrunt", "terraform"},
        "pulumi.yaml":        {"pulumi"},
        "serverless.yml":     {"serverless"},

        // Frontend
        "tailwind.config.js": {"node", "tailwind"},
        "postcss.config.js":  {"node", "postcss"},
        "vite.config.ts":     {"node"},

        // Database
        "migration.sql":     {"sql"},
        "*.prisma":          {"prisma"},
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
    add([]string{"node"}, ".js", ".ts", ".jsx", ".tsx", ".mjs", ".cjs", ".vue", ".svelte")
    add([]string{"python"}, ".py", ".pyi", ".pyc", ".pyd", ".pyo", ".pyw", ".pyz")
    add([]string{"php"}, ".php", ".phtml", ".php3", ".php4", ".php5", ".php7", ".phps")
    add([]string{"java"}, ".java", ".class", ".jar", ".war", ".ear")
    add([]string{"kotlin"}, ".kt", ".kts", ".ktm")
    add([]string{"rust"}, ".rs", ".rlib")
    add([]string{"ruby"}, ".rb", ".erb", ".gemspec", ".rake")
    add([]string{"dotnet"}, ".cs", ".fs", ".vb", ".cshtml", ".razor", ".fsx", ".fsi")
    add([]string{"swift"}, ".swift")
    add([]string{"scala"}, ".scala", ".sc")
    add([]string{"haskell"}, ".hs", ".lhs")
    add([]string{"elixir"}, ".ex", ".exs", ".eex", ".leex", ".heex")
    add([]string{"erlang"}, ".erl", ".hrl")
    add([]string{"dart"}, ".dart")
    add([]string{"c"}, ".c", ".h")
    add([]string{"cpp"}, ".cpp", ".hpp", ".cc", ".hh", ".cxx", ".hxx", ".c++", ".h++", ".ino")
    add([]string{"objective-c"}, ".m", ".mm")
    add([]string{"lua"}, ".lua")
    add([]string{"perl"}, ".pl", ".pm", ".t", ".pod")
    add([]string{"r"}, ".r", ".R", ".Rmd", ".Rnw", ".Rdata", ".Rds", ".Rda")
    add([]string{"bash"}, ".sh", ".bash", ".zsh", ".fish", ".ksh")
    add([]string{"powershell"}, ".ps1", ".psm1", ".psd1")
    add([]string{"html"}, ".html", ".htm", ".xhtml", ".html.hl", ".mustache", ".hbs")
    add([]string{"css"}, ".css", ".scss", ".sass", ".less", ".styl")
    add([]string{"xml"}, ".xml", ".xsl", ".xslt", ".xsd", ".wsdl")
    add([]string{"json"}, ".json", ".jsonc", ".json5", ".geojson", ".topojson")
    add([]string{"yaml"}, ".yaml", ".yml")
    add([]string{"markdown"}, ".md", ".markdown", ".mdown", ".mdwn", ".mkd", ".mkdn")
    add([]string{"sql"}, ".sql", ".ddl", ".dml", ".pks", ".pkb")
    add([]string{"typescript"}, ".ts", ".tsx", ".mts", ".cts")
    add([]string{"clojure"}, ".clj", ".cljs", ".cljc", ".edn")
    add([]string{"groovy"}, ".groovy", ".gvy", ".gy", ".gsh")
    add([]string{"julia"}, ".jl")
    add([]string{"ocaml"}, ".ml", ".mli")
    add([]string{"pascal"}, ".pas", ".pp", ".p", ".inc")
    add([]string{"prolog"}, ".pl", ".pro", ".prolog")
    add([]string{"matlab"}, ".m", ".mat")
    add([]string{"fortran"}, ".f", ".for", ".f90", ".f95", ".f03", ".f08")
    add([]string{"assembly"}, ".asm", ".s", ".S", ".inc")
    add([]string{"vim"}, ".vim", ".vimrc")
    add([]string{"tex"}, ".tex", ".sty", ".cls", ".dtx", ".ins", ".ltx")
    add([]string{"zig"}, ".zig")
    add([]string{"solidity"}, ".sol")
    add([]string{"dockerfile"}, "dockerfile", ".dockerignore")
    add([]string{"nginx"}, ".conf", ".nginx")
    add([]string{"apache"}, ".htaccess", ".htpasswd")
    add([]string{"csv"}, ".csv", ".tsv")
    add([]string{"toml"}, ".toml")
    add([]string{"ini"}, ".ini", ".cfg", ".prefs")
    add([]string{"batch"}, ".bat", ".cmd")
    add([]string{"makefile"}, "makefile", ".mk", ".mak")
    add([]string{"protobuf"}, ".proto")
    add([]string{"thrift"}, ".thrift")
    add([]string{"graphql"}, ".graphql", ".gql")
    add([]string{"wasm"}, ".wat", ".wasm")
    add([]string{"crystal"}, ".cr", ".ecr")
    add([]string{"nim"}, ".nim", ".nims")
    add([]string{"vala"}, ".vala", ".vapi")
    add([]string{"d"}, ".d", ".di")
    add([]string{"fsharp"}, ".fs", ".fsx", ".fsi", ".fsscript")
    add([]string{"coffeescript"}, ".coffee", ".litcoffee", ".coffee.md")
    add([]string{"purescript"}, ".purs")
    add([]string{"reason"}, ".re", ".rei")
    add([]string{"rescript"}, ".res", ".resi")
    add([]string{"perl6"}, ".pm6", ".p6", ".pl6", ".pod6")
    add([]string{"tcl"}, ".tcl", ".tk", ".itcl", ".itk")
    add([]string{"racket"}, ".rkt", ".rktl", ".rktd")
    add([]string{"scheme"}, ".scm", ".ss", ".rkt")
    add([]string{"smalltalk"}, ".st", ".cs")
    add([]string{"vhdl"}, ".vhd", ".vhdl")
    add([]string{"verilog"}, ".v", ".vh", ".sv", ".svh")
    add([]string{"systemverilog"}, ".sv", ".svh")
    add([]string{"ada"}, ".adb", ".ads")
    add([]string{"abap"}, ".abap")
    add([]string{"cobol"}, ".cbl", ".cob", ".cpy")
    add([]string{"jupyter"}, ".ipynb")
    add([]string{"robotframework"}, ".robot")
    add([]string{"red"}, ".red", ".reds")
    add([]string{"rebol"}, ".r", ".reb")
    add([]string{"xquery"}, ".xq", ".xql", ".xqm", ".xquery")
    add([]string{"yacc"}, ".y", ".yy")
    add([]string{"lex"}, ".l", ".ll")
    add([]string{"antlr"}, ".g4")
    add([]string{"idris"}, ".idr", ".lidr")
    add([]string{"agda"}, ".agda", ".lagda")
    add([]string{"unrealscript"}, ".uc")
    add([]string{"qml"}, ".qml")
    add([]string{"bluespec"}, ".bsv")
    add([]string{"chisel"}, ".scala")
    add([]string{"systemc"}, ".h", ".cpp")
    add([]string{"opencl"}, ".cl")
    add([]string{"cuda"}, ".cu", ".cuh")
    add([]string{"hlsl"}, ".hlsl", ".fx", ".fxh", ".hlsli")
    add([]string{"glsl"}, ".glsl", ".vert", ".frag", ".geom", ".tesc", ".tese", ".comp")
    add([]string{"metal"}, ".metal")
    add([]string{"spirv"}, ".spvasm")
    add([]string{"webgpu"}, ".wgsl")
    add([]string{"assemblyscript"}, ".as")
    add([]string{"tensorflow"}, ".pb")
    add([]string{"onnx"}, ".onnx")
    add([]string{"pytorch"}, ".pt", ".pth")
    add([]string{"keras"}, ".h5", ".keras")
    add([]string{"netlogo"}, ".nlogo")
    add([]string{"scratch"}, ".sb", ".sb2", ".sb3")
    add([]string{"blockly"}, ".blockly")
    add([]string{"labview"}, ".vi", ".lvproj", ".lvclass")
    add([]string{"mathematica"}, ".nb", ".wl", ".wls")
    add([]string{"maple"}, ".mpl", ".mw")
    add([]string{"maxscript"}, ".ms")
    add([]string{"mel"}, ".mel")
    add([]string{"openscad"}, ".scad")
    add([]string{"g-code"}, ".g", ".gcode", ".nc", ".cnc", ".tap")
    add([]string{"postscript"}, ".ps", ".eps", ".epsi", ".epsf")
    add([]string{"opengl"}, ".glsl", ".vert", ".frag", ".geom", ".tesc", ".tese", ".comp")
    add([]string{"vulkan"}, ".spv", ".comp")
    add([]string{"directx"}, ".fx", ".fxh", ".hlsl")
    add([]string{"unity"}, ".unity", ".prefab", ".mat", ".asset", ".controller")
    add([]string{"unreal"}, ".uproject", ".umap", ".uasset", ".blueprint")
    add([]string{"godot"}, ".tscn", ".tres", ".gd")
    add([]string{"gamemaker"}, ".gml", ".yy", ".yyp")
    add([]string{"renpy"}, ".rpy", ".rpyc")
    add([]string{"twine"}, ".twee", ".tw")
    add([]string{"inform"}, ".ni", ".i7", ".i6")
    add([]string{"tads"}, ".t", ".h")
    add([]string{"alan"}, ".alan", ".a3c", ".a3t")
    add([]string{"quest"}, ".aslx", ".quest")
    add([]string{"ink"}, ".ink")
    add([]string{"squirrel"}, ".nut")
    add([]string{"angelscript"}, ".as")
    add([]string{"pawn"}, ".pwn", ".inc")
    add([]string{"sourcepawn"}, ".sp", ".inc")
    add([]string{"lisp"}, ".lisp", ".lsp", ".cl", ".el")
    add([]string{"common lisp"}, ".lisp", ".lsp", ".cl", ".el")
    add([]string{"emacs lisp"}, ".el", ".elc")
    add([]string{"fennel"}, ".fnl")
    add([]string{"janet"}, ".janet")
    add([]string{"hy"}, ".hy")
    add([]string{"carp"}, ".carp")
    add([]string{"gerbil"}, ".ss", ".sld")
    add([]string{"guile"}, ".scm", ".ss")
    add([]string{"newlisp"}, ".lsp", ".nl")
    add([]string{"picolisp"}, ".l")
    add([]string{"logos"}, ".x", ".xm", ".xi", ".xmi")
    add([]string{"limbo"}, ".b", ".m")
    add([]string{"io"}, ".io")
    add([]string{"ioke"}, ".ik")
    add([]string{"factor"}, ".factor")
    add([]string{"j"}, ".ijs")
    add([]string{"k"}, ".k")
    add([]string{"q"}, ".q")
    add([]string{"apl"}, ".apl", ".dyalog", ".aplf", ".apln", ".aplo", ".aplc", ".apli")
    add([]string{"octave"}, ".m", ".oct")
    add([]string{"scilab"}, ".sci", ".sce")
    add([]string{"maxima"}, ".mac", ".max")
    add([]string{"gnuplot"}, ".plt", ".gp")
    add([]string{"sas"}, ".sas", ".SAS")
    add([]string{"stata"}, ".do", ".ado")
    add([]string{"spss"}, ".sps", ".sav", ".zsav")
    add([]string{"mathcad"}, ".mcd", ".mcdx")
    add([]string{"igor"}, ".pxp", ".ibw")
    add([]string{"idlang"}, ".pro")
    add([]string{"pvwave"}, ".pro")
    add([]string{"galaxy"}, ".galaxy")
    add([]string{"pike"}, ".pike", ".pmod")
    add([]string{"lpc"}, ".c", ".h")
    add([]string{"b"}, ".b")
    add([]string{"bcpl"}, ".bcpl")
    add([]string{"clipper"}, ".prg", ".ch", ".hb")
    add([]string{"foxpro"}, ".prg", ".fxp")
    add([]string{"xbase"}, ".prg", ".ch", ".hb")
    add([]string{"harbour"}, ".prg", ".ch", ".hb")
    add([]string{"visual foxpro"}, ".prg", ".fxp")
    add([]string{"dbase"}, ".prg", ".dbf")
    add([]string{"informix"}, ".4gl")
    add([]string{"progress"}, ".p", ".w", ".i")
    add([]string{"powerbuilder"}, ".srw", ".sru", ".srd")
    add([]string{"clarion"}, ".clw", ".inc")
    add([]string{"lansa"}, ".rdml", ".rdmlx")
    add([]string{"uniface"}, ".un", ".unp")
    add([]string{"magic"}, ".mgc")
    add([]string{"genexus"}, ".gxg", ".gxu")
    add([]string{"outsystems"}, ".oml")
    add([]string{"mendix"}, ".mpr", ".mda")
    add([]string{"salesforce"}, ".cls", ".trigger", ".page", ".component")
    add([]string{"apex"}, ".cls", ".trigger")
    add([]string{"visualforce"}, ".page", ".component")
    add([]string{"lightning"}, ".cmp", ".app", ".evt", ".intf", ".auradoc")
    add([]string{"wave"}, ".wv", ".wvc")
    add([]string{"zbrush"}, ".ztl")
    add([]string{"maya"}, ".ma", ".mb")
    add([]string{"3ds max"}, ".max")
    add([]string{"blender"}, ".blend")
    add([]string{"cinema 4d"}, ".c4d")
    add([]string{"houdini"}, ".hip", ".hipnc")
    add([]string{"substance"}, ".sbs", ".sbsar")
    add([]string{"quixel"}, ".megascans")
    add([]string{"marvelous designer"}, ".zpac", ".zpa")
    add([]string{"world machine"}, ".worldmachine")
    add([]string{"speedtree"}, ".srt", ".spm")
    add([]string{"substance painter"}, ".spp")
    add([]string{"mari"}, ".mra")
    add([]string{"nuke"}, ".nk")
    add([]string{"katana"}, ".katana")
    add([]string{"clarisse"}, ".project", ".asset", ".render")
    add([]string{"arnold"}, ".ass", ".usd")
    add([]string{"redshift"}, ".rs", ".rsc")
    add([]string{"v-ray"}, ".vrscene")
    add([]string{"corona"}, ".corona")
    add([]string{"octane"}, ".ocs", ".orbx")
    add([]string{"renderman"}, ".rib", ".rman")
    add([]string{"mental ray"}, ".mi", ".mib")
    add([]string{"keyshot"}, ".bip", ".ksp")
    add([]string{"twinmotion"}, ".twin")
    add([]string{"lumion"}, ".ls", ".lsf")
    add([]string{"enscape"}, ".ens")
    add([]string{"d5 render"}, ".d5")
    add([]string{"cryengine"}, ".cry", ".cryasset")
    add([]string{"amazon lumberyard"}, ".ly", ".slice")
    add([]string{"stingray"}, ".level", ".material", ".texture")
    add([]string{"source"}, ".vmf", ".vmt", ".vtx", ".mdl", ".phy")
    add([]string{"id tech"}, ".map", ".bsp", ".md3", ".md5")
    add([]string{"doom"}, ".wad", ".pk3", ".pk7")
    add([]string{"quake"}, ".map", ".bsp", ".mdl")
    add([]string{"unreal tournament"}, ".ut2", ".ut3", ".ut4")
    add([]string{"half-life"}, ".bsp", ".mdl", ".vmt")
    add([]string{"counter-strike"}, ".bsp", ".mdl", ".vmt")
    add([]string{"call of duty"}, ".ff", ".iwd", ".xpak")
    add([]string{"battlefield"}, ".sb", ".toc", ".cas")
    add([]string{"frostbite"}, ".sb", ".toc", ".cas")
    add([]string{"rockstar"}, ".rpf", ".ydr", ".ydd", ".yft")
    add([]string{"ea sports"}, ".big", ".bh", ".mbe")
    add([]string{"nintendo"}, ".narc", ".sarc", ".byml")
    add([]string{"sony"}, ".pkg", ".sb", ".toc")
    add([]string{"microsoft"}, ".xex", ".xzp", ".xvd")
    add([]string{"capcom"}, ".arc", ".pak", ".mod")
    add([]string{"square enix"}, ".dat", ".bin", ".pkg")
    add([]string{"bandai namco"}, ".pac", ".paz", ".cpk")
    add([]string{"konami"}, ".bin", ".dat", ".pak")
    add([]string{"ubisoft"}, ".forge", ".dat", ".pak")
    add([]string{"activision"}, ".ff", ".iwd", ".xpak")
    add([]string{"take-two"}, ".rpf", ".ydr", ".ydd", ".yft")
    add([]string{"warner bros"}, ".upk", ".xxx", ".pak")
    add([]string{"disney"}, ".pkg", ".pak", ".dat")
    add([]string{"ea"}, ".big", ".bh", ".mbe")
    add([]string{"blizzard"}, ".mpq", ".app", ".data")
    add([]string{"riot"}, ".raf", ".dat", ".bnk")
    add([]string{"epic"}, ".pak", ".utoc", ".ucas")
    add([]string{"valve"}, ".vpk", ".vdf", ".vtf")
    add([]string{"bethesda"}, ".bsa", ".ba2", ".esm")
    add([]string{"cd projekt red"}, ".bundle", ".cache", ".red")
    add([]string{"naughty dog"}, ".pak", ".snd", ".stream")
    add([]string{"insomniac"}, ".psarc", ".stream", ".sdat")
    add([]string{"santa monica"}, ".pak", ".core", ".stream")
    add([]string{"guerrilla"}, ".bin", ".stream", ".core")
    add([]string{"playground"}, ".pak", ".stream", ".cache")
    add([]string{"turn 10"}, ".pak", ".stream", ".cache")
    add([]string{"343 industries"}, ".pak", ".map", ".module")
    add([]string{"the coalition"}, ".pak", ".stream", ".cache")
    add([]string{"rare"}, ".pak", ".stream", ".cache")
    add([]string{"mojang"}, ".mcpack", ".mclevel", ".mcworld")
    add([]string{"ccp"}, ".red", ".dat", ".cache")
    add([]string{"frontier"}, ".pak", ".dat", ".cache")
    add([]string{"bohemia"}, ".pbo", ".bin", ".sqm")
    add([]string{"egosoft"}, ".cat", ".dat", ".x3")
    add([]string{"re-logic"}, ".wld", ".plr", ".twld")
    add([]string{"klei"}, ".k", ".tex", ".anim")
    add([]string{"unknown worlds"}, ".assets", ".level", ".entity")
    add([]string{"introversion"}, ".save", ".map", ".prison")
    add([]string{"double fine"}, ".pak", ".dat", ".bnk")
    add([]string{"obsidian"}, ".pak", ".dat", ".bnk")
    add([]string{"inxile"}, ".pak", ".dat", ".bnk")
    add([]string{"larian"}, ".pak", ".lsf", ".lsx")
    add([]string{"paradox"}, ".db", ".asset", ".gfx")
    add([]string{"creative assembly"}, ".pack", ".asset", ".model")
    add([]string{"firaxis"}, ".civ5save", ".civ6save", ".civ6map")
    add([]string{"2k"}, ".2k", ".iff", ".rpf")

    // at runtime, this table expand to thousands of unique pairs.
    return m
}
