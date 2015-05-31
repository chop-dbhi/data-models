FROM golang:onbuild

RUN make install
RUN make build

EXPOSE 8123

ENTRYPOINT ["app", "-port", "8123", "-host", "0.0.0.0"]

CMD ["-log", "error", "-path", "/opt/repos", "-repo", "https://github.com/chop-dbhi/data-models"]
