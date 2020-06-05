# Setup

Should work with any Asterisk >= 13.0 with PJSIP.

* install Asterisk like ```apt install asterisk```
* copy configs from this folder to /etc/asterisk
* adapt pjsip.conf to your needs
  * adapt transport section (see line 123 and below)
  * adapt acl to your needs (see line 414 and below)
  * adapt upstream endpoint and registration to your needs (see line 944 and below)

# Dialplan

* if your PBX does set a reasonable name for the caller you might want to comment the line with CALLERID(name) in extensions.conf
* if you want to set another language than english adapt and uncomment the line with LANGUAGE in extensions.conf
  * make sure you have the correct prompts in the correct folder
