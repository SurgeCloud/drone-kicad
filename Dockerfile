FROM surgecloud/kicad-base:latest

MAINTAINER tomasz@napierala.org

COPY drone-kicad /bin/
COPY kicad-ci-scripts /bin/ci-scripts
COPY PcbDraw /bin/PcbDraw
ENTRYPOINT ["/bin/drone-kicad"]
