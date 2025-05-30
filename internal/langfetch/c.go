package langfetch

import (
	"errors"
	"strings"
)

func init() { register("c", fetchC); register("cpp", fetchC) }

func fetchC() (LangInfo, error) {
	gcc := runSilently("gcc", "--version")
	gpp := runSilently("g++", "--version")
	if gcc == "" && gpp == "" {
		return LangInfo{}, errors.New("gcc/g++ not found")
	}
	ver := strings.SplitN(gcc, "\n", 2)[0]
	return LangInfo{"C/C++", ver, []string{gpp}}, nil
}
