<<<<<<< HEAD
FROM surgecloud/kicad-base:latest

MAINTAINER tomasz@napierala.org

COPY drone-kicad /bin/
=======
FROM toroid/kicad-base
ADD drone-kicad /bin/
>>>>>>> upstream/master
COPY kicad-ci-scripts /bin/ci-scripts
COPY PcbDraw /bin/PcbDraw
ENTRYPOINT ["/bin/drone-kicad"]
