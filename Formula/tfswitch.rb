class Tfswitch < Formula
  desc "The tfswitch command lets you switch between terraform versions."
  homepage "https://warren-veerasingam.github.io/terraform-switcher/"
  url "https://github.com/warren-veerasingam/terraform-switcher/archive/0.2.173.tar.gz"
  head "https://github.com/warren-veerasingam/terraform-switcher.git"
  version "0.2.173"
  sha256 "56b1d320c97a639fd03701e31b8bfe706777885d868d95c7920e6edb3bc22d1e"
  
  depends_on "git"
  depends_on "make" => :build
  depends_on "gcc" => :build
  depends_on "go" => :build
  
  conflicts_with "terraform"

  def install
    bin.install "tfswitch"
  end

  def caveats; <<~EOS
    Type 'tfswitch' on your command line and choose the terraform version that you want from the dropdown. This command currently only works on MacOs and Linux
  EOS
  end

  test do
    system "#{bin}/tfswitch --version"
  end
end
