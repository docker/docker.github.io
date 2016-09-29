FROM bagel/universe:337f873f4f23f4b2603972229ae3519c5f61f6d7

ENV ENV local
ENV NODE_ENV local

COPY ./app /opt/hub/app
COPY ./Makefile /opt/hub/Makefile
COPY ./_webpack /opt/hub/_webpack
COPY ./gulpfile.js /opt/hub/gulpfile.js
COPY ./gulp-tasks /opt/hub/gulp-tasks
COPY ./app-server /opt/hub/app-server
COPY ./.eslintrc /opt/hub/.eslintrc

RUN make server-prod-target
RUN make server-extras
RUN make js-local
RUN make images-prod
RUN make docker-font-prod
RUN gulp images::prod
RUN make styles-base-prod
RUN make stats-dir
RUN make css-stats
