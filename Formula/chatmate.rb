class Chatmate < Formula
  desc "AI-powered CLI for managing VS Code Copilot Chat agents (chatmates)"
  homepage "https://github.com/jonassiebler/chatmate"
  url "https://github.com/jonassiebler/chatmate/archive/refs/tags/v1.0.1.tar.gz"
  sha256 "8777c8bdeac0f5fcd607d358b9f3a7a5b0bc0027652aa120b75042017b222f83"
  license "MIT"
  head "https://github.com/jonassiebler/chatmate.git", branch: "dev"

  depends_on "go" => :build

  def install
    # Set version information via ldflags during build
    ldflags = %W[
      -X github.com/jonassiebler/chatmate/cmd.version=#{version}
      -X github.com/jonassiebler/chatmate/cmd.commit=homebrew-#{version}
      -X github.com/jonassiebler/chatmate/cmd.date=#{Date.today.iso8601}
    ]
    
    system "go", "build", *std_go_args(output: bin/"chatmate", ldflags: ldflags), "."
    man1.install "docs/man/chatmate.1"
  end

  test do
    assert_match "ChatMate", shell_output("#{bin}/chatmate --help")
  end
end
