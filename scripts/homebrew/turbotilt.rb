class Turbotilt < Formula
  desc "CLI pour environnements dev cloud-native"
  homepage "https://github.com/Fascinax/turbotilt"
  version "0.1.0"  # Change this to your current version
  license "MIT"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/Fascinax/turbotilt/releases/download/v#{version}/turbotilt-#{version}-darwin-arm64.zip"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"  # To be updated on release
    else
      url "https://github.com/Fascinax/turbotilt/releases/download/v#{version}/turbotilt-#{version}-darwin-amd64.zip"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"  # To be updated on release
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/Fascinax/turbotilt/releases/download/v#{version}/turbotilt-#{version}-linux-arm64.zip"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"  # To be updated on release
    else
      url "https://github.com/Fascinax/turbotilt/releases/download/v#{version}/turbotilt-#{version}-linux-amd64.zip"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"  # To be updated on release
    end
  end

  depends_on "docker" => :recommended
  depends_on "docker-compose" => :recommended
  depends_on "openjdk" => :recommended

  def install
    bin.install "turbotilt"
  end

  test do
    assert_match "Turbotilt #{version}", shell_output("#{bin}/turbotilt --version")
  end
  
  def caveats
    <<~EOS
      Turbotilt est installé!
      
      Pour commencer un nouveau projet:
        turbotilt init
        
      Pour vérifier votre environnement:
        turbotilt doctor
    EOS
  end
end
