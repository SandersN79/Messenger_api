#!/bin/bash

echo "

============================================================================
    <###########-STARTING GENERATOR API INSTALLATION-##############>
============================================================================

"
echo "

<========###PROCESS 1 OF 3: STARTING GENERATOR API SETUP PROCESS###========>

"

### Installer Configuration Variables ###
export GitUser="git"
export GitPass="g!tpassW0rd123"
export ConfigUser="config"
export RepoName=`pwd`
export OSName=`awk -F= '/^NAME/{print $2}' /etc/os-release`
export APIUser=$USER
export DIRECTORY="/srv/git/default.git"

### Run updates according to the OS the API is running on ###
sudo apt-get -y update && sudo apt-get -y upgrade
sudo apt-get install -y git-all wget
sudo apt install -y build-essential

### Install and Configure Go v1.14.13 ###
cd ~ || exit
export HomeDir=`pwd`
wget https://dl.google.com/go/go1.14.13.linux-amd64.tar.gz
sudo tar -C $HomeDir -xzf go1.14.13.linux-amd64.tar.gz
export PATH=$PATH:$RepoName/bin
export PATH=$PATH:$RepoName/pkg
export PATH=$PATH:$RepoName/src
export GOPATH=$RepoName

# Configure the Go Variables for both ubuntu and centos as needed
if ! grep -Fxq "export PATH=$PATH:'$RepoName'/bin" $HomeDir/.profile ; then
  echo 'export PATH=$PATH:'$HOME'/go/bin' >> $HomeDir/.profile
  echo 'export PATH=$PATH:'$HOME'/go/pkg' >> $HomeDir/.profile
  echo 'export PATH=$PATH:'$HOME'/go/src' >> $HomeDir/.profile
  echo 'export PATH=$PATH:'$RepoName'/bin' >> $HomeDir/.profile
  echo 'export PATH=$PATH:'$RepoName'/pkg' >> $HomeDir/.profile
  echo 'export PATH=$PATH:'$RepoName'/src' >> $HomeDir/.profile
  echo 'export GOPATH='$RepoName >> $HomeDir/.profile
fi
source $HomeDir/.profile

### Get APIs Go Dependencies ###
cd $RepoName || exit
go get -u "github.com/gorilla/mux"
go get -u "github.com/gorilla/handlers"
go get -u "go.mongodb.org/mongo-driver/mongo"
go get -u "go.mongodb.org/mongo-driver/mongo/options"
go get -u "go.mongodb.org/mongo-driver/bson"
go get -u "go.mongodb.org/mongo-driver/bson/primitive"
go get -u "github.com/dgrijalva/jwt-go"
go get -u "github.com/gofrs/uuid"
go get -u "golang.org/x/crypto/bcrypt"
go get -u "gopkg.in/src-d/go-git.v4/..."
go get -u "github.com/JECSand/fetch"
rm -f ~/go1.14.13.linux-amd64.tar.gz

### Install and Configure Local Testing Mongo ###
# Run updates according to the OS the API is running on
if ! sudo grep -Fxq 'deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/4.0 multiverse' /etc/apt/sources.list.d/mongodb.list ; then
  sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 9DA31620334BD75D9DCB49F368818C72E52529D4
  echo 'deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/4.0 multiverse' | sudo tee /etc/apt/sources.list.d/mongodb.list
  sudo apt -y update && sudo apt -y install mongodb-org
  sudo systemctl enable mongod
  sudo systemctl start mongod
fi
echo "

<========###PROCESS 1 OF 3: GENERATOR API SETUP COMPLETE###========>

"

if ! [ -d "$DIRECTORY" ]; then

echo "

<========###PROCESS 2 OF 3: STARTING GIT SERVER INSTALLATION PROCESS###========>

"

## CREATE GIT USER ##
echo "
<--------##STEP 1 of 5: CREATING GIT USER##-------->

"

sudo adduser --home /home/git --shell /bin/bash git --gecos "First Last,RoomNumber,WorkPhone,HomePhone" --disabled-password
echo "$GitUser:$GitPass" | sudo chpasswd
sudo -u git mkdir /home/git/.ssh && sudo chmod 700 /home/git/.ssh
sudo -u git touch /home/git/.ssh/authorized_keys && sudo chmod 600 /home/git/.ssh/authorized_keys
sudo -u git ssh-keygen -o -t rsa -N "" -b 4096 -C emample@synercloud.io -f /home/git/.ssh/id_rsa
sudo cat /home/git/.ssh/id_rsa.pub | sudo tee /home/git/.ssh/authorized_keys
echo "
<--------##STEP 1 of 5: COMPLETE##-------->

"

## CONFIGURE GIT SERVICE ##
echo "
<--------##STEP 2 of 5: CONFIGURING GIT SERVICE##-------->

"

