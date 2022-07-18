sed -i '' 's/http:\/\/localhost\:8080/https:\/\/filetransfer\.safespace\.dev/g' .env
rsync -ayp ./ tau@116.203.84.88:~/filetransfer
sed -i '' 's/https:\/\/filetransfer\.safespace\.dev/http:\/\/localhost\:8080/g' .env
