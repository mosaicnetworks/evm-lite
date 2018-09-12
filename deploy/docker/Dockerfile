FROM alpine
COPY entry.sh /entry.sh
COPY keys/client.pub /root/.ssh/authorized_keys
RUN apk add --no-cache openssh \
    && ssh-keygen -A \
    && echo "PermitRootLogin yes" >> /etc/ssh/sshd_config 
ADD evm-lite /usr/local/bin
EXPOSE 22 8080 8000
ENV HOME=/
ENTRYPOINT ["/entry.sh"]
CMD []