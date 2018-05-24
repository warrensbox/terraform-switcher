class Tfswitch < Formula
  desc "The tfswitch command lets you switch between terraform versions."
  homepage "https://warren-veerasingam.github.io/terraform-switcher/"
  url "https://github.com/warren-veerasingam/terraform-switcher/archive/0.2.180.tar.gz"
  head "https://github.com/warren-veerasingam/terraform-switcher.git"
  version "0.2.180"
  sha256 "4e6403fc0046f25355650b06095b6b03ec11f7bb85e2319785f36926a9f89618"
  
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
