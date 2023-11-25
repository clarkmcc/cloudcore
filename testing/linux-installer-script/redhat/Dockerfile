FROM redhat/ubi8 AS runtime

COPY scripts/linux/install.sh /install.sh

RUN chmod +x /install.sh && /install.sh

CMD ["cloudcored", "version"]