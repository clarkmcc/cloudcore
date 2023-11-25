FROM debian:bullseye-slim AS runtime

RUN apt-get update && apt-get install -y curl vim ca-certificates

COPY scripts/linux/install.sh /install.sh

RUN chmod +x /install.sh && /install.sh

CMD ["cloudcored", "version"]