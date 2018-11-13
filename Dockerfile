FROM surgecloud/kicad-base:latest
COPY drone-kicad /bin/
COPY kicad-ci-scripts /bin/ci-scripts
COPY PcbDraw /bin/PcbDraw
ENTRYPOINT ["/bin/drone-kicad"]
