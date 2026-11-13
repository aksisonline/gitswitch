# Template only — a build-from-source formula in the shape Homebrew core
# expects, kept here as the starting point for a homebrew-core submission.
#
# The formulae people actually install live in aksisonline/homebrew-tap and
# are generated on release:
#   gitswitch       stable — written by GoReleaser (.goreleaser.yaml)
#   gitswitch-beta  beta   — written by .github/scripts/beta-formula.sh
#
# This file is not updated on release; do not install from it directly.
#
# What keeps gitswitch acceptable as a Homebrew formula: it never updates
# itself when Homebrew installed it. `gitswitch upgrade`, the in-app update
# prompt and `gitswitch beta` / `stable` all run `brew` instead
# (internal/version/install.go). Keep it that way.
class Gitswitch < Formula
  desc "Switch git identities and GitHub accounts on one machine"
  homepage "https://gitswitch.dev"
  url "https://github.com/aksisonline/gitswitch/archive/refs/tags/v0.5.0.tar.gz"
  sha256 "TODO_SHA256"
  license "Apache-2.0"
  head "https://github.com/aksisonline/gitswitch.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=v#{version}"), "./cmd/gitswitch"
    generate_completions_from_executable(bin/"gitswitch", "completion")
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/gitswitch --version")
  end
end
