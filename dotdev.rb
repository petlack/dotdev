class Dotdev < Formula
  desc "Lightweight Web server for static HTML with built-in live reload written in Go."
  homepage "https://github.com/petlack/dotdev"
  url "https://github.com/petlack/dotdev/archive/refs/tags/v0.0.1.20250224.01.tar.gz"
  sha256 "cb74122d8a8c4f60c30ed194026afb64f8ffcb6d5b179b05d302601f6107ea40"
  license "GPL-3.0-or-later"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args
  end

  test do
    assert_match "dotdev", shell_output("#{bin}/dotdev --version")
  end
end
