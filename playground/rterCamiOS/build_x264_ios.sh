#!/bin/sh

# This script is originally based off of the one by Gabriel Handford
# Original scripts can be found here: https://github.com/gabriel/ffmpeg-iphone-build
# Modified by Roderick Buenviaje
# Builds versions of the VideoLAN x264 for armv6 and armv7
# Combines the two libraries into a single one

trap exit ERR

export DIR=./x264
export DEST6=armv6/
export DEST7=armv7/
export DESTi386=i386/
#specify the version of iOS you want to build against
export VERSION="6.1"

#mkdir -p ./x264

#git clone git://git.videolan.org/x264.git x264

cd $DIR

# echo "Building armv6..."

# CC=/Developer/Platforms/iPhoneOS.platform/Developer/usr/bin/llvm-gcc \
# ./configure \
# --host=arm-apple-darwin \
# --sysroot=/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS${VERSION}.sdk \
# --prefix=$DEST6 \
# --extra-cflags='-arch armv6' \
# --extra-ldflags='-L/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS5.0.sdk/usr/lib/system -arch armv6' \
# --enable-pic \
# # --disable-asm \
# --enable-static

# make && make install

# echo "Installed: $DEST6"


echo "Building armv7..."

CC=/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/usr/bin/llvm-gcc \
./configure \
--host=arm-apple-darwin \
--sysroot=/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS${VERSION}.sdk \
--prefix=$DEST7 \
--extra-cflags='-arch armv7' \
--extra-ldflags='-arch armv7 -L/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS6.1.sdk/usr/lib/system  ' \
--enable-pic \
--enable-static \
--enable-shared

# make clean
# make && make install

# echo "Installed: $DEST7"

# echo "Building i386..."

# CC=/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/usr/bin/gcc \
# ./configure \
# --host=i386-apple-darwin \
# --sysroot=/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/SDKs/iPhoneSimulator6.1.sdk \
# --prefix=$DESTi386 \
# --extra-cflags="" \
# #'-arch i386' \
# --extra-ldflags="" \
# #'-arch i386 -L/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/SDKs/iPhoneSimulator6.1.sdk' \
# --enable-static

# make clean
# make && make install

# echo "Installed: $DESTi386"


# echo "Combining Libraries"
# ARCHS="armv6 armv7"

# BUILD_LIBS="libx264.a"

# OUTPUT_DIR="x264-uarch"
# mkdir $OUTPUT_DIR
# mkdir $OUTPUT_DIR/lib
# mkdir $OUTPUT_DIR/include

# for LIB in $BUILD_LIBS; do
#   LIPO_CREATE=""
#   for ARCH in $ARCHS; do
#     LIPO_CREATE="$LIPO_CREATE-arch $ARCH $ARCH/lib/$LIB "
#   done
#   OUTPUT="$OUTPUT_DIR/lib/$LIB"
#   echo "Creating: $OUTPUT"
#   lipo -create $LIPO_CREATE -output $OUTPUT
#   lipo -info $OUTPUT
# done

# cp -f $ARCH/include/*.* $OUTPUT_DIR/include/