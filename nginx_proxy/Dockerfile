FROM nginx

RUN apt-get update
RUN apt-get install --assume-yes python-dev python-setuptools python-pip
RUN apt-get install --assume-yes --reinstall ca-certificates
RUN pip install j2cli

ADD nginx.conf.j2 /
ADD nginx-site.conf.j2 /
ADD docker-entrypoint.sh /
RUN chmod a+rx /docker-entrypoint.sh
ENTRYPOINT ["/bin/bash", "/docker-entrypoint.sh"]
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
