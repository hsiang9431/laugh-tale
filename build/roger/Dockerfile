FROM alpine:latest

COPY roger /usr/bin/roger
RUN chmod +x /usr/bin/roger

COPY poneglyph /work/poneglyph
COPY poneglyph.sh /work/poneglyph.sh
RUN chmod +x /work/*

VOLUME ["/input", "/home/roger", "/output"]

ENTRYPOINT [ "roger" ]