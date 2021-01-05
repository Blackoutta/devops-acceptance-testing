docker build -t ecp-test .
docker run -e CONFIGFILE --add-host=devops.test.cq.iot.chinamobile.com:10.12.4.9 ecp-test /TESTRUN/grpctest