sudo mkdir /srv/git
cd /srv/git && sudo mkdir -p default.git
cd /srv/git/default.git && sudo git init --bare

if ! sudo grep -Fxq "$(which git-shell)" /etc/shells ; then
  which git-shell | sudo tee -a /etc/shells
fi

sudo chsh git -s "$(which git-shell)"

echo "
no-port-forwarding,no-X11-forwarding,no-agent-forwarding,no-pty" | sudo tee -a /home/git/.ssh/authorized_keys

echo "[Unit]
Description=Start Git Daemon

[Service]
ExecStart=/usr/bin/git daemon --reuseaddr --base-path=/srv/git/
Restart=always
RestartSec=500ms

StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=git-daemon
User=git
Group=git

[Install]
WantedBy=multi-user.target" | sudo tee /etc/systemd/system/git-daemon.service

sudo systemctl enable git-daemon
sudo systemctl start git-daemon
sudo systemctl stop git-daemon

cd /srv/git/default.git && sudo touch git-daemon-export-ok
echo "
<--------##STEP 2 of 5: COMPLETE##-------->

"

## SETUP LOCAL APACHE FOR GIT SERVICE ##
echo "
<--------##STEP 3 of 5: CONFIGURING LOCAL APACHE##-------->

"

sudo apt-get -y install apache2 apache2-utils
sudo a2enmod cgi alias env
sudo chgrp -R www-data /srv/git

if ! sudo grep -Fxq "SetEnv GIT_PROJECT_ROOT /srv/git" /etc/apache2/apache2.conf ; then
echo "

<Files \"git-http-backend\">
    AuthType Basic
    AuthName \"Git Access\"
    AuthUserFile /srv/git/.htpasswd
    Require expr !(%{QUERY_STRING} -strmatch '*service=git-receive-pack*' || %{REQUEST_URI} =~ m#/git-receive-pack\$#)
    Require valid-user
</Files>

SetEnv GIT_PROJECT_ROOT /srv/git
SetEnv GIT_HTTP_EXPORT_ALL
ScriptAlias /git/ /usr/lib/git-core/git-http-backend/

" | sudo tee -a /etc/apache2/apache2.conf
fi

sudo systemctl restart apache2
sudo htpasswd -b -c /srv/git/.htpasswd $GitUser $GitPass
sudo usermod -a -G www-data git
sudo usermod -s /usr/bin/git-shell git
sudo chown git:www-data /srv/git/ -R
sudo chmod 777 /srv/git/ -R
echo "
<--------##STEP 3 of 5: COMPLETE##-------->

"

## SETUP GIT SHELL COMMANDS ##
echo "
<--------##STEP 4 of 5: SETTING UP GIT SHELL COMMANDS##-------->

"
sudo cp /usr/share/doc/git/contrib/git-shell-commands /home/git -R
sudo chown git:git /home/git/git-shell-commands/ -R
sudo chmod +x /home/git/git-shell-commands/help
sudo chmod +x /home/git/git-shell-commands/list
echo "
<--------##STEP 4 of 5: COMPLETE##-------->

"

## SETUP default.git REPO ##
echo "
<--------##STEP 5 of 5: SETTING UP DEFAULT.GIT REPO##-------->

"
mkdir "$HomeDir"/git_config
cd "$HomeDir"/git_config && git clone http://127.0.1.1/git/default.git
cd "$HomeDir"/git_config/default && touch README.md
echo "#Default Repo From Which All Other Repos Will Descend" >> "$HomeDir"/git_config/default/README.md
cd "$HomeDir"/git_config/default && git config --global user.name "Git Installer"
cd "$HomeDir"/git_config/default && git config --global user.email "config@synercloud.io"
cd "$HomeDir"/git_config/default && git add .
cd "$HomeDir"/git_config/default && git commit -m "Initial Commit"
cd "$HomeDir"/git_config/default && git push http://$GitUser:$GitPass@127.0.1.1/git/default.git --all
rm -rf "$HomeDir"/git_config
cd $RepoName || exit
echo "
<--------##STEP 5 of 5: COMPLETE##-------->

"
echo "

<========###PROCESS 2 OF 3: GIT SERVER INSTALLATION COMPLETE###========>

"
else
   echo "

<========###PROCESS 2 OF 3: EXISTING GIT SERVER DETECTED.... SKIPPING PROCESS...###========>

"
fi
echo "

<========###PROCESS 3 OF 3: STARTING SYSTEMD SETUP PROCESS###========>

"

### Setup systemd service ###
sudo sh $RepoName/install/setup_service_ubuntu.sh $RepoName $APIUser
echo "

<========###PROCESS 3 OF 3: SYSTEMD SETUP COMPLETE###========>

"
echo "

============================================================================
    <###########-GENERATOR API INSTALLATION COMPLETE-##############>
============================================================================

"
