#!/bin/bash
#
# Script to generate APT repo metadata
mkdir -p dist/apt
cp *.deb dist/apt/
cd dist/apt
dpkg-scanpackages . /dev/null | gzip -9c > Packages.gz
