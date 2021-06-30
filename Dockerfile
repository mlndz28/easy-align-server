FROM golang:1.15-buster

RUN apt-get update &&  apt-get install -y \
    git \
    wget \
    praat \
    meson \
    gcc-multilib \
    && rm -rf /var/lib/apt/lists/*


WORKDIR /tmp
RUN git clone https://github.com/open-speech/HTK.git \
    && cd HTK \
    && ./configure --disable-hslab \
    && make htklib \
    && make htktools \
    && make install-htktools \
    && rm -rf /tmp/HTK

RUN git clone https://github.com/TALP-UPC/saga.git \
    && cd saga \
    && meson builddir \
    && cd builddir \
    && ninja test \
    && ninja install

RUN wget -qO- https://raw.githubusercontent.com/mlndz28/praat-easy-align-linux/master/installer.sh | bash

COPY main.go /go/src/easyalignserver/main.go

RUN go get github.com/mlndz28/praatgo \
    && export GOBIN=/go/bin \
    && go install easyalignserver

CMD easyalignserver