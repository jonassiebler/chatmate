class Chatmate < Formula
  desc "AI-powered CLI for managing VS Code Copilot Chat agents (chatmates)"
  homepage "https://github.com/jonassiebler/chatmate"
  url "https://github.com/jonassiebler/chatmate.git",
      using:    :git,
      revision: "cd8718b"
  version "20250817113552"
  license "MIT"
  head "https://github.com/jonassiebler/chatmate.git", branch: "dev"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(output: bin/"chatmate"), "."
    man1.install "docs/man/chatmate.1"
  end

  test do
    assert_match "ChatMate", shell_output("#{bin}/chatmate --help")
  end
end
