make image
echo $DOCKER_PASSWORD | docker login docker.pkg.github.com -u $DOCKER_USERNAME --password-stdin
make push-image

eval "$(ssh-agent -s)"
chmod 600 deploy/deploy_rsa
ssh-add deploy/deploy_rsa

ssh $USER@$IP << EOF
    cd $DEPLOY_DIR
    docker-compose down
    docker-compose pull
    docker-compose up -d
EOF
