class Turbotilt < Formula
  desc "CLI pour environnements dev cloud-native"
  homepage "https://github.com/Fascinax/turbotilt"
  version "0.1.0"  # Change this to your current version
  license "MIT"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/Fascinax/turbotilt/releases/download/v#{version}/turbotilt-#{version}-darwin-arm64.zip"
      sha256 "24F3972EE9B564D7361BFB0FDCC8309488667F834249D06FE921D084C2281557"  # darwin-arm64
    else
      url "https://github.com/Fascinax/turbotilt/releases/download/v#{version}/turbotilt-#{version}-darwin-amd64.zip"
      sha256 "94E9F0C43BDB7E473A1A53B0876D186267EAFAA64D2B84F2E68CB59E46E9E2EA"  # darwin-amd64
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/Fascinax/turbotilt/releases/download/v#{version}/turbotilt-#{version}-linux-arm64.zip"
      sha256 "8D3C0E96087E4FA93578E7B724281FF3CBAE93FBE4D7A51B998A1B854D389735"  # linux-arm64
    else
      url "https://github.com/Fascinax/turbotilt/releases/download/v#{version}/turbotilt-#{version}-linux-amd64.zip"
      sha256 "6BAD512411BF018828FA13DBC98947B8340C61FDB0D3E9D61266ADE9DBDE8F7B"  # linux-amd64
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
