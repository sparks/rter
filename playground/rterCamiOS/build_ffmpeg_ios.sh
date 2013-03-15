# the purpose of this scrip is to build ffmpeg for iOS and the iOS simulator in Xcode
#
# compiles and seems to work with sdk 6.1, Xcode 4.6, iOS 6.0.x, osx 8.2.x
#
# comment out parts which are not needed
#
# put together using:
# * a github repo I can't find right now
# * http://www.tangentsoftworks.com/blog/2012/11/12/how-to-prepare-your-mac-for-ios-development-with-ffmpeg-libraries/
# * http://blog.fh-kaernten.at/wehr/?p=1764
# * https://github.com/rodisbored/ffmpeg_x264_iOS5_build

#change the below line to point to the where libx264 is
export X264ROOT=../x264/armv7
export X264LIB=$X264ROOT/lib
export X264INCLUDE=$X264ROOT/include

rm -r ./compiled
 
echo Configure for armv7 build
./configure \
--cc=/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/usr/bin/gcc \
--as='gas-preprocessor.pl /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/usr/bin/gcc' \
--nm="/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/usr/bin/nm" \
--sysroot=/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS6.1.sdk \
--target-os=darwin \
--arch=arm \
--cpu=cortex-a8 \
--enable-pic \
--extra-cflags='-I$X264INCLUDE -arch armv7 -mfpu=neon -pipe -Os -gdwarf-2 -miphoneos-version-min=5.0' \
--extra-ldflags='-L$X264LIB -arch armv7 -miphoneos-version-min=5.0 -isysroot /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS6.1.sdk' \
--prefix=compiled/armv7 \
--enable-cross-compile \
--enable-nonfree \
--enable-gpl \
--disable-doc \
--disable-ffmpeg \
--disable-ffplay \
--disable-ffprobe \
--disable-ffserver \
--enable-libx264 \
--enable-encoder=libx264
#--disable-armv5te \
#--disable-swscale-alpha \

#--disable-asm \	# we want assembly level optimization
# --disable-debug \
# libx264

 
# make clean
# make && make install
 
# echo Configure for i386
# ./configure \
# --cc=/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/usr/bin/gcc \
# --as='gas-preprocessor.pl /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/usr/bin/gcc' \
# --nm="/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/usr/bin/nm" \
# --sysroot=/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/SDKs/iPhoneSimulator6.1.sdk \
# --target-os=darwin \
# --arch=i386 \
# --cpu=i386 \
# --extra-cflags='-arch i386 -miphoneos-version-min=5.0 -mdynamic-no-pic' \
# --extra-ldflags='-arch i386 -miphoneos-version-min=5.0 -isysroot /Applications/Xcode.app/Contents/Developer/Platforms/iPhoneSimulator.platform/Developer/SDKs/iPhoneSimulator6.1.sdk' \
# --prefix=compiled/i386 \
# --enable-cross-compile \
# --enable-nonfree \
# --enable-gpl \
# --disable-armv5te \
# --disable-swscale-alpha \
# --disable-doc \
# --disable-ffmpeg \
# --disable-ffplay \
# --disable-ffprobe \
# --disable-ffserver \
# --disable-asm \
# --disable-debug
 
# make clean
# make && make install
 
 
# # # make fat (universal) libs
# mkdir -p ./compiled/fat/lib
 
# lipo -output ./compiled/fat/lib/libavcodec.a  -create \
# -arch armv7 ./compiled/armv7/lib/libavcodec.a \
# -arch i386 ./compiled/i386/lib/libavcodec.a
 
# lipo -output ./compiled/fat/lib/libavdevice.a  -create \
# -arch armv7 ./compiled/armv7/lib/libavdevice.a \
# -arch i386 ./compiled/i386/lib/libavdevice.a
 
# lipo -output ./compiled/fat/lib/libavformat.a  -create \
# -arch armv7 ./compiled/armv7/lib/libavformat.a \
# -arch i386 ./compiled/i386/lib/libavformat.a
 
# lipo -output ./compiled/fat/lib/libavutil.a  -create \
# -arch armv7 ./compiled/armv7/lib/libavutil.a \
# -arch i386 ./compiled/i386/lib/libavutil.a
 
# lipo -output ./compiled/fat/lib/libswresample.a  -create \
# -arch armv7 ./compiled/armv7/lib/libswresample.a \
# -arch i386 ./compiled/i386/lib/libswresample.a
 
# lipo -output ./compiled/fat/lib/libpostproc.a  -create \
# -arch armv7 ./compiled/armv7/lib/libpostproc.a \
# -arch i386 ./compiled/i386/lib/libpostproc.a
 
# lipo -output ./compiled/fat/lib/libswscale.a  -create \
# -arch armv7 ./compiled/armv7/lib/libswscale.a \
# -arch i386 ./compiled/i386/lib/libswscale.a
 
# lipo -output ./compiled/fat/lib/libavfilter.a  -create \
# -arch armv7 ./compiled/armv7/lib/libavfilter.a \
# -arch i386 ./compiled/i386/lib/libavfilter.a