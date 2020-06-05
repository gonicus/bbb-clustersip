# What is it?

astregbackend is a small http service which can be used as realtime backend for Asterisk.

It generates PJSIP-Endpoints (endpoint, aor, auth) within a specific range of numbers when Asterisk asks for it.
It takes a parameter 'digits'.
Whith digits=3, it generates configuration for Endpoints 000-999.
Whith digits=5, it generates configuration for Endpoints 00000-99999.

After generating the configuration it is stored in a redis.
This is neccessary to serve Asterisk a list of "all" generated endpoints.
When executing 'pjsip show endpoints' for example.

# Installation

```
# install golang on ubuntu 16.04
apt install golang-1.10 redis-server

# get and build astregbackend
/usr/lib/go-1.10/bin/go get github.com/denzs/astregbackend

cd ~/go/src/github.com/denzs/astregbackend
cp ~/go/bin/astregbackend /usr/local/sbin/
cp astregbackend.sample /etc/default/astregbackend
cp astregbackend.service /etc/systemd/system/
cp astregbackend.conf /etc
systemctl daemon-reload
systemctl enable astregbackend
systemctl start astregbackend
```

See extconfig.conf and sorcery.conf for examples on how to use with Asterisk.

# Operating

* ensure that the RedisExpiration is bigger than your SIP registration expirationtime!
* ensure that astregbackend parameter *digits* matches your BBB setting *defaultNumDigitsForTelVoice*
