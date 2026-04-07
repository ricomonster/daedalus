# daedalus
My utility toolbox

## Installation
 
**1. Clone and enter the repo**
```bash
git clone <this-repo>
cd release-cli
```
 
**2. Build and install**
 
Pick whichever option fits your setup:
 
```bash
# Recommended — installs to ~/.local/bin (no sudo required)
make install
 
# If you use Homebrew (no sudo, already in PATH)
make install INSTALL_DIR=$(brew --prefix)/bin
 
# If you have Go installed (~/go/bin is usually already in PATH)
make install INSTALL_DIR=$(go env GOPATH)/bin
 
# Custom directory
make install INSTALL_DIR=~/bin
```
 
> If the install directory isn't in your `$PATH` yet, the installer will print the exact line to add to your `~/.zshrc` or `~/.bashrc`.
 
**3. Verify**
```bash
release-cli --version
```
