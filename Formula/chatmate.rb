class Chatmate < Formula
  desc "AI-powered CLI for managing VS Code Copilot Chat agents (chatmates)"
  homepage "https://github.com/jonassiebler/chatmate"
  url "https://github.com/jonassiebler/chatmate/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "6c625da12eabf2762ae8d6059c5c5a2448d21090039f66cc344252c77961abd5"
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
