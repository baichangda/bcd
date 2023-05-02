rsync -e 'ssh -p 60022' -avz --exclude={.idea,.gradle,logs,*.log,temp} /Users/baichangda/bcd/bcd/bcd_react/build/* root@cdbai.cn:/app/software/nginx-1.24.0/html
