#!/usr/bin/env bash

echo "Run this script from the testdata folder cyclonedx/cyclonedx-go/testdata to"
echo "download valid and invalid JSON and XML files from CycloneDX/specification."
echo "This create the following folders: "
echo ""
echo "tree -d testdata/"
echo "  testdata/"
echo "  ├── 1.0"
echo "  ├── 1.1"
echo "  ├── 1.2"
echo "  ├── 1.3"
echo "  ├── 1.4"
echo "  ├── 1.5"
echo "  ├── 1.6"
echo "  ├── ext"
echo "  └── snapshots"
echo ""

# latest CycloneDX/specification release tag
CYCLONEDX_VERSION=1.6
TEST_DATA_FOLDER="../"

echo "remove temporary folder tmp/"
rm -r ./tmp

echo "create a temporary folder tmp/"
mkdir tmp

echo "cd into tmp/"
cd tmp

echo "get CycloneDX/specification v$CYCLONEDX_VERSION"
wget https://github.com/CycloneDX/specification/archive/refs/tags/$CYCLONEDX_VERSION.tar.gz

echo "untar $CYCLONEDX_VERSION.tar.gz"
tar -zxvf $CYCLONEDX_VERSION.tar.gz

echo "Copy content of specification-$CYCLONEDX_VERSION/tools/src/test/resources/ to local testdata folder"
cp -r specification-$CYCLONEDX_VERSION/tools/src/test/resources/* $TEST_DATA_FOLDER/

echo "go back to testdata folder"
cd $TEST_DATA_FOLDER

echo "remove temporary folder tmp/"
rm -r ./tmp