#!/bin/bash
dnf install gcc pkg-config wayland-devel libX11-devel libxkbcommon-x11-devel mesa-libGLES-devel mesa-libEGL-devel libXcursor-devel xclip

# because spirv-cross is required to run the shader generators for gio and it's not available except on ubuntu 21 and
# later, we are keeping this compatible with ubuntu 20 by making it from scratch

git clone https://github.com/KhronosGroup/SPIRV-Cross.git
cd SPIRV-Cross
make -j$(nproc)
chmod +x spirv-cross
sudo mv spirv-cross /usr/local/bin/
cd ..
rm -rf SPIRV-Cross