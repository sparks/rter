## Deps

    sudo apt-get install subversion git cmake build-essential yasm libqt4-dev kdelibs5-dev libsdl1.2-dev libsdl-image1.2-dev libxml2-dev libx264-dev libtheora-dev libxvidcore-dev libogg-dev libvorbis-dev libschroedinger-dev libmp3lame-dev libfaac-dev libfaad-dev libgsm1-dev libopencore-amrnb-dev libopencore-amrwb-dev libsamplerate0-dev libjack-dev libsox-dev ladspa-sdk swh-plugins libmad0-dev libpango1.0-dev

    apt-get install subversion unzip frei0r-plugins-dev libdc1394-22-dev libfaac-dev libmp3lame-dev libx264-dev libdirac-dev libxvidcore4-dev libfreetype6-dev libvorbis-dev libgsm1-dev libopencore-amrnb-dev libopencore-amrwb-dev libopenjpeg-dev librtmp-dev libschroedinger-dev libspeex-dev libtheora-dev libva-dev libvpx-dev libvo-amrwbenc-dev libvo-aacenc-dev libaacplus-dev libbz2-dev libgnutls-dev libssl-dev libopenal-dev libv4l-dev libpulse-dev libmodplug-dev libass-dev libcdio-dev libcdio-cdda-dev libcdio-paranoia-dev libvdpau-dev libxfixes-dev libxext-dev libbluray-dev

## Config

debian recomended
    ./configure --enable-avresample --enable-vda --enable-libx264 --enable-libfaac --enable-libmp3lame --enable-libxvid --enable-libtheora --enable-libvorbis --enable-libvpx

rter
    ./configure --enable-gpl --enable-libmp3lame --enable-libvorbis --enable-libfaac --enable-libopencore-amrnb --enable-libvpx --enable-libopencore-amrwb --enable-libtheora --enable-libx264 --enable-nonfree --enable-version3 --enable-avresample --enable-vda --enable-libxvid

brew 
    --prefix=/usr/local/Cellar/ffmpeg/1.2 --enable-shared --enable-pthreads --enable-gpl --enable-version3 --enable-nonfree --enable-hardcoded-tables --enable-avresample --enable-vda --cc=cc --host-cflags= --host-ldflags= --enable-libx264 --enable-libfaac --enable-libmp3lame --enable-libxvid --enable-libtheora --enable-libvorbis --enable-libvpx
