#!/bin/bash

if [ "$SHELL" == "/bin/zsh" ]; then
	echo "Detected zsh"
	sh=1
elif [ "$SHELL" == "/bin/bash" ]; then
	echo "Detected bash"
	sh=2
else
	echo "shell not detected"
	sh=0
fi

go get -u -v
go build -o whsub 
if [ $? -ne 0 ]; then
	exit 1
fi

if test -f "./whsub"; then
	echo "Binary './whsub' built successfully"
else
	echo "Building binary failed!"
	exit 1
fi

# Install autocompletion

mkdir -p $HOME/.config/autocompletion/
if [ $sh -eq 1 ]; then
	./whsub --completion-script-zsh > $HOME/.config/autocompletion/whsub.compl
	if [ -z "$(cat $HOME/.zshrc | grep 'source $HOME/.config/autocompletion/whsub.compl')" ]; then
		echo "source $HOME/.config/autocompletion/whsub.compl" >> $HOME/.zshrc	
	fi
	echo "autocompletion installed successfully"
elif [ $sh -eq 2 ]; then
	./whsub --completion-script-bash > $HOME/.config/autocompletion/whsub.compl
	if [ -z "$(cat $HOME/.bashrc | grep 'source $HOME/.config/autocompletion/whsub.compl')" ]; then
		echo "source $HOME/.config/autocompletion/whsub.compl" >> $HOME/.bashrc
	fi
	echo "autocompletion installed successfully"
else 
	echo "No autocompletion installed. Bash '$SHELL' not supported"
fi

#Install binary

sudo cp whsub /usr/local/bin/whsub &&
sudo chmod +x /usr/local/bin/whsub &&
sudo touch /usr/local/man/man1/whsub.1 &&
sudo chown jojii:jojii /usr/local/man/man1/whsub.1 &&
sudo /usr/local/bin/whsub --help-man > /usr/local/man/man1/whsub.1 &&
sudo chown root:root /usr/local/man/man1/whsub.1 &&
echo "Man entry installed"
