// Package skill embeds gitswitch's Agent Skills package (SKILL.md), living
// at this repo-root conventional path so it's discoverable both by Go's
// go:embed (which can't reach outside this directory) and by skills.sh's
// repo crawl (npx skills add aksisonline/gitswitch).
package skill

import _ "embed"

//go:embed SKILL.md
var Content []byte
