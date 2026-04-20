class Gitswitch < Formula
  desc "Manage multiple git profiles locally"
  homepage "https://github.com/aksisonline/gitswitch"
  url "https://github.com/aksisonline/gitswitch/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "TODO_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", "gitswitch", "./cmd/gitswitch"
    bin.install "gitswitch"
  end

  test do
    system "#{bin}/gitswitch", "current"
  end
end
