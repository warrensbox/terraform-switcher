class Tfswitch < Formula
  desc "The tfswitch command lets you switch between terraform versions."
  homepage "https://warren-veerasingam.github.io/terraform-switcher/"
  url "https://github.com/warren-veerasingam/terraform-switcher/archive/0.2.168.tar.gz"
  head "https://github.com/warren-veerasingam/terraform-switcher.git"
  version "0.2.168"
  sha256 "5eb2de36846e031817f90bfdb89f2b522423636dd3e24dbd618e11761483ff4c"
  
  depends_on "git"
  depends_on "make" => :build
  depends_on "gcc" => :build
  
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
