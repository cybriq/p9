#!/bin/bash
sudo apt install -y git wget curl build-essential gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev \
  libxkbcommon-x11-dev libgles2-mesa-dev \
  libegl1-mesa-dev libffi-dev libxcursor-dev xclip glslang-tools glslang-dev \
  spirv-tools spirv-headers \
  wine wine-development

# because spirv-cross is required to run the shader generators for gio and it's not available except on ubuntu 21 and
# later, we are keeping this compatible with ubuntu 20 by making it from scratch

git clone https://github.com/KhronosGroup/SPIRV-Cross.git
cd SPIRV-Cross
make -j$(nproc)
chmod +x spirv-cross
sudo mv spirv-cross /usr/local/bin/
cd ..
rm -rf SPIRV-Cross