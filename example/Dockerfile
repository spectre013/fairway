#######################################
# Release
#######################################
FROM tb-common-java as release
#run apk update && apk add supervisor && apk add bash
# Copy our static executable
COPY .skel /
COPY .skel /skel
COPY site /data/fairway/site

COPY example /data/fairway/bin/static

RUN chmod u+xw /data/fairway/run.sh
CMD [ "supervisord", "-n" ]